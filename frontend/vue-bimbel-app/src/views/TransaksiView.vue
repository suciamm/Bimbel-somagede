<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import {
  createMidtransTransactionApi,
  deleteTransaksiApi,
  getMuridAktifApi,
  getPaketAktifApi,
  getTransaksiApi,
  tambahTransaksiApi
} from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showCreateModal = ref(false)
const searchTransaksi = ref('')
const currentPageTransaksi = ref(1)
const PAGE_SIZE = 10

const transaksiList = ref([])
const paketAktif = ref([])
const muridAktif = ref([])
const MIDTRANS_CLIENT_KEY = import.meta.env.VITE_MIDTRANS_CLIENT_KEY || 'Mid-client-7r84InjwqcNGIMqF'

const form = reactive({
  id_murid: '',
  id_paket: '',
  tanggal_mulai: '',
  tanggal_selesai: '',
  tanggal_bayar: '',
  jumlah_bayar: 0,
  metode_bayar: 'transfer'
})

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function getStatus(item) {
  return item.status || '-'
}

function formatMuridDisplay(item) {
  const kode = item?.kode_murid || item?.KodeMurid || ''
  const nama = item?.nama_murid || item?.NamaMurid || ''
  if (kode && nama) return `(${kode})${nama}`
  return nama || '-'
}

function getSelectedPaketHarga() {
  const selectedId = Number(form.id_paket)
  if (!selectedId) return 0

  const paket = paketAktif.value.find((item) => Number(item.id_paket) === selectedId)
  return Number(paket?.harga || 0)
}

function getSelectedPaketDurasiHari() {
  const selectedId = Number(form.id_paket)
  if (!selectedId) return 0

  const paket = paketAktif.value.find((item) => Number(item.id_paket) === selectedId)
  const durasiBulan = Number(paket?.durasi_bulan || 0)
  const durasiHari = Number(paket?.durasi_hari || 0)
  return durasiHari + durasiBulan * 30
}

function hitungTotalHariInklusif(tanggalMulai, tanggalSelesai) {
  if (!tanggalMulai || !tanggalSelesai) return 0

  const mulai = new Date(tanggalMulai)
  const selesai = new Date(tanggalSelesai)
  if (Number.isNaN(mulai.getTime()) || Number.isNaN(selesai.getTime())) return 0

  const msPerHari = 24 * 60 * 60 * 1000
  const selisih = Math.floor((selesai - mulai) / msPerHari) + 1
  return Math.max(0, selisih)
}

function calculateJumlahBayar() {
  const hargaPaket = getSelectedPaketHarga()
  const durasiPaketHari = getSelectedPaketDurasiHari()
  if (!hargaPaket || !durasiPaketHari) return 0

  const totalHari = hitungTotalHariInklusif(form.tanggal_mulai, form.tanggal_selesai)
  if (totalHari <= 0) return hargaPaket

  const kelipatanPaket = Math.max(1, Math.ceil(totalHari / durasiPaketHari))
  return kelipatanPaket * hargaPaket
}

function formatCurrency(value) {
  return `Rp ${Number(value || 0).toLocaleString('id-ID')}`
}

function getSelectedMurid() {
  const selectedId = Number(form.id_murid)
  if (!selectedId) return null
  return muridAktif.value.find((item) => Number(item.id_murid) === selectedId) || null
}

function getSelectedPaket() {
  const selectedId = Number(form.id_paket)
  if (!selectedId) return null
  return paketAktif.value.find((item) => Number(item.id_paket) === selectedId) || null
}

function loadMidtransSnapScript() {
  if (window.snap) return Promise.resolve()

  return new Promise((resolve, reject) => {
    const existing = document.getElementById('midtrans-snap-script')
    if (existing) {
      existing.addEventListener('load', () => resolve())
      existing.addEventListener('error', () => reject(new Error('Gagal memuat Midtrans Snap'))) 
      return
    }

    const script = document.createElement('script')
    script.id = 'midtrans-snap-script'
    script.src = 'https://app.sandbox.midtrans.com/snap/snap.js'
    script.setAttribute('data-client-key', MIDTRANS_CLIENT_KEY)
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Gagal memuat Midtrans Snap'))
    document.body.appendChild(script)
  })
}

