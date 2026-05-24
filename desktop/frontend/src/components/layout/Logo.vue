<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  size?: number | string
  animated?: boolean
}>(), {
  size: 24,
  animated: true
})

const sizeStyle = computed(() => {
  const s = typeof props.size === 'number' ? `${props.size}px` : props.size
  return { width: s, height: s }
})
</script>

<template>
  <div class="relative shrink-0 flex items-center justify-center select-none" :style="sizeStyle">
    <svg
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="w-full h-full"
    >
      <defs>
        <!-- 底板渐变 -->
        <linearGradient id="ccx-bg" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="#4F46E5" />
          <stop offset="100%" stop-color="#7C3AED" />
        </linearGradient>
        <!-- 高光 -->
        <linearGradient id="ccx-shine" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stop-color="rgba(255,255,255,0.25)" />
          <stop offset="50%" stop-color="rgba(255,255,255,0.02)" />
          <stop offset="100%" stop-color="rgba(0,0,0,0.08)" />
        </linearGradient>
        <!-- 图标线条渐变 -->
        <linearGradient id="ccx-line" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stop-color="rgba(255,255,255,0.95)" />
          <stop offset="100%" stop-color="rgba(255,255,255,0.7)" />
        </linearGradient>
        <!-- 阴影 -->
        <filter id="ccx-shadow" x="-20%" y="-10%" width="140%" height="130%">
          <feDropShadow dx="0" dy="3" stdDeviation="4" flood-color="rgba(0,0,0,0.2)" />
        </filter>
      </defs>

      <!-- 圆角方形底板 -->
      <rect x="6" y="6" width="88" height="88" rx="22"
        fill="url(#ccx-bg)" filter="url(#ccx-shadow)" />
      <!-- 高光覆盖 -->
      <rect x="6" y="6" width="88" height="88" rx="22"
        fill="url(#ccx-shine)" />

      <!-- 路由节点拓扑图 -->
      <g :class="{ 'ccx-pulse': animated }" transform-origin="50 50">
        <!-- 连接线 -->
        <line x1="50" y1="50" x2="28" y2="26" stroke="url(#ccx-line)" stroke-width="2.5" stroke-linecap="round" opacity="0.45" />
        <line x1="50" y1="50" x2="74" y2="30" stroke="url(#ccx-line)" stroke-width="2.5" stroke-linecap="round" opacity="0.4" />
        <line x1="50" y1="50" x2="26" y2="68" stroke="url(#ccx-line)" stroke-width="2.5" stroke-linecap="round" opacity="0.4" />
        <line x1="50" y1="50" x2="74" y2="72" stroke="url(#ccx-line)" stroke-width="2.5" stroke-linecap="round" opacity="0.45" />

        <!-- 外圈节点 -->
        <circle cx="28" cy="26" r="5" fill="url(#ccx-line)" opacity="0.75" />
        <circle cx="74" cy="30" r="5" fill="url(#ccx-line)" opacity="0.65" />
        <circle cx="26" cy="68" r="5" fill="url(#ccx-line)" opacity="0.65" />
        <circle cx="74" cy="72" r="5" fill="url(#ccx-line)" opacity="0.75" />

        <!-- 中心节点 -->
        <circle cx="50" cy="50" r="8" fill="url(#ccx-line)" opacity="0.95" />

        <!-- 轨道环 -->
        <circle cx="50" cy="50" r="18"
          stroke="url(#ccx-line)" stroke-width="2"
          stroke-dasharray="5 3 2 4" fill="none"
          opacity="0.3"
          :class="{ 'ccx-orbit': animated }" />
      </g>
    </svg>
  </div>
</template>

<style scoped>
.ccx-pulse {
  animation: ccx-pulse 3s infinite ease-in-out;
}
@keyframes ccx-pulse {
  0%, 100% { transform: scale(0.97); opacity: 0.92; }
  50%      { transform: scale(1.02); opacity: 1; }
}
.ccx-orbit {
  transform-origin: 50px 50px;
  animation: ccx-spin 18s linear infinite;
}
@keyframes ccx-spin {
  to { transform: rotate(360deg); }
}
</style>
