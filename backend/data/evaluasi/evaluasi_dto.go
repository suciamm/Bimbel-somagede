package evaluasi

type CreateEvaluasiRequest struct {
	IDMurid           uint   `json:"id_murid" binding:"required"`
	IDPembimbing      uint   `json:"id_pembimbing" binding:"required"`
	EvaluasiKe        uint8  `json:"evaluasi_ke" binding:"required"`
	Nilai             string `json:"nilai" binding:"required"`
	CatatanPembimbing string `json:"catatan_pembimbing"`
	TanggalEvaluasi   string `json:"tanggal_evaluasi" binding:"required"` // YYYY-MM-DD
}

type EvaluasiPerMuridResponse struct {
	IDMurid            uint    `json:"id_murid"`
	KodeMurid          string  `json:"kode_murid"`
	NamaMurid          string  `json:"nama_murid"`
	IDPembimbing       *uint   `json:"id_pembimbing,omitempty"`
	NamaPembimbing     *string `json:"nama_pembimbing,omitempty"`
	Evaluasi1Nilai     *string `json:"evaluasi_1_nilai,omitempty"`
	Evaluasi1Catatan   *string `json:"evaluasi_1_catatan,omitempty"`
	Evaluasi1Tanggal   *string `json:"evaluasi_1_tanggal,omitempty"`
	Evaluasi2Nilai     *string `json:"evaluasi_2_nilai,omitempty"`
	Evaluasi2Catatan   *string `json:"evaluasi_2_catatan,omitempty"`
	Evaluasi2Tanggal   *string `json:"evaluasi_2_tanggal,omitempty"`
	Evaluasi3Nilai     *string `json:"evaluasi_3_nilai,omitempty"`
	Evaluasi3Catatan   *string `json:"evaluasi_3_catatan,omitempty"`
	Evaluasi3Tanggal   *string `json:"evaluasi_3_tanggal,omitempty"`
	TahapBerikutnya    uint8   `json:"tahap_berikutnya"`
	BolehInputTahap2   bool    `json:"boleh_input_tahap_2"`
	BolehInputTahap3   bool    `json:"boleh_input_tahap_3"`
	StatusProgressText string  `json:"status_progress_text"`
}
