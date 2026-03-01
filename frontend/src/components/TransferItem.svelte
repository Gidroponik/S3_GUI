<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Transfer } from '../lib/types';
  import { formatBytes } from '../lib/utils';

  export let transfer: Transfer;

  const dispatch = createEventDispatcher();

  $: isActive = transfer.status === 'pending' || transfer.status === 'in_progress';
  $: statusColor = {
    pending: 'text-surface-500',
    in_progress: 'text-accent',
    completed: 'text-success',
    failed: 'text-danger',
    cancelled: 'text-surface-500',
  }[transfer.status];
</script>

<div class="flex items-center gap-4 px-7 py-5 text-[15px] border-t border-surface-800/40">
  <div class="shrink-0 {statusColor}">
    {#if transfer.type === 'upload'}
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
      </svg>
    {:else}
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
      </svg>
    {/if}
  </div>

  <div class="flex-1 min-w-0">
    <div class="flex items-center justify-between mb-2">
      <span class="truncate text-surface-200 text-[15px] font-medium">{transfer.fileName}</span>
      <span class="shrink-0 ml-3 text-[14px] font-bold {statusColor}">
        {#if transfer.status === 'in_progress'}
          {transfer.percentage.toFixed(0)}%
        {:else if transfer.status === 'completed'}
          Done
        {:else if transfer.status === 'failed'}
          Failed
        {:else if transfer.status === 'cancelled'}
          Cancelled
        {:else}
          Queued
        {/if}
      </span>
    </div>

    {#if isActive}
      <div class="w-full h-3 bg-surface-800 rounded-full overflow-hidden">
        <div class="h-full progress-gradient rounded-full transition-all duration-300 ease-out"
          style="width: {transfer.percentage}%"></div>
      </div>
    {/if}

    {#if transfer.bytesTotal > 0 && isActive}
      <div class="text-[13px] text-surface-500 mt-2 font-medium">
        {formatBytes(transfer.bytesDone)} / {formatBytes(transfer.bytesTotal)}
      </div>
    {/if}

    {#if transfer.error}
      <div class="text-[13px] text-danger/80 mt-1.5 truncate" title={transfer.error}>{transfer.error}</div>
    {/if}
  </div>

  {#if isActive}
    <button class="shrink-0 p-2.5 rounded-xl text-surface-600 hover:text-surface-200 hover:bg-surface-700/60 transition-all"
      on:click={() => dispatch('cancel', transfer.id)} title="Cancel">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </button>
  {/if}
</div>
