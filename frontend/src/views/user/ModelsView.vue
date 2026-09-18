<template>
  <AppLayout>
    <div class="space-y-6">
      <section class="space-y-2">
        <h1 class="text-2xl font-semibold text-txt-primary">{{ t('modelsPage.title') }}</h1>
        <p class="text-sm text-txt-secondary">
          {{ t('modelsPage.description') }}
        </p>
      </section>

      <section class="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
        <SearchInput
          v-model="searchQuery"
          :placeholder="t('modelsPage.searchPlaceholder')"
          class="w-full xl:w-80"
        />

        <div class="flex flex-wrap gap-3">
          <button
            type="button"
            class="btn btn-secondary"
            :disabled="loading || !hasEntries || allExpanded"
            @click="expandAll"
          >
            {{ t('modelsPage.expandAll') }}
          </button>
          <button
            type="button"
            class="btn btn-secondary"
            :disabled="!hasExpandedGroups"
            @click="collapseAll"
          >
            {{ t('modelsPage.collapseAll') }}
          </button>
          <button type="button" class="btn btn-secondary" :disabled="loading" @click="reloadAll">
            <Icon name="refresh" size="sm" class="mr-2" :class="loading ? 'animate-spin' : ''" />
            {{ t('common.refresh') }}
          </button>
        </div>
      </section>

      <div
        v-if="searchBackfillLoading"
        class="flex items-center gap-3 rounded-lg border border-blue-200 bg-blue-50 px-4 py-3 text-sm text-blue-700 dark:bg-blue-900/20 dark:text-blue-200"
      >
        <Icon name="refresh" size="sm" class="animate-spin" />
        <span>{{ t('modelsPage.searchLoading') }}</span>
      </div>

      <div v-if="loading && modelEntries.length === 0" class="flex items-center justify-center py-12">
        <LoadingSpinner />
      </div>

      <EmptyState
        v-else-if="!loading && filteredEntries.length === 0 && modelEntries.length === 0"
        :title="t('modelsPage.emptyTitle')"
        :description="t('modelsPage.emptyDescription')"
      />

      <EmptyState
        v-else-if="!loading && filteredEntries.length === 0"
        :title="t('modelsPage.noMatchTitle')"
        :description="t('modelsPage.noMatchDescription')"
      />

      <section v-else class="grid gap-4 lg:grid-cols-2">
        <article
          v-for="entry in filteredEntries"
          :key="entry.group.id"
          class="rounded-lg border border-cyber-border bg-cyber-surface p-5"
        >
          <div class="flex flex-wrap items-start justify-between gap-3">
            <div class="min-w-0 space-y-2">
              <GroupBadge
                :name="entry.group.name"
                :platform="entry.group.platform"
                :subscription-type="entry.group.subscription_type"
                :rate-multiplier="entry.group.rate_multiplier"
              />
              <p v-if="entry.group.description" class="text-sm text-txt-secondary">
                {{ entry.group.description }}
              </p>
            </div>

            <div class="flex flex-wrap gap-2">
              <button
                type="button"
                class="btn btn-secondary btn-sm"
                @click="toggleGroup(entry.group.id)"
              >
                <Icon
                  :name="isExpanded(entry.group.id) ? 'chevronDown' : 'chevronRight'"
                  size="xs"
                  class="mr-1"
                />
                {{
                  isExpanded(entry.group.id)
                    ? t('modelsPage.collapseModels')
                    : t('modelsPage.viewModels')
                }}
              </button>

              <button
                v-if="isExpanded(entry.group.id)"
                type="button"
                class="btn btn-secondary btn-sm"
                :disabled="entry.loading"
                @click="reloadGroup(entry.group.id)"
              >
                <Icon
                  name="refresh"
                  size="xs"
                  class="mr-1"
                  :class="entry.loading ? 'animate-spin' : ''"
                />
                {{ t('common.refresh') }}
              </button>
            </div>
          </div>

          <div
            v-if="isExpanded(entry.group.id)"
            class="mt-4 space-y-4 border-t border-cyber-border pt-4"
          >
            <div
              v-if="entry.loading"
              class="flex items-center gap-2 rounded-lg border border-cyber-border bg-cyber-bg px-4 py-3 text-sm text-txt-secondary"
            >
              <Icon name="refresh" size="sm" class="animate-spin" />
              <span>{{ t('common.loading') }}</span>
            </div>

            <div
              v-else-if="entry.error"
              class="space-y-3 rounded-lg border border-red-200 bg-red-50 px-4 py-4 text-sm text-red-700"
            >
              <p>{{ entry.error }}</p>
              <div>
                <button type="button" class="btn btn-secondary btn-sm" @click="reloadGroup(entry.group.id)">
                  {{ t('common.retry') }}
                </button>
              </div>
            </div>

            <template v-else-if="entry.data">
              <div class="flex flex-wrap items-center gap-2 text-sm">
                <span class="rounded-md bg-primary-50 px-2.5 py-1 font-medium text-primary-700 dark:bg-primary-900/20 dark:text-primary-300">
                  {{ t('modelsPage.modelCount', { count: entry.data.models.length }) }}
                </span>
                <span
                  v-if="entry.data.is_default"
                  class="rounded-md bg-blue-50 px-2.5 py-1 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300"
                >
                  {{ t('modelsPage.defaultTag') }}
                </span>
              </div>

              <p
                v-if="entry.data.is_default"
                class="rounded-lg border border-blue-200 bg-blue-50 px-4 py-3 text-sm text-blue-700 dark:bg-blue-900/20 dark:text-blue-200"
              >
                {{ t('modelsPage.defaultHint') }}
              </p>

              <div
                v-if="entry.data.models.length > 0"
                class="grid gap-2 sm:grid-cols-2"
              >
                <div
                  v-for="model in entry.data.models"
                  :key="model"
                  class="rounded-md border border-cyber-border bg-cyber-bg px-3 py-2 font-mono text-sm text-txt-primary break-all"
                >
                  {{ model }}
                </div>
              </div>

              <div
                v-else
                class="rounded-lg border border-cyber-border bg-cyber-bg px-4 py-6 text-center text-sm text-txt-secondary"
              >
                {{ t('modelsPage.noModels') }}
              </div>
            </template>
          </div>
        </article>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import AppLayout from '@/components/layout/AppLayout.vue'
