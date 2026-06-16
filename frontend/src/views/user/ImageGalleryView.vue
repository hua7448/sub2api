<template>
  <div class="flex min-h-[60vh] items-center justify-center text-sm text-gray-500 dark:text-gray-400">
    正在打开生图广场...
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  window.location.assign('/image-playground/index.html')
})
</script>
