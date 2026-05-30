<script setup lang="ts">
import { computed } from 'vue'
import { useStatus } from '@/composables/useStatus'
import { useUpdater } from '@/composables/useUpdater'
import Logo from '@/components/layout/Logo.vue'
import {
  Activity,
  Settings,
  Sliders,
  Globe,
  Play,
  Square,
  Power,
  RefreshCw,
  Network
} from 'lucide-vue-next'
import type { TabValue } from '@/types'

const modelValue = defineModel<TabValue>({ required: true })
const props = defineProps<{ theme: 'dark' | 'light' }>()

const { status, loading, autostartEnabled, startService, stopService, setAutostart } = useStatus()
const { state: updaterState, check: checkUpdate } = useUpdater()

const menuItems = [
  { id: 'status', label: '网关监控', icon: Activity, desc: '实时状态及核心日志' },
  { id: 'agent', label: 'Agent 配置', icon: Settings, desc: '本地开发代理控制' },
  { id: 'channels', label: '渠道中心', icon: Network, desc: '一键添加上游渠道' },
  { id: 'env', label: '环境参数', icon: Sliders, desc: '网关配置文件编辑' },
  { id: 'web', label: '管理界面', icon: Globe, desc: 'ai-trun Web 控制面板' }
] as const

const statusLabel = computed(() => {
  if (status.value.running) return '运行正常'
  if (status.value.starting) return '网关启动中'
  return '服务已断开'
})

const statusGlowClass = computed(() => {
  if (status.value.running) return 'shadow-glow-green bg-emerald-500'
  if (status.value.starting) return 'shadow-glow-orange bg-amber-500'
  return 'shadow-glow-red bg-rose-500 animate-pulse'
})

const handleDaemonAction = async () => {
  if (loading.value) return
  if (status.value.running) {
    await stopService()
  } else {
    await startService()
  }
}
</script>

<template>
  <aside class="w-68 flex flex-col h-full bg-card/40 border-r border-border backdrop-blur-3xl shrink-0 select-none">
    <!-- macOS 交通灯避让区 & 标题栏拖拽区域 -->
    <div class="h-14 w-full flex items-center justify-end px-5 shrink-0" data-wails-drag>
      <div class="flex items-center gap-2.5 mt-2.5">
        <Logo :size="32" />
        <span class="text-sm font-bold tracking-wider text-foreground/80">
          AI-TRUN
        </span>
      </div>
    </div>

    <!-- 导航菜单 -->
    <nav class="flex-1 px-3 py-6 space-y-1.5 overflow-y-auto">
      <button
        v-for="item in menuItems"
        :key="item.id"
        @click="modelValue = item.id"
        :class="[
          'w-full flex items-center gap-3.5 px-4 py-3 rounded-lg text-left transition-all duration-300 relative group overflow-hidden',
          modelValue === item.id
            ? 'bg-blue-600/10 text-blue-400 border border-blue-500/15'
            : 'text-muted-foreground hover:text-foreground hover:bg-background/60 border border-transparent'
        ]"
      >
        <!-- 侧栏滑块小霓虹指示器 -->
        <div
          v-if="modelValue === item.id"
          class="absolute left-0 top-3 bottom-3 w-1 rounded-r-full bg-blue-500 shadow-[0_0_10px_rgba(59,130,246,0.6)]"
        />

        <component
          :is="item.icon"
          :class="[
            'w-4.5 h-4.5 shrink-0 transition-transform duration-300 group-hover:scale-110',
            modelValue === item.id ? 'text-blue-400' : 'text-muted-foreground/60 group-hover:text-muted-foreground'
          ]"
        />
        <div class="flex flex-col min-w-0">
          <span class="text-sm font-medium leading-tight">{{ item.label }}</span>
          <span class="text-[10px] text-muted-foreground/60 mt-0.5 truncate group-hover:text-muted-foreground transition-colors">
            {{ item.desc }}
          </span>
        </div>
      </button>
    </nav>

    <!-- 底部常驻迷你服务守护面板 -->
    <div class="p-4 border-t border-border bg-background/30 shrink-0">
      <div class="p-3.5 rounded-xl border border-border bg-background/40 hover:bg-background/60 transition-colors">
        <div class="flex items-center justify-between mb-3.5">
          <div class="flex items-center gap-2 min-w-0">
            <!-- 霓虹呼吸指示灯 -->
            <div :class="['w-2 h-2 rounded-full transition-all duration-500 shrink-0', statusGlowClass]" />
            <span class="text-xs font-semibold text-foreground/80 truncate">{{ statusLabel }}</span>
          </div>
          <!-- 迷你开关控制，可快速启停 -->
          <button
            @click="handleDaemonAction"
            :disabled="loading"
            :class="[
              'p-1.5 rounded-lg border text-xs transition-all duration-200 hover:scale-105 active:scale-95 shrink-0 cursor-pointer',
              status.running
                ? 'bg-rose-500/10 text-rose-400 border-rose-500/20 hover:bg-rose-500/20'
                : 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20 hover:bg-emerald-500/20'
            ]"
          >
            <component :is="status.running ? Square : Play" class="w-3 h-3" />
          </button>
        </div>

        <!-- 详细物理信息 -->
        <div class="space-y-1.5 text-[10px] font-mono text-muted-foreground">
          <div class="flex justify-between items-center">
            <span>网关端口</span>
            <span class="text-foreground/70 bg-background px-1.5 py-0.5 rounded border border-border">
              {{ status.port || '——' }}
            </span>
          </div>
          <div class="flex justify-between items-center" v-if="status.pid">
            <span>守护 PID</span>
            <span class="text-foreground/70 bg-background px-1.5 py-0.5 rounded border border-border">
              {{ status.pid }}
            </span>
          </div>
          <div class="flex justify-between items-center">
            <span>开机自启</span>
            <button
              @click="setAutostart(!autostartEnabled)"
              :class="[
                'flex items-center gap-1 px-1.5 py-0.5 rounded border transition-all duration-200 cursor-pointer',
                autostartEnabled
                  ? 'bg-blue-500/10 text-blue-400 border-blue-500/20'
                  : 'bg-background text-muted-foreground border-border hover:text-foreground'
              ]"
            >
              <Power class="w-2.5 h-2.5" />
              <span>{{ autostartEnabled ? '已开启' : '已关闭' }}</span>
            </button>
          </div>
          <div class="flex justify-between items-center">
            <span>当前版本</span>
            <button
              @click="checkUpdate()"
              :disabled="updaterState.checking"
              :class="[
                'flex items-center gap-1 px-1.5 py-0.5 rounded border transition-all duration-200 cursor-pointer',
                'bg-background text-foreground/70 border-border hover:text-blue-400 hover:border-blue-500/20',
                updaterState.checking && 'opacity-60 cursor-wait'
              ]"
              :title="updaterState.checking ? '检查中…' : '点击检查更新'"
            >
              <RefreshCw class="w-2.5 h-2.5" :class="updaterState.checking && 'animate-spin'" />
              <span>v{{ updaterState.version?.version || '—' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>
