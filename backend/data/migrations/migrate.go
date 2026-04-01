package migrations

import (
	"data/config"
	"fmt"
)

func Migrate() {
	db := config.DB
	if db == nil {
		panic("❌ Database belum diinisialisasi — panggil config.ConnectDB() dulu!")
	}

	// err := db.AutoMigrate()

	err := db.AutoMigrate(
		&config.User{},  // parent
		&config.Murid{}, // child
		&config.Jadwal{},
		&config.Absensi{},
		&config.Materi{},
		&config.PerkembanganMurid{},
		&config.EvaluasiMurid{},
		&config.PaketBimbingan{},
		&config.PaketBimbel{},
		&config.Langganan{},
		&config.Pembayaran{},
	)
	if err != nil {
		panic(fmt.Sprintf("Gagal migrasi: %v", err))
	}

	fmt.Println("✅ Migrasi database selesai.")
}
