# OAuth Google Login API Documentation

## Overview
API ini menggunakan OAuth 2.0 untuk autentikasi menggunakan akun Google. Proses login melibatkan 2 endpoint utama yang bekerja secara berurutan.

## Flow OAuth
```
1. Frontend → GET /auth/google/login
2. Backend → Generate state dan set cookie
3. Backend → Redirect ke Google OAuth page
4. User → Login di Google dan authorize app
5. Google → Redirect ke /auth/google/callback?code=xxx&state=xxx
6. Backend → Exchange code untuk access token
7. Backend → Fetch user data dari Google API
8. Backend → Create/update user di database
9. Backend → Generate JWT token
10. Backend → Redirect ke success page dengan token
```

---

## 1. Initiate Google Login

### Endpoint
`GET /auth/google/login`

### Description
Endpoint ini memulai proses OAuth dengan Google. Backend akan generate state string untuk keamanan CSRF, menyimpannya dalam cookie, lalu redirect user ke halaman login Google.

### Request
Tidak memerlukan body atau query parameters. Cukup redirect user ke endpoint ini.

### Response
**Status Code**: `307 Temporary Redirect`

Backend akan:
1. Generate random state string (32 bytes, encoded base64)
2. Set cookie `oauthstate` untuk CSRF protection
3. Redirect ke Google OAuth authorization URL

**Cookie yang di-set:**
```
oauthstate: "base64_encoded_random_string"
  - Path: /
  - Expires: 10 minutes
  - HttpOnly: true
```

### Error Response

#### Failed to Generate State
**Status Code**: `307 Temporary Redirect`

Akan redirect ke `/failed?error=state_generation_failed`

**Penyebab:**
- Gagal generate random bytes untuk state string

### Frontend Implementation Example
```javascript
// Redirect user ke endpoint login
window.location.href = 'https://your-api.com/auth/google/login';
```

---

## 2. OAuth Callback

### Endpoint
`GET /auth/google/callback`

### Description
Endpoint ini dipanggil oleh Google setelah user berhasil login dan authorize aplikasi. Backend akan memproses authorization code, exchange untuk access token, fetch user data dari Google, dan generate JWT token.

### Query Parameters
- **code** (string, required): Authorization code dari Google
- **state** (string, required): State string untuk validasi CSRF

**Note:** Parameters ini dikirim otomatis oleh Google, bukan dari frontend.

### Success Response
**Status Code**: `307 Temporary Redirect`

Redirect ke: `/success?token={jwt_token}`

Backend akan:
1. Validasi state parameter
2. Exchange authorization code dengan access token
3. Fetch user info dari Google API
4. Create atau update user di database
5. Generate JWT token
6. Redirect ke success page dengan token

### Error Responses

#### 1. Invalid State Parameter
**Status Code**: `307 Temporary Redirect`

Redirect ke: `/failed?error=invalid_state`

**Penyebab:**
- Cookie `oauthstate` tidak ditemukan
- State parameter dari Google tidak match dengan cookie

#### 2. Token Exchange Failed
**Status Code**: `307 Temporary Redirect`

Redirect ke: `/failed?error=token_exchange_failed`

**Penyebab:**
- Gagal exchange authorization code dengan access token
- Authorization code tidak valid atau expired

#### 3. Failed to Fetch User Data
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "oauth authentication failed",
  "data": null,
  "errors": ["failed to get user info"]
}
```

**Penyebab:**
- Gagal mengambil data user dari Google API (`https://www.googleapis.com/oauth2/v2/userinfo`)
- Access token tidak valid

#### 4. Failed to Parse User Info
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "oauth authentication failed",
  "data": null,
  "errors": ["failed to parse user info"]
}
```

**Penyebab:**
- Gagal unmarshal response JSON dari Google API
- Format response tidak sesuai ekspektasi

#### 5. Database Operation Failed
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "internal server error",
  "data": null,
  "errors": ["database operation failed"]
}
```

**Penyebab:**
- Error saat find atau create user di database
- Database connection error

#### 6. Failed to Generate JWT Token
**Status Code**: `500 Internal Server Error`

```json
{
  "status": "error",
  "message": "internal server error",
  "data": null,
  "errors": ["failed to generate token"]
}
```

