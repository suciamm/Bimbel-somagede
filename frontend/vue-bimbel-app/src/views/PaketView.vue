<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import { deletePaketApi, editPaketApi, getPaketApi, tambahPaketApi } from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const paketList = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const searchPaket = ref('')
const currentPagePaket = ref(1)
const PAGE_SIZE = 10

const form = reactive({
  nama_paket: '',
  harga: '',
  durasi_hari: '',
  durasi_bulan: '',
  deskripsi: ''
})

const editForm = reactive({
  id_paket: '',
  nama_paket: '',
  harga: '',
  durasi_hari: 0,
  durasi_bulan: '',
  deskripsi: '',
  status: 'aktif'
})

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function formatDurasiPart(value, unit) {
  const raw = String(value ?? '').trim()
  if (!raw) return ''

  if (/^\d+$/.test(raw)) {
    if (Number(raw) <= 0) return ''
    return `${raw} ${unit}`
  }

  return raw
    .replace(/(\d)([A-Za-z])/g, '$1 $2')
    .replace(/([A-Za-z])(\d)/g, '$1 $2')
    .replace(/\s+/g, ' ')
    .trim()
}

function formatDurasi(item) {
  const bulan = formatDurasiPart(item?.durasi_bulan, 'bulan')
  const hari = formatDurasiPart(item?.durasi_hari, 'hari')
  const result = [bulan, hari].filter(Boolean).join(' ')
  return result || '-'
}

