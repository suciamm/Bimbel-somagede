<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import { approvalOrtuApi, deleteOrtuApi, getDaftarOrtuApi, getPengajuanOrtuApi, updateOrtuApi } from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const ortuList = ref([])
const pengajuanOrtuList = ref([])
const showEditModal = ref(false)
const viewMode = ref('daftar')
const searchOrtu = ref('')
const searchPengajuanOrtu = ref('')
const currentPageOrtu = ref(1)
const currentPagePengajuanOrtu = ref(1)
const PAGE_SIZE = 10

const editForm = reactive({
  id_user: '',
  nama_lengkap: '',
  no_telp: '',
  alamat: '',
  status: 'true'
})

function mapStatusToString(value) {
  if (value === true || value === '1' || value === 'Aktif' || value === 'aktif') return 'true'
  return 'false'
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

const filteredOrtuList = computed(() => {
  const source = toArray(ortuList.value)
  const keyword = normalizeText(searchOrtu.value).trim()
  if (!keyword) return source

  return source.filter((item) => {
    const joined = [
      item.id_user,
      item.nama_lengkap,
      item.no_telp,
      mapStatusToString(item.status) === 'true' ? 'aktif' : 'nonaktif'
    ]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const totalPagesOrtu = computed(() => Math.max(1, Math.ceil(filteredOrtuList.value.length / PAGE_SIZE)))

const paginatedOrtuList = computed(() => {
  const start = (currentPageOrtu.value - 1) * PAGE_SIZE
  return filteredOrtuList.value.slice(start, start + PAGE_SIZE)
})

const filteredPengajuanOrtuList = computed(() => {
  const source = toArray(pengajuanOrtuList.value)
  const keyword = normalizeText(searchPengajuanOrtu.value).trim()
  if (!keyword) return source

  return source.filter((item) => {
    const joined = [item.id_user, item.username, item.nama_lengkap, item.no_telp].map(normalizeText).join(' ')
    return joined.includes(keyword)
  })
})

const totalPagesPengajuanOrtu = computed(() =>
  Math.max(1, Math.ceil(filteredPengajuanOrtuList.value.length / PAGE_SIZE))
)

const paginatedPengajuanOrtuList = computed(() => {
  const start = (currentPagePengajuanOrtu.value - 1) * PAGE_SIZE
  return filteredPengajuanOrtuList.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(page, index) {
  return (page - 1) * PAGE_SIZE + index + 1
}

function goPrevPage() {
  if (currentPageOrtu.value > 1) currentPageOrtu.value -= 1
}

function goNextPage() {
  if (currentPageOrtu.value < totalPagesOrtu.value) currentPageOrtu.value += 1
}

function goPrevPengajuanPage() {
  if (currentPagePengajuanOrtu.value > 1) currentPagePengajuanOrtu.value -= 1
}

function goNextPengajuanPage() {
  if (currentPagePengajuanOrtu.value < totalPagesPengajuanOrtu.value) currentPagePengajuanOrtu.value += 1
}

function downloadOrtuPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredOrtuList.value.map((item, index) => [
    index + 1,
    item.id_user || '-',
    item.nama_lengkap || '-',
    item.no_telp || '-',
    mapStatusToString(item.status) === 'true' ? 'aktif' : 'nonaktif'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Orang Tua', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'ID Orang Tua', 'Nama Lengkap', 'No Telepon', 'Status']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-orang-tua.pdf')
}

function downloadPengajuanOrtuPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredPengajuanOrtuList.value.map((item, index) => [
    index + 1,
    item.id_user || '-',
    item.username || '-',
    item.nama_lengkap || '-',
    item.no_telp || '-'
  ])

  doc.setFontSize(12)
  doc.text('Daftar Pengajuan Orang Tua', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'ID Orang Tua', 'Username', 'Nama Lengkap', 'No Telepon']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-pengajuan-orang-tua.pdf')
}

function toggleViewMode() {
  viewMode.value = viewMode.value === 'daftar' ? 'pengajuan' : 'daftar'
}

async function loadOrtu() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [responseDaftar, responsePengajuan] = await Promise.all([getDaftarOrtuApi(), getPengajuanOrtuApi()])
    ortuList.value = toArray(responseDaftar?.data || responseDaftar)
    pengajuanOrtuList.value = toArray(responsePengajuan?.data || responsePengajuan)
    currentPageOrtu.value = 1
    currentPagePengajuanOrtu.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(searchOrtu, () => {
  currentPageOrtu.value = 1
})

watch(searchPengajuanOrtu, () => {
  currentPagePengajuanOrtu.value = 1
})

watch(filteredOrtuList, () => {
  if (currentPageOrtu.value > totalPagesOrtu.value) {
    currentPageOrtu.value = totalPagesOrtu.value
  }
})

watch(filteredPengajuanOrtuList, () => {
  if (currentPagePengajuanOrtu.value > totalPagesPengajuanOrtu.value) {
    currentPagePengajuanOrtu.value = totalPagesPengajuanOrtu.value
  }
})

function selectForEdit(item) {
  editForm.id_user = String(item.id_user || '')
  editForm.nama_lengkap = item.nama_lengkap || ''
  editForm.no_telp = item.no_telp || ''
  editForm.alamat = item.alamat || ''
  editForm.status = mapStatusToString(item.status)
  showEditModal.value = true
}

async function submitUpdateOrtu() {
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await updateOrtuApi(editForm.id_user, {
      nama_lengkap: editForm.nama_lengkap,
      no_telp: editForm.no_telp,
      alamat: editForm.alamat,
      status: editForm.status === 'true'
    })

    successMessage.value = 'Data orang tua berhasil diperbarui.'
    showEditModal.value = false
    await loadOrtu()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function handleDelete(idUser) {
  if (!window.confirm('Hapus data orang tua ini?')) return

  errorMessage.value = ''
  successMessage.value = ''
  try {
    await deleteOrtuApi(idUser)
    successMessage.value = 'Data orang tua berhasil dihapus.'
    await loadOrtu()
  } catch (error) {
    errorMessage.value = `${error.message}`
  }
}

async function approvalOrtu(idUser, status) {
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await approvalOrtuApi({ id_user: Number(idUser), status })
    successMessage.value = `Pengajuan orang tua berhasil di-${status}.`
    await loadOrtu()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(loadOrtu)
</script>

<template>
  <section>
    <h1 class="page-title">Data Orang Tua</h1>
    <p class="page-subtitle">Kelola akun orang tua: lihat, edit profil, dan hapus data.</p>

    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
    <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

    <section class="panel block">
      <header class="block-header">
        <h2>{{ viewMode === 'daftar' ? 'Daftar Orang Tua' : 'Daftar Pengajuan Orang Tua Murid' }}</h2>
      </header>

      <div class="table-tools">
        <input
          :value="viewMode === 'daftar' ? searchOrtu : searchPengajuanOrtu"
          :placeholder="viewMode === 'daftar' ? 'Cari data orang tua...' : 'Cari pengajuan orang tua...'"
          @input="viewMode === 'daftar' ? (searchOrtu = $event.target.value) : (searchPengajuanOrtu = $event.target.value)"
        />
        <div class="tools-actions">
          <button class="btn btn-secondary" type="button" @click="toggleViewMode">
            {{ viewMode === 'daftar' ? 'Daftar Pengajuan' : 'Kembali ke Daftar Orang Tua' }}
          </button>
        <button
          class="btn btn-secondary btn-icon btn-pdf"
          type="button"
          title="Download PDF"
          aria-label="Download PDF data orang tua"
          @click="viewMode === 'daftar' ? downloadOrtuPdf() : downloadPengajuanOrtuPdf()"
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
              <th>ID Orang Tua</th>
              <th>Nama Lengkap</th>
              <th>No Telepon</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedOrtuList" :key="item.id_user || index">
              <td>{{ nomorUrut(currentPageOrtu, index) }}</td>
              <td>{{ item.id_user }}</td>
              <td>{{ item.nama_lengkap }}</td>
              <td>{{ item.no_telp || '-' }}</td>
              <td>
                <span class="pill" :class="mapStatusToString(item.status) === 'true' ? 'pill-success' : 'pill-warning'">
                  {{ mapStatusToString(item.status) === 'true' ? 'aktif' : 'nonaktif' }}
                </span>
              </td>
              <td class="actions">
                <button class="btn btn-secondary btn-edit-action" @click="selectForEdit(item)">Edit</button>
                <button class="btn btn-danger" @click="handleDelete(item.id_user)">Hapus</button>
              </td>
            </tr>
            <tr v-if="paginatedOrtuList.length === 0">
              <td colspan="6">Tidak ada data orang tua.</td>
            </tr>
          </tbody>
        </table>
        <table v-else>
          <thead>
            <tr>
              <th>No</th>
              <th>ID Orang Tua</th>
              <th>Username</th>
              <th>Nama Lengkap</th>
              <th>No Telepon</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedPengajuanOrtuList" :key="item.id_user || index">
              <td>{{ nomorUrut(currentPagePengajuanOrtu, index) }}</td>
              <td>{{ item.id_user || '-' }}</td>
              <td>{{ item.username || '-' }}</td>
              <td>{{ item.nama_lengkap || '-' }}</td>
              <td>{{ item.no_telp || '-' }}</td>
              <td class="actions">
                <button class="btn btn-primary" @click="approvalOrtu(item.id_user, 'setujui')">Setujui</button>
                <button class="btn btn-danger" @click="approvalOrtu(item.id_user, 'tolak')">Tolak</button>
              </td>
            </tr>
            <tr v-if="paginatedPengajuanOrtuList.length === 0">
              <td colspan="6">Tidak ada data pengajuan orang tua.</td>
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
          @click="viewMode === 'daftar' ? goPrevPage() : goPrevPengajuanPage()"
          :disabled="viewMode === 'daftar' ? currentPageOrtu === 1 : currentPagePengajuanOrtu === 1"
        >
          <<
        </button>
        <span v-if="viewMode === 'daftar'">
          Halaman {{ currentPageOrtu }} / {{ totalPagesOrtu }} ({{ filteredOrtuList.length }} data)
        </span>
        <span v-else>
          Halaman {{ currentPagePengajuanOrtu }} / {{ totalPagesPengajuanOrtu }} ({{ filteredPengajuanOrtuList.length }} data)
        </span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="viewMode === 'daftar' ? goNextPage() : goNextPengajuanPage()"
          :disabled="viewMode === 'daftar' ? currentPageOrtu === totalPagesOrtu : currentPagePengajuanOrtu === totalPagesPengajuanOrtu"
        >
          >>
        </button>
      </div>
    </section>

    <BaseModal :show="showEditModal" title="Ubah Data Orang Tua" @close="showEditModal = false">
      <form class="form-grid" @submit.prevent="submitUpdateOrtu">
        <div class="field">
          <label>ID Orang Tua</label>
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
