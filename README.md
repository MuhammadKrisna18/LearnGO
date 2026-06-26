# Modular Monolith System

Sebuah sistem terintegrasi dengan arsitektur **Modular Monolith**, dikembangkan menggunakan **Go (Fiber)** di sisi backend dan **SvelteKit (Svelte 5 Runes)** di sisi frontend. Proyek ini dirancang untuk memberikan performa yang tinggi, struktur kode yang bersih, serta desain antarmuka pengguna (UI) modern berbasis *Glassmorphism*.

## 🚀 Fitur Utama

- **Modular Architecture**: Struktur backend dibagi berdasarkan modul (contoh: modul `auth` dan `matakuliah`) yang membuat kode lebih terisolasi dan mudah di-maintain.
- **Frontend Modular Design**: Pemisahan antarmuka (UI) menjadi komponen-komponen terisolasi di `src/lib/components` dipadukan dengan manajemen *Global State* Svelte 5 Runes untuk optimasi reaktivitas dan caching data.
- **Modern Authentication**: Implementasi JWT (JSON Web Token) dengan kapabilitas *Role-based Access Control* (RBAC).
- **Custom Performance Logger**: Middleware khusus untuk mengukur latensi eksekusi setiap request API, dikategorikan secara *real-time* ke dalam tiga level: `FAST`, `MODERATE`, dan `SLOW`.
- **Role-based Dashboards**: Tampilan dashboard spesifik untuk masing-masing role. Serta pemisahan Sidebar Navigation untuk fitur Admin.
- **Admin Management**: Fitur khusus bagi admin untuk mendaftarkan akun baru (seperti akun dosen) serta fitur **Manajemen Mata Kuliah** lengkap dengan validasi unik.
- **Svelte 5 Runes**: Memanfaatkan fitur *reactivity* modern dari Svelte 5 (`$state`) yang membuat manajemen *state* frontend lebih efisien.
- **Database & Cache**: Terhubung dengan **PostgreSQL** (melalui GORM) untuk database persisten dan **Redis** untuk manajemen *cache*.
	- **Swagger Documentation**: Dokumentasi endpoint API otomatis yang dapat diakses dengan mudah untuk kebutuhan *development*.

---

## 🏆 Quality Attributes

Sistem ini dibangun dengan memprioritaskan NFR (Non-Functional Requirements) berikut:
- **Maintainability (Keterpeliharaan)**: Diwujudkan melalui pemisahan domain pada *backend* (Modular Monolith) dan pemisahan komponen secara atomik pada *frontend* (SvelteKit).
- **Performance (Kinerja)**: Backend Go Fiber sangat cepat dalam menangani *request* HTTP, didukung dengan *caching* Redis, serta manajemen *state* reaktif menggunakan *Runes* pada Svelte 5.
- **Security (Keamanan)**: Diimplementasikan menggunakan autentikasi *stateless* (JWT), enkripsi *password* menggunakan Bcrypt, dan sistem RBAC (Role-Based Access Control) yang ketat di level API maupun UI.
- **Usability (Kebergunaan)**: Pengalaman antarmuka (*User Experience*) difokuskan dengan desain visual *Glassmorphism*, transisi antar halaman instan (*Single Page Application* feel), serta *feedback* yang responsif (pesan sukses/error).
- **Scalability (Skalabilitas)**: Arsitektur *stateless* memungkinkan aplikasi untuk mudah di-*scale* (dikembangkan) ke depannya, baik dari sisi fungsionalitas (menambah modul baru) maupun dari sisi infrastruktur (menambah *instance* server).

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Go Fiber (v2)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Cache**: Redis
- **Security**: JWT & Bcrypt
- **API Docs**: Swagger (Swaggo)

### Frontend
- **Framework**: SvelteKit (Svelte 5)
- **Styling**: Vanilla CSS (Glassmorphism & Modern Animations)
- **Build Tool**: Vite

---

## 📂 Struktur Proyek

