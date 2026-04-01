package jadwalmengajar

import (
	"data/config"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetAllJadwalModel(db *gorm.DB) ([]JadwalResponse, error) {
	var jadwalList []JadwalResponse

	err := config.DB.
		Table("jadwal j").
		Select(`
			j.id_jadwal,
			j.hari_bimbingan,
			j.waktu_mulai,
			j.waktu_selesai,
			j.ruangan,

			u.kode_pembimbing,
			u.nama_lengkap AS nama_pembimbing,

			m.kode_murid,
			m.nama_murid
		`).
		Joins("JOIN users u ON j.id_pembimbing = u.id_user").
		Joins("JOIN murid m ON j.id_murid = m.id_murid").
		Order("j.hari_bimbingan ASC, j.waktu_mulai ASC").
		Scan(&jadwalList).Error

	return jadwalList, err

}

func GetJadwalByPembimbingModel(idUser int, db *gorm.DB) ([]JadwalResponse, error) {
	var jadwalList []JadwalResponse

	err := config.DB.
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
		Table("jadwal j").
		Select(`
			j.id_jadwal,
			j.hari_bimbingan,
			j.waktu_mulai,
			j.waktu_selesai,
			j.ruangan,

			u.kode_pembimbing,
			u.nama_lengkap AS nama_pembimbing,

			m.kode_murid,
			m.nama_murid
		`).
		Joins("JOIN users u ON j.id_pembimbing = u.id_user").
		Joins("JOIN murid m ON j.id_murid = m.id_murid").
		Where("j.id_pembimbing = ?", idUser).
		Order("j.hari_bimbingan ASC, j.waktu_mulai ASC").
		Scan(&jadwalList).Error

	return jadwalList, err

}

func CreateJadwalModel(data config.Jadwal) error {
	if err := config.DB.Create(&data).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("Jadwal dengan data yang sama sudah ada")
		}
		return errors.New("Gagal menambahkan jadwal ke database")
	}
	return nil
}

func UpdateJadwalModel(jadwal config.Jadwal) error {
	// Save akan melakukan update pada semua field berdasarkan Primary Key (IDJadwal)
	if err := config.DB.Save(&jadwal).Error; err != nil {
		return errors.New("Gagal mengupdate jadwal ke database")
	}
	return nil
}

func DeleteJadwalModel(id_jadwal uint) error {
	// Menggunakan Unscoped() jika ingin benar-benar menghapus permanen (Hard Delete)
	result := config.DB.Unscoped().Delete(&config.Jadwal{}, id_jadwal)

	if result.Error != nil {
		return errors.New("Gagal menghapus jadwal")
	}

	// Tambahan validasi jika RowsAffected kosong (antisipasi double hit)
	if result.RowsAffected == 0 {
		return errors.New("Jadwal tidak ditemukan")
	}

	return nil
}

func GetJadwalByPembimbing(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID Pembimbing dari URL parameter
	params := mux.Vars(r)
	idUser, err := strconv.ParseUint(params["id_user"], 10, 32)
	if err != nil {
		http.Error(w, "ID Pembimbing tidak valid", http.StatusBadRequest)
		return
	}

	var jadwalList []JadwalResponse

	result := config.DB.
		Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
		Table("jadwal j").
		Select(`
			j.id_jadwal,
			j.hari_bimbingan,
			j.waktu_mulai,
			j.waktu_selesai,
			j.ruangan,

			u.kode_pembimbing,
			u.nama_lengkap AS nama_pembimbing,

			m.kode_murid,
			m.nama_murid
		`).
		Joins("JOIN users u ON j.id_pembimbing = u.id_user").
		Joins("JOIN murid m ON j.id_murid = m.id_murid").
		Where("j.id_pembimbing = ?", idUser).
		Order("j.hari_bimbingan ASC, j.waktu_mulai ASC").
		Scan(&jadwalList)

	if result.Error != nil {
		http.Error(w, "Gagal mengambil data jadwal", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "Data jadwal berhasil diambil",
		"total":   len(jadwalList),
		"data":    jadwalList,
	})
}
