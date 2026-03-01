<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import type { Settings } from '../lib/types';
  import { GetSettings, SaveSettings, ExportConnections, ImportConnections } from '../../wailsjs/go/main/App';
  import { addToast } from '../lib/stores/toasts';

  const dispatch = createEventDispatcher();

  let maxParallel = 3;
  let saving = false;
  let exporting = false;
  let importing = false;

  onMount(async () => {
    try {
      const s = await GetSettings();
      maxParallel = s.maxParallel;
    } catch (e: any) {
      addToast('error', 'Failed to load settings: ' + e);
    }
  });

  function clamp(v: number): number {
    return Math.max(1, Math.min(10, Math.round(v)));
  }

  async function handleSave() {
    saving = true;
    try {
      maxParallel = clamp(maxParallel);
      await SaveSettings({ maxParallel } as Settings);
      addToast('success', 'Settings saved');
      dispatch('close');
    } catch (e: any) {
      addToast('error', 'Failed to save settings: ' + e);
    } finally {
      saving = false;
    }
  }

  async function handleExport() {
    exporting = true;
    try {
      await ExportConnections();
      addToast('success', 'Connections exported');
    } catch (e: any) {
      if (e && String(e) !== '') addToast('error', 'Export failed: ' + e);
    } finally {
      exporting = false;
    }
  }

  async function handleImport() {
    importing = true;
    try {
      const added = await ImportConnections();
      if (added > 0) {
        addToast('success', `Imported ${added} connection(s)`);
        dispatch('imported');
      } else {
        addToast('info', 'No new connections imported (all duplicates or empty file)');
      }
    } catch (e: any) {
      if (e && String(e) !== '') addToast('error', 'Import failed: ' + e);
    } finally {
      importing = false;
    }
  }

  function close() {
    dispatch('close');
  }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center" style="animation: fadeIn 0.15s ease-out">
  <div class="absolute inset-0 bg-black/75 backdrop-blur-lg"></div>

  <div class="relative w-[520px] bg-surface-850 rounded-3xl shadow-2xl shadow-black/60 border border-surface-700/40 flex flex-col overflow-hidden ring-1 ring-surface-700/20"
    style="animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)">

    <!-- Header -->
    <div class="px-8 pt-8 pb-6 shrink-0">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-5">
          <div class="w-14 h-14 rounded-2xl bg-accent/15 flex items-center justify-center ring-1 ring-accent/20 shadow-lg shadow-accent/10">
            <svg class="w-7 h-7 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
          </div>
          <div>
            <h2 class="text-[22px] font-bold text-surface-100">Settings</h2>
            <p class="text-[15px] text-surface-500 mt-1">Configure app preferences</p>
          </div>
        </div>
        <button on:click={close}
          class="w-10 h-10 rounded-lg flex items-center justify-center text-surface-500 hover:text-surface-300 hover:bg-surface-700/50 transition-all">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </div>

    <div class="mx-8 border-t border-surface-700/40"></div>

    <!-- Body -->
    <div class="px-8 py-6 space-y-8">

      <!-- Transfers section -->
      <div>
        <div class="flex items-center gap-2.5 mb-4">
          <svg class="w-4.5 h-4.5 text-surface-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/>
          </svg>
          <span class="text-[13px] font-bold text-surface-400 uppercase tracking-[0.12em]">Transfers</span>
        </div>
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <div>
          <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Parallel Transfers</label>
          <div class="flex items-center gap-4">
            <input
              bind:value={maxParallel}
              type="number"
              min="1"
              max="10"
              on:blur={() => maxParallel = clamp(maxParallel)}
              class="w-24 bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all text-center"
            />
            <span class="text-[14px] text-surface-500">simultaneous uploads / downloads (1–10)</span>
          </div>
        </div>
      </div>

      <div class="border-t border-surface-700/30"></div>

      <!-- Profiles section -->
      <div>
        <div class="flex items-center gap-2.5 mb-4">
          <svg class="w-4.5 h-4.5 text-surface-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"/>
          </svg>
          <span class="text-[13px] font-bold text-surface-400 uppercase tracking-[0.12em]">Connection Profiles</span>
        </div>
        <div class="flex gap-3">
          <button
            on:click={handleExport}
            disabled={exporting}
            class="flex-1 py-3.5 px-5 text-[15px] font-medium rounded-xl border border-surface-700/60 text-surface-300 hover:bg-surface-800 hover:border-surface-600 transition-all disabled:opacity-50 flex items-center justify-center gap-2.5"
          >
            <svg class="w-4.5 h-4.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
            </svg>
            {exporting ? 'Exporting...' : 'Export'}
          </button>
          <button
            on:click={handleImport}
            disabled={importing}
            class="flex-1 py-3.5 px-5 text-[15px] font-medium rounded-xl border border-surface-700/60 text-surface-300 hover:bg-surface-800 hover:border-surface-600 transition-all disabled:opacity-50 flex items-center justify-center gap-2.5"
          >
            <svg class="w-4.5 h-4.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
            </svg>
            {importing ? 'Importing...' : 'Import'}
          </button>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="px-8 py-5 border-t border-surface-700/40 flex items-center justify-end gap-3 bg-surface-900/50 shrink-0">
      <button on:click={close}
        class="px-6 py-3 text-[15px] font-medium text-surface-400 hover:text-surface-200 hover:bg-surface-800 rounded-xl transition-all">
        Cancel
      </button>
      <button on:click={handleSave} disabled={saving}
        class="px-8 py-3 text-[15px] font-semibold rounded-xl transition-all disabled:opacity-50
          bg-gradient-to-r from-accent to-accent-dim
          hover:from-accent-hover hover:to-accent
          text-white shadow-lg shadow-accent/20 hover:shadow-xl hover:shadow-accent/30
          active:scale-[0.98]">
        {saving ? 'Saving...' : 'Save'}
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
    from { opacity: 0; transform: scale(0.95) translateY(8px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }
</style>
