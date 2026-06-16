import { useEffect, useMemo, useState } from 'react'
import { createPortal } from 'react-dom'
import { loadPromptLibraryChunk } from '../data/promptLibrary/loaders'
import type { PromptLibraryLocale, PromptPresetDetail, PromptPresetSource, PromptPresetSummary } from '../data/promptLibrary/types'
import { getHostLocale, hostText } from '../lib/sub2apiHost'
import { CloseIcon, CopyIcon, LinkIcon } from './icons'

interface PromptLibraryModalProps {
  open: boolean
  onClose: () => void
  onUse: (prompt: string, mode: 'replace' | 'append') => void
}

const ALL_CATEGORY = '__all__'
const MAX_VISIBLE_PRESETS = 250

type PromptLibraryModule = typeof import('../data/promptLibrary')

function normalizeText(value: string) {
  return value.toLowerCase().replace(/\s+/g, ' ').trim()
}

function getPresetSearchText(preset: PromptPresetSummary) {
  return normalizeText([
    preset.title,
    preset.category,
    preset.description,
    preset.promptPreview,
    preset.sourceLabel,
    preset.sourceLanguage,
    ...preset.tags,
  ].filter(Boolean).join(' '))
}

function getSourceName(sourceId: string, sources: PromptPresetSource[]) {
  return sources.find((source) => source.id === sourceId)?.name ?? sourceId
}

function PreviewText({ text }: { text: string }) {
  const compact = text.replace(/\s+/g, ' ').trim()
  return (
    <p className="line-clamp-3 text-xs leading-relaxed text-gray-500 dark:text-gray-400">
      {compact}
    </p>
  )
}

