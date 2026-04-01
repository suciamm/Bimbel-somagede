package pembimbing

import (
	"data/config"
	"errors"
)

func GetAllPembimbingModel() ([]PembimbingResponse, error) {
	var result []PembimbingResponse

	// Menggunakan Select untuk mengambil field tertentu saja
	err := config.DB.
		Table("users").
		Select("id_user, kode_pembimbing, username, nama_lengkap, no_telp, status").
		Where("role = ? AND status = ?", "pembimbing", true).
		Scan(&result).Error

	return result, err
}

func DeletePembimbingModel(idUser uint) error {
    // Melakukan Hard Delete (menghapus permanen)
    result := config.DB.Unscoped().Delete(&config.User{}, idUser)
    
    if result.Error != nil {
        return errors.New("Gagal menghapus data pembimbing")
    }

    // Jika tidak ada baris yang terhapus
    if result.RowsAffected == 0 {
        return errors.New("Pembimbing tidak ditemukan untuk dihapus")
    }

    return nil
}


func GetPengajuanPembimbingModel() ([]PengajuanPembimbingResponse, error) {

	var pembimbingList []config.User // sesuaikan dengan struct model User kamu
	var response []PengajuanPembimbingResponse

	// Ambil user role pembimbing dan status = false
	result := config.DB.
		Where("role = ? AND status = ?", "pembimbing", false).
		Find(&pembimbingList)

	if result.Error != nil {
		return nil, result.Error
	}

	// Mapping ke response + konversi status
	for _, pb := range pembimbingList {
		status := "Tidak Aktif"
		if pb.Status {
			status = "Aktif"
		}

		response = append(response, PengajuanPembimbingResponse{
			IdUser:      pb.IDUser,
			Username:    pb.Username,
			NamaLengkap: pb.NamaLengkap,
			NoTelpon:    pb.NoTelp,
			Alamat:      pb.Alamat,
			Status:      status,
		})
	}

	return response, nil
}

func GetPengajuanOrtuModel() ([]PengajuanPembimbingResponse, error) {
	var orangtuaList []config.User
	var response []PengajuanPembimbingResponse

	// Ambil user role pembimbing dan status = false
	result := config.DB.
		Where("role = ? AND status = ?", "orangtua", false).
		Find(&orangtuaList)

	if result.Error != nil {
		return nil, result.Error
	}

	// Mapping ke response + konversi status
	for _, pb := range orangtuaList {
		status := "Tidak Aktif"
		if pb.Status {
			status = "Aktif"
		}

		response = append(response, PengajuanPembimbingResponse{
			IdUser:      pb.IDUser,
			Username:    pb.Username,
			NamaLengkap: pb.NamaLengkap,
			NoTelpon:    pb.NoTelp,
			Alamat:      pb.Alamat,
			Status:      status,
		})
	}

	return response, nil

}

func UbahStatusPembimbingModel(user config.User, action string) error {
	if action == "setujui" {
		// Save mengupdate semua field termasuk status dan kode_pembimbing yang baru
		if err := config.DB.Save(&user).Error; err != nil {
			return errors.New("Gagal memperbarui status pembimbing")
		}
	} else if action == "tolak" {
		// Menghapus data user
		if err := config.DB.Delete(&user).Error; err != nil {
			return errors.New("Gagal menghapus data pembimbing")
		}
	}
	return nil
}

func UbahStatusOrangtuaModel(user config.User, action string) error {
	if action == "setujui" {
		// Save mengupdate semua field termasuk status dan kode_pembimbing yang baru
		if err := config.DB.Save(&user).Error; err != nil {
			return errors.New("Gagal memperbarui status orangtua")
		}
	} else if action == "tolak" {
		// Menghapus data user
		if err := config.DB.Delete(&user).Error; err != nil {
			return errors.New("Gagal menghapus data orangtua")
		}
	}
	return nil
}

func UpdatePembimbingModel(id_user uint, updateData map[string]interface{}) error {
	if err := config.DB.Model(&config.User{}).Where("id_user = ?", id_user).Updates(updateData).Error; err != nil {
		return errors.New("Gagal memperbarui data pembimbing")
	}
	return nil
}
