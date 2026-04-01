package payment

type TambahPaketRequest struct {
	NamaPaket   string  `json:"nama_paket"`
	Harga       float64 `json:"harga"`
	DurasiHari  int     `json:"durasi_hari"`
	DurasiBulan int     `json:"durasi_bulan"`
	Deskripsi   string  `json:"deskripsi"`
}

type EditPaketRequest struct {
	NamaPaket   string  `json:"nama_paket"`
	Harga       float64 `json:"harga"`
	DurasiHari  int     `json:"durasi_hari"`
	DurasiBulan int     `json:"durasi_bulan"`
	Deskripsi   string  `json:"deskripsi"`
	Status      string  `json:"status"`
}

type TambahTransaksiRequest struct {
	IDMurid        uint    `json:"id_murid"`
	IDPaket        uint    `json:"id_paket"`
	TanggalMulai   string  `json:"tanggal_mulai"`
	TanggalSelesai string  `json:"tanggal_selesai"`
	TanggalBayar   string  `json:"tanggal_bayar"`
	JumlahBayar    float64 `json:"jumlah_bayar"`
	MetodeBayar    string  `json:"metode_bayar"`
}

type TambahPaketBimbinganInput struct {
	NamaPaket   string  `json:"nama_paket"`
	Harga       float64 `json:"harga"`
	DurasiBulan int     `json:"durasi_bulan"`
	Deskripsi   string  `json:"deskripsi"`
}

type PaketBimbingan struct {
	IdPaket     int     `json:"id_paket"`
	NamaPaket   string  `json:"nama_paket"`
	Harga       float64 `json:"harga"`
	DurasiBulan int     `json:"durasi_bulan"`
	DurasiHari  int     `json:"durasi_hari"`
	Deskripsi   string  `json:"deskripsi"`
	Status      string  `json:"status"`
}

type MidtransRequest struct {
	IDPaket     int    `json:"id_paket"`
	IDMurid     int    `json:"id_murid"`
	JumlahBayar int64  `json:"jumlah_bayar"`
	NamaPaket   string `json:"nama_paket"`
	NamaUser    string `json:"nama_user"`
}

type PerpanjangRequest struct {
	IDMurid   uint   `json:"id_murid"`
	NamaUser  string `json:"nama_user"`
	NamaPaket string `json:"nama_paket"`
}

type KonfirmasiPerpanjangRequest struct {
	IDMurid     uint    `json:"id_murid"`
	OrderID     string  `json:"order_id"`
	JumlahBayar float64 `json:"jumlah_bayar"`
}

type MidtransResponse struct {
	Token   string
	OrderID string
}

type PerpanjangMidtransResponse struct {
	Token   string
	OrderID string
	Harga   float64
}

type TransaksiPembayaran struct {
	IdTransaksi uint    `json:"id_transaksi"`
	IdPaket     uint    `json:"id_paket"`
	IdMurid     uint    `json:"id_murid"`
	KodeMurid   string  `json:"kode_murid"`
	NamaMurid   string  `json:"nama_murid"`
	NamaPaket   string  `json:"nama_paket"`
	TglMulai    string  `json:"tgl_mulai"`
	TglSelesai  string  `json:"tgl_selesai"`
	JumlahBayar float64 `json:"jumlah_bayar"`
	MetodeBayar string  `json:"metode_bayar"`
	Status      string  `json:"status"`
}