```text
📦 Project Root
 ┣ 📂 cmd/api          # Main entry point aplikasi Go
 ┣ 📂 config           # Konfigurasi environment (DB, Redis, JWT)
 ┣ 📂 docs             # Dokumentasi auto-generated (Swagger)
 ┣ 📂 frontend         # Direktori frontend SvelteKit
 ┃  ┣ 📂 src/lib       # Komponen reusable, API Services, dan Global Stores (Svelte 5 Runes)
 ┃  ┣ 📂 src/routes    # Halaman aplikasi dengan Layout Sidebar
 ┃  ┗ 📜 app.css       # File core CSS dengan variable desain
 ┣ 📂 internal         # Core logic dari backend Modular Monolith
 ┃  ┣ 📂 app           # Registrasi aplikasi dan middleware (Fiber)
 ┃  ┣ 📂 middleware    # Custom Middleware (Performance Logger, JWT)
 ┃  ┣ 📂 modules       # Folder modular domain (contoh: Auth module, Matakuliah module)
 ┃  ┗ 📂 shared        # Logic shared antar module (DB, Cache, Error Handler)
 ┗ 📜 README.md
```

---

## 💻 Cara Menjalankan Aplikasi (Local Development)

### 1. Persyaratan Sistem
Pastikan sistem Anda telah menginstall:
- Go (v1.20+)
- Node.js (v18+)
- PostgreSQL (berjalan di port default 5432)
- Redis (berjalan di port default 6379)

### 2. Konfigurasi Environment Backend
Buat file `.env` di *root directory* atau sesuaikan variabel yang ada di dalam `config/config.go`:
```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=supersecret-jwt-key
```

### 3. Menjalankan Backend (Go)
Buka terminal di *root directory* dan jalankan perintah:
```bash
go run cmd/api/main.go
```
*Backend akan berjalan di `http://localhost:8080`*

*(Opsional) Jika ingin melihat dokumentasi Swagger, akses: `http://localhost:8080/api/v1/swagger/`*

### 4. Menjalankan Frontend (SvelteKit)
Buka terminal baru, masuk ke direktori `frontend`, lakukan instalasi dependensi, lalu jalankan development server:
```bash
cd frontend
npm install
npm run dev
```
*Frontend akan berjalan di `http://localhost:5173`*

---

## 🚦 Flow Aplikasi (Login to Dashboard)

1. **Akses**: Pengguna membuka `http://localhost:5173/login`.
2. **Autentikasi**: Pengguna memasukkan Email dan Password.
3. **Validasi**: Data dikirim ke API `/api/v1/auth/login`. Jika sukses, Backend mengembalikan JWT dan *Role* (Admin/User).
4. **Performa**: Secara bersamaan di console backend, performa request akan dilog (`FAST` / `MODERATE` / `SLOW`). Di console browser, tercatat latensi dari sisi klien.
5. **Splash Screen**: Frontend menyimpan token dan role ke Global Store memory dan langsung memunculkan Splash Screen dinamis (*Admin Splash* atau *Default Splash*).
6. **Redirect**: Setelah simulasi *loading* animasi (1.5 detik), pengguna otomatis diarahkan ke `http://localhost:5173/dashboard`.
7. **Dashboard**: Mengambil data Profile ke API `/api/v1/auth/me` menggunakan token JWT (hasil di-_cache_ di Global Store Svelte 5 sehingga navigasi antar menu instan).
8. **Role UI & Sidebar**: 
   - Jika pengguna adalah Admin, Admin melihat menu **Manajemen Dosen** dan **Mata Kuliah** di *Sidebar*.
   - Admin dapat menambahkan Dosen dan Mata Kuliah melalui rute tersendiri yang menampilkan *form registration*.
   - Admin melihat rekap Daftar Dosen dan Daftar Mata Kuliah di halaman depan (Dashboard).
   - Jika pengguna adalah Dosen, halaman akan menampilkan antarmuka bersih (Empty State).
9. **Clean UI**: Semua instruksi di dalam file program bersifat bersih (telah di-*strip* dari semua komentar developer) sehingga *codebase* sangat *clean*.

---

## 📜 Lisensi
Aplikasi ini dikembangkan sebagai contoh penerapan arsitektur *Modular Monolith* tingkat lanjut dipadukan dengan standar *web development* yang estetik dan cepat.
