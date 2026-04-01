import http from './http'

export async function loginApi(payload) {
  const request = {
    username: String(payload?.username || '').trim(),
    password: String(payload?.password || '')
  }

  const { data } = await http.post('/api/auth/login', request)
  return data
}

export async function registerApi(payload) {
  const allowedRoles = ['orangtua', 'pembimbing']
  const role = String(payload?.role || '').trim().toLowerCase()

  if (!allowedRoles.includes(role)) {
    throw new Error('Role registrasi hanya untuk orangtua atau pembimbing.')
  }

  const request = {
    username: String(payload?.username || '').trim(),
    password: String(payload?.password || ''),
    nama_lengkap: String(payload?.nama_lengkap || '').trim(),
    role,
    no_telp: String(payload?.no_telp || '').trim(),
    alamat: String(payload?.alamat || '').trim()
  }

  const { data } = await http.post('/api/auth/register', request)
  return data
}

export async function getMuridAktifApi() {
  const { data } = await http.get('/api/murid/viewDaftarMuridAktif')
  return data
}

export async function getMuridTidakAktifApi() {
  const { data } = await http.get('/api/murid/viewDaftarMuridTidakAktif')
  return data
}

export async function getMuridByOrtuApi(idUser) {
  const { data } = await http.get(`/api/murid/viewDaftarMuridByOrtu/${idUser}`)
  return data
}

export async function getRekapMuridBulananApi({ tanggalMulai, tanggalSelesai, bulan, tahun } = {}) {
  const params = new URLSearchParams()
  if (tanggalMulai) params.set('tanggal_mulai', tanggalMulai)
  if (tanggalSelesai) params.set('tanggal_selesai', tanggalSelesai)
  if (bulan) params.set('bulan', bulan)
  if (tahun) params.set('tahun', tahun)
  const query = params.toString()

  const { data } = await http.get(`/api/murid/viewRekapBulanan${query ? `?${query}` : ''}`)
  return data
}

export async function tambahMuridApi(payload) {
  const { data } = await http.post('/api/murid/tambahMurid', payload)
  return data
}

export async function updateMuridApi(id, payload) {
  const { data } = await http.put(`/api/murid/editMurid/${id}`, payload)
  return data
}

export async function deleteMuridApi(id) {
  const { data } = await http.delete(`/api/murid/deleteMurid/${id}`)
  return data
}

export async function getDaftarOrtuApi() {
  const { data } = await http.get('/api/murid/viewDaftarOrtu')
  return data
}

export async function getMuridByPembimbingApi(idUser) {
  const { data } = await http.get(`/api/murid/viewMuridByPembimbing/${idUser}`)
  return data
}

export async function createEvaluasiApi(payload) {
  const { data } = await http.post('/api/evaluasi/tambah', payload)
  return data
}

export async function getEvaluasiByPembimbingApi(idUser) {
  const { data } = await http.get(`/api/evaluasi/viewByPembimbing/${idUser}`)
  return data
}

export async function getEvaluasiByOrtuApi(idUser) {
  const { data } = await http.get(`/api/evaluasi/viewByOrtu/${idUser}`)
  return data
}

export async function updateOrtuApi(id, payload) {
  const { data } = await http.put(`/api/murid/editdataOrtu/${id}`, payload)
  return data
}

export async function deleteOrtuApi(id) {
  const { data } = await http.delete(`/api/murid/deleteOrtu/${id}`)
  return data
}

export async function getJadwalApi() {
  const { data } = await http.get('/api/jadwal/viewJadwal')
  return data
}

export async function getJadwalByPembimbingApi(idUser) {
  const { data } = await http.get(`/api/jadwal/viewJadwalByPembimbing/${idUser}`)
  return data
}

export async function getBimbinganByOrtuApi(idUser) {
  const { data } = await http.get(`/api/murid/viewBimbinganByOrtu/${idUser}`)
  return data
}

