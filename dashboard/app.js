'use strict';

// ---------------------------------------------------------------------------
// App init — view switching, SSE connection
// ---------------------------------------------------------------------------

let profilesInterval = null;

function switchView(view) {
  document.querySelectorAll('.view-btn').forEach(b => b.classList.remove('active'));
  document.querySelector('[data-view="' + view + '"]').classList.add('active');
  document.getElementById('feed-view').style.display = view === 'feed' ? 'flex' : 'none';
  document.getElementById('profiles-view').style.display = view === 'profiles' ? 'flex' : 'none';
  document.getElementById('live-view').style.display = view === 'live' ? 'flex' : 'none';
  document.getElementById('settings-view').style.display = view === 'settings' ? 'block' : 'none';

  if (view === 'live') refreshTabs();
  if (view === 'profiles') loadProfiles();
  if (view === 'settings') loadSettings();

  // Auto-refresh profiles every 3s while on that view
  if (profilesInterval) { clearInterval(profilesInterval); profilesInterval = null; }
  if (view === 'profiles') {
    profilesInterval = setInterval(loadProfiles, 3000);
  }
}

function openInstanceDirect(port) {
  window.open('http://localhost:' + port + '/dashboard', '_blank');
}

// Boot
connect();
