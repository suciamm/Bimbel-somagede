<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import { createAbsensiApi, getAbsensiRekapApi } from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const absensiList = ref([])
const showCreateModal = ref(false)
const searchAbsensi = ref('')
const currentPageAbsensi = ref(1)
const PAGE_SIZE = 10

const filter = reactive({
  bulan: String(new Date().getMonth() + 1),
  tahun: String(new Date().getFullYear())
})

const form = reactive({
  id_jadwal: '',
  id_murid: '',
  id_pembimbing: '',
  tanggal_sesi: '',
  status_hadir: 'hadir',
  keterangan: ''
})

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function getTanggal(item) {
  return (item.tanggal_sesi || item.TanggalSesi || '').slice?.(0, 10) || '-'
}

function getStatus(item) {
  return item.status_hadir || item.StatusHadir || '-'
}

const filteredAbsensiList = computed(() => {
  const keyword = normalizeText(searchAbsensi.value).trim()
  if (!keyword) return absensiList.value

  return absensiList.value.filter((item) => {
    const joined = [
      item.id_absensi || item.IDAbsensi,
      item.id_murid || item.IDMurid,
      item.id_pembimbing || item.IDPembimbing,
      getTanggal(item),
      getStatus(item),
      item.keterangan || item.Keterangan
    ]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const totalPagesAbsensi = computed(() => Math.max(1, Math.ceil(filteredAbsensiList.value.length / PAGE_SIZE)))

const paginatedAbsensiList = computed(() => {
  const start = (currentPageAbsensi.value - 1) * PAGE_SIZE
  return filteredAbsensiList.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(index) {
  return (currentPageAbsensi.value - 1) * PAGE_SIZE + index + 1
}

function goPrevPage() {
  if (currentPageAbsensi.value > 1) currentPageAbsensi.value -= 1
}

function goNextPage() {
  if (currentPageAbsensi.value < totalPagesAbsensi.value) currentPageAbsensi.value += 1
}

function downloadAbsensiPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredAbsensiList.value.map((item, index) => [
    index + 1,
    item.id_absensi || item.IDAbsensi || '-',
    item.id_murid || item.IDMurid || '-',
    item.id_pembimbing || item.IDPembimbing || '-',
    getTanggal(item),
    getStatus(item),
    item.keterangan || item.Keterangan || '-'
  ])

  doc.setFontSize(12)
  doc.text('Rekap Absensi', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'ID Absensi', 'ID Murid', 'ID Pembimbing', 'Tanggal', 'Status', 'Keterangan']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('rekap-absensi.pdf')
}

async function loadRekap() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await getAbsensiRekapApi({
      bulan: filter.bulan,
      tahun: filter.tahun
    })
    absensiList.value = response?.data || response || []
    currentPageAbsensi.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(searchAbsensi, () => {
  currentPageAbsensi.value = 1
})

watch(filteredAbsensiList, () => {
  if (currentPageAbsensi.value > totalPagesAbsensi.value) {
    currentPageAbsensi.value = totalPagesAbsensi.value
  }
})

async function submitAbsensi() {
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await createAbsensiApi(Number(form.id_jadwal), {
      id_murid: Number(form.id_murid),
      id_pembimbing: Number(form.id_pembimbing),
      tanggal_sesi: form.tanggal_sesi,
      status_hadir: form.status_hadir,
      keterangan: form.keterangan
    })

    successMessage.value = 'Data absensi berhasil dicatat.'
    form.id_jadwal = ''
    form.id_murid = ''
    form.id_pembimbing = ''
    form.tanggal_sesi = ''
    form.status_hadir = 'hadir'
    form.keterangan = ''
    showCreateModal.value = false

    await loadRekap()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(loadRekap)
</script>

<template>
  <section>
    <h1 class="page-title">Rekap Absensi</h1>
    <p class="page-subtitle">Input kehadiran sesi dan lihat rekap absensi berdasarkan bulan/tahun.</p>

    <BaseModal :show="showCreateModal" title="Input Data Absensi" @close="showCreateModal = false">
      <form class="form-grid" @submit.prevent="submitAbsensi">
        <div class="field">
          <label>ID Jadwal</label>
          <input v-model="form.id_jadwal" type="number" min="1" required />
        </div>

        <div class="field">
          <label>ID Murid</label>
          <input v-model="form.id_murid" type="number" min="1" required />
        </div>

        <div class="field">
          <label>ID Pembimbing</label>
          <input v-model="form.id_pembimbing" type="number" min="1" required />
        </div>

        <div class="field">
          <label>Tanggal Sesi</label>
          <input v-model="form.tanggal_sesi" type="date" required />
        </div>

        <div class="field">
          <label>Status Kehadiran</label>
          <select v-model="form.status_hadir">
            <option value="hadir">hadir</option>
            <option value="izin">izin</option>
            <option value="alpa">alpa</option>
          </select>
        </div>

        <div class="field" style="grid-column: 1 / -1">
          <label>Keterangan</label>
          <textarea v-model="form.keterangan" rows="2" />
        </div>

        <div>
          <button class="btn btn-primary" type="submit">Simpan Absensi</button>
        </div>
      </form>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>
    </BaseModal>

    <section class="panel block">
      <header class="block-header">
        <h2>Data Rekap Absensi</h2>
        <div class="filter-row">
          <input v-model="filter.bulan" type="number" min="1" max="12" placeholder="Bulan" />
          <input v-model="filter.tahun" type="number" min="2020" placeholder="Tahun" />
          <button class="btn btn-secondary" @click="loadRekap">Terapkan Filter</button>
        </div>
      </header>

      <div class="table-tools">
        <input v-model="searchAbsensi" placeholder="Cari data absensi..." />
        <div class="tools-actions">
          <button class="btn btn-primary" type="button" @click="showCreateModal = true">Tambah Absensi</button>
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF rekap absensi"
            @click="downloadAbsensiPdf"
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
              <th>ID Absensi</th>
              <th>ID Murid</th>
              <th>ID Pembimbing</th>
              <th>Tanggal</th>
              <th>Status</th>
              <th>Keterangan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedAbsensiList" :key="item.id_absensi || item.IDAbsensi || index">
              <td>{{ nomorUrut(index) }}</td>
              <td>{{ item.id_absensi || item.IDAbsensi || '-' }}</td>
              <td>{{ item.id_murid || item.IDMurid || '-' }}</td>
              <td>{{ item.id_pembimbing || item.IDPembimbing || '-' }}</td>
              <td>{{ (item.tanggal_sesi || item.TanggalSesi || '').slice?.(0, 10) || '-' }}</td>
              <td>
                <span class="pill" :class="(item.status_hadir || item.StatusHadir) === 'hadir' ? 'pill-success' : (item.status_hadir || item.StatusHadir) === 'izin' ? 'pill-warning' : 'pill-danger'">
                  {{ item.status_hadir || item.StatusHadir || '-' }}
                </span>
              </td>
              <td>{{ item.keterangan || item.Keterangan || '-' }}</td>
            </tr>
            <tr v-if="paginatedAbsensiList.length === 0">
              <td colspan="7">Tidak ada data absensi.</td>
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
          :disabled="currentPageAbsensi === 1"
        >
          <<
        </button>
        <span>Halaman {{ currentPageAbsensi }} / {{ totalPagesAbsensi }} ({{ filteredAbsensiList.length }} data)</span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="goNextPage"
          :disabled="currentPageAbsensi === totalPagesAbsensi"
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

.filter-row {
  display: flex;
  gap: 8px;
}

.filter-row input {
  width: 110px;
  border: 1px solid var(--border);
  background: var(--surface-soft);
  border-radius: 10px;
  padding: 8px 10px;
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
  .block-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .filter-row {
    width: 100%;
    flex-wrap: wrap;
  }

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
