package payment

import (
	"data/config"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TambahPaketBimbinganController(c *gin.Context) {
	var req TambahPaketRequest

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
	err := TambahPaketBimbinganService(req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "wajib diisi") {
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
		"message": "Paket bimbingan berhasil ditambahkan",
	})
}

func GetPaketBimbinganController(c *gin.Context) {
	// Panggil Model langsung
	data, err := GetPaketBimbinganModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data paket bimbingan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetTransaksiPembayaranController(c *gin.Context) {
	// Panggil Model langsung
	data, err := GetTransaksiPembayaranModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data transaksi pembayaran",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func GetPaketBimbinganAktifController(c *gin.Context) {
	// Panggil Model langsung
	data, err := GetPaketBimbinganAktifModel(config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data paket bimbingan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})
}

func EditPaketBimbinganController(c *gin.Context) {
	id := c.Param("id_paket")

	// 1. Validasi ID
	idPaket, err := strconv.Atoi(id)
	if err != nil || idPaket < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID paket tidak valid"})
		return
	}

	var req EditPaketRequest

	// 2. Binding JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Input data tidak valid"})
		return
	}

	// 3. Panggil Service
	err = EditPaketBimbinganService(uint(idPaket), req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 4. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Paket bimbingan berhasil diperbarui",
	})
}

func DeletePaketBimbinganController(c *gin.Context) {
	id := c.Param("id_paket")

	// 1. Validasi format ID
	idPaket, err := strconv.Atoi(id)
	if err != nil || idPaket < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID paket tidak valid"})
		return
	}

	if idPaket < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id tidak boleh negatif"})
		return
	}

	// 2. Panggil Service
	err = DeletePaketBimbinganService(uint(idPaket))
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak ditemukan") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "Tidak dapat menghapus") {
			// Conflict error: ada jadwal mengajar terkait
			status = http.StatusConflict
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
		"message": "Paket bimbingan berhasil dihapus",
	})
}

func HapusTransaksiPembayaranController(c *gin.Context) {
	id := c.Param("id_transaksi")

	// 1. Validasi format ID
	IDTransaksi, err := strconv.Atoi(id)
	if err != nil || IDTransaksi < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID paket tidak valid"})
		return
	}

	if IDTransaksi < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id tidak boleh negatif"})
		return
	}

	// 2. Panggil Service
	err = HapusTransaksiPembayaranService(uint(IDTransaksi))
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
		"message": "Transaksi pembayaran berhasil dihapus",
	})
}

func CreateMidtransTransactionController(c *gin.Context) {
	var req MidtransRequest

	// 1. Bind & Validasi JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = getErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  out,
				"message": "Validasi gagal",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	// 2. Panggil Service
	resp, err := CreateMidtransTransactionService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 3. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"snap_token": resp.Token,
		"order_id":   resp.OrderID,
	})
}

func PerpanjangPaketMidtransController(c *gin.Context) {
	// Ambil param id_paket
	idPaketParam := c.Param("id_paket")
	idPaket, err := strconv.ParseUint(idPaketParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID paket tidak valid",
		})
		return
	}

	var req PerpanjangRequest

	// 1. Bind & Validasi JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = getErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  out,
				"message": "Validasi gagal",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	// 2. Panggil Service
	resp, err := PerpanjangPaketMidtransService(uint(idPaket), req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "tidak") ||
			strings.Contains(err.Error(), "tersedia") ||
			strings.Contains(err.Error(), "ditemukan") {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 3. Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"snap_token": resp.Token,
		"order_id":   resp.OrderID,
		"harga":      resp.Harga,
	})
}

func KonfirmasiPerpanjangMidtransController(c *gin.Context) {
	idPaketParam := c.Param("id_paket")
	idPaket, err := strconv.ParseUint(idPaketParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID paket tidak valid",
		})
		return
	}

	var req KonfirmasiPerpanjangRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Format JSON tidak valid",
		})
		return
	}

	err = KonfirmasiPerpanjangMidtransService(uint(idPaket), req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "wajib diisi") || strings.Contains(err.Error(), "tidak valid") {
			status = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "tidak ditemukan") || strings.Contains(err.Error(), "tidak tersedia") {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Perpanjangan paket berhasil diproses",
	})
}

func TambahTransaksiController(c *gin.Context) {
	var req TambahTransaksiRequest

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
	err := TambahTransaksiService(req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "wajib diisi") || strings.Contains(err.Error(), "Format tanggal") {
			status = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "masih memiliki paket bimbingan aktif") {
			status = http.StatusConflict
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
		"message": "Transaksi pembayaran berhasil ditambahkan",
	})
}

func GetTransaksiByMuridController(c *gin.Context) {

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
	data, err := GetTransaksiByMuridModel(id_u, config.DB)

	// Cek apakah ada error di database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data transaksi pembayaran",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"total":   len(data),
	})

}

func MidtransCallbackController(c *gin.Context) {
	log.Println("midtrans-callback-hit")

	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid payload",
		})
		return
	}

	// Panggil Service (error hanya di-log)
	if err := MidtransCallbackService(notification); err != nil {
		log.Println("CALLBACK SERVICE ERROR:", err)
		// ⚠️ tetap return 200 agar Midtrans tidak retry
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
