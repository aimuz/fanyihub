<script lang="ts">
  import Modal from './Modal.svelte'
  import ProviderCard from './ProviderCard.svelte'
  import ProviderModal from './ProviderModal.svelte'
  import { setDefaultLanguage } from '../services/wails'
  import type { Provider } from '../types'

  type Props = {
    providers: Provider[]
    defaultLanguages: Record<string, string>
    onClose: () => void
    onProvidersChange: () => void
    onLanguagesChange: () => void
    onToast: (message: string, type?: 'info' | 'error' | 'success') => void
  }

  let {
    providers,
    defaultLanguages,
    onClose,
    onProvidersChange,
    onLanguagesChange,
    onToast,
  }: Props = $props()

  // State
  let showAddProvider = $state(false)
  let editingProvider = $state<Provider | null>(null)
  let defaultZhTarget = $state('en')
  let defaultEnTarget = $state('zh')

  // Sync defaults when props change
  $effect(() => {
    defaultZhTarget = defaultLanguages['zh'] || 'en'
    defaultEnTarget = defaultLanguages['en'] || 'zh'
  })

  // Save default languages
  async function saveDefaultLanguages() {
    try {
      await setDefaultLanguage('zh', defaultZhTarget)
      await setDefaultLanguage('en', defaultEnTarget)
      onLanguagesChange()
      onToast('默认翻译语言设置已保存', 'success')
    } catch (error) {
      onToast(String(error), 'error')
    }
  }

  // Handle provider modal close
  function handleProviderModalClose() {
    showAddProvider = false
    editingProvider = null
  }

  // Handle provider saved
  function handleProviderSaved() {
    onProvidersChange()
    handleProviderModalClose()
  }
</script>

<Modal title="翻译服务设置" {onClose}>
  {#snippet children()}
    <div class="settings-section">
      <h3>默认翻译语言</h3>
      <p class="settings-description">当检测到以下语言时，自动设置目标语言</p>
      <div class="default-language-settings">
        <div class="form-group">
          <label for="default-zh-target">检测到中文时，翻译为：</label>
          <select id="default-zh-target" bind:value={defaultZhTarget}>
            <option value="en">英语</option>
            <option value="ja">日语</option>
            <option value="ko">韩语</option>
            <option value="fr">法语</option>
            <option value="de">德语</option>
            <option value="es">西班牙语</option>
            <option value="ru">俄语</option>
            <option value="auto">自动</option>
          </select>
        </div>
        <div class="form-group">
          <label for="default-en-target">检测到英语时，翻译为：</label>
          <select id="default-en-target" bind:value={defaultEnTarget}>
            <option value="zh">中文</option>
            <option value="ja">日语</option>
            <option value="ko">韩语</option>
            <option value="fr">法语</option>
            <option value="de">德语</option>
            <option value="es">西班牙语</option>
            <option value="ru">俄语</option>
            <option value="auto">自动</option>
          </select>
        </div>
        <button class="btn btn-primary" onclick={saveDefaultLanguages}>保存默认语言设置</button>
      </div>
    </div>

    <div class="settings-section">
      <h3>翻译提供商</h3>
      <div class="providers-container">
        {#if providers.length === 0}
          <div class="empty-state">还没有添加任何翻译提供商</div>
        {:else}
          {#each providers as provider (provider.name)}
            <ProviderCard
              {provider}
              onEdit={() => (editingProvider = provider)}
              onChange={onProvidersChange}
              {onToast}
            />
          {/each}
        {/if}
      </div>
      <button class="add-provider-btn" onclick={() => (showAddProvider = true)}
        >添加 LLM 提供商</button
      >
    </div>
  {/snippet}
</Modal>

{#if showAddProvider}
  <ProviderModal onClose={handleProviderModalClose} onSave={handleProviderSaved} {onToast} />
{/if}

{#if editingProvider}
  <ProviderModal
    provider={editingProvider}
    onClose={handleProviderModalClose}
    onSave={handleProviderSaved}
    {onToast}
  />
{/if}

<style>
  .settings-section {
    margin-bottom: 30px;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: 20px;
  }

  .settings-section:last-child {
    border-bottom: none;
    margin-bottom: 0;
  }

  .settings-section h3 {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 12px;
    color: var(--color-text);
  }

  .settings-description {
    font-size: 14px;
    color: var(--color-text-secondary);
    margin-bottom: 16px;
  }

  .default-language-settings {
    background: var(--color-surface);
    padding: 16px;
    border-radius: var(--radius-lg);
    margin-bottom: 16px;
  }

  .providers-container {
    margin-bottom: 20px;
  }

  .empty-state {
    text-align: center;
    padding: 40px 20px;
    color: var(--color-text-secondary);
    font-size: 14px;
  }

  .add-provider-btn {
    width: 100%;
    padding: 12px;
    background: var(--color-primary);
    color: #fff;
    border: none;
    border-radius: var(--radius-lg);
    font-size: 14px;
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .add-provider-btn:hover {
    background: var(--color-primary-hover);
  }
</style>
