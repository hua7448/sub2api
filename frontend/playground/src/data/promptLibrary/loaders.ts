import type { PromptPresetDetail } from './types'

type PromptChunkModule = { PROMPT_LIBRARY_CHUNK: PromptPresetDetail[] }

export async function loadPromptLibraryChunk(chunkId: string): Promise<PromptPresetDetail[]> {
  switch (chunkId) {
    case 'chunk-0':
      return ((await import('./chunks/chunk-0')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK
    case 'chunk-1':
      return ((await import('./chunks/chunk-1')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK
    case 'chunk-2':
      return ((await import('./chunks/chunk-2')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK
    case 'chunk-3':
      return ((await import('./chunks/chunk-3')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK
    case 'chunk-4':
      return ((await import('./chunks/chunk-4')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK
    case 'chunk-5':
      return ((await import('./chunks/chunk-5')) as PromptChunkModule).PROMPT_LIBRARY_CHUNK
    default:
      return []
  }
}
