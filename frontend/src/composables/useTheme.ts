import { useTheme as useVuetifyTheme } from 'vuetify'

// Precision Glass 主题配置
export const GLASS_THEME = {
  font: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
  monoFont: "'JetBrains Mono', 'SF Mono', 'Fira Code', monospace",
  radius: {
    sm: '8px',
    md: '12px',
    lg: '16px',
    full: '999px'
  }
}

export function useAppTheme() {
  const _vuetifyTheme = useVuetifyTheme()

  // 应用 Precision Glass 主题
  function applyGlassTheme() {
    document.documentElement.style.setProperty('--app-font', GLASS_THEME.font)
    document.documentElement.style.setProperty('--app-font-mono', GLASS_THEME.monoFont)
    document.documentElement.style.setProperty('--app-radius-sm', GLASS_THEME.radius.sm)
    document.documentElement.style.setProperty('--app-radius-md', GLASS_THEME.radius.md)
    document.documentElement.style.setProperty('--app-radius-lg', GLASS_THEME.radius.lg)
    document.documentElement.style.setProperty('--app-radius-full', GLASS_THEME.radius.full)
  }

  // 初始化
  function init() {
    applyGlassTheme()
  }

  return {
    init,
    applyGlassTheme
  }
}
