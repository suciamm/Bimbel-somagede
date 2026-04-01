<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from '../components/AppSidebar.vue'
import AppHeader from '../components/AppHeader.vue'

const route = useRoute()
const isSidebarOpen = ref(false)

watch(
  () => route.fullPath,
  () => {
    isSidebarOpen.value = false
  }
)
</script>

<template>
  <div class="layout">
    <AppSidebar :is-open="isSidebarOpen" @close="isSidebarOpen = false" />

    <main class="main-area">
      <AppHeader @toggle-menu="isSidebarOpen = !isSidebarOpen" />
      <section class="content panel">
        <RouterView />
      </section>
    </main>
  </div>
</template>

<style scoped>
.layout {
  height: 100vh;
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
  padding: 16px;
  overflow: hidden;
}

.main-area {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
  min-height: 0;
}

.content {
  padding: 20px;
  height: calc(100vh - 132px);
  min-height: calc(100vh - 132px);
  overflow: auto;
  scrollbar-gutter: stable;
}

@media (max-width: 980px) {
  .layout {
    height: auto;
    min-height: 100vh;
    grid-template-columns: 1fr;
    overflow: visible;
  }

  .content {
    height: auto;
    min-height: auto;
    overflow: visible;
  }
}
</style>