async function startMidtransPayment() {
  const selectedMurid = getSelectedMurid()
  const selectedPaket = getSelectedPaket()

  if (!selectedMurid || !selectedPaket) {
    throw new Error('Murid dan paket wajib dipilih sebelum pembayaran Midtrans')
  }

  const payload = {
    id_murid: Number(form.id_murid),
    id_paket: Number(form.id_paket),
    jumlah_bayar: Math.round(Number(form.jumlah_bayar || 0)),
    nama_paket: selectedPaket.nama_paket || 'Paket Bimbingan',
    nama_user: selectedMurid.nama_murid || 'Murid'
  }

  let snapToken
  try {
    const response = await createMidtransTransactionApi(payload)
    snapToken = response?.snap_token || response?.data?.snap_token
  } catch (err) {
    const errorMsg = err?.response?.data?.message || err.message || 'Gagal membuat transaksi Midtrans'
    throw new Error(errorMsg)
  }

  if (!snapToken) {
    throw new Error('Token Midtrans tidak ditemukan. Periksa konfigurasi server atau silakan coba lagi.')
  }

  await loadMidtransSnapScript()

  await new Promise((resolve, reject) => {
    window.snap.pay(snapToken, {
      onSuccess: async () => {
        await tambahTransaksiApi({
          id_murid: Number(form.id_murid),
          id_paket: Number(form.id_paket),
          tanggal_mulai: form.tanggal_mulai,
          tanggal_selesai: form.tanggal_selesai,
          tanggal_bayar: form.tanggal_bayar,
          jumlah_bayar: Number(form.jumlah_bayar),
          metode_bayar: 'midtrans'
        })

        successMessage.value = 'Pembayaran Midtrans berhasil diproses.'
        form.id_murid = ''
        form.id_paket = ''
        form.tanggal_mulai = ''
        form.tanggal_selesai = ''
        form.tanggal_bayar = ''
        form.jumlah_bayar = 0
        form.metode_bayar = 'transfer'
        showCreateModal.value = false
        await loadData()
        resolve()
      },
      onPending: async () => {
        successMessage.value = 'Pembayaran Midtrans pending. Data transaksi belum disimpan.'
        resolve()
      },
      onError: (result) => {
        reject(new Error(result?.status_message || 'Pembayaran Midtrans gagal'))
      },
      onClose: () => {
        // User tutup popup tanpa selesaikan - tidak perlu error, resolve saja
        resolve()
      }
    })
  })
}

function formatTanggal(value) {
  const raw = String(value || '').trim()
  if (!raw) return '-'

  const datePart = raw.split('T')[0]
  if (/^\d{4}-\d{2}-\d{2}$/.test(datePart)) {
    const [year, month, day] = datePart.split('-')
    return `${day}-${month}-${year}`
  }

  return raw
}

function formatPeriode(item) {
  return `${formatTanggal(item?.tgl_mulai)} - ${formatTanggal(item?.tgl_selesai)}`
}

