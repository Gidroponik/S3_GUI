<script lang="ts">
  import { toasts } from '../lib/stores/toasts';
</script>

{#if $toasts.length > 0}
  <div class="fixed bottom-6 right-6 z-[100] flex flex-col gap-4 pointer-events-none">
    {#each $toasts as toast (toast.id)}
      <div class="pointer-events-auto flex items-center gap-3 px-6 py-5 rounded-xl shadow-xl text-[15px] backdrop-blur-md
        {toast.type === 'success' ? 'bg-success/10 text-success border border-success/15 shadow-success/5' :
         toast.type === 'error' ? 'bg-danger/10 text-danger border border-danger/15 shadow-danger/5' :
         'bg-accent/10 text-accent border border-accent/15 shadow-accent/5'}"
        style="animation: toastSlide 0.25s cubic-bezier(0.16, 1, 0.3, 1)">
        {#if toast.type === 'success'}
          <svg class="w-5.5 h-5.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
          </svg>
        {:else if toast.type === 'error'}
          <svg class="w-5.5 h-5.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
        {:else}
          <svg class="w-5.5 h-5.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        {/if}
        <span class="font-medium">{toast.message}</span>
      </div>
    {/each}
  </div>
{/if}

<style>
  @keyframes toastSlide {
    from {
      transform: translateX(120%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }
</style>
