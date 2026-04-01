package api

import (
	"data/evaluasi"
	jadwalmengajar "data/jadwal_mengajar"
	"data/murid"
	"data/payment"
	"data/pembimbing"
	"data/users"
	"log"

	"github.com/gin-gonic/gin"
)

func ApiAllRoleWithToken(r *gin.Engine) {
	log.Println("public api")
	r.POST("/api/auth/register", users.RegisterUser)
	r.POST("/api/auth/login", users.LoginUser)

	// Menu kelola data Murid -
	r.GET("/api/murid/viewDaftarMuridAktif", murid.GetMuridAktifController)
	r.GET("/api/murid/viewDaftarMuridByOrtu/:id_user", murid.GetMuridByOrtuController)
	r.GET("/api/murid/viewDaftarMuridTidakAktif", murid.GetMuridTidakAktifController)
	r.GET("/api/murid/viewRekapBulanan", murid.GetRekapMuridBulananController)
	r.POST("/api/murid/tambahMurid", murid.CreateMuridController)
	r.PUT("/api/murid/editMurid/:id", murid.UpdateMuridController)

	r.PUT("/api/murid/editdataOrtu/:id", murid.UpdateOrtuController)
	r.DELETE("/api/murid/deleteMurid/:id", murid.DeleteMuridController)
	r.GET("/api/murid/viewMuridByPembimbing/:id_user", murid.GetMuridByPembimbingController)
	r.GET("/api/murid/viewDaftarOrtu", murid.GetDaftarOrtuController)
	r.DELETE("/api/murid/deleteOrtu/:id", murid.DeleteOrtuController)
	r.GET("/api/murid/viewBimbinganByOrtu/:id_user", murid.GetDataBimbinganByOrtuController)

	// Menu evaluasi pembimbing
	r.POST("/api/evaluasi/tambah", evaluasi.CreateEvaluasiController)
	r.GET("/api/evaluasi/viewByPembimbing/:id_user", evaluasi.GetEvaluasiByPembimbingController)
	r.GET("/api/evaluasi/viewByOrtu/:id_user", evaluasi.GetEvaluasiByOrtuController)

	r.GET("/api/pembimbing/viewPembimbing", pembimbing.GetAllPembimbingController)
	r.GET("/api/pembimbing/viewDaftarPengajuanPembimbing", pembimbing.GetPengajuanPembimbingController)
	r.GET("/api/pembimbing/viewDaftarPengajuanOrtu", pembimbing.GetPengajuanOrtugController)
	r.PUT("/api/pembimbing/approvalRegistrasiPembimbing", pembimbing.UbahStatusPembimbingController)
	r.PUT("/api/pembimbing/approvalRegistrasiOrangtua", pembimbing.UbahStatusOrangtuaController)
	r.PUT("/api/pembimbing/editDataPembimbing/:id", pembimbing.UpdatePembimbingController)
	r.DELETE("/api/pembimbing/deletePembimbing/:id", pembimbing.DeletePembimbingController)

	r.GET("/api/jadwal/viewJadwal", jadwalmengajar.GetAllJadwalController)
	r.GET("/api/jadwal/viewJadwalByPembimbing/:id_user", jadwalmengajar.GetJadwalByPembimbingController)
	r.POST("/api/jadwal/tambahJadwal", jadwalmengajar.CreateJadwalController)
	r.PUT("/api/jadwal/editJadwal/:id", jadwalmengajar.UpdateJadwalController)
	r.DELETE("/api/jadwal/deleteJadwal/:id_jadwal", jadwalmengajar.DeleteJadwalController)

	r.POST("/api/payment/tambahPaketBimbingan", payment.TambahPaketBimbinganController)
	r.GET("/api/payment/viewPaketBimbingan", payment.GetPaketBimbinganController)
	r.GET("/api/payment/viewPaketBimbinganAktif", payment.GetPaketBimbinganAktifController)
	r.PUT("/api/payment/editPaketBimbingan/:id_paket", payment.EditPaketBimbinganController)
	r.DELETE("/api/payment/deletePaketBimbingan/:id_paket", payment.DeletePaketBimbinganController)
	r.POST("/api/payment/tambahTransaksiPembayaran", payment.TambahTransaksiController)

	r.GET("/api/payment/viewTransaksiPembayaran", payment.GetTransaksiPembayaranController)
	r.DELETE("/api/payment/hapusTransaksiPembayaran/:id_transaksi", payment.HapusTransaksiPembayaranController)
	r.GET("/api/payment/viewTransaksiByOrtu/:id_user", payment.GetTransaksiByMuridController)

	// Midtrans endpoints
	r.POST("/api/payment/createMidtransTransaction", payment.CreateMidtransTransactionController)
	r.POST("/api/payment/perpanjangPaketMidtrans/:id_paket", payment.PerpanjangPaketMidtransController)
	r.POST("/api/payment/konfirmasiPerpanjangMidtrans/:id_paket", payment.KonfirmasiPerpanjangMidtransController)
	r.POST("/api/midtrans-callback", payment.MidtransCallbackController)

}
