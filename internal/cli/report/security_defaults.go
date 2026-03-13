package report

import (
	"reflect"
	"strings"

	"github.com/pinchtab/pinchtab/internal/config"
)

func ApplyRecommendedSecurityDefaults(fc *config.FileConfig) {
	defaults := config.DefaultFileConfig()
	if fc == nil {
		return
	}
	fc.Server.Bind = defaults.Server.Bind
	fc.Security = defaults.Security
	if strings.TrimSpace(fc.Server.Token) == "" {
		token, err := config.GenerateAuthToken()
		if err == nil {
			fc.Server.Token = token
		}
	}
}

func applyRecommendedSecurityDefaults(fc *config.FileConfig) {
	ApplyRecommendedSecurityDefaults(fc)
}

func RestoreSecurityDefaults() (string, bool, error) {
	fc, configPath, err := config.LoadFileConfig()
	if err != nil {
		return "", false, err
	}
	before := securityDefaultsSnapshot(fc)
	ApplyRecommendedSecurityDefaults(fc)
	after := securityDefaultsSnapshot(fc)
	if reflect.DeepEqual(before, after) {
		return configPath, false, nil
	}
	if err := config.SaveFileConfig(fc, configPath); err != nil {
		return "", false, err
	}
	return configPath, true, nil
}

func restoreSecurityDefaults() (string, bool, error) {
	return RestoreSecurityDefaults()
}

func RecommendedSecurityDefaultLines(cfg *config.RuntimeConfig) []string {
	posture := AssessSecurityPosture(cfg)
	ordered := []string{
		"server.bind = 127.0.0.1",
		"security.allowEvaluate = false",
		"security.allowMacro = false",
		"security.allowScreencast = false",
		"security.allowDownload = false",
		"security.allowUpload = false",
		"security.attach.enabled = false",
		"security.attach.allowHosts = 127.0.0.1,localhost,::1",
		"security.attach.allowSchemes = ws,wss",
		"security.idpi.enabled = true",
		"security.idpi.allowedDomains = 127.0.0.1,localhost,::1",
		"security.idpi.strictMode = true",
		"security.idpi.scanContent = true",
		"security.idpi.wrapContent = true",
	}
	needed := make(map[string]bool, len(ordered))

	for _, check := range posture.Checks {
		if check.Passed {
			continue
		}
		switch check.ID {
		case "bind_loopback":
			needed["server.bind = 127.0.0.1"] = true
		case "sensitive_endpoints_disabled":
			for _, line := range []string{
				"security.allowEvaluate = false",
				"security.allowMacro = false",
				"security.allowScreencast = false",
				"security.allowDownload = false",
				"security.allowUpload = false",
			} {
				needed[line] = true
			}
		case "attach_disabled", "attach_local_only":
			for _, line := range []string{
				"security.attach.enabled = false",
				"security.attach.allowHosts = 127.0.0.1,localhost,::1",
				"security.attach.allowSchemes = ws,wss",
			} {
				needed[line] = true
			}
		case "idpi_whitelist_scoped", "idpi_strict_mode", "idpi_content_protection":
			for _, line := range []string{
				"security.idpi.enabled = true",
				"security.idpi.allowedDomains = 127.0.0.1,localhost,::1",
				"security.idpi.strictMode = true",
				"security.idpi.scanContent = true",
				"security.idpi.wrapContent = true",
			} {
				needed[line] = true
			}
		}
	}

	lines := make([]string, 0, len(needed))
	for _, line := range ordered {
		if needed[line] {
			lines = append(lines, line)
		}
	}
	return lines
}

type securityDefaultsState struct {
	Bind     string
	Token    string
	Security securityConfigValues
}

type securityConfigValues struct {
	AllowEvaluate   bool
	AllowMacro      bool
	AllowScreencast bool
	AllowDownload   bool
	AllowUpload     bool
	MaxRedirects    int
	AttachEnabled   bool
	IDPI            config.IDPIConfig
}

func securityDefaultsSnapshot(fc *config.FileConfig) securityDefaultsState {
	if fc == nil {
		return securityDefaultsState{}
	}
	s := securityDefaultsState{
		Bind:  fc.Server.Bind,
		Token: fc.Server.Token,
		Security: securityConfigValues{
			IDPI: fc.Security.IDPI,
		},
	}
	if fc.Security.AllowEvaluate != nil {
		s.Security.AllowEvaluate = *fc.Security.AllowEvaluate
	}
	if fc.Security.AllowMacro != nil {
		s.Security.AllowMacro = *fc.Security.AllowMacro
	}
	if fc.Security.AllowScreencast != nil {
		s.Security.AllowScreencast = *fc.Security.AllowScreencast
	}
	if fc.Security.AllowDownload != nil {
		s.Security.AllowDownload = *fc.Security.AllowDownload
	}
	if fc.Security.AllowUpload != nil {
		s.Security.AllowUpload = *fc.Security.AllowUpload
	}
	if fc.Security.MaxRedirects != nil {
		s.Security.MaxRedirects = *fc.Security.MaxRedirects
	}
	if fc.Security.Attach.Enabled != nil {
		s.Security.AttachEnabled = *fc.Security.Attach.Enabled
	}
	return s
}
