/**
 * Khaojai IT - Clean Architecture Frontend Application
 * Interacts with FastAPI Backend API (/api/v1)
 */

const API_BASE = '/api/v1';

// State Management
const state = {
  accessToken: localStorage.getItem('khaojai_access_token') || null,
  refreshToken: localStorage.getItem('khaojai_refresh_token') || null,
  user: null,
};

// DOM Elements
const authSection = document.getElementById('authSection');
const dashboardSection = document.getElementById('dashboardSection');
const loginForm = document.getElementById('loginForm');
const registerForm = document.getElementById('registerForm');
const tabLoginBtn = document.getElementById('tabLoginBtn');
const tabRegisterBtn = document.getElementById('tabRegisterBtn');
const healthBadge = document.getElementById('healthBadge');
const healthDot = document.getElementById('healthDot');
const healthText = document.getElementById('healthText');
const apiResponseConsole = document.getElementById('apiResponseConsole');
const responseStatusBadge = document.getElementById('responseStatusBadge');
const toastContainer = document.getElementById('toastContainer');

// Dashboard Elements
const dashFullName = document.getElementById('dashFullName');
const dashUserId = document.getElementById('dashUserId');
const dashEmail = document.getElementById('dashEmail');
const dashFirstName = document.getElementById('dashFirstName');
const dashLastName = document.getElementById('dashLastName');
const dashAccessToken = document.getElementById('dashAccessToken');
const dashRefreshToken = document.getElementById('dashRefreshToken');
const refreshTokensBtn = document.getElementById('refreshTokensBtn');
const logoutBtn = document.getElementById('logoutBtn');

// Initialize
document.addEventListener('DOMContentLoaded', () => {
  setupEventListeners();
  checkApiHealth();
  checkExistingSession();
});

function setupEventListeners() {
  // Tab Switcher
  tabLoginBtn.addEventListener('click', () => switchTab('login'));
  tabRegisterBtn.addEventListener('click', () => switchTab('register'));

  // Form Submissions
  loginForm.addEventListener('submit', handleLogin);
  registerForm.addEventListener('submit', handleRegister);

  // Dashboard Actions
  refreshTokensBtn.addEventListener('click', handleRefreshToken);
  logoutBtn.addEventListener('click', handleLogout);
}

function switchTab(tab) {
  if (tab === 'login') {
    tabLoginBtn.className = 'flex-1 py-2 text-sm font-semibold rounded-lg bg-indigo-600 text-white shadow transition';
    tabRegisterBtn.className = 'flex-1 py-2 text-sm font-semibold rounded-lg text-slate-400 hover:text-white transition';
    loginForm.classList.remove('hidden');
    registerForm.classList.add('hidden');
  } else {
    tabRegisterBtn.className = 'flex-1 py-2 text-sm font-semibold rounded-lg bg-brand-600 text-white shadow transition';
    tabLoginBtn.className = 'flex-1 py-2 text-sm font-semibold rounded-lg text-slate-400 hover:text-white transition';
    registerForm.classList.remove('hidden');
    loginForm.classList.add('hidden');
  }
}

// API Health Check
async function checkApiHealth() {
  try {
    const res = await fetch(`${API_BASE}/health`);
    const data = await res.json();
    if (res.ok && data.status === 'healthy') {
      healthDot.className = 'w-2 h-2 rounded-full bg-emerald-400';
      healthText.textContent = `API Online (${data.database.type})`;
      healthText.className = 'text-emerald-300 font-medium';
    } else {
      healthDot.className = 'w-2 h-2 rounded-full bg-amber-400';
      healthText.textContent = 'Degraded';
    }
  } catch (err) {
    healthDot.className = 'w-2 h-2 rounded-full bg-rose-400';
    healthText.textContent = 'API Offline';
    healthText.className = 'text-rose-300 font-medium';
  }
}

// Check if user is already logged in
async function checkExistingSession() {
  if (state.accessToken) {
    try {
      await fetchUserProfile();
    } catch (err) {
      // Access token may be expired, attempt refresh
      if (state.refreshToken) {
        try {
          await handleRefreshToken();
          return;
        } catch (e) {
          clearSession();
        }
      } else {
        clearSession();
      }
    }
  } else {
    showAuthView();
  }
}

// Register Handler
async function handleRegister(e) {
  e.preventDefault();
  const payload = {
    first_name: document.getElementById('regFirstName').value.trim(),
    last_name: document.getElementById('regLastName').value.trim(),
    email: document.getElementById('regEmail').value.trim(),
    password: document.getElementById('regPassword').value,
  };

  const submitBtn = document.getElementById('regSubmitBtn');
  submitBtn.disabled = true;
  submitBtn.innerHTML = '<span>Creating Account...</span>';

  try {
    const res = await fetch(`${API_BASE}/users/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    const data = await res.json();
    logConsoleResponse(res.status, data);

    if (res.status === 201) {
      showToast('Registration successful! Please sign in.', 'success');
      document.getElementById('loginEmail').value = payload.email;
      switchTab('login');
      registerForm.reset();
    } else {
      const errCode = data.detail?.code || data.code || 'users-001';
      const errMsg = data.detail?.message || data.message || 'Registration failed';
      showToast(`[${errCode}] ${errMsg}`, 'error');
    }
  } catch (err) {
    showToast(`Network Error: ${err.message}`, 'error');
  } finally {
    submitBtn.disabled = false;
    submitBtn.innerHTML = '<span>Create Account (สมัครสมาชิก)</span>';
  }
}

// Login Handler
async function handleLogin(e) {
  e.preventDefault();
  const payload = {
    email: document.getElementById('loginEmail').value.trim(),
    password: document.getElementById('loginPassword').value,
  };

  const submitBtn = document.getElementById('loginSubmitBtn');
  submitBtn.disabled = true;
  submitBtn.innerHTML = '<span>Signing In...</span>';

  try {
    const res = await fetch(`${API_BASE}/users/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });

    const data = await res.json();
    logConsoleResponse(res.status, data);

    if (res.status === 200) {
      setSession(data.token.access_token, data.token.refresh_token, data.user);
      showToast(`Welcome back, ${data.user.first_name}!`, 'success');
      loginForm.reset();
    } else {
      const errCode = data.detail?.code || data.code || 'users-002';
      const errMsg = data.detail?.message || data.message || 'Login failed';
      showToast(`[${errCode}] ${errMsg}`, 'error');
    }
  } catch (err) {
    showToast(`Network Error: ${err.message}`, 'error');
  } finally {
    submitBtn.disabled = false;
    submitBtn.innerHTML = '<span>Sign In to Account</span>';
  }
}

