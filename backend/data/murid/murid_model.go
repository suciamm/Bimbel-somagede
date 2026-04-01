package murid

import (
	"data/config"
	"errors"
	"time"

	"gorm.io/gorm"
)

type rekapNamaRow struct {
	KodeMurid string    `gorm:"column:kode_murid"`
	NamaMurid string    `gorm:"column:nama_murid"`
	TglMasuk  time.Time `gorm:"column:tgl_masuk"`
}

func GetMuridAktifModel(db *gorm.DB) ([]MuridListAktifResponse, error) {
	var result []MuridListAktifResponse

	err := config.DB.
		Table("murid").
		Select(`
			murid.id_murid,
			murid.kode_murid,
			murid.nama_murid,
			users.nama_lengkap AS nama_ortu,
			murid.tgl_lahir,
			murid.tgl_masuk,
			murid.tgl_keluar,
			users.alamat,
			murid.status_murid
		`).
		Joins("LEFT JOIN users ON users.id_user = murid.id_user").
		Where("murid.status_murid = ?", "aktif").
		Order("murid.tgl_masuk DESC").
		Scan(&result).Error

	return result, err
}

func GetMuridByOrtuModel(id_user int, db *gorm.DB) ([]MuridListByOrtu2Response, error) {
	var result []MuridListByOrtu2Response

	err := config.DB.
		Table("murid").
		Select(`
			murid.id_murid,
			murid.kode_murid,
			murid.nama_murid,
			users.nama_lengkap AS nama_ortu,
			murid.tgl_lahir,
			murid.tgl_masuk,
			murid.tgl_keluar,
			users.alamat,
			murid.status_murid
		`).
		Joins("LEFT JOIN users ON users.id_user = murid.id_user").
		Where("murid.id_user = ? AND murid.status_murid = ?", id_user, "aktif").
		Order("murid.tgl_masuk DESC").
		Scan(&result).Error

	return result, err
}

func GetMuridTidakAktifModel(db *gorm.DB) ([]MuridListAktifResponse, error) {
	var result []MuridListAktifResponse

	err := config.DB.
		Table("murid").
		Select(`
			murid.id_murid,
			murid.kode_murid,
			murid.nama_murid,
			users.nama_lengkap AS nama_ortu,
			murid.tgl_lahir,
			murid.tgl_masuk,
			murid.tgl_keluar,
			users.alamat,
			murid.status_murid
		`).
		Joins("LEFT JOIN users ON users.id_user = murid.id_user").
		Where("murid.status_murid != ?", "aktif").
		Order("murid.tgl_masuk DESC").
		Scan(&result).Error
	return result, err

}

func GetMuridByPembimbingModel(id_user int, db *gorm.DB) ([]MuridListAktifResponse, error) {
	var result []MuridListAktifResponse

	err := config.DB.
		Table("murid m").
		Select(`
            m.id_murid,
            m.kode_murid,
            m.nama_murid,
            ortu.nama_lengkap AS nama_ortu,
            m.tgl_lahir,
            m.tgl_masuk,
            m.tgl_keluar,
            ortu.alamat,
            m.status_murid
        `).
		Joins("JOIN jadwal j ON j.id_murid = m.id_murid").
		Joins("LEFT JOIN users ortu ON ortu.id_user = m.id_user").
		Where("j.id_pembimbing = ? AND m.status_murid = ?", id_user, "aktif").
		Group("m.id_murid, m.kode_murid, m.nama_murid, ortu.nama_lengkap, m.tgl_lahir, m.tgl_masuk, m.tgl_keluar, ortu.alamat, m.status_murid").
		Order("m.tgl_masuk DESC").
		Scan(&result).Error

	return result, err
}

func CreateMuridModel(data config.Murid) error {
	// Simpan ke database
	if err := config.DB.Create(&data).Error; err != nil {
		return errors.New("gagal menyimpan data ke database")
	}
	return nil
}

func UpdateMuridModel(id int, updateData map[string]interface{}) error {
	// Model(&config.Murid{}) memberitahu GORM tabel mana yang dipakai
	// Where("id_murid = ?", id) menentukan baris mana yang diupdate
	if err := config.DB.Model(&config.Murid{}).Where("id_murid = ?", id).Updates(updateData).Error; err != nil {
		return errors.New("gagal memperbarui database")
	}
	return nil
}

func UpdateOrtuModel(id int, updateData map[string]interface{}) error {
	// Updates menggunakan map akan mengabaikan field yang tidak ada di map
	if err := config.DB.Model(&config.User{}).Where("id_user = ?", id).Updates(updateData).Error; err != nil {
		return errors.New("Gagal memperbarui data user ke database")
	}
	return nil
}

// untuk menghitung jumlah NamaAkun yang terikat
func CheckIfMuridHasLanggananModel(db *gorm.DB, id uint) (int64, error) {
	var count int64
	// Kita cek ke tabel 'langganan' (pastikan nama tabelnya sesuai di DB)
	err := db.Table("langganan").Where("id_murid = ? AND status", id, "aktif").Count(&count).Error
	return count, err
}

func DeleteMuridModel(db *gorm.DB, id uint) error {
	// db.Delete akan menghapus record berdasarkan ID yang diberikan
	if err := db.Delete(&config.Murid{}, id).Error; err != nil {
		return err
	}
	return nil
}

func CheckIfOruHasMuridModel(db *gorm.DB, id uint) (int64, error) {
	var count int64
	// Kita cek ke tabel 'murid' (pastikan nama tabelnya sesuai di DB)
	err := db.Table("murid").Where("id_user = ?", id).Count(&count).Error
	return count, err
}

