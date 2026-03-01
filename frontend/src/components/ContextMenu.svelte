<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';

  export let x: number;
  export let y: number;
  export let items: Array<{ label: string; icon?: string; danger?: boolean; action: () => void }>;

  const dispatch = createEventDispatcher();

  let menuEl: HTMLDivElement;

  onMount(() => {
    const rect = menuEl.getBoundingClientRect();
    if (rect.right > window.innerWidth) x = window.innerWidth - rect.width - 8;
    if (rect.bottom > window.innerHeight) y = window.innerHeight - rect.height - 8;
  });
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="fixed inset-0 z-50" on:click={() => dispatch('close')} on:contextmenu|preventDefault={() => dispatch('close')}>
  <div
    bind:this={menuEl}
    class="absolute bg-surface-850 border border-surface-700/50 rounded-2xl shadow-2xl shadow-black/60 py-3 min-w-[240px] backdrop-blur-xl ring-1 ring-surface-700/20"
    style="left: {x}px; top: {y}px; animation: menuPop 0.15s cubic-bezier(0.34, 1.56, 0.64, 1)"
  >
    {#each items as item, i}
      {#if i > 0 && item.danger}
        <div class="mx-4 my-2 border-t border-surface-700/40"></div>
      {/if}
      <button
        class="w-[calc(100%-12px)] mx-1.5 px-5 py-3.5 text-[15px] text-left transition-all duration-150 flex items-center gap-3 rounded-xl
          {item.danger ? 'text-danger/80 hover:text-danger hover:bg-danger/12' : 'text-surface-300 hover:bg-surface-750/80 hover:text-surface-100'}"
        on:click|stopPropagation={() => { item.action(); dispatch('close'); }}
      >
        {item.label}
      </button>
    {/each}
  </div>
</div>

<style>
  @keyframes menuPop {
    from { opacity: 0; transform: scale(0.92) translateY(-4px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }
</style>
