package murid

import "time"

type MuridListResponse struct {
	IDMurid     uint      `json:"id_murid"`
	NamaMurid   string    `json:"nama_murid"`
	NamaOrtu    string    `json:"nama_ortu"`
	TglMasuk    time.Time `json:"tgl_masuk"`
	Alamat      string    `json:"alamat"`
	StatusMurid string    `json:"status_murid"`
}

type MuridListByOrtuResponse struct {
	IDMurid        uint      `json:"id_murid"`
	NamaMurid      string    `json:"nama_murid"`
	NamaOrtu       string    `json:"nama_ortu"`
	TglMasuk       time.Time `json:"tgl_masuk"`
	HariBimbingan  string    `json:"hari_bimbingan"`
	WaktuMulai     string    `json:"waktu_mulai"`
	WaktuSelesai   string    `json:"waktu_selesai"`
	Ruangan        string    `json:"ruangan"`
	NamaPembimbing string    `json:"nama_pembimbing"`
	StatusMurid    string    `json:"status_murid"`
}

type MuridListByOrtu2Response struct {
	IDMurid        uint      `json:"id_murid"`
	KodeMurid      string    `json:"kode_murid"`
	NamaMurid      string    `json:"nama_murid"`
	NamaOrtu       string    `json:"nama_ortu"`
	TglMasuk       time.Time `json:"tgl_masuk"`
	HariBimbingan  string    `json:"hari_bimbingan"`
	WaktuMulai     string    `json:"waktu_mulai"`
	WaktuSelesai   string    `json:"waktu_selesai"`
	Ruangan        string    `json:"ruangan"`
	NamaPembimbing string    `json:"nama_pembimbing"`
	StatusMurid    string    `json:"status_murid"`
}

type MuridListAktifResponse struct {
	IDMurid     uint       `json:"id_murid"`
	KodeMurid   string     `json:"kode_murid"`
	NamaMurid   string     `json:"nama_murid"`
	NamaOrtu    string     `json:"nama_ortu"`
	TglLahir    time.Time  `json:"tgl_lahir"`
	TglMasuk    time.Time  `json:"tgl_masuk"`
	TglKeluar   *time.Time `json:"tgl_keluar"`
	Alamat      string     `json:"alamat"`
	StatusMurid string     `json:"status_murid"`
}

type OrtuListResponse struct {
	IDUser      uint   `json:"id_user"`
	NamaLengkap string `json:"nama_lengkap"`
	NoTelp      string `json:"no_telp"`
	Alamat      string `json:"alamat"`
	Status      string `json:"status"`
}

type RekapMuridDetailItem struct {
	KodeMurid string    `json:"kode_murid"`
	NamaMurid string    `json:"nama_murid"`
	TglMasuk  time.Time `json:"tgl_masuk"`
}

type RekapMuridBulananResponse struct {
	TanggalMulai      string                 `json:"tanggal_mulai"`
	TanggalSelesai    string                 `json:"tanggal_selesai"`
	JumlahMuridAktif  int64                  `json:"jumlah_murid_aktif"`
	MuridMasukBaru    int64                  `json:"murid_masuk_baru"`
	MuridKeluar       int64                  `json:"murid_keluar"`
	DaftarMuridAktif  []RekapMuridDetailItem `json:"daftar_murid_aktif" gorm:"-"`
	DaftarMuridMasuk  []RekapMuridDetailItem `json:"daftar_murid_masuk_baru" gorm:"-"`
	DaftarMuridKeluar []RekapMuridDetailItem `json:"daftar_murid_keluar" gorm:"-"`
}

type CreateMuridRequest struct {
	NamaMurid string `json:"nama_murid" binding:"required"`
	TglLahir  string `json:"tgl_lahir" binding:"required"` // YYYY-MM-DD
	IDUser    uint   `json:"id_user" binding:"required"`
	Alamat    string `json:"alamat" binding:"required"`
	TglMasuk  string `json:"tgl_masuk" binding:"required"` // YYYY-MM-DD
}

type UpdateMuridInput struct {
	NamaMurid   *string `json:"nama_murid"`
	TglLahir    *string `json:"tgl_lahir"`
	NamaOrtu    *string `json:"nama_ortu"`
	Alamat      *string `json:"alamat"`
	TglMasuk    *string `json:"tgl_masuk"`
	TglKeluar   *string `json:"tgl_keluar"`
	StatusMurid *string `json:"status_murid"`
}

type UpdateMuridRequest struct {
	NamaMurid   *string `json:"nama_murid"`
	TglLahir    *string `json:"tgl_lahir"`
	IDUser      *uint   `json:"id_user"`
	Alamat      *string `json:"alamat"`
	TglMasuk    *string `json:"tgl_masuk"`
	TglKeluar   *string `json:"tgl_keluar"`
	StatusMurid *string `json:"status_murid"`
}

type UpdateOrtuRequest struct {
	NamaLengkap *string `json:"nama_lengkap"`
	NoTelp      *string `json:"no_telp"`
	Alamat      *string `json:"alamat"`
	Status      *bool   `json:"status"`
}
