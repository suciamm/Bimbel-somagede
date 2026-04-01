import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { loginApi } from '../services/api'

const STORAGE_KEY = 'bimbel_auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const loading = ref(false)

  const token = computed(() => user.value?.token || null)
  const isAuthenticated = computed(() => Boolean(user.value?.id_user))

  function readAuthStorage() {
    try {
      return sessionStorage.getItem(STORAGE_KEY)
    } catch {
      return null
    }
  }

  function writeAuthStorage(value) {
    try {
      sessionStorage.setItem(STORAGE_KEY, value)
    } catch {
      // Ignore browser storage errors.
    }
  }

  function removeAuthStorage() {
    try {
      sessionStorage.removeItem(STORAGE_KEY)
    } catch {
      // Ignore browser storage errors.
    }
  }

  function hydrate() {
    const raw = readAuthStorage()
    if (!raw) return
    try {
      user.value = JSON.parse(raw)
    } catch {
      user.value = null
      removeAuthStorage()
    }
  }

  async function login(payload) {
    loading.value = true
    try {
      const response = await loginApi(payload)

      const isApproved =
        response?.status === true ||
        response?.status === 1 ||
        String(response?.status || '').toLowerCase() === 'true' ||
        String(response?.status || '') === '1' ||
        String(response?.status || '').toLowerCase() === 'aktif'

      // Akun non-admin wajib menunggu persetujuan admin sebelum bisa login.
      if (response?.role !== 'admin' && !isApproved) {
        user.value = null
        removeAuthStorage()
        throw new Error('Akun Anda sedang menunggu verifikasi admin. Silakan coba login kembali setelah disetujui.')
      }

      const profile = {
        id_user: response.id_user,
        username: response.username,
        role: response.role,
        nama_lengkap: response.nama_lengkap,
        status: response.status
      }
      user.value = profile
      writeAuthStorage(JSON.stringify(profile))
      return profile
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    removeAuthStorage()
  }

  hydrate()

  return {
    user,
    token,
    loading,
    isAuthenticated,
    login,
    logout,
    hydrate
  }
})
