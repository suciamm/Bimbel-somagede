-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Mar 27, 2026 at 05:53 AM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.0.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `bimbel`
--

-- --------------------------------------------------------

--
-- Table structure for table `absensi`
--

CREATE TABLE `absensi` (
  `id_absensi` bigint(20) UNSIGNED NOT NULL,
  `id_murid` bigint(20) UNSIGNED NOT NULL,
  `id_pembimbing` bigint(20) UNSIGNED NOT NULL,
  `tanggal_sesi` date NOT NULL,
  `status_hadir` enum('hadir','izin','alpa') NOT NULL,
  `keterangan` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `jadwal`
--

CREATE TABLE `jadwal` (
  `id_jadwal` bigint(20) UNSIGNED NOT NULL,
  `id_pembimbing` bigint(20) UNSIGNED NOT NULL,
  `id_murid` bigint(20) UNSIGNED NOT NULL,
  `hari_bimbingan` enum('senin','selasa','rabu','kamis','jumat','sabtu','minggu') NOT NULL,
  `waktu_mulai` datetime(3) NOT NULL,
  `waktu_selesai` datetime(3) NOT NULL,
  `ruangan` varchar(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `jadwal`
--

INSERT INTO `jadwal` (`id_jadwal`, `id_pembimbing`, `id_murid`, `hari_bimbingan`, `waktu_mulai`, `waktu_selesai`, `ruangan`) VALUES
(1, 6, 4, 'kamis', '2000-01-01 20:00:00.000', '2000-01-01 23:00:00.000', 're'),
(2, 4, 4, 'rabu', '2000-01-01 22:00:00.000', '2000-01-01 23:00:00.000', 'wdwdwdw'),
(3, 4, 4, 'minggu', '2000-01-01 15:00:00.000', '2000-01-01 16:00:00.000', 'sss'),
(4, 4, 4, 'kamis', '2000-01-01 15:00:00.000', '2000-01-01 16:00:00.000', 'cd'),
(5, 4, 5, 'kamis', '2000-01-01 15:00:00.000', '2000-01-01 16:00:00.000', 'mk'),
(6, 4, 6, 'kamis', '2000-01-01 15:00:00.000', '2000-01-01 16:00:00.000', 'mk'),
(7, 4, 5, 'kamis', '2000-01-01 15:00:00.000', '2000-01-01 16:00:00.000', 're');

-- --------------------------------------------------------

--
-- Table structure for table `langganan`
--

CREATE TABLE `langganan` (
  `id_langganan` bigint(20) UNSIGNED NOT NULL,
  `id_murid` bigint(20) UNSIGNED NOT NULL,
  `id_paket` bigint(20) UNSIGNED NOT NULL,
  `tgl_mulai` date NOT NULL,
  `tgl_perpanjang` date DEFAULT NULL,
  `tgl_selesai` date NOT NULL,
  `status` enum('aktif','habis','batal') DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `langganan`
--

INSERT INTO `langganan` (`id_langganan`, `id_murid`, `id_paket`, `tgl_mulai`, `tgl_perpanjang`, `tgl_selesai`, `status`) VALUES
(1, 4, 2, '2026-03-23', '2026-03-26', '2026-05-30', 'aktif'),
(2, 4, 3, '2026-03-10', NULL, '2026-04-03', 'aktif'),
(3, 4, 2, '2026-01-21', NULL, '2026-04-01', 'aktif'),
(4, 5, 2, '2026-01-01', '2026-03-26', '2026-06-28', 'aktif');

-- --------------------------------------------------------

--
-- Table structure for table `materi`
--

CREATE TABLE `materi` (
  `id_materi` bigint(20) UNSIGNED NOT NULL,
  `judul_materi` varchar(255) NOT NULL,
  `deskripsi` text DEFAULT NULL,
  `kategori` varchar(50) DEFAULT NULL,
  `tipe_file` varchar(20) DEFAULT NULL,
  `lokasi_file` varchar(255) NOT NULL,
  `uploaded_by` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `murid`
--

CREATE TABLE `murid` (
  `id_murid` bigint(20) UNSIGNED NOT NULL,
  `kode_murid` varchar(20) NOT NULL,
  `nama_murid` varchar(100) NOT NULL,
  `tgl_lahir` date NOT NULL,
  `alamat` text DEFAULT NULL,
  `tgl_masuk` date NOT NULL,
  `tgl_keluar` date DEFAULT NULL,
  `status_murid` enum('aktif','lulus','keluar') NOT NULL DEFAULT 'aktif',
  `id_user` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `murid`
--

INSERT INTO `murid` (`id_murid`, `kode_murid`, `nama_murid`, `tgl_lahir`, `alamat`, `tgl_masuk`, `tgl_keluar`, `status_murid`, `id_user`) VALUES
(1, 'MRD001', 'adam smith', '2026-03-01', 'bjynthbgrvfcd', '2026-03-16', NULL, 'lulus', 2),
(4, 'MRD002', 'gui', '2021-12-12', 'fwexs', '2026-01-31', NULL, 'aktif', 3),
(5, 'MRD003', 'sandi', '2023-07-13', 'defcrv', '2026-03-30', NULL, 'aktif', 3),
(6, 'MRD004', 'kris', '2023-10-09', 'cwervervrevre', '2026-03-16', NULL, 'aktif', 5);

-- --------------------------------------------------------

--
-- Table structure for table `paket_bimbel`
--

CREATE TABLE `paket_bimbel` (
  `id_paket` bigint(20) UNSIGNED NOT NULL,
  `nama_paket` varchar(100) NOT NULL,
  `durasi` bigint(20) NOT NULL,
  `harga` bigint(20) NOT NULL,
  `keterangan` longtext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `paket_bimbingan`
--

CREATE TABLE `paket_bimbingan` (
  `id_paket` bigint(20) UNSIGNED NOT NULL,
  `nama_paket` varchar(100) NOT NULL,
  `harga` double NOT NULL,
  `durasi_hari` bigint(20) NOT NULL,
  `durasi_bulan` bigint(20) NOT NULL,
  `deskripsi` text DEFAULT NULL,
  `status` enum('aktif','tidak aktif') DEFAULT 'aktif'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `paket_bimbingan`
--

INSERT INTO `paket_bimbingan` (`id_paket`, `nama_paket`, `harga`, `durasi_hari`, `durasi_bulan`, `deskripsi`, `status`) VALUES
(2, 'a', 125000, 0, 2, 'frevfdcsx', 'aktif'),
(3, 'bb', 100000, 1, 12, 'dewcecdc', 'aktif');

-- --------------------------------------------------------

--
-- Table structure for table `pembayaran`
--

CREATE TABLE `pembayaran` (
  `id_pembayaran` bigint(20) UNSIGNED NOT NULL,
  `id_langganan` bigint(20) UNSIGNED NOT NULL,
  `tanggal_bayar` date NOT NULL,
  `jumlah_bayar` bigint(20) NOT NULL,
  `metode_bayar` varchar(50) NOT NULL,
  `status` enum('pending','lunas','gagal') DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `pembayaran`
--

INSERT INTO `pembayaran` (`id_pembayaran`, `id_langganan`, `tanggal_bayar`, `jumlah_bayar`, `metode_bayar`, `status`) VALUES
(1, 1, '2026-03-31', 1000000, 'cash', 'lunas'),
(2, 2, '2026-03-10', 100000, 'midtrans', 'lunas'),
(3, 3, '2026-03-26', 250000, 'midtrans', 'lunas'),
(4, 1, '2026-03-26', 125000, 'midtrans:RENEW-4-1774514312', 'lunas'),
(5, 4, '2026-03-26', 125000, 'midtrans', 'lunas'),
(6, 4, '2026-03-26', 125000, 'midtrans:RENEW-5-1774518963', 'lunas'),
(7, 4, '2026-03-26', 125000, 'midtrans:RENEW-5-1774519905', 'lunas');

-- --------------------------------------------------------

--
-- Table structure for table `perkembangan_murid`
--

CREATE TABLE `perkembangan_murid` (
  `id_perkembangan` bigint(20) UNSIGNED NOT NULL,
  `id_murid` bigint(20) UNSIGNED NOT NULL,
  `id_pembimbing` bigint(20) UNSIGNED NOT NULL,
  `tanggal_catat` date NOT NULL,
  `catatan_kualitatif` text NOT NULL,
  `rekomendasi_next` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id_user` bigint(20) UNSIGNED NOT NULL,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `nama_lengkap` varchar(100) NOT NULL,
  `role` enum('admin','pembimbing','orangtua') NOT NULL,
  `alamat` text DEFAULT NULL,
  `no_telp` varchar(15) DEFAULT NULL,
  `status` tinyint(1) DEFAULT 0,
  `kode_pembimbing` varchar(20) DEFAULT NULL,
  `tanggal_pendaftaran` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id_user`, `username`, `password`, `nama_lengkap`, `role`, `alamat`, `no_telp`, `status`, `kode_pembimbing`, `tanggal_pendaftaran`) VALUES
(1, 'admin', '$2a$10$vGUdEsdV6tuvh0sy8UMTjORw3DY/W3Zhe1UI6qhhePDxB1hpBjhoS', 'admin bimbel', 'admin', 'jhjdxcdcd', '067283764', 0, '', '2026-03-25 23:10:19'),
(3, 'orangtua', '$2a$10$vMGe5SLnOFie07yeLWwZAu5njLs8Z0./.0UkGylBShiEY687N1lji', 'Suci', 'orangtua', 'veccfsdc', '7483764783', 1, '', '2026-03-26 01:09:52'),
(4, 'pbb', '$2a$10$/Ty78tYbfsjq5eRyEAg5ReG.D3mX1kTbIksiP7lC.xG3IAC6GjOtm', 'pembimbing one', 'pembimbing', '', '067382763', 1, 'PBM001', '2026-03-26 01:29:47'),
(5, 'orangtua2', '$2a$10$bguGTxS06ScHBPersj0.7uIAfQX4WdDh9I0pLjOOjG5RMr2LVeiWm', 'sitiamin', 'orangtua', 'uhgfregtyn', '8675676756', 1, '', '2026-03-26 13:55:29'),
(6, 'pbb2', '$2a$10$Gpl.JvpgmL0n5uub33H6IefzLEFFLI1h9NPp9vu.T7tgTnnbqLDOa', 'nafisa', 'pembimbing', 'vfdsaferw', '67564356', 1, 'PBM002', '2026-03-26 16:43:11');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `absensi`
--
ALTER TABLE `absensi`
  ADD PRIMARY KEY (`id_absensi`);

--
-- Indexes for table `jadwal`
--
ALTER TABLE `jadwal`
  ADD PRIMARY KEY (`id_jadwal`),
  ADD KEY `fk_id_pembimbing` (`id_pembimbing`),
  ADD KEY `fk_id_murid` (`id_murid`);

--
-- Indexes for table `langganan`
--
ALTER TABLE `langganan`
  ADD PRIMARY KEY (`id_langganan`);

--
-- Indexes for table `materi`
--
ALTER TABLE `materi`
  ADD PRIMARY KEY (`id_materi`);

--
-- Indexes for table `murid`
--
ALTER TABLE `murid`
  ADD PRIMARY KEY (`id_murid`),
  ADD UNIQUE KEY `uni_murid_kode_murid` (`kode_murid`),
  ADD KEY `idx_murid_id_user` (`id_user`);

--
-- Indexes for table `paket_bimbel`
--
ALTER TABLE `paket_bimbel`
  ADD PRIMARY KEY (`id_paket`);

--
-- Indexes for table `paket_bimbingan`
--
ALTER TABLE `paket_bimbingan`
  ADD PRIMARY KEY (`id_paket`);

--
-- Indexes for table `pembayaran`
--
ALTER TABLE `pembayaran`
  ADD PRIMARY KEY (`id_pembayaran`);

--
-- Indexes for table `perkembangan_murid`
--
ALTER TABLE `perkembangan_murid`
  ADD PRIMARY KEY (`id_perkembangan`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id_user`),
  ADD UNIQUE KEY `uni_users_username` (`username`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `absensi`
--
ALTER TABLE `absensi`
  MODIFY `id_absensi` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `jadwal`
--
ALTER TABLE `jadwal`
  MODIFY `id_jadwal` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `langganan`
--
ALTER TABLE `langganan`
  MODIFY `id_langganan` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `materi`
--
ALTER TABLE `materi`
  MODIFY `id_materi` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `murid`
--
ALTER TABLE `murid`
  MODIFY `id_murid` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `paket_bimbel`
--
ALTER TABLE `paket_bimbel`
  MODIFY `id_paket` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `paket_bimbingan`
--
ALTER TABLE `paket_bimbingan`
  MODIFY `id_paket` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `pembayaran`
--
ALTER TABLE `pembayaran`
  MODIFY `id_pembayaran` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `perkembangan_murid`
--
ALTER TABLE `perkembangan_murid`
  MODIFY `id_perkembangan` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id_user` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `jadwal`
--
ALTER TABLE `jadwal`
  ADD CONSTRAINT `fk_id_murid` FOREIGN KEY (`id_murid`) REFERENCES `murid` (`id_murid`),
  ADD CONSTRAINT `fk_id_pembimbing` FOREIGN KEY (`id_pembimbing`) REFERENCES `users` (`id_user`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
