package config

import (
	"time"
)

// Ubah variabel ini untuk mengganti hari input evaluasi (default: Minggu)
var AllowedEvaluasiWeekday = time.Sunday
var AllowedEvaluasiWeekdayLabel = "Minggu"

// 1. Struct User (Tabel users)
type User struct {
	IDUser         uint   `gorm:"primaryKey;autoIncrement"`
	Username       string `gorm:"type:varchar(50);unique;not null"`
	Password       string `gorm:"type:varchar(255);not null"`
	NamaLengkap    string `gorm:"type:varchar(100);not null"`
	Role           string `gorm:"type:enum('admin','pembimbing','orangtua');not null"`
	Alamat         string `gorm:"type:text"`
	NoTelp         string `gorm:"type:varchar(15)"`
	Status         bool   `gorm:"default:false"`
	KodePembimbing string `gorm:"type:varchar(20);null"`

	TanggalPendaftaran time.Time `gorm:"type:datetime;null"`
}

func (User) TableName() string {
	return "users"
}

// 2. Struct Murid (Tabel murid)
type Murid struct {
	IDMurid     uint       `gorm:"primaryKey;autoIncrement"`
	KodeMurid   string     `gorm:"type:varchar(20);unique;not null"`
	NamaMurid   string     `gorm:"type:varchar(100);not null"`
	TglLahir    time.Time  `gorm:"type:date;not null"`
	Alamat      string     `gorm:"type:text"`
	TglMasuk    time.Time  `gorm:"type:date;not null"`
	TglKeluar   *time.Time `gorm:"type:date"`
	StatusMurid string     `gorm:"type:enum('aktif','lulus','keluar');default:'aktif';not null"`

	// FK ke users (orang tua)
	IDUser uint `gorm:"not null;column:id_user;index"`
	User   User `gorm:"-"`
}

func (Murid) TableName() string {
	return "murid"
}

