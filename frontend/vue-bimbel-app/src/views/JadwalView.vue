<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import BaseModal from '../components/BaseModal.vue'
import { deleteJadwalApi, editJadwalApi, getJadwalApi, getMuridAktifApi, getPembimbingApi, tambahJadwalApi } from '../services/api'

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const jadwalList = ref([])
const pembimbingList = ref([])
const muridList = ref([])
const showCreateModal = ref(false)
const showEditModal = ref(false)
const searchJadwal = ref('')
const currentPageJadwal = ref(1)
const PAGE_SIZE = 10
const editingJadwalId = ref('')

const form = reactive({
  id_pembimbing: '',
  jumlah_murid: 1,
  id_murid_list: [''],
  hari_bimbingan: 'senin',
  waktu_mulai: '15:00',
  waktu_selesai: '16:00',
  ruangan: ''
})

const editForm = reactive({
  id_pembimbing: '',
  id_murid: '',
  hari_bimbingan: 'senin',
  waktu_mulai: '15:00',
  waktu_selesai: '16:00',
  ruangan: ''
})

function normalizeText(value) {
  return String(value || '').toLowerCase()
}

function toArray(value) {
  if (Array.isArray(value)) return value
  if (Array.isArray(value?.data)) return value.data
  return []
}

function toTimeWithSeconds(value) {
  const raw = String(value || '').trim()
  if (!raw) return ''
  if (/^\d{2}:\d{2}:\d{2}$/.test(raw)) return raw
  if (/^\d{2}:\d{2}$/.test(raw)) return `${raw}:00`
  return raw
}

function formatJam(value) {
  const raw = String(value || '').trim()
  if (!raw) return '-'

  const match = raw.match(/(\d{2}:\d{2})/)
  return match ? match[1] : raw
}

function formatRentangWaktu(item) {
  const mulai = formatJam(item?.waktu_mulai)
  const selesai = formatJam(item?.waktu_selesai)

  if (mulai === '-' && selesai === '-') return '-'
  return `${mulai} - ${selesai}`
}

function getIdPembimbing(item) {
  return item?.id_pembimbing || item?.IDPembimbing || item?.id_user || ''
}

function getIdMurid(item) {
  return item?.id_murid || item?.IDMurid || ''
}

function resolvePembimbingId(item) {
  const directId = getIdPembimbing(item)
  if (directId) return String(directId)

  const kode = String(item?.kode_pembimbing || item?.KodePembimbing || '').trim().toLowerCase()
  const nama = String(item?.nama_pembimbing || item?.NamaPembimbing || '').trim().toLowerCase()

  const found = pembimbingList.value.find((p) => {
    const kodeP = String(p?.kode_pembimbing || p?.KodePembimbing || '').trim().toLowerCase()
    const namaP = String(p?.nama_lengkap || p?.NamaLengkap || p?.nama_pembimbing || '').trim().toLowerCase()
    return (kode && kodeP && kode === kodeP) || (nama && namaP && nama === namaP)
  })

  return found ? String(found.id_user || found.id_pembimbing || '') : ''
}

function resolveMuridId(item) {
  const directId = getIdMurid(item)
  if (directId) return String(directId)

  const kode = String(item?.kode_murid || item?.KodeMurid || '').trim().toLowerCase()
  const nama = String(item?.nama_murid || item?.NamaMurid || '').trim().toLowerCase()

  const found = muridList.value.find((m) => {
    const kodeM = String(m?.kode_murid || m?.KodeMurid || '').trim().toLowerCase()
    const namaM = String(m?.nama_murid || m?.NamaMurid || '').trim().toLowerCase()
    return (kode && kodeM && kode === kodeM) || (nama && namaM && nama === namaM)
  })

  return found ? String(found.id_murid || found.IDMurid || '') : ''
}

