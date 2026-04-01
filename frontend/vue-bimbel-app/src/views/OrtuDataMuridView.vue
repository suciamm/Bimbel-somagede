<script setup>
import { computed, onMounted, ref } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { getMuridByOrtuApi } from '../services/api'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const errorMessage = ref('')
const muridList = ref([])
const searchMurid = ref('')

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

const filteredMurid = computed(() => {
  const keyword = normalizeText(searchMurid.value).trim()
  if (!keyword) return muridList.value

  return muridList.value.filter((item) => {
    const text = [item?.kode_murid, item?.nama_murid, item?.nama_ortu, item?.status_murid].map(normalizeText).join(' ')
    return text.includes(keyword)
  })
})

function downloadPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredMurid.value.map((item, index) => [
    index + 1,
    item?.kode_murid || '-',
    item?.nama_murid || '-',
    item?.nama_ortu || '-',
    toDateOnly(item?.tgl_masuk),
    item?.status_murid || '-'
  ])

  doc.setFontSize(12)
  doc.text('Data Murid Aktif', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Kode', 'Nama Murid', 'Nama Orang Tua', 'Tgl Masuk', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('orangtua-data-murid.pdf')
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''
  try {
    const idUser = authStore.user?.id_user
    if (!idUser) {
      throw new Error('ID pengguna tidak ditemukan')
    }
    const response = await getMuridByOrtuApi(idUser)
    muridList.value = toArray(response)
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
    <h1 class="page-title">Data Murid</h1>
    <p class="page-subtitle">Daftar murid aktif Anda.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Murid Aktif</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchMurid" placeholder="Cari murid..." />
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
              <th>Kode</th>
              <th>Nama Murid</th>
              <th>Nama Orang Tua</th>
              <th>Tgl Masuk</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredMurid" :key="item.id_murid || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.kode_murid || '-' }}</td>
              <td>{{ item.nama_murid || '-' }}</td>
              <td>{{ item.nama_ortu || '-' }}</td>
              <td>{{ toDateOnly(item.tgl_masuk) }}</td>
              <td>{{ item.status_murid || '-' }}</td>
            </tr>
            <tr v-if="!loading && filteredMurid.length === 0">
              <td colspan="6">Tidak ada data murid aktif.</td>
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
