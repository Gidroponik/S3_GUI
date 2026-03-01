<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { S3Object } from '../lib/types';
  import { formatBytes, formatDate, getFileIcon } from '../lib/utils';

  export let obj: S3Object;
  export let selected = false;

  const dispatch = createEventDispatcher();

  function handleClick(e: MouseEvent) {
    dispatch('select', { key: obj.key, shiftKey: e.shiftKey, ctrlKey: e.ctrlKey || e.metaKey });
  }

  function handleDblClick() {
    if (obj.isFolder) {
      dispatch('navigate', obj.key);
    } else {
      dispatch('download', obj.key);
    }
  }

  function handleContextMenu(e: MouseEvent) {
    e.preventDefault();
    dispatch('contextmenu', { x: e.clientX, y: e.clientY, object: obj });
  }

  const iconType = getFileIcon(obj.name, obj.isFolder);
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
  class="flex items-center px-7 py-[18px] cursor-pointer transition-all duration-150 group border-b border-surface-800/40
    {selected
      ? 'bg-accent/10 border-l-[4px] border-l-accent/70 pl-[24px]'
      : 'hover:bg-surface-850/80 border-l-[4px] border-l-transparent'}"
  on:click={handleClick}
  on:dblclick={handleDblClick}
  on:contextmenu={handleContextMenu}
  role="row"
  tabindex="0"
>
  <div class="w-12 shrink-0 flex items-center">
    <input type="checkbox" checked={selected}
      on:click|stopPropagation
      on:change={() => dispatch('select', { key: obj.key, shiftKey: false, ctrlKey: true })}
      class="rounded" />
  </div>

  <div class="w-11 shrink-0 flex items-center justify-center">
    {#if obj.isFolder}
      <svg class="w-8 h-8 text-accent drop-shadow-[0_0_6px_rgba(124,92,252,0.3)]" fill="currentColor" viewBox="0 0 24 24">
        <path d="M10 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/>
      </svg>
    {:else if iconType === 'image'}
      <svg class="w-7 h-7 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
      </svg>
    {:else if iconType === 'video'}
      <svg class="w-7 h-7 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"/>
      </svg>
    {:else if iconType === 'music'}
      <svg class="w-7 h-7 text-pink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"/>
      </svg>
    {:else if iconType === 'archive'}
      <svg class="w-7 h-7 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"/>
      </svg>
    {:else if iconType === 'file-code'}
      <svg class="w-7 h-7 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
      </svg>
    {:else if iconType === 'file-text'}
      <svg class="w-7 h-7 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
      </svg>
    {:else}
      <svg class="w-7 h-7 text-surface-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.6">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
      </svg>
    {/if}
  </div>

  <div class="flex-1 min-w-0 truncate text-[16px] ml-1
    {obj.isFolder ? 'font-semibold text-surface-100' : 'text-surface-300'}">
    {obj.name}
  </div>

  <div class="w-32 text-right text-surface-500 text-[14px] font-medium">
    {#if !obj.isFolder}
      {formatBytes(obj.size)}
    {/if}
  </div>

  <div class="w-52 text-right text-surface-500 text-[14px]">
    {#if !obj.isFolder && obj.lastModified}
      {formatDate(obj.lastModified)}
    {/if}
  </div>

  <div class="w-12 flex justify-end">
    <button class="p-2.5 rounded-xl text-surface-600 hover:text-surface-200 hover:bg-surface-700/60 opacity-0 group-hover:opacity-100 transition-all"
      on:click|stopPropagation={(e) => dispatch('contextmenu', { x: e.clientX, y: e.clientY, object: obj })}>
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
        <path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"/>
      </svg>
    </button>
  </div>
</div>
