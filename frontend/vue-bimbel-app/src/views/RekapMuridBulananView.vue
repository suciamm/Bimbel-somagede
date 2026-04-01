<script setup>
import { computed, reactive, ref } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import { getRekapMuridBulananApi } from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const showPrintModal = ref(false)
const hasLoadedSummary = ref(false)
const selectedDetailKey = ref('aktif')
const summary = reactive({
  tanggal_mulai: '-',
  tanggal_selesai: '-',
  jumlah_murid_aktif: 0,
  murid_masuk_baru: 0,
  murid_keluar: 0,
  daftar_murid_aktif: [],
  daftar_murid_masuk_baru: [],
  daftar_murid_keluar: []
})

const filter = reactive({
  tanggalMulai: '',
  tanggalSelesai: ''
})

const pdfFilter = reactive({
  bulan: String(new Date().getMonth() + 1),
  tahun: String(new Date().getFullYear())
})

const isPdfFilterValid = computed(() => {
  const bulan = Number(pdfFilter.bulan)
  const tahun = Number(pdfFilter.tahun)

  const bulanValid = Number.isInteger(bulan) && bulan >= 1 && bulan <= 12
  const tahunValid = Number.isInteger(tahun) && tahun >= 2000 && tahun <= 2100

  return bulanValid && tahunValid
})

const detailRows = computed(() => {
  if (selectedDetailKey.value === 'masuk') return summary.daftar_murid_masuk_baru
  if (selectedDetailKey.value === 'keluar') return summary.daftar_murid_keluar
  return summary.daftar_murid_aktif
})

const detailTitle = computed(() => {
  if (selectedDetailKey.value === 'masuk') return 'Detail Murid Masuk Baru'
  if (selectedDetailKey.value === 'keluar') return 'Detail Murid Keluar'
  return 'Detail Murid Aktif'
})

function toDateOnly(value) {
  const iso = String(value || '').match(/\d{4}-\d{2}-\d{2}/)?.[0]
  if (!iso) return '-'
  const [year, month, day] = iso.split('-')
  return `${day}-${month}-${year}`
}

function applySummary(data) {
  summary.tanggal_mulai = data?.tanggal_mulai || '-'
  summary.tanggal_selesai = data?.tanggal_selesai || '-'
  summary.jumlah_murid_aktif = Number(data?.jumlah_murid_aktif || 0)
  summary.murid_masuk_baru = Number(data?.murid_masuk_baru || 0)
  summary.murid_keluar = Number(data?.murid_keluar || 0)
  summary.daftar_murid_aktif = Array.isArray(data?.daftar_murid_aktif)
    ? data.daftar_murid_aktif.filter((item) => item && typeof item === 'object')
    : []
  summary.daftar_murid_masuk_baru = Array.isArray(data?.daftar_murid_masuk_baru)
    ? data.daftar_murid_masuk_baru.filter((item) => item && typeof item === 'object')
    : []
  summary.daftar_murid_keluar = Array.isArray(data?.daftar_murid_keluar)
    ? data.daftar_murid_keluar.filter((item) => item && typeof item === 'object')
    : []
}

