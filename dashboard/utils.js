'use strict';

// ---------------------------------------------------------------------------
// Shared utilities
// ---------------------------------------------------------------------------

function esc(s) {
  if (!s) return '';
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}

function timeAgo(d) {
  const s = Math.floor((Date.now() - d.getTime()) / 1000);
  if (s < 5) return 'just now';
  if (s < 60) return s + 's ago';
  if (s < 3600) return Math.floor(s / 60) + 'm ago';
  return Math.floor(s / 3600) + 'h ago';
}

function appConfirm(message, title, isDanger) {
  return new Promise((resolve) => {
    document.getElementById('confirm-title').textContent = title || 'Confirm';
    document.getElementById('confirm-message').textContent = message;
    const okBtn = document.getElementById('confirm-ok');
    okBtn.textContent = 'Confirm';
    okBtn.style.display = '';
    okBtn.className = isDanger !== false ? 'danger' : '';
    document.getElementById('confirm-cancel').textContent = 'Cancel';
    document.getElementById('confirm-modal').classList.add('open');
    const cleanup = () => { document.getElementById('confirm-modal').classList.remove('open'); };
    okBtn.onclick = () => { cleanup(); resolve(true); };
    document.getElementById('confirm-cancel').onclick = () => { cleanup(); resolve(false); };
  });
}

function appAlert(message, title) {
  return new Promise((resolve) => {
    document.getElementById('confirm-title').textContent = title || 'Notice';
    document.getElementById('confirm-message').textContent = message;
    document.getElementById('confirm-ok').style.display = 'none';
    document.getElementById('confirm-cancel').textContent = 'OK';
    document.getElementById('confirm-modal').classList.add('open');
    document.getElementById('confirm-cancel').onclick = () => {
      document.getElementById('confirm-modal').classList.remove('open');
      resolve();
    };
  });
}

function closeModal() {
  document.getElementById('modal').classList.remove('open');
}