**Penyebab:**
- Error saat generate JWT token
- JWT secret key tidak valid

---

## Google User Info Response

Data yang diambil dari Google API:

```json
{
  "id": "1234567890",
  "email": "user@example.com",
  "name": "John Doe"
}
```

**Fields:**
- **id**: Google user ID (unique identifier)
- **email**: Email address user
- **name**: Display name user

---

## Frontend Integration Guide

### Step 1: Initiate Login
```javascript
function loginWithGoogle() {
  // Redirect user ke endpoint OAuth
  window.location.href = 'https://your-api.com/auth/google/login';
}
```

### Step 2: Handle Success Callback
Backend akan redirect ke `/success?token={jwt_token}` setelah login berhasil.

#### Option A: Handle di Success Page
```javascript
// Di success page (e.g., /success)
useEffect(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const token = urlParams.get('token');
  
  if (token) {
    // Save token
    localStorage.setItem('authToken', token);
    
    // Optional: Decode JWT untuk get user_id
    const payload = JSON.parse(atob(token.split('.')[1]));
    localStorage.setItem('userId', payload.user_id);
    
    // Redirect ke dashboard
    window.location.href = '/dashboard';
  }
}, []);
```

#### Option B: Popup Window (Recommended)
```javascript
function loginWithGooglePopup() {
  // Open OAuth flow di popup window
  const width = 600;
  const height = 700;
  const left = (screen.width - width) / 2;
  const top = (screen.height - height) / 2;
  
  const popup = window.open(
    'https://your-api.com/auth/google/login',
    'Google Login',
    `width=${width},height=${height},left=${left},top=${top}`
  );
  
  // Check popup untuk URL changes
  const checkPopup = setInterval(() => {
    try {
      if (popup.closed) {
        clearInterval(checkPopup);
        return;
      }
      
      // Check jika redirect ke success page
      if (popup.location.href.includes('/success')) {
        const url = new URL(popup.location.href);
        const token = url.searchParams.get('token');
        
        if (token) {
          // Save token
          localStorage.setItem('authToken', token);
          
          // Close popup
          popup.close();
          clearInterval(checkPopup);
          
          // Redirect atau update UI
          window.location.href = '/dashboard';
        }
      }
    } catch (e) {
      // Cross-origin error, popup masih di domain lain
    }
  }, 500);
}
```

### Step 3: Handle Failed Callback
Backend akan redirect ke `/failed?error={error_type}` jika terjadi error.

```javascript
// Di failed page (e.g., /failed)
useEffect(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const error = urlParams.get('error');
  
  if (error) {
    let errorMessage = 'Login failed. Please try again.';
    
    switch (error) {
      case 'invalid_state':
        errorMessage = 'Security validation failed. Please try again.';
        break;
      case 'token_exchange_failed':
        errorMessage = 'Failed to complete authentication. Please try again.';
        break;
      case 'state_generation_failed':
        errorMessage = 'System error. Please try again later.';
        break;
    }
    
    // Show error message
    alert(errorMessage);
    
    // Redirect ke login page
    setTimeout(() => {
      window.location.href = '/login';
    }, 2000);
  }
}, []);
```

### Step 4: Use JWT Token
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

## Example Complete Flow

### 1. User clicks "Login with Google"
```javascript
<button onClick={() => window.location.href = '/auth/google/login'}>
  Login with Google
</button>
```

### 2. Backend generates state and redirects
```
User → Backend: GET /auth/google/login
Backend → Sets cookie: oauthstate=random_string
Backend → Redirects to: https://accounts.google.com/o/oauth2/auth?...
```

### 3. User authorizes on Google
```
User → Google: Login and authorize
Google → Validates credentials
Google → Redirects to: /auth/google/callback?code=xxx&state=xxx
```

### 4. Backend processes callback
```
Backend → Validates state
Backend → Exchanges code for token
Backend → Fetches user info from Google
Backend → Creates/updates user in database
Backend → Generates JWT token
Backend → Redirects to: /success?token=jwt_token
```

### 5. Frontend saves token
```javascript
// On success page
const token = new URLSearchParams(window.location.search).get('token');
localStorage.setItem('authToken', token);
window.location.href = '/dashboard';
```

---
