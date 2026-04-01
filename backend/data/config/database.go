package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql" // Driver MySQL untuk GORM
	"gorm.io/gorm"         // Package utama GORM
	// Jangan lupa install package ini: go get gorm.io/gorm gorm.io/driver/mysql
)

// DB adalah koneksi global, TIPE-nya HARUS *gorm.DB
var DB *gorm.DB

// ConnectDB membuka koneksi ke database MySQL dan inisialisasi GORM
func ConnectDB() {
	var err error

	// Format DSN: username:password@tcp(host:port)/nama_db?charset=utf8mb4&parseTime=True&loc=Local
	// Pastikan nama_db kamu (misalnya: db_bimba_somagede) sudah dibuat
	dsn := "root:@tcp(127.0.0.1:3306)/bimbel?charset=utf8mb4&parseTime=True&loc=Local"

	// Menggunakan gorm.Open untuk koneksi
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		log.Fatalf("❌ Gagal konek ke database GORM: %v", err)
	}

	// Ping tidak diperlukan secara eksplisit di GORM, koneksi sudah diverifikasi oleh Open()
	fmt.Println("✅ Berhasil konek ke database GORM")

	// Opsional: Lakukan Migrasi Otomatis di sini
	// Kamu perlu import package models
	// db.AutoMigrate(&models.User{}, &models.Murid{}, ...)
}
