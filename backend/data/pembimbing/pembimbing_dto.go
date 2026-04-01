package pembimbing

import "time"

// Struct untuk menampung hasil query (tanpa password)
type PembimbingResponse struct {
	IDUser         uint   `json:"id_user"`
	KodePembimbing string `json:"kode_pembimbing"`
	Username       string `json:"username"`
	NamaLengkap    string `json:"nama_lengkap"`
	NoTelp         string `json:"no_telp"`
	Status         bool   `json:"status"`
}

type User struct {
	IDUser             uint
	Username           string
	Password           string
	NamaLengkap        string
	Role               string
	Alamat             string
	NoTelp             string
	Status             bool
	KodePembimbing     string
	TanggalPendaftaran time.Time
}

type PengajuanPembimbingResponse struct {
	IdUser      uint   `json:"id_user"`
	Username    string `json:"username"`
	NamaLengkap string `json:"nama_lengkap"`
	NoTelpon    string `json:"no_telp"`
	Alamat      string `json:"alamat"`
	Status      string `json:"status"`
}

type RequestUbahStatus struct {
	IdUser uint   `json:"id_user"`
	Status string `json:"status"` // 'tolak' atau 'setujui'
}
