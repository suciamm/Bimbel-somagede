package evaluasi

import (
	"data/config"

	"gorm.io/gorm"
)

func CreateEvaluasiModel(data config.EvaluasiMurid, db *gorm.DB) error {
	if db == nil {
		db = config.DB
	}
	return db.Create(&data).Error
}

func GetEvaluasiByPembimbingModel(idPembimbing int, db *gorm.DB) ([]EvaluasiPerMuridResponse, error) {
	var result []EvaluasiPerMuridResponse
	if db == nil {
		db = config.DB
	}

	err := db.
		Table("jadwal j").
		Select(`
			m.id_murid,
			m.kode_murid,
			m.nama_murid,
			MAX(e.id_pembimbing) AS id_pembimbing,
			MAX(u.nama_lengkap) AS nama_pembimbing,
			MAX(CASE WHEN e.evaluasi_ke = 1 THEN e.nilai END) AS evaluasi_1_nilai,
			MAX(CASE WHEN e.evaluasi_ke = 1 THEN e.catatan_pembimbing END) AS evaluasi_1_catatan,
			MAX(CASE WHEN e.evaluasi_ke = 1 THEN DATE_FORMAT(e.tanggal_evaluasi, '%Y-%m-%d') END) AS evaluasi_1_tanggal,
			MAX(CASE WHEN e.evaluasi_ke = 2 THEN e.nilai END) AS evaluasi_2_nilai,
			MAX(CASE WHEN e.evaluasi_ke = 2 THEN e.catatan_pembimbing END) AS evaluasi_2_catatan,
			MAX(CASE WHEN e.evaluasi_ke = 2 THEN DATE_FORMAT(e.tanggal_evaluasi, '%Y-%m-%d') END) AS evaluasi_2_tanggal,
			MAX(CASE WHEN e.evaluasi_ke = 3 THEN e.nilai END) AS evaluasi_3_nilai,
			MAX(CASE WHEN e.evaluasi_ke = 3 THEN e.catatan_pembimbing END) AS evaluasi_3_catatan,
			MAX(CASE WHEN e.evaluasi_ke = 3 THEN DATE_FORMAT(e.tanggal_evaluasi, '%Y-%m-%d') END) AS evaluasi_3_tanggal
		`).
		Joins("JOIN murid m ON m.id_murid = j.id_murid").
		Joins("LEFT JOIN evaluasi_murid e ON e.id_murid = m.id_murid AND e.id_pembimbing = ?", idPembimbing).
		Joins("LEFT JOIN users u ON u.id_user = e.id_pembimbing").
		Where("j.id_pembimbing = ?", idPembimbing).
		Group("m.id_murid, m.kode_murid, m.nama_murid").
		Order("m.nama_murid ASC").
		Scan(&result).Error

	return result, err
}

func GetEvaluasiByOrtuModel(idOrtu int, db *gorm.DB) ([]EvaluasiPerMuridResponse, error) {
	var result []EvaluasiPerMuridResponse
	if db == nil {
		db = config.DB
	}

	err := db.
		Table("jadwal j").
		Select(`
			       m.id_murid,
			       m.kode_murid,
			       m.nama_murid,
			       j.id_pembimbing,
			       u.nama_lengkap AS nama_pembimbing,
			       MAX(CASE WHEN e.evaluasi_ke = 1 THEN e.nilai END) AS evaluasi_1_nilai,
			       MAX(CASE WHEN e.evaluasi_ke = 1 THEN e.catatan_pembimbing END) AS evaluasi_1_catatan,
			       MAX(CASE WHEN e.evaluasi_ke = 1 THEN DATE_FORMAT(e.tanggal_evaluasi, '%Y-%m-%d') END) AS evaluasi_1_tanggal,
			       MAX(CASE WHEN e.evaluasi_ke = 2 THEN e.nilai END) AS evaluasi_2_nilai,
			       MAX(CASE WHEN e.evaluasi_ke = 2 THEN e.catatan_pembimbing END) AS evaluasi_2_catatan,
			       MAX(CASE WHEN e.evaluasi_ke = 2 THEN DATE_FORMAT(e.tanggal_evaluasi, '%Y-%m-%d') END) AS evaluasi_2_tanggal,
			       MAX(CASE WHEN e.evaluasi_ke = 3 THEN e.nilai END) AS evaluasi_3_nilai,
			       MAX(CASE WHEN e.evaluasi_ke = 3 THEN e.catatan_pembimbing END) AS evaluasi_3_catatan,
			       MAX(CASE WHEN e.evaluasi_ke = 3 THEN DATE_FORMAT(e.tanggal_evaluasi, '%Y-%m-%d') END) AS evaluasi_3_tanggal
		       `).
		Joins("JOIN murid m ON m.id_murid = j.id_murid").
		Joins("LEFT JOIN evaluasi_murid e ON e.id_murid = m.id_murid AND e.id_pembimbing = j.id_pembimbing").
		Joins("LEFT JOIN users u ON u.id_user = j.id_pembimbing").
		Where("m.id_user = ?", idOrtu).
		Group("m.id_murid, m.kode_murid, m.nama_murid, j.id_pembimbing, u.nama_lengkap").
		Order("m.nama_murid ASC").
		Scan(&result).Error

	return result, err
}
