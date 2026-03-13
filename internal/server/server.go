package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pinchtab/pinchtab/internal/cli"
	"github.com/pinchtab/pinchtab/internal/config"
	"github.com/pinchtab/pinchtab/internal/dashboard"
	"github.com/pinchtab/pinchtab/internal/handlers"
	"github.com/pinchtab/pinchtab/internal/orchestrator"
	"github.com/pinchtab/pinchtab/internal/profiles"
	"github.com/pinchtab/pinchtab/internal/scheduler"
	"github.com/pinchtab/pinchtab/internal/strategy"
	"github.com/pinchtab/pinchtab/internal/web"

	// Register strategies
	_ "github.com/pinchtab/pinchtab/internal/strategy/alwayson"
	_ "github.com/pinchtab/pinchtab/internal/strategy/autorestart"
	_ "github.com/pinchtab/pinchtab/internal/strategy/explicit"
	_ "github.com/pinchtab/pinchtab/internal/strategy/simple"
)

func RunDashboard(cfg *config.RuntimeConfig, version string) {
	dashPort := cfg.Port
	if dashPort == "" {
		dashPort = "9870"
	}
	startedAt := time.Now()

	profilesDir := cfg.ProfilesBaseDir
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		slog.Error("cannot create profiles dir", "err", err)
		os.Exit(1)
	}

	profMgr := profiles.NewProfileManager(profilesDir)
	dash := dashboard.NewDashboard(nil)
	orch := orchestrator.NewOrchestrator(profilesDir)
	orch.ApplyRuntimeConfig(cfg)
	orch.SetProfileManager(profMgr)
	dash.SetInstanceLister(orch)
	dash.SetMonitoringSource(orch)
	dash.SetServerMetricsProvider(func() dashboard.MonitoringServerMetrics {
		snapshot := handlers.SnapshotMetrics()
		return dashboard.MonitoringServerMetrics{
			GoHeapAllocMB:   MetricFloat(snapshot["goHeapAllocMB"]),
			GoNumGoroutine:  MetricInt(snapshot["goNumGoroutine"]),
			RateBucketHosts: MetricInt(snapshot["rateBucketHosts"]),
		}
	})
	configAPI := dashboard.NewConfigAPI(cfg, orch, profMgr, orch, version, startedAt)

	// Wire up instance events to SSE broadcast
	orch.OnEvent(func(evt orchestrator.InstanceEvent) {
		dash.BroadcastSystemEvent(dashboard.SystemEvent{
			Type:     evt.Type,
			Instance: evt.Instance,
		})
	})

	mux := http.NewServeMux()

	dash.RegisterHandlers(mux)
	configAPI.RegisterHandlers(mux)
	profMgr.RegisterHandlers(mux)

	var activeStrategy strategy.Strategy
	stratName := "manual"
	if cfg.Strategy != "" {
		strat, err := strategy.New(cfg.Strategy)
		if err != nil {
			slog.Warn("unknown strategy, falling back to default", "strategy", cfg.Strategy, "err", err)
		} else {
			if runtimeAware, ok := strat.(strategy.RuntimeConfigAware); ok {
				runtimeAware.SetRuntimeConfig(cfg)
			}
			if setter, ok := strat.(strategy.OrchestratorAware); ok {
				setter.SetOrchestrator(orch)
			}
			strat.RegisterRoutes(mux)
			activeStrategy = strat
			stratName = strat.Name()
		}
	}

	allocPolicy := cfg.AllocationPolicy
	if allocPolicy == "" {
		allocPolicy = "none"
	}

	listenStatus := "starting"
	if cli.IsDaemonRunning() && CheckPinchTabRunning(dashPort, cfg.Token) {
		listenStatus = "running"
	}

	cli.PrintStartupBanner(cfg, cli.StartupBannerOptions{
		Mode:         "server",
		ListenAddr:   cfg.Bind + ":" + dashPort,
		ListenStatus: listenStatus,
		PublicURL:    fmt.Sprintf("http://localhost:%s", dashPort),
		Strategy:     stratName,
		Allocation:   allocPolicy,
	})

	if listenStatus == "running" {
		fmt.Println(cli.StyleStdout(cli.WarningStyle, fmt.Sprintf("  pinchtab already running as a daemon on port %s", dashPort)))
		fmt.Println(cli.StyleStdout(cli.MutedStyle, "  Stop the daemon first with `pinchtab daemon stop` to run in the foreground."))
		fmt.Println()
		os.Exit(0)
	}

	slog.Info("orchestration", "strategy", stratName, "allocation", allocPolicy)

	if activeStrategy == nil {
		orch.RegisterHandlers(mux)
		RegisterDefaultProxyRoutes(mux, orch)
	}

	var sched *scheduler.Scheduler
	if cfg.Scheduler.Enabled {
		schedCfg := scheduler.DefaultConfig()
		schedCfg.Enabled = true
		if cfg.Scheduler.Strategy != "" {
			schedCfg.Strategy = cfg.Scheduler.Strategy
		}
		if cfg.Scheduler.MaxQueueSize > 0 {
			schedCfg.MaxQueueSize = cfg.Scheduler.MaxQueueSize
		}
		if cfg.Scheduler.MaxPerAgent > 0 {
			schedCfg.MaxPerAgent = cfg.Scheduler.MaxPerAgent
		}
		if cfg.Scheduler.MaxInflight > 0 {
			schedCfg.MaxInflight = cfg.Scheduler.MaxInflight
		}
		if cfg.Scheduler.MaxPerAgentFlight > 0 {
			schedCfg.MaxPerAgentFlight = cfg.Scheduler.MaxPerAgentFlight
		}
		if cfg.Scheduler.ResultTTLSec > 0 {
			schedCfg.ResultTTL = time.Duration(cfg.Scheduler.ResultTTLSec) * time.Second
		}
		if cfg.Scheduler.WorkerCount > 0 {
			schedCfg.WorkerCount = cfg.Scheduler.WorkerCount
		}

		resolver := &scheduler.ManagerResolver{Mgr: orch.InstanceManager()}
		sched = scheduler.New(schedCfg, resolver)
		sched.RegisterHandlers(mux)
		sched.Start()
		slog.Info("scheduler enabled", "strategy", schedCfg.Strategy, "workers", schedCfg.WorkerCount)
	}

	mux.HandleFunc("GET /health", configAPI.HandleHealth)
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		web.JSON(w, 200, map[string]any{"metrics": handlers.SnapshotMetrics()})
	})

	handler := handlers.LoggingMiddleware(handlers.CorsMiddleware(handlers.AuthMiddleware(cfg, mux)))
	cli.LogSecurityWarnings(cfg)

	srv := &http.Server{
		Addr:              cfg.Bind + ":" + dashPort,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	strategyHandlesLaunch := false
	if activeStrategy != nil {
		if err := activeStrategy.Start(context.Background()); err != nil {
			slog.Error("strategy start failed", "strategy", activeStrategy.Name(), "err", err)
		} else if launchAware, ok := activeStrategy.(strategy.LaunchAware); ok && launchAware.HandlesLaunch() {
			strategyHandlesLaunch = true
		}
	}

	autoLaunch := strings.EqualFold(os.Getenv("PINCHTAB_AUTO_LAUNCH"), "1") ||
		strings.EqualFold(os.Getenv("PINCHTAB_AUTO_LAUNCH"), "true") ||
		strings.EqualFold(os.Getenv("PINCHTAB_AUTO_LAUNCH"), "yes")
	if autoLaunch && !strategyHandlesLaunch {
		defaultProfile := os.Getenv("PINCHTAB_DEFAULT_PROFILE")
		defaultProfileExplicit := defaultProfile != ""
		defaultPort := os.Getenv("PINCHTAB_DEFAULT_PORT")

		go func() {
			time.Sleep(500 * time.Millisecond)
			profileToLaunch := defaultProfile
			if !defaultProfileExplicit {
				list, err := profMgr.List()
				if err != nil {
					slog.Warn("auto-launch profile list failed", "err", err)
				}
				if len(list) > 0 {
					profileToLaunch = list[0].Name
				} else {
					profileToLaunch = "default"
					if err := os.MkdirAll(filepath.Join(profilesDir, profileToLaunch, "Default"), 0755); err != nil {
						slog.Warn("failed to create auto-launch profile dir", "profile", profileToLaunch, "err", err)
					}
				}
			}

			headlessDefault := os.Getenv("PINCHTAB_HEADED") == ""
			inst, err := orch.Launch(profileToLaunch, defaultPort, headlessDefault, nil)
			if err != nil {
				slog.Warn("auto-launch failed", "profile", profileToLaunch, "err", err)
				return
			}
			slog.Info("auto-launched instance", "profile", profileToLaunch, "id", inst.ID, "port", inst.Port, "headless", headlessDefault)
		}()
	}

	shutdownOnce := &sync.Once{}
	doShutdown := func() {
		shutdownOnce.Do(func() {
			slog.Info("shutting down dashboard...")
			if activeStrategy != nil {
				if err := activeStrategy.Stop(); err != nil {
					slog.Warn("strategy stop failed", "err", err)
				}
			}
			if sched != nil {
				sched.Stop()
			}
			dash.Shutdown()
			orch.Shutdown()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				slog.Error("shutdown http", "err", err)
			}
		})
	}

	mux.HandleFunc("POST /shutdown", func(w http.ResponseWriter, r *http.Request) {
		web.JSON(w, 200, map[string]string{"status": "shutting down"})
		go doShutdown()
	})

	go func() {
		sig := make(chan os.Signal, 2)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		go doShutdown()
		<-sig
		slog.Warn("force shutdown requested")
		orch.ForceShutdown()
		os.Exit(130)
	}()

	slog.Info("dashboard started", "port", dashPort)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server", "err", err)
		os.Exit(1)
	}
}