async function loadSummary() {
  loading.value = true
  errorMessage.value = ''
  hasLoadedSummary.value = false

  try {
    const response = await getRekapMuridBulananApi({
      tanggalMulai: filter.tanggalMulai,
      tanggalSelesai: filter.tanggalSelesai
    })

    applySummary(response?.data || response)
    selectedDetailKey.value = 'aktif'
    hasLoadedSummary.value = true
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

async function printPdfByMonth() {
  errorMessage.value = ''

  if (!isPdfFilterValid.value) {
    errorMessage.value = 'Bulan harus 1-12 dan tahun harus valid sebelum print PDF.'
    return
  }

  try {
    const response = await getRekapMuridBulananApi({
      bulan: pdfFilter.bulan,
      tahun: pdfFilter.tahun
    })

    const data = response?.data || response || {}

    const doc = new jsPDF({ orientation: 'portrait' })
    doc.setFontSize(13)
    doc.text('Laporan Rekap Murid Bulanan', 14, 14)
    doc.setFontSize(10)
    doc.text(`Periode: ${data.tanggal_mulai || '-'} s/d ${data.tanggal_selesai || '-'}`, 14, 21)

    autoTable(doc, {
      startY: 28,
      head: [['Metrik', 'Jumlah']],
      body: [
        ['Jumlah Murid Aktif', Number(data.jumlah_murid_aktif || 0)],
        ['Murid Masuk Baru', Number(data.murid_masuk_baru || 0)],
        ['Murid Keluar', Number(data.murid_keluar || 0)]
      ],
      styles: { fontSize: 10 }
    })

    doc.save(`rekap-murid-${pdfFilter.tahun}-${String(pdfFilter.bulan).padStart(2, '0')}.pdf`)
    showPrintModal.value = false
  } catch (error) {
    errorMessage.value = error.message
  }
}

function openPrintModal() {
  showPrintModal.value = true
}

function showDetail(key) {
  selectedDetailKey.value = key
}
</script>

<template>
  <section>
    <h1 class="page-title">Rekap Murid Bulanan</h1>
    <p class="page-subtitle">Ringkasan jumlah murid aktif, murid masuk baru, dan murid keluar per periode.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Filter Tanggal</h2>
      </header>

      <div class="filter-row">
        <div class="field-inline">
          <label>Tanggal Mulai</label>
          <input v-model="filter.tanggalMulai" type="date" />
        </div>

        <div class="field-inline">
          <label>Tanggal Selesai</label>
          <input v-model="filter.tanggalSelesai" type="date" />
        </div>

        <div class="tools-actions">
          <button class="btn btn-secondary" type="button" @click="loadSummary">Terapkan Filter</button>
          <button class="btn btn-secondary btn-icon btn-pdf" type="button" @click="openPrintModal">
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>

      <p v-if="!loading && !hasLoadedSummary" class="empty-detail">
        Pilih tanggal mulai dan tanggal selesai, lalu klik Terapkan Filter untuk melihat rekap.
      </p>

      <div class="stats-grid" v-if="!loading && hasLoadedSummary">
        <article class="stat-card mini">
          <p class="label">Murid Aktif</p>
          <p class="value">{{ summary.jumlah_murid_aktif }}</p>
          <button
            class="btn btn-secondary btn-detail"
            type="button"
            :class="selectedDetailKey === 'aktif' ? 'active' : ''"
            @click="showDetail('aktif')"
          >
            Detail
          </button>
        </article>

        <article class="stat-card mini">
          <p class="label">Murid Masuk Baru</p>
          <p class="value">{{ summary.murid_masuk_baru }}</p>
          <button
            class="btn btn-secondary btn-detail"
            type="button"
            :class="selectedDetailKey === 'masuk' ? 'active' : ''"
            @click="showDetail('masuk')"
          >
            Detail
          </button>
        </article>

        <article class="stat-card mini">
          <p class="label">Murid Keluar</p>
          <p class="value">{{ summary.murid_keluar }}</p>
          <button
            class="btn btn-secondary btn-detail"
            type="button"
            :class="selectedDetailKey === 'keluar' ? 'active' : ''"
            @click="showDetail('keluar')"
          >
            Detail
          </button>
        </article>
      </div>

      <p class="period-info" v-if="!loading && hasLoadedSummary">
        Periode data: <strong>{{ summary.tanggal_mulai }}</strong> sampai <strong>{{ summary.tanggal_selesai }}</strong>
      </p>

      <section class="panel detail-table-wrap" v-if="!loading && hasLoadedSummary">
        <header class="block-header">
          <h3>{{ detailTitle }}</h3>
        </header>

        <div class="table-wrap">
          <table>
            <thead>
              <tr>
                <th style="width: 80px">No</th>
                <th style="width: 160px">Kode Murid</th>
                <th>Nama Murid</th>
                <th style="width: 150px">Tgl Masuk</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(item, index) in detailRows"
                :key="`${selectedDetailKey}-${index}-${item.kode_murid || item.nama_murid}`"
              >
                <td>{{ index + 1 }}</td>
                <td>{{ item.kode_murid || '-' }}</td>
                <td>{{ item.nama_murid || '-' }}</td>
                <td>{{ toDateOnly(item.tgl_masuk) }}</td>
              </tr>
              <tr v-if="detailRows.length === 0">
                <td colspan="4">Tidak ada data untuk kategori ini.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </section>

    <BaseModal :show="showPrintModal" title="Print PDF Rekap Murid Bulanan" @close="showPrintModal = false">
      <form class="form-grid" @submit.prevent="printPdfByMonth">
        <div class="field-inline short">
          <label>Bulan</label>
          <input v-model="pdfFilter.bulan" type="number" min="1" max="12" placeholder="1-12" required />
        </div>

        <div class="field-inline short">
          <label>Tahun</label>
          <input v-model="pdfFilter.tahun" type="number" min="2000" max="2100" placeholder="YYYY" required />
        </div>

        <div>
          <button class="btn btn-primary" type="submit" :disabled="!isPdfFilterValid">Print PDF</button>
        </div>
      </form>

      <p class="modal-hint">Isi bulan dan tahun laporan yang ingin dicetak.</p>
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

.filter-row {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: end;
}

.tools-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.field-inline {
  display: grid;
  gap: 6px;
  min-width: 210px;
}

.field-inline.short {
  min-width: 120px;
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

.modal-hint {
  margin-top: 10px;
  color: var(--text-muted);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.stats-grid {
  margin-top: 14px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px;
}

.stat-card.mini {
  padding: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #ffffff;
}

.stat-card .label {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.stat-card .value {
  margin: 6px 0 0;
  font-size: 1.7rem;
  font-weight: 700;
  color: #0f172a;
}

.btn-detail {
  margin-top: 10px;
  padding: 4px 10px;
  font-size: 0.8rem;
  line-height: 1.2;
  min-height: auto;
}

.btn-detail.active {
  background: #1e3a8a;
  border-color: #1e3a8a;
  color: #fff;
}

.period-info {
  margin-top: 12px;
  color: #334155;
}

.detail-table-wrap {
  margin-top: 16px;
  padding: 12px;
  border: 1px solid #e2e8f0;
}

.detail-table-wrap h3 {
  margin: 0 0 10px;
  font-size: 0.95rem;
  color: #0f172a;
}

.empty-detail {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.9rem;
}
</style>
