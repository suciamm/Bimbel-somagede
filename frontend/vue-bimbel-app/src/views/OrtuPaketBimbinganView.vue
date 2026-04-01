<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import BaseModal from '../components/BaseModal.vue'
import { useAuthStore } from '../stores/auth'
import {
  createMidtransTransactionApi,
  getBimbinganByOrtuApi,
  getPaketAktifApi,
  getTransaksiByOrtuApi,
  tambahTransaksiApi
} from '../services/api'

const authStore = useAuthStore()
const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const paketList = ref([])
const muridOptions = ref([])
const transaksiByOrtu = ref([])
const searchPaket = ref('')
const showModal = ref(false)
const selectedPaket = ref(null)
const MIDTRANS_CLIENT_KEY = import.meta.env.VITE_MIDTRANS_CLIENT_KEY || 'Mid-client-7r84InjwqcNGIMqF'
const isTanggalSelesaiManual = ref(false)

const form = reactive({
  id_murid: '',
  tanggal_mulai: '',
  tanggal_selesai: '',
  jumlah_bayar: 0
})

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

const filteredPaket = computed(() => {
  const keyword = normalizeText(searchPaket.value).trim()
  if (!keyword) return paketList.value

  return paketList.value.filter((item) => {
    const text = [item?.nama_paket, item?.deskripsi, item?.status].map(normalizeText).join(' ')
    return text.includes(keyword)
  })
})

function formatHarga(value) {
  const number = Number(value || 0)
  return `Rp ${number.toLocaleString('id-ID')}`
}

function toDateInput(value) {
  return String(value || '').match(/\d{4}-\d{2}-\d{2}/)?.[0] || ''
}

function hitungTanggalSelesaiOtomatis(tanggalMulai, paket) {
  if (!tanggalMulai || !paket) return ''

  const totalDurasiHari = getDurasiPaketHari(paket)
  if (totalDurasiHari <= 0) return ''

  const mulai = new Date(`${tanggalMulai}T00:00:00`)
  if (Number.isNaN(mulai.getTime())) return ''

  // Periode inklusif: jika durasi 1 hari, tanggal selesai sama dengan tanggal mulai.
  mulai.setDate(mulai.getDate() + totalDurasiHari - 1)
  return toDateInput(mulai.toISOString())
}

function toDateOnly(value) {
  return String(value || '').match(/\d{4}-\d{2}-\d{2}/)?.[0] || ''
}

function formatDateDisplay(value) {
  const iso = toDateOnly(value)
  if (!iso) return '-'
  const [year, month, day] = iso.split('-')
  return `${day}-${month}-${year}`
}

