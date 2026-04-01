# Vue Frontend - Bimbel Somagede

Frontend untuk backend Bimbel Somagede menggunakan Vue 3 + Vite.

## Tech Stack

- Vue 3
- Vue Router
- Pinia
- Axios
- Vite

## Struktur Utama

- `src/router`: konfigurasi route + proteksi route
- `src/stores`: state global (`auth`)
- `src/services`: API client dan endpoint service
- `src/views`: halaman per modul
- `src/layouts`: kerangka dashboard
- `src/components`: komponen UI reusable

## Setup

1. Install Node.js 18+ (disarankan Node 20 LTS)
2. Install dependency:

```bash
npm install
```

3. Copy env:

```bash
cp .env.example .env
```

4. Jalankan development server:

```bash
npm run dev
```

Aplikasi berjalan di `http://localhost:3000`.

## Integrasi Backend

Base URL API dibaca dari `VITE_API_BASE_URL`.
Contoh backend lokal: `http://localhost:8080`.

## Halaman yang Sudah Disiapkan

- Login
- Register
- Dashboard
- Data Murid
- Jadwal Mengajar
- Paket Bimbingan

## Catatan

Karena endpoint backend saat ini belum menggunakan JWT bearer token, autentikasi frontend menyimpan profil user di localStorage untuk session UI.