import SearchInput from '@/components/common/SearchInput.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import Icon from '@/components/icons/Icon.vue'
import { userGroupsAPI, type GroupModelsResponse } from '@/api'
import type { Group } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

interface GroupModelsEntry {
  group: Group
  loading: boolean
  error: string
  data: GroupModelsResponse | null
}

const loading = ref(false)
const searchQuery = ref('')
const searchBackfillLoading = ref(false)
const modelEntries = ref<GroupModelsEntry[]>([])
const expandedGroupIds = ref<number[]>([])

let modelEntriesVersion = 0
let groupLoadSequence = 0
let searchBackfillPromise: Promise<void> | null = null
let activeSearchBackfillToken: symbol | null = null

const pendingLoads = new Map<number, Promise<void>>()
const latestLoadIds = new Map<number, number>()
const pendingLoadTokens = new Map<number, symbol>()

const hasEntries = computed(() => modelEntries.value.length > 0)
const hasExpandedGroups = computed(() => expandedGroupIds.value.length > 0)
const allExpanded = computed(() => {
  return hasEntries.value && modelEntries.value.every((entry) => expandedGroupIds.value.includes(entry.group.id))
})

const filteredEntries = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return modelEntries.value

  return modelEntries.value.filter((entry) => {
    if (entry.group.name.toLowerCase().includes(query)) return true
    if ((entry.group.description || '').toLowerCase().includes(query)) return true
    return (entry.data?.models || []).some((model) => model.toLowerCase().includes(query))
  })
})

const getEntry = (groupId: number) => {
  return modelEntries.value.find((entry) => entry.group.id === groupId) || null
}

const shouldLoadModels = (entry: GroupModelsEntry | null) => {
  return !!entry && entry.data === null && !entry.loading && !entry.error
}

const updateEntry = (groupId: number, patch: Partial<GroupModelsEntry>, version = modelEntriesVersion) => {
  if (version !== modelEntriesVersion) return

  const index = modelEntries.value.findIndex((entry) => entry.group.id === groupId)
  if (index === -1) return

  modelEntries.value[index] = {
    ...modelEntries.value[index],
    ...patch
  }
}

const isExpanded = (groupId: number) => expandedGroupIds.value.includes(groupId)

const expandGroup = (groupId: number) => {
  if (isExpanded(groupId)) return
  expandedGroupIds.value = [...expandedGroupIds.value, groupId]
}

const collapseGroup = (groupId: number) => {
  expandedGroupIds.value = expandedGroupIds.value.filter((id) => id !== groupId)
}

