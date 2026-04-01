package users

import (
	"data/config"
	"database/sql"
	"time"

	"gorm.io/gorm"
)


// ----------------- VERSI TANPA JOIN
// UserResponse dipakai untuk contoh "select sebagian kolom".
type UserResponse struct {
	IDUser      int       `json:"id_user"`
	NamaLengkap string    `json:"nama_lengkap"`
	Email       string    `json:"email"`
	NoHP        string    `json:"no_hp"`
	Alamat      string    `json:"alamat"`
	Role        string    `json:"role"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// 1) Tanpa GORM (manual SQL + scan)
func GetUserDetailManualSQL(id_user int, db *sql.DB) (UserResponse, error) {
	query := `
		SELECT id_user, nama_lengkap, email, no_hp, alamat, role, status
		FROM users
		WHERE id_user = ?
	`

	row := db.QueryRow(query, id_user)
	var u UserResponse
	err := row.Scan(
		&u.IDUser,
		&u.NamaLengkap,
		&u.Email,
		&u.NoHP,
		&u.Alamat,
		&u.Role,
		&u.Status,
	)
	if err != nil {
		return UserResponse{}, err
	}

	return u, nil
}

// 2) Dengan GORM (select kolom tertentu)
func GetUserDetailGormSelected(id_user int, db *gorm.DB) (UserResponse, error) {
	var user UserResponse
	err := db.Table("users").
		Select("id_user, nama_lengkap, email, no_hp, alamat, role, status, created_at").
		Where("id_user = ?", id_user).
		First(&user).Error

	return user, err
}

// 3) Dengan GORM (ambil full data model langsung panggil nama tabelnya)
func GetUserDetailGormFull(id_user int, db *gorm.DB) (config.User, error) {
	var user config.User
	err := db.Where("id_user = ?", id_user).First(&user).Error
	return user, err
}

// ------------------ VERSI JOIN
type MuridListByOrtuResponse struct {
	IDMurid        int       `json:"id_murid"`
	NamaMurid      string    `json:"nama_murid"`
	NamaOrtu       string    `json:"nama_ortu"`
	TglMasuk       time.Time `json:"tgl_masuk"`
	HariBimbingan  string    `json:"hari_bimbingan"`
	WaktuMulai     string    `json:"waktu_mulai"`
	WaktuSelesai   string    `json:"waktu_selesai"`
	Ruangan        string    `json:"ruangan"`
	NamaPembimbing string    `json:"nama_pembimbing"`
	StatusMurid    string    `json:"status_murid"`
}

// 1) Join tanpa GORM
func GetDataBimbinganByOrtuManualSQL(id_user int, db *sql.DB) ([]MuridListByOrtuResponse, error) {
	query := `
		SELECT
			murid.id_murid,
			murid.nama_murid,
			users.nama_lengkap AS nama_ortu,
			murid.tgl_masuk,
			jadwal.hari_bimbingan,
			jadwal.waktu_mulai,
			jadwal.waktu_selesai,
			jadwal.ruangan,
			pembimbing.nama_lengkap AS nama_pembimbing,
			murid.status_murid
		FROM murid
		JOIN users ON users.id_user = murid.id_user
		JOIN jadwal ON jadwal.id_murid = murid.id_murid
		LEFT JOIN users AS pembimbing ON pembimbing.id_user = jadwal.id_pembimbing
		WHERE murid.id_user = ? AND users.role = ? AND murid.status_murid = ?
		ORDER BY murid.nama_murid ASC, jadwal.hari_bimbingan ASC, jadwal.waktu_mulai ASC`

	rows, err := db.Query(query, id_user, "orangtua", "aktif")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []MuridListByOrtuResponse
	for rows.Next() {
		var r MuridListByOrtuResponse
		err := rows.Scan(
			&r.IDMurid,
			&r.NamaMurid,
			&r.NamaOrtu,
			&r.TglMasuk,
			&r.HariBimbingan,
			&r.WaktuMulai,
			&r.WaktuSelesai,
			&r.Ruangan,
			&r.NamaPembimbing,
			&r.StatusMurid,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// 2) Join dengan GORM
func GetDataBimbinganByOrtuGorm(id_user int, db *gorm.DB) ([]MuridListByOrtuResponse, error) {
	var result []MuridListByOrtuResponse

	err := db.
		Table("murid").
		Select(`
			murid.id_murid,
			murid.nama_murid,
			users.nama_lengkap AS nama_ortu,
			murid.tgl_masuk,
			jadwal.hari_bimbingan,
			jadwal.waktu_mulai,
			jadwal.waktu_selesai,
			jadwal.ruangan,
			pembimbing.nama_lengkap AS nama_pembimbing,
			murid.status_murid
		`).
		Joins("JOIN users ON users.id_user = murid.id_user").
		Joins("JOIN jadwal ON jadwal.id_murid = murid.id_murid").
		Joins("LEFT JOIN users AS pembimbing ON pembimbing.id_user = jadwal.id_pembimbing").
		Where("murid.id_user = ? AND users.role = ? AND murid.status_murid = ?", id_user, "orangtua", "aktif").
		Order("murid.nama_murid ASC, jadwal.hari_bimbingan ASC, jadwal.waktu_mulai ASC").
		Scan(&result).Error

	return result, err
}
