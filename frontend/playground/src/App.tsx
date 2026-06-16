import { useEffect, useState } from 'react'
import { initStore } from './store'
import { useStore } from './store'
import { buildSettingsFromUrlParams, clearUrlSettingParams, hasUrlSettingParams } from './lib/urlSettings'
import type { AppSettings } from './types'
import Header from './components/Header'
import SearchBar from './components/SearchBar'
import TaskGrid from './components/TaskGrid'
import AgentWorkspace from './components/AgentWorkspace'
import InputBar from './components/InputBar'
import DetailModal from './components/DetailModal'
import Lightbox from './components/Lightbox'
import SettingsModal from './components/SettingsModal'
import ConfirmDialog from './components/ConfirmDialog'
import Toast from './components/Toast'
import MaskEditorModal from './components/MaskEditorModal'
import ImageContextMenu from './components/ImageContextMenu'
import LoginRequiredDialog from './components/LoginRequiredDialog'
import { FavoriteCollectionPickerModal, FavoriteCollectionsView, ManageCollectionsModal } from './components/FavoriteCollections'
import { useGlobalClickSuppression } from './lib/clickSuppression'
import { applySub2APISettings, fetchSub2APIEligibleKeys, fetchSub2APISettings } from './lib/sub2api'
import { applyHostDocumentChrome } from './lib/sub2apiHost'
import { subscribePlaygroundLocaleChange, translateStaticUi } from './lib/playgroundI18n'

export default function App() {
  const setSettings = useStore((s) => s.setSettings)
  const appMode = useStore((s) => s.appMode)
  const filterFavorite = useStore((s) => s.filterFavorite)
  const activeFavoriteCollectionId = useStore((s) => s.activeFavoriteCollectionId)
  const [, setHostLocaleVersion] = useState(0)
  const [loginRequired, setLoginRequired] = useState(false)
  useGlobalClickSuppression()

  useEffect(() => {
    const updateHostChrome = () => {
      applyHostDocumentChrome()
      setHostLocaleVersion((version) => version + 1)
      window.requestAnimationFrame(() => translateStaticUi())
    }
    updateHostChrome()
    return subscribePlaygroundLocaleChange(updateHostChrome)
  }, [])

  useEffect(() => {
    translateStaticUi()
  })

  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search)

    const applyUrlSettings = (baseSettings: Partial<AppSettings>) => {
      const nextSettings = buildSettingsFromUrlParams(baseSettings, searchParams)
      return Object.keys(nextSettings).length ? nextSettings : baseSettings
    }

    const clearAppliedUrlSettings = () => {
      if (!hasUrlSettingParams(searchParams)) return

      clearUrlSettingParams(searchParams)

      const nextSearch = searchParams.toString()
      const nextUrl = `${window.location.pathname}${nextSearch ? `?${nextSearch}` : ''}${window.location.hash}`
      window.history.replaceState(null, '', nextUrl)
    }

    void (async () => {
      try {
        await initStore()
        const [remoteSettings, keys] = await Promise.all([fetchSub2APISettings(), fetchSub2APIEligibleKeys()])
        const baseSettings = applySub2APISettings(useStore.getState().settings, remoteSettings, keys)
        setSettings(applyUrlSettings(baseSettings))
        clearAppliedUrlSettings()
      } catch (error) {
        const message = error instanceof Error ? error.message : String(error ?? '')
        console.warn('Failed to initialize sub2api playground:', error)
        if (message.includes('401') || message.toLowerCase().includes('unauthorized')) {
          setLoginRequired(true)
        }
      }
    })()
  }, [setSettings])

  useEffect(() => {
    const preventPageImageDrag = (e: DragEvent) => {
      if ((e.target as HTMLElement | null)?.closest('img')) {
        e.preventDefault()
      }
    }

    document.addEventListener('dragstart', preventPageImageDrag)
    return () => document.removeEventListener('dragstart', preventPageImageDrag)
  }, [])

  return (
    <>
      <Header />
      {appMode === 'agent' ? (
        <AgentWorkspace />
      ) : (
        <main data-home-main data-drag-select-surface className="pb-48">
          <div className="safe-area-x max-w-7xl mx-auto">
            <SearchBar />
            {filterFavorite && !activeFavoriteCollectionId ? <FavoriteCollectionsView /> : <TaskGrid />}
          </div>
        </main>
      )}
      <InputBar />
      <DetailModal />
      <Lightbox />
      <SettingsModal />
      <ConfirmDialog />
      <FavoriteCollectionPickerModal />
      <ManageCollectionsModal />
      <Toast />
      <MaskEditorModal />
      <ImageContextMenu />
      <LoginRequiredDialog open={loginRequired} />
    </>
  )
}
