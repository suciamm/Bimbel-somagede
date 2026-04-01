package murid

import (
	"data/config"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetMuridAktifController(c *gin.Context) {
	// Panggil Model langsung
	data, err := GetMuridAktifModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data murid aktif",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetMuridByOrtuController(c *gin.Context) {
	// ambil parameter dari URL
	idUser := c.Param("id_user")
	if idUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh kosong"})
		return
	}

	// validasi harus angka
	id_u, err := strconv.Atoi(idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_u < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh negatif"})
		return
	}

	// Panggil Model langsung
	data, err := GetMuridByOrtuModel(id_u, config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data murid aktif",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetMuridTidakAktifController(c *gin.Context) {

	// Panggil Model langsung
	data, err := GetMuridTidakAktifModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data murid tidak aktif",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})

}

func GetRekapMuridBulananController(c *gin.Context) {
	layout := "2006-01-02"
	startDateParam := strings.TrimSpace(c.Query("tanggal_mulai"))
	endDateParam := strings.TrimSpace(c.Query("tanggal_selesai"))
	bulanParam := strings.TrimSpace(c.Query("bulan"))
	tahunParam := strings.TrimSpace(c.Query("tahun"))

	var startDate time.Time
	var endDate time.Time

	hasDateRange := startDateParam != "" || endDateParam != ""
	hasMonthYear := bulanParam != "" || tahunParam != ""

	switch {
	case hasDateRange:
		if startDateParam == "" || endDateParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tanggal_mulai dan tanggal_selesai wajib diisi bersamaan"})
			return
		}

		parsedStartDate, err := time.Parse(layout, startDateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format tanggal_mulai harus YYYY-MM-DD"})
			return
		}

		parsedEndDate, err := time.Parse(layout, endDateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format tanggal_selesai harus YYYY-MM-DD"})
			return
		}

		startDate = parsedStartDate
		endDate = parsedEndDate

	case hasMonthYear:
		if bulanParam == "" || tahunParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bulan dan tahun wajib diisi bersamaan"})
			return
		}

		bulan, err := strconv.Atoi(bulanParam)
		if err != nil || bulan < 1 || bulan > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bulan harus angka 1 sampai 12"})
			return
		}

		tahun, err := strconv.Atoi(tahunParam)
		if err != nil || tahun < 2000 || tahun > 2100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tahun harus angka valid"})
			return
		}

		startDate = time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, -1)

	default:
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, -1)
	}

	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tanggal_selesai tidak boleh lebih kecil dari tanggal_mulai"})
		return
	}

	data, err := GetRekapMuridBulananModel(startDate, endDate, config.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil rekap murid bulanan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// func GetDaftarOrtu(c *gin.Context) {
// 	var result []OrtuListResponse

// 	// 1. Query Database
// 	err := config.DB.
// 		Table("users").
// 		Select(`
//             users.id_user,
//             users.nama_lengkap,
//             users.no_telp,
//             users.alamat,
//             users.status AS status
//         `).
// 		Where("users.role = ?", "orangtua").
// 		Group("users.id_user").
// 		Order("users.nama_lengkap ASC").
// 		Scan(&result).Error

// 	// 2. Error Handling (Gin Style)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Gagal mengambil data orang tua",
// 		})
// 		return
// 	}

// 	// 3. Mapping Status (tetap sama, tapi lebih rapi)
// 	for i := range result {
// 		if result[i].Status != "0" {
// 			result[i].Status = "Aktif"
// 		} else {
// 			result[i].Status = "Tidak Aktif"
// 		}
// 	}

// 	// 4. Response Sukses
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Data orang tua berhasil diambil",
// 		"data":    result,
// 		"total":   len(result),
// 	})
// }

func GetDaftarOrtuController(c *gin.Context) {
	// Panggil Model langsung
	data, err := GetDaftarOrtuModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data orangtua",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetMuridByPembimbingController(c *gin.Context) {

	// 1. Ambil ID Pembimbing dari URL parameter (Gin style)
	idUser := c.Param("id_user")

	//validasi harus di isi
	if idUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh negatif"})
		return
	}

	// validasi harus angka
	id_u, err := strconv.Atoi(idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_u < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh negatif"})
		return
	}

	// Panggil Model langsung
	data, err := GetMuridByPembimbingModel(id_u, config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data murid",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})

}

func GetDataBimbinganByOrtuController(c *gin.Context) {

	// 1. Ambil ID Pembimbing dari URL parameter (Gin style)
	idUser := c.Param("id_user")

	//validasi harus di isi
	if idUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh negatif"})
		return
	}

	// validasi harus angka
	id_u, err := strconv.Atoi(idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_u < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh negatif"})
		return
	}

	// Panggil Model langsung
	data, err := GetDataBimbinganByOrtuModel(id_u, config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data bimbingan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})

}

func CreateMuridController(c *gin.Context) {
	var req CreateMuridRequest

	// 1. Validasi JSON & Required Fields
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = getErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": out, "message": "Validasi gagal"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format JSON tidak valid"})
		return
	}

	// 2. Panggil Service
	err := CreateMuridService(req)
	if err != nil {
		// Kita bisa bedakan status code berdasarkan pesan error jika perlu
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Data murid berhasil ditambahkan",
	})
}

func UpdateMuridController(c *gin.Context) {
	id := c.Param("id")

	//validasi harus di isi
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return

	}

	// validasi harus angka
	id_m, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_m < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	var req UpdateMuridRequest // Menggunakan struct eksternal seperti kemauanmu

	// 1. Validasi JSON & Required Fields menggunakan validator
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = getErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": out, "message": "Validasi gagal"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format JSON tidak valid"})
		return
	}

	// 2. Panggil Service
	err = UpdateMuridService(id_m, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 3. Response Sukses
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Data murid berhasil diperbarui",
	})
}

func UpdateOrtuController(c *gin.Context) {
	id := c.Param("id")

	//validasi harus di isi
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	// validasi harus angka
	id_ot, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_ot < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	var req UpdateOrtuRequest // Menggunakan struct eksternal agar rapi

	// 1. Validasi JSON & Binding
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = getErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": out, "message": "Validasi gagal"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Format JSON tidak valid"})
		return
	}

	// 2. Panggil Service
	err = UpdateOrtuService(id_ot, req)
	if err != nil {
		// Logika Status Code (jika tidak ditemukan beri 404, selain itu 500)
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 3. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data user berhasil diperbarui",
	})
}

func DeleteOrtuController(c *gin.Context) {
	id := c.Param("id")

	//validasi harus di isi
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	// validasi harus angka
	id_murid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_murid < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	// Panggil Service untuk pengecekan logika
	err = DeleteOrtuService(uint(id_murid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data orangtua berhasil dihapus",
	})
}

func DeleteMuridController(c *gin.Context) {
	id := c.Param("id")

	//validasi harus di isi
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	// validasi harus angka
	id_murid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id harus berupa angka"})
		return
	}

	// pastikan non-negative sebelum cast ke uint
	if id_murid < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak boleh negatif"})
		return
	}

	// Panggil Service untuk pengecekan logika
	err = DeleteMuridService(uint(id_murid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data murid berhasil dihapus",
	})
}
