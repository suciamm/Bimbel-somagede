<script setup>
defineProps({
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  tone: {
    type: String,
    default: 'default'
  }
})
</script>

<template>
  <article :class="['stat-card', 'panel', tone]">
    <p>{{ label }}</p>
    <strong>{{ value }}</strong>
  </article>
</template>

<style scoped>
.stat-card {
  padding: 16px;
}

.stat-card p {
  margin: 0;
  color: var(--text-muted);
  font-size: 0.9rem;
}

.stat-card strong {
  display: block;
  margin-top: 8px;
  font-size: 1.8rem;
  letter-spacing: -0.02em;
}

.stat-card.ocean {
  background: linear-gradient(135deg, #2563eb, #3b82f6);
  color: #fff;
}

.stat-card.ocean p {
  color: #dbeafe;
}

.stat-card.violet {
  background: linear-gradient(135deg, #7c3aed, #8b5cf6);
  color: #fff;
}

.stat-card.violet p {
  color: #ede9fe;
}

.stat-card.sunset {
  background: linear-gradient(135deg, #f59e0b, #f97316);
  color: #fff;
}

.stat-card.sunset p {
  color: #ffedd5;
}

.stat-card.ruby {
  background: linear-gradient(135deg, #ef4444, #f43f5e);
  color: #fff;
}

.stat-card.ruby p {
  color: #ffe4e6;
}
</style>
