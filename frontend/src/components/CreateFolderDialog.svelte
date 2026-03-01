<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';

  export let parentPath: string = '';

  const dispatch = createEventDispatcher();

  let folderName = '';
  let inputEl: HTMLInputElement;

  onMount(() => {
    inputEl?.focus();
  });

  function handleCreate() {
    if (folderName.trim()) {
      dispatch('create', folderName.trim());
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Enter') handleCreate();
    if (e.key === 'Escape') dispatch('cancel');
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="fixed inset-0 z-50 flex items-center justify-center" style="animation: fadeIn 0.15s ease-out">
  <div class="absolute inset-0 bg-black/75 backdrop-blur-lg"></div>

  <div class="relative bg-surface-850 border border-surface-700/40 rounded-3xl shadow-2xl shadow-black/60 w-[520px] p-10 ring-1 ring-surface-700/20"
    style="animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)">
    <div class="w-14 h-14 rounded-2xl bg-accent/12 flex items-center justify-center mb-6 ring-1 ring-accent/20">
      <svg class="w-7 h-7 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"/>
      </svg>
    </div>
    <h3 class="text-[22px] font-bold text-surface-100 mb-2">Create Folder</h3>
    {#if parentPath}
      <p class="text-[14px] text-surface-500 mb-6">Inside <span class="text-surface-300 font-medium">{parentPath}</span></p>
    {:else}
      <p class="text-[14px] text-surface-500 mb-6">In the root directory</p>
    {/if}

    <!-- svelte-ignore a11y-label-has-associated-control -->
    <div class="mb-8">
      <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Folder Name</label>
      <input
        bind:this={inputEl}
        bind:value={folderName}
        placeholder="my-folder"
        on:keydown={handleKeydown}
        class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all"
      />
    </div>

    <div class="flex justify-end gap-3">
      <button class="px-8 py-3.5 text-[15px] font-medium text-surface-400 hover:text-surface-100 hover:bg-surface-800 rounded-2xl transition-all"
        on:click={() => dispatch('cancel')}>
        Cancel
      </button>
      <button
        class="px-8 py-3.5 text-[15px] font-bold rounded-2xl transition-all disabled:opacity-40 bg-gradient-to-r from-accent to-accent-dim hover:from-accent-hover hover:to-accent text-white shadow-xl shadow-accent/25 hover:shadow-2xl hover:shadow-accent/35 active:scale-[0.97]"
        disabled={!folderName.trim()}
        on:click={handleCreate}>
        Create
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
