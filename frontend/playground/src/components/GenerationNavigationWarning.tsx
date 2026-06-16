import { useEffect } from 'react'
import { useStore } from '../store'
import { text } from '../lib/playgroundI18n'

function useHasActiveGeneration() {
  return useStore((state) => (
    state.tasks.some((task) => task.status === 'running') ||
    state.agentConversations.some((conversation) =>
      conversation.rounds.some((round) => round.status === 'running'),
    )
  ))
}

export default function GenerationNavigationWarning() {
  const hasActiveGeneration = useHasActiveGeneration()

  useEffect(() => {
    if (!hasActiveGeneration) return

    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      const message = text(
        '图片生成仍在进行中，请等待完成后再切换页面或关闭窗口。',
        'Image generation is still running. Please wait until it finishes before switching pages or closing this window.',
      )
      event.preventDefault()
      event.returnValue = message
      return message
    }

    window.addEventListener('beforeunload', handleBeforeUnload)
    return () => window.removeEventListener('beforeunload', handleBeforeUnload)
  }, [hasActiveGeneration])

  if (!hasActiveGeneration) return null

  return (
    <div className="pointer-events-none fixed left-1/2 top-[calc(var(--safe-area-top)+4.25rem)] z-40 w-[min(calc(100vw-1.5rem),44rem)] -translate-x-1/2 px-3">
      <div className="rounded-lg border border-amber-300/80 bg-amber-50/95 px-4 py-2.5 text-center text-sm font-medium leading-5 text-amber-900 shadow-[0_10px_30px_rgba(146,64,14,0.16)] ring-1 ring-amber-900/5 backdrop-blur dark:border-amber-400/25 dark:bg-amber-950/90 dark:text-amber-100">
        {text(
          '图片生成中，请不要切换页面或关闭窗口，完成后再离开。',
          'Image generation is running. Do not switch pages or close this window until it finishes.',
        )}
      </div>
    </div>
  )
}
