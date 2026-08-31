import { ref, computed } from 'vue';

const CLIENT_ID = '388639280752296451';
const AUTH_ENDPOINT = 'https://auth.pdnode.com/oauth/v2/authorize';
const TOKEN_ENDPOINT = 'https://auth.pdnode.com/oauth/v2/token';
const USERINFO_ENDPOINT = 'https://auth.pdnode.com/oidc/v1/userinfo';
const END_SESSION_ENDPOINT = 'https://auth.pdnode.com/oidc/v1/end_session';

const USER_STORAGE_KEY = 'liarsdeck_auth_user';
const TOKEN_STORAGE_KEY = 'liarsdeck_auth_tokens';

// Global singleton reactive state
const user = ref(loadSavedUser());
const tokens = ref(loadSavedTokens());
const isLoggingIn = ref(false);
const authError = ref('');

function loadSavedUser() {
  try {
    const raw = localStorage.getItem(USER_STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch (_) {
    return null;
  }
}

function loadSavedTokens() {
  try {
    const raw = localStorage.getItem(TOKEN_STORAGE_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch (_) {
    return null;
  }
}

// ── PKCE Cryptographic Helpers ──

function base64UrlEncode(arrayBuffer) {
  const bytes = new Uint8Array(arrayBuffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary)
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}

function generateRandomString(length = 64) {
  const validChars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~';
  const array = new Uint8Array(length);
  window.crypto.getRandomValues(array);
  return Array.from(array, (byte) => validChars[byte % validChars.length]).join('');
}

async function generateCodeChallenge(verifier) {
  const encoder = new TextEncoder();
  const data = encoder.encode(verifier);
  const digest = await window.crypto.subtle.digest('SHA-256', data);
  return base64UrlEncode(digest);
}

function getRedirectUri() {
  const host = window.location.hostname;
  if (host === 'localhost' || host === '127.0.0.1') {
    return `${window.location.origin}/callback`;
  }
  return 'https://liarsbar.games.pdnode.com/callback';
}

function parseJwt(token) {
  try {
    if (!token || typeof token !== 'string') return {};
    const parts = token.split('.');
    if (parts.length < 2) return {};
    const base64Url = parts[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch (_) {
    return {};
  }
}

// 提取最佳昵称（优先使用 nickname / preferred_username / username，绝不直接使用邮箱地址）
function extractBestNickname(userData = {}, idTokenClaims = {}) {
  const candidates = [
    userData.nickname,
    userData.preferred_username,
    userData.username,
    idTokenClaims.nickname,
    idTokenClaims.preferred_username,
    idTokenClaims.username,
    userData.given_name,
    userData.name,
    idTokenClaims.name,
  ];

  // 1. 优先寻找不含 '@' 的有效纯昵称/用户名
  for (const c of candidates) {
    if (c && typeof c === 'string' && c.trim() && !c.includes('@')) {
      return c.trim();
    }
  }

  // 2. 若字段仅包含邮箱格式，提取 '@' 之前的前缀作为昵称
  for (const c of candidates) {
    if (c && typeof c === 'string' && c.trim()) {
      if (c.includes('@')) {
        return c.split('@')[0].trim();
      }
      return c.trim();
    }
  }

  if (userData.email && typeof userData.email === 'string') {
    return userData.email.split('@')[0].trim();
  }

  return 'Player';
}

export function useAuth() {
  const isAuthenticated = computed(() => !!user.value && (!!user.value.name || !!user.value.nickname));

  const nickname = computed(() => {
    if (!user.value) return '';
    return extractBestNickname(user.value);
  });

  const avatar = computed(() => (user.value ? user.value.picture || '' : ''));

  async function login() {
    isLoggingIn.value = true;
    authError.value = '';

    try {
      const verifier = generateRandomString(96);
      const challenge = await generateCodeChallenge(verifier);
      const state = generateRandomString(32);
      const redirectUri = getRedirectUri();

      sessionStorage.setItem('liarsdeck_pkce_verifier', verifier);
      sessionStorage.setItem('liarsdeck_oauth_state', state);
      sessionStorage.setItem('liarsdeck_auth_redirect_uri', redirectUri);

      const params = new URLSearchParams({
        response_type: 'code',
        client_id: CLIENT_ID,
        redirect_uri: redirectUri,
        scope: 'openid profile email',
        state: state,
        code_challenge: challenge,
        code_challenge_method: 'S256',
      });

      window.location.href = `${AUTH_ENDPOINT}?${params.toString()}`;
    } catch (err) {
      isLoggingIn.value = false;
      authError.value = err.message || '初始化登录授权失败';
    }
  }

  async function handleCallback() {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    const state = urlParams.get('state');
    const error = urlParams.get('error');
    const errorDescription = urlParams.get('error_description');

    if (error) {
      authError.value = errorDescription || error;
      cleanUrlQuery();
      return false;
    }

    if (!code) {
      return false;
    }

    const savedState = sessionStorage.getItem('liarsdeck_oauth_state');
    const verifier = sessionStorage.getItem('liarsdeck_pkce_verifier');
    const redirectUri = sessionStorage.getItem('liarsdeck_auth_redirect_uri') || getRedirectUri();

    if (savedState && state && savedState !== state) {
      authError.value = 'OAuth 状态校验失败 (State mismatch)，可能受到跨站请求伪造攻击';
      cleanUrlQuery();
      return false;
    }

    if (!verifier) {
      authError.value = 'PKCE 校验凭据丢失，请重新发起登录';
      cleanUrlQuery();
      return false;
    }

    isLoggingIn.value = true;

    try {
      // 1. 兑换 Access Token
      const body = new URLSearchParams({
        grant_type: 'authorization_code',
        client_id: CLIENT_ID,
        code: code,
        redirect_uri: redirectUri,
        code_verifier: verifier,
      });

      const tokenRes = await fetch(TOKEN_ENDPOINT, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: body.toString(),
      });

      if (!tokenRes.ok) {
        const errJson = await tokenRes.json().catch(() => ({}));
        throw new Error(errJson.error_description || errJson.error || `HTTP ${tokenRes.status}`);
      }

      const tokenData = await tokenRes.json();
      tokens.value = tokenData;
      localStorage.setItem(TOKEN_STORAGE_KEY, JSON.stringify(tokenData));

      // 解析 ID Token 中的 Claims
      const idTokenClaims = parseJwt(tokenData.id_token);

      // 2. 拉取 Userinfo
      let userData = {};
      try {
        const userRes = await fetch(USERINFO_ENDPOINT, {
          headers: {
            Authorization: `Bearer ${tokenData.access_token}`,
          },
        });
        if (userRes.ok) {
          userData = await userRes.json();
        }
      } catch (_) {}

      // 3. 解析最优实名昵称（避免邮箱作为昵称）
      const finalNickname = extractBestNickname(userData, idTokenClaims);

      const userProfile = {
        sub: userData.sub || idTokenClaims.sub || '',
        name: finalNickname,
        nickname: userData.nickname || idTokenClaims.nickname || finalNickname,
        preferred_username: userData.preferred_username || idTokenClaims.preferred_username || '',
        username: userData.username || idTokenClaims.username || '',
        email: userData.email || idTokenClaims.email || '',
        picture: userData.picture || idTokenClaims.picture || '',
      };

      user.value = userProfile;
      localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(userProfile));

      // 清理临时凭证
      sessionStorage.removeItem('liarsdeck_pkce_verifier');
      sessionStorage.removeItem('liarsdeck_oauth_state');
      sessionStorage.removeItem('liarsdeck_auth_redirect_uri');

      cleanUrlQuery();
      isLoggingIn.value = false;
      return true;
    } catch (err) {
      isLoggingIn.value = false;
      authError.value = '授权登录失败: ' + (err.message || '未知错误');
      cleanUrlQuery();
      return false;
    }
  }

  function cleanUrlQuery() {
    const url = new URL(window.location.href);
    url.searchParams.delete('code');
    url.searchParams.delete('state');
    url.searchParams.delete('error');
    url.searchParams.delete('error_description');
    url.searchParams.delete('session_state');
    if (url.pathname.endsWith('/callback')) {
      window.history.replaceState({}, document.title, '/' + url.search);
    } else {
      window.history.replaceState({}, document.title, url.pathname + url.search);
    }
  }

  function logout() {
    user.value = null;
    tokens.value = null;
    localStorage.removeItem(USER_STORAGE_KEY);
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    localStorage.removeItem('liarsdeck_active_session');

    const redirectUri = window.location.origin;
    const endSessionUrl = `${END_SESSION_ENDPOINT}?client_id=${CLIENT_ID}&post_logout_redirect_uri=${encodeURIComponent(
      redirectUri
    )}`;

    // 跳转 OIDC 单点登出
    window.location.href = endSessionUrl;
  }

  return {
    user,
    tokens,
    isAuthenticated,
    isLoggingIn,
    authError,
    nickname,
    avatar,
    login,
    handleCallback,
    logout,
  };
}