function isDateOnOrAfter(left, right) {
  if (!left || !right) return false
  return left >= right
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

function getDurasiPaketHari(paket) {
  const durasiBulan = Number(paket?.durasi_bulan || 0)
  const durasiHari = Number(paket?.durasi_hari || 0)
  return durasiHari + durasiBulan * 30
}

function calculateJumlahBayar() {
  const paket = selectedPaket.value
  const hargaPaket = Number(paket?.harga || 0)
  const durasiPaketHari = getDurasiPaketHari(paket)
  if (!hargaPaket || !durasiPaketHari) return 0

  const totalHari = hitungTotalHariInklusif(form.tanggal_mulai, form.tanggal_selesai)
  if (totalHari <= 0) return hargaPaket

  const kelipatanPaket = Math.max(1, Math.ceil(totalHari / durasiPaketHari))
  return kelipatanPaket * hargaPaket
}

function getSelectedMurid() {
  const selectedId = Number(form.id_murid)
  if (!selectedId) return null
  return muridOptions.value.find((item) => Number(item.id_murid) === selectedId) || null
}

function hasActivePackageForMurid(idMurid) {
  const selectedId = Number(idMurid)
  if (!selectedId) return false

  const today = toDateInput(new Date().toISOString())
  return transaksiByOrtu.value.some((item) => {
    if (Number(item?.id_murid) !== selectedId) return false

    const tglSelesai = toDateOnly(item?.tgl_selesai)
    const statusPembayaran = normalizeText(item?.status).trim()
    const belumGagal = statusPembayaran !== 'gagal' && statusPembayaran !== 'failed' && statusPembayaran !== 'expire'

    return belumGagal && isDateOnOrAfter(tglSelesai, today)
  })
}

const activePackageWarning = computed(() => {
  if (!form.id_murid) return ''
  if (!hasActivePackageForMurid(form.id_murid)) return ''
  return 'Murid yang dipilih masih memiliki paket bimbingan aktif. Tunggu paket saat ini selesai sebelum menambah paket baru.'
})

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
  const selected = selectedPaket.value

  if (!selectedMurid || !selected) {
    throw new Error('Murid dan paket wajib dipilih sebelum pembayaran Midtrans')
  }

  const payload = {
    id_murid: Number(form.id_murid),
    id_paket: Number(selected.id_paket),
    jumlah_bayar: Math.round(Number(form.jumlah_bayar || 0)),
    nama_paket: selected.nama_paket || 'Paket Bimbingan',
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
    throw new Error('Token Midtrans tidak ditemukan. Periksa konfigurasi server atau coba lagi.')
  }

  await loadMidtransSnapScript()

  await new Promise((resolve, reject) => {
    window.snap.pay(snapToken, {
      onSuccess: async () => {
        await tambahTransaksiApi({
          id_murid: Number(form.id_murid),
          id_paket: Number(selected.id_paket),
          tanggal_mulai: form.tanggal_mulai,
          tanggal_selesai: form.tanggal_selesai,
          tanggal_bayar: toDateInput(new Date().toISOString()),
          jumlah_bayar: Number(form.jumlah_bayar),
          metode_bayar: 'midtrans'
        })

        successMessage.value = 'Pembayaran Midtrans berhasil diproses.'
        showModal.value = false
        resolve()
      },
      onPending: () => {
        successMessage.value = 'Pembayaran Midtrans pending. Silakan selesaikan pembayaran Anda.'
        resolve()
      },
      onError: (result) => {
        reject(new Error(result?.status_message || 'Pembayaran Midtrans gagal'))
      },
      onClose: () => {
        resolve()
      }
    })
  })
}

function openModal(paket) {
  selectedPaket.value = paket
  form.id_murid = ''
  form.tanggal_mulai = ''
  form.tanggal_selesai = ''
  form.jumlah_bayar = Number(paket?.harga || 0)
  isTanggalSelesaiManual.value = false
  errorMessage.value = ''
  successMessage.value = ''
  showModal.value = true
}

function handleTanggalSelesaiInput() {
  isTanggalSelesaiManual.value = true
}

function handleTanggalMulaiInput() {
  isTanggalSelesaiManual.value = false
}

function resetTanggalSelesaiOtomatis() {
  isTanggalSelesaiManual.value = false
  form.tanggal_selesai = hitungTanggalSelesaiOtomatis(form.tanggal_mulai, selectedPaket.value)
}

async function submitTambahPaket() {
  if (!selectedPaket.value) return

  if (!form.id_murid || !form.tanggal_mulai || !form.tanggal_selesai) {
    errorMessage.value = 'Murid, tanggal mulai, dan tanggal selesai wajib diisi.'
    return
  }

  if (activePackageWarning.value) {
    errorMessage.value = activePackageWarning.value
    return
  }

  if (new Date(form.tanggal_selesai) < new Date(form.tanggal_mulai)) {
    errorMessage.value = 'Tanggal selesai tidak boleh lebih kecil dari tanggal mulai.'
    return
  }

  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await startMidtransPayment()
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    submitting.value = false
  }
}

watch(
  () => [selectedPaket.value?.id_paket, form.tanggal_mulai, form.tanggal_selesai],
  () => {
    if (!isTanggalSelesaiManual.value) {
      form.tanggal_selesai = hitungTanggalSelesaiOtomatis(form.tanggal_mulai, selectedPaket.value)
    }
    form.jumlah_bayar = calculateJumlahBayar()
  }
)