const loadGroupModels = async (groupId: number, options: { force?: boolean } = {}) => {
  const entry = getEntry(groupId)
  if (!entry) return

  if (!options.force) {
    if (entry.data) return

    const existingRequest = pendingLoads.get(groupId)
    if (existingRequest) {
      return existingRequest
    }
  }

  const version = modelEntriesVersion
  const requestId = ++groupLoadSequence
  const pendingToken = Symbol(`group-${groupId}`)
  latestLoadIds.set(groupId, requestId)
  pendingLoadTokens.set(groupId, pendingToken)
  updateEntry(groupId, { loading: true, error: '' }, version)

  const request = (async () => {
    try {
      const data = await userGroupsAPI.getModels(groupId)
      if (version !== modelEntriesVersion || latestLoadIds.get(groupId) !== requestId) return

      updateEntry(groupId, {
        data,
        loading: false,
        error: ''
      }, version)
    } catch (error) {
      if (version !== modelEntriesVersion || latestLoadIds.get(groupId) !== requestId) return

      updateEntry(groupId, {
        loading: false,
        error: t('modelsPage.failedToLoad')
      }, version)
    } finally {
      if (pendingLoadTokens.get(groupId) === pendingToken) {
        pendingLoadTokens.delete(groupId)
        pendingLoads.delete(groupId)
      }
    }
  })()

  pendingLoads.set(groupId, request)
  return request
}

const hydrateSearchModels = async () => {
  if (!searchQuery.value.trim()) return

  const missingGroupIds = modelEntries.value
    .filter((entry) => shouldLoadModels(entry))
    .map((entry) => entry.group.id)

  if (missingGroupIds.length === 0) {
    searchBackfillLoading.value = false
    return
  }

  if (searchBackfillPromise) {
    return searchBackfillPromise
  }

  const version = modelEntriesVersion
  const backfillToken = Symbol('search-backfill')
  activeSearchBackfillToken = backfillToken
  searchBackfillLoading.value = true

  const request = (async () => {
    try {
      await Promise.allSettled(missingGroupIds.map((groupId) => loadGroupModels(groupId)))
    } finally {
      if (version === modelEntriesVersion && activeSearchBackfillToken === backfillToken) {
        searchBackfillLoading.value = false
        searchBackfillPromise = null
        activeSearchBackfillToken = null
      }
    }
  })()

  searchBackfillPromise = request
  return request
}

const loadAvailableGroups = async () => {
  const version = ++modelEntriesVersion
  pendingLoads.clear()
  latestLoadIds.clear()
  pendingLoadTokens.clear()
  searchBackfillPromise = null
  activeSearchBackfillToken = null
  searchBackfillLoading.value = false
  modelEntries.value = []
  loading.value = true

  try {
    const groups = await userGroupsAPI.getAvailable()
    if (version !== modelEntriesVersion) return

    modelEntries.value = groups.map((group) => ({
      group,
      loading: false,
      error: '',
      data: null
    }))

    expandedGroupIds.value = groups.map((group) => group.id)
    await Promise.allSettled(groups.map((group) => loadGroupModels(group.id)))

    if (searchQuery.value.trim()) {
      void hydrateSearchModels()
    }
  } catch (error) {
    if (version !== modelEntriesVersion) return
    console.error('Failed to load available groups:', error)
    appStore.showError(t('modelsPage.failedToLoad'))
  } finally {
    if (version === modelEntriesVersion) {
      loading.value = false
    }
  }
}

const toggleGroup = (groupId: number) => {
  if (isExpanded(groupId)) {
    collapseGroup(groupId)
    return
  }

  expandGroup(groupId)

  const entry = getEntry(groupId)
  if (shouldLoadModels(entry)) {
    void loadGroupModels(groupId)
  }
}

const expandAll = async () => {
  expandedGroupIds.value = modelEntries.value.map((entry) => entry.group.id)
  await Promise.allSettled(
    modelEntries.value
      .filter((entry) => shouldLoadModels(entry))
      .map((entry) => loadGroupModels(entry.group.id))
  )
}

const collapseAll = () => {
  expandedGroupIds.value = []
}

const reloadAll = async () => {
  expandedGroupIds.value = []
  await loadAvailableGroups()
}

const reloadGroup = async (groupId: number) => {
  expandGroup(groupId)
  await loadGroupModels(groupId, { force: true })
}

watch(searchQuery, (value) => {
  if (!value.trim()) {
    searchBackfillLoading.value = false
    return
  }

  void hydrateSearchModels()
})

onMounted(() => {
  void loadAvailableGroups()
})
</script>
