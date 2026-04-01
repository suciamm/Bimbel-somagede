<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const emit = defineEmits(['toggle-menu'])

const router = useRouter()
const authStore = useAuthStore()

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header class="header panel">
    <div class="identity">
      <button class="btn btn-secondary menu-toggle" @click="emit('toggle-menu')">
        <span class="material-symbols-rounded">menu</span>
      </button>
      <div>
        <p class="hello">Selamat datang,</p>
        <h2>{{ authStore.user?.nama_lengkap || 'Pengguna' }}</h2>
      </div>
    </div>

    <div class="meta">
      <span class="role">{{ authStore.user?.role || '-' }}</span>
      <button class="btn btn-secondary" @click="handleLogout">Keluar</button>
    </div>
  </header>
</template>

<style scoped>
.header {
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  background: linear-gradient(100deg, #ffffff 0%, #f7f9ff 100%);
}

.hello {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.86rem;
}

.identity {
  display: flex;
  align-items: center;
  gap: 10px;
}

.menu-toggle {
  display: none;
  width: 38px;
  height: 38px;
  padding: 0;
  border-radius: 10px;
}

.menu-toggle .material-symbols-rounded {
  font-size: 20px;
}

h2 {
  margin: 2px 0 0;
  font-size: 1.15rem;
  font-family: 'Sora', sans-serif;
}

.meta {
  display: flex;
  align-items: center;
  gap: 10px;
}

.role {
  background: #ede9fe;
  color: #5b21b6;
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 0.8rem;
  text-transform: uppercase;
  font-weight: 700;
}

@media (max-width: 640px) {
  .header {
    flex-direction: column;
    align-items: flex-start;
  }
}

@media (max-width: 980px) {
  .menu-toggle {
    display: inline-grid;
    place-items: center;
  }
}
</style>
