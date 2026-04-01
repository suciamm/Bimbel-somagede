<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import { useAuthStore } from '../stores/auth'
import { getJadwalByPembimbingApi } from '../services/api'

const authStore = useAuthStore()
const loading = ref(false)
const errorMessage = ref('')
const jadwalList = ref([])
const searchJadwal = ref('')
const openDays = ref([])

const dayOrder = ['senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu', 'minggu']

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function formatJam(value) {
  const raw = String(value || '')
  const match = raw.match(/(\d{2}:\d{2})/)
  return match ? match[1] : '-'
}

function formatRentangWaktu(item) {
  return `${formatJam(item?.waktu_mulai)} - ${formatJam(item?.waktu_selesai)}`
}

function formatMuridDisplay(item) {
  const kode = String(item?.kode_murid || '').trim()
  const nama = String(item?.nama_murid || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  return nama || '-'
}

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function dayRank(day) {
  const idx = dayOrder.indexOf(normalizeText(day).trim())
  return idx === -1 ? 999 : idx
}

function sortByTime(a, b) {
  const timeA = formatJam(a?.waktu_mulai)
  const timeB = formatJam(b?.waktu_mulai)
  return timeA.localeCompare(timeB)
}

function makeSessionKey(item) {
  return [
    normalizeText(item?.hari_bimbingan).trim(),
    formatJam(item?.waktu_mulai),
    formatJam(item?.waktu_selesai),
    normalizeText(item?.ruangan).trim()
  ].join('|')
}

const filteredJadwalList = computed(() => {
  const keyword = normalizeText(searchJadwal.value).trim()
  if (!keyword) return jadwalList.value

  return jadwalList.value.filter((item) => {
    const joined = [item.hari_bimbingan, formatRentangWaktu(item), item.ruangan, formatMuridDisplay(item)]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const groupedJadwalByDay = computed(() => {
  const dayMap = new Map()

  filteredJadwalList.value.forEach((item) => {
    const dayName = String(item?.hari_bimbingan || '-').trim() || '-'
    const dayKey = normalizeText(dayName)

    if (!dayMap.has(dayKey)) {
      dayMap.set(dayKey, {
        key: dayKey,
        label: dayName,
        sessions: new Map(),
        muridTotal: 0
      })
    }

    const dayGroup = dayMap.get(dayKey)
    const sessionKey = makeSessionKey(item)

    if (!dayGroup.sessions.has(sessionKey)) {
      dayGroup.sessions.set(sessionKey, {
        id: item?.id_jadwal || sessionKey,
        waktu_mulai: item?.waktu_mulai,
        waktu_selesai: item?.waktu_selesai,
        ruangan: item?.ruangan || '-',
        muridList: []
      })
    }

    const session = dayGroup.sessions.get(sessionKey)
    const muridName = formatMuridDisplay(item)
    if (!session.muridList.includes(muridName)) {
      session.muridList.push(muridName)
      dayGroup.muridTotal += 1
    }
  })

  return [...dayMap.values()]
    .map((group) => ({
      ...group,
      sessionCount: group.sessions.size,
      sessions: [...group.sessions.values()].sort(sortByTime)
    }))
    .sort((a, b) => {
      const dayDiff = dayRank(a.label) - dayRank(b.label)
      if (dayDiff !== 0) return dayDiff
      return a.label.localeCompare(b.label)
    })
})

function toggleDay(dayKey) {
  if (openDays.value.includes(dayKey)) {
    openDays.value = openDays.value.filter((key) => key !== dayKey)
    return
  }
  openDays.value = [...openDays.value, dayKey]
}

function isDayOpen(dayKey) {
  return openDays.value.includes(dayKey)
}

watch(
  groupedJadwalByDay,
  (groups) => {
    const allowed = new Set(groups.map((group) => group.key))
    openDays.value = openDays.value.filter((key) => allowed.has(key))
    if (!openDays.value.length && groups.length) {
      openDays.value = [groups[0].key]
    }
  },
  { immediate: true }
)

function downloadJadwalSayaPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = groupedJadwalByDay.value.flatMap((group) =>
    group.sessions.map((session, index) => [
      index === 0 ? group.label : '',
      formatRentangWaktu(session),
      session.ruangan || '-',
      session.muridList.join(', ') || '-'
    ])
  )

  doc.setFontSize(12)
  doc.text('Jadwal Mengajar Saya', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['Hari', 'Waktu', 'Ruangan', 'Murid']],
    body: rows.length ? rows : [['-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('jadwal-mengajar-saya.pdf')
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await getJadwalByPembimbingApi(authStore.user?.id_user)
    jadwalList.value = toArray(response)
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
    <h1 class="page-title">Jadwal Mengajar Saya</h1>
    <p class="page-subtitle">Menampilkan seluruh jadwal bimbingan yang ditangani pembimbing.</p>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Jadwal Mengajar</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchJadwal" placeholder="Cari jadwal mengajar..." />
        <div class="tools-actions">
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF jadwal mengajar"
            @click="downloadJadwalSayaPdf"
          >
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>

      <div class="grouped-wrap" v-else>
        <article v-for="group in groupedJadwalByDay" :key="group.key" class="day-group">
          <button type="button" class="day-header" @click="toggleDay(group.key)">
            <div>
              <h3>{{ group.label }}</h3>
              <p>{{ group.sessionCount }} sesi • {{ group.muridTotal }} murid</p>
            </div>
            <span class="chevron" :class="{ open: isDayOpen(group.key) }" aria-hidden="true">⌄</span>
          </button>

          <div v-if="isDayOpen(group.key)" class="session-list">
            <section v-for="session in group.sessions" :key="session.id" class="session-card">
              <header class="session-head">
                <strong>{{ formatRentangWaktu(session) }}</strong>
                <span>{{ session.ruangan || '-' }}</span>
              </header>

              <div class="murid-list">
                <span v-for="murid in session.muridList" :key="murid" class="murid-chip">
                  {{ murid }}
                </span>
              </div>
            </section>
          </div>
        </article>

        <p v-if="!loading && groupedJadwalByDay.length === 0" class="empty-state">
          Tidak ada jadwal mengajar.
        </p>
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

.grouped-wrap {
  display: grid;
  gap: 12px;
}

.day-group {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #ffffff;
  overflow: hidden;
}

.day-header {
  width: 100%;
  border: none;
  background: linear-gradient(135deg, #f8fafc, #eff6ff);
  padding: 14px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  text-align: left;
  cursor: pointer;
}

.day-header h3 {
  margin: 0;
  font-size: 1rem;
  color: #0f172a;
}

.day-header p {
  margin: 4px 0 0;
  font-size: 0.84rem;
  color: #475569;
}

.chevron {
  font-size: 1.2rem;
  color: #334155;
  transition: transform 0.18s ease;
}

.chevron.open {
  transform: rotate(180deg);
}

.session-list {
  display: grid;
  gap: 10px;
  padding: 12px;
}

.session-card {
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 10px;
  background: #f8fafc;
}

.session-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.session-head strong {
  color: #0f172a;
}

.session-head span {
  font-size: 0.85rem;
  color: #475569;
}

.murid-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.murid-chip {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  border: 1px solid #bfdbfe;
  background: #eff6ff;
  color: #1e3a8a;
  font-size: 0.82rem;
  padding: 4px 10px;
}

.empty-state {
  text-align: center;
  margin: 0;
  color: #64748b;
  padding: 20px 10px;
}

@media (max-width: 700px) {
  .session-head {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
