package users

import (
	"data/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser digunakan untuk mendaftarkan akun Admin atau Pembimbing baru.
func RegisterUser(c *gin.Context) {
	// 1. Definisikan struct untuk input (Binding)
	var input struct {
		Username    string `json:"username" binding:"required"`
		Password    string `json:"password" binding:"required"`
		NamaLengkap string `json:"nama_lengkap" binding:"required"`
		Role        string `json:"role" binding:"required"`
		NoTelp      string `json:"no_telp"`
		Alamat      string `json:"alamat"`
	}

	// 2. Parsing & Validasi Input
	// ShouldBindJSON akan otomatis cek 'required' tag di atas
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Data tidak lengkap atau format salah",
			"details": err.Error(),
		})
		return
	}

	// 3. Enkripsi Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	// 4. Mapping ke Model User
	newUser := config.User{
		Username:           input.Username,
		Password:           string(hashedPassword),
		NamaLengkap:        input.NamaLengkap,
		Role:               input.Role,
		NoTelp:             input.NoTelp,
		Alamat:             input.Alamat,
		Status:             false,      // Default sesuai kebutuhanmu
		TanggalPendaftaran: time.Now(), // Set waktu pendaftaran
	}

	// 5. Simpan ke Database
	if err := config.DB.Create(&newUser).Error; err != nil {
		// Cek jika error karena username duplikat
		c.JSON(http.StatusConflict, gin.H{
			"error": "Gagal membuat akun. Username mungkin sudah terdaftar.",
		})
		return
	}

	// 6. Respon Sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil!",
		"user":    newUser.Username,
	})
}

func LoginUser(c *gin.Context) {
	log.Println("1")
	// 1. Definisikan struct input
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 2. Parsing & Validasi Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan Password wajib diisi"})
		return
	}
	log.Println("2")

	var user config.User
	// 3. Cari User berdasarkan Username menggunakan GORM
	// First() akan mencari data pertama yang cocok
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		// Jika tidak ditemukan, GORM mengembalikan error 'record not found'
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password salah"})
		return
	}

	// 4. Verifikasi Password dengan bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// Password tidak cocok
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau Password salah"})
		return
	}

	log.Println("3")

	// 5. Beri Respon Sukses
	// Kamu tidak perlu manual Encode, cukup pakai c.JSON
	c.JSON(http.StatusOK, gin.H{
		"message":      "Login berhasil!",
		"id_user":      user.IDUser,
		"username":     user.Username,
		"role":         user.Role,
		"nama_lengkap": user.NamaLengkap,
		"status":       user.Status,
	})
}
