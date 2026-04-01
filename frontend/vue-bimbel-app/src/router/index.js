import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/RegisterView.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/',
    component: () => import('../layouts/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('../views/DashboardView.vue'),
        meta: { roles: ['admin', 'pembimbing', 'orangtua'] }
      },
      {
        path: 'murid',
        name: 'murid',
        component: () => import('../views/MuridView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'ortu',
        name: 'ortu',
        component: () => import('../views/OrtuView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'pembimbing',
        name: 'pembimbing',
        component: () => import('../views/PembimbingView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'jadwal',
        name: 'jadwal',
        component: () => import('../views/JadwalView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'paket',
        name: 'paket',
        component: () => import('../views/PaketView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'transaksi',
        name: 'transaksi',
        component: () => import('../views/TransaksiView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'absensi',
        name: 'absensi',
        component: () => import('../views/AbsensiView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'rekap-murid',
        name: 'rekap-murid',
        component: () => import('../views/RekapMuridBulananView.vue'),
        meta: { roles: ['admin'] }
      },
      {
        path: 'pembimbing/jadwal-hari-ini',
        name: 'pembimbing-jadwal-hari-ini',
        component: () => import('../views/PembimbingJadwalHariIniView.vue'),
        meta: { roles: ['pembimbing'] }
      },
      {
        path: 'pembimbing/jadwal-saya',
        name: 'pembimbing-jadwal-saya',
        component: () => import('../views/PembimbingJadwalSayaView.vue'),
        meta: { roles: ['pembimbing'] }
      },
      {
        path: 'pembimbing/murid-saya',
        name: 'pembimbing-murid-saya',
        component: () => import('../views/PembimbingMuridSayaView.vue'),
        meta: { roles: ['pembimbing'] }
      },
      {
        path: 'pembimbing/evaluasi',
        name: 'pembimbing-evaluasi',
        component: () => import('../views/PembimbingEvaluasiView.vue'),
        meta: { roles: ['pembimbing'] }
      },
      {
        path: 'orangtua/jadwal-bimbingan',
        name: 'orangtua-jadwal-bimbingan',
        component: () => import('../views/OrtuJadwalBimbinganView.vue'),
        meta: { roles: ['orangtua'] }
      },
      {
        path: 'orangtua/data-murid',
        name: 'orangtua-data-murid',
        component: () => import('../views/OrtuDataMuridView.vue'),
        meta: { roles: ['orangtua'] }
      },
      {
        path: 'orangtua/evaluasi',
        name: 'orangtua-evaluasi',
        component: () => import('../views/OrtuEvaluasiView.vue'),
        meta: { roles: ['orangtua'] }
      },
      {
        path: 'orangtua/paket-bimbingan',
        name: 'orangtua-paket-bimbingan',
        component: () => import('../views/OrtuPaketBimbinganView.vue'),
        meta: { roles: ['orangtua'] }
      },
      {
        path: 'orangtua/paket-bimbingan-saya',
        name: 'orangtua-paket-bimbingan-saya',
        component: () => import('../views/OrtuPaketBimbinganSayaView.vue'),
        meta: { roles: ['orangtua'] }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../views/NotFoundView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return { name: 'dashboard' }
  }

  const requiredRoles = to.meta?.roles
  const currentRole = authStore.user?.role
  if (requiredRoles && requiredRoles.length > 0 && !requiredRoles.includes(currentRole)) {
    return { name: 'dashboard' }
  }

  return true
})

export default router
