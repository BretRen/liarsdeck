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

export function useAuth() {
  const isAuthenticated = computed(() => !!user.value && !!user.value.name);
  const nickname = computed(() => {
    if (!user.value) return '';
    return user.value.nickname || user.value.preferred_username || user.value.name || '';
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

      // 2. 拉取 Userinfo
      const userRes = await fetch(USERINFO_ENDPOINT, {
        headers: {
          Authorization: `Bearer ${tokenData.access_token}`,
        },
      });

      if (!userRes.ok) {
        throw new Error(`获取用户信息失败: HTTP ${userRes.status}`);
      }

      const userData = await userRes.json();
      const finalName = userData.nickname || userData.preferred_username || userData.name || 'Player';

      const userProfile = {
        sub: userData.sub,
        name: finalName,
        preferred_username: userData.preferred_username || '',
        nickname: userData.nickname || '',
        email: userData.email || '',
        picture: userData.picture || '',
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