export default function PromptLibraryModal({ open, onClose, onUse }: PromptLibraryModalProps) {
  const [query, setQuery] = useState('')
  const [locale, setLocale] = useState<PromptLibraryLocale>(getHostLocale())
  const [category, setCategory] = useState(ALL_CATEGORY)
  const [activeId, setActiveId] = useState<string | null>(null)
  const [presets, setPresets] = useState<PromptPresetSummary[]>([])
  const [sources, setSources] = useState<PromptPresetSource[]>([])
  const [detailsById, setDetailsById] = useState<Record<string, PromptPresetDetail>>({})
  const [loadingChunkId, setLoadingChunkId] = useState<string | null>(null)
  const [loadError, setLoadError] = useState<string | null>(null)

  useEffect(() => {
    if (!open) return
    setLocale(getHostLocale())
    setCategory(ALL_CATEGORY)
    setActiveId(null)
  }, [open])

  useEffect(() => {
    if (!open || presets.length > 0 || loadError) return
    let cancelled = false
    import('../data/promptLibrary')
      .then((module: PromptLibraryModule) => {
        if (cancelled) return
        setPresets(module.PROMPT_LIBRARY_INDEX)
        setSources(module.PROMPT_LIBRARY_SOURCES)
      })
      .catch((error) => {
        if (cancelled) return
        setLoadError(error instanceof Error ? error.message : String(error))
      })
    return () => {
      cancelled = true
    }
  }, [loadError, open, presets.length])

  const localeFiltered = useMemo(() => {
    return presets.filter((preset) => locale === 'all' || preset.locale === locale)
  }, [locale, presets])

  const categories = useMemo(() => {
    const counts = new Map<string, number>()
    localeFiltered.forEach((preset) => counts.set(preset.category, (counts.get(preset.category) ?? 0) + 1))
    return Array.from(counts.entries()).sort((a, b) => b[1] - a[1] || a[0].localeCompare(b[0]))
  }, [localeFiltered])

  const filtered = useMemo(() => {
    const q = normalizeText(query)
    return localeFiltered
      .filter((preset) => category === ALL_CATEGORY || preset.category === category)
      .filter((preset) => !q || getPresetSearchText(preset).includes(q))
  }, [category, localeFiltered, query])

  const displayedPresets = useMemo(() => filtered.slice(0, MAX_VISIBLE_PRESETS), [filtered])

  const activePreset = useMemo(() => {
    if (activeId) {
      const found = filtered.find((preset) => preset.id === activeId)
      if (found) return found
    }
    return null
  }, [activeId, filtered])

  const activeDetail = activePreset ? detailsById[activePreset.id] : null
  const activeSource = activePreset ? sources.find((source) => source.id === activePreset.sourceId) : null

  useEffect(() => {
    if (!open || !activePreset || detailsById[activePreset.id]) return
    let cancelled = false
    setLoadingChunkId(activePreset.chunkId)
    loadPromptLibraryChunk(activePreset.chunkId)
      .then((chunk) => {
        if (cancelled) return
        setDetailsById((current) => {
          const next = { ...current }
          chunk.forEach((detail) => {
            next[detail.id] = detail
          })
          return next
        })
      })
      .catch((error) => {
        if (cancelled) return
        setLoadError(error instanceof Error ? error.message : String(error))
      })
      .finally(() => {
        if (!cancelled) setLoadingChunkId(null)
      })
    return () => {
      cancelled = true
    }
  }, [activePreset, detailsById, open])

  if (!open) return null

  const usePrompt = (mode: 'replace' | 'append') => {
    if (!activeDetail?.prompt) return
    onUse(activeDetail.prompt, mode)
    onClose()
  }

  const isPromptLoading = Boolean(activePreset && loadingChunkId === activePreset.chunkId && !activeDetail)

  return createPortal(
    <div data-no-drag-select className="fixed inset-0 z-[80] flex items-end justify-center p-0 sm:items-center sm:p-4">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm animate-overlay-in" onClick={onClose} />
      <div className="relative z-10 flex h-[88vh] w-full flex-col overflow-hidden rounded-t-3xl border border-white/50 bg-white/95 shadow-2xl ring-1 ring-black/5 animate-modal-in dark:border-white/[0.08] dark:bg-gray-950/95 dark:ring-white/10 sm:h-[78vh] sm:max-w-6xl sm:rounded-3xl">
        <div className="flex items-center justify-between gap-4 border-b border-gray-200/70 px-4 py-3 dark:border-white/[0.08] sm:px-5">
          <div className="min-w-0">
            <h3 className="truncate text-base font-bold text-gray-800 dark:text-gray-100">
              {hostText('提示词参考', 'Prompt Library')}
            </h3>
            <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
              {hostText('内置多源提示词库；完整正文按需加载，插入后可继续修改。', 'Built-in multi-source prompt library; full prompt text loads on demand.')}
            </p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="rounded-full p-2 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-white/[0.08] dark:hover:text-gray-200"
            aria-label={hostText('关闭', 'Close')}
          >
            <CloseIcon className="h-5 w-5" />
          </button>
        </div>

        <div className="flex min-h-0 flex-1 flex-col sm:flex-row">
          <aside className="border-b border-gray-200/70 p-4 dark:border-white/[0.08] sm:w-72 sm:border-b-0 sm:border-r">
            <input
              value={query}
              onChange={(event) => {
                setQuery(event.target.value)
                setActiveId(null)
              }}
              placeholder={hostText('搜索标题、分类、作者或关键词', 'Search title, category, author, or keyword')}
              className="w-full rounded-xl border border-gray-200/70 bg-white/70 px-3 py-2.5 text-sm text-gray-700 outline-none transition focus:border-blue-300 focus:ring-1 focus:ring-blue-300/30 dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-gray-200 dark:focus:border-blue-500/50"
            />

            <div className="mt-3 grid grid-cols-3 gap-1 rounded-xl bg-gray-100/80 p-1 dark:bg-white/[0.05]">
              {([
                ['zh', hostText('中文', 'ZH')],
                ['en', hostText('英文', 'EN')],
                ['all', hostText('全部', 'All')],
              ] as const).map(([value, label]) => (
                <button
                  key={value}
                  type="button"
                  onClick={() => {
                    setLocale(value)
                    setCategory(ALL_CATEGORY)
                    setActiveId(null)
                  }}
                  className={`rounded-lg px-2 py-1.5 text-xs font-medium transition-colors ${
                    locale === value
                      ? 'bg-white text-blue-600 shadow-sm dark:bg-white/[0.12] dark:text-blue-300'
                      : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>

            <div className="mt-4 max-h-32 overflow-y-auto custom-scrollbar sm:max-h-[calc(78vh-238px)]">
              <button
                type="button"
                onClick={() => {
                  setCategory(ALL_CATEGORY)
                  setActiveId(null)
                }}
                className={`mb-1 flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-xs transition-colors ${
                  category === ALL_CATEGORY
                    ? 'bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300'
                    : 'text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-white/[0.06]'
                }`}
              >
                <span>{hostText('全部分类', 'All categories')}</span>
                <span>{localeFiltered.length}</span>
              </button>
              {categories.map(([name, count]) => (
                <button
                  key={name}
                  type="button"
                  onClick={() => {
                    setCategory(name)
                    setActiveId(null)
                  }}
                  className={`mb-1 flex w-full items-center justify-between gap-2 rounded-xl px-3 py-2 text-left text-xs transition-colors ${
                    category === name
                      ? 'bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-300'
                      : 'text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-white/[0.06]'
                  }`}
                >
                  <span className="truncate">{name}</span>
                  <span className="shrink-0">{count}</span>
                </button>
              ))}
            </div>

            <div className="mt-3 rounded-xl bg-gray-50/80 px-3 py-2 text-[11px] leading-relaxed text-gray-400 dark:bg-white/[0.03] dark:text-gray-500">
              {hostText('来源：YouMind 公开 README、davidwuw、freestylefly。', 'Sources: YouMind public README, davidwuw, and freestylefly.')}
            </div>
          </aside>

          <div className="grid min-h-0 flex-1 grid-rows-[minmax(0,1fr)_auto] sm:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)] sm:grid-rows-1">
            <div className="min-h-0 overflow-y-auto border-b border-gray-200/70 p-3 custom-scrollbar dark:border-white/[0.08] sm:border-b-0 sm:border-r">
              <div className="mb-2 px-1 text-xs text-gray-500 dark:text-gray-400">
                {loadError
                  ? hostText('提示词加载失败', 'Failed to load prompts')
                  : presets.length === 0
                  ? hostText('正在加载提示词索引...', 'Loading prompt index...')
                  : `${hostText('共', 'Showing')} ${filtered.length} ${hostText('条参考', 'prompts')}${filtered.length > displayedPresets.length ? hostText(`，显示前 ${displayedPresets.length} 条`, `, first ${displayedPresets.length} shown`) : ''}`}
              </div>
              <div className="space-y-2">
                {presets.length === 0 && !loadError && (
                  <div className="rounded-2xl border border-gray-200/70 bg-white/60 p-8 text-center text-sm text-gray-400 dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-gray-500">
                    {hostText('正在加载提示词参考库', 'Loading prompt library')}
                  </div>
                )}
                {loadError && (
                  <div className="rounded-2xl border border-red-200/70 bg-red-50/70 p-8 text-center text-sm text-red-500 dark:border-red-500/20 dark:bg-red-500/10 dark:text-red-300">
                    {loadError}
                  </div>
                )}
                {presets.length > 0 && displayedPresets.map((preset) => (
                  <button
                    key={preset.id}
                    type="button"
                    onClick={() => setActiveId(preset.id)}
                    className={`block w-full rounded-2xl border px-3 py-3 text-left transition-colors ${
                      activePreset?.id === preset.id
                        ? 'border-blue-200 bg-blue-50/80 dark:border-blue-500/30 dark:bg-blue-500/10'
                        : 'border-gray-200/70 bg-white/60 hover:bg-gray-50 dark:border-white/[0.08] dark:bg-white/[0.03] dark:hover:bg-white/[0.06]'
                    }`}
                  >
                    <div className="mb-1.5 flex items-center gap-2">
                      <span className="min-w-0 flex-1 truncate text-sm font-semibold text-gray-800 dark:text-gray-100">{preset.title}</span>
                      <span className="shrink-0 rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] text-gray-500 dark:bg-white/[0.08] dark:text-gray-400">{preset.category}</span>
                    </div>
                    <PreviewText text={preset.description || preset.promptPreview} />
                    <div className="mt-2 flex flex-wrap gap-1.5">
                      <span className="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] text-gray-500 dark:bg-white/[0.08] dark:text-gray-400">{preset.locale.toUpperCase()}</span>
                      <span className="rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] text-gray-500 dark:bg-white/[0.08] dark:text-gray-400">{getSourceName(preset.sourceId, sources)}</span>
                      {preset.featured && <span className="rounded-md bg-amber-100 px-1.5 py-0.5 text-[10px] text-amber-700 dark:bg-amber-500/10 dark:text-amber-300">Featured</span>}
                      {preset.needsReference && <span className="rounded-md bg-purple-100 px-1.5 py-0.5 text-[10px] text-purple-700 dark:bg-purple-500/10 dark:text-purple-300">{hostText('需参考图', 'Reference')}</span>}
                    </div>
                  </button>
                ))}
                {presets.length > 0 && filtered.length > displayedPresets.length && (
                  <div className="rounded-2xl border border-dashed border-gray-200 p-4 text-center text-xs text-gray-400 dark:border-white/[0.08] dark:text-gray-500">
                    {hostText('结果较多，请继续搜索或选择分类缩小范围。', 'Many results. Search or choose a category to narrow the list.')}
                  </div>
                )}
                {presets.length > 0 && filtered.length === 0 && (
                  <div className="rounded-2xl border border-dashed border-gray-200 p-8 text-center text-sm text-gray-400 dark:border-white/[0.08] dark:text-gray-500">
                    {hostText('没有匹配的提示词', 'No matching prompts')}
                  </div>
                )}
              </div>
            </div>

            <section className="flex min-h-0 flex-col">
              <div className="min-h-0 flex-1 overflow-y-auto p-4 custom-scrollbar sm:p-5">
                {activePreset ? (
                  <>
                    <div className="mb-3 flex flex-wrap items-center gap-2">
                      <span className="rounded-md bg-blue-50 px-2 py-1 text-xs font-medium text-blue-600 dark:bg-blue-500/10 dark:text-blue-300">{activePreset.category}</span>
                      <span className="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-500 dark:bg-white/[0.08] dark:text-gray-400">{activePreset.locale.toUpperCase()}</span>
                      {activePreset.sourceLanguage && <span className="rounded-md bg-gray-100 px-2 py-1 text-xs text-gray-500 dark:bg-white/[0.08] dark:text-gray-400">{activePreset.sourceLanguage}</span>}
                      {activePreset.needsReference && <span className="rounded-md bg-purple-100 px-2 py-1 text-xs font-medium text-purple-700 dark:bg-purple-500/10 dark:text-purple-300">{hostText('参考图', 'Reference image')}</span>}
                    </div>
                    <h4 className="text-base font-bold text-gray-800 dark:text-gray-100">{activePreset.title}</h4>
                    {activePreset.description && (
                      <p className="mt-2 text-sm leading-relaxed text-gray-500 dark:text-gray-400">{activePreset.description}</p>
                    )}
                    <pre className="mt-4 whitespace-pre-wrap break-words rounded-2xl border border-gray-200/70 bg-gray-50/80 p-4 text-xs leading-relaxed text-gray-700 dark:border-white/[0.08] dark:bg-black/20 dark:text-gray-200">
                      {activeDetail?.prompt ?? (isPromptLoading ? hostText('正在加载完整提示词...', 'Loading full prompt...') : activePreset.promptPreview)}
                    </pre>
                    {(activeDetail?.author || activeDetail?.publishedAt || activePreset.sourceLabel) && (
                      <div className="mt-3 text-[11px] leading-relaxed text-gray-400 dark:text-gray-500">
                        {[activeDetail?.author || activePreset.sourceLabel, activeDetail?.publishedAt].filter(Boolean).join(' · ')}
                      </div>
                    )}
                  </>
                ) : (
                  <div className="flex h-full items-center justify-center text-sm text-gray-400 dark:text-gray-500">
                    {hostText('选择一条提示词查看内容', 'Select a prompt to preview')}
                  </div>
                )}
              </div>

              <div className="border-t border-gray-200/70 p-4 dark:border-white/[0.08]">
                <div className="mb-3 flex flex-wrap items-center gap-x-2 gap-y-1 text-[11px] leading-relaxed text-gray-400 dark:text-gray-500">
                  <span>{hostText('来源：', 'Source: ')}</span>
                  <a className="font-medium text-gray-500 underline-offset-2 hover:underline dark:text-gray-400" href={activeSource?.url ?? 'https://github.com/YouMind-OpenLab/awesome-gpt-image-2'} target="_blank" rel="noopener noreferrer">
                    {activeSource?.name ?? hostText('多源提示词库', 'Multi-source prompt library')}
                  </a>
                  {activeSource && (
                    <>
                      <span>·</span>
                      <a className="font-medium text-gray-500 underline-offset-2 hover:underline dark:text-gray-400" href={activeSource.licenseUrl} target="_blank" rel="noopener noreferrer">
                        {activeSource.license}
                      </a>
                    </>
                  )}
                  {activeDetail?.sourceUrl && (
                    <>
                      <span>·</span>
                      <a className="inline-flex items-center gap-1 font-medium text-gray-500 underline-offset-2 hover:underline dark:text-gray-400" href={activeDetail.sourceUrl} target="_blank" rel="noopener noreferrer">
                        <LinkIcon className="h-3 w-3" />
                        {hostText('原链接', 'Original')}
                      </a>
                    </>
                  )}
                  <span>·</span>
                  <span>{hostText('已做格式整理', 'Adapted for in-app browsing')}</span>
                </div>
                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => usePrompt('append')}
                    disabled={!activeDetail?.prompt}
                    className="flex flex-1 items-center justify-center gap-2 rounded-xl border border-gray-200/70 bg-white/80 px-4 py-2.5 text-sm font-medium text-gray-700 transition-colors hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-gray-200 dark:hover:bg-white/[0.06]"
                  >
                    <CopyIcon className="h-4 w-4" />
                    {hostText('追加', 'Append')}
                  </button>
                  <button
                    type="button"
                    onClick={() => usePrompt('replace')}
                    disabled={!activeDetail?.prompt}
                    className="flex flex-1 items-center justify-center rounded-xl bg-blue-500 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-blue-600 disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    {isPromptLoading ? hostText('加载中', 'Loading') : hostText('使用', 'Use')}
                  </button>
                </div>
              </div>
            </section>
          </div>
        </div>
      </div>
    </div>,
    document.body,
  )
}
