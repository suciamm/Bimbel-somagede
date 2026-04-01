<template>
  <div class="not-found panel">
    <p>404</p>
    <h1>Halaman tidak ditemukan</h1>
    <RouterLink class="btn btn-primary" to="/">Kembali ke dashboard</RouterLink>
  </div>
</template>

<style scoped>
.not-found {
  margin: 16px;
  min-height: calc(100vh - 32px);
  display: grid;
  place-content: center;
  gap: 10px;
  text-align: center;
}

p {
  margin: 0;
  color: var(--accent);
  font-weight: 800;
  font-size: 1.2rem;
}

h1 {
  margin: 0;
}
</style>
