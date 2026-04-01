<script setup>
import { computed, onMounted, ref } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { useAuthStore } from '../stores/auth'
import { getTransaksiByOrtuApi, konfirmasiPerpanjangMidtransApi, perpanjangPaketMidtransApi } from '../services/api'

const authStore = useAuthStore()
const loading = ref(false)
const processingId = ref(null)
const errorMessage = ref('')
const successMessage = ref('')
const transaksiList = ref([])
const searchTransaksi = ref('')
const MIDTRANS_CLIENT_KEY = import.meta.env.VITE_MIDTRANS_CLIENT_KEY || 'Mid-client-7r84InjwqcNGIMqF'

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function toDateOnly(value) {
  const raw = String(value || '')
  const match = raw.match(/\d{4}-\d{2}-\d{2}/)
  if (!match) return '-'
  const [year, month, day] = match[0].split('-')
  return `${day}-${month}-${year}`
}

function formatHarga(value) {
  const number = Number(value || 0)
  return `Rp ${number.toLocaleString('id-ID')}`
}

function formatMetodeBayar(value) {
  const raw = String(value || '').trim().toLowerCase()
  if (!raw) return '-'
  if (raw.startsWith('midtrans')) return 'midtrans'
  return raw
}

function downloadPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredTransaksi.value.map((item, index) => [
    index + 1,
    item?.nama_murid || '-',
    item?.nama_paket || '-',
    `${toDateOnly(item?.tgl_mulai)} s/d ${toDateOnly(item?.tgl_selesai)}`,
    formatMetodeBayar(item?.metode_bayar),
    formatHarga(item?.jumlah_bayar),
    item?.status || '-'
  ])

  doc.setFontSize(12)
  doc.text('Paket Bimbingan Saya', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Murid', 'Paket', 'Periode', 'Metode', 'Bayar', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('orangtua-paket-bimbingan-saya.pdf')
}

const filteredTransaksi = computed(() => {
  const keyword = normalizeText(searchTransaksi.value).trim()
  if (!keyword) return transaksiList.value

  return transaksiList.value.filter((item) => {
    const text = [item?.nama_murid, item?.nama_paket, formatMetodeBayar(item?.metode_bayar), item?.status]
      .map(normalizeText)
      .join(' ')
    return text.includes(keyword)
  })
})

async function handlePerpanjang(item) {
  processingId.value = item.id_transaksi
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const response = await perpanjangPaketMidtransApi(item.id_paket, {
      id_murid: Number(item.id_murid),
      nama_user: authStore.user?.nama_lengkap || 'Orangtua',
      nama_paket: item.nama_paket || ''
    })

    const token = response?.snap_token || response?.token
    const orderId = response?.order_id || response?.data?.order_id || ''
    const jumlahBayar = Number(response?.harga || item?.jumlah_bayar || 0)

    if (!token) {
      throw new Error('Token Midtrans tidak tersedia untuk perpanjangan')
    }

    await loadMidtransSnapScript()

    await new Promise((resolve, reject) => {
      window.snap.pay(token, {
        onSuccess: async () => {
          await konfirmasiPerpanjangMidtransApi(item.id_paket, {
            id_murid: Number(item.id_murid),
            order_id: orderId,
            jumlah_bayar: jumlahBayar
          })

          successMessage.value = 'Perpanjangan berhasil. Periode paket telah diperbarui.'
          await loadData()
          resolve()
        },
        onPending: () => {
          successMessage.value = 'Pembayaran perpanjang masih pending. Silakan selesaikan transaksi Anda.'
          resolve()
        },
        onError: (result) => {
          reject(new Error(result?.status_message || 'Pembayaran perpanjang gagal'))
        },
        onClose: () => {
          resolve()
        }
      })
    })
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    processingId.value = null
  }
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

async function loadData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await getTransaksiByOrtuApi(authStore.user?.id_user)
    transaksiList.value = toArray(response)
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
    <h1 class="page-title">Paket Bimbingan Saya</h1>
    <p class="page-subtitle">Daftar paket bimbingan yang sedang atau pernah diikuti.</p>

    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
    <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Riwayat Paket Bimbingan</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchTransaksi" placeholder="Cari paket saya..." />
        <div class="tools-actions">
          <button class="btn btn-secondary btn-icon btn-pdf" type="button" @click="downloadPdf">
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
              <th>Metode</th>
              <th>Bayar</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredTransaksi" :key="item.id_transaksi || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.nama_murid || '-' }}</td>
              <td>{{ item.nama_paket || '-' }}</td>
              <td>{{ toDateOnly(item.tgl_mulai) }} s/d {{ toDateOnly(item.tgl_selesai) }}</td>
              <td>{{ formatMetodeBayar(item.metode_bayar) }}</td>
              <td>{{ formatHarga(item.jumlah_bayar) }}</td>
              <td>{{ item.status || '-' }}</td>
              <td>
                <button class="btn btn-secondary" type="button" @click="handlePerpanjang(item)" :disabled="processingId === item.id_transaksi">
                  {{ processingId === item.id_transaksi ? 'Proses...' : 'Perpanjang' }}
                </button>
              </td>
            </tr>
            <tr v-if="!loading && filteredTransaksi.length === 0">
              <td colspan="8">Tidak ada data paket bimbingan.</td>
            </tr>
          </tbody>
        </table>
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

.tools-actions {
  display: flex;
  gap: 8px;
  align-items: center;
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
