<template>
  <section class="user-layout">
    <aside class="sidebar">
      <slot name="sidebar" />
    </aside>

    <main class="content">
      <slot />
    </main>

    <aside class="actions">
      <slot name="actions" />
    </aside>
  </section>
</template>

<style scoped>
.user-layout {
  display: grid;
  grid-template-columns: 280px minmax(0, 1fr) 220px;
  gap: 24px;
  align-items: start;
}

@media (max-width: 1280px) {
  .user-layout {
    grid-template-columns: 260px minmax(0, 1fr);
  }
  .actions {
    display: none;
  }
}

@media (max-width: 768px) {
  .user-layout {
    grid-template-columns: 1fr;
  }
}
</style>
