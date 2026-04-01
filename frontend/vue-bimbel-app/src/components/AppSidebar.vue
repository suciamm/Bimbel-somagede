<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import logoBimbel from '../assets/logo-bimbel.svg'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const route = useRoute()
const authStore = useAuthStore()

const adminNavItems = [
  { to: '/', label: 'Dashboard', icon: 'space_dashboard' },
  { to: '/murid', label: 'Data Murid', icon: 'school' },
  { to: '/ortu', label: 'Data Orang Tua', icon: 'family_restroom' },
  { to: '/pembimbing', label: 'Data Pembimbing', icon: 'co_present' },
  { to: '/jadwal', label: 'Jadwal Mengajar', icon: 'calendar_month' },
  { to: '/paket', label: 'Paket Bimbingan', icon: 'inventory_2' },
  { to: '/transaksi', label: 'Transaksi Pembayaran', icon: 'payments' },
  { to: '/absensi', label: 'Rekap Absensi', icon: 'fact_check' },
  { to: '/rekap-murid', label: 'Rekap Murid Bulanan', icon: 'assessment' }
]

const pembimbingNavItems = [
  { to: '/', label: 'Dashboard', icon: 'space_dashboard' },
  { to: '/pembimbing/jadwal-hari-ini', label: 'Jadwal Mengajar Hari Ini', icon: 'today' },
  { to: '/pembimbing/jadwal-saya', label: 'Jadwal Mengajar Saya', icon: 'calendar_month' },
  { to: '/pembimbing/murid-saya', label: 'Data Murid Saya', icon: 'school' },
  { to: '/pembimbing/evaluasi', label: 'Evaluasi Murid', icon: 'grading' }
]

const orangtuaNavItems = [
  { to: '/', label: 'Dashboard', icon: 'space_dashboard' },
  { to: '/orangtua/jadwal-bimbingan', label: 'Jadwal Bimbingan', icon: 'event_note' },
  { to: '/orangtua/data-murid', label: 'Data Murid', icon: 'school' },
  { to: '/orangtua/evaluasi', label: 'Hasil Evaluasi', icon: 'fact_check' },
  { type: 'section', label: 'Data Pembayaran' },
  { to: '/orangtua/paket-bimbingan', label: 'Paket Bimbingan', icon: 'inventory_2' },
  { to: '/orangtua/paket-bimbingan-saya', label: 'Paket Bimbingan Saya', icon: 'payments' }
]

const navItems = computed(() => {
  if (authStore.user?.role === 'pembimbing') return pembimbingNavItems
  if (authStore.user?.role === 'orangtua') return orangtuaNavItems
  return adminNavItems
})

const consoleTitle = computed(() => {
  if (authStore.user?.role === 'orangtua') return 'Orang Tua'
  return authStore.user?.role === 'pembimbing' ? 'Pembimbing' : 'Admin'
})

const activePath = computed(() => route.path)

function handleNavClick() {
  if (window.innerWidth <= 980) {
    emit('close')
  }
}
</script>

<template>
  <div :class="['sidebar-backdrop', props.isOpen ? 'show' : '']" @click="emit('close')"></div>

  <aside :class="['sidebar', 'panel', props.isOpen ? 'open' : '']">
    <button class="btn btn-danger mobile-close" @click="emit('close')">
      <span class="material-symbols-rounded">close</span>
    </button>

    <div class="brand">
      <img :src="logoBimbel" alt="Logo Bimbel" class="brand-logo" />
      <h1>{{ consoleTitle }}</h1>
      <p class="brand-caption">Sistem manajemen operasional bimbingan belajar.</p>
    </div>

    <nav class="menu">
      <template v-for="item in navItems" :key="item.to || item.label">
        <p v-if="item.type === 'section'" class="menu-section">{{ item.label }}</p>
        <RouterLink
          v-else
          :to="item.to"
          :class="['menu-item', activePath === item.to ? 'active' : '']"
          @click="handleNavClick"
        >
          <span class="material-symbols-rounded nav-icon">{{ item.icon }}</span>
          <span>{{ item.label }}</span>
        </RouterLink>
      </template>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar-backdrop {
  display: none;
}

.sidebar {
  height: calc(100vh - 32px);
  min-height: calc(100vh - 32px);
  padding: 20px;
  position: sticky;
  top: 16px;
  background: linear-gradient(180deg, #ffffff 0%, #f8faff 100%);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.mobile-close {
  display: none;
}

.brand {
  margin-bottom: 24px;
  position: sticky;
  top: 0;
  z-index: 2;
  background: linear-gradient(180deg, #ffffff 0%, #f8faff 100%);
}

.brand-logo {
  width: 72px;
  height: 72px;
  object-fit: contain;
  display: block;
  margin-bottom: 10px;
}

.brand h1 {
  margin: 6px 0 0;
  font-size: 1.35rem;
  line-height: 1.2;
  font-family: 'Sora', sans-serif;
}

.brand-caption {
  margin: 8px 0 0;
  color: var(--text-muted);
  font-size: 0.82rem;
}

.menu {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-bottom: 12px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  scrollbar-gutter: stable;
}

.menu-item {
  padding: 10px 12px;
  border-radius: 10px;
  font-weight: 600;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  gap: 10px;
  border: 1px solid transparent;
}

.menu-item.active,
.menu-item:hover {
  background: #eef2ff;
  color: #1e3a8a;
  border-color: #bfdbfe;
}

.menu-section {
  margin: 10px 6px 2px;
  font-size: 0.72rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  font-weight: 700;
}

.nav-icon {
  background: #dbeafe;
  color: #1d4ed8;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: inline-grid;
  place-items: center;
  font-size: 16px;
}

@media (max-width: 980px) {
  .sidebar-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    z-index: 998;
    background: rgba(15, 23, 42, 0.32);
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s ease;
  }

  .sidebar-backdrop.show {
    opacity: 1;
    pointer-events: auto;
  }

  .sidebar {
    min-height: 100vh;
    height: 100dvh;
    position: fixed;
    top: 0;
    left: 0;
    width: min(320px, 88vw);
    border-radius: 0 16px 16px 0;
    z-index: 999;
    transform: translateX(-110%);
    transition: transform 0.24s ease;
    overflow-y: auto;
    overscroll-behavior: contain;
    -webkit-overflow-scrolling: touch;
    padding-bottom: 22px;
  }

  .sidebar.open {
    transform: translateX(0);
  }

  .mobile-close {
    display: inline-grid;
    place-items: center;
    width: 34px;
    height: 34px;
    padding: 0;
    border-radius: 8px;
    margin-left: auto;
  }

  .menu {
    flex-direction: column;
    flex-wrap: nowrap;
    min-height: max-content;
    overflow: visible;
  }
}
</style>
