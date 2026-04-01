<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import {
  deleteMuridApi,
  getDaftarOrtuApi,
  getMuridAktifApi,
  getMuridTidakAktifApi,
  tambahMuridApi,
  updateMuridApi
} from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const muridAktif = ref([])
const muridTidakAktif = ref([])
const daftarOrtu = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingMuridId = ref('')
const searchAktif = ref('')
const searchTidakAktif = ref('')
const currentPageAktif = ref(1)
const currentPageTidakAktif = ref(1)
const PAGE_SIZE = 10

const form = reactive({
  nama_murid: '',
  tgl_lahir: '',
  id_user: '',
  alamat: '',
  tgl_masuk: ''
})

const editForm = reactive({
  nama_murid: '',
  tgl_lahir: '',
  alamat: '',
  tgl_masuk: '',
  tgl_keluar: '',
  status_murid: 'aktif'
})

function toDateOnly(value) {
  const iso = String(value || '').match(/\d{4}-\d{2}-\d{2}/)?.[0]
  if (!iso) return '-'
  const [year, month, day] = iso.split('-')
  return `${day}-${month}-${year}`
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

const filteredMuridAktif = computed(() => {
  const keyword = normalizeText(searchAktif.value).trim()
  if (!keyword) return muridAktif.value

  return muridAktif.value.filter((item) => {
    const joined = [
      item.kode_murid,
      item.nama_murid,
      item.nama_ortu,
      item.status_murid,
      toDateOnly(item.tgl_lahir),
      toDateOnly(item.tgl_masuk)
    ]
      .map(normalizeText)
      .join(' ')

    return joined.includes(keyword)
  })
})

const filteredMuridTidakAktif = computed(() => {
  const keyword = normalizeText(searchTidakAktif.value).trim()
  if (!keyword) return muridTidakAktif.value

  return muridTidakAktif.value.filter((item) => {
    const joined = [item.nama_murid, item.nama_ortu, item.status_murid, toDateOnly(item.tgl_keluar)]
      .map(normalizeText)
      .join(' ')

    return joined.includes(keyword)
  })
})

const totalPagesAktif = computed(() => Math.max(1, Math.ceil(filteredMuridAktif.value.length / PAGE_SIZE)))
const totalPagesTidakAktif = computed(() => Math.max(1, Math.ceil(filteredMuridTidakAktif.value.length / PAGE_SIZE)))

const paginatedMuridAktif = computed(() => {
  const start = (currentPageAktif.value - 1) * PAGE_SIZE
  return filteredMuridAktif.value.slice(start, start + PAGE_SIZE)
})

const paginatedMuridTidakAktif = computed(() => {
  const start = (currentPageTidakAktif.value - 1) * PAGE_SIZE
  return filteredMuridTidakAktif.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(page, index) {
  return (page - 1) * PAGE_SIZE + index + 1
}

function goPrevAktifPage() {
  if (currentPageAktif.value > 1) currentPageAktif.value -= 1
}

function goNextAktifPage() {
  if (currentPageAktif.value < totalPagesAktif.value) currentPageAktif.value += 1
}

function goPrevTidakAktifPage() {
  if (currentPageTidakAktif.value > 1) currentPageTidakAktif.value -= 1
}

function goNextTidakAktifPage() {
  if (currentPageTidakAktif.value < totalPagesTidakAktif.value) currentPageTidakAktif.value += 1
}

function downloadMuridAktifPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredMuridAktif.value.map((item, index) => [
    index + 1,
    item.kode_murid || '-',
    item.nama_murid || '-',
    item.nama_ortu || '-',
    toDateOnly(item.tgl_lahir),
    toDateOnly(item.tgl_masuk),
    item.status_murid || '-'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Murid Aktif', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Kode', 'Nama Murid', 'Nama Ortu', 'Tgl Lahir', 'Tgl Masuk', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-murid-aktif.pdf')
}

function downloadMuridTidakAktifPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredMuridTidakAktif.value.map((item, index) => [
    index + 1,
    item.nama_murid || '-',
    item.nama_ortu || '-',
    toDateOnly(item.tgl_keluar),
    item.status_murid || '-'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Murid Tidak Aktif', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Nama Murid', 'Nama Ortu', 'Tgl Keluar', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-murid-tidak-aktif.pdf')
}

async function loadMurid() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [aktif, nonAktif, ortu] = await Promise.all([
      getMuridAktifApi(),
      getMuridTidakAktifApi(),
      getDaftarOrtuApi()
    ])

    const aktifRaw = aktif?.data || aktif || []
    const nonAktifRaw = nonAktif?.data || nonAktif || []
    const ortuRaw = ortu?.data || ortu || []

    muridAktif.value = Array.isArray(aktifRaw)
      ? aktifRaw.filter((item) => item && typeof item === 'object')
      : []

    muridTidakAktif.value = Array.isArray(nonAktifRaw)
      ? nonAktifRaw.filter((item) => item && typeof item === 'object')
      : []

    daftarOrtu.value = Array.isArray(ortuRaw)
      ? ortuRaw.filter((item) => item && typeof item === 'object' && item.id_user)
      : []

    currentPageAktif.value = 1
    currentPageTidakAktif.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(searchAktif, () => {
  currentPageAktif.value = 1
})

watch(searchTidakAktif, () => {
  currentPageTidakAktif.value = 1
})

watch(filteredMuridAktif, () => {
  if (currentPageAktif.value > totalPagesAktif.value) {
    currentPageAktif.value = totalPagesAktif.value
  }
})

watch(filteredMuridTidakAktif, () => {
  if (currentPageTidakAktif.value > totalPagesTidakAktif.value) {
    currentPageTidakAktif.value = totalPagesTidakAktif.value
  }
})

async function submitMurid() {
  successMessage.value = ''
  errorMessage.value = ''

  try {
    await tambahMuridApi({
      ...form,
      id_user: Number(form.id_user)
    })

    successMessage.value = 'Data murid berhasil ditambahkan.'
    form.nama_murid = ''
    form.tgl_lahir = ''
    form.id_user = ''
    form.alamat = ''
    form.tgl_masuk = ''
    showCreateModal.value = false

    await loadMurid()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function handleDeleteMurid(idMurid) {
  if (!window.confirm('Hapus data murid ini?')) return

  successMessage.value = ''
  errorMessage.value = ''

  try {
    await deleteMuridApi(idMurid)
    successMessage.value = 'Data murid berhasil dihapus.'
    await loadMurid()
  } catch (error) {
    errorMessage.value = error.message
  }
}

function openEditMurid(item) {
  editingMuridId.value = item.id_murid || ''
  editForm.nama_murid = item.nama_murid || ''
  editForm.tgl_lahir = (item.tgl_lahir || '').slice?.(0, 10) || ''
  editForm.alamat = item.alamat || ''
  editForm.tgl_masuk = (item.tgl_masuk || '').slice?.(0, 10) || ''
  editForm.tgl_keluar = (item.tgl_keluar || '').slice?.(0, 10) || ''
  editForm.status_murid = item.status_murid || 'aktif'
  showEditModal.value = true
}

async function submitEditMurid() {
  successMessage.value = ''
  errorMessage.value = ''

  if (!editingMuridId.value) {
    errorMessage.value = 'Pilih data murid dari tombol Edit pada tabel.'
    return
  }

  try {
    await updateMuridApi(editingMuridId.value, {
      nama_murid: editForm.nama_murid,
      tgl_lahir: editForm.tgl_lahir,
      alamat: editForm.alamat,
      tgl_masuk: editForm.tgl_masuk,
      tgl_keluar: editForm.tgl_keluar || null,
      status_murid: editForm.status_murid
    })

    successMessage.value = 'Data murid berhasil diperbarui.'
    showEditModal.value = false
    editingMuridId.value = ''
    await loadMurid()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(loadMurid)
</script>

<template>
  <section>
    <h1 class="page-title">Manajemen Murid</h1>
    <p class="page-subtitle">Tambah murid baru dan pantau status keaktifannya.</p>

    <BaseModal :show="showCreateModal" title="Tambah Data Murid" @close="showCreateModal = false">
      <form class="form-grid" @submit.prevent="submitMurid">
        <div class="field">
          <label>Nama Murid</label>
          <input v-model="form.nama_murid" required />
        </div>
        <div class="field">
          <label>Tanggal Lahir</label>
          <input v-model="form.tgl_lahir" type="date" required />
        </div>
        <div class="field">
          <label>Orang Tua</label>
          <select v-model="form.id_user" required>
            <option disabled value="">Pilih orang tua</option>
            <option v-for="ortu in daftarOrtu" :key="ortu.id_user" :value="ortu.id_user">
              ({{ ortu.id_user }}){{ ortu.nama_lengkap }}
            </option>
          </select>
        </div>
        <div class="field">
          <label>Tanggal Masuk</label>
          <input v-model="form.tgl_masuk" type="date" required />
        </div>
        <div class="field" style="grid-column: 1 / -1">
          <label>Alamat</label>
          <textarea v-model="form.alamat" rows="2" required />
        </div>
        <div>
          <button class="btn btn-primary btn-add" type="submit">Simpan Murid</button>
        </div>
      </form>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>
    </BaseModal>

    <BaseModal :show="showEditModal" title="Ubah Data Murid" @close="showEditModal = false">
      <form class="form-grid" @submit.prevent="submitEditMurid">
        <div class="field">
          <label>Nama Murid</label>
          <input v-model="editForm.nama_murid" required />
        </div>
        <div class="field">
          <label>Tanggal Lahir</label>
          <input v-model="editForm.tgl_lahir" type="date" required />
        </div>
        <div class="field">
          <label>Tanggal Masuk</label>
          <input v-model="editForm.tgl_masuk" type="date" required />
        </div>
        <div class="field">
          <label>Tanggal Keluar</label>
          <input v-model="editForm.tgl_keluar" type="date" />
        </div>
        <div class="field">
          <label>Status Murid</label>
          <select v-model="editForm.status_murid">
            <option value="aktif">aktif</option>
            <option value="lulus">lulus</option>
            <option value="keluar">keluar</option>
          </select>
        </div>
        <div class="field" style="grid-column: 1 / -1">
          <label>Alamat</label>
          <textarea v-model="editForm.alamat" rows="2" required />
        </div>
        <div>
          <button class="btn btn-primary btn-edit-action" type="submit">Simpan Perubahan</button>
        </div>
      </form>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>
    </BaseModal>

    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
    <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

    <section class="panel block">
      <h2>Daftar Murid Aktif</h2>
      <div class="table-tools">
        <input v-model="searchAktif" placeholder="Cari data murid aktif..." />
        <div class="tools-actions">
          <button class="btn btn-primary btn-add" type="button" @click="showCreateModal = true">Tambah Murid</button>
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF murid aktif"
            @click="downloadMuridAktifPdf"
          >
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>
      <p v-if="loading">Memuat data murid...</p>
      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Kode</th>
              <th>Nama Murid</th>
              <th>Nama Ortu</th>
              <th>Tgl Lahir</th>
              <th>Tgl Masuk</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedMuridAktif" :key="item.id_murid || index">
              <td>{{ nomorUrut(currentPageAktif, index) }}</td>
              <td>{{ item.kode_murid || '-' }}</td>
              <td>{{ item.nama_murid }}</td>
              <td>{{ item.nama_ortu || '-' }}</td>
              <td>{{ item.tgl_lahir?.slice?.(0, 10) || '-' }}</td>
              <td>{{ item.tgl_masuk?.slice?.(0, 10) || '-' }}</td>
              <td><span class="pill pill-success">{{ item.status_murid }}</span></td>
              <td class="actions">
                <button class="btn btn-secondary btn-edit-action" @click="openEditMurid(item)">Edit</button>
                <button class="btn btn-danger" @click="handleDeleteMurid(item.id_murid)">Hapus</button>
              </td>
            </tr>
            <tr v-if="paginatedMuridAktif.length === 0">
              <td colspan="8">Belum ada data murid aktif.</td>
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
          @click="goPrevAktifPage"
          :disabled="currentPageAktif === 1"
        >
          <<
        </button>
        <span>Halaman {{ currentPageAktif }} / {{ totalPagesAktif }} ({{ filteredMuridAktif.length }} data)</span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="goNextAktifPage"
          :disabled="currentPageAktif === totalPagesAktif"
        >
          >>
        </button>
      </div>
    </section>

    <section class="panel block">
      <h2>Daftar Murid Tidak Aktif</h2>
      <div class="table-tools">
        <input v-model="searchTidakAktif" placeholder="Cari data murid tidak aktif..." />
        <button
          class="btn btn-secondary btn-icon btn-pdf"
          type="button"
          title="Download PDF"
          aria-label="Download PDF murid tidak aktif"
          @click="downloadMuridTidakAktifPdf"
        >
          <span class="pdf-icon" aria-hidden="true">&#128424;</span>
          <span>PDF</span>
        </button>
      </div>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Nama Murid</th>
              <th>Nama Ortu</th>
              <th>Tgl Keluar</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedMuridTidakAktif" :key="item.id_murid || index">
              <td>{{ nomorUrut(currentPageTidakAktif, index) }}</td>
              <td>{{ item.nama_murid }}</td>
              <td>{{ item.nama_ortu || '-' }}</td>
              <td>{{ item.tgl_keluar?.slice?.(0, 10) || '-' }}</td>
              <td><span class="pill pill-warning">{{ item.status_murid }}</span></td>
            </tr>
            <tr v-if="paginatedMuridTidakAktif.length === 0">
              <td colspan="5">Belum ada data murid tidak aktif.</td>
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
          @click="goPrevTidakAktifPage"
          :disabled="currentPageTidakAktif === 1"
        >
          <<
        </button>
        <span>
          Halaman {{ currentPageTidakAktif }} / {{ totalPagesTidakAktif }} ({{ filteredMuridTidakAktif.length }} data)
        </span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="goNextTidakAktifPage"
          :disabled="currentPageTidakAktif === totalPagesTidakAktif"
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

h2 {
  margin-top: 0;
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

.btn-add {
  background-color: #1f9d55;
  border-color: #1f9d55;
  color: #fff;
}

.btn-add:hover,
.btn-add:focus,
.btn-add:active {
  background-color: #168246;
  border-color: #168246;
  color: #fff;
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

@media (max-width: 700px) {
  .table-tools {
    flex-direction: column;
    align-items: stretch;
  }

  .table-tools input {
    max-width: none;
  }

  .tools-actions {
    justify-content: flex-end;
  }
}
</style>
