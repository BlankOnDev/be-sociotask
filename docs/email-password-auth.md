# Autentikasi Email & Password

## Ringkasan Endpoint

- **Register**: `POST /users`
- **Login**: `POST /login`

## Register via Email

- **Endpoint**: `POST /users`
- **Handler**: `UserHandler.HandleCreateUser()`
- **Input body**:

```json
{
  "username": "janedoe",
  "email": "jane@example.com",
  "password": "pa55word",
  "bio": "Penjelasan singkat" // opsional
}
```

- **Validasi utama** (`validateRegisterRequest()`):
  - `username` wajib dan maksimal 50 karakter.
  - `email` wajib & harus mengikuti regex `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`.
  - `password` wajib dan minimal 8 karakter.

### Respons Berhasil (201)

```json
{
  "status": "success",
  "message": "registration successful",
  "data": {
    "user": {
      "id": 123,
      "username": "janedoe",
      "email": "jane@example.com",
      "bio": "Penjelasan singkat",
      "wallet_address": {
        "String": "",
        "Valid": false
      },
      "x_id": {
        "String": "",
        "Valid": false
      },
      "created_at": "2025-10-03T09:00:00Z",
      "updated_at": "2025-10-03T09:00:00Z"
    }
  }
}
```

> Catatan: Field `wallet_address` dan `x_id` menggunakan `sql.NullString` sehingga tampil sebagai objek dengan properti `String` dan `Valid`.

### Respons Gagal

- **Validasi gagal (400)**

```json
{
  "status": "error",
  "message": "validation failed",
  "errors": ["password must be at least 8 characters long"]
}
```

- **Email sudah digunakan / error DB (500)**

```json
{
  "status": "error",
  "message": "registration failed",
  "errors": null
}
```

## Login via Email

- **Endpoint**: `POST /login`
- **Handler**: `UserHandler.HandleLoginUser()`
- **Input body**:

```json
{
  "email": "jane@example.com",
  "password": "pa55word"
}
```

- **Validasi utama** (`validateLoginRequest()`):
  - `email` wajib & harus valid sesuai regex di atas.
  - `password` wajib.

### Respons Berhasil (200)

```json
{
  "status": "success",
  "message": "login successful",
  "data": {
    "token": "<JWT>"
  }
}
```

Token dibuat via `auth.GenerateJWTToken()` dengan masa berlaku 24 jam.

### Respons Gagal

- **Validasi gagal (400)**

```json
{
  "status": "error",
  "message": "validation failed",
  "errors": ["email is required"]
}
```

- **User tidak ditemukan (404)**

```json
{
  "status": "error",
  "message": "resource not found",
  "errors": null
}
```

- **Password salah (401)**

```json
{
  "status": "error",
  "message": "invalid credentials",
  "errors": null
}
```
