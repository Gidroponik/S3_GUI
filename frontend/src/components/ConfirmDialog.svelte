<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let title: string;
  export let message: string;

  const dispatch = createEventDispatcher();
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center" style="animation: fadeIn 0.15s ease-out">
  <div class="absolute inset-0 bg-black/75 backdrop-blur-lg"></div>

  <div class="relative bg-surface-850 border border-surface-700/40 rounded-3xl shadow-2xl shadow-black/60 w-[520px] p-10 ring-1 ring-surface-700/20"
    style="animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)">
    <div class="w-16 h-16 rounded-2xl bg-danger/12 flex items-center justify-center mb-6 ring-1 ring-danger/15">
      <svg class="w-8 h-8 text-danger" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
      </svg>
    </div>
    <h3 class="text-[22px] font-bold text-surface-100 mb-3">{title}</h3>
    <p class="text-[16px] text-surface-400 mb-8 leading-relaxed whitespace-pre-line">{message}</p>
    <div class="flex justify-end gap-3">
      <button class="px-8 py-3.5 text-[15px] font-medium text-surface-400 hover:text-surface-100 hover:bg-surface-800 rounded-2xl transition-all"
        on:click={() => dispatch('cancel')}>
        Cancel
      </button>
      <button class="px-8 py-3.5 text-[15px] font-bold bg-danger hover:bg-danger-hover text-white rounded-2xl shadow-xl shadow-danger/25 hover:shadow-2xl hover:shadow-danger/35 transition-all active:scale-[0.97]"
        on:click={() => dispatch('confirm')}>
        Delete
      </button>
    </div>
  </div>
</div>

<style>
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  @keyframes scaleIn {
    from { opacity: 0; transform: scale(0.92) translateY(12px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }
</style>
