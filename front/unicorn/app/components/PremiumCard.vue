<template>
  <div class="space-y-4">
    <!-- Premium Badge -->
    <div v-if="isPremium" class="flex items-center gap-2">
      <span class="inline-block px-2 py-1 bg-yellow-100 text-yellow-700 text-xs font-semibold rounded-full">
        ⭐ PREMIUM
      </span>
    </div>

    <!-- Card with Premium Styling -->
    <div
      class="border rounded-lg p-4 transition-all"
      :style="isPremium ? {
        borderColor: colorCode || '#FFD700',
        borderLeftWidth: '4px',
        background: 'linear-gradient(to right, rgba(255, 215, 0, 0.05), transparent)'
      } : {}"
    >
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  isPremium?: boolean
  colorCode?: string
}

withDefaults(defineProps<Props>(), {
  isPremium: false,
  colorCode: '#FFD700'
})
</script>
