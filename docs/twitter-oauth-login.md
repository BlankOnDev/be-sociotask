# Twitter OAuth Login

## Ringkasan Endpoint
- **Mulai login**: `GET /login/twitter`
- **Callback**: `GET /login/twitter/callback`

## Alur Singkat
1. **Redirect ke backend**. Frontend mengarahkan pengguna (biasanya `window.location` atau membuka tab/popup) ke `GET /login/twitter`.
2. **Backend siapkan sesi OAuth**. Handler `HandleTwitterLogin()` menyimpan `oauth2_state` dan `oauth2_verifier` dalam cookie `HttpOnly`, `Secure`, berlaku 15 menit, lalu mengarahkan pengguna ke halaman otorisasi X/Twitter.
3. **Pengguna mengotorisasi**. Setelah izin diberikan, X mengarahkan kembali ke `/login/twitter/callback?state=...&code=...`.
4. **Validasi & tukar kode**. Handler `HandleTwitterCallback()` memverifikasi cookie `oauth2_state`, mengambil kode PKCE dari cookie `oauth2_verifier`, dan menukar authorization code dengan access token Twitter.
5. **Ambil data user**. Backend memanggil `https://api.twitter.com/2/users/me` untuk mendapatkan `id`, `name`, `username`, dan (jika ada) `email`. Jika email kosong, backend membuat placeholder `<username>@twitter.user` untuk memenuhi syarat unik di database.
6. **Sinkronisasi user & JWT**. Jika user belum ada, backend membuat akun baru dengan password acak. Lalu backend membuat JWT melalui `auth.GenerateJWTToken()` dan mengembalikannya ke frontend.

## Format Respons
Berhasil:
```json
{
 "status": "success",
 "message": "oauth authentication successful",
 "data": {
  "token": "<JWT>",
  "user_id": 123
 }
}
```

Gagal (contoh):
```json
{
 "status": "error",
 "message": "oauth authentication failed",
 "errors": null
}
```

## Panduan Frontend
- **Mulai proses** dengan membuka `/login/twitter` (tab baru atau popup) agar user dapat menyelesaikan login di X.
- **Terima respons JSON** dari `/login/twitter/callback`. Jika memakai popup, tutup popup setelah menerima respons dan kirim `token` + `user_id` ke aplikasi utama.
- **Simpan token** untuk kebutuhan autentikasi berikutnya (mis. `localStorage` atau mekanisme lain sesuai kebutuhan aplikasi).
- **Tangani error** dengan menampilkan pesan dari field `message` atau `errors`.
- Pastikan aplikasi berjalan melalui HTTPS; cookie `Secure` akan diabaikan jika backend diakses tanpa TLS.
