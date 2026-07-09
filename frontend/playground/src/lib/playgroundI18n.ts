import { getHostLocale, setHostLocale, subscribeHostLocaleChange } from './sub2apiHost'

export type PlaygroundLocale = 'en' | 'zh'

export function getPlaygroundLocale(): PlaygroundLocale {
  return getHostLocale()
}

export function setPlaygroundLocale(locale: PlaygroundLocale) {
  setHostLocale(locale)
  window.dispatchEvent(new CustomEvent('sub2api-playground-locale-change', { detail: { locale } }))
}

export function subscribePlaygroundLocaleChange(callback: () => void): () => void {
  const localHandler = () => callback()
  const unsubscribeStorage = subscribeHostLocaleChange(callback)
  window.addEventListener('sub2api-playground-locale-change', localHandler)
  return () => {
    unsubscribeStorage()
    window.removeEventListener('sub2api-playground-locale-change', localHandler)
  }
}

export function text(zh: string, en: string): string {
  return getPlaygroundLocale() === 'zh' ? zh : en
}

export function useLocaleText(zh: string, en: string, localeVersion = 0): string {
  void localeVersion
  return text(zh, en)
}

const EXACT_EN: Record<string, string> = {
  '占位': '',
  '画廊': 'Gallery',
  '控制台': 'Console',
  '操作指南': 'Help',
  '设置': 'Settings',
  '历史任务': 'History',
  '搜索聊天...': 'Search chats...',
  '新对话': 'New conversation',
  '折叠左侧边栏': 'Collapse sidebar',
  '展开对话列表': 'Expand conversations',
  '滚动到底部': 'Scroll to bottom',
  '切换语言': 'Switch language',
  '下拉展示顶栏': 'Swipe down to show toolbar',
  '全部': 'All',
  '已完成': 'Done',
  '生成中': 'Running',
  '生成中...': 'Generating...',
  '正在生成……': 'Generating...',
  '输入内容将在响应完成时接收': 'Input will appear when the response completes',
  '失败': 'Failed',
  '已停止': 'Stopped',
  '重连中': 'Reconnecting',
  '连接已断开，正在从服务端恢复任务结果。': 'Connection lost. Recovering the server job result.',
  '已停止等待，任务可能仍在服务端完成并计费。': 'Stopped waiting. The job may still finish on the server and be billed.',
  '已停止等待服务端任务': 'Stopped waiting for the server job',
  '已停止生成。': 'Generation stopped.',
  '已停止生成': 'Generation stopped',
  '停止生成': 'Stop generation',
  '上传过程中断，未能确认服务端是否创建任务。请重新上传后再试。': 'Upload interrupted. The server job could not be confirmed. Upload again and retry.',
  '与服务器的连接已断开，任务仍会在服务端继续运行。': 'Connection to the server was lost. The job will continue running on the server.',
  '提交响应丢失，已通过任务 ID 找回服务端任务。': 'Submit response was lost. The server job was found by client task ID.',
  '任务已完成，但没有返回可下载的图片。': 'The job completed but did not return downloadable images.',
  '服务端任务失败。': 'Server job failed.',
  '预览': 'Preview',
  '局部重绘': 'Masked edit',
  '透明背景': 'Transparent',
  '质量': 'Quality',
  '尺寸': 'Size',
  '格式': 'Format',
  '数量': 'Count',
  '参数配置': 'Parameters',
  '来源': 'Source',
  '输入内容': 'Input',
  '参考图': 'Reference images',
  '清空文本': 'Clear text',
  '选择尺寸': 'Choose size',
  '编辑': 'Edit',
  '编辑参考图': 'Edit reference images',
  '清空全部输入图': 'Clear all input images',
  '清空参考图': 'Clear reference images',
  '清空遮罩主图、参考图和遮罩': 'Clear mask target, references, and mask',
  '清空全部参考图': 'Clear all reference images',
  '由模型自主选择': 'Selected by the model',
  '由模型自主选择，可能包含其他图片': 'Selected by the model, may include other images',
  '复用配置': 'Reuse',
  '编辑输出': 'Edit outputs',
  '删除任务': 'Delete task',
  '删除对话': 'Delete conversation',
  '删除轮次': 'Delete round',
  '删除消息': 'Delete message',
  '切换到画廊模式？': 'Switch to Gallery?',
  '切换并复用': 'Switch and reuse',
  '收藏任务': 'Favorite',
  '编辑收藏夹': 'Edit collections',
  '重试任务': 'Retry',
  '下载图片': 'Download image',
  '下载全部': 'Download all',
  '下载原图': 'Download original',
  '复制完整报错': 'Copy full error',
  '查看原始响应': 'View raw response',
  '复制图片链接': 'Copy image link',
  '下载中间步骤图': 'Download partial images',
  '提示词已复制': 'Prompt copied',
  '复制提示词': 'Copy prompt',
  '提示词参考': 'Prompt Library',
  '关闭': 'Close',
  '确认': 'Confirm',
  '取消': 'Cancel',
  '确认删除': 'Confirm delete',
  '清除': 'Clear',
  '没有失败记录': 'No failed records',
  '管理收藏夹': 'Manage collections',
  '返回收藏夹': 'Back to collection',
  '退出收藏夹': 'Exit collections',
  '收藏夹': 'Collections',
  '多选任务': 'Multi-select tasks',
  '多选收藏夹': 'Multi-select collections',
  '批量操作': 'Batch actions',
  '自动': 'Auto',
  '按比例': 'By ratio',
  '自定义宽高': 'Custom size',
  '设置图像尺寸': 'Image size',
  '确定': 'Apply',
  '正在载入图片...': 'Loading image...',
  '编辑遮罩': 'Edit mask',
  '移除遮罩': 'Remove mask',
  '保存中...': 'Saving...',
  '保存': 'Save',
  '画笔': 'Brush',
  '橡皮': 'Eraser',
  '撤销': 'Undo',
  '重做': 'Redo',
  '重置视图': 'Reset view',
  '清空遮罩': 'Clear mask',
  '复制': 'Copy',
  '下载': 'Download',
  '下载成功': 'Downloaded',
  '下载失败': 'Download failed',
  '复制成功': 'Copied',
  '复制失败': 'Copy failed',
  '加载中': 'Loading',
  '没有找到匹配的任务': 'No matching tasks',
  '输入提示词开始生成图片': 'Enter a prompt to start generating',
  '正在生成': 'Generating',
  '点击复制完整报错': 'Click to copy full error',
  '使用': 'Use',
  '追加': 'Append',
  '原链接': 'Original',
  '参考图数量已达上限（16 张），无法继续添加': 'Reference image limit reached (16)',
}

