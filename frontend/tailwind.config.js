/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // 主色调 - Anthropic-inspired clay
        primary: {
          50: '#fbf4ec',
          100: '#f4e4d6',
          200: '#e9c5ad',
          300: '#daa17f',
          400: '#c97b57',
          500: '#b85f3f',
          600: '#9d4b33',
          700: '#7f3c2d',
          800: '#663229',
          900: '#522b24',
          950: '#2f1712'
        },
        // 全局暖灰，覆盖默认冷灰以统一界面基调
        gray: {
          50: '#faf8f4',
          100: '#f1ede6',
          200: '#e4ddd2',
          300: '#d0c6b8',
          400: '#a99d8d',
          500: '#7d7164',
          600: '#5e554b',
          700: '#474038',
          800: '#302c27',
          900: '#1f1d1a',
          950: '#141210'
        },
        // 辅助色 - charcoal neutral
        accent: {
          50: '#faf8f4',
          100: '#f1ede6',
          200: '#e4ddd2',
          300: '#d0c6b8',
          400: '#a99d8d',
          500: '#7d7164',
          600: '#5e554b',
          700: '#474038',
          800: '#302c27',
          900: '#1f1d1a',
          950: '#141210'
        },
        // 深色模式背景 - warm charcoal
        dark: {
          50: '#faf8f4',
          100: '#f1ede6',
          200: '#e4ddd2',
          300: '#d0c6b8',
          400: '#a99d8d',
          500: '#7d7164',
          600: '#5e554b',
          700: '#474038',
          800: '#302c27',
          900: '#1f1d1a',
          950: '#141210'
        }
      },
      fontFamily: {
        sans: [
          'system-ui',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'PingFang SC',
          'Hiragino Sans GB',
          'Microsoft YaHei',
          'sans-serif'
        ],
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Monaco', 'Consolas', 'monospace']
      },
      boxShadow: {
        glass: '0 8px 32px rgba(0, 0, 0, 0.08)',
        'glass-sm': '0 4px 16px rgba(0, 0, 0, 0.06)',
        glow: '0 0 20px rgba(184, 95, 63, 0.22)',
        'glow-lg': '0 0 40px rgba(184, 95, 63, 0.32)',
        card: '0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.06)',
        'card-hover': '0 10px 40px rgba(0, 0, 0, 0.08)',
        'inner-glow': 'inset 0 1px 0 rgba(255, 255, 255, 0.1)'
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-primary': 'linear-gradient(135deg, #c97b57 0%, #9d4b33 100%)',
        'gradient-dark': 'linear-gradient(135deg, #302c27 0%, #141210 100%)',
        'gradient-glass':
          'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
        'mesh-gradient':
          'radial-gradient(at 35% 18%, rgba(184, 95, 63, 0.10) 0px, transparent 48%), radial-gradient(at 82% 4%, rgba(71, 64, 56, 0.08) 0px, transparent 50%), radial-gradient(at 0% 52%, rgba(201, 123, 87, 0.08) 0px, transparent 50%)'
      },
      animation: {
        'fade-in': 'fadeIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.3s ease-out',
        'slide-in-right': 'slideInRight 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        shimmer: 'shimmer 2s linear infinite',
        glow: 'glow 2s ease-in-out infinite alternate'
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
        glow: {
          '0%': { boxShadow: '0 0 20px rgba(184, 95, 63, 0.22)' },
          '100%': { boxShadow: '0 0 30px rgba(184, 95, 63, 0.36)' }
        }
      },
      backdropBlur: {
        xs: '2px'
      },
      borderRadius: {
        '4xl': '2rem'
      }
    }
  },
  plugins: []
}
