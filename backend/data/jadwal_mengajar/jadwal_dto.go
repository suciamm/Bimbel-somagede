package jadwalmengajar

type JadwalResponse struct {
	IDJadwal      uint   `json:"id_jadwal"`
	HariBimbingan string `json:"hari_bimbingan"`
	WaktuMulai    string `json:"waktu_mulai"`
	WaktuSelesai  string `json:"waktu_selesai"`
	Ruangan       string `json:"ruangan"`

	KodePembimbing string `json:"kode_pembimbing"`
	NamaPembimbing string `json:"nama_pembimbing"`

	KodeMurid   string `json:"kode_murid"`
	NamaMurid string `json:"nama_murid"`
}

type JadwalCreateRequest struct {
	IDPembimbing  uint   `json:"id_pembimbing"`
	IDMurid       uint   `json:"id_murid"`
	HariBimbingan string `json:"hari_bimbingan"`
	WaktuMulai    string `json:"waktu_mulai"`
	WaktuSelesai  string `json:"waktu_selesai"`
	Ruangan       string `json:"ruangan"`
}

type JadwalUpdateRequest struct {
	IDPembimbing  uint   `json:"id_pembimbing" binding:"required"`
	IDMurid       uint   `json:"id_murid" binding:"required"`
	HariBimbingan string `json:"hari_bimbingan" binding:"required"`
	WaktuMulai    string `json:"waktu_mulai" binding:"required"`
	WaktuSelesai  string `json:"waktu_selesai" binding:"required"`
	Ruangan       string `json:"ruangan" binding:"required"`
}
