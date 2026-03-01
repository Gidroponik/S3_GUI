<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { currentPrefix, searchQuery, selectedKeys, loading } from '../lib/stores/browser';
  import { breadcrumbs } from '../lib/stores/browser';

  const dispatch = createEventDispatcher();

</script>

<div class="h-[72px] bg-surface-900 border-b border-surface-750/50 flex items-center px-7 gap-4 shrink-0">
  <!-- Breadcrumbs -->
  <div class="flex items-center gap-2 min-w-0 flex-1">
    <button class="text-surface-500 hover:text-accent transition-colors shrink-0 p-3 rounded-xl hover:bg-surface-800/80"
      on:click={() => dispatch('navigateTo', '')} title="Root">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
      </svg>
    </button>
    {#each $breadcrumbs as crumb, i}
      <svg class="w-4.5 h-4.5 text-surface-600 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
      </svg>
      <button class="text-[15px] font-medium text-surface-400 hover:text-accent transition-colors truncate max-w-[180px] px-2.5 py-1.5 rounded-xl hover:bg-surface-800/80"
        on:click={() => dispatch('navigateTo', crumb.prefix)}>
        {crumb.name}
      </button>
    {/each}

    {#if $currentPrefix}
      <button class="ml-2 p-3 rounded-xl text-surface-500 hover:text-surface-200 hover:bg-surface-800/80 transition-all shrink-0" on:click={() => dispatch('navigateUp')} title="Go up (Backspace)">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
        </svg>
      </button>
    {/if}
  </div>

  <!-- Search -->
  <div class="relative w-72">
    <svg class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-surface-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
    </svg>
    <input bind:value={$searchQuery} placeholder="Search files..."
      class="w-full bg-surface-850 border border-surface-700/50 rounded-2xl pl-12 pr-4 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/15 transition-all" />
  </div>

  <!-- Bulk actions -->
  {#if $selectedKeys.size > 0}
    <div class="flex items-center gap-2.5 text-[15px] text-surface-400 border-l border-surface-700/50 pl-5 ml-1">
      <span class="text-accent font-bold text-[16px]">{$selectedKeys.size}</span>
      <span>selected</span>
      <button class="p-3 hover:bg-surface-800 rounded-xl text-surface-400 hover:text-surface-200 transition-all"
        on:click={() => dispatch('downloadSelected')} title="Download selected">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
        </svg>
      </button>
      <button class="p-3 hover:bg-danger/15 rounded-xl text-surface-400 hover:text-danger transition-all"
        on:click={() => dispatch('deleteSelected')} title="Delete selected">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
        </svg>
      </button>
    </div>
  {/if}

  <!-- Separator -->
  <div class="border-l border-surface-700/40 h-8 mx-1"></div>

  <!-- Actions -->
  <button class="p-3 hover:bg-surface-800 rounded-xl text-surface-500 hover:text-surface-200 transition-all"
    on:click={() => dispatch('refresh')} disabled={$loading} title="Refresh (Ctrl+R)">
    <svg class="w-6 h-6 {$loading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
    </svg>
  </button>

  <button class="p-3 hover:bg-surface-800 rounded-xl text-surface-500 hover:text-surface-200 transition-all"
    on:click={() => dispatch('createFolder')} title="New Folder">
    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
      <path stroke-linecap="round" stroke-linejoin="round" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
    </svg>
  </button>

  <button class="p-3 hover:bg-surface-800 rounded-xl text-surface-500 hover:text-surface-200 transition-all"
    on:click={() => dispatch('uploadFolder')} title="Upload Folder">
    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
      <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
    </svg>
  </button>

  <button class="px-7 py-3.5 text-[15px] font-bold rounded-2xl transition-all duration-200 flex items-center gap-3
    bg-gradient-to-r from-accent to-accent-dim hover:from-accent-hover hover:to-accent text-white
    shadow-lg shadow-accent/25 hover:shadow-xl hover:shadow-accent/40 active:scale-[0.97] hover:-translate-y-0.5"
    on:click={() => dispatch('upload')}>
    <svg class="w-5.5 h-5.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
      <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
    </svg>
    Upload
  </button>
</div>
