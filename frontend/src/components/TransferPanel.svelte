<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import TransferItem from './TransferItem.svelte';
  import { transfers, transferPanelOpen } from '../lib/stores/transfers';

  const dispatch = createEventDispatcher();

  $: activeCount = $transfers.filter(t => t.status === 'pending' || t.status === 'in_progress').length;
  $: hasTransfers = $transfers.length > 0;

  $: if (activeCount > 0 && !$transferPanelOpen) {
    transferPanelOpen.set(true);
  }
</script>

{#if hasTransfers}
  <div class="border-t border-surface-800/80 bg-surface-950 shrink-0">
    <button class="w-full flex items-center justify-between px-6 py-4 hover:bg-surface-800/30 transition-colors"
      on:click={() => transferPanelOpen.update(v => !v)}>
      <div class="flex items-center gap-3 text-[14px] font-medium text-surface-400">
        <svg class="w-5.5 h-5.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
          <path stroke-linecap="round" stroke-linejoin="round" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/>
        </svg>
        <span>Transfers</span>
        {#if activeCount > 0}
          <span class="bg-accent/20 text-accent px-3 py-0.5 rounded-full text-[12px] font-bold">{activeCount}</span>
        {/if}
      </div>
      <div class="flex items-center gap-4">
        {#if $transfers.some(t => t.status === 'completed' || t.status === 'failed' || t.status === 'cancelled')}
          <button class="text-[12px] text-surface-500 hover:text-surface-300 font-medium px-3 py-1.5 rounded-lg hover:bg-surface-800/50 transition-all"
            on:click|stopPropagation={() => dispatch('clear')}>
            Clear
          </button>
        {/if}
        <svg class="w-4 h-4 text-surface-600 transition-transform duration-200 {$transferPanelOpen ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
        </svg>
      </div>
    </button>

    {#if $transferPanelOpen}
      <div class="max-h-56 overflow-y-auto">
        {#each $transfers as transfer (transfer.id)}
          <TransferItem {transfer} on:cancel />
        {/each}
      </div>
    {/if}
  </div>
{/if}
