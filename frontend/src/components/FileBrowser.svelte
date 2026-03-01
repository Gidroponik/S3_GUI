<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import FileRow from './FileRow.svelte';
  import { filteredObjects, selectedKeys, loading, sortField, sortDir, objects } from '../lib/stores/browser';
  import type { S3Object, SortField } from '../lib/types';

  export let isDragOver = false;

  const dispatch = createEventDispatcher();

  function toggleSort(field: SortField) {
    if ($sortField === field) {
      sortDir.update(d => d === 'asc' ? 'desc' : 'asc');
    } else {
      sortField.set(field);
      sortDir.set('asc');
    }
  }

  function handleSelect(e: CustomEvent<{ key: string; shiftKey: boolean; ctrlKey: boolean }>) {
    const { key, shiftKey, ctrlKey } = e.detail;
    selectedKeys.update(keys => {
      const next = new Set(keys);
      if (ctrlKey) {
        if (next.has(key)) next.delete(key);
        else next.add(key);
      } else if (shiftKey && keys.size > 0) {
        const allKeys = $filteredObjects.map(o => o.key);
        const lastSelected = Array.from(keys).pop()!;
        const lastIdx = allKeys.indexOf(lastSelected);
        const curIdx = allKeys.indexOf(key);
        const [start, end] = lastIdx < curIdx ? [lastIdx, curIdx] : [curIdx, lastIdx];
        for (let i = start; i <= end; i++) {
          next.add(allKeys[i]);
        }
      } else {
        next.clear();
        next.add(key);
      }
      return next;
    });
  }

  function handleSelectAll() {
    if ($selectedKeys.size === $filteredObjects.length) {
      selectedKeys.set(new Set());
    } else {
      selectedKeys.set(new Set($filteredObjects.map(o => o.key)));
    }
  }

  function sortIcon(field: SortField): string {
    if ($sortField !== field) return '';
    return $sortDir === 'asc' ? ' ↑' : ' ↓';
  }

  function handleBgContextMenu(e: MouseEvent) {
    e.preventDefault();
    const target = e.target as HTMLElement;
    // Only fire on the browser background, not on file rows
    if (target.closest('[role="row"]')) return;
    dispatch('bgContextmenu', { x: e.clientX, y: e.clientY });
  }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="flex-1 overflow-auto relative" style="--wails-drop-target: drop"
  on:contextmenu={handleBgContextMenu}>
  {#if $loading && $objects.length === 0}
    <!-- Loading skeleton -->
    <div class="p-6 space-y-0.5">
      {#each Array(10) as _, i}
        <div class="flex items-center gap-5 py-5 animate-pulse border-b border-surface-800/30" style="opacity: {1 - i * 0.08}">
          <div class="w-5.5 h-5.5 bg-surface-800/50 rounded ml-5"></div>
          <div class="w-8 h-8 bg-surface-800/50 rounded-lg"></div>
          <div class="h-4.5 bg-surface-800/50 rounded-lg" style="width: {150 + Math.random() * 170}px"></div>
          <div class="flex-1"></div>
          <div class="h-4.5 bg-surface-800/50 rounded-lg w-24"></div>
          <div class="h-4.5 bg-surface-800/50 rounded-lg w-36"></div>
        </div>
      {/each}
    </div>
  {:else if $filteredObjects.length === 0}
    <div class="flex-1 flex items-center justify-center h-full">
      <div class="text-center py-24">
        <div class="w-28 h-28 rounded-[28px] bg-surface-850 flex items-center justify-center mx-auto mb-6 ring-1 ring-surface-750/50 shadow-xl shadow-surface-950/50">
          <svg class="w-14 h-14 text-surface-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
          </svg>
        </div>
        <p class="text-[20px] font-semibold text-surface-300 mb-2">This folder is empty</p>
        <p class="text-[15px] text-surface-500 mb-8">Upload files or create a subfolder to get started</p>
        <div class="flex items-center gap-4">
          <button class="px-8 py-4 text-[15px] font-bold rounded-2xl transition-all duration-200 inline-flex items-center gap-3 bg-gradient-to-r from-accent to-accent-dim hover:from-accent-hover hover:to-accent text-white shadow-xl shadow-accent/25 hover:shadow-2xl hover:shadow-accent/40 active:scale-[0.97] hover:-translate-y-0.5"
            on:click={() => dispatch('upload')}>
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
            </svg>
            Upload Files
          </button>
          <button class="px-8 py-4 text-[15px] font-bold rounded-2xl transition-all duration-200 inline-flex items-center gap-3 bg-surface-800 hover:bg-surface-750 text-surface-200 ring-1 ring-surface-700/50 hover:ring-surface-600/50 active:scale-[0.97]"
            on:click={() => dispatch('createFolder')}>
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
            </svg>
            Create Folder
          </button>
        </div>
      </div>
    </div>
  {:else}
    <!-- Header -->
    <div class="sticky top-0 bg-surface-900/95 backdrop-blur-xl border-b border-surface-750/50 z-10">
      <div class="flex items-center px-7 py-4 text-[13px] text-surface-500 font-bold uppercase tracking-[0.08em]">
        <div class="w-12 shrink-0 flex items-center">
          <input type="checkbox"
            checked={$selectedKeys.size === $filteredObjects.length && $filteredObjects.length > 0}
            on:change={handleSelectAll}
            class="rounded border-surface-600 bg-surface-800 text-accent" />
        </div>
        <div class="w-11 shrink-0"></div>
        <button class="flex-1 text-left hover:text-surface-200 transition-colors py-1.5 ml-1" on:click={() => toggleSort('name')}>
          Name{sortIcon('name')}
        </button>
        <button class="w-32 text-right hover:text-surface-200 transition-colors py-1.5" on:click={() => toggleSort('size')}>
          Size{sortIcon('size')}
        </button>
        <button class="w-52 text-right hover:text-surface-200 transition-colors py-1.5" on:click={() => toggleSort('lastModified')}>
          Modified{sortIcon('lastModified')}
        </button>
        <div class="w-12"></div>
      </div>
    </div>

    <!-- Rows -->
    <div>
      {#each $filteredObjects as obj (obj.key)}
        <FileRow
          {obj}
          selected={$selectedKeys.has(obj.key)}
          on:select={handleSelect}
          on:navigate
          on:contextmenu
          on:download
        />
      {/each}
    </div>
  {/if}

  <!-- Drop zone overlay -->
  {#if isDragOver}
    <div class="absolute inset-0 z-20 flex items-center justify-center bg-accent/8 border-[3px] border-dashed border-accent/50 rounded-2xl m-3 backdrop-blur-sm pointer-events-none">
      <div class="text-center">
        <div class="w-20 h-20 rounded-3xl bg-accent/15 flex items-center justify-center mx-auto mb-5 ring-1 ring-accent/30">
          <svg class="w-10 h-10 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
          </svg>
        </div>
        <p class="text-[18px] font-bold text-accent mb-2">Drop files here</p>
        <p class="text-[14px] text-surface-400">Files will be uploaded to the current folder</p>
      </div>
    </div>
  {/if}
</div>
