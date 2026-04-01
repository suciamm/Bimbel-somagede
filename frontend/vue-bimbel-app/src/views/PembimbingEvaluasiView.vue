<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import BaseModal from '../components/BaseModal.vue'
import { createEvaluasiApi, getEvaluasiByPembimbingApi } from '../services/api'
import { useAuthStore } from '../stores/auth'

// Ubah nilai ini jika ingin ganti hari input evaluasi.
// getDay(): Minggu=0, Senin=1, ..., Sabtu=6
const ALLOWED_EVALUASI_DAY_INDEX = 0
const ALLOWED_EVALUASI_DAY_LABEL = 'Minggu'

const authStore = useAuthStore()
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const showCreateModal = ref(false)
const evaluasiList = ref([])

const form = reactive({
  id_murid: '',
  evaluasi_ke: 1,
  nilai: 'A',
  catatan_pembimbing: '',
  tanggal_evaluasi: new Date().toISOString().slice(0, 10)
})

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

const searchKeyword = ref('')

const filteredList = computed(() => {
  const keyword = normalizeText(searchKeyword.value).trim()
  if (!keyword) return evaluasiList.value

  return evaluasiList.value.filter((item) => {
    const joined = [
      item.kode_murid,
      item.nama_murid,
      item.status_progress_text,
      item.evaluasi_1_nilai,
      item.evaluasi_2_nilai,
      item.evaluasi_3_nilai
    ]
      .map(normalizeText)
      .join(' ')

    return joined.includes(keyword)
  })
})

const isTodayEvaluasiDay = computed(() => new Date().getDay() === ALLOWED_EVALUASI_DAY_INDEX)

const selectedMurid = computed(() => {
  const id = Number(form.id_murid)
  if (!id) return null
  return evaluasiList.value.find((item) => Number(item.id_murid) === id) || null
})

const canOpenInput = computed(() => isTodayEvaluasiDay.value)

const submitDisabled = computed(() => {
  if (!canOpenInput.value) return true
  if (!form.id_murid || !form.tanggal_evaluasi || !form.nilai) return true
  const tahap = Number(form.evaluasi_ke)
  return ![1, 2, 3].includes(tahap)
})

const muridOptions = computed(() => {
  return evaluasiList.value.map((item) => ({
    id_murid: item.id_murid,
    kode_murid: item.kode_murid,
    nama_murid: item.nama_murid,
    tahap_berikutnya: item.tahap_berikutnya || 1
  }))
})

function openCreateModal() {
  errorMessage.value = ''
  successMessage.value = ''
  if (!canOpenInput.value) {
    errorMessage.value = `Input evaluasi hanya bisa dilakukan pada hari ${ALLOWED_EVALUASI_DAY_LABEL}.`
    return
  }
  showCreateModal.value = true
}

function onMuridChange() {
  if (!selectedMurid.value) {
    form.evaluasi_ke = 1
    return
  }
  form.evaluasi_ke = Number(selectedMurid.value.tahap_berikutnya || 1)
}