const filteredPaketList = computed(() => {
  const keyword = normalizeText(searchPaket.value).trim()
  if (!keyword) return paketList.value

  return paketList.value.filter((item) => {
    const joined = [item.nama_paket, item.harga, item.durasi_hari, item.durasi_bulan, item.status, item.deskripsi]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const totalPagesPaket = computed(() => Math.max(1, Math.ceil(filteredPaketList.value.length / PAGE_SIZE)))

const paginatedPaketList = computed(() => {
  const start = (currentPagePaket.value - 1) * PAGE_SIZE
  return filteredPaketList.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(index) {
  return (currentPagePaket.value - 1) * PAGE_SIZE + index + 1
}

function goPrevPage() {
  if (currentPagePaket.value > 1) currentPagePaket.value -= 1
}

function goNextPage() {
  if (currentPagePaket.value < totalPagesPaket.value) currentPagePaket.value += 1
}

function downloadPaketPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredPaketList.value.map((item, index) => [
    index + 1,
    item.nama_paket || '-',
    `Rp ${Number(item.harga || 0).toLocaleString('id-ID')}`,
    formatDurasi(item),
    item.status || '-'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Paket Bimbingan', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Nama Paket', 'Harga', 'Durasi', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-paket.pdf')
}

async function loadPaket() {
  loading.value = true
  errorMessage.value = ''
  try {
    const response = await getPaketApi()
    paketList.value = response?.data || response || []
    currentPagePaket.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(searchPaket, () => {
  currentPagePaket.value = 1
})

watch(filteredPaketList, () => {
  if (currentPagePaket.value > totalPagesPaket.value) {
    currentPagePaket.value = totalPagesPaket.value
  }
})

async function submitPaket() {
  successMessage.value = ''
  errorMessage.value = ''

  try {
    await tambahPaketApi({
      nama_paket: form.nama_paket,
      harga: Number(form.harga),
      durasi_hari: Number(form.durasi_hari),
      durasi_bulan: Number(form.durasi_bulan),
      deskripsi: form.deskripsi
    })

    successMessage.value = 'Paket berhasil ditambahkan.'
    form.nama_paket = ''
    form.harga = ''
    form.durasi_hari = ''
    form.durasi_bulan = ''
    form.deskripsi = ''
    showCreateModal.value = false

    await loadPaket()
  } catch (error) {
    errorMessage.value = error.message
  }
}

function selectPaket(item) {
  editForm.id_paket = String(item.id_paket || item.IdPaket || '')
  editForm.nama_paket = item.nama_paket || item.NamaPaket || ''
  editForm.harga = item.harga || item.Harga || 0
  editForm.durasi_hari = item.durasi_hari || item.DurasiHari || 0
  editForm.durasi_bulan = item.durasi_bulan || item.DurasiBulan || 1
  editForm.deskripsi = item.deskripsi || item.Deskripsi || ''
  editForm.status = item.status || item.Status || 'aktif'
  showEditModal.value = true
}

async function submitEditPaket() {
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await editPaketApi(editForm.id_paket, {
      nama_paket: editForm.nama_paket,
      harga: Number(editForm.harga),
      durasi_hari: Number(editForm.durasi_hari),
      durasi_bulan: Number(editForm.durasi_bulan),
      deskripsi: editForm.deskripsi,
      status: editForm.status
    })

    successMessage.value = 'Paket berhasil diperbarui.'
    showEditModal.value = false
    await loadPaket()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function handleDeletePaket(idPaket) {
  if (!window.confirm('Hapus paket bimbingan ini?')) return

  errorMessage.value = ''
  successMessage.value = ''

  try {
    await deletePaketApi(idPaket)
    successMessage.value = 'Paket berhasil dihapus.'
    await loadPaket()
  } catch (error) {
    // Error message from backend includes detailed information
    errorMessage.value = error.message || 'Gagal menghapus paket bimbingan'
  }
}

onMounted(loadPaket)
</script>

<template>
  <section>
    <h1 class="page-title">Paket Bimbingan</h1>
    <p class="page-subtitle">Kelola produk bimbingan yang ditawarkan ke orang tua murid.</p>

    <BaseModal :show="showCreateModal" title="Tambah Paket Bimbingan" @close="showCreateModal = false">
      <form class="form-grid" @submit.prevent="submitPaket">
        <div class="field">
          <label>Nama Paket</label>
          <input v-model="form.nama_paket" required />
        </div>

        <div class="field">
          <label>Harga</label>
          <input v-model="form.harga" type="number" min="0" required />
        </div>

        <div class="field">
          <label>Durasi (hari)</label>
          <input v-model="form.durasi_hari" type="number" min="0" required />
        </div>

        <div class="field">
          <label>Durasi (bulan)</label>
          <input v-model="form.durasi_bulan" type="number" min="1" required />
        </div>

        <div class="field" style="grid-column: 1 / -1">
          <label>Deskripsi</label>
          <textarea v-model="form.deskripsi" rows="2" />
        </div>

        <div>
          <button class="btn btn-primary" type="submit">Simpan Paket</button>
        </div>
      </form>
      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>
    </BaseModal>

    <BaseModal :show="showEditModal" title="Ubah Paket Bimbingan" @close="showEditModal = false">
      <form class="form-grid" @submit.prevent="submitEditPaket">
        <div class="field">
          <label>Nama Paket</label>
          <input v-model="editForm.nama_paket" required />
        </div>

        <div class="field">
          <label>Harga</label>
          <input v-model="editForm.harga" type="number" min="0" required />
        </div>

        <div class="field">
          <label>Durasi Hari</label>
          <input v-model="editForm.durasi_hari" type="number" min="0" required />
        </div>

        <div class="field">
          <label>Durasi Bulan</label>
          <input v-model="editForm.durasi_bulan" type="number" min="1" required />
        </div>

        <div class="field">
          <label>Status</label>
          <select v-model="editForm.status">
            <option value="aktif">aktif</option>
            <option value="tidak aktif">tidak aktif</option>
          </select>
        </div>

        <div class="field" style="grid-column: 1 / -1">
          <label>Deskripsi</label>
          <textarea v-model="editForm.deskripsi" rows="2" />
        </div>

        <div>
          <button class="btn btn-primary btn-edit-action" type="submit">Update Paket</button>
        </div>
      </form>
    </BaseModal>

    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
    <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Paket</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchPaket" placeholder="Cari data paket..." />
        <div class="tools-actions">
          <button class="btn btn-primary" type="button" @click="showCreateModal = true">Tambah Paket</button>
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF data paket"
            @click="downloadPaketPdf"
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
              <th>Nama Paket</th>
              <th>Harga</th>
              <th>Durasi</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedPaketList" :key="item.id_paket || index">
              <td>{{ nomorUrut(index) }}</td>
              <td>{{ item.nama_paket }}</td>
              <td>Rp {{ Number(item.harga || 0).toLocaleString('id-ID') }}</td>
              <td>{{ formatDurasi(item) }}</td>
              <td>
                <span
                  class="pill"
                  :class="item.status === 'aktif' ? 'pill-success' : 'pill-warning'"
                >
                  {{ item.status || '-' }}
                </span>
              </td>
              <td class="actions">
                <button class="btn btn-secondary btn-edit-action" @click="selectPaket(item)">Edit</button>
                <button class="btn btn-danger" @click="handleDeletePaket(item.id_paket)">Hapus</button>
              </td>
            </tr>
            <tr v-if="paginatedPaketList.length === 0">
              <td colspan="6">Tidak ada data paket.</td>
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
          :disabled="currentPagePaket === 1"
        >
          <<
        </button>
        <span>Halaman {{ currentPagePaket }} / {{ totalPagesPaket }} ({{ filteredPaketList.length }} data)</span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="goNextPage"
          :disabled="currentPagePaket === totalPagesPaket"
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
  margin: 0;
  font-size: 1.05rem;
}

.actions {
  display: flex;
  gap: 8px;
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
