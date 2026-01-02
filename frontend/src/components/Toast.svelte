<script lang="ts">
  type Props = {
    message: string
    type: 'info' | 'error' | 'success'
    visible: boolean
  }

  let { message, type, visible }: Props = $props()
</script>

{#if visible}
  <div class="toast {type}">
    {message}
  </div>
{/if}

<style>
  .toast {
    position: fixed;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    padding: 12px 24px;
    background: #333;
    color: #fff;
    border-radius: var(--radius-lg);
    font-size: 14px;
    z-index: 10000;
    animation: slideUp 0.3s ease;
  }

  .toast.error {
    background: rgba(255, 59, 48, 0.9);
  }

  .toast.success {
    background: rgba(52, 199, 89, 0.9);
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateX(-50%) translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateX(-50%) translateY(0);
    }
  }
</style>