// 3. Struct Jadwal (Tabel jadwal)
type Jadwal struct {
	IDJadwal      uint      `gorm:"primaryKey;autoIncrement" json:"id_jadwal"` // Primary Key
	IDPembimbing  uint      `gorm:"not null" json:"id_pembimbing"`             // Foreign Key ke User
	IDMurid       uint      `gorm:"not null" json:"id_murid"`                  // Foreign Key ke Murid
	HariBimbingan string    `gorm:"type:enum('senin','selasa','rabu','kamis','jumat','sabtu','minggu');not null" json:"hari_bimbingan"`
	WaktuMulai    time.Time `gorm:"type:time;not null" json:"waktu_mulai"`
	WaktuSelesai  time.Time `gorm:"type:time;not null" json:"waktu_selesai"`
	Ruangan       string    `gorm:"type:varchar(50)" json:"ruangan"`

	// Relasi (Opsional, untuk memudahkan query GORM)
	Pembimbing User  `gorm:"foreignKey:IDPembimbing"`
	Murid      Murid `gorm:"foreignKey:IDMurid"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}

// 4. Struct Absensi (Tabel absensi)
type Absensi struct {
	IDAbsensi    uint      `gorm:"primaryKey;autoIncrement" json:"id_absensi"` // Primary Key
	IDMurid      uint      `gorm:"not null" json:"id_murid"`
	IDPembimbing uint      `gorm:"not null" json:"id_pembimbing"`
	TanggalSesi  time.Time `gorm:"type:date;not null" json:"tanggal_sesi"`
	StatusHadir  string    `gorm:"type:enum('hadir','izin','alpa');not null" json:"status_hadir"`
	Keterangan   string    `gorm:"type:text" json:"keterangan"`

	// Relasi (Opsional)
	Murid      Murid `gorm:"foreignKey:IDMurid"`
	Pembimbing User  `gorm:"foreignKey:IDPembimbing"`
}

func (Absensi) TableName() string {
	return "absensi"
}

// 5. Struct Materi (Tabel materi)
type Materi struct {
	IDMateri    uint   `gorm:"primaryKey;autoIncrement" json:"id_materi"` // Primary Key
	JudulMateri string `gorm:"type:varchar(255);not null" json:"judul_materi"`
	Deskripsi   string `gorm:"type:text" json:"deskripsi"`
	Kategori    string `gorm:"type:varchar(50)" json:"kategori"`
	TipeFile    string `gorm:"type:varchar(20)" json:"tipe_file"`
	LokasiFile  string `gorm:"type:varchar(255);not null" json:"lokasi_file"`
	UploadedBy  uint   `gorm:"not null" json:"uploaded_by"`

	// Relasi (Opsional)
	Uploader User `gorm:"foreignKey:UploadedBy"`
}

func (Materi) TableName() string {
	return "materi"
}

// 6. Struct PerkembanganMurid (Tabel perkembangan_murid)
type PerkembanganMurid struct {
	IDPerkembangan    uint      `gorm:"primaryKey;autoIncrement" json:"id_perkembangan"` // Primary Key
	IDMurid           uint      `gorm:"not null" json:"id_murid"`
	IDPembimbing      uint      `gorm:"not null" json:"id_pembimbing"`
	TanggalCatat      time.Time `gorm:"type:date;not null" json:"tanggal_catat"`
	CatatanKualitatif string    `gorm:"type:text;not null" json:"catatan_kualitatif"`
	RekomendasiNext   string    `gorm:"type:text" json:"rekomendasi_next"`

	// Relasi (Opsional)
	Murid      Murid `gorm:"foreignKey:IDMurid"`
	Pembimbing User  `gorm:"foreignKey:IDPembimbing"`
}

func (PerkembanganMurid) TableName() string {
	return "perkembangan_murid"
}

type EvaluasiMurid struct {
	IDEvaluasi        uint      `gorm:"primaryKey;autoIncrement" json:"id_evaluasi"`
	IDMurid           uint      `gorm:"not null;index:idx_evaluasi_murid" json:"id_murid"`
	IDPembimbing      uint      `gorm:"not null;index:idx_evaluasi_pembimbing" json:"id_pembimbing"`
	EvaluasiKe        uint8     `gorm:"type:tinyint unsigned;not null;index:idx_evaluasi_ke" json:"evaluasi_ke"` // 1,2,3
	Nilai             string    `gorm:"type:enum('A','B','C','D','E');not null" json:"nilai"`
	CatatanPembimbing string    `gorm:"type:text" json:"catatan_pembimbing"`
	TanggalEvaluasi   time.Time `gorm:"type:date;not null;index:idx_tanggal_evaluasi" json:"tanggal_evaluasi"`

	Murid      Murid `gorm:"foreignKey:IDMurid;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Pembimbing User  `gorm:"foreignKey:IDPembimbing;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (EvaluasiMurid) TableName() string {
	return "evaluasi_murid"
}

type PaketBimbel struct {
	IDPaket    uint   `gorm:"primaryKey;autoIncrement"`
	NamaPaket  string `gorm:"type:varchar(100);not null"`
	Durasi     int    `gorm:"not null"` // dalam bulan
	Harga      int64  `gorm:"not null"`
	Keterangan string
}

func (PaketBimbel) TableName() string {
	return "paket_bimbel"
}

type Langganan struct {
	IDLangganan   uint       `gorm:"primaryKey;autoIncrement"`
	IDMurid       uint       `gorm:"not null"`
	IDPaket       uint       `gorm:"not null"`
	TglMulai      time.Time  `gorm:"type:date;not null"`
	TglPerpanjang *time.Time `gorm:"type:date;null"`
	TglSelesai    time.Time  `gorm:"type:date;not null"`
	Status        string     `gorm:"type:enum('aktif','habis','batal')"`
}

func (Langganan) TableName() string {
	return "langganan"
}

type Pembayaran struct {
	IDPembayaran uint      `gorm:"primaryKey;autoIncrement"`
	IDLangganan  uint      `gorm:"not null"`
	TanggalBayar time.Time `gorm:"type:date;not null"`
	JumlahBayar  int64     `gorm:"not null"`
	MetodeBayar  string    `gorm:"type:varchar(50);not null"`
	Status       string    `gorm:"type:enum('pending','lunas','gagal')"`
}

func (Pembayaran) TableName() string {
	return "pembayaran"
}

type PaketBimbingan struct {
	IdPaket     uint    `gorm:"primaryKey;autoIncrement"`
	NamaPaket   string  `gorm:"type:varchar(100);not null"`
	Harga       float64 `gorm:"not null"`
	DurasiHari  int     `gorm:"not null"` // dalam hari
	DurasiBulan int     `gorm:"not null"` // dalam bulan
	Deskripsi   string  `gorm:"type:text"`
	Status      string  `gorm:"type:enum('aktif','tidak aktif');default:'aktif'"`
}

func (PaketBimbingan) TableName() string {
	return "paket_bimbingan"
}


type EvalluasiPembimbing struct{

}