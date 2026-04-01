package evaluasi

import (
	"data/config"
	"errors"
	"strings"
	"time"
)

var allowedNilai = map[string]bool{
	"A": true,
	"B": true,
	"C": true,
	"D": true,
	"E": true,
}

func CreateEvaluasiService(req CreateEvaluasiRequest) error {
	if req.EvaluasiKe < 1 || req.EvaluasiKe > 3 {
		return errors.New("evaluasi_ke hanya boleh 1, 2, atau 3")
	}

	nilai := strings.ToUpper(strings.TrimSpace(req.Nilai))
	if !allowedNilai[nilai] {
		return errors.New("nilai evaluasi tidak valid, gunakan A/B/C/D/E")
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60)
	}

	tanggalEvaluasi, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(req.TanggalEvaluasi), loc)
	if err != nil {
		return errors.New("format tanggal_evaluasi harus YYYY-MM-DD")
	}

	if tanggalEvaluasi.Weekday() != config.AllowedEvaluasiWeekday {
		return errors.New("evaluasi hanya boleh diinput pada hari " + config.AllowedEvaluasiWeekdayLabel)
	}

	var murid config.Murid
	if err := config.DB.First(&murid, req.IDMurid).Error; err != nil {
		return errors.New("murid tidak ditemukan")
	}

	var pembimbing config.User
	if err := config.DB.Where("id_user = ? AND role = ?", req.IDPembimbing, "pembimbing").First(&pembimbing).Error; err != nil {
		return errors.New("pembimbing tidak ditemukan")
	}

	var relasiCount int64
	if err := config.DB.Model(&config.Jadwal{}).Where("id_murid = ? AND id_pembimbing = ?", req.IDMurid, req.IDPembimbing).Count(&relasiCount).Error; err != nil {
		return errors.New("gagal memeriksa relasi pembimbing dan murid")
	}
	if relasiCount == 0 {
		return errors.New("murid ini bukan murid bimbingan pembimbing terkait")
	}

	var evaluasiSamaTahap config.EvaluasiMurid
	err = config.DB.Where("id_murid = ? AND evaluasi_ke = ?", req.IDMurid, req.EvaluasiKe).First(&evaluasiSamaTahap).Error
	if err == nil {
		return errors.New("evaluasi tahap ini sudah pernah diinput")
	}

	if req.EvaluasiKe > 1 {
		var evaluasiSebelumnya config.EvaluasiMurid
		err := config.DB.Where("id_murid = ? AND evaluasi_ke = ?", req.IDMurid, req.EvaluasiKe-1).First(&evaluasiSebelumnya).Error
		if err != nil {
			return errors.New("evaluasi tahap sebelumnya belum ada")
		}
		if strings.ToUpper(strings.TrimSpace(evaluasiSebelumnya.Nilai)) != "A" {
			return errors.New("tidak bisa lanjut ke tahap berikutnya karena nilai tahap sebelumnya belum A")
		}
	}

	newData := config.EvaluasiMurid{
		IDMurid:           req.IDMurid,
		IDPembimbing:      req.IDPembimbing,
		EvaluasiKe:        req.EvaluasiKe,
		Nilai:             nilai,
		CatatanPembimbing: strings.TrimSpace(req.CatatanPembimbing),
		TanggalEvaluasi:   tanggalEvaluasi,
	}

	if err := CreateEvaluasiModel(newData, config.DB); err != nil {
		return errors.New("gagal menyimpan evaluasi")
	}

	return nil
}

func EnrichProgressInfo(items []EvaluasiPerMuridResponse) []EvaluasiPerMuridResponse {
	for i := range items {
		items[i].TahapBerikutnya = 1
		items[i].BolehInputTahap2 = false
		items[i].BolehInputTahap3 = false
		items[i].StatusProgressText = "Belum ada evaluasi"

		if items[i].Evaluasi1Nilai != nil {
			if strings.ToUpper(strings.TrimSpace(*items[i].Evaluasi1Nilai)) == "A" {
				items[i].TahapBerikutnya = 2
				items[i].BolehInputTahap2 = true
				items[i].StatusProgressText = "Tahap 1 lulus"
			} else {
				items[i].TahapBerikutnya = 1
				items[i].StatusProgressText = "Tahap 1 belum lulus (butuh A)"
			}
		}

		if items[i].Evaluasi2Nilai != nil {
			if strings.ToUpper(strings.TrimSpace(*items[i].Evaluasi2Nilai)) == "A" {
				items[i].TahapBerikutnya = 3
				items[i].BolehInputTahap3 = true
				items[i].StatusProgressText = "Tahap 2 lulus"
			} else {
				items[i].TahapBerikutnya = 2
				items[i].BolehInputTahap3 = false
				items[i].StatusProgressText = "Tahap 2 belum lulus (butuh A)"
			}
		}

		if items[i].Evaluasi3Nilai != nil {
			items[i].TahapBerikutnya = 3
			items[i].BolehInputTahap2 = false
			items[i].BolehInputTahap3 = false
			items[i].StatusProgressText = "Evaluasi tahap 1-3 selesai"
		}
	}

	return items
}
