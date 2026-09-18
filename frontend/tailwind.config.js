/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  // 深/浅主题通过 <html> 上的 .dark-theme 类切换（见 src/stores/theme.ts），
  // 这样 dark: 变体与下方的 CSS 变量令牌都由同一个类驱动，而不再跟随操作系统。
  darkMode: ['class', '.dark-theme'],
  theme: {
    extend: {
      colors: {
        // 科技风背景层级 - 由 CSS 变量驱动，浅/深主题自动翻转（见 src/style.css）
        cyber: {
          bg: 'rgb(var(--cyber-bg) / <alpha-value>)',
          surface: 'rgb(var(--cyber-surface) / <alpha-value>)',
          elevated: 'rgb(var(--cyber-elevated) / <alpha-value>)',
          border: 'rgb(var(--cyber-border) / <alpha-value>)',
          muted: 'rgb(var(--cyber-muted) / <alpha-value>)'
        },
        // 蓝色系主色调（品牌色，浅/深主题通用）
        neon: {
          cyan: '#3b82f6',
          purple: '#6366f1',
          green: '#10b981',
          pink: '#ec4899',
          amber: '#f59e0b'
        },
        // 文字色阶 - 由 CSS 变量驱动，浅/深主题自动翻转
        txt: {
          primary: 'rgb(var(--txt-primary) / <alpha-value>)',
          secondary: 'rgb(var(--txt-secondary) / <alpha-value>)',
          dim: 'rgb(var(--txt-dim) / <alpha-value>)'
        },
        // 兼容旧代码中的 dark-* 引用（现在用于深色模式）
        dark: {
          300: '#d4deea',
          400: '#b8c6d8',
          500: '#93a5bb',
          600: '#70829a',
          700: '#566983',
          800: '#41526d',
          900: '#2f415d',
          950: '#223149'
        },
        // 兼容旧代码中的 primary-* 引用 - 蓝色系
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554'
        },
        // 兼容旧代码中的 accent-* 引用 - 靛蓝色
        accent: {
          50: '#eef2ff',
          100: '#e0e7ff',
          200: '#c7d2fe',
          300: '#a5b4fc',
          400: '#818cf8',
          500: '#6366f1',
          600: '#4f46e5',
          700: '#4338ca',
          800: '#3730a3',
          900: '#312e81',
          950: '#1e1b4b'
        }
      },
      fontFamily: {
        sans: [
          'Inter',
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: ['JetBrains Mono Variable', 'JetBrains Mono', 'Fira Code', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'monospace']
      },
      boxShadow: {
        'neon-cyan': '0 4px 20px rgba(59, 130, 246, 0.3)',
        'neon-cyan-lg': '0 8px 30px rgba(59, 130, 246, 0.4)',
        'neon-purple': '0 4px 14px 0 rgba(99, 102, 241, 0.25)',
        'neon-pink': '0 4px 14px 0 rgba(236, 72, 153, 0.25)',
        'neon-green': '0 4px 14px 0 rgba(16, 185, 129, 0.25)',
        'neon-amber': '0 4px 14px 0 rgba(245, 158, 11, 0.25)',
        glass: '0 4px 20px rgba(59, 130, 246, 0.08)',
        'glass-sm': '0 2px 10px rgba(59, 130, 246, 0.05)',
        card: '0 1px 3px rgba(0, 0, 0, 0.1)',
        'card-hover': '0 8px 30px rgba(59, 130, 246, 0.15)',
        glow: '0 0 20px rgba(59, 130, 246, 0.4)',
        'glow-lg': '0 0 40px rgba(59, 130, 246, 0.5)',
        'inner-glow': 'inset 0 1px 0 rgba(59, 130, 246, 0.1)'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-cyber': 'linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 50%, #f8fafc 100%)',
        'gradient-primary': 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
        'gradient-dark': 'linear-gradient(135deg, #1e293b 0%, #0f172a 100%)',
        'gradient-glass': 'linear-gradient(135deg, rgba(255,255,255,0.7) 0%, rgba(255,255,255,0.5) 100%)',
        'grid-pattern': 'linear-gradient(rgba(59,130,246,0.03) 1px, transparent 1px), linear-gradient(90deg, rgba(59,130,246,0.03) 1px, transparent 1px)',
        'scan-lines': 'repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(59,130,246,0.02) 2px, rgba(59,130,246,0.02) 4px)',
        'mesh-gradient': 'radial-gradient(at 20% 30%, rgba(59, 130, 246, 0.08) 0px, transparent 50%), radial-gradient(at 80% 20%, rgba(96, 165, 250, 0.06) 0px, transparent 50%), radial-gradient(at 40% 80%, rgba(59, 130, 246, 0.05) 0px, transparent 50%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glowPulse 2s ease-in-out infinite alternate',
        'neon-flicker': 'neonFlicker 3s ease-in-out infinite',
        'scan-line': 'scanLine 8s linear infinite'
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' }
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideDown: {
          '0%': { opacity: '0', transform: 'translateY(-10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' }
        },
        slideInRight: {
          '0%': { opacity: '0', transform: 'translateX(20px)' },
          '100%': { opacity: '1', transform: 'translateX(0)' }
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.95)' },
          '100%': { opacity: '1', transform: 'scale(1)' }
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' }
        },
        glowPulse: {
          '0%': { boxShadow: '0 0 20px rgba(59, 130, 246, 0.3)' },
          '100%': { boxShadow: '0 0 30px rgba(59, 130, 246, 0.5)' }
        },
        neonFlicker: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.8' }
        },
        scanLine: {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': '2rem'
      },
      backgroundSize: {
        'grid-64': '64px 64px'
      }
    }
  },
  plugins: []
}
