<script setup>
const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Form'
  }
})

const emit = defineEmits(['close'])

function onBackdropClick(event) {
  if (event.target === event.currentTarget) {
    emit('close')
  }
}
</script>

<template>
  <div v-if="props.show" class="modal-overlay" @click="onBackdropClick">
    <div class="modal-card panel" role="dialog" aria-modal="true" :aria-label="props.title">
      <header class="modal-header">
        <h2>{{ props.title }}</h2>
        <button class="btn btn-danger close-btn" @click="emit('close')">Tutup</button>
      </header>

      <div class="modal-body">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(4px);
  display: grid;
  place-items: center;
  padding: 16px;
}

.modal-card {
  width: min(840px, 100%);
  max-height: 90vh;
  overflow: auto;
  padding: 16px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.1rem;
}

.close-btn {
  padding: 8px 12px;
}
</style>
