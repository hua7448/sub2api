import { useEffect, useRef, useState } from 'react'
import { createPortal } from 'react-dom'
import type { AppMode } from '../types'
import { useCloseOnEscape } from '../hooks/useCloseOnEscape'
import { usePreventBackgroundScroll } from '../hooks/usePreventBackgroundScroll'
import { hostText } from '../lib/sub2apiHost'

interface HelpModalProps {
  appMode: AppMode
  isFavoriteCollectionOverview?: boolean
  onClose: () => void
}

function useIsMobile() {
  const [isMobile, setIsMobile] = useState(window.innerWidth < 640)
  useEffect(() => {
    const onResize = () => setIsMobile(window.innerWidth < 640)
    window.addEventListener('resize', onResize)
    return () => window.removeEventListener('resize', onResize)
  }, [])
  return isMobile
}

function HelpList({ items }: { items: string[] }) {
  return (
    <ul className="list-disc space-y-2 pl-4">
      {items.map((item) => (
        <li key={item}>{item}</li>
      ))}
    </ul>
  )
}

export default function HelpModal({ appMode, isFavoriteCollectionOverview = false, onClose }: HelpModalProps) {
  const isMobile = useIsMobile()
  const modalRef = useRef<HTMLDivElement>(null)
  const isAgentMode = appMode === 'agent'
  useCloseOnEscape(true, onClose)
  usePreventBackgroundScroll(true, modalRef)

  const selectionItems = isMobile
    ? [hostText('在卡片上左右滑动即可选中或取消选中。', 'Swipe left or right on cards to select or unselect them.')]
    : [
        hostText('在空白处拖拽框选卡片。', 'Drag on empty space to box-select cards.'),
        hostText('按住 Ctrl 或 Command 并点击卡片，可添加或移除单项。', 'Hold Ctrl or Command and click a card to add or remove one item.'),
        hostText('再次框选已选中的卡片会将其取消选中。', 'Box-selecting selected cards again unselects them.'),
        hostText('点击卡片外空白处可取消所有选择。', 'Click empty space outside cards to clear the selection.'),
      ]

  return createPortal(
    <div
      data-no-drag-select
      className="fixed inset-0 z-[100] flex items-center justify-center p-4"
      onClick={onClose}
    >
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm animate-overlay-in" />
      <div
        ref={modalRef}
        className="relative z-10 flex max-h-[85vh] w-full max-w-md flex-col rounded-2xl border border-amber-200/80 bg-stone-50/95 p-5 shadow-2xl ring-1 ring-amber-900/5 animate-modal-in dark:border-amber-400/20 dark:bg-stone-950/95 dark:ring-white/10"
        onClick={(event) => event.stopPropagation()}
      >
        <div className="mb-5 flex items-center justify-between gap-4">
          <h3 className="flex items-center gap-2 text-base font-semibold text-stone-800 dark:text-stone-100">
            <svg className="h-5 w-5 text-orange-500" fill="none" stroke="currentColor" strokeWidth={2} strokeLinecap="round" strokeLinejoin="round" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" />
              <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3" />
              <path d="M12 17h.01" />
            </svg>
            {hostText('操作指南', 'Help')}
          </h3>
          <button
            onClick={onClose}
            className="rounded-full p-1 text-stone-400 transition hover:bg-amber-50 hover:text-stone-700 dark:hover:bg-white/[0.06] dark:hover:text-stone-200"
            aria-label={hostText('关闭', 'Close')}
          >
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="custom-scrollbar mb-2 flex-1 space-y-6 overflow-y-auto overscroll-contain pr-2 text-sm leading-6 text-stone-600 dark:text-stone-300">
          {isAgentMode ? (
            <section>
              <HelpList
                items={[
                  hostText('Agent 模式需要使用 Responses API 配置。', 'Agent mode requires a Responses API configuration.'),
                  hostText('如需联网搜索，可在设置的 Agent 配置中开启网络搜索。', 'Enable web search in Agent settings when internet search is needed.'),
                  hostText('输入 @ 可引用参考图或前面轮次生成的图片。', 'Type @ to reference uploaded images or images from earlier rounds.'),
                  hostText('编辑或重新生成某轮消息会产生可切换的分支。', 'Editing or regenerating a round creates switchable branches.'),
                  hostText('生成的图片会同步到画廊；删除对话默认不会删除画廊任务。', 'Generated images are synced to Gallery; deleting a conversation does not delete gallery tasks by default.'),
                ]}
              />
            </section>
          ) : (
            <>
              <section>
                <h4 className="mb-3 text-sm font-semibold text-stone-800 dark:text-stone-100">
                  {isFavoriteCollectionOverview
                    ? hostText('多选收藏夹', 'Multi-select collections')
                    : hostText('多选任务', 'Multi-select tasks')}
                </h4>
                <HelpList items={selectionItems} />
              </section>
              <section>
                <h4 className="mb-3 text-sm font-semibold text-stone-800 dark:text-stone-100">
                  {hostText('批量操作', 'Batch actions')}
                </h4>
                <p>
                  {isFavoriteCollectionOverview
                    ? hostText('选中一个或多个收藏夹后，底部操作栏支持取消选择、全选、反选、下载选中和删除选中。', 'After selecting one or more collections, the bottom bar supports clear selection, select all, invert, download selected, and delete selected.')
                    : hostText('选中一个或多个任务后，底部操作栏支持取消选择、全选、反选、编辑收藏夹、下载选中和删除选中。', 'After selecting one or more tasks, the bottom bar supports clear selection, select all, invert, edit collections, download selected, and delete selected.')}
                </p>
              </section>
            </>
          )}
        </div>
      </div>
    </div>,
    document.body,
  )
}
