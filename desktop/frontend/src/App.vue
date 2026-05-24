<script setup lang="ts">
import { ref, watch } from 'vue'
import { Sun, Moon } from 'lucide-vue-next'
import Sidebar from '@/components/layout/Sidebar.vue'
import StatusTab from '@/components/status/StatusTab.vue'
import AgentTab from '@/components/agent/AgentTab.vue'
import EnvTab from '@/components/env/EnvTab.vue'
import WebUITab from '@/components/webui/WebUITab.vue'
import ChannelTab from '@/components/channel/ChannelTab.vue'
import UpdateDialog from '@/components/update/UpdateDialog.vue'
import { useStatus } from '@/composables/useStatus'
import { useUpdater } from '@/composables/useUpdater'
import { useWailsEvents } from '@/composables/useWailsEvents'

import type { TabValue } from '@/types'

const activeTab = ref<TabValue>('status')
const { status, actionError, syncStatus } = useStatus()
useUpdater()

useWailsEvents(activeTab, actionError, syncStatus)

const switchToWeb = () => {
  activeTab.value = 'web'
}

// 主题管理
const theme = ref<'dark' | 'light'>(
  (typeof localStorage !== 'undefined' && (localStorage.getItem('ccx-desktop-theme') as 'dark' | 'light')) || 'dark'
)

const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  localStorage.setItem('ccx-desktop-theme', theme.value)

  // 同步主题到 iframe 内的 Web UI
  const iframe = document.querySelector('iframe[src*="localhost"], iframe[src*="127.0.0.1"]') as HTMLIFrameElement | null
  iframe?.contentWindow?.postMessage(
    { type: 'ccx-theme-change', theme: theme.value },
    '*'
  )
}

// 选项卡标题映射
const tabTitles: Record<TabValue, string> = {
  status: '网关状态监控',
  agent: 'Agent 代理配置',
  channels: '渠道中心',
  env: '环境参数管理',
  web: '内置控制台 Web UI'
}
</script>

<template>
  <div :class="[theme === 'light' ? 'light' : '', 'flex h-screen w-screen bg-background text-foreground overflow-hidden font-sans']">
    <!-- 常驻左侧高级磨砂侧边栏 -->
    <Sidebar v-model="activeTab" :theme="theme" />

    <!-- 右侧内容主展区 -->
    <main class="flex-1 flex flex-col min-w-0 h-full relative">
      <!-- 右侧顶部精细页眉 -->
      <header class="h-14 border-b border-border bg-background/60 backdrop-blur-md flex items-center justify-between px-8 shrink-0" data-wails-drag>
        <div class="flex items-center gap-3">
          <span class="text-xs bg-blue-500/10 text-blue-400 font-semibold px-2 py-0.5 rounded border border-blue-500/15">
            CCX CORE
          </span>
          <h2 class="text-sm font-bold text-foreground tracking-wide uppercase">
            {{ tabTitles[activeTab] }}
          </h2>
        </div>

        <div class="flex items-center gap-3">
          <!-- 主题切换按钮 -->
          <button
            @click="toggleTheme"
            class="p-1.5 rounded-lg border border-border text-muted-foreground hover:text-foreground hover:bg-background transition-colors"
            :title="theme === 'dark' ? '切换到亮色模式' : '切换到暗色模式'"
          >
            <Sun v-if="theme === 'dark'" class="w-4 h-4" />
            <Moon v-else class="w-4 h-4" />
          </button>

          <!-- 实时网关状态指示微标 -->
          <span
            v-if="status.running"
            class="text-[10px] bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 font-bold px-2 py-0.5 rounded-full"
          >
            GATEWAY ONLINE
          </span>
          <span
            v-else-if="status.starting"
            class="text-[10px] bg-amber-500/10 text-amber-400 border border-amber-500/20 font-bold px-2 py-0.5 rounded-full animate-pulse"
          >
            CONNECTING...
          </span>
          <span
            v-else
            class="text-[10px] bg-rose-500/10 text-rose-400 border border-rose-500/20 font-bold px-2 py-0.5 rounded-full"
          >
            GATEWAY OFFLINE
          </span>
        </div>
      </header>

      <!-- 独立内容滚动区域 -->
      <div class="flex-1 overflow-y-auto px-8 py-7">
        <div class="max-w-5xl mx-auto h-full">
          <!-- v-show 常驻缓存各 Tab，切换时保留内部状态与滚动位置 -->
          <div v-show="activeTab === 'status'" class="h-full">
            <StatusTab @switch-to-web="switchToWeb" />
          </div>
          <div v-show="activeTab === 'agent'" class="h-full">
            <AgentTab />
          </div>
          <div v-show="activeTab === 'channels'" class="h-full">
            <ChannelTab />
          </div>
          <div v-show="activeTab === 'env'" class="h-full">
            <EnvTab />
          </div>
          <div v-show="activeTab === 'web'" class="h-full">
            <WebUITab :status="status" :loading="false" />
          </div>
        </div>
      </div>
    </main>

    <UpdateDialog />
  </div>
</template>
