<script lang="ts">
  import { onMount } from 'svelte'
  import TranslationPanel from './components/TranslationPanel.svelte'
  import SettingsModal from './components/SettingsModal.svelte'
  import Toast from './components/Toast.svelte'
  import { getProviders, getDefaultLanguages } from './services/wails'
  import type { Provider } from './types'

  // Global state using Svelte 5 runes
  let providers = $state<Provider[]>([])
  let defaultLanguages = $state<Record<string, string>>({})
  let showSettings = $state(false)
  let toastMessage = $state('')
  let toastType = $state<'info' | 'error' | 'success'>('info')
  let toastVisible = $state(false)

  // Toast helper
  function showToast(message: string, type: 'info' | 'error' | 'success' = 'info') {
    toastMessage = message
    toastType = type
    toastVisible = true
    setTimeout(() => {
      toastVisible = false
    }, 3000)
  }

  // Load initial data
  async function loadData() {
    try {
      providers = await getProviders()
      defaultLanguages = await getDefaultLanguages()
    } catch (error) {
      console.error('Failed to load data:', error)
      showToast(String(error), 'error')
    }
  }

  // Reload providers
  async function reloadProviders() {
    providers = await getProviders()
  }

  // Reload default languages
  async function reloadDefaultLanguages() {
    defaultLanguages = await getDefaultLanguages()
  }

  onMount(() => {
    loadData()

    // Listen for clipboard events from backend
    if (window.runtime) {
      window.runtime.EventsOn('set-clipboard-text', (text: string) => {
        // Dispatch custom event that TranslationPanel can listen to
        window.dispatchEvent(new CustomEvent('clipboard-text', { detail: text }))
      })
    }
  })
</script>

<div class="app">
  <div class="drag-region" data-wails-drag></div>

  <main class="container">
    <TranslationPanel {defaultLanguages} onToast={showToast} />
  </main>

  <footer class="footer">
    <div class="version">FanyiHub v1.0</div>
    <button class="settings-btn" onclick={() => (showSettings = true)}>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="3"></circle>
        <path
          d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82-.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
        ></path>
      </svg>
    </button>
  </footer>

  {#if showSettings}
    <SettingsModal
      {providers}
      {defaultLanguages}
      onClose={() => (showSettings = false)}
      onProvidersChange={reloadProviders}
      onLanguagesChange={reloadDefaultLanguages}
      onToast={showToast}
    />
  {/if}

  <Toast message={toastMessage} type={toastType} visible={toastVisible} />
</div>

<style>
  .app {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .container {
    flex: 1;
    display: flex;
    flex-direction: column;
    padding: 0 16px 60px;
    max-width: 1200px;
    margin: 0 auto;
    width: 100%;
    height: 100%;
  }

  .footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 12px 20px;
    background: #fff;
    border-top: 1px solid var(--color-border);
    display: flex;
    justify-content: space-between;
    align-items: center;
    z-index: 100;
  }

  .version {
    color: var(--color-text-secondary);
    font-size: 12px;
  }

  .settings-btn {
    color: var(--color-text-secondary);
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    border-radius: var(--radius-md);
    transition: all var(--transition-fast);
    display: flex;
    align-items: center;
  }

  .settings-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }
</style>
