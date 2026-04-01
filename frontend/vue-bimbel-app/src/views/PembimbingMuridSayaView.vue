<script setup>
import { computed, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { useAuthStore } from '../stores/auth'
import { getMuridByPembimbingApi } from '../services/api'

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

function formatTanggal(value) {
  const raw = String(value || '').trim()
  if (!raw) return '-'
  const datePart = raw.split('T')[0]
  if (!/^\d{4}-\d{2}-\d{2}$/.test(datePart)) return raw
  const [year, month, day] = datePart.split('-')
  return `${day}-${month}-${year}`
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

const filteredMuridList = computed(() => {
  const keyword = normalizeText(searchMurid.value).trim()
  if (!keyword) return muridList.value

  return muridList.value.filter((item) => {
    const joined = [item.kode_murid, item.nama_murid, item.tgl_lahir, item.tgl_masuk, item.status_murid]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

function downloadMuridSayaPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredMuridList.value.map((item, index) => [
    index + 1,
    item.kode_murid || '-',
    item.nama_murid || '-',
    formatTanggal(item.tgl_lahir),
    formatTanggal(item.tgl_masuk),
    item.status_murid || '-'
  ])

  doc.setFontSize(12)
  doc.text('Data Murid Saya', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Kode Murid', 'Nama Murid', 'Tanggal Lahir', 'Tanggal Masuk', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('data-murid-saya.pdf')
}

async function loadData() {
  const idUser = Number(authStore.user?.id_user)
  if (!idUser) {
    muridList.value = []
    errorMessage.value = 'Sesi pembimbing belum siap. Silakan muat ulang halaman.'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await getMuridByPembimbingApi(idUser)
    muridList.value = toArray(response)
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(
  () => authStore.user?.id_user,
  (id) => {
    if (id) loadData()
  },
  { immediate: true }
)
</script>

<template>
  <section>
    <h1 class="page-title">Data Murid Saya</h1>
    <p class="page-subtitle">Menampilkan daftar murid aktif yang menjadi tanggung jawab pembimbing.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Murid Aktif</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchMurid" placeholder="Cari data murid..." />
        <div class="tools-actions">
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF data murid"
            @click="downloadMuridSayaPdf"
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
              <th>Kode Murid</th>
              <th>Nama Murid</th>
              <th>Tanggal Lahir</th>
              <th>Tanggal Masuk</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredMuridList" :key="item.id_murid || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.kode_murid || '-' }}</td>
              <td>{{ item.nama_murid || '-' }}</td>
              <td>{{ formatTanggal(item.tgl_lahir) }}</td>
              <td>{{ formatTanggal(item.tgl_masuk) }}</td>
              <td>
                <span class="pill" :class="item.status_murid === 'aktif' ? 'pill-success' : 'pill-warning'">
                  {{ item.status_murid || '-' }}
                </span>
              </td>
            </tr>
            <tr v-if="!loading && filteredMuridList.length === 0">
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
