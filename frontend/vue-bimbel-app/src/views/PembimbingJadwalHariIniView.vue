<script setup>
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { getJadwalByPembimbingApi } from '../services/api'

const authStore = useAuthStore()
const loading = ref(false)
const errorMessage = ref('')
const jadwalList = ref([])
const searchJadwal = ref('')

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function hariIniKey() {
  const hari = new Intl.DateTimeFormat('id-ID', { weekday: 'long' }).format(new Date()).toLowerCase()
  return hari.replace(/\s+/g, '')
}

function normalizeHari(value) {
  return String(value || '').toLowerCase().replace(/\s+/g, '')
}

function formatJam(value) {
  const raw = String(value || '')
  const match = raw.match(/(\d{2}:\d{2})/)
  return match ? match[1] : '-'
}

function formatRentangWaktu(item) {
  return `${formatJam(item?.waktu_mulai)} - ${formatJam(item?.waktu_selesai)}`
}

function formatMuridDisplay(item) {
  const kode = String(item?.kode_murid || '').trim()
  const nama = String(item?.nama_murid || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  return nama || '-'
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

const jadwalHariIni = computed(() => {
  const key = hariIniKey()
  return jadwalList.value.filter((item) => normalizeHari(item?.hari_bimbingan) === key)
})

const filteredJadwalHariIni = computed(() => {
  const keyword = normalizeText(searchJadwal.value).trim()
  if (!keyword) return jadwalHariIni.value

  return jadwalHariIni.value.filter((item) => {
    const joined = [item.hari_bimbingan, formatRentangWaktu(item), item.ruangan, formatMuridDisplay(item)]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

function downloadJadwalHariIniPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredJadwalHariIni.value.map((item, index) => [
    index + 1,
    item.hari_bimbingan || '-',
    formatRentangWaktu(item),
    item.ruangan || '-',
    formatMuridDisplay(item)
  ])

  doc.setFontSize(12)
  doc.text('Jadwal Mengajar Hari Ini', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Hari', 'Waktu', 'Ruangan', 'Murid']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('jadwal-hari-ini.pdf')
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await getJadwalByPembimbingApi(authStore.user?.id_user)
    jadwalList.value = toArray(response)
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
    <h1 class="page-title">Jadwal Mengajar Hari Ini</h1>
    <p class="page-subtitle">Menampilkan daftar bimbingan hari ini untuk pembimbing aktif.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Jadwal Hari Ini</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchJadwal" placeholder="Cari jadwal hari ini..." />
        <div class="tools-actions">
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF jadwal hari ini"
            @click="downloadJadwalHariIniPdf"
          >
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>

      <div class="table-wrap" v-else>
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Hari</th>
              <th>Waktu</th>
              <th>Ruangan</th>
              <th>Murid</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredJadwalHariIni" :key="item.id_jadwal || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.hari_bimbingan || '-' }}</td>
              <td>{{ formatRentangWaktu(item) }}</td>
              <td>{{ item.ruangan || '-' }}</td>
              <td>{{ formatMuridDisplay(item) }}</td>
            </tr>
            <tr v-if="!loading && filteredJadwalHariIni.length === 0">
              <td colspan="5">Tidak ada jadwal mengajar hari ini.</td>
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
