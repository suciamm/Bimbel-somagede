package payment

import (
	"data/config"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "oneof":
		return "Nilai tidak valid, pilih salah satu: " + fe.Param()
	case "numeric":
		return "Harus berupa angka"
	}
	return "Format tidak valid"

}

func TambahPaketBimbinganService(data TambahPaketRequest) error {

	// 1. Validasi field wajib (Business Logic)
	if data.NamaPaket == "" || data.Harga < 0 || data.DurasiHari < 0 || data.DurasiBulan <= 0 {
		return errors.New("Nama paket, harga, durasi hari, dan durasi bulan wajib diisi dengan nilai yang valid")
	}

	// 2. Siapkan data paket bimbingan baru
	paketBaru := config.PaketBimbingan{
		NamaPaket:   data.NamaPaket,
		Harga:       data.Harga,
		DurasiHari:  data.DurasiHari,
		DurasiBulan: data.DurasiBulan,
		Deskripsi:   data.Deskripsi,
	}

	// 3. Panggil Model
	return TambahPaketBimbinganModel(paketBaru)
}

func EditPaketBimbinganService(id uint, data EditPaketRequest) error {

	// 1. Cari paket bimbingan berdasarkan ID
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, id).Error; err != nil {
		return errors.New("Paket bimbingan tidak ditemukan")
	}

	// 2. Update data objek (Logic sesuai kode lama)
	paket.NamaPaket = data.NamaPaket
	paket.Harga = data.Harga
	paket.DurasiHari = data.DurasiHari
	paket.DurasiBulan = data.DurasiBulan
	paket.Deskripsi = data.Deskripsi
	paket.Status = data.Status

	// 3. Panggil Model
	return EditPaketBimbinganModel(paket)
}

func DeletePaketBimbinganService(id uint) error {
	// 1. Cari paket bimbingan berdasarkan ID untuk memastikan datanya ada
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, id).Error; err != nil {
		return errors.New("Paket bimbingan tidak ditemukan")
	}

	// 2. Check apakah ada jadwal mengajar untuk murid dengan langganan aktif paket ini
	var jadwalCount int64
	err := config.DB.Table("jadwal j").
		Joins("JOIN langganan l ON j.id_murid = l.id_murid").
		Where("l.id_paket = ?", id).
		Where("l.status IN ?", []string{"aktif", "perpanjang"}).
		Count(&jadwalCount).Error

	if err != nil {
		return errors.New("Gagal memeriksa jadwal mengajar")
	}

	if jadwalCount > 0 {
		return errors.New(fmt.Sprintf("Tidak dapat menghapus paket bimbingan ini karena masih ada %d jadwal mengajar untuk murid yang berlangganan paket ini. Silakan hapus atau ubah jadwal tersebut terlebih dahulu", jadwalCount))
	}

	// 3. Panggil Model untuk eksekusi hapus
	return DeletePaketBimbinganModel(paket)
}

func HapusTransaksiPembayaranService(id uint) error {
	// 1. Cari transaksi berdasarkan ID untuk memastikan datanya ada
	var transaksi config.Langganan
	if err := config.DB.First(&transaksi, id).Error; err != nil {
		return errors.New("Transaksi pembayaran tidak ditemukan")
	}
	// 2. Panggil Model untuk eksekusi hapus
	return HapusTransaksiPembayaranModel(transaksi)
}

