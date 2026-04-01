package absensi

import (
	"data/config"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type AbsensiInput struct {
	IDMurid      uint   `json:"id_murid"`
	IDPembimbing uint   `json:"id_pembimbing"`
	TanggalSesi  string `json:"tanggal_sesi"` // Terima sebagai string "YYYY-MM-DD"
	StatusHadir  string `json:"status_hadir"` // enum: 'hadir', 'izin', 'alpa'
	Keterangan   string `json:"keterangan"`
}

// CreateAbsensi mencatat kehadiran murid dalam satu sesi
func CreateAbsensi(w http.ResponseWriter, r *http.Request) {
	var input AbsensiInput

	// 1. Parsing Body Request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Input data tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Validasi dan Parsing Tanggal
	tanggalSesiParsed, err := time.Parse("2006-01-02", input.TanggalSesi)
	if err != nil {
		http.Error(w, "Format tanggal sesi tidak valid. Gunakan YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	// 3. Persiapkan Data Absensi
	absensiBaru := config.Absensi{
		IDMurid:      input.IDMurid,
		IDPembimbing: input.IDPembimbing,
		TanggalSesi:  tanggalSesiParsed,
		StatusHadir:  input.StatusHadir,
		Keterangan:   input.Keterangan,
	}

	// 4. Simpan ke config
	// Opsional: Cek apakah absensi untuk murid tersebut di tanggal yang sama sudah ada
	var existingAbsensi config.Absensi
	check := config.DB.Where("id_murid = ? AND tanggal_sesi = ?", input.IDMurid, tanggalSesiParsed).First(&existingAbsensi)
	if check.RowsAffected > 0 {
		http.Error(w, "Absensi untuk murid ini di tanggal tersebut sudah dicatat.", http.StatusConflict)
		return
	}

	result := config.DB.Create(&absensiBaru)
	if result.Error != nil {
		http.Error(w, "Gagal menyimpan data absensi.", http.StatusInternalServerError)
		return
	}

	// 5. Beri Respon Sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Absensi berhasil dicatat!",
		"data":    absensiBaru,
	})
}

func GetRekapAbsensi(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter query seperti 'bulan' dan 'tahun'
	// Contoh: /api/v1/absensi/rekap?tahun=2025&bulan=12

	bulan := r.URL.Query().Get("bulan")
	tahun := r.URL.Query().Get("tahun")

	// Default filter jika tidak ada parameter
	if tahun == "" {
		tahun = time.Now().Format("2006")
	}

	var absensiList []config.Absensi
	query := config.DB.Preload("Murid").Preload("Pembimbing")

	// Filter berdasarkan bulan dan tahun (jika bulan diisi)
	if bulan != "" {
		// Contoh filter MySQL untuk bulan: MONTH(tanggal_sesi) = ? AND YEAR(tanggal_sesi) = ?
		query = query.Where("MONTH(tanggal_sesi) = ? AND YEAR(tanggal_sesi) = ?", bulan, tahun)
	} else {
		// Hanya filter berdasarkan tahun
		query = query.Where("YEAR(tanggal_sesi) = ?", tahun)
	}

	result := query.Find(&absensiList)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Gagal mengambil rekap absensi.", http.StatusInternalServerError)
		return
	}

	// Beri Respon Sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Rekap absensi berhasil diambil",
		"data":    absensiList,
		"total":   len(absensiList),
	})
}

func UpdateAbsensi(w http.ResponseWriter, r *http.Request) {
	idAbsensiStr := chi.URLParam(r, "id")
	idAbsensi, err := strconv.ParseUint(idAbsensiStr, 10, 32)
	if err != nil {
		http.Error(w, "ID Absensi tidak valid", http.StatusBadRequest)
		return
	}

	var existingAbsensi config.Absensi
	// 1. Cari Absensi yang akan diupdate
	if err := config.DB.First(&existingAbsensi, idAbsensi).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Data absensi tidak ditemukan untuk diupdate", http.StatusNotFound)
			return
		}
		http.Error(w, "Gagal mencari data absensi", http.StatusInternalServerError)
		return
	}

	var input struct {
		StatusHadir string `json:"status_hadir"`
		Keterangan  string `json:"keterangan"`
		// IDMurid dan TanggalSesi tidak boleh diubah melalui endpoint update ini
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Input data update tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Siapkan data update
	updateData := map[string]interface{}{}

	if input.StatusHadir != "" {
		updateData["status_hadir"] = input.StatusHadir
	}
	if input.Keterangan != "" {
		updateData["keterangan"] = input.Keterangan
	}

	// 3. Lakukan update
	result := config.DB.Model(&existingAbsensi).Updates(updateData)

	if result.Error != nil {
		http.Error(w, "Gagal memperbarui data absensi", http.StatusInternalServerError)
		return
	}

	// 4. Beri Respon Sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Absensi berhasil diperbarui"})
}

func DeleteAbsensi(w http.ResponseWriter, r *http.Request) {
	idAbsensiStr := chi.URLParam(r, "id")
	idAbsensi, err := strconv.ParseUint(idAbsensiStr, 10, 32)
	if err != nil {
		http.Error(w, "ID Absensi tidak valid", http.StatusBadRequest)
		return
	}

	// 1. Lakukan Hard Delete
	result := config.DB.Delete(&config.Absensi{}, idAbsensi)

	if result.Error != nil {
		http.Error(w, "Gagal menghapus absensi", http.StatusInternalServerError)
		return
	}

	// 2. Cek RowsAffected
	if result.RowsAffected == 0 {
		http.Error(w, "Absensi tidak ditemukan untuk dihapus", http.StatusNotFound)
		return
	}

	// 3. Beri Respon Sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Absensi berhasil dihapus"})
}