func DeleteOrtuModel(db *gorm.DB, id uint) error {
	// db.Delete akan menghapus record berdasarkan ID yang diberikan
	if err := db.Delete(&config.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func GetDataBimbinganByOrtuModel(id_user int, db *gorm.DB) ([]MuridListByOrtuResponse, error) {
	var result []MuridListByOrtuResponse
	err := db.
		Table("murid").
		Select(`
            murid.id_murid,
            murid.nama_murid,
            users.nama_lengkap AS nama_ortu,
            murid.tgl_masuk,
            jadwal.hari_bimbingan,
            jadwal.waktu_mulai,
            jadwal.waktu_selesai,
            jadwal.ruangan,
            pembimbing.nama_lengkap AS nama_pembimbing,
            murid.status_murid
        `).
		Joins("JOIN users ON users.id_user = murid.id_user").
		Joins("JOIN jadwal ON jadwal.id_murid = murid.id_murid").
		Joins("LEFT JOIN users AS pembimbing ON pembimbing.id_user = jadwal.id_pembimbing").
		Where("murid.id_user = ? AND users.role = ? AND murid.status_murid = ?", id_user, "orangtua", "aktif").
		Order("murid.nama_murid ASC, jadwal.hari_bimbingan ASC, jadwal.waktu_mulai ASC").
		Scan(&result).Error

	return result, err
}

func GetDaftarOrtuModel(db *gorm.DB) ([]OrtuListResponse, error) {
	var result []OrtuListResponse
	err := config.DB.
		Table("users").
		Select(`
            users.id_user,
            users.nama_lengkap,
            users.no_telp,
            users.alamat,
            users.status AS status
        `).
		Where("users.role = ? AND users.status = ?", "orangtua", "1").
		Group("users.id_user").
		Order("users.nama_lengkap ASC").
		Scan(&result).Error

	// Mapping Status (tetap sama, tapi lebih rapi)
	for i := range result {
		if result[i].Status != "0" {
			result[i].Status = "Aktif"
		} else {
			result[i].Status = "Tidak Aktif"
		}
	}

	return result, err

}

func GetRekapMuridBulananModel(startDate time.Time, endDate time.Time, db *gorm.DB) (RekapMuridBulananResponse, error) {
	var result RekapMuridBulananResponse
	muridBaruStartDate := endDate.AddDate(0, -1, 0)

	var muridAktifRows []rekapNamaRow
	var muridMasukRows []rekapNamaRow
	var muridKeluarRows []rekapNamaRow

	err := db.Table("murid").
		Select(`
			COALESCE(SUM(CASE
				WHEN murid.status_murid = 'aktif'
					AND murid.tgl_masuk <= ?
					AND (murid.tgl_keluar IS NULL OR murid.tgl_keluar > ?)
				THEN 1 ELSE 0 END), 0) AS jumlah_murid_aktif,
			COALESCE(SUM(CASE
				WHEN DATE(murid.tgl_masuk) BETWEEN ? AND ?
				THEN 1 ELSE 0 END), 0) AS murid_masuk_baru,
			COALESCE(SUM(CASE
				WHEN murid.tgl_keluar IS NOT NULL
					AND DATE(murid.tgl_keluar) BETWEEN ? AND ?
				THEN 1 ELSE 0 END), 0) AS murid_keluar
		`, endDate, endDate, muridBaruStartDate, endDate, startDate, endDate).
		Scan(&result).Error
	if err != nil {
		return result, err
	}

	err = db.Table("murid").
		Select("murid.kode_murid, murid.nama_murid, murid.tgl_masuk").
		Where("murid.status_murid = ? AND murid.tgl_masuk <= ? AND (murid.tgl_keluar IS NULL OR murid.tgl_keluar > ?)", "aktif", endDate, endDate).
		Order("murid.nama_murid ASC").
		Scan(&muridAktifRows).Error
	if err != nil {
		return result, err
	}

	err = db.Table("murid").
		Select("murid.kode_murid, murid.nama_murid, murid.tgl_masuk").
		Where("DATE(murid.tgl_masuk) BETWEEN ? AND ?", muridBaruStartDate, endDate).
		Order("murid.nama_murid ASC").
		Scan(&muridMasukRows).Error
	if err != nil {
		return result, err
	}

	err = db.Table("murid").
		Select("murid.kode_murid, murid.nama_murid, murid.tgl_masuk").
		Where("murid.tgl_keluar IS NOT NULL AND DATE(murid.tgl_keluar) BETWEEN ? AND ?", startDate, endDate).
		Order("murid.nama_murid ASC").
		Scan(&muridKeluarRows).Error
	if err != nil {
		return result, err
	}

	result.TanggalMulai = startDate.Format("2006-01-02")
	result.TanggalSelesai = endDate.Format("2006-01-02")
	result.DaftarMuridAktif = toRekapMuridDetailList(muridAktifRows)
	result.DaftarMuridMasuk = toRekapMuridDetailList(muridMasukRows)
	result.DaftarMuridKeluar = toRekapMuridDetailList(muridKeluarRows)

	return result, err
}

func toRekapMuridDetailList(rows []rekapNamaRow) []RekapMuridDetailItem {
	result := make([]RekapMuridDetailItem, 0, len(rows))
	for _, row := range rows {
		if row.NamaMurid == "" {
			continue
		}
		result = append(result, RekapMuridDetailItem{
			KodeMurid: row.KodeMurid,
			NamaMurid: row.NamaMurid,
			TglMasuk:  row.TglMasuk,
		})
	}
	return result
}
