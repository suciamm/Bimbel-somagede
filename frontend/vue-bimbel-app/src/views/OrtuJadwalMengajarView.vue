<script setup>
import { computed, onMounted, ref } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { getJadwalApi } from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const jadwalList = ref([])
const searchJadwal = ref('')

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function formatJam(value) {
  const raw = String(value || '')
  const match = raw.match(/(\d{2}:\d{2})/)
  return match ? match[1] : '-'
}

function formatMuridDisplay(item) {
  const kode = String(item?.kode_murid || '').trim()
  const nama = String(item?.nama_murid || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  return nama || '-'
}

function formatPembimbing(item) {
  const kode = String(item?.kode_pembimbing || '').trim()
  const nama = String(item?.nama_pembimbing || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  return nama || '-'
}

function formatWaktu(item) {
  return `${formatJam(item?.waktu_mulai)} - ${formatJam(item?.waktu_selesai)}`
}

const filteredJadwal = computed(() => {
  const keyword = normalizeText(searchJadwal.value).trim()
  if (!keyword) return jadwalList.value

  return jadwalList.value.filter((item) => {
    const text = [
      item?.hari_bimbingan,
      formatWaktu(item),
      item?.ruangan,
      formatMuridDisplay(item),
      formatPembimbing(item)
    ]
      .map(normalizeText)
      .join(' ')
    return text.includes(keyword)
  })
})

function downloadPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredJadwal.value.map((item, index) => [
    index + 1,
    item?.hari_bimbingan || '-',
    formatWaktu(item),
    item?.ruangan || '-',
    formatPembimbing(item),
    formatMuridDisplay(item)
  ])

  doc.setFontSize(12)
  doc.text('Jadwal Mengajar', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Hari', 'Waktu', 'Ruangan', 'Pembimbing', 'Murid']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('orangtua-jadwal-mengajar.pdf')
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await getJadwalApi()
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
    <h1 class="page-title">Jadwal Mengajar</h1>
    <p class="page-subtitle">Daftar semua jadwal bimbingan siswa.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Jadwal</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchJadwal" placeholder="Cari jadwal..." />
        <div class="tools-actions">
          <button class="btn btn-secondary btn-icon btn-pdf" type="button" @click="downloadPdf">
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
              <th>Pembimbing</th>
              <th>Murid</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredJadwal" :key="item.id_jadwal || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.hari_bimbingan || '-' }}</td>
              <td>{{ formatWaktu(item) }}</td>
              <td>{{ item.ruangan || '-' }}</td>
              <td>{{ formatPembimbing(item) }}</td>
              <td>{{ formatMuridDisplay(item) }}</td>
            </tr>
            <tr v-if="!loading && filteredJadwal.length === 0">
              <td colspan="6">Tidak ada data jadwal.</td>
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