// Fetch Profile
async function fetchUserProfile() {
  const res = await fetch(`${API_BASE}/users/profile`, {
    headers: {
      'Authorization': `Bearer ${state.accessToken}`,
    },
  });

  const data = await res.json();
  logConsoleResponse(res.status, data);

  if (res.ok) {
    state.user = data;
    showDashboardView();
  } else {
    throw new Error('Profile fetch failed');
  }
}

// Refresh Token Handler
async function handleRefreshToken() {
  if (!state.refreshToken) {
    showToast('No refresh token available', 'error');
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/users/refresh-token`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: state.refreshToken }),
    });

    const data = await res.json();
    logConsoleResponse(res.status, data);

    if (res.ok) {
      setSession(data.token.access_token, data.token.refresh_token, data.user);
      showToast('Tokens refreshed successfully!', 'success');
    } else {
      const errCode = data.detail?.code || 'users-003';
      showToast(`[${errCode}] Token refresh failed, signing out`, 'error');
      clearSession();
    }
  } catch (err) {
    showToast(`Refresh error: ${err.message}`, 'error');
  }
}

// Logout Handler
async function handleLogout() {
  try {
    if (state.refreshToken) {
      await fetch(`${API_BASE}/users/logout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: state.refreshToken }),
      });
    }
  } catch (err) {
    console.warn('Logout API error:', err);
  } finally {
    clearSession();
    showToast('Logged out successfully', 'info');
  }
}

// Session Helpers
function setSession(accessToken, refreshToken, user) {
  state.accessToken = accessToken;
  state.refreshToken = refreshToken;
  state.user = user;

  localStorage.setItem('khaojai_access_token', accessToken);
  localStorage.setItem('khaojai_refresh_token', refreshToken);

  showDashboardView();
}

function clearSession() {
  state.accessToken = null;
  state.refreshToken = null;
  state.user = null;

  localStorage.removeItem('khaojai_access_token');
  localStorage.removeItem('khaojai_refresh_token');

  showAuthView();
}

function showAuthView() {
  authSection.classList.remove('hidden');
  dashboardSection.classList.add('hidden');
}

function showDashboardView() {
  authSection.classList.add('hidden');
  dashboardSection.classList.remove('hidden');

  if (state.user) {
    dashFullName.textContent = state.user.full_name || `${state.user.first_name} ${state.user.last_name}`;
    dashUserId.textContent = `#${state.user.id}`;
    dashEmail.textContent = state.user.email;
    dashFirstName.textContent = state.user.first_name;
    dashLastName.textContent = state.user.last_name;
  }

  dashAccessToken.value = state.accessToken || '';
  dashRefreshToken.value = state.refreshToken || '';
}

// Live Endpoint Tester
async function testEndpoint(type) {
  if (type === 'health') {
    const res = await fetch(`${API_BASE}/health`);
    const data = await res.json();
    logConsoleResponse(res.status, data);
  } else if (type === 'profile') {
    if (!state.accessToken) {
      logConsoleResponse(401, { code: 'users-003', message: 'No Bearer token available. Please sign in first.' });
      return;
    }
    const res = await fetch(`${API_BASE}/users/profile`, {
      headers: { 'Authorization': `Bearer ${state.accessToken}` },
    });
    const data = await res.json();
    logConsoleResponse(res.status, data);
  }
}

function logConsoleResponse(status, data) {
  responseStatusBadge.textContent = `HTTP ${status}`;
  if (status >= 200 && status < 300) {
    responseStatusBadge.className = 'font-mono font-bold px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30';
  } else {
    responseStatusBadge.className = 'font-mono font-bold px-2 py-0.5 rounded bg-rose-500/20 text-rose-400 border border-rose-500/30';
  }
  apiResponseConsole.textContent = JSON.stringify(data, null, 2);
}

function copyToClipboard(elementId) {
  const el = document.getElementById(elementId);
  if (el) {
    navigator.clipboard.writeText(el.value);
    showToast('Copied token to clipboard!', 'info');
  }
}

function showToast(message, type = 'info') {
  const toast = document.createElement('div');
  const borderCol = type === 'success' ? 'border-emerald-500/50 bg-emerald-950/90 text-emerald-200' :
                    type === 'error' ? 'border-rose-500/50 bg-rose-950/90 text-rose-200' :
                    'border-indigo-500/50 bg-slate-900/90 text-slate-200';

  toast.className = `p-4 rounded-xl border ${borderCol} shadow-2xl text-xs font-semibold backdrop-blur pointer-events-auto flex items-center justify-between space-x-3 transition transform duration-200 ease-out`;
  toast.innerHTML = `
    <span>${message}</span>
    <button onclick="this.parentElement.remove()" class="text-slate-400 hover:text-white">&times;</button>
  `;

  toastContainer.appendChild(toast);
  setTimeout(() => {
    toast.remove();
  }, 4000);
}
