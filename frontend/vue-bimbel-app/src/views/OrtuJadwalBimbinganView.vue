<script setup>
import { computed, onMounted, ref } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { useAuthStore } from '../stores/auth'
import { getBimbinganByOrtuApi } from '../services/api'

const authStore = useAuthStore()
const loading = ref(false)
const errorMessage = ref('')
const bimbinganList = ref([])
const searchBimbingan = ref('')

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

function formatWaktu(item) {
  return `${formatJam(item?.waktu_mulai)} - ${formatJam(item?.waktu_selesai)}`
}

const filteredBimbingan = computed(() => {
  const keyword = normalizeText(searchBimbingan.value).trim()
  if (!keyword) return bimbinganList.value

  return bimbinganList.value.filter((item) => {
    const text = [item?.nama_murid, item?.hari_bimbingan, formatWaktu(item), item?.ruangan, item?.nama_pembimbing]
      .map(normalizeText)
      .join(' ')
    return text.includes(keyword)
  })
})

function downloadPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredBimbingan.value.map((item, index) => [
    index + 1,
    item?.nama_murid || '-',
    item?.hari_bimbingan || '-',
    formatWaktu(item),
    item?.ruangan || '-',
    item?.nama_pembimbing || '-'
  ])

  doc.setFontSize(12)
  doc.text('Jadwal Bimbingan Anak', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Murid', 'Hari', 'Waktu', 'Ruangan', 'Pembimbing']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('orangtua-jadwal-bimbingan.pdf')
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await getBimbinganByOrtuApi(authStore.user?.id_user)
    bimbinganList.value = toArray(response)
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
    <h1 class="page-title">Jadwal Bimbingan</h1>
    <p class="page-subtitle">Daftar bimbingan anak orangtua sesuai akun login.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Bimbingan Anak</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchBimbingan" placeholder="Cari jadwal bimbingan..." />
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
              <th>Murid</th>
              <th>Hari</th>
              <th>Waktu</th>
              <th>Ruangan</th>
              <th>Pembimbing</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in filteredBimbingan"
              :key="`${item.id_murid || 'm'}-${item.hari_bimbingan || 'h'}-${item.waktu_mulai || 'wm'}-${index}`"
            >
              <td>{{ index + 1 }}</td>
              <td>{{ item.nama_murid || '-' }}</td>
              <td>{{ item.hari_bimbingan || '-' }}</td>
              <td>{{ formatWaktu(item) }}</td>
              <td>{{ item.ruangan || '-' }}</td>
              <td>{{ item.nama_pembimbing || '-' }}</td>
            </tr>
            <tr v-if="!loading && filteredBimbingan.length === 0">
              <td colspan="6">Tidak ada data jadwal bimbingan.</td>
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
