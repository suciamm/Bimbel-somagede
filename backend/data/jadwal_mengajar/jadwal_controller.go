package jadwalmengajar

import (
	"data/config"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetAllJadwalController(c *gin.Context) {
	// Panggil Model langsung
	data, err := GetAllJadwalModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data jadwal aktif",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetJadwalByPembimbingController(c *gin.Context) {
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
	data, err := GetJadwalByPembimbingModel(id_u, config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data jadwal pembimbing",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func CreateJadwalController(c *gin.Context) {
	var req JadwalCreateRequest

	// 1. Validasi JSON & Required Fields (Binding)
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
	err := CreateJadwalService(req)
	if err != nil {
		// Bedakan status code: 400 untuk validasi bisnis, 500 untuk database
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak valid") ||
			strings.Contains(err.Error(), "wajib") ||
			strings.Contains(err.Error(), "ditemukan") ||
			strings.Contains(err.Error(), "sudah memiliki") ||
			strings.Contains(err.Error(), "bentrok") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 3. Response Sukses
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Jadwal berhasil ditambahkan",
	})
}

func UpdateJadwalController(c *gin.Context) {
	id := c.Param("id") // Sesuaikan dengan nama parameter di router kamu

	// 1. Validasi ID
	idJadwal, err := strconv.Atoi(id)
	if err != nil || idJadwal < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID Jadwal tidak valid"})
		return
	}

	var req JadwalUpdateRequest

	// 2. Binding JSON
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

	// 3. Panggil Service
	err = UpdateJadwalService(uint(idJadwal), req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "valid") || strings.Contains(err.Error(), "kosong") || strings.Contains(err.Error(), "sudah memiliki") || strings.Contains(err.Error(), "bentrok") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 4. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jadwal berhasil diperbarui",
	})
}

func DeleteJadwalController(c *gin.Context) {
	id := c.Param("id_jadwal")

	// 1. Validasi format ID
	id_jadwal, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id harus berupa angka"})
		return
	}

	if id_jadwal < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id tidak boleh negatif"})
		return
	}

	// 2. Panggil Service
	err = DeleteJadwalService(uint(id_jadwal))
	if err != nil {
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
		"message": "Jadwal berhasil dihapus",
	})
}