const REGEX_EN: Array<[RegExp, (match: RegExpMatchArray) => string]> = [
  [/^第 (\d+) 张生成失败$/, (m) => `Image ${m[1]} failed`],
  [/^创建于 (.+)$/, (m) => `Created at ${m[1]}`],
  [/^耗时 (.+)$/, (m) => `Elapsed ${m[1]}`],
  [/^(.+) 条任务$/, (m) => `${m[1]} tasks`],
  [/^清除 (\d+) 条失败记录$/, (m) => `Clear ${m[1]} failed records`],
  [/^搜索收藏夹名称\.\.\.$/, () => 'Search collection names...'],
  [/^搜索提示词、参数\.\.\.$/, () => 'Search prompts and parameters...'],
  [/^描述你想生成的图片，可输入 @ 来指定参考图\.\.\.$/, () => 'Describe the image. Type @ to reference images...'],
  [/^共 (\d+) 条参考$/, (m) => `${m[1]} prompts`],
  [/^图片 (\d+)$/, (m) => `Image ${m[1]}`],
  [/^第 (.+) 轮$/, (m) => `Round ${m[1]}`],
  [/^确定要删除这个任务吗？关联的图片资源也会被清理（如果没有其他任务引用）。$/, () => 'Delete this task? Related image assets will also be cleaned up if no other task references them.'],
  [/^确定要删除这个 Agent 对话吗？$/, () => 'Delete this Agent conversation?'],
  [/^复用参数会应用到画廊输入区。切换到画廊模式后，当前 Agent 对话仍会保留。$/, () => 'The reused parameters will be applied to the Gallery input. The current Agent conversation will remain after switching.'],
  [/^确定要删除选中的 (\d+) 个任务吗？$/, (m) => `Delete ${m[1]} selected tasks?`],
  [/^确定要删除选中的 (\d+) 个收藏夹吗？$/, (m) => `Delete ${m[1]} selected collections?`],
  [/^确定要删除收藏夹「(.+)」吗？$/, (m) => `Delete collection "${m[1]}"?`],
  [/^确定要将默认收藏夹从「(.+)」改为「(.+)」吗？$/, (m) => `Change the default collection from "${m[1]}" to "${m[2]}"?`],
  [/^清除失败记录$/, () => 'Clear failed records'],
  [/^请输入提示词$/, () => 'Enter a prompt'],
  [/^请输入消息$/, () => 'Enter a message'],
  [/^服务端任务已恢复，共 (\d+) 张图片$/, (m) => `Server job recovered with ${m[1]} image${m[1] === '1' ? '' : 's'}`],
  [/^生成完成，共 (\d+) 张图片$/, (m) => `Generation complete, ${m[1]} image${m[1] === '1' ? '' : 's'}`],
  [/^生成完成：成功 (\d+) 张，失败 (\d+) 张$/, (m) => `Generation complete: ${m[1]} succeeded, ${m[2]} failed`],
  [/^文件太大：单文件上限为 (.+)$/, (m) => `File too large: per-file limit is ${m[1]}`],
  [/^请求体太大：当前约 (.+)$/, (m) => `Request body too large: current estimate is ${m[1]}`],
  [/^图片下载失败：HTTP (\d+)$/, (m) => `Image download failed: HTTP ${m[1]}`],
]

