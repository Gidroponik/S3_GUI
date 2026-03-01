<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { connections, activeConnectionId, showConnectionForm, editingConnection } from '../lib/stores/connections';
  import type { Connection } from '../lib/types';
  import logoSrc from '../assets/s3_gui.png';

  const dispatch = createEventDispatcher();

  function handleAdd() {
    editingConnection.set(null);
    showConnectionForm.set(true);
  }

  function handleConnect(conn: Connection) {
    if ($activeConnectionId === conn.id) {
      dispatch('disconnect');
    } else {
      dispatch('connect', conn.id);
    }
  }
</script>

<aside class="w-[320px] bg-surface-950 border-r border-surface-750/50 flex flex-col shrink-0">
  <!-- Header (draggable title bar) -->
  <div class="h-[76px] flex items-center px-6 border-b border-surface-750/50 shrink-0 gap-4" style="--wails-draggable: drag">
    <img src={logoSrc} alt="S3 GUI" class="w-11 h-11 rounded-[14px] shrink-0 shadow-lg shadow-accent/10" />
    <div class="flex-1">
      <span class="text-[18px] font-bold text-surface-100 tracking-tight">S3 Bucket</span>
      <span class="text-[18px] font-normal text-surface-500 tracking-tight ml-1.5">Manager</span>
    </div>
    <button
      class="w-9 h-9 rounded-xl flex items-center justify-center text-surface-500 hover:text-surface-200 hover:bg-surface-700/50 transition-all shrink-0"
      style="--wails-draggable: none"
      on:click={() => dispatch('settings')}
      title="Settings"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
        <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
      </svg>
    </button>
  </div>

  <!-- Connections label -->
  <div class="px-6 pt-8 pb-3">
    <span class="text-[12px] font-bold text-surface-500 uppercase tracking-[0.12em]">Connections</span>
  </div>

  <!-- Connection list -->
  <div class="flex-1 overflow-y-auto px-4 space-y-2">
    {#each $connections as conn (conn.id)}
      <button
        class="w-full px-5 py-4 flex items-center gap-4 text-left rounded-2xl transition-all duration-200 group
          {$activeConnectionId === conn.id
            ? 'bg-accent/12 ring-1 ring-accent/25 shadow-md shadow-accent/8'
            : 'bg-surface-900/50 hover:bg-surface-850 ring-1 ring-surface-750/30 hover:ring-surface-700/50'}"
        on:click={() => handleConnect(conn)}
      >
        <div class="w-3.5 h-3.5 rounded-full shrink-0 transition-all duration-300
          {$activeConnectionId === conn.id ? 'bg-success animate-pulse-glow' : 'bg-surface-600 ring-1 ring-surface-500/40'}">
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-[15px] font-semibold truncate transition-colors
            {$activeConnectionId === conn.id ? 'text-surface-100' : 'text-surface-200 group-hover:text-surface-100'}">{conn.name}</div>
          <div class="text-[13px] text-surface-500 truncate mt-1.5">{conn.host}:{conn.port}/{conn.bucket}</div>
        </div>
        <button
          class="opacity-0 group-hover:opacity-100 text-surface-500 hover:text-surface-200 p-2 rounded-xl hover:bg-surface-700/60 transition-all"
          on:click|stopPropagation={() => dispatch('edit', conn)}
          title="Edit"
        >
          <svg class="w-4.5 h-4.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
          </svg>
        </button>
      </button>
    {/each}

    {#if $connections.length === 0}
      <div class="px-4 py-16 text-center">
        <div class="w-16 h-16 rounded-2xl bg-surface-800/50 flex items-center justify-center mx-auto mb-5 ring-1 ring-surface-700/40">
          <svg class="w-8 h-8 text-surface-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 12h14M12 5v14"/>
          </svg>
        </div>
        <p class="text-[15px] text-surface-400 font-medium mb-1.5">No connections yet</p>
        <p class="text-[14px] text-surface-600 leading-relaxed">Add one to get started</p>
      </div>
    {/if}
  </div>

  <!-- New Connection button -->
  <div class="p-5 shrink-0">
    <button
      class="w-full py-4 px-5 rounded-2xl text-[16px] font-bold transition-all duration-200 flex items-center justify-center gap-3 bg-gradient-to-r from-accent to-accent-dim hover:from-accent-hover hover:to-accent text-white shadow-xl shadow-accent/25 hover:shadow-2xl hover:shadow-accent/35 active:scale-[0.97] hover:-translate-y-0.5"
      on:click={handleAdd}
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
      </svg>
      New Connection
    </button>
  </div>
</aside>
