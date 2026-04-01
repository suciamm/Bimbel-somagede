<script setup>
import { computed, onMounted, ref } from 'vue'
import { getEvaluasiByOrtuApi } from '../services/api'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const errorMessage = ref('')
const evaluasiList = ref([])
const searchKeyword = ref('')

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

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
      item.evaluasi_3_nilai,
      item.nama_pembimbing
    ]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

function formatTahap(item, tahap) {
  const nilai = item?.[`evaluasi_${tahap}_nilai`] || '-'
  const tanggal = item?.[`evaluasi_${tahap}_tanggal`] || '-'
  const catatan = item?.[`evaluasi_${tahap}_catatan`] || '-'
  return { nilai, tanggal, catatan }
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const idUser = Number(authStore.user?.id_user)
    if (!idUser) {
      throw new Error('ID orangtua tidak ditemukan')
    }

    const response = await getEvaluasiByOrtuApi(idUser)
    evaluasiList.value = toArray(response)
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
    <h1 class="page-title">Hasil Evaluasi Anak</h1>
    <p class="page-subtitle">Lihat perkembangan evaluasi tahap 1-3 dari pembimbing untuk setiap murid Anda.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Riwayat Evaluasi</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchKeyword" placeholder="Cari murid / pembimbing / progres..." />
      </div>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>

      <div class="table-wrap" v-else>
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Kode</th>
              <th>Nama Murid</th>
              <th>Pembimbing</th>
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
              <td>{{ item.nama_pembimbing || '-' }}</td>
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
              <td colspan="8">Belum ada data evaluasi untuk ditampilkan.</td>
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
</style>