export async function tambahJadwalApi(payload) {
  const { data } = await http.post('/api/jadwal/tambahJadwal', payload)
  return data
}

export async function editJadwalApi(idJadwal, payload) {
  const { data } = await http.put(`/api/jadwal/editJadwal/${idJadwal}`, payload)
  return data
}

export async function deleteJadwalApi(idJadwal) {
  const { data } = await http.delete(`/api/jadwal/deleteJadwal/${idJadwal}`)
  return data
}

export async function getPembimbingApi() {
  const { data } = await http.get('/api/pembimbing/viewPembimbing')
  return data
}

export async function getPengajuanPembimbingApi() {
  const { data } = await http.get('/api/pembimbing/viewDaftarPengajuanPembimbing')
  return data
}

export async function getPengajuanOrtuApi() {
  const { data } = await http.get('/api/pembimbing/viewDaftarPengajuanOrtu')
  return data
}

export async function approvalPembimbingApi(payload) {
  const { data } = await http.put('/api/pembimbing/approvalRegistrasiPembimbing', payload)
  return data
}

export async function approvalOrtuApi(payload) {
  const { data } = await http.put('/api/pembimbing/approvalRegistrasiOrangtua', payload)
  return data
}

export async function updatePembimbingApi(id, payload) {
  const { data } = await http.put(`/api/pembimbing/editDataPembimbing/${id}`, payload)
  return data
}

export async function deletePembimbingApi(id) {
  const { data } = await http.delete(`/api/pembimbing/deletePembimbing/${id}`)
  return data
}

export async function getPaketApi() {
  const { data } = await http.get('/api/payment/viewPaketBimbingan')
  return data
}

export async function getPaketAktifApi() {
  const { data } = await http.get('/api/payment/viewPaketBimbinganAktif')
  return data
}

export async function tambahPaketApi(payload) {
  const { data } = await http.post('/api/payment/tambahPaketBimbingan', payload)
  return data
}

export async function editPaketApi(idPaket, payload) {
  const { data } = await http.put(`/api/payment/editPaketBimbingan/${idPaket}`, payload)
  return data
}

export async function deletePaketApi(idPaket) {
  const { data } = await http.delete(`/api/payment/deletePaketBimbingan/${idPaket}`)
  return data
}

export async function getTransaksiApi() {
  const { data } = await http.get('/api/payment/viewTransaksiPembayaran')
  return data
}

export async function tambahTransaksiApi(payload) {
  const { data } = await http.post('/api/payment/tambahTransaksiPembayaran', payload)
  return data
}

export async function createMidtransTransactionApi(payload) {
  const { data } = await http.post('/api/payment/createMidtransTransaction', payload)
  return data
}

export async function getTransaksiByOrtuApi(idUser) {
  const { data } = await http.get(`/api/payment/viewTransaksiByOrtu/${idUser}`)
  return data
}

export async function perpanjangPaketMidtransApi(idPaket, payload) {
  const { data } = await http.post(`/api/payment/perpanjangPaketMidtrans/${idPaket}`, payload)
  return data
}

export async function konfirmasiPerpanjangMidtransApi(idPaket, payload) {
  const { data } = await http.post(`/api/payment/konfirmasiPerpanjangMidtrans/${idPaket}`, payload)
  return data
}

export async function deleteTransaksiApi(id) {
  const { data } = await http.delete(`/api/payment/hapusTransaksiPembayaran/${id}`)
  return data
}

export async function getAbsensiRekapApi({ bulan, tahun } = {}) {
  const params = new URLSearchParams()
  if (bulan) params.set('bulan', bulan)
  if (tahun) params.set('tahun', tahun)
  const query = params.toString()

  const { data } = await http.get(`/api/v1/absensi${query ? `?${query}` : ''}`)
  return data
}

export async function createAbsensiApi(idJadwal, payload) {
  const { data } = await http.post(`/api/v1/createAbsensi/${idJadwal}`, payload)
  return data
}
