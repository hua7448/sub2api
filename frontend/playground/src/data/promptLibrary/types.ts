export type PromptLibraryLocale = 'all' | 'zh' | 'en'

export interface PromptPresetSummary {
  id: string
  locale: 'zh' | 'en'
  title: string
  category: string
  description: string
  promptPreview: string
  sourceId: string
  sourceLanguage?: string
  sourceLabel?: string
  tags: string[]
  featured: boolean
  needsReference: boolean
  chunkId: string
}

export interface PromptPresetDetail {
  id: string
  prompt: string
  author?: string
  sourceUrl?: string
  publishedAt?: string
}

export interface PromptPresetSource {
  id: string
  name: string
  url: string
  license: string
  licenseUrl: string
  commit: string
}
