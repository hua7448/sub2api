import fs from 'node:fs'
import path from 'node:path'
import { execFileSync } from 'node:child_process'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const playgroundRoot = path.resolve(__dirname, '..')
const outputDir = path.join(playgroundRoot, 'src/data/promptLibrary')
const chunksDir = path.join(outputDir, 'chunks')

const sourceRoots = {
  youmind: process.env.YOUMIND_PROMPTS_DIR || '/tmp/awesome-gpt-image-2',
  davidwuw: process.env.DAVIDWUW_PROMPTS_DIR || '/tmp/awesome-gpt-image2-prompts',
  freestyle: process.env.FREESTYLE_PROMPTS_DIR || '/tmp/freestyle-awesome-gpt-image-2',
}

const sourceDefs = [
  {
    id: 'youmind',
    name: 'YouMind OpenLab / Awesome GPT Image 2 Prompts',
    url: 'https://github.com/YouMind-OpenLab/awesome-gpt-image-2',
    license: 'CC BY 4.0',
    licenseUrl: 'https://creativecommons.org/licenses/by/4.0/',
    commit: getGitCommit(sourceRoots.youmind),
  },
  {
    id: 'davidwuw',
    name: 'davidwuw0811-boop / awesome-gpt-image2-prompts',
    url: 'https://github.com/davidwuw0811-boop/awesome-gpt-image2-prompts',
    license: 'MIT',
    licenseUrl: 'https://github.com/davidwuw0811-boop/awesome-gpt-image2-prompts',
    commit: getGitCommit(sourceRoots.davidwuw),
  },
  {
    id: 'freestylefly',
    name: 'freestylefly / awesome-gpt-image-2',
    url: 'https://github.com/freestylefly/awesome-gpt-image-2',
    license: 'MIT',
    licenseUrl: 'https://github.com/freestylefly/awesome-gpt-image-2/blob/main/LICENSE',
    commit: getGitCommit(sourceRoots.freestyle),
  },
]

const freestyleCategoryMap = {
  'Architecture & Spaces': '建筑与空间',
  'Brand & Logos': '品牌与标志',
  'Characters & People': '角色与人物',
  'Charts & Infographics': '图表与信息可视化',
  'Documents & Publishing': '文档与出版',
  'History & Classical Themes': '历史与古典主题',
  'Illustration & Art': '插画与艺术',
  'Other Use Cases': '其他案例',
  'Photography & Realism': '摄影与写实',
  'Posters & Typography': '海报与排版',
  'Products & E-commerce': '商品与电商',
  'Scenes & Storytelling': '场景与叙事',
  'UI & Interfaces': 'UI 与界面',
}

const entries = [
  ...readYouMindEntries(sourceRoots.youmind),
  ...readFreestyleEntries(sourceRoots.freestyle),
  ...readFreestyleTemplates(sourceRoots.freestyle),
  ...readDavidWuwEntries(sourceRoots.davidwuw),
]

const safetyFiltered = entries.filter(isSuitableForPromptLibrary)
const deduped = dedupeEntries(safetyFiltered)
writePromptLibrary(deduped)

console.log(`Prompt library generated: ${deduped.length} entries, ${sourceDefs.length} sources, ${entries.length - safetyFiltered.length} filtered`)

function getGitCommit(root) {
  try {
    return execFileSync('git', ['-C', root, 'rev-parse', 'HEAD'], {
      encoding: 'utf8',
      stdio: ['ignore', 'pipe', 'ignore'],
    }).trim()
  } catch {
    // Fall through to direct .git/HEAD parsing for environments without git.
  }

  try {
    const head = fs.readFileSync(path.join(root, '.git/HEAD'), 'utf8').trim()
    if (head.startsWith('ref: ')) {
      const refPath = head.replace('ref: ', '')
      return fs.readFileSync(path.join(root, '.git', refPath), 'utf8').trim()
    }
    return head
  } catch {
    return 'unknown'
  }
}

function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, 'utf8'))
}

function readText(filePath) {
  return fs.readFileSync(filePath, 'utf8')
}

function readYouMindEntries(root) {
  const sources = [
    { file: 'README_zh.md', locale: 'zh', prefix: 'youmind-zh' },
    { file: 'README.md', locale: 'en', prefix: 'youmind-en' },
  ]
  return sources.flatMap(({ file, locale, prefix }) => {
    const markdown = readText(path.join(root, file))
    return parseYouMindMarkdown(markdown, locale).map((item, index) => ({
      id: `${prefix}-${item.galleryId || index + 1}`,
      locale,
      title: item.title,
      category: inferCategory(item.title, item.featured, locale),
      description: item.description || item.prompt.slice(0, 180),
      prompt: item.prompt,
      promptPreview: makePreview(item.prompt),
      sourceId: 'youmind',
      sourceLanguage: item.sourceLanguage,
      sourceUrl: item.tryUrl || item.sourceUrl || sourceDefs[0].url,
      sourceLabel: item.author || 'YouMind OpenLab',
      author: item.author,
      publishedAt: item.publishedAt,
      tags: item.featured ? ['Featured'] : [],
      featured: item.featured,
      needsReference: /reference|参考图|上传图片|provided image|uploaded image/i.test(item.prompt),
    }))
  })
}

