<script setup>
import { computed, onMounted, ref } from 'vue'
import StatCard from '../components/StatCard.vue'
import {
  getBimbinganByOrtuApi,
  getJadwalApi,
  getJadwalByPembimbingApi,
  getMuridAktifApi,
  getMuridByPembimbingApi,
  getPaketAktifApi,
  getTransaksiApi,
  getTransaksiByOrtuApi
} from '../services/api'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const loading = ref(true)
const errorMessage = ref('')

const role = computed(() => authStore.user?.role || 'admin')
const currentUserId = computed(() => Number(authStore.user?.id_user || 0))

const statCards = ref([])
const listTitle = ref('')
const listColumns = ref([])
const listRows = ref([])

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function formatCurrency(value) {
  return `Rp ${Number(value || 0).toLocaleString('id-ID')}`
}

function extractDateOnly(value) {
  const raw = String(value || '')
  return raw.match(/\d{4}-\d{2}-\d{2}/)?.[0] || '-'
}

function formatDateOnly(value) {
  const iso = extractDateOnly(value)
  if (!iso || iso === '-') return '-'
  const [year, month, day] = iso.split('-')
  return `${day}-${month}-${year}`
}

function formatJam(value) {
  const raw = String(value || '')
  return raw.match(/(\d{2}:\d{2})/)?.[1] || '-'
}

function formatWaktuRange(item) {
  return `${formatJam(item?.waktu_mulai)} - ${formatJam(item?.waktu_selesai)}`
}

function getTodayHariName() {
  const map = ['minggu', 'senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu']
  return map[new Date().getDay()]
}

async function loadDashboard() {
  loading.value = true
  errorMessage.value = ''

  try {
    if (role.value === 'admin') {
      const [murid, jadwal, paket, transaksi] = await Promise.all([
        getMuridAktifApi(),
        getJadwalApi(),
        getPaketAktifApi(),
        getTransaksiApi()
      ])

      const muridList = toArray(murid)
      const jadwalList = toArray(jadwal)
      const paketList = toArray(paket)
      const transaksiList = toArray(transaksi)

      statCards.value = [
        { label: 'Murid Aktif', value: muridList.length, tone: 'ocean' },
        { label: 'Total Jadwal', value: jadwalList.length, tone: 'violet' },
        { label: 'Paket Aktif', value: paketList.length, tone: 'sunset' },
        { label: 'Transaksi', value: transaksiList.length, tone: 'ruby' }
      ]

      listTitle.value = 'Transaksi Terbaru'
      listColumns.value = ['Murid', 'Paket', 'Jumlah Bayar', 'Status']
      listRows.value = transaksiList.slice(0, 6).map((item) => [
        item?.nama_murid || '-',
        item?.nama_paket || '-',
        formatCurrency(item?.jumlah_bayar),
        item?.status || '-'
      ])
      return
    }

    if (role.value === 'pembimbing') {
      const [jadwalResp, muridResp] = await Promise.all([
        getJadwalByPembimbingApi(currentUserId.value),
        getMuridByPembimbingApi(currentUserId.value)
      ])

      const jadwalList = toArray(jadwalResp)
      const muridList = toArray(muridResp)
      const hariIni = getTodayHariName()
      const jadwalHariIni = jadwalList.filter((item) => normalizeText(item?.hari_bimbingan) === hariIni)

      statCards.value = [
        { label: 'Jadwal Saya', value: jadwalList.length, tone: 'violet' },
        { label: 'Jadwal Hari Ini', value: jadwalHariIni.length, tone: 'sunset' },
        { label: 'Murid Binaan', value: muridList.length, tone: 'ocean' },
        { label: 'Hari Ini', value: hariIni.charAt(0).toUpperCase() + hariIni.slice(1), tone: 'ruby' }
      ]

      listTitle.value = 'Jadwal Hari Ini'
      listColumns.value = ['Waktu', 'Ruangan', 'Murid']
      listRows.value = jadwalHariIni.slice(0, 6).map((item) => [
        formatWaktuRange(item),
        item?.ruangan || '-',
        item?.nama_murid || '-'
      ])
      return
    }

    const [bimbinganResp, transaksiResp] = await Promise.all([
      getBimbinganByOrtuApi(currentUserId.value),
      getTransaksiByOrtuApi(currentUserId.value)
    ])

    const bimbinganList = toArray(bimbinganResp)
    const transaksiList = toArray(transaksiResp)
    const uniqueMurid = new Set(bimbinganList.map((item) => Number(item?.id_murid || 0)).filter(Boolean))
    const today = extractDateOnly(new Date().toISOString())
    const paketAktifCount = transaksiList.filter((item) => {
      const status = normalizeText(item?.status).trim()
      const selesai = extractDateOnly(item?.tgl_selesai)
      return status !== 'gagal' && status !== 'failed' && status !== 'expire' && selesai >= today
    }).length

    statCards.value = [
      { label: 'Anak Terdaftar', value: uniqueMurid.size, tone: 'ocean' },
      { label: 'Jadwal Bimbingan', value: bimbinganList.length, tone: 'violet' },
      { label: 'Paket Aktif', value: paketAktifCount, tone: 'sunset' },
      { label: 'Riwayat Transaksi', value: transaksiList.length, tone: 'ruby' }
    ]

    listTitle.value = 'Paket Bimbingan Terbaru'
    listColumns.value = ['Murid', 'Periode', 'Bayar', 'Status']
    listRows.value = transaksiList.slice(0, 6).map((item) => [
      item?.nama_murid || '-',
      `${formatDateOnly(item?.tgl_mulai)} s/d ${formatDateOnly(item?.tgl_selesai)}`,
      formatCurrency(item?.jumlah_bayar),
      item?.status || '-'
    ])
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)

</script>

<template>
  <section>
    <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>

    <div class="stats-grid">
      <StatCard v-for="card in statCards" :key="card.label" :label="card.label" :value="card.value" :tone="card.tone" />
    </div>

    <section class="panel latest">
      <header>
        <h2>{{ listTitle }}</h2>
        <button class="btn btn-secondary" @click="loadDashboard">Muat Ulang</button>
      </header>

      <p v-if="loading">Memuat data dashboard...</p>

      <div v-else class="table-wrap">
        <table>
          <thead>
            <tr>
              <th v-for="column in listColumns" :key="column">{{ column }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, idx) in listRows" :key="idx">
              <td v-for="(cell, cidx) in row" :key="`${idx}-${cidx}`">{{ cell }}</td>
            </tr>
            <tr v-if="listRows.length === 0">
              <td :colspan="Math.max(1, listColumns.length)">Belum ada data untuk ditampilkan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </section>
</template>

<style scoped>
.stats-grid {
  margin-top: 18px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.latest {
  margin-top: 18px;
  padding: 16px;
}

.latest header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 12px;
}

.latest h2 {
  margin: 0;
  font-size: 1.05rem;
}

@media (max-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 620px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .latest header {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
