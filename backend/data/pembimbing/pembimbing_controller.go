package pembimbing

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetAllPembimbingController(c *gin.Context) {
	// 1. Panggil fungsi model
	result, err := GetAllPembimbingModel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data pembimbing",
			"success": false,
		})
		return
	}

	// 2. Beri Respon Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"total":   len(result),
	})
}

func GetPengajuanPembimbingController(c *gin.Context) {
	// 1. Panggil fungsi model
	result, err := GetPengajuanPembimbingModel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil pengajuan data pembimbing",
			"success": false,
		})
		return
	}

	// 2. Beri Respon Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"total":   len(result),
	})
}

func GetPengajuanOrtugController(c *gin.Context) {
	// 1. Panggil fungsi model
	result, err := GetPengajuanOrtuModel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil pengajuan data orangtua",
			"success": false,
		})
		return
	}

	// 2. Beri Respon Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
		"total":   len(result),
	})
}

func UbahStatusPembimbingController(c *gin.Context) {
	var req RequestUbahStatus

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
	err := UbahStatusPembimbingService(req)
	if err != nil {
		// Cek jika error karena data tidak ditemukan (404) atau bad request (400)
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "Status tidak valid") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 3. Response Sukses (Menyesuaikan pesan kode lama)
	msg := "Pembimbing berhasil disetujui"
	if req.Status == "tolak" {
		msg = "Data pembimbing berhasil dihapus"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": msg,
	})
}

func UbahStatusOrangtuaController(c *gin.Context) {
	var req RequestUbahStatus

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
	err := UbahStatusOrangtuaService(req)
	if err != nil {
		// Cek jika error karena data tidak ditemukan (404) atau bad request (400)
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "Status tidak valid") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 3. Response Sukses (Menyesuaikan pesan kode lama)
	msg := "Orangtua berhasil disetujui"
	if req.Status == "tolak" {
		msg = "Data Orangtua berhasil dihapus"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": msg,
	})
}

type UpdatePembimbingRequest struct {
	NamaLengkap *string `json:"nama_lengkap"`
	NoTelp      *string `json:"no_telp"`
	Alamat      *string `json:"alamat"`
	Status      *bool   `json:"status"`
}

func UpdatePembimbingController(c *gin.Context) {
	idUser := c.Param("id")

	// 1. Validasi format ID
	id_u, err := strconv.Atoi(idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID harus berupa angka"})
		return
	}

	if id_u < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID tidak boleh negatif"})
		return
	}

	// 2. Binding JSON
	var req UpdatePembimbingRequest
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
	err = UpdatePembimbingService(uint(id_u), req)
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

	// 4. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pembimbing berhasil diperbarui",
	})
}

func DeletePembimbingController(c *gin.Context) {
	id := c.Param("id")

	// 1. Validasi format ID
	id_u, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id harus berupa angka"})
		return
	}

	if id_u < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id tidak boleh negatif"})
		return
	}

	// 2. Panggil Service
	err = DeletePembimbingService(uint(id_u))
	if err != nil {
		// Tentukan status code berdasarkan pesan error
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "Akses ditolak") {
			status = http.StatusForbidden
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
		"message": "Data pembimbing berhasil dihapus",
	})
}
