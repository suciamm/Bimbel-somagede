package evaluasi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "oneof":
		return "Nilai tidak valid, pilih salah satu: " + fe.Param()
	case "numeric":
		return "Harus berupa angka"
	}
	return "Format tidak valid"
}

func CreateEvaluasiController(c *gin.Context) {
	var req CreateEvaluasiRequest

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

	if err := CreateEvaluasiService(req); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak valid") ||
			strings.Contains(err.Error(), "wajib") ||
			strings.Contains(err.Error(), "belum") ||
			strings.Contains(err.Error(), "tidak ditemukan") ||
			strings.Contains(err.Error(), "bukan") ||
			strings.Contains(err.Error(), "hanya boleh") ||
			strings.Contains(err.Error(), "sudah") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Evaluasi berhasil disimpan",
	})
}

func GetEvaluasiByPembimbingController(c *gin.Context) {
	idUser := c.Param("id_user")
	if idUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh kosong"})
		return
	}

	idPembimbing, err := strconv.Atoi(idUser)
	if err != nil || idPembimbing < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user harus berupa angka positif"})
		return
	}

	data, err := GetEvaluasiByPembimbingModel(idPembimbing, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data evaluasi pembimbing",
		})
		return
	}

	// DEBUG: log hasil query
	for i, d := range data {
		// log 1 baris saja jika banyak
		if i == 0 {
			fmt.Printf("[DEBUG] EvaluasiByPembimbing: %+v\n", d)
		}
	}

	data = EnrichProgressInfo(data)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetEvaluasiByOrtuController(c *gin.Context) {
	idUser := c.Param("id_user")
	if idUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user tidak boleh kosong"})
		return
	}

	idOrtu, err := strconv.Atoi(idUser)
	if err != nil || idOrtu < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id user harus berupa angka positif"})
		return
	}

	data, err := GetEvaluasiByOrtuModel(idOrtu, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data evaluasi untuk orangtua",
		})
		return
	}

	// DEBUG: log hasil query
	for i, d := range data {
		if i == 0 {
			fmt.Printf("[DEBUG] EvaluasiByOrtu: %+v\n", d)
		}
	}

	data = EnrichProgressInfo(data)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}
