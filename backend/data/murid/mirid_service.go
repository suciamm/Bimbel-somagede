package murid

import (
	"data/config"
	"errors"
	"fmt"
	"log"
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

func CreateMuridService(req CreateMuridRequest) error {
	// 1. Parse tanggal
	tglLahir, err := time.Parse("2006-01-02", req.TglLahir)
	if err != nil {
		return errors.New("format tgl_lahir harus YYYY-MM-DD")
	}

	tglMasuk, err := time.Parse("2006-01-02", req.TglMasuk)
	if err != nil {
		return errors.New("format tgl_masuk harus YYYY-MM-DD")
	}

	// 2. Validasi user (orang tua)
	var user config.User
	if err := config.DB.First(&user, req.IDUser).Error; err != nil {
		return errors.New("data orang tua tidak ditemukan")
	}

	// --- TAMBAHAN VALIDASI DUPLIKAT ---
	// Cek apakah orang tua ini sudah mendaftarkan anak dengan nama yang sama
	var count int64
	config.DB.Model(&config.Murid{}).
		Where("id_user = ? AND nama_murid = ?", req.IDUser, req.NamaMurid).
		Count(&count)

	if count > 0 {
		return errors.New("murid dengan nama tersebut sudah terdaftar untuk orang tua ini")
	}
	// ----------------------------------

	// 3. Generate kode murid (Logika bisnis)
	var lastMurid config.Murid
	kode := "MRD001"
	if err := config.DB.Order("kode_murid DESC").First(&lastMurid).Error; err == nil {
		var num int
		fmt.Sscanf(lastMurid.KodeMurid, "MRD%03d", &num)
		kode = fmt.Sprintf("MRD%03d", num+1)
	}

	// 4. Masukkan ke Model
	muridData := config.Murid{
		KodeMurid:   kode,
		NamaMurid:   req.NamaMurid,
		TglLahir:    tglLahir,
		IDUser:      req.IDUser,
		Alamat:      req.Alamat,
		TglMasuk:    tglMasuk,
		StatusMurid: "aktif",
	}

	return CreateMuridModel(muridData)
}

func UpdateMuridService(id int, req UpdateMuridRequest) error {
	// 1. Cek keberadaan murid
	var existingMurid config.Murid
	if err := config.DB.First(&existingMurid, id).Error; err != nil {
		return errors.New("Data murid tidak ditemukan")
	}

	// 2. Siapkan map untuk update (Partial Update)
	updateData := make(map[string]interface{})

	if req.NamaMurid != nil {
		updateData["nama_murid"] = *req.NamaMurid
	}
	if req.IDUser != nil {
		updateData["id_user"] = *req.IDUser
	}
	if req.Alamat != nil {
		updateData["alamat"] = *req.Alamat
	}
	if req.StatusMurid != nil {
		updateData["status_murid"] = *req.StatusMurid
	}

	// Helper internal untuk parse tanggal
	parseDate := func(dateStr *string, fieldName string) error {
		if dateStr != nil {
			if *dateStr == "" {
				updateData[fieldName] = nil
			} else {
				parsed, err := time.Parse("2006-01-02", *dateStr)
				if err != nil {
					return errors.New("Format " + fieldName + " salah")
				}
				updateData[fieldName] = parsed
			}
		}
		return nil
	}

	if err := parseDate(req.TglLahir, "tgl_lahir"); err != nil {
		return err
	}
	if err := parseDate(req.TglMasuk, "tgl_masuk"); err != nil {
		return err
	}
	if err := parseDate(req.TglKeluar, "tgl_keluar"); err != nil {
		return err
	}

	// 3. Kirim ke Model
	return UpdateMuridModel(id, updateData)
}

func UpdateOrtuService(id int, req UpdateOrtuRequest) error {
	// 1. Cek keberadaan User (Ortu)
	var existingUser config.User
	if err := config.DB.First(&existingUser, id).Error; err != nil {
		return errors.New("Data user tidak ditemukan")
	}

	// 2. Siapkan map untuk partial update
	updateData := make(map[string]interface{})

	if req.NamaLengkap != nil {
		updateData["nama_lengkap"] = *req.NamaLengkap
	}
	if req.NoTelp != nil {
		updateData["no_telp"] = *req.NoTelp
	}
	if req.Alamat != nil {
		updateData["alamat"] = *req.Alamat
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	// 3. Kirim ke Model
	return UpdateOrtuModel(id, updateData)
}

func DeleteMuridService(id uint) error {
	// 1. Pastikan muridnya ada
	var murid config.Murid
	if err := config.DB.First(&murid, id).Error; err != nil {
		return errors.New("data murid tidak ditemukan")
	}

	// 2. Cek apakah masih ada langganan aktif
	count, err := CheckIfMuridHasLanggananModel(config.DB, id)
	if err != nil {
		// ⚠️ JANGAN return error ke user
		// cukup log, karena ini error sistem
		log.Println("ERROR cek langganan murid:", err)
	}

	if count > 0 {
		return errors.New("murid ini masih memiliki bimbingan aktif dan tidak boleh dihapus")
	}

	// 3. Hapus murid
	return DeleteMuridModel(config.DB, id)
}

func DeleteOrtuService(id uint) error {
	// 1. Pastikan data orangtua ada
	var ortu config.User
	if err := config.DB.First(&ortu, id).Error; err != nil {
		return errors.New("data orangtua tidak ditemukan")
	}

	// 2. Cek apakah masih punya murid
	count, err := CheckIfOruHasMuridModel(config.DB, id)
	if err != nil {
		log.Println("ERROR cek murid ortu:", err)
		return errors.New("gagal mengecek relasi murid orangtua")
	}

	if count > 0 {
		return errors.New("orangtua ini masih memiliki murid dan tidak boleh dihapus")
	}

	// 3. Hapus orangtua
	return DeleteOrtuModel(config.DB, id)
}