func TambahTransaksiService(req TambahTransaksiRequest) error {
	// 1. Validasi field wajib
	if req.IDMurid == 0 || req.IDPaket == 0 || req.TanggalMulai == "" || req.TanggalSelesai == "" {
		return errors.New("ID Murid, ID Paket, Tanggal Mulai, dan Tanggal Selesai wajib diisi dengan nilai yang valid")
	}

	// 2. Parse Tanggal
	tglMulai, errM := time.Parse("2006-01-02", req.TanggalMulai)
	tglSelesai, errS := time.Parse("2006-01-02", req.TanggalSelesai)
	if errM != nil || errS != nil {
		return errors.New("Format tanggal tidak valid. Gunakan format YYYY-MM-DD")
	}

	if tglSelesai.Before(tglMulai) {
		return errors.New("Tanggal selesai tidak boleh lebih kecil dari tanggal mulai")
	}

	// Cegah tambah paket baru saat murid masih punya paket aktif pada periode yang overlap.
	var activeCount int64
	errActive := config.DB.Table("langganan").
		Where("id_murid = ?", req.IDMurid).
		Where("status IN ?", []string{"aktif", "perpanjang"}).
		Where("tgl_mulai <= ? AND tgl_selesai >= ?", tglSelesai, tglMulai).
		Count(&activeCount).Error
	if errActive != nil {
		return errors.New("Gagal memeriksa status paket aktif murid")
	}
	if activeCount > 0 {
		return errors.New("Murid masih memiliki paket bimbingan aktif. Tidak dapat menambah paket baru")
	}

	var tglBayar time.Time
	if req.TanggalBayar != "" {
		var errB error
		tglBayar, errB = time.Parse("2006-01-02", req.TanggalBayar)
		if errB != nil {
			return errors.New("Format tanggal bayar tidak valid. Gunakan format YYYY-MM-DD")
		}
	} else {
		tglBayar = time.Now()
	}

	// 3. Ambil data paket untuk hitung jumlah bayar otomatis
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, req.IDPaket).Error; err != nil {
		return errors.New("Paket bimbingan tidak ditemukan")
	}

	durasiPaketHari := (paket.DurasiBulan * 30) + paket.DurasiHari
	if durasiPaketHari <= 0 {
		return errors.New("Durasi paket tidak valid")
	}

	totalHari := int(tglSelesai.Sub(tglMulai).Hours()/24) + 1
	if totalHari <= 0 {
		totalHari = 1
	}

	kelipatanPaket := int(math.Ceil(float64(totalHari) / float64(durasiPaketHari)))
	if kelipatanPaket < 1 {
		kelipatanPaket = 1
	}

	jumlahBayar := int64(math.Round(paket.Harga * float64(kelipatanPaket)))

	// 4. Logika Penentuan Status Pembayaran
	statusPembayaran := "lunas"
	if req.MetodeBayar == "midtrans_pending" {
		statusPembayaran = "pending"
	}

	// 5. Siapkan Data untuk Model
	langganan := config.Langganan{
		IDMurid:    req.IDMurid,
		IDPaket:    req.IDPaket,
		TglMulai:   tglMulai,
		TglSelesai: tglSelesai,
		Status:     "aktif",
	}

	pembayaran := config.Pembayaran{
		TanggalBayar: tglBayar,
		JumlahBayar:  jumlahBayar,
		MetodeBayar:  req.MetodeBayar,
		Status:       statusPembayaran,
	}

	// 6. Kirim ke Model (Model akan menangani Transaction)
	return TambahTransaksiModel(langganan, pembayaran)
}