export function translateUiString(value: string): string {
  if (getPlaygroundLocale() !== 'en') return value
  const normalized = value.replace(/\s+/g, ' ').trim()
  if (!normalized) return value
  const exact = EXACT_EN[normalized]
  if (exact !== undefined) return exact
  for (const [pattern, render] of REGEX_EN) {
    const match = normalized.match(pattern)
    if (match) return render(match)
  }
  return value
}

function shouldSkipNode(node: Node): boolean {
  const element = node.nodeType === Node.ELEMENT_NODE
    ? node as Element
    : node.parentElement
  if (!element) return false
  return Boolean(element.closest('[data-selectable-text], [data-skip-ui-translate], .markdown-renderer, textarea, input, code, pre'))
}

function translateAttributes(root: ParentNode) {
  root.querySelectorAll<HTMLElement>('[title], [aria-label], [placeholder]').forEach((element) => {
    if (element.closest('[data-selectable-text], [data-skip-ui-translate], .markdown-renderer, code, pre')) return
    for (const attr of ['title', 'aria-label', 'placeholder']) {
      const value = element.getAttribute(attr)
      if (!value) continue
      const translated = translateUiString(value)
      if (translated !== value) element.setAttribute(attr, translated)
    }
  })
}

function translateTextNodes(root: ParentNode) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT)
  const updates: Array<[Text, string]> = []
  while (walker.nextNode()) {
    const node = walker.currentNode as Text
    if (shouldSkipNode(node)) continue
    const value = node.nodeValue ?? ''
    const translated = translateUiString(value)
    if (translated !== value) updates.push([node, translated])
  }
  updates.forEach(([node, value]) => {
    node.nodeValue = value
  })
}

export function translateStaticUi(root: ParentNode = document.body) {
  if (getPlaygroundLocale() !== 'en') return
  translateAttributes(root)
  translateTextNodes(root)
}