const filteredTransaksiList = computed(() => {
  const keyword = normalizeText(searchTransaksi.value).trim()
  if (!keyword) return transaksiList.value

  return transaksiList.value.filter((item) => {
    const joined = [
      item.kode_murid,
      item.nama_murid,
      item.nama_paket,
      item.tgl_mulai,
      item.tgl_selesai,
      item.jumlah_bayar,
      getStatus(item)
    ]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const totalPagesTransaksi = computed(() => Math.max(1, Math.ceil(filteredTransaksiList.value.length / PAGE_SIZE)))

const paginatedTransaksiList = computed(() => {
  const start = (currentPageTransaksi.value - 1) * PAGE_SIZE
  return filteredTransaksiList.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(index) {
  return (currentPageTransaksi.value - 1) * PAGE_SIZE + index + 1
}

function goPrevPage() {
  if (currentPageTransaksi.value > 1) currentPageTransaksi.value -= 1
}

function goNextPage() {
  if (currentPageTransaksi.value < totalPagesTransaksi.value) currentPageTransaksi.value += 1
}

function downloadTransaksiPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredTransaksiList.value.map((item, index) => [
    index + 1,
    formatMuridDisplay(item),
    item.nama_paket || '-',
    formatPeriode(item),
    formatCurrency(item.jumlah_bayar),
    getStatus(item)
  ])

  doc.setFontSize(12)
  doc.text('Daftar Transaksi Pembayaran', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Murid', 'Paket', 'Periode', 'Jumlah Bayar', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-transaksi.pdf')
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [trxResp, paketResp, muridResp] = await Promise.all([
      getTransaksiApi(),
      getPaketAktifApi(),
      getMuridAktifApi()
    ])

    transaksiList.value = trxResp?.data || trxResp || []
    paketAktif.value = paketResp?.data || paketResp || []
    muridAktif.value = muridResp?.data || muridResp || []
    currentPageTransaksi.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(searchTransaksi, () => {
  currentPageTransaksi.value = 1
})

watch(filteredTransaksiList, () => {
  if (currentPageTransaksi.value > totalPagesTransaksi.value) {
    currentPageTransaksi.value = totalPagesTransaksi.value
  }
})

watch(
  () => [form.id_paket, form.tanggal_mulai, form.tanggal_selesai],
  () => {
    form.jumlah_bayar = calculateJumlahBayar()
  }
)

async function submitTransaksi() {
  errorMessage.value = ''
  successMessage.value = ''

  try {
    if (form.metode_bayar === 'midtrans') {
      await startMidtransPayment()
      return
    }

    await tambahTransaksiApi({
      id_murid: Number(form.id_murid),
      id_paket: Number(form.id_paket),
      tanggal_mulai: form.tanggal_mulai,
      tanggal_selesai: form.tanggal_selesai,
      tanggal_bayar: form.tanggal_bayar,
      jumlah_bayar: Number(form.jumlah_bayar),
      metode_bayar: form.metode_bayar
    })

    successMessage.value = 'Transaksi berhasil ditambahkan.'
    form.id_murid = ''
    form.id_paket = ''
    form.tanggal_mulai = ''
    form.tanggal_selesai = ''
    form.tanggal_bayar = ''
    form.jumlah_bayar = 0
    form.metode_bayar = 'transfer'
    showCreateModal.value = false

    await loadData()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function handleDelete(idTransaksi) {
  if (!window.confirm('Hapus transaksi pembayaran ini?')) return

  errorMessage.value = ''
  successMessage.value = ''

  try {
    await deleteTransaksiApi(idTransaksi)
    successMessage.value = 'Transaksi berhasil dihapus.'
    await loadData()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(loadData)
</script>

<template>
  <section>
    <h1 class="page-title">Transaksi Pembayaran</h1>
    <p class="page-subtitle">Kelola transaksi bimbingan: tambah data pembayaran dan hapus transaksi.</p>

    <BaseModal :show="showCreateModal" title="Tambah Transaksi Pembayaran" @close="showCreateModal = false">
      <form class="form-grid" @submit.prevent="submitTransaksi">
        <div class="field">
          <label>Murid</label>
          <select v-model="form.id_murid" required>
            <option disabled value="">Pilih murid</option>
            <option v-for="murid in muridAktif" :key="murid.id_murid" :value="murid.id_murid">
              {{ formatMuridDisplay(murid) }}
            </option>
          </select>
        </div>

        <div class="field">
          <label>Paket</label>
          <select v-model="form.id_paket" required>
            <option disabled value="">Pilih paket</option>
            <option v-for="paket in paketAktif" :key="paket.id_paket" :value="paket.id_paket">
              {{ paket.nama_paket }}
            </option>
          </select>
        </div>

        <div class="field">
          <label>Tanggal Mulai</label>
          <input v-model="form.tanggal_mulai" type="date" required />
        </div>

        <div class="field">
          <label>Tanggal Selesai</label>
          <input v-model="form.tanggal_selesai" type="date" required />
        </div>

        <div class="field">
          <label>Tanggal Bayar</label>
          <input v-model="form.tanggal_bayar" type="date" required />
        </div>

        <div class="field">
          <label>Jumlah Bayar</label>
          <input :value="formatCurrency(form.jumlah_bayar)" type="text" readonly />
        </div>

        <div class="field">
          <label>Metode Bayar</label>
          <select v-model="form.metode_bayar">
            <option value="transfer">transfer</option>
            <option value="cash">cash</option>
            <option value="midtrans">midtrans</option>
          </select>
        </div>

        <div>
          <button class="btn btn-primary" type="submit">
            {{ form.metode_bayar === 'midtrans' ? 'Lanjut Bayar Midtrans' : 'Simpan Transaksi' }}
          </button>
        </div>
      </form>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>
    </BaseModal>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Transaksi</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchTransaksi" placeholder="Cari data transaksi..." />
        <div class="tools-actions">
          <button class="btn btn-primary" type="button" @click="showCreateModal = true">Tambah Transaksi</button>
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF data transaksi"
            @click="downloadTransaksiPdf"
          >
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Murid</th>
              <th>Paket</th>
              <th>Periode</th>
              <th>Jumlah Bayar</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedTransaksiList" :key="item.id_transaksi || index">
              <td>{{ nomorUrut(index) }}</td>
              <td>{{ formatMuridDisplay(item) }}</td>
              <td>{{ item.nama_paket || '-' }}</td>
              <td>{{ formatPeriode(item) }}</td>
              <td>{{ formatCurrency(item.jumlah_bayar) }}</td>
              <td>
                <span class="pill" :class="item.status === 'lunas' ? 'pill-success' : item.status === 'pending' ? 'pill-warning' : 'pill-danger'">
                  {{ item.status || '-' }}
                </span>
              </td>
              <td>
                <button class="btn btn-danger" @click="handleDelete(item.id_transaksi)">Hapus</button>
              </td>
            </tr>
            <tr v-if="paginatedTransaksiList.length === 0">
              <td colspan="7">Tidak ada data transaksi.</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="pagination-wrap" v-if="!loading">
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Sebelumnya"
          aria-label="Halaman sebelumnya"
          @click="goPrevPage"
          :disabled="currentPageTransaksi === 1"
        >
          <<
        </button>
        <span>
          Halaman {{ currentPageTransaksi }} / {{ totalPagesTransaksi }} ({{ filteredTransaksiList.length }} data)
        </span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="goNextPage"
          :disabled="currentPageTransaksi === totalPagesTransaksi"
        >
          >>
        </button>
      </div>
    </section>
  </section>
</template>

<style scoped>
.block {
  margin-top: 16px;
  padding: 16px;
}

.block-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

h2 {
  margin: 0 0 12px;
  font-size: 1.05rem;
}

.table-tools {
  display: flex;
  gap: 10px;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.table-tools input {
  width: 100%;
  max-width: 420px;
}

.tools-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.pagination-wrap {
  margin-top: 10px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  flex-wrap: wrap;
}

.btn-icon {
  min-width: 40px;
  text-align: center;
  font-weight: 700;
  padding: 4px 8px;
  background: transparent;
  border: none;
  color: inherit;
  box-shadow: none;
}

.btn-icon:hover,
.btn-icon:focus,
.btn-icon:active {
  background: transparent;
  color: inherit;
  box-shadow: none;
}

.btn-icon:not(:disabled):hover {
  cursor: pointer;
}

.btn-pdf {
  background-color: #334155;
  border: 1px solid #334155;
  color: #fff;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
}

.btn-pdf:hover,
.btn-pdf:focus,
.btn-pdf:active {
  background-color: #1f2937;
  border-color: #1f2937;
  color: #fff;
}

.pdf-icon {
  font-size: 14px;
  line-height: 1;
}

@media (max-width: 900px) {
  .table-tools {
    flex-direction: column;
    align-items: stretch;
  }

  .table-tools input {
    max-width: none;
  }

  .tools-actions {
    justify-content: flex-end;
    flex-wrap: wrap;
  }
}
</style>