function parseYouMindMarkdown(markdown, locale) {
  const lines = markdown.split(/\r?\n/)
  const items = []
  let featured = false
  let i = 0

  while (i < lines.length) {
    const line = lines[i]
    if (/^##\s+🔥/.test(line)) featured = true
    if (/^##\s+📋/.test(line)) featured = false

    const match = line.match(/^### No\. \d+: (.+)$/)
    if (!match) {
      i += 1
      continue
    }

    const title = match[1].trim()
    const section = []
    i += 1
    while (i < lines.length && !/^### No\. \d+: /.test(lines[i]) && !/^##\s+(🤝|📄|🙏|⭐|📚)/.test(lines[i])) {
      section.push(lines[i])
      i += 1
    }

    const description = extractMarkdownBlock(section, /^#### 📖/)
    const prompt = extractFencedAfterHeader(section, /^#### 📝/)
    if (!prompt) continue

    const text = section.join('\n')
    const galleryId = text.match(/gpt-image-2-prompts\?id=(\d+)/)?.[1]
    const tryUrl = text.match(/\*\*\[.*?\]\((https:\/\/youmind\.com[^)]+)\)\*\*/)?.[1]
    const sourceUrl = text.match(/-\s+\*\*(?:Source|来源):\*\*\s+\[[^\]]+\]\(([^)]+)\)/)?.[1]
    const author = text.match(/-\s+\*\*(?:Author|作者):\*\*\s+\[([^\]]+)\]\([^)]+\)/)?.[1]
    const publishedAt = text.match(/-\s+\*\*(?:Published|发布时间):\*\*\s+(.+)/)?.[1]?.trim()
    const sourceLanguage = text.match(/-\s+\*\*(?:Languages|多语言):\*\*\s+(.+)/)?.[1]?.trim()?.toUpperCase()

    items.push({ title, description, prompt, galleryId, tryUrl, sourceUrl, author, publishedAt, sourceLanguage, featured })
  }

  return items
}