async function loadData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [paketResp, muridResp] = await Promise.all([
      getPaketAktifApi(),
      getBimbinganByOrtuApi(authStore.user?.id_user)
    ])
    paketList.value = toArray(paketResp)

    const unikMurid = new Map()
    toArray(muridResp).forEach((item) => {
      const id = Number(item?.id_murid)
      if (!id || unikMurid.has(id)) return
      unikMurid.set(id, {
        id_murid: id,
        nama_murid: item?.nama_murid || `Murid ${id}`
      })
    })
    muridOptions.value = [...unikMurid.values()]

    const transaksiResp = await getTransaksiByOrtuApi(authStore.user?.id_user)
    transaksiByOrtu.value = toArray(transaksiResp)
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <section>
    <h1 class="page-title">Paket Bimbingan</h1>
    <p class="page-subtitle">Pilih paket aktif untuk didaftarkan ke murid Anda.</p>

    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
    <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Paket Bimbingan Aktif</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchPaket" placeholder="Cari paket bimbingan..." />
      </div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Nama Paket</th>
              <th>Harga</th>
              <th>Durasi</th>
              <th>Deskripsi</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredPaket" :key="item.id_paket || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.nama_paket || '-' }}</td>
              <td>{{ formatHarga(item.harga) }}</td>
              <td>{{ item.durasi_bulan || 0 }} bulan {{ item.durasi_hari || 0 }} hari</td>
              <td>{{ item.deskripsi || '-' }}</td>
              <td>
                <button class="btn btn-primary" type="button" @click="openModal(item)">Tambah Paket</button>
              </td>
            </tr>
            <tr v-if="!loading && filteredPaket.length === 0">
              <td colspan="6">Tidak ada paket aktif.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <BaseModal :show="showModal" title="Tambah Paket Bimbingan" @close="showModal = false">
      <form class="form-grid" @submit.prevent="submitTambahPaket">
        <div class="field" style="grid-column: 1 / -1">
          <label>Paket</label>
          <input :value="selectedPaket?.nama_paket || '-'" readonly />
        </div>

        <div class="field">
          <label>Pilih Murid</label>
          <select v-model="form.id_murid" required>
            <option disabled value="">Pilih Murid</option>
            <option v-for="murid in muridOptions" :key="murid.id_murid" :value="murid.id_murid">
              {{ murid.nama_murid }}
            </option>
          </select>
          <small v-if="activePackageWarning" class="field-warning">{{ activePackageWarning }}</small>
        </div>

        <div class="field">
          <label>Tanggal Mulai</label>
          <input v-model="form.tanggal_mulai" type="date" required @input="handleTanggalMulaiInput" />
        </div>

        <div class="field">
          <label>Tanggal Selesai</label>
          <div class="date-input-wrap">
            <input v-model="form.tanggal_selesai" type="date" required @input="handleTanggalSelesaiInput" />
            <button
              class="btn btn-secondary btn-reset-auto"
              type="button"
              @click="resetTanggalSelesaiOtomatis"
              :disabled="!form.tanggal_mulai"
            >
              Reset ke otomatis
            </button>
          </div>
          <small class="field-note">Otomatis dihitung dari tanggal mulai dan durasi paket, tetap bisa Anda ubah.</small>
        </div>

        <div class="field" style="grid-column: 1 / -1">
          <label>Total Bayar</label>
          <input :value="formatHarga(form.jumlah_bayar)" type="text" readonly />
          <small class="field-note">
            Total dihitung dari harga paket dan durasi periode {{ formatDateDisplay(form.tanggal_mulai) }} s/d {{ formatDateDisplay(form.tanggal_selesai) }}.
          </small>
        </div>

        <div class="form-actions">
          <button class="btn btn-secondary" type="button" @click="showModal = false">Batal</button>
          <button class="btn btn-primary" type="submit" :disabled="submitting || Boolean(activePackageWarning)">
            {{ submitting ? 'Memproses...' : 'Lanjut Bayar Midtrans' }}
          </button>
        </div>
      </form>
    </BaseModal>
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
  margin: 0;
  font-size: 1.05rem;
}

.table-tools {
  display: flex;
  gap: 10px;
  justify-content: space-between;
  align-items: center;
  margin: 10px 0;
}

.table-tools input {
  width: 100%;
  max-width: 420px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  grid-column: 1 / -1;
}

.field-note {
  display: block;
  margin-top: 6px;
  color: #64748b;
  font-size: 0.8rem;
}

.field-warning {
  display: block;
  margin-top: 6px;
  color: #b91c1c;
  font-size: 0.8rem;
  font-weight: 600;
}

.date-input-wrap {
  display: flex;
  gap: 8px;
  align-items: center;
}

.date-input-wrap input {
  flex: 1;
}

.btn-reset-auto {
  white-space: nowrap;
  padding: 6px 10px;
  font-size: 0.78rem;
}

@media (max-width: 900px) {
  .table-tools {
    flex-direction: column;
    align-items: stretch;
  }

  .table-tools input {
    max-width: none;
  }

  .date-input-wrap {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
