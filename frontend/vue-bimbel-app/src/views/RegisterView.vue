<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { registerApi } from '../services/api'
import logoBimbel from '../assets/logo-bimbel.svg'

const router = useRouter()

const form = reactive({
  username: '',
  password: '',
  nama_lengkap: '',
  role: 'orangtua',
  no_telp: '',
  alamat: ''
})

const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

async function submitRegister() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const response = await registerApi(form)
    successMessage.value = response?.message || 'Registrasi berhasil.'
    setTimeout(() => router.push('/login'), 700)
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <section class="auth-shell panel">
      <aside class="auth-branding">
        <div class="shape shape-blue" aria-hidden="true"></div>
        <div class="shape shape-green" aria-hidden="true"></div>
        <div class="shape shape-orange" aria-hidden="true"></div>
        <div class="shape shape-red" aria-hidden="true"></div>
        <div class="shape shape-soft" aria-hidden="true"></div>

        <div class="logo-frame">
          <img :src="logoBimbel" alt="Logo Bimbel Somagede" class="brand-logo" />
          <p class="brand-name">Bimbel Somagede</p>
          <p class="brand-slogan">Belajar Terarah, Prestasi Cerah</p>
        </div>
      </aside>

      <div class="auth-card">
        <div class="auth-card-inner">
          <p class="badge">Portal Akademik</p>
          <h1>Registrasi Akun</h1>
          <p class="subtitle">Lengkapi data diri Anda untuk melanjutkan proses pendaftaran.</p>

          <form @submit.prevent="submitRegister">
            <div class="form-grid">
              <div class="field">
                <label>Username</label>
                <input v-model="form.username" required />
              </div>

              <div class="field">
                <label>Password</label>
                <input v-model="form.password" type="password" required />
              </div>

              <div class="field">
                <label>Nama Lengkap</label>
                <input v-model="form.nama_lengkap" required />
              </div>

              <div class="field">
                <label>Role</label>
                <select v-model="form.role">
                  <option value="pembimbing">pembimbing</option>
                  <option value="orangtua">orangtua</option>
                </select>
              </div>

              <div class="field">
                <label>No. Telepon</label>
                <input v-model="form.no_telp" />
              </div>

              <div class="field">
                <label>Alamat</label>
                <input v-model="form.alamat" />
              </div>
            </div>

            <button class="btn btn-primary full" :disabled="loading" type="submit">
              {{ loading ? 'Memproses...' : 'Daftar' }}
            </button>
          </form>

          <p v-if="errorMessage" class="message message-error">{{ errorMessage }}</p>
          <p v-if="successMessage" class="message message-success">{{ successMessage }}</p>

          <p class="footnote">
            Sudah punya akun?
            <RouterLink to="/login">Masuk sekarang</RouterLink>
          </p>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  padding: 0;
}

.auth-shell {
  width: 100%;
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1.25fr 0.75fr;
  border: none;
  border-radius: 0;
  box-shadow: none;
}

.auth-branding {
  padding: 0;
  background: #ffffff;
  color: #f8fbff;
  position: relative;
  overflow: hidden;
  display: grid;
  place-items: center;
}

.shape {
  position: absolute;
  border-radius: 28px;
  opacity: 0.85;
}

.shape-blue {
  width: 40%;
  height: 45%;
  left: -7%;
  top: -6%;
  border-radius: 0 0 110px 0;
  background: linear-gradient(160deg, rgba(125, 211, 252, 0.62) 0%, rgba(59, 130, 246, 0.36) 100%);
}

.shape-green {
  width: 34%;
  height: 32%;
  right: -4%;
  top: 17%;
  border-radius: 120px 0 0 120px;
  background: linear-gradient(155deg, rgba(74, 222, 128, 0.56) 0%, rgba(34, 197, 94, 0.34) 100%);
}

.shape-orange {
  width: 44%;
  height: 24%;
  left: 11%;
  bottom: -4%;
  border-radius: 120px 120px 0 0;
  background: linear-gradient(155deg, rgba(253, 186, 116, 0.54) 0%, rgba(251, 146, 60, 0.32) 100%);
}

.shape-red {
  width: 26%;
  height: 26%;
  right: 10%;
  bottom: 12%;
  border-radius: 999px;
  background: linear-gradient(155deg, rgba(248, 113, 113, 0.5) 0%, rgba(239, 68, 68, 0.32) 100%);
}

.shape-soft {
  width: 72%;
  height: 72%;
  border-radius: 999px;
  background: radial-gradient(circle at center, rgba(148, 163, 184, 0.14) 0%, rgba(148, 163, 184, 0) 70%);
  opacity: 0.6;
}

.logo-frame {
  position: relative;
  z-index: 2;
  width: min(390px, 76%);
  aspect-ratio: 1 / 1;
  border-radius: 999px;
  border: 10px solid #ffffff;
  background:
    linear-gradient(170deg, rgba(255, 255, 255, 0.98) 0%, rgba(246, 250, 255, 0.98) 100%),
    #ffffff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 26px 48px rgba(15, 23, 42, 0.12);
}

.brand-logo {
  width: min(185px, 47%);
  object-fit: contain;
  filter: drop-shadow(0 10px 16px rgba(30, 64, 175, 0.18));
  animation: logoFloat 3.2s ease-in-out infinite;
}

.brand-name {
  margin: 0;
  color: #1e3a8a;
  font-size: clamp(1.25rem, 2.2vw, 1.7rem);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  font-weight: 800;
}

.brand-slogan {
  margin: 0;
  color: #475569;
  font-size: clamp(0.86rem, 1.25vw, 0.98rem);
  font-weight: 600;
  letter-spacing: 0.02em;
}

.auth-card {
  background: linear-gradient(180deg, #ffffff 0%, #f4f8ff 100%);
  display: grid;
  place-items: center;
  padding: 24px;
}

.auth-card-inner {
  width: min(700px, 100%);
  background: #ffffff;
  border: 1px solid #dbe4f6;
  border-radius: 22px;
  padding: 28px;
  box-shadow: 0 18px 36px rgba(30, 41, 59, 0.08);
}

.badge {
  margin: 0;
  color: #1d4ed8;
  font-size: 0.78rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-weight: 800;
}

h1 {
  margin: 10px 0 6px;
  font-size: 1.62rem;
  font-family: 'Sora', 'Plus Jakarta Sans', sans-serif;
  letter-spacing: -0.02em;
}

.subtitle {
  margin: 0 0 18px;
  color: var(--text-muted);
}

form {
  display: grid;
  gap: 14px;
}

.full {
  width: 100%;
  margin-top: 4px;
}

.footnote {
  margin: 14px 0 0;
  font-size: 0.9rem;
}

.footnote a {
  color: var(--primary);
  font-weight: 700;
}

@keyframes logoFloat {
  0% {
    transform: translateY(0);
  }

  50% {
    transform: translateY(-7px);
  }

  100% {
    transform: translateY(0);
  }
}

@media (max-width: 980px) {
  .auth-shell {
    grid-template-columns: 1fr;
  }

  .auth-branding {
    min-height: 280px;
  }

  .auth-card {
    padding: 28px 20px;
  }

  .auth-card-inner {
    padding: 24px;
  }

  .logo-frame {
    width: min(280px, 68%);
    border-width: 8px;
  }
}

@media (max-width: 640px) {
  .auth-branding {
    min-height: 220px;
  }
}
</style>
