package pembimbing

import (
	"data/config"
	"errors"
	"fmt"

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

func UbahStatusPembimbingService(req RequestUbahStatus) error {
	// 1. Cari pembimbing berdasarkan id user
	var pembimbing config.User
	if err := config.DB.First(&pembimbing, req.IdUser).Error; err != nil {
		return errors.New("Pembimbing tidak ditemukan")
	}

	switch req.Status {
	case "setujui":
		// 2. Logika Generate Kode (Hanya jika role pembimbing)
		if pembimbing.Role == "pembimbing" {
			var lastPembimbing config.User
			config.DB.
				Where("kode_pembimbing IS NOT NULL AND kode_pembimbing != ''").
				Order("kode_pembimbing DESC").
				First(&lastPembimbing)

			var kode string
			if lastPembimbing.KodePembimbing == "" {
				kode = "PBM001"
			} else {
				var num int
				fmt.Sscanf(lastPembimbing.KodePembimbing, "PBM%03d", &num)
				num++
				kode = fmt.Sprintf("PBM%03d", num)
			}
			pembimbing.KodePembimbing = kode
		}

		pembimbing.Status = true
		// Panggil model untuk Update/Save
		return UbahStatusPembimbingModel(pembimbing, "setujui")

	case "tolak":
		// Panggil model untuk Delete
		return UbahStatusPembimbingModel(pembimbing, "tolak")

	default:
		return errors.New("Status tidak valid. Gunakan 'setujui' atau 'tolak'.")
	}
}




func UbahStatusOrangtuaService(req RequestUbahStatus) error {
	// 1. Cari pembimbing berdasarkan id user
	var orangtua config.User
	if err := config.DB.First(&orangtua, req.IdUser).Error; err != nil {
		return errors.New("Orangtua tidak ditemukan")
	}

	switch req.Status {
	case "setujui":

		orangtua.Status = true
		// Panggil model untuk Update/Save
		return UbahStatusOrangtuaModel(orangtua, "setujui")

	case "tolak":
		// Panggil model untuk Delete
		return UbahStatusOrangtuaModel(orangtua, "tolak")

	default:
		return errors.New("Status tidak valid. Gunakan 'setujui' atau 'tolak'.")
	}
}


func UpdatePembimbingService(id_user uint, req UpdatePembimbingRequest) error {
    // 1. Cari Pembimbing yang akan diupdate
    var existingUser config.User
    if err := config.DB.First(&existingUser, id_user).Error; err != nil {
        return errors.New("Data pembimbing tidak ditemukan untuk diupdate")
    }

    // 2. Siapkan data update (menggunakan Map agar field nil tidak terupdate)
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

    // 3. Panggil Model
    return UpdatePembimbingModel(id_user, updateData)
}


func DeletePembimbingService(idUser uint) error {
    // 1. Cari user untuk memastikan datanya ada dan cek role-nya
    var userToDelete config.User
    if err := config.DB.First(&userToDelete, idUser).Error; err != nil {
        return errors.New("User tidak ditemukan.")
    }

    // 2. Cek role (Logika bisnis: dilarang hapus admin di sini)
    if userToDelete.Role == "admin" {
        return errors.New("Akses ditolak: Tidak dapat menghapus user Admin melalui endpoint ini.")
    }

    // 3. Panggil Model untuk eksekusi hapus
    return DeletePembimbingModel(idUser)
}