func CreateMidtransTransactionService(req MidtransRequest) (*MidtransResponse, error) {
	// 1. Validasi Paket
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, req.IDPaket).Error; err != nil {
		return nil, errors.New("Paket bimbingan tidak ditemukan")
	}

	if paket.Status != "aktif" {
		return nil, errors.New("Paket bimbingan tidak aktif")
	}

	// 2. Validasi Murid
	var murid config.Murid
	if err := config.DB.First(&murid, req.IDMurid).Error; err != nil {
		return nil, errors.New("Murid tidak ditemukan")
	}

	// 3. Init Midtrans Client
	var s snap.Client
	if midtrans.ServerKey == "" {
		return nil, errors.New("Midtrans Server Key belum dikonfigurasi")
	}
	s.New(midtrans.ServerKey, midtrans.Sandbox)

	// 4. Generate Order ID
	orderID := fmt.Sprintf("ORDER-%d", time.Now().Unix())

	// 3. Build Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: req.JumlahBayar,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.NamaUser,
		},
		CustomField1: "beli_baru",
		CustomField2: fmt.Sprintf("%d", req.IDMurid),
		CustomField3: fmt.Sprintf("%d", req.IDPaket),
		Items: &[]midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("PAKET-%d", req.IDPaket),
				Name:  req.NamaPaket,
				Price: req.JumlahBayar,
				Qty:   1,
			},
		},
	}

	// 5. Call Midtrans
	resp, err := s.CreateTransaction(snapReq)
	if err != nil {
		log.Println("MIDTRANS ERROR:", err)
		return nil, fmt.Errorf("Gagal membuat transaksi Midtrans: %v", err)
	}

	if resp.Token == "" {
		return nil, errors.New("Midtrans tidak mengembalikan token")
	}

	return &MidtransResponse{
		Token:   resp.Token,
		OrderID: orderID,
	}, nil
}

func PerpanjangPaketMidtransService(idPaket uint, req PerpanjangRequest) (*PerpanjangMidtransResponse, error) {
	// 1. Validasi Paket
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, idPaket).Error; err != nil {
		return nil, errors.New("Paket tidak ditemukan")
	}

	if paket.Status != "aktif" {
		return nil, errors.New("Paket tidak tersedia")
	}

	// 2. Validasi Langganan
	var langganan config.Langganan
	if err := config.DB.
		Where("id_murid = ? AND id_paket = ?", req.IDMurid, idPaket).
		First(&langganan).Error; err != nil {
		return nil, errors.New("Data langganan tidak ditemukan")
	}

	// 3. Init Midtrans
	var s snap.Client
	if midtrans.ServerKey == "" {
		return nil, errors.New("Midtrans Server Key belum dikonfigurasi")
	}
	s.New(midtrans.ServerKey, midtrans.Sandbox)

	// Order ID khusus perpanjang
	orderID := fmt.Sprintf("RENEW-%d-%d", req.IDMurid, time.Now().Unix())

	// 4. Build Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(paket.Harga),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.NamaUser,
		},
		CustomField1: "perpanjang",
		CustomField2: fmt.Sprintf("%d", req.IDMurid),
		CustomField3: fmt.Sprintf("%d", idPaket),
		Items: &[]midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("PAKET-%d", idPaket),
				Name:  "Perpanjang: " + paket.NamaPaket,
				Price: int64(paket.Harga),
				Qty:   1,
			},
		},
	}

	// 5. Create Transaksi
	resp, err := s.CreateTransaction(snapReq)
	if err != nil {
		return nil, errors.New("Gagal generate token Midtrans")
	}

	return &PerpanjangMidtransResponse{
		Token:   resp.Token,
		OrderID: orderID,
		Harga:   paket.Harga,
	}, nil
}

func KonfirmasiPerpanjangMidtransService(idPaket uint, req KonfirmasiPerpanjangRequest) error {
	if req.IDMurid == 0 {
		return errors.New("ID murid wajib diisi")
	}

	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, idPaket).Error; err != nil {
		return errors.New("Paket tidak ditemukan")
	}

	if paket.Status != "aktif" {
		return errors.New("Paket tidak tersedia")
	}

	jumlahBayar := req.JumlahBayar
	if jumlahBayar <= 0 {
		jumlahBayar = paket.Harga
	}

	metodeBayar := "midtrans"
	if req.OrderID != "" {
		metodeBayar = "midtrans:" + req.OrderID
	}

	return ProsesPerpanjangLangganan(req.IDMurid, idPaket, jumlahBayar, metodeBayar)
}

