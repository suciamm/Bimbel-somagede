package payment

import (
	"data/config"
	"errors"

	"gorm.io/gorm"
)

func TambahPaketBimbinganModel(data config.PaketBimbingan) error {
	if err := config.DB.Create(&data).Error; err != nil {
		return errors.New("Gagal menyimpan data paket bimbingan ke database")
	}
	return nil
}

func GetPaketBimbinganModel(db *gorm.DB) ([]PaketBimbingan, error) {
	var paketList []PaketBimbingan

	err := config.DB.
		Table("paket_bimbingan").
		Select(`
			id_paket,
			nama_paket,
			harga,
			durasi_bulan,
			durasi_hari,
			deskripsi,
			status
		`).
		Order("id_paket ASC").
		Scan(&paketList).Error

	return paketList, err
}

func GetTransaksiPembayaranModel(db *gorm.DB) ([]TransaksiPembayaran, error) {
	var transaksiList []TransaksiPembayaran

	err := config.DB.
		Table("langganan tp").
		Select(`
			tp.id_langganan AS id_transaksi,
			mp.kode_murid,
			mp.nama_murid,
			pb.nama_paket,
			tp.tgl_mulai,
			tp.tgl_selesai,
			pmb.jumlah_bayar,
			pmb.metode_bayar,
			pmb.tanggal_bayar,
			pmb.status
		`).
		Joins("JOIN murid mp ON tp.id_murid = mp.id_murid").
		Joins("JOIN paket_bimbingan pb ON tp.id_paket = pb.id_paket").
		Joins("JOIN pembayaran pmb ON pmb.id_langganan = tp.id_langganan").
		Order("tp.tgl_mulai DESC").
		Scan(&transaksiList).Error

	return transaksiList, err

}

func GetPaketBimbinganAktifModel(db *gorm.DB) ([]PaketBimbingan, error) {
	var paketList []PaketBimbingan
	err := config.DB.
		Table("paket_bimbingan").
		Select(`
			id_paket,
			nama_paket,
			harga,
			durasi_bulan,
			durasi_hari,
			deskripsi,
			status
		`).
		Where("status = ?", "aktif").
		Order("id_paket ASC").
		Scan(&paketList).Error

	return paketList, err
}

func EditPaketBimbinganModel(paket config.PaketBimbingan) error {
	if err := config.DB.Save(&paket).Error; err != nil {
		return errors.New("Gagal memperbarui data paket bimbingan")
	}
	return nil
}

func DeletePaketBimbinganModel(paket config.PaketBimbingan) error {
	// Eksekusi delete berdasarkan objek paket yang ditemukan
	if err := config.DB.Delete(&paket).Error; err != nil {
		return errors.New("Gagal menghapus data paket bimbingan")
	}
	return nil
}

func HapusTransaksiPembayaranModel(transaksi config.Langganan) error {
	//Hapus transaksi pembayaran
	if err := config.DB.Delete(&transaksi).Error; err != nil {
		return errors.New("Gagal menghapus data transaksi pembayaran")
	}
	return nil
}

func TambahTransaksiModel(langganan config.Langganan, pembayaran config.Pembayaran) error {
	// Memulai Transaksi Database
	tx := config.DB.Begin()

	// 1. Simpan ke tabel Langganan
	if err := tx.Create(&langganan).Error; err != nil {
		tx.Rollback()
		return errors.New("Gagal menyimpan data transaksi pembayaran (langganan)")
	}

	// 2. Set IDLangganan ke objek pembayaran (mengambil ID yang baru digenerate)
	pembayaran.IDLangganan = langganan.IDLangganan

	// 3. Simpan ke tabel Pembayaran
	if err := tx.Create(&pembayaran).Error; err != nil {
		tx.Rollback()
		return errors.New("Gagal menyimpan data transaksi pembayaran (pembayaran)")
	}

	// Jika semua oke, commit
	return tx.Commit().Error
}

func GetTransaksiByMuridModel(id_user int, db *gorm.DB) ([]TransaksiPembayaran, error) {
	var transaksiList []TransaksiPembayaran

	err := db.
		Table("langganan tp").
		Select(`
			tp.id_langganan AS id_transaksi,
			pb.id_paket,
			mp.id_murid,
			mp.kode_murid,
			mp.nama_murid,
			pb.nama_paket,
			tp.tgl_mulai,
			tp.tgl_selesai,
			COALESCE(pmb_sum.total_jumlah_bayar, 0) AS jumlah_bayar,
			pmb.metode_bayar,
			pmb.tanggal_bayar,
			pmb.status
		`).
		Joins("JOIN murid mp ON tp.id_murid = mp.id_murid").
		Joins("JOIN paket_bimbingan pb ON tp.id_paket = pb.id_paket").
		Joins(`JOIN (
			SELECT id_langganan, MAX(id_pembayaran) AS latest_id_pembayaran
			FROM pembayaran
			GROUP BY id_langganan
		) pmb_latest ON pmb_latest.id_langganan = tp.id_langganan`).
		Joins(`LEFT JOIN (
			SELECT id_langganan, SUM(jumlah_bayar) AS total_jumlah_bayar
			FROM pembayaran
			GROUP BY id_langganan
		) pmb_sum ON pmb_sum.id_langganan = tp.id_langganan`).
		Joins("JOIN pembayaran pmb ON pmb.id_pembayaran = pmb_latest.latest_id_pembayaran").
		Where("mp.id_user = ?", id_user).
		Order("tp.tgl_mulai DESC").
		Scan(&transaksiList).Error

	return transaksiList, err
}

// Mencari paket berdasarkan ID
func GetPaketByIdModel(id uint) (config.PaketBimbingan, error) {
	var paket config.PaketBimbingan
	err := config.DB.First(&paket, id).Error
	return paket, err
}

// Mencari data langganan murid di paket tertentu
func GetLanggananByMuridPaketModel(idMurid, idPaket uint) (config.Langganan, error) {
	var langganan config.Langganan
	err := config.DB.Where("id_murid = ? AND id_paket = ?", idMurid, idPaket).First(&langganan).Error
	return langganan, err
}
