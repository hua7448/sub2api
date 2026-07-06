import {
  Terminal,
  Key,
  Wallet,
  Gift,
  HelpCircle,
  MessageCircle,
  ShieldCheck,
  Code2,
  type LucideIcon,
} from 'lucide-react'

export type NavItem = {
  id: string
  label: string
  icon: LucideIcon
}

export const NAV_ITEMS: NavItem[] = [
  { id: 'overview', label: '概览', icon: Terminal },
  { id: 'quick-start', label: '快速开始', icon: Code2 },
  { id: 'pricing', label: '充值套餐', icon: Wallet },
  { id: 'billing', label: '计费说明', icon: ShieldCheck },
  { id: 'api-key', label: '使用 API 密钥', icon: Key },
  { id: 'codex-config', label: 'Codex 配置', icon: Terminal },
  { id: 'claude-config', label: 'Claude Code / GLM 配置', icon: Terminal },
  { id: 'kimi-config', label: 'Kimi 配置', icon: Code2 },
  { id: 'invite', label: '邀请奖励', icon: Gift },
  { id: 'faq', label: '常见问题', icon: HelpCircle },
  { id: 'contact', label: '联系支持', icon: MessageCircle },
]

export const SITE = {
  home: 'https://smartapi.mynatapp.cc/',
  baseUrl: 'https://smartapi.mynatapp.cc',
  clients: ['Claude Code', 'Codex', 'Kimi CLI', 'GLM（via Claude Code）'],
  models: ['Claude', 'Codex', 'Kimi', 'GLM'],
  qqGroup: '665886623',
  hours: '工作日 9:00-24:00',
}

export type PricingTier = {
  usd: number
  cny: string
  discount: string
  highlight?: boolean
}

export const PRICING: PricingTier[] = [
  { usd: 20, cny: '19.9 元', discount: '不打折' },
  { usd: 90, cny: '85 元', discount: '9.4 折' },
  { usd: 150, cny: '135 元', discount: '9 折' },
  { usd: 200, cny: '176 元', discount: '8.8 折' },
  { usd: 500, cny: '425 元', discount: '8.5 折' },
  { usd: 1000, cny: '800 元', discount: '8 折', highlight: true },
]

export type FaqItem = {
  q: string
  a: string
}

export const FAQ: FaqItem[] = [
  {
    q: 'SmartQ 支持哪些工具？',
    a: '支持 Claude Code、Codex、Kimi CLI 与 GLM。',
  },
  {
    q: '为什么创建了 API Key 但看不到使用配置？',
    a: '请检查 API Key 是否已选择分组。没有分组的密钥无法生成对应客户端配置。',
  },
  {
    q: 'Codex 配置后仍无法使用怎么办？',
    a: '检查 config.toml、auth.json、base_url、OPENAI_API_KEY、密钥状态、分组是否正确，修改配置后重新打开终端或重启 Codex。',
  },
  {
    q: 'Claude Code 提示认证失败怎么办？',
    a: '检查 ANTHROPIC_BASE_URL、ANTHROPIC_AUTH_TOKEN、API Key 状态、Claude 分组，以及环境变量是否在同一个终端窗口中设置。',
  },
  {
    q: '为什么 Claude 扣费比预期高？',
    a: 'Claude 按 1.8X 倍率扣除额度。长上下文、反复读取文件、大量工具调用、长时间连续任务都会增加消耗。',
  },
  {
    q: '充值后没有到账怎么办？',
    a: '充值通过站内支付页面完成，如遇到账延迟请联系 QQ 群：665886623。',
  },
  {
    q: 'API Key 泄露了怎么办？',
    a: '立即进入"API 密钥"页面删除泄露的密钥，然后重新创建新密钥。',
  },
]