func MidtransCallbackService(notification map[string]interface{}) error {
	status, _ := notification["transaction_status"].(string)
	orderID, _ := notification["order_id"].(string)

	// Ambil custom field
	aksi, _ := notification["custom_field1"].(string)
	idMuridStr, _ := notification["custom_field2"].(string)
	idPaketStr, _ := notification["custom_field3"].(string)

	idMurid, _ := strconv.ParseUint(idMuridStr, 10, 32)
	idPaket, _ := strconv.ParseUint(idPaketStr, 10, 32)

	// Ambil nominal
	grossAmountStr, _ := notification["gross_amount"].(string)
	jumlahBayar, _ := strconv.ParseFloat(grossAmountStr, 64)

	// Hanya proses jika sukses
	if status != "settlement" && status != "capture" {
		log.Printf("INFO: Status %s untuk Order %s, tidak diproses", status, orderID)
		return nil
	}

	metodeBayar := "midtrans:" + orderID

	if aksi == "perpanjang" {
		log.Printf("INFO: PERPANJANG Order %s", orderID)
		return ProsesPerpanjangLangganan(
			uint(idMurid),
			uint(idPaket),
			jumlahBayar,
			metodeBayar,
		)
	}

	log.Printf("INFO: TAMBAH BARU Order %s", orderID)
	return ProsesTambahLanggananBaru(
		uint(idMurid),
		uint(idPaket),
		jumlahBayar,
		metodeBayar,
	)
}

func ProsesPerpanjangLangganan(
	idMurid uint,
	idPaket uint,
	jumlahBayar float64,
	metodeBayar string,
) error {

	// 1. Ambil Paket
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, idPaket).Error; err != nil {
		return err
	}
	if paket.Status != "aktif" {
		return errors.New("paket bimbingan tidak aktif")
	}

	// 2. Ambil Langganan
	var langganan config.Langganan
	if err := config.DB.
		Where("id_murid = ? AND id_paket = ?", idMurid, idPaket).
		First(&langganan).Error; err != nil {
		return err
	}

	// 3. Hitung Tanggal: pertahankan tanggal mulai, hanya geser tanggal selesai.
	oldEnd := langganan.TglSelesai
	newEnd := oldEnd.AddDate(0, paket.DurasiBulan, paket.DurasiHari)
	now := time.Now()
	langganan.TglPerpanjang = &now
	langganan.TglSelesai = newEnd
	langganan.Status = "aktif"

	// 4. Simpan Langganan
	if err := config.DB.Save(&langganan).Error; err != nil {
		return err
	}

	// 5. Simpan Pembayaran
	pembayaran := config.Pembayaran{
		IDLangganan:  langganan.IDLangganan,
		TanggalBayar: time.Now(),
		JumlahBayar:  int64(jumlahBayar),
		MetodeBayar:  metodeBayar,
		Status:       "lunas",
	}

	return config.DB.Create(&pembayaran).Error
}

func ProsesTambahLanggananBaru(
	idMurid uint,
	idPaket uint,
	jumlahBayar float64,
	metodeBayar string,
) error {

	// 1. Ambil Paket
	var paket config.PaketBimbingan
	if err := config.DB.First(&paket, idPaket).Error; err != nil {
		return err
	}

	// 2. Hitung Tanggal
	tglMulai := time.Now()
	tglSelesai := tglMulai.AddDate(0, paket.DurasiBulan, paket.DurasiHari)

	// 3. Simpan Langganan
	langganan := config.Langganan{
		IDMurid:    idMurid,
		IDPaket:    idPaket,
		TglMulai:   tglMulai,
		TglSelesai: tglSelesai,
		Status:     "aktif",
	}

	if err := config.DB.Create(&langganan).Error; err != nil {
		return err
	}

	// 4. Simpan Pembayaran
	pembayaran := config.Pembayaran{
		IDLangganan:  langganan.IDLangganan,
		TanggalBayar: time.Now(),
		JumlahBayar:  int64(jumlahBayar),
		MetodeBayar:  metodeBayar,
		Status:       "lunas",
	}

	return config.DB.Create(&pembayaran).Error
}
