package jadwalmengajar

import (
	"data/config"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Helper untuk menerjemahkan tag error
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

func CreateJadwalService(req JadwalCreateRequest) error {
	// 1. Validasi Hari
	hariValid := map[string]bool{
		"senin": true, "selasa": true, "rabu": true,
		"kamis": true, "jumat": true, "sabtu": true, "minggu": true,
	}
	if !hariValid[strings.ToLower(req.HariBimbingan)] {
		return errors.New("Hari bimbingan tidak valid")
	}

	// 2. Parsing Waktu (Mulai & Selesai)
	layoutTime := "15:04:05"
	loc, locErr := time.LoadLocation("Asia/Jakarta")
	if locErr != nil {
		loc = time.FixedZone("WIB", 7*60*60)
	}

	waktuMulaiParsed, err := time.ParseInLocation(layoutTime, req.WaktuMulai, loc)
	if err != nil {
		return errors.New("Format Waktu Mulai tidak valid. Gunakan HH:mm:ss")
	}

	waktuSelesaiParsed, err := time.ParseInLocation(layoutTime, req.WaktuSelesai, loc)
	if err != nil {
		return errors.New("Format Waktu Selesai tidak valid. Gunakan HH:mm:ss")
	}

	// Normalisasi tanggal agar bisa dibandingkan (hanya jam yang penting)
	waktuMulaiParsed = time.Date(2000, 1, 1, waktuMulaiParsed.Hour(), waktuMulaiParsed.Minute(), waktuMulaiParsed.Second(), 0, loc)
	waktuSelesaiParsed = time.Date(2000, 1, 1, waktuSelesaiParsed.Hour(), waktuSelesaiParsed.Minute(), waktuSelesaiParsed.Second(), 0, loc)

	if waktuSelesaiParsed.Before(waktuMulaiParsed) || waktuSelesaiParsed.Equal(waktuMulaiParsed) {
		return errors.New("Waktu selesai harus setelah waktu mulai")
	}

	roomName := strings.TrimSpace(req.Ruangan)
	if roomName == "" {
		return errors.New("Ruangan wajib diisi")
	}

	// 3. Validasi Keberadaan Pembimbing & Murid
	var pembimbing config.User
	if err := config.DB.First(&pembimbing, req.IDPembimbing).Error; err != nil {
		return errors.New("Pembimbing tidak ditemukan")
	}

	var murid config.Murid
	if err := config.DB.First(&murid, req.IDMurid).Error; err != nil {
		return errors.New("Murid tidak ditemukan")
	}

	// 4.5 Validasi: pembimbing yang sama tidak boleh punya jadwal bentrok di hari yang sama
	var countPembimbing int64
	err = config.DB.Model(&config.Jadwal{}).
		Where("id_pembimbing = ? AND hari_bimbingan = ?", req.IDPembimbing, strings.ToLower(req.HariBimbingan)).
		Where("? < waktu_selesai AND ? > waktu_mulai", waktuMulaiParsed, waktuSelesaiParsed).
		Count(&countPembimbing).Error
	if err != nil {
		return errors.New("Gagal memeriksa jadwal pembimbing")
	}
	if countPembimbing > 0 {
		return errors.New("Jadwal pembimbing bentrok di hari dan jam yang sama")
	}

	// 4.6 Validasi: ruangan yang sama tidak boleh dipakai pembimbing lain di jam overlap pada hari yang sama
	var countRuangan int64
	err = config.DB.Model(&config.Jadwal{}).
		Where("LOWER(TRIM(ruangan)) = LOWER(TRIM(?)) AND hari_bimbingan = ? AND id_pembimbing <> ?", roomName, strings.ToLower(req.HariBimbingan), req.IDPembimbing).
		Where("? < waktu_selesai AND ? > waktu_mulai", waktuMulaiParsed, waktuSelesaiParsed).
		Count(&countRuangan).Error
	if err != nil {
		return errors.New("Gagal memeriksa bentrok ruangan")
	}
	if countRuangan > 0 {
		return errors.New("Ruangan sudah dipakai pembimbing lain pada hari dan jam yang sama")
	}

	// 5. Validasi: murid yang sama tidak boleh punya jadwal bentrok di hari yang sama
	var count int64
	err = config.DB.Model(&config.Jadwal{}).
		Where("id_murid = ? AND hari_bimbingan = ?", req.IDMurid, strings.ToLower(req.HariBimbingan)).
		Where("? < waktu_selesai AND ? > waktu_mulai", waktuMulaiParsed, waktuSelesaiParsed).
		Count(&count).Error
	if err != nil {
		return errors.New("Gagal memeriksa jadwal murid")
	}
	if count > 0 {
		return errors.New("Jadwal murid bentrok di hari dan jam yang sama")
	}

	// 6. Siapkan Data untuk Model
	jadwalData := config.Jadwal{
		IDPembimbing:  req.IDPembimbing,
		IDMurid:       req.IDMurid,
		HariBimbingan: strings.ToLower(req.HariBimbingan),
		WaktuMulai:    waktuMulaiParsed,
		WaktuSelesai:  waktuSelesaiParsed,
		Ruangan:       roomName,
	}

	return CreateJadwalModel(jadwalData)
}

func UpdateJadwalService(id uint, req JadwalUpdateRequest) error {
	// 1. Validasi Field Wajib (Logic Bisnis)
	if req.IDPembimbing == 0 || req.IDMurid == 0 || req.HariBimbingan == "" || req.WaktuMulai == "" || req.WaktuSelesai == "" {
		return errors.New("Field wajib tidak boleh kosong")
	}

	// 2. Cek data jadwal exist
	var jadwal config.Jadwal
	if err := config.DB.First(&jadwal, id).Error; err != nil {
		return errors.New("Data jadwal tidak ditemukan")
	}

	// 3. Parsing Waktu
	layoutTime := "15:04:05"
	loc, locErr := time.LoadLocation("Asia/Jakarta")
	if locErr != nil {
		loc = time.FixedZone("WIB", 7*60*60)
	}

	waktuMulaiParsed, errM := time.ParseInLocation(layoutTime, req.WaktuMulai, loc)
	waktuSelesaiParsed, errS := time.ParseInLocation(layoutTime, req.WaktuSelesai, loc)

	if errM != nil || errS != nil {
		return errors.New("Format waktu tidak valid. Gunakan HH:mm:ss")
	}

	// Normalisasi tanggal untuk perbandingan jam
	wM := time.Date(2000, 1, 1, waktuMulaiParsed.Hour(), waktuMulaiParsed.Minute(), waktuMulaiParsed.Second(), 0, loc)
	wS := time.Date(2000, 1, 1, waktuSelesaiParsed.Hour(), waktuSelesaiParsed.Minute(), waktuSelesaiParsed.Second(), 0, loc)

	if wS.Before(wM) || wS.Equal(wM) {
		return errors.New("Waktu selesai harus setelah waktu mulai")
	}

	roomName := strings.TrimSpace(req.Ruangan)
	if roomName == "" {
		return errors.New("Ruangan wajib diisi")
	}

	// 4.5 Validasi: pembimbing yang sama tidak boleh punya jadwal bentrok di hari yang sama
	var countPembimbing int64
	errP := config.DB.Model(&config.Jadwal{}).
		Where("id_pembimbing = ? AND hari_bimbingan = ? AND id_jadwal <> ?", req.IDPembimbing, strings.ToLower(req.HariBimbingan), id).
		Where("? < waktu_selesai AND ? > waktu_mulai", wM, wS).
		Count(&countPembimbing).Error
	if errP != nil {
		return errors.New("Gagal memeriksa jadwal pembimbing")
	}
	if countPembimbing > 0 {
		return errors.New("Jadwal pembimbing bentrok di hari dan jam yang sama")
	}

	// 4.6 Validasi: ruangan yang sama tidak boleh dipakai pembimbing lain di jam overlap pada hari yang sama
	var countRuangan int64
	errR := config.DB.Model(&config.Jadwal{}).
		Where("LOWER(TRIM(ruangan)) = LOWER(TRIM(?)) AND hari_bimbingan = ? AND id_jadwal <> ? AND id_pembimbing <> ?", roomName, strings.ToLower(req.HariBimbingan), id, req.IDPembimbing).
		Where("? < waktu_selesai AND ? > waktu_mulai", wM, wS).
		Count(&countRuangan).Error
	if errR != nil {
		return errors.New("Gagal memeriksa bentrok ruangan")
	}
	if countRuangan > 0 {
		return errors.New("Ruangan sudah dipakai pembimbing lain pada hari dan jam yang sama")
	}

	// 5. Validasi: murid yang sama tidak boleh punya jadwal bentrok di hari yang sama
	var count int64
	err := config.DB.Model(&config.Jadwal{}).
		Where("id_murid = ? AND hari_bimbingan = ? AND id_jadwal <> ?", req.IDMurid, strings.ToLower(req.HariBimbingan), id).
		Where("? < waktu_selesai AND ? > waktu_mulai", wM, wS).
		Count(&count).Error
	if err != nil {
		return errors.New("Gagal memeriksa jadwal murid")
	}
	if count > 0 {
		return errors.New("Jadwal murid bentrok di hari dan jam yang sama")
	}

	// 6. Update Struct Object
	jadwal.IDPembimbing = req.IDPembimbing
	jadwal.IDMurid = req.IDMurid
	jadwal.HariBimbingan = strings.ToLower(req.HariBimbingan)
	jadwal.WaktuMulai = wM
	jadwal.WaktuSelesai = wS
	jadwal.Ruangan = roomName

	// 7. Kirim ke Model
	return UpdateJadwalModel(jadwal)
}

func DeleteJadwalService(id_jadwal uint) error {
	// 1. Cek apakah jadwal ada di database
	var jadwal config.Jadwal
	if err := config.DB.First(&jadwal, id_jadwal).Error; err != nil {
		return errors.New("Jadwal tidak ditemukan untuk dihapus")
	}

	// 2. Panggil Model untuk eksekusi delete
	return DeleteJadwalModel(id_jadwal)
}