function extractMarkdownBlock(lines, headerPattern) {
  const start = lines.findIndex((line) => headerPattern.test(line))
  if (start < 0) return ''
  const out = []
  for (let i = start + 1; i < lines.length; i += 1) {
    if (/^#### /.test(lines[i])) break
    const line = lines[i].trim()
    if (line) out.push(line)
  }
  return out.join('\n').trim()
}

function extractFencedAfterHeader(lines, headerPattern) {
  const start = lines.findIndex((line) => headerPattern.test(line))
  if (start < 0) return ''
  let fenceStart = -1
  for (let i = start + 1; i < lines.length; i += 1) {
    if (/^```/.test(lines[i])) {
      fenceStart = i
      break
    }
    if (/^#### /.test(lines[i])) return ''
  }
  if (fenceStart < 0) return ''
  const out = []
  for (let i = fenceStart + 1; i < lines.length; i += 1) {
    if (/^```/.test(lines[i])) break
    out.push(lines[i])
  }
  return out.join('\n').trim()
}

function readDavidWuwEntries(root) {
  const prompts = readJson(path.join(root, 'prompts.json'))
  return prompts.map((item) => ({
    id: `davidwuw-${item.id}`,
    locale: 'zh',
    title: item.title_cn || item.title_en || `Prompt ${item.id}`,
    category: item.category_cn || item.category || '提示词',
    description: item.note || item.title_en || item.category_cn || '',
    prompt: item.prompt,
    promptPreview: makePreview(item.prompt),
    sourceId: 'davidwuw',
    sourceLanguage: detectLanguage(item.prompt).toUpperCase(),
    sourceUrl: sourceDefs[1].url,
    sourceLabel: item.source || 'davidwuw prompt collection',
    author: item.author,
    tags: [item.category, item.source].filter(Boolean),
    needsReference: Boolean(item.needs_ref),
  }))
}

function readFreestyleEntries(root) {
  const data = readJson(path.join(root, 'data/cases.json'))
  return data.cases.map((item) => ({
    id: `freestylefly-${item.id}`,
    locale: detectLanguage(`${item.title}\n${item.prompt}`),
    title: item.title,
    category: freestyleCategoryMap[item.category] || item.category || '案例',
    description: item.promptPreview || item.imageAlt || '',
    prompt: item.prompt,
    promptPreview: item.promptPreview || makePreview(item.prompt),
    sourceId: 'freestylefly',
    sourceLanguage: detectLanguage(item.prompt).toUpperCase(),
    sourceUrl: item.githubUrl || item.sourceUrl || sourceDefs[2].url,
    sourceLabel: item.sourceLabel || 'freestylefly case',
    author: item.sourceLabel,
    tags: [item.category, ...(item.styles || []), ...(item.scenes || [])].filter(Boolean),
    featured: Boolean(item.featured),
    needsReference: /reference|provided|uploaded|参考图|改图|using the provided/i.test(item.prompt),
  }))
}

function readFreestyleTemplates(root) {
  const markdown = readText(path.join(root, 'docs/templates.md'))
  const lines = markdown.split(/\r?\n/)
  const items = []
  let category = '工业级模板'
  let templateTitle = ''
  let counter = 1

  for (let i = 0; i < lines.length; i += 1) {
    const categoryMatch = lines[i].match(/^###\s+(.+)$/)
    if (categoryMatch) {
      category = categoryMatch[1].trim()
      templateTitle = ''
      continue
    }

    const titleMatch = lines[i].match(/^\*\*(.+?)\*\*$/)
    if (titleMatch) {
      templateTitle = titleMatch[1].trim()
      continue
    }

    const fenceMatch = lines[i].match(/^```(\w+)?/)
    if (!fenceMatch) continue

    const block = []
    i += 1
    while (i < lines.length && !/^```/.test(lines[i])) {
      block.push(lines[i])
      i += 1
    }
    const prompt = block.join('\n').trim()
    if (!prompt || prompt.length < 30) continue

    const title = templateTitle || `${category}模板`
    items.push({
      id: `freestylefly-template-${counter}`,
      locale: detectLanguage(`${title}\n${prompt}`),
      title: `${category} - ${title}`,
      category: '工业级模板',
      description: `来自 freestylefly 工业级提示词模板与防坑指南的 ${category} 模板。`,
      prompt,
      promptPreview: makePreview(prompt),
      sourceId: 'freestylefly',
      sourceLanguage: detectLanguage(prompt).toUpperCase(),
      sourceUrl: `${sourceDefs[2].url}/blob/main/docs/templates.md`,
      sourceLabel: 'freestylefly templates',
      author: 'freestylefly',
      tags: ['template', category, fenceMatch[1]].filter(Boolean),
      featured: false,
      needsReference: /reference|参考图|上传图片|provided image|uploaded image/i.test(prompt),
    })
    counter += 1
  }

  return items
}

function inferCategory(title, featured, locale) {
  if (featured) return locale === 'zh' ? '精选提示词' : 'Featured Prompts'
  const separator = title.includes(' - ') ? ' - ' : title.includes('：') ? '：' : null
  if (!separator) return locale === 'zh' ? '综合提示词' : 'General Prompts'
  return title.split(separator)[0].trim()
}

function detectLanguage(value) {
  const chinese = (value.match(/[\u4e00-\u9fff]/g) || []).length
  const latin = (value.match(/[A-Za-z]/g) || []).length
  return chinese > latin * 0.08 ? 'zh' : 'en'
}

function makePreview(value) {
  return value.replace(/\s+/g, ' ').trim().slice(0, 160)
}

function makeDescription(value) {
  return (value || '').replace(/\s+/g, ' ').trim().slice(0, 140)
}

function normalizeForDedupe(value) {
  return value
    .toLowerCase()
    .replace(/\[[^\]]{1,20}\]/g, '')
    .replace(/\{argument name="[^"]+" default="([^"]*)"\}/g, '$1')
    .replace(/[^\p{L}\p{N}]+/gu, ' ')
    .trim()
    .slice(0, 1200)
}

function dedupeEntries(items) {
  const seen = new Set()
  const result = []
  for (const item of items) {
    if (!item.prompt || item.prompt.trim().length < 20) continue
    const key = normalizeForDedupe(item.prompt)
    if (seen.has(key)) continue
    seen.add(key)
    result.push(item)
  }
  return result
}

function isSuitableForPromptLibrary(item) {
  const text = `${item.title}\n${item.description}\n${item.prompt}`.toLowerCase()
  const blocked = [
    /porn/,
    /explicit nudity/,
    /deep cleavage/,
    /lingerie/,
    /underwear/,
    /camisole/,
    /lace-trim camisole/,
    /black garter/,
    /garter straps?/,
    /thigh[-\s]?high/,
    /bikini/,
    /swimsuit/,
    /wet .*swimsuit/,
    /mini skirt/,
    /spaghetti[-\s]?strap/,
    /bodysuit/,
    /off[-\s]?shoulder/,
    /bare shoulders?/,
    /bare back/,
    /slipping off one shoulder/,
    /seductive .*pose/,
    /seductive expression/,
    /subtly seductive/,
    /softly seductive/,
    /sultry/,
    /teenage/,
    /late teens/,
    /high school girl/,
    /high school dxd/,
    /schoolgirl/,
    /school uniform/,
    /诱人.*姿势/,
    /诱人/,
    /性感.*姿势/,
    /丝袜.*内衣/,
    /内衣.*直播/,
    /黑色丝袜/,
    /吊带袜/,
    /深邃乳沟/,
    /乳沟/,
    /极小的.*迷你裙/,
    /色情化/,
    /十几岁/,
    /未成年/,
    /吊带背心/,
    /泳衣/,
    /高叉泳衣/,
    /露背/,
    /裸露/,
    /胸部/,
    /身体线条/,
    /丰腴/,
  ]
  return !blocked.some((pattern) => pattern.test(text))
}

function writePromptLibrary(items) {
  fs.rmSync(outputDir, { recursive: true, force: true })
  fs.mkdirSync(chunksDir, { recursive: true })

  const sorted = items.sort((a, b) => {
    const featuredDiff = Number(Boolean(b.featured)) - Number(Boolean(a.featured))
    if (featuredDiff) return featuredDiff
    return a.sourceId.localeCompare(b.sourceId) || a.category.localeCompare(b.category) || a.title.localeCompare(b.title)
  })

  const chunkSize = 160
  const chunks = []
  const index = sorted.map((item, idx) => {
    const chunkId = `chunk-${Math.floor(idx / chunkSize)}`
    if (!chunks.includes(chunkId)) chunks.push(chunkId)
    return {
      id: item.id,
      locale: item.locale,
      title: item.title,
      category: item.category,
      description: makeDescription(item.description),
      promptPreview: makePreview(item.promptPreview || item.prompt),
      sourceId: item.sourceId,
      sourceLanguage: item.sourceLanguage,
      sourceLabel: item.sourceLabel,
      tags: item.tags || [],
      featured: Boolean(item.featured),
      needsReference: Boolean(item.needsReference),
      chunkId,
    }
  })

  chunks.forEach((chunkId) => {
    const detailRows = sorted
      .filter((_, idx) => `chunk-${Math.floor(idx / chunkSize)}` === chunkId)
      .map((item) => ({
        id: item.id,
        prompt: item.prompt,
        author: item.author,
        sourceUrl: item.sourceUrl,
        publishedAt: item.publishedAt,
      }))
    writeTsFile(path.join(chunksDir, `${chunkId}.ts`), `import type { PromptPresetDetail } from '../types'\n\nexport const PROMPT_LIBRARY_CHUNK: PromptPresetDetail[] = ${stableJson(detailRows)}\n`)
  })

  writeTsFile(path.join(outputDir, 'loaders.ts'), `import type { PromptPresetDetail } from './types'\n\ntype PromptChunkModule = { PROMPT_LIBRARY_CHUNK: PromptPresetDetail[] }\n\nexport async function loadPromptLibraryChunk(chunkId: string): Promise<PromptPresetDetail[]> {\n  switch (chunkId) {\n${chunks.map((chunkId) => `    case '${chunkId}':\n      return ((await import('./chunks/${chunkId}')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK`).join('\n')}\n    default:\n      return []\n  }\n}\n`)

  writeTsFile(path.join(outputDir, 'types.ts'), `export type PromptLibraryLocale = 'all' | 'zh' | 'en'\n\nexport interface PromptPresetSummary {\n  id: string\n  locale: 'zh' | 'en'\n  title: string\n  category: string\n  description: string\n  promptPreview: string\n  sourceId: string\n  sourceLanguage?: string\n  sourceLabel?: string\n  tags: string[]\n  featured: boolean\n  needsReference: boolean\n  chunkId: string\n}\n\nexport interface PromptPresetDetail {\n  id: string\n  prompt: string\n  author?: string\n  sourceUrl?: string\n  publishedAt?: string\n}\n\nexport interface PromptPresetSource {\n  id: string\n  name: string\n  url: string\n  license: string\n  licenseUrl: string\n  commit: string\n}\n`)

  writeTsFile(path.join(outputDir, 'index.ts'), `import type { PromptPresetSource, PromptPresetSummary } from './types'\n\nexport const PROMPT_LIBRARY_SOURCES: PromptPresetSource[] = ${stableJson(sourceDefs)}\n\nexport const PROMPT_LIBRARY_CHUNK_IDS = ${stableJson(chunks)} as const\n\nexport const PROMPT_LIBRARY_INDEX: PromptPresetSummary[] = ${stableJson(index)}\n`)
}

function writeTsFile(filePath, content) {
  fs.writeFileSync(filePath, `${content.trimEnd()}\n`, 'utf8')
}

function stableJson(value) {
  return JSON.stringify(value, null, 2)
}