async function loadEvaluasi() {
  const idUser = Number(authStore.user?.id_user)
  if (!idUser) {
    evaluasiList.value = []
    errorMessage.value = 'ID pembimbing tidak ditemukan pada sesi login.'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const response = await getEvaluasiByPembimbingApi(idUser)
    console.log('[DEBUG] API getEvaluasiByPembimbingApi response:', response)
    evaluasiList.value = toArray(response)
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

async function submitEvaluasi() {
  if (submitDisabled.value) return

  errorMessage.value = ''
  successMessage.value = ''

  try {
    await createEvaluasiApi({
      id_murid: Number(form.id_murid),
      id_pembimbing: Number(authStore.user?.id_user),
      evaluasi_ke: Number(form.evaluasi_ke),
      nilai: String(form.nilai || '').toUpperCase(),
      catatan_pembimbing: form.catatan_pembimbing,
      tanggal_evaluasi: form.tanggal_evaluasi
    })

    successMessage.value = 'Evaluasi berhasil disimpan.'
    showCreateModal.value = false
    form.id_murid = ''
    form.evaluasi_ke = 1
    form.nilai = 'A'
    form.catatan_pembimbing = ''
    form.tanggal_evaluasi = new Date().toISOString().slice(0, 10)

    await loadEvaluasi()
  } catch (error) {
    errorMessage.value = error.message
  }
}

function formatTahap(item, tahap) {
  const nilai = item?.[`evaluasi_${tahap}_nilai`] || '-'
  const tanggal = item?.[`evaluasi_${tahap}_tanggal`] || '-'
  const catatan = item?.[`evaluasi_${tahap}_catatan`] || '-'
  return { nilai, tanggal, catatan }
}

onMounted(loadEvaluasi)
</script>

<template>
  <section>
    <h1 class="page-title">Evaluasi Murid</h1>
    <p class="page-subtitle">Input evaluasi tahap 1-3 untuk murid bimbingan. Tahap selanjutnya terbuka jika tahap sebelumnya bernilai A.</p>

    <BaseModal :show="showCreateModal" title="Input Evaluasi Murid" @close="showCreateModal = false">
      <form class="form-grid" @submit.prevent="submitEvaluasi">
        <div class="field">
          <label>Murid</label>
          <select v-model="form.id_murid" required @change="onMuridChange">
            <option disabled value="">Pilih murid</option>
            <option v-for="item in muridOptions" :key="item.id_murid" :value="item.id_murid">
              {{ item.kode_murid }} - {{ item.nama_murid }} (Tahap {{ item.tahap_berikutnya }})
            </option>
          </select>
        </div>

        <div class="field">
          <label>Evaluasi Ke</label>
          <input v-model.number="form.evaluasi_ke" type="number" min="1" max="3" readonly />
        </div>

        <div class="field">
          <label>Nilai</label>
          <select v-model="form.nilai" required>
            <option value="A">A</option>
            <option value="B">B</option>
            <option value="C">C</option>
            <option value="D">D</option>
            <option value="E">E</option>
          </select>
        </div>

        <div class="field">
          <label>Tanggal Evaluasi</label>
          <input v-model="form.tanggal_evaluasi" type="date" required />
        </div>

        <div class="field" style="grid-column: 1 / -1">
          <label>Catatan Pembimbing</label>
          <textarea v-model="form.catatan_pembimbing" rows="3" placeholder="Contoh: bisa membaca kata sederhana" />
        </div>

        <div>
          <button class="btn btn-primary" type="submit" :disabled="submitDisabled">Simpan Evaluasi</button>
        </div>
      </form>

      <p class="modal-hint" v-if="selectedMurid">Murid dipilih akan otomatis diarahkan ke tahap berikutnya.</p>
      <p class="message message-error" v-if="errorMessage">{{ errorMessage }}</p>
    </BaseModal>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Evaluasi Murid Bimbingan</h2>
        <button class="btn btn-primary" type="button" :disabled="!canOpenInput" @click="openCreateModal">Input Evaluasi</button>
      </header>

      <p v-if="!isTodayEvaluasiDay" class="message message-warning">
        Hari ini bukan {{ ALLOWED_EVALUASI_DAY_LABEL }}. Input evaluasi dikunci di frontend.
      </p>
      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

      <div class="table-tools">
        <input v-model="searchKeyword" placeholder="Cari murid / progres evaluasi..." />
      </div>

      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Kode</th>
              <th>Nama Murid</th>
              <th>Evaluasi 1</th>
              <th>Evaluasi 2</th>
              <th>Evaluasi 3</th>
              <th>Progress</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in filteredList" :key="item.id_murid || index">
              <td>{{ index + 1 }}</td>
              <td>{{ item.kode_murid || '-' }}</td>
              <td>{{ item.nama_murid || '-' }}</td>
              <td>
                <div class="eval-cell">
                  <strong>{{ formatTahap(item, 1).nilai }}</strong>
                  <small>{{ formatTahap(item, 1).tanggal }}</small>
                  <small>{{ formatTahap(item, 1).catatan }}</small>
                </div>
              </td>
              <td>
                <div class="eval-cell">
                  <strong>{{ formatTahap(item, 2).nilai }}</strong>
                  <small>{{ formatTahap(item, 2).tanggal }}</small>
                  <small>{{ formatTahap(item, 2).catatan }}</small>
                </div>
              </td>
              <td>
                <div class="eval-cell">
                  <strong>{{ formatTahap(item, 3).nilai }}</strong>
                  <small>{{ formatTahap(item, 3).tanggal }}</small>
                  <small>{{ formatTahap(item, 3).catatan }}</small>
                </div>
              </td>
              <td>{{ item.status_progress_text || '-' }}</td>
            </tr>
            <tr v-if="!loading && filteredList.length === 0">
              <td colspan="7">Belum ada data evaluasi murid.</td>
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
  margin: 10px 0;
}

.table-tools input {
  width: 100%;
  max-width: 420px;
}

.eval-cell {
  display: grid;
  gap: 2px;
}

.eval-cell small {
  color: var(--text-muted);
  font-size: 0.78rem;
}

.modal-hint {
  margin-top: 8px;
  color: var(--text-muted);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
