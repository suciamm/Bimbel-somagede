<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import {
  approvalPembimbingApi,
  deletePembimbingApi,
  getPembimbingApi,
  getPengajuanPembimbingApi,
  updatePembimbingApi
} from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const pembimbingList = ref([])
const pengajuanPembimbingList = ref([])

const showEditModal = ref(false)
const viewMode = ref('daftar')
const searchPembimbing = ref('')
const searchPengajuanPembimbing = ref('')
const currentPagePembimbing = ref(1)
const currentPagePengajuanPembimbing = ref(1)

const PAGE_SIZE = 10

const editForm = reactive({
  id_user: '',
  nama_lengkap: '',
  no_telp: '',
  alamat: '',
  status: 'true'
})

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

const filteredPembimbingList = computed(() => {
  const source = toArray(pembimbingList.value)
  const keyword = normalizeText(searchPembimbing.value).trim()
  if (!keyword) return source

  return source.filter((item) => {
    const joined = [
      item.id_user,
      item.kode_pembimbing,
      item.username,
      item.nama_lengkap,
      item.no_telp,
      item.status ? 'aktif' : 'nonaktif'
    ]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const filteredPengajuanPembimbingList = computed(() => {
  const source = toArray(pengajuanPembimbingList.value)
  const keyword = normalizeText(searchPengajuanPembimbing.value).trim()
  if (!keyword) return source

  return source.filter((item) => {
    const joined = [item.id_user, item.username, item.nama_lengkap, item.no_telp].map(normalizeText).join(' ')
    return joined.includes(keyword)
  })
})

const totalPagesPembimbing = computed(() => Math.max(1, Math.ceil(filteredPembimbingList.value.length / PAGE_SIZE)))
const totalPagesPengajuanPembimbing = computed(() =>
  Math.max(1, Math.ceil(filteredPengajuanPembimbingList.value.length / PAGE_SIZE))
)

const paginatedPembimbingList = computed(() => {
  const start = (currentPagePembimbing.value - 1) * PAGE_SIZE
  return filteredPembimbingList.value.slice(start, start + PAGE_SIZE)
})

const paginatedPengajuanPembimbingList = computed(() => {
  const start = (currentPagePengajuanPembimbing.value - 1) * PAGE_SIZE
  return filteredPengajuanPembimbingList.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(page, index) {
  return (page - 1) * PAGE_SIZE + index + 1
}

function goPrevDaftarPage() {
  if (currentPagePembimbing.value > 1) currentPagePembimbing.value -= 1
}

function goNextDaftarPage() {
  if (currentPagePembimbing.value < totalPagesPembimbing.value) currentPagePembimbing.value += 1
}

function goPrevPengajuanPage() {
  if (currentPagePengajuanPembimbing.value > 1) currentPagePengajuanPembimbing.value -= 1
}

function goNextPengajuanPage() {
  if (currentPagePengajuanPembimbing.value < totalPagesPengajuanPembimbing.value) {
    currentPagePengajuanPembimbing.value += 1
  }
}

function toggleViewMode() {
  viewMode.value = viewMode.value === 'daftar' ? 'pengajuan' : 'daftar'
}

function downloadPembimbingPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredPembimbingList.value.map((item, index) => [
    index + 1,
    item.id_user || '-',
    item.kode_pembimbing || '-',
    item.username || '-',
    item.nama_lengkap || '-',
    item.no_telp || '-',
    item.status ? 'aktif' : 'nonaktif'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Pembimbing', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'ID Pembimbing', 'Kode Pembimbing', 'Username', 'Nama Lengkap', 'No Telepon', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-pembimbing.pdf')
}

function downloadPengajuanPembimbingPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredPengajuanPembimbingList.value.map((item, index) => [
    index + 1,
    item.id_user || '-',
    item.username || '-',
    item.nama_lengkap || '-',
    item.no_telp || '-'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Pengajuan Pembimbing', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'ID Pembimbing', 'Username', 'Nama Lengkap', 'No Telepon']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-pengajuan-pembimbing.pdf')
}

async function loadPembimbing() {
  loading.value = true
  errorMessage.value = ''

  try {
    const [responseDaftar, responsePengajuan] = await Promise.all([getPembimbingApi(), getPengajuanPembimbingApi()])
    pembimbingList.value = toArray(responseDaftar?.data || responseDaftar)
    pengajuanPembimbingList.value = toArray(responsePengajuan?.data || responsePengajuan)
    currentPagePembimbing.value = 1
    currentPagePengajuanPembimbing.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

function selectPembimbing(item) {
  editForm.id_user = String(item.id_user || '')
  editForm.nama_lengkap = item.nama_lengkap || ''
  editForm.no_telp = item.no_telp || ''
  editForm.alamat = item.alamat || ''
  editForm.status = item.status ? 'true' : 'false'
  showEditModal.value = true
}

async function submitEditPembimbing() {
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await updatePembimbingApi(editForm.id_user, {
      nama_lengkap: editForm.nama_lengkap,
      no_telp: editForm.no_telp,
      alamat: editForm.alamat,
      status: editForm.status === 'true'
    })

    successMessage.value = 'Data pembimbing berhasil diperbarui.'
    showEditModal.value = false
    await loadPembimbing()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function handleDeletePembimbing(idUser) {
  if (!window.confirm('Hapus data pembimbing ini?')) return

  errorMessage.value = ''
  successMessage.value = ''
  try {
    await deletePembimbingApi(idUser)
    successMessage.value = 'Data pembimbing berhasil dihapus.'
    await loadPembimbing()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function approvalPembimbing(idUser, status) {
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await approvalPembimbingApi({ id_user: Number(idUser), status })
    successMessage.value = `Pengajuan pembimbing berhasil di-${status}.`
    await loadPembimbing()
  } catch (error) {
    errorMessage.value = error.message
  }
}

watch(searchPembimbing, () => {
  currentPagePembimbing.value = 1
})

watch(searchPengajuanPembimbing, () => {
  currentPagePengajuanPembimbing.value = 1
})

watch(filteredPembimbingList, () => {
  if (currentPagePembimbing.value > totalPagesPembimbing.value) {
    currentPagePembimbing.value = totalPagesPembimbing.value
  }
})

watch(filteredPengajuanPembimbingList, () => {
  if (currentPagePengajuanPembimbing.value > totalPagesPengajuanPembimbing.value) {
    currentPagePengajuanPembimbing.value = totalPagesPengajuanPembimbing.value
  }
})

onMounted(loadPembimbing)
</script>

<template>
  <section>
    <h1 class="page-title">Data Pembimbing</h1>
    <p class="page-subtitle">Kelola akun pembimbing dan proses approval registrasi pembimbing.</p>

    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
    <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

    <section class="panel block">
      <header class="block-header">
        <h2>{{ viewMode === 'daftar' ? 'Daftar Pembimbing' : 'Daftar Pengajuan Pembimbing' }}</h2>
      </header>

      <div class="table-tools">
        <input
          :value="viewMode === 'daftar' ? searchPembimbing : searchPengajuanPembimbing"
          :placeholder="viewMode === 'daftar' ? 'Cari data pembimbing...' : 'Cari pengajuan pembimbing...'"
          @input="viewMode === 'daftar' ? (searchPembimbing = $event.target.value) : (searchPengajuanPembimbing = $event.target.value)"
        />
        <div class="tools-actions">
          <button class="btn btn-secondary" type="button" @click="toggleViewMode">
            {{ viewMode === 'daftar' ? 'Daftar Pengajuan' : 'Kembali ke Daftar Pembimbing' }}
          </button>
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF data pembimbing"
            @click="viewMode === 'daftar' ? downloadPembimbingPdf() : downloadPengajuanPembimbingPdf()"
          >
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>

      <div class="table-wrap">
        <table v-if="viewMode === 'daftar'">
          <thead>
            <tr>
              <th>No</th>
              <th>ID Pembimbing</th>
              <th>Kode Pembimbing</th>
              <th>Username</th>
              <th>Nama Lengkap</th>
              <th>No Telepon</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedPembimbingList" :key="item.id_user || index">
              <td>{{ nomorUrut(currentPagePembimbing, index) }}</td>
              <td>{{ item.id_user }}</td>
              <td>{{ item.kode_pembimbing || '-' }}</td>
              <td>{{ item.username }}</td>
              <td>{{ item.nama_lengkap }}</td>
              <td>{{ item.no_telp || '-' }}</td>
              <td><span class="pill" :class="item.status ? 'pill-success' : 'pill-warning'">{{ item.status ? 'aktif' : 'nonaktif' }}</span></td>
              <td class="actions">
                <button class="btn btn-secondary btn-edit-action" @click="selectPembimbing(item)">Edit</button>
                <button class="btn btn-danger" @click="handleDeletePembimbing(item.id_user)">Hapus</button>
              </td>
            </tr>
            <tr v-if="paginatedPembimbingList.length === 0">
              <td colspan="8">Tidak ada data pembimbing.</td>
            </tr>
          </tbody>
        </table>

        <table v-else>
          <thead>
            <tr>
              <th>No</th>
              <th>ID Pembimbing</th>
              <th>Username</th>
              <th>Nama Lengkap</th>
              <th>No Telepon</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedPengajuanPembimbingList" :key="item.id_user || index">
              <td>{{ nomorUrut(currentPagePengajuanPembimbing, index) }}</td>
              <td>{{ item.id_user || '-' }}</td>
              <td>{{ item.username || '-' }}</td>
              <td>{{ item.nama_lengkap || '-' }}</td>
              <td>{{ item.no_telp || '-' }}</td>
              <td class="actions">
                <button class="btn btn-primary" @click="approvalPembimbing(item.id_user, 'setujui')">Setujui</button>
                <button class="btn btn-danger" @click="approvalPembimbing(item.id_user, 'tolak')">Tolak</button>
              </td>
            </tr>
            <tr v-if="paginatedPengajuanPembimbingList.length === 0">
              <td colspan="6">Tidak ada data pengajuan pembimbing.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination-wrap" v-if="!loading">
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          @click="viewMode === 'daftar' ? goPrevDaftarPage() : goPrevPengajuanPage()"
          :disabled="viewMode === 'daftar' ? currentPagePembimbing === 1 : currentPagePengajuanPembimbing === 1"
        >
          <<
        </button>
        <span v-if="viewMode === 'daftar'">
          Halaman {{ currentPagePembimbing }} / {{ totalPagesPembimbing }} ({{ filteredPembimbingList.length }} data)
        </span>
        <span v-else>
          Halaman {{ currentPagePengajuanPembimbing }} / {{ totalPagesPengajuanPembimbing }}
          ({{ filteredPengajuanPembimbingList.length }} data)
        </span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          @click="viewMode === 'daftar' ? goNextDaftarPage() : goNextPengajuanPage()"
          :disabled="viewMode === 'daftar' ? currentPagePembimbing === totalPagesPembimbing : currentPagePengajuanPembimbing === totalPagesPengajuanPembimbing"
        >
          >>
        </button>
      </div>
    </section>

    <BaseModal :show="showEditModal" title="Ubah Data Pembimbing" @close="showEditModal = false">
      <form class="form-grid" @submit.prevent="submitEditPembimbing">
        <div class="field">
          <label>ID Pembimbing</label>
          <input v-model="editForm.id_user" type="number" min="1" required />
        </div>
        <div class="field">
          <label>Nama Lengkap</label>
          <input v-model="editForm.nama_lengkap" required />
        </div>
        <div class="field">
          <label>No Telepon</label>
          <input v-model="editForm.no_telp" />
        </div>
        <div class="field">
          <label>Status</label>
          <select v-model="editForm.status">
            <option value="true">aktif</option>
            <option value="false">nonaktif</option>
          </select>
        </div>
        <div class="field" style="grid-column: 1 / -1">
          <label>Alamat</label>
          <textarea v-model="editForm.alamat" rows="2" />
        </div>
        <div>
          <button class="btn btn-primary btn-edit-action" type="submit">Simpan Perubahan</button>
        </div>
      </form>
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

.actions {
  display: flex;
  gap: 8px;
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
