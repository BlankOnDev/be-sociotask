# OAuth X (Twitter) Login API Documentation

## Overview
API ini menggunakan OAuth 2.0 dengan PKCE (Proof Key for Code Exchange) untuk autentikasi menggunakan akun X (Twitter). Proses login melibatkan 2 endpoint utama yang bekerja secara berurutan.

## Flow OAuth
```
1. Frontend → GET /auth/twitter/login
2. Backend → Redirect ke X OAuth page
3. User → Login di X dan authorize app
4. X → Redirect ke /auth/twitter/callback?code=xxx&state=xxx
5. Backend → Exchange code untuk access token
6. Backend → Fetch user data dari X API
7. Backend → Create/update user di database
8. Backend → Generate JWT token
9. Backend → Return JWT token ke frontend
```

---

## 1. Initiate Twitter Login

### Endpoint
`GET /auth/twitter/login`

### Description
Endpoint ini memulai proses OAuth dengan X (Twitter). Backend akan generate state dan PKCE verifier untuk keamanan, menyimpannya dalam cookies, lalu redirect user ke halaman login X.

### Request
Tidak memerlukan body atau query parameters. Cukup redirect user ke endpoint ini.

### Response
**Status Code**: `307 Temporary Redirect`

Backend akan:
1. Generate random state string (32 karakter)
2. Generate PKCE verifier
3. Set 2 cookies:
   - `oauth2_state`: State string untuk CSRF protection
   - `oauth2_verifier`: PKCE verifier untuk security
4. Redirect ke X OAuth authorization URL

**Cookies yang di-set:**
```
oauth2_state: "random_32_character_string"
  - Path: /
  - Expires: 15 minutes
  - HttpOnly: true
  - Secure: true
  - SameSite: Lax

oauth2_verifier: "pkce_verifier_string"
  - Path: /
  - Expires: 15 minutes
  - HttpOnly: true
  - Secure: true
  - SameSite: Lax
```

### Frontend Implementation Example
```javascript
// Redirect user ke endpoint login
window.location.href = 'https://your-api.com/auth/twitter/login';
```

---

## 2. OAuth Callback

### Endpoint
`GET /auth/twitter/callback`

### Description
Endpoint ini dipanggil oleh X setelah user berhasil login dan authorize aplikasi. Backend akan memproses authorization code, exchange untuk access token, fetch user data dari X, dan generate JWT token.

### Query Parameters
- **code** (string, required): Authorization code dari X
- **state** (string, required): State string untuk validasi CSRF

**Note:** Parameters ini dikirim otomatis oleh X, bukan dari frontend.

### Success Response
**Status Code**: `200 OK`

```json
{
  "status": "success",
  "message": "OAuth authentication successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 123
  },
  "errors": null
}
```

**Response Fields:**
- **token**: JWT token yang bisa digunakan untuk authenticated requests
- **user_id**: ID user di database sistem kita

### Error Responses

#### 1. Invalid State Parameter
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "OAuth authentication failed",
  "data": null,
  "errors": null
}
```

**Penyebab:**
- Cookie `oauth2_state` tidak ditemukan
- State parameter dari X tidak match dengan cookie

#### 2. Missing Authorization Code
**Status Code**: `400 Bad Request`

```json
{
  "status": "error",
  "message": "OAuth authentication failed",
  "data": null,
  "errors": null
}
```

**Penyebab:**
- Parameter `code` tidak ada dalam query string

#### 3. Token Exchange Failed
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "OAuth authentication failed",
  "data": null,
  "errors": null
}
```

**Penyebab:**
- Gagal exchange authorization code dengan access token
- PKCE verifier tidak valid

#### 4. Failed to Fetch User Data
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "OAuth authentication failed",
  "data": null,
  "errors": null
}
```

**Penyebab:**
- Gagal mengambil data user dari X API
- X API mengembalikan error

#### 5. User Creation Failed (Conflict)
**Status Code**: `409 Conflict`

```json
{
  "status": "error",
  "message": "Registration failed",
  "data": null,
  "errors": null
}
```

**Penyebab:**
- Username atau email sudah terdaftar (unique constraint violation)

#### 6. Internal Server Error
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "Internal server error",
  "data": null,
  "errors": null
}
```

**Penyebab:**
- Error saat create user baru di database
- Error saat generate JWT token
- Error internal lainnya

---


## Frontend Integration Guide

### Step 1: Initiate Login
```javascript
function loginWithTwitter() {
  // Redirect user ke endpoint OAuth
  window.location.href = 'https://your-api.com/auth/twitter/login';
}
```

### Step 2: Handle Callback
Backend akan otomatis handle callback dari X. Frontend perlu setup route untuk menerima response.

#### Option A: Popup Window (Recommended)
```javascript
function loginWithTwitterPopup() {
  // Open OAuth flow di popup window
  const width = 600;
  const height = 700;
  const left = (screen.width - width) / 2;
  const top = (screen.height - height) / 2;
  
  const popup = window.open(
    'https://your-api.com/auth/twitter/login',
    'Twitter Login',
    `width=${width},height=${height},left=${left},top=${top}`
  );
  
  // Listen untuk message dari popup
  window.addEventListener('message', (event) => {
    if (event.origin !== 'https://your-api.com') return;
    
    if (event.data.token) {
      // Save token
      localStorage.setItem('authToken', event.data.token);
      localStorage.setItem('userId', event.data.user_id);
      
      // Close popup
      popup.close();
      
      // Redirect atau update UI
      window.location.href = '/dashboard';
    }
  });
}
```

**Backend perlu tambahan di callback handler untuk kirim postMessage:**
```javascript
// Di callback success response, tambahkan HTML yang kirim postMessage
const html = `
  <script>
    window.opener.postMessage(
      { token: '${token}', user_id: ${userId} },
      'https://your-frontend.com'
    );
    window.close();
  </script>
`;
```

#### Option B: Full Page Redirect
```javascript
// Setup callback route di frontend (e.g., /auth/callback)
// Backend redirect ke route ini dengan token sebagai query param

// Di callback route component:
useEffect(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const token = urlParams.get('token');
  const userId = urlParams.get('user_id');
  
  if (token) {
    // Save token
    localStorage.setItem('authToken', token);
    localStorage.setItem('userId', userId);
    
    // Redirect ke dashboard
    navigate('/dashboard');
  }
}, []);
```

### Step 3: Use JWT Token
```javascript
// Include token di setiap authenticated request
const token = localStorage.getItem('authToken');

fetch('https://your-api.com/api/protected-endpoint', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
});
```

---

## Security Considerations

### CSRF Protection
- Sistem menggunakan random state string untuk protect dari CSRF attacks
- State disimpan di HttpOnly cookie dan divalidasi saat callback

### PKCE (Proof Key for Code Exchange)
- Menggunakan PKCE flow untuk additional security
- Verifier disimpan di HttpOnly cookie dan di-exchange saat token exchange

### Cookie Security
- `HttpOnly`: true - Tidak bisa diakses via JavaScript
- `Secure`: true - Hanya dikirim via HTTPS
- `SameSite`: Lax - Protection dari CSRF
- `Expires`: 15 minutes - Limited lifetime

### Token Security
- JWT token berisi user ID dan role
- Token harus disimpan dengan aman (localStorage atau httpOnly cookie)
- Token harus di-include di Authorization header untuk authenticated requests

---
