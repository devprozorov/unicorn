<template>
  <div class="skeleton">
    <div class="line title"></div>
    <div class="line"></div>
    <div class="line"></div>
  </div>
</template>

<style scoped>
.skeleton {
  animation: pulse 1.5s infinite;
}

.line {
  height: 14px;
  background: #e5e7eb;
  border-radius: 8px;
  margin-bottom: 12px;
}

.title {
  height: 24px;
  width: 40%;
}

@keyframes pulse {
  0% { opacity: 1 }
  50% { opacity: .4 }
  100% { opacity: 1 }
}
</style>