function formatPembimbingOption(item) {
  const kode = String(item?.kode_pembimbing || item?.KodePembimbing || '').trim()
  const nama = String(item?.nama_lengkap || item?.NamaLengkap || item?.nama_pembimbing || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  if (nama) return nama
  return `ID ${item?.id_user || item?.id_pembimbing || '-'}`
}

function formatMuridOption(item) {
  const kode = String(item?.kode_murid || item?.KodeMurid || '').trim()
  const nama = String(item?.nama_murid || item?.NamaMurid || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  if (nama) return nama
  return `ID ${item?.id_murid || '-'}`
}

function formatMuridDisplay(item) {
  const kode = String(item?.kode_murid || item?.KodeMurid || '').trim()
  const nama = String(item?.nama_murid || item?.NamaMurid || '').trim()
  if (kode && nama) return `(${kode})${nama}`
  if (nama) return nama
  return '-'
}

function getPembimbingById(idPembimbing) {
  const id = String(idPembimbing || '').trim()
  if (!id) return null
  return pembimbingList.value.find((item) => String(item?.id_user || item?.id_pembimbing || '') === id) || null
}

function getMuridById(idMurid) {
  const id = String(idMurid || '').trim()
  if (!id) return null
  return muridList.value.find((item) => String(item?.id_murid || '') === id) || null
}

function toMinutes(value) {
  const jam = formatJam(value)
  const match = String(jam || '').match(/^(\d{2}):(\d{2})$/)
  if (!match) return null
  return Number(match[1]) * 60 + Number(match[2])
}

function isTimeOverlap(startA, endA, startB, endB) {
  const aStart = toMinutes(startA)
  const aEnd = toMinutes(endA)
  const bStart = toMinutes(startB)
  const bEnd = toMinutes(endB)
  if (aStart === null || aEnd === null || bStart === null || bEnd === null) return false
  return aStart < bEnd && aEnd > bStart
}

function muridBentrokJadwal(idMurid, hariBimbingan, waktuMulai, waktuSelesai) {
  const murid = getMuridById(idMurid)
  if (!murid) return false

  const kodeMurid = String(murid?.kode_murid || '').trim().toLowerCase()
  const namaMurid = String(murid?.nama_murid || '').trim().toLowerCase()
  const targetHari = String(hariBimbingan || '').trim().toLowerCase()

  return jadwalList.value.some((jadwal) => {
    const hari = String(jadwal?.hari_bimbingan || '').trim().toLowerCase()
    const kode = String(jadwal?.kode_murid || '').trim().toLowerCase()
    const nama = String(jadwal?.nama_murid || '').trim().toLowerCase()

    const sameMurid = (kodeMurid && kode && kodeMurid === kode) || (namaMurid && nama && namaMurid === nama)
    if (!(sameMurid && hari === targetHari)) return false
    return isTimeOverlap(waktuMulai, waktuSelesai, jadwal?.waktu_mulai, jadwal?.waktu_selesai)
  })
}

function pembimbingBentrokJadwal(idPembimbing, hariBimbingan, waktuMulai, waktuSelesai) {
  const pembimbing = getPembimbingById(idPembimbing)
  if (!pembimbing) return false

  const kodePembimbing = String(pembimbing?.kode_pembimbing || '').trim().toLowerCase()
  const namaPembimbing = String(pembimbing?.nama_lengkap || '').trim().toLowerCase()
  const targetHari = String(hariBimbingan || '').trim().toLowerCase()

  return jadwalList.value.some((jadwal) => {
    const hari = String(jadwal?.hari_bimbingan || '').trim().toLowerCase()
    const kode = String(jadwal?.kode_pembimbing || '').trim().toLowerCase()
    const nama = String(jadwal?.nama_pembimbing || '').trim().toLowerCase()

    const samePembimbing = (kodePembimbing && kode && kodePembimbing === kode) || (namaPembimbing && nama && namaPembimbing === nama)
    if (!(samePembimbing && hari === targetHari)) return false
    return isTimeOverlap(waktuMulai, waktuSelesai, jadwal?.waktu_mulai, jadwal?.waktu_selesai)
  })
}

function ruanganBentrokJadwal(idPembimbing, ruangan, hariBimbingan, waktuMulai, waktuSelesai) {
  const room = String(ruangan || '').trim().toLowerCase()
  const targetHari = String(hariBimbingan || '').trim().toLowerCase()
  if (!room || !targetHari) return false

  const selectedPembimbing = getPembimbingById(idPembimbing)
  const selectedKode = String(selectedPembimbing?.kode_pembimbing || '').trim().toLowerCase()
  const selectedNama = String(selectedPembimbing?.nama_lengkap || '').trim().toLowerCase()

  return jadwalList.value.some((jadwal) => {
    const hari = String(jadwal?.hari_bimbingan || '').trim().toLowerCase()
    const roomJadwal = String(jadwal?.ruangan || '').trim().toLowerCase()
    if (!(hari === targetHari && roomJadwal === room)) return false

    const kodeJadwal = String(jadwal?.kode_pembimbing || '').trim().toLowerCase()
    const namaJadwal = String(jadwal?.nama_pembimbing || '').trim().toLowerCase()
    const isSamePembimbing =
      (selectedKode && kodeJadwal && selectedKode === kodeJadwal) ||
      (selectedNama && namaJadwal && selectedNama === namaJadwal)

    if (isSamePembimbing) return false
    return isTimeOverlap(waktuMulai, waktuSelesai, jadwal?.waktu_mulai, jadwal?.waktu_selesai)
  })
}

function ruanganBentrokJadwalForEdit(idJadwal, idPembimbing, ruangan, hariBimbingan, waktuMulai, waktuSelesai) {
  const targetId = String(idJadwal || '').trim()
  const room = String(ruangan || '').trim().toLowerCase()
  const targetHari = String(hariBimbingan || '').trim().toLowerCase()
  if (!targetId || !room || !targetHari) return false

  const selectedPembimbing = getPembimbingById(idPembimbing)
  const selectedKode = String(selectedPembimbing?.kode_pembimbing || '').trim().toLowerCase()
  const selectedNama = String(selectedPembimbing?.nama_lengkap || '').trim().toLowerCase()

  return jadwalList.value.some((jadwal) => {
    const jadwalId = String(jadwal?.id_jadwal || jadwal?.IDJadwal || '').trim()
    if (jadwalId === targetId) return false

    const hari = String(jadwal?.hari_bimbingan || '').trim().toLowerCase()
    const roomJadwal = String(jadwal?.ruangan || '').trim().toLowerCase()
    if (!(hari === targetHari && roomJadwal === room)) return false

    const kodeJadwal = String(jadwal?.kode_pembimbing || '').trim().toLowerCase()
    const namaJadwal = String(jadwal?.nama_pembimbing || '').trim().toLowerCase()
    const isSamePembimbing =
      (selectedKode && kodeJadwal && selectedKode === kodeJadwal) ||
      (selectedNama && namaJadwal && selectedNama === namaJadwal)

    if (isSamePembimbing) return false
    return isTimeOverlap(waktuMulai, waktuSelesai, jadwal?.waktu_mulai, jadwal?.waktu_selesai)
  })
}

const maxJumlahMurid = computed(() => {
  const total = muridList.value.length
  if (!total) return 1
  return Math.min(total, 10)
})

const jumlahMuridOptions = computed(() => {
  return Array.from({ length: maxJumlahMurid.value }, (_, i) => i + 1)
})

const muridConflictIndexes = computed(() => {
  return form.id_murid_list.reduce((acc, idMurid, index) => {
    if (!String(idMurid || '').trim()) return acc
    if (muridBentrokJadwal(idMurid, form.hari_bimbingan, form.waktu_mulai, form.waktu_selesai)) {
      acc.push(index)
    }
    return acc
  }, [])
})

const muridConflictMessages = computed(() => {
  return muridConflictIndexes.value.map((index) => {
    const murid = getMuridById(form.id_murid_list[index])
    return `${formatMuridOption(murid)} bentrok dengan jadwal lain pada ${form.hari_bimbingan} di jam tersebut.`
  })
})

const pembimbingConflict = computed(() => {
  return pembimbingBentrokJadwal(form.id_pembimbing, form.hari_bimbingan, form.waktu_mulai, form.waktu_selesai)
})

const ruanganConflict = computed(() => {
  return ruanganBentrokJadwal(form.id_pembimbing, form.ruangan, form.hari_bimbingan, form.waktu_mulai, form.waktu_selesai)
})

const editRuanganConflict = computed(() => {
  return ruanganBentrokJadwalForEdit(
    editingJadwalId.value,
    editForm.id_pembimbing,
    editForm.ruangan,
    editForm.hari_bimbingan,
    editForm.waktu_mulai,
    editForm.waktu_selesai
  )
})

function isMuridConflict(index) {
  return muridConflictIndexes.value.includes(index)
}

function syncMuridSlots() {
  const target = Number(form.jumlah_murid) || 1
  while (form.id_murid_list.length < target) {
    form.id_murid_list.push('')
  }
  while (form.id_murid_list.length > target) {
    form.id_murid_list.pop()
  }
}

function getMuridOptionsForSlot(index) {
  const selectedOther = new Set(
    form.id_murid_list.filter((value, i) => i !== index && String(value || '').trim() !== '')
  )

  return muridList.value.filter((item) => {
    const id = String(item?.id_murid || '')
    const currentValue = String(form.id_murid_list[index] || '')
    return !selectedOther.has(id) || id === currentValue
  })
}

const filteredJadwal = computed(() => {
  const keyword = normalizeText(searchJadwal.value).trim()
  if (!keyword) return jadwalList.value

  return jadwalList.value.filter((item) => {
    const joined = [
      item.hari_bimbingan,
      item.waktu_mulai,
      item.waktu_selesai,
      item.ruangan,
      item.nama_pembimbing,
      item.kode_murid,
      item.nama_murid
    ]
      .map(normalizeText)
      .join(' ')
    return joined.includes(keyword)
  })
})

const groupedJadwalByHari = computed(() => {
  const hariUrut = ['senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu', 'minggu']
  const hariMap = new Map()

  filteredJadwal.value.forEach((item) => {
    const hari = String(item?.hari_bimbingan || '-').toLowerCase()
    if (!hariMap.has(hari)) {
      hariMap.set(hari, {
        hari,
        sessions: new Map(),
        totalMurid: 0
      })
    }

    const group = hariMap.get(hari)
    const sessionKey = [
      hari,
      formatJam(item?.waktu_mulai),
      formatJam(item?.waktu_selesai),
      String(item?.ruangan || '').trim().toLowerCase(),
      String(item?.nama_pembimbing || '').trim().toLowerCase()
    ].join('|')

    if (!group.sessions.has(sessionKey)) {
      group.sessions.set(sessionKey, {
        key: sessionKey,
        waktu: formatRentangWaktu(item),
        ruangan: item?.ruangan || '-',
        pembimbing: item?.nama_pembimbing || '-',
        rows: []
      })
    }

    group.sessions.get(sessionKey).rows.push(item)
    group.totalMurid += 1
  })

  return [...hariMap.values()]
    .map((group) => ({
      ...group,
      hariLabel: group.hari === '-' ? '-' : group.hari.charAt(0).toUpperCase() + group.hari.slice(1),
      sessions: [...group.sessions.values()]
    }))
    .sort((a, b) => {
      const idxA = hariUrut.indexOf(a.hari)
      const idxB = hariUrut.indexOf(b.hari)
      return (idxA === -1 ? 999 : idxA) - (idxB === -1 ? 999 : idxB)
    })
})

const totalPagesJadwal = computed(() => Math.max(1, Math.ceil(filteredJadwal.value.length / PAGE_SIZE)))

const paginatedJadwal = computed(() => {
  const start = (currentPageJadwal.value - 1) * PAGE_SIZE
  return filteredJadwal.value.slice(start, start + PAGE_SIZE)
})

function nomorUrut(index) {
  return (currentPageJadwal.value - 1) * PAGE_SIZE + index + 1
}

function goPrevPage() {
  if (currentPageJadwal.value > 1) currentPageJadwal.value -= 1
}

function goNextPage() {
  if (currentPageJadwal.value < totalPagesJadwal.value) currentPageJadwal.value += 1
}

function downloadJadwalPdf() {
  const doc = new jsPDF({ orientation: 'landscape' })
  const rows = filteredJadwal.value.map((item, index) => [
    index + 1,
    item.hari_bimbingan || '-',
    formatRentangWaktu(item),
    item.ruangan || '-',
    item.nama_pembimbing || '-',
    formatMuridDisplay(item)
  ])

  doc.setFontSize(12)
  doc.text('Daftar Jadwal Mengajar', 14, 14)
  autoTable(doc, {
    startY: 20,
    head: [['No', 'Hari', 'Waktu', 'Ruangan', 'Pembimbing', 'Murid']],
    body: rows.length ? rows : [['-', '-', '-', '-', '-', '-']],
    styles: { fontSize: 9 }
  })

  doc.save('daftar-jadwal.pdf')
}

async function loadJadwal() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [jadwalResponse, pembimbingResponse, muridResponse] = await Promise.all([
      getJadwalApi(),
      getPembimbingApi(),
      getMuridAktifApi()
    ])
    jadwalList.value = toArray(jadwalResponse)
    pembimbingList.value = toArray(pembimbingResponse)
    muridList.value = toArray(muridResponse)
    currentPageJadwal.value = 1
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

watch(searchJadwal, () => {
  currentPageJadwal.value = 1
})

watch(filteredJadwal, () => {
  if (currentPageJadwal.value > totalPagesJadwal.value) {
    currentPageJadwal.value = totalPagesJadwal.value
  }
})

watch(
  () => form.jumlah_murid,
  () => {
    syncMuridSlots()
  }
)

watch(
  () => muridList.value.length,
  () => {
    if (form.jumlah_murid > maxJumlahMurid.value) {
      form.jumlah_murid = maxJumlahMurid.value
    }
    syncMuridSlots()
  }
)

async function submitJadwal() {
  successMessage.value = ''
  errorMessage.value = ''
  try {
    if (!form.id_pembimbing) {
      throw new Error('Pembimbing wajib dipilih.')
    }

    if (pembimbingConflict.value) {
      throw new Error('Jadwal pembimbing bentrok pada hari dan jam yang sama.')
    }

    if (ruanganConflict.value) {
      throw new Error('Ruangan sudah dipakai pembimbing lain pada hari dan jam yang sama.')
    }

    const muridIds = form.id_murid_list
      .map((id) => String(id || '').trim())
      .filter(Boolean)

    if (!muridIds.length || muridIds.length !== Number(form.jumlah_murid)) {
      throw new Error('Semua pilihan murid wajib diisi.')
    }

    const uniqueIds = [...new Set(muridIds)]
    if (uniqueIds.length !== muridIds.length) {
      throw new Error('Murid tidak boleh duplikat dalam satu jadwal.')
    }

    if (muridConflictIndexes.value.length > 0) {
      throw new Error('Ada murid yang bentrok pada hari dan jam yang sama. Silakan periksa kembali pilihan murid.')
    }

    const payloadBase = {
      id_pembimbing: Number(form.id_pembimbing),
      hari_bimbingan: form.hari_bimbingan,
      waktu_mulai: toTimeWithSeconds(form.waktu_mulai),
      waktu_selesai: toTimeWithSeconds(form.waktu_selesai),
      ruangan: form.ruangan
    }

    const results = await Promise.allSettled(
      uniqueIds.map((idMurid) =>
        tambahJadwalApi({
          ...payloadBase,
          id_murid: Number(idMurid)
        })
      )
    )

    const successCount = results.filter((r) => r.status === 'fulfilled').length
    const failed = results.filter((r) => r.status === 'rejected')

    if (!successCount) {
      const firstError = failed[0]?.reason?.message || 'Gagal menambahkan semua jadwal.'
      throw new Error(firstError)
    }

    if (failed.length) {
      const firstError = failed[0]?.reason?.message || 'Sebagian data gagal diproses.'
      successMessage.value = `Berhasil menambahkan ${successCount} jadwal, ${failed.length} gagal. ${firstError}`
    } else {
      successMessage.value = `${successCount} jadwal baru berhasil ditambahkan.`
    }

    form.id_pembimbing = ''
    form.jumlah_murid = 1
    form.id_murid_list = ['']
    form.hari_bimbingan = 'senin'
    form.waktu_mulai = '15:00'
    form.waktu_selesai = '16:00'
    form.ruangan = ''
    showCreateModal.value = false

    await loadJadwal()
  } catch (error) {
    errorMessage.value = error.message
  }
}

function openEditModal(item) {
  editingJadwalId.value = String(item?.id_jadwal || item?.IDJadwal || '')
  editForm.id_pembimbing = resolvePembimbingId(item)
  editForm.id_murid = resolveMuridId(item)
  editForm.hari_bimbingan = String(item?.hari_bimbingan || 'senin').toLowerCase()
  editForm.waktu_mulai = formatJam(item?.waktu_mulai) === '-' ? '15:00' : formatJam(item?.waktu_mulai)
  editForm.waktu_selesai = formatJam(item?.waktu_selesai) === '-' ? '16:00' : formatJam(item?.waktu_selesai)
  editForm.ruangan = item?.ruangan || ''
  showEditModal.value = true
}

async function submitEditJadwal() {
  successMessage.value = ''
  errorMessage.value = ''
  try {
    if (!editingJadwalId.value) {
      throw new Error('ID jadwal tidak ditemukan.')
    }
    if (!editForm.id_pembimbing || !editForm.id_murid) {
      throw new Error('Pembimbing dan murid wajib dipilih.')
    }

    if (editRuanganConflict.value) {
      throw new Error('Ruangan sudah dipakai pembimbing lain pada hari dan jam yang sama.')
    }

    await editJadwalApi(editingJadwalId.value, {
      id_pembimbing: Number(editForm.id_pembimbing),
      id_murid: Number(editForm.id_murid),
      hari_bimbingan: editForm.hari_bimbingan,
      waktu_mulai: toTimeWithSeconds(editForm.waktu_mulai),
      waktu_selesai: toTimeWithSeconds(editForm.waktu_selesai),
      ruangan: editForm.ruangan
    })

    successMessage.value = 'Jadwal berhasil diperbarui.'
    showEditModal.value = false
    editingJadwalId.value = ''
    await loadJadwal()
  } catch (error) {
    errorMessage.value = error.message
  }
}

async function handleDeleteJadwal(idJadwal) {
  if (!window.confirm('Hapus jadwal bimbingan ini?')) return

  successMessage.value = ''
  errorMessage.value = ''
  try {
    await deleteJadwalApi(idJadwal)
    successMessage.value = 'Jadwal berhasil dihapus.'
    await loadJadwal()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(loadJadwal)
</script>

<template>
  <section>
    <h1 class="page-title">Jadwal Mengajar</h1>
    <p class="page-subtitle">Atur alokasi jadwal pembimbing dan murid dengan rapi.</p>

    <BaseModal :show="showCreateModal" title="Tambah Jadwal Mengajar" @close="showCreateModal = false">
      <form class="form-grid" @submit.prevent="submitJadwal">
        <div class="form-section" style="grid-column: 1 / -1">
          <div class="section-title">Detail Jadwal</div>
          <div class="section-grid">
            <div class="field">
              <label>Pembimbing</label>
              <select v-model="form.id_pembimbing" required>
                <option disabled value="">Pilih pembimbing</option>
                <option v-for="item in pembimbingList" :key="item.id_user || item.id_pembimbing" :value="String(item.id_user || item.id_pembimbing)">
                  {{ formatPembimbingOption(item) }}
                </option>
              </select>
              <p v-if="pembimbingConflict" class="conflict-inline">
                Pembimbing ini bentrok dengan jadwal lain di hari {{ form.hari_bimbingan }} pada jam tersebut.
              </p>
            </div>

            <div class="field">
              <label>Hari Bimbingan</label>
              <select v-model="form.hari_bimbingan">
                <option value="senin">senin</option>
                <option value="selasa">selasa</option>
                <option value="rabu">rabu</option>
                <option value="kamis">kamis</option>
                <option value="jumat">jumat</option>
                <option value="sabtu">sabtu</option>
                <option value="minggu">minggu</option>
              </select>
            </div>

            <div class="field">
              <label>Ruangan</label>
              <input v-model="form.ruangan" required />
              <p v-if="ruanganConflict" class="conflict-inline">
                Ruangan ini sudah dipakai pembimbing lain di hari {{ form.hari_bimbingan }} pada jam tersebut.
              </p>
            </div>

            <div class="field">
              <label>Waktu Mulai (HH:mm)</label>
              <input v-model="form.waktu_mulai" type="time" required />
            </div>

            <div class="field">
              <label>Waktu Selesai (HH:mm)</label>
              <input v-model="form.waktu_selesai" type="time" required />
            </div>
          </div>
        </div>

        <div class="form-section form-section-student" style="grid-column: 1 / -1">
          <div class="section-title">Pengaturan Murid</div>
          <p class="section-hint">
            Pilih jumlah murid lalu tentukan siapa saja murid yang akan masuk ke jadwal ini.
            (Terpilih: {{ form.id_murid_list.filter((id) => String(id || '').trim() !== '').length }} / {{ form.jumlah_murid }})
          </p>

          <div class="field student-count-field">
            <label>Jumlah Murid</label>
            <select v-model.number="form.jumlah_murid">
              <option v-for="jumlah in jumlahMuridOptions" :key="jumlah" :value="jumlah">{{ jumlah }}</option>
            </select>
          </div>

          <div v-if="muridConflictMessages.length > 0" class="conflict-box">
            <strong>Konflik Jadwal Terdeteksi:</strong>
            <p v-for="(msg, idx) in muridConflictMessages" :key="`conflict-${idx}`" class="conflict-item">{{ msg }}</p>
          </div>

          <div class="student-picker-grid">
            <div class="field" :class="{ 'field-conflict': isMuridConflict(index) }" v-for="(_, index) in form.id_murid_list" :key="`murid-slot-${index}`">
              <label>Murid {{ index + 1 }}</label>
              <select v-model="form.id_murid_list[index]" required>
                <option disabled value="">Pilih murid</option>
                <option v-for="item in getMuridOptionsForSlot(index)" :key="item.id_murid" :value="String(item.id_murid)">
                  {{ formatMuridOption(item) }}
                </option>
              </select>
              <p v-if="isMuridConflict(index)" class="conflict-inline">
                Murid ini bentrok dengan jadwal lain di hari {{ form.hari_bimbingan }} pada jam tersebut.
              </p>
            </div>
          </div>
        </div>

        <div>
          <button
            class="btn btn-primary"
            type="submit"
            :disabled="loading || muridConflictIndexes.length > 0 || pembimbingConflict || ruanganConflict"
          >
            Simpan Jadwal
          </button>
        </div>
      </form>

      <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
      <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>
    </BaseModal>

    <BaseModal :show="showEditModal" title="Ubah Jadwal Mengajar" @close="showEditModal = false">
      <form class="form-grid" @submit.prevent="submitEditJadwal">
        <div class="field">
          <label>Pembimbing</label>
          <select v-model="editForm.id_pembimbing" required>
            <option disabled value="">Pilih pembimbing</option>
            <option v-for="item in pembimbingList" :key="item.id_user || item.id_pembimbing" :value="String(item.id_user || item.id_pembimbing)">
              {{ formatPembimbingOption(item) }}
            </option>
          </select>
        </div>

        <div class="field">
          <label>Murid</label>
          <select v-model="editForm.id_murid" required>
            <option disabled value="">Pilih murid</option>
            <option v-for="item in muridList" :key="item.id_murid" :value="String(item.id_murid)">
              {{ formatMuridOption(item) }}
            </option>
          </select>
        </div>

        <div class="field">
          <label>Hari Bimbingan</label>
          <select v-model="editForm.hari_bimbingan">
            <option value="senin">senin</option>
            <option value="selasa">selasa</option>
            <option value="rabu">rabu</option>
            <option value="kamis">kamis</option>
            <option value="jumat">jumat</option>
            <option value="sabtu">sabtu</option>
            <option value="minggu">minggu</option>
          </select>
        </div>

        <div class="field">
          <label>Ruangan</label>
          <input v-model="editForm.ruangan" required />
          <p v-if="editRuanganConflict" class="conflict-inline">
            Ruangan ini sudah dipakai pembimbing lain di hari {{ editForm.hari_bimbingan }} pada jam tersebut.
          </p>
        </div>

        <div class="field">
          <label>Waktu Mulai (HH:mm)</label>
          <input v-model="editForm.waktu_mulai" type="time" required />
        </div>

        <div class="field">
          <label>Waktu Selesai (HH:mm)</label>
          <input v-model="editForm.waktu_selesai" type="time" required />
        </div>

        <div>
          <button class="btn btn-primary btn-edit-action" type="submit" :disabled="editRuanganConflict">Update Jadwal</button>
        </div>
      </form>
    </BaseModal>

    <section class="panel block">
      <header class="block-header">
        <h2>Daftar Jadwal</h2>
      </header>

      <div class="table-tools">
        <input v-model="searchJadwal" placeholder="Cari jadwal..." />
        <div class="tools-actions">
          <button class="btn btn-primary" type="button" @click="showCreateModal = true">Tambah Jadwal</button>
          <button
            class="btn btn-secondary btn-icon btn-pdf"
            type="button"
            title="Download PDF"
            aria-label="Download PDF jadwal"
            @click="downloadJadwalPdf"
          >
            <span class="pdf-icon" aria-hidden="true">&#128424;</span>
            <span>PDF</span>
          </button>
        </div>
      </div>

      <div class="grouped-wrap">
        <article v-for="group in groupedJadwalByHari" :key="group.hari" class="day-group">
          <header class="day-group-head">
            <h3>{{ group.hariLabel }}</h3>
            <p>{{ group.sessions.length }} sesi · {{ group.totalMurid }} murid</p>
          </header>

          <div class="session-list">
            <section v-for="session in group.sessions" :key="session.key" class="session-card">
              <div class="session-head">
                <strong>{{ session.waktu }}</strong>
                <span class="session-room">{{ session.ruangan }}</span>
                <span class="session-tutor">{{ session.pembimbing }}</span>
              </div>

              <div class="murid-row-list">
                <div v-for="item in session.rows" :key="item.id_jadwal" class="murid-row-item">
                  <span class="murid-name">{{ formatMuridDisplay(item) }}</span>
                  <div class="actions">
                    <button class="btn btn-secondary btn-edit-action" @click="openEditModal(item)">Edit</button>
                    <button class="btn btn-danger" @click="handleDeleteJadwal(item.id_jadwal)">Hapus</button>
                  </div>
                </div>
              </div>
            </section>
          </div>
        </article>

        <p v-if="groupedJadwalByHari.length === 0" class="empty-text">Tidak ada data jadwal.</p>
      </div>

      <details class="detail-table-wrap">
        <summary>Lihat versi tabel detail</summary>
        <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>No</th>
              <th>Hari</th>
              <th>Waktu</th>
              <th>Ruangan</th>
              <th>Pembimbing</th>
              <th>Murid</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in paginatedJadwal" :key="item.id_jadwal || index">
              <td>{{ nomorUrut(index) }}</td>
              <td>{{ item.hari_bimbingan }}</td>
              <td>{{ formatRentangWaktu(item) }}</td>
              <td>{{ item.ruangan }}</td>
              <td>{{ item.nama_pembimbing || '-' }}</td>
              <td>{{ formatMuridDisplay(item) }}</td>
              <td class="actions">
                <button class="btn btn-secondary btn-edit-action" @click="openEditModal(item)">Edit</button>
                <button class="btn btn-danger" @click="handleDeleteJadwal(item.id_jadwal)">Hapus</button>
              </td>
            </tr>
            <tr v-if="paginatedJadwal.length === 0">
              <td colspan="7">Tidak ada data jadwal.</td>
            </tr>
          </tbody>
        </table>
        </div>
      </details>
      <div class="pagination-wrap" v-if="!loading">
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Sebelumnya"
          aria-label="Halaman sebelumnya"
          @click="goPrevPage"
          :disabled="currentPageJadwal === 1"
        >
          <<
        </button>
        <span>Halaman {{ currentPageJadwal }} / {{ totalPagesJadwal }} ({{ filteredJadwal.length }} data)</span>
        <button
          class="btn btn-secondary btn-icon"
          type="button"
          title="Berikutnya"
          aria-label="Halaman berikutnya"
          @click="goNextPage"
          :disabled="currentPageJadwal === totalPagesJadwal"
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

.form-section {
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 12px;
  background: #f8fbff;
}

.form-section-student {
  border-color: #bfdbfe;
  background: linear-gradient(180deg, #eff6ff 0%, #f8fbff 100%);
}

.section-title {
  margin-bottom: 8px;
  font-weight: 800;
  color: #1e3a8a;
}

.section-hint {
  margin: 0 0 10px;
  color: var(--text-muted);
  font-size: 0.86rem;
}

.section-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.student-count-field {
  max-width: 260px;
  margin-bottom: 10px;
}

.student-picker-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.conflict-box {
  margin: 0 0 10px;
  padding: 10px 12px;
  border: 1px solid #fca5a5;
  border-radius: 10px;
  background: #fff1f2;
  color: #9f1239;
}

.conflict-item {
  margin: 4px 0 0;
  font-size: 0.86rem;
}

.field-conflict select {
  border-color: #ef4444;
  background: #fff5f5;
}

.conflict-inline {
  margin: 4px 0 0;
  font-size: 0.8rem;
  color: #b91c1c;
}

.actions {
  display: flex;
  gap: 8px;
}

.grouped-wrap {
  display: grid;
  gap: 12px;
}

.day-group {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  overflow: hidden;
}

.day-group-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  background: linear-gradient(135deg, #f8fafc, #eff6ff);
  border-bottom: 1px solid #e2e8f0;
}

.day-group-head h3 {
  margin: 0;
  color: #0f172a;
  font-size: 1rem;
}

.day-group-head p {
  margin: 0;
  color: #475569;
  font-size: 0.84rem;
}

.session-list {
  display: grid;
  gap: 10px;
  padding: 10px;
}

.session-card {
  border: 1px solid #dbeafe;
  border-radius: 10px;
  background: #f8fbff;
  padding: 10px;
}

.session-head {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

.session-room,
.session-tutor {
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 3px 9px;
  font-size: 0.78rem;
}

.session-room {
  background: #e0f2fe;
  color: #075985;
}

.session-tutor {
  background: #ede9fe;
  color: #5b21b6;
}

.murid-row-list {
  display: grid;
  gap: 6px;
}

.murid-row-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #e5e7eb;
}

.murid-name {
  font-weight: 600;
  color: #1f2937;
}

.detail-table-wrap {
  margin-top: 12px;
}

.detail-table-wrap summary {
  cursor: pointer;
  color: #334155;
  font-weight: 700;
  margin-bottom: 8px;
}

.empty-text {
  text-align: center;
  color: #64748b;
  margin: 8px 0;
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
  .section-grid {
    grid-template-columns: 1fr;
  }

  .student-picker-grid {
    grid-template-columns: 1fr;
  }

  .student-count-field {
    max-width: none;
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
  }
}
</style>
