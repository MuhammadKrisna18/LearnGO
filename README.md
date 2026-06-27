# Modular Monolith System

Sebuah sistem terintegrasi dengan arsitektur **Modular Monolith**, dikembangkan menggunakan **Go (Fiber)** di sisi backend dan **SvelteKit (Svelte 5 Runes)** di sisi frontend. Proyek ini dirancang untuk memberikan performa yang tinggi, struktur kode yang bersih, serta desain antarmuka pengguna (UI) modern berbasis *Glassmorphism*.

## 🚀 Fitur Utama

- **Modular Architecture**: Struktur backend dibagi berdasarkan modul (contoh: modul `auth` dan `matakuliah`) yang membuat kode lebih terisolasi dan mudah di-maintain.
- **Frontend Modular Design**: Pemisahan antarmuka (UI) menjadi komponen-komponen terisolasi di `src/lib/components` dipadukan dengan manajemen *Global State* Svelte 5 Runes untuk optimasi reaktivitas dan caching data.
- **Modern Authentication**: Implementasi JWT (JSON Web Token) dengan kapabilitas *Role-based Access Control* (RBAC).
- **Custom Performance Logger**: Middleware khusus untuk mengukur latensi eksekusi setiap request API, dikategorikan secara *real-time* ke dalam tiga level: `FAST`, `MODERATE`, dan `SLOW`.
- **Role-based Dashboards**: Tampilan dashboard spesifik untuk masing-masing role dengan tata letak (layout) grid 2-kolom yang elegan. Sidebar menggunakan antarmuka *Glassmorphism* modern dengan ikon vektor (SVG) profesional.
- **Admin Management**: Fitur khusus bagi admin untuk mendaftarkan akun baru (seperti akun dosen) serta mengelola **Permintaan Ganti Email** (Request Email).
- **Manajemen Kelas Dinamis**: Mendukung pembuatan jadwal kelas (*Multi-schedule*) yang fleksibel. Satu ruangan kelas fisik dapat digunakan untuk hari dan jam yang berbeda dengan sistem proteksi *conflict* (jadwal ganda) terintegrasi, serta tabel UI canggih yang secara otomatis dikelompokkan berdasarkan Program Studi, lengkap dengan fitur pencarian dan pengurutan (jam/hari/nama kelas).
- **Sistem Profil & Keamanan**: Halaman profil terdedikasi (`/dashboard/profile`). Data vital seperti Nomor Induk Dosen (NID) bersifat permanen (tidak bisa diubah/dihapus), dan penggantian email memerlukan _approval_ Admin.
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
   - Dashboard utama (`/dashboard`) hanya menampilkan informasi *overview* bersih. Untuk Admin, halaman ini digunakan untuk meninjau **Permintaan Ganti Email**.
   - Navigasi khusus Admin tersedia di *Sidebar* elegan (ikon vektor, tanpa list berantakan).
   - Pengaturan Profil terpisah secara rapi di halaman khusus (`/dashboard/profile`) yang dapat diakses dengan mengklik Avatar di pojok kanan atas.
   - Manajemen Dosen (`/dashboard/dosen`), Manajemen Mata Kuliah (`/dashboard/matakuliah`), dan Manajemen Kelas (`/dashboard/kelas`) menggunakan sistem Grid dua kolom, menyatukan formulir pendaftaran dan tabel daftar data dalam satu antarmuka yang bersih.
   - Pada halaman **Daftar Kelas**, data kelas secara cerdas dikelompokkan per Program Studi dan memiliki dukungan kolom pencarian manual dan opsi filter pengurutan (*Sorting* berdasarkan Hari, Jam, dan Kelas).
   - Jika pengguna adalah Dosen, halaman dashboard utama menampilkan antarmuka bersih (Empty State).
9. **Manajemen Akun Dosen**:
   - Dosen yang baru terdaftar memiliki NID (Nomor Induk Dosen) otomatis sepanjang 5 digit yang bersifat **permanen**.
   - Dosen yang ingin mengganti email harus mengirim permintaan melalui sistem (*Email Request System*) ke Admin dengan memasukkan *username* saja (sistem otomatis menambahkan `@DosenGO.id`).
10. **Clean UI**: Semua instruksi di dalam file program bersifat bersih (telah di-*strip* dari semua komentar developer) sehingga *codebase* sangat *clean*.

---

## 📜 Lisensi
Aplikasi ini dikembangkan sebagai contoh penerapan arsitektur *Modular Monolith* tingkat lanjut dipadukan dengan standar *web development* yang estetik dan cepat.
