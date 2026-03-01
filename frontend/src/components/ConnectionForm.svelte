<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Connection } from '../lib/types';
  import { SaveConnection, TestConnection, DeleteConnection } from '../../wailsjs/go/main/App';
  import { addToast } from '../lib/stores/toasts';

  export let connection: Connection | null = null;

  const dispatch = createEventDispatcher();

  let form: Connection = connection ? { ...connection } : {
    id: '',
    name: '',
    host: '',
    port: 8333,
    accessKey: '',
    secretKey: '',
    region: 'us-east-1',
    bucket: '',
    useSSL: false,
    pathStyle: true,
  };

  let showSecret = false;
  let testing = false;
  let testSuccess: boolean | null = null;
  let saving = false;

  async function handleTest() {
    testing = true;
    testSuccess = null;
    try {
      await TestConnection(form);
      testSuccess = true;
      addToast('success', 'Connection successful!');
    } catch (e: any) {
      testSuccess = false;
      addToast('error', 'Connection failed: ' + e);
    } finally {
      testing = false;
    }
  }

  async function handleSave() {
    if (!form.name || !form.host || !form.accessKey || !form.secretKey || !form.bucket) {
      addToast('error', 'Please fill all required fields');
      return;
    }
    saving = true;
    try {
      await SaveConnection(form);
      addToast('success', connection ? 'Connection updated' : 'Connection saved');
      dispatch('save');
    } catch (e: any) {
      addToast('error', 'Failed to save: ' + e);
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!connection?.id) return;
    try {
      await DeleteConnection(connection.id);
      addToast('success', 'Connection deleted');
      dispatch('save');
    } catch (e: any) {
      addToast('error', 'Failed to delete: ' + e);
    }
  }

  function close() {
    dispatch('close');
  }
</script>

<!-- Backdrop — does NOT close on click -->
<div class="fixed inset-0 z-50 flex items-center justify-center" style="animation: fadeIn 0.15s ease-out">
  <div class="absolute inset-0 bg-black/75 backdrop-blur-lg"></div>

  <!-- Modal card -->
  <div class="relative w-[720px] max-h-[92vh] bg-surface-850 rounded-3xl shadow-2xl shadow-black/60 border border-surface-700/40 flex flex-col overflow-hidden ring-1 ring-surface-700/20"
    style="animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)">

    <!-- Header -->
    <div class="px-8 pt-8 pb-6 shrink-0">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-5">
          <div class="w-14 h-14 rounded-2xl bg-accent/15 flex items-center justify-center ring-1 ring-accent/20 shadow-lg shadow-accent/10">
            <svg class="w-7 h-7 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 12h14M12 5v14"/>
            </svg>
          </div>
          <div>
            <h2 class="text-[22px] font-bold text-surface-100">
              {connection ? 'Edit Connection' : 'New Connection'}
            </h2>
            <p class="text-[15px] text-surface-500 mt-1">Configure your S3-compatible storage endpoint</p>
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

    <!-- Divider -->
    <div class="mx-8 border-t border-surface-700/40"></div>

    <!-- Form body -->
    <div class="flex-1 overflow-y-auto px-8 py-6 space-y-6">

      <!-- Connection Name -->
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <div>
        <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Connection Name <span class="text-danger/60">*</span></label>
        <input bind:value={form.name} placeholder="My S3 Server"
          class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all" />
      </div>

      <!-- Host + Port -->
      <div class="grid grid-cols-4 gap-4">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <div class="col-span-3">
          <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Endpoint Host <span class="text-danger/60">*</span></label>
          <input bind:value={form.host} placeholder="s3.example.com"
            class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all" />
        </div>
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <div>
          <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Port</label>
          <input bind:value={form.port} type="number"
            class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all" />
        </div>
      </div>

      <!-- Divider -->
      <div class="border-t border-surface-700/30 !mt-8 !mb-8"></div>

      <!-- Access Key -->
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <div>
        <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Access Key <span class="text-danger/60">*</span></label>
        <input bind:value={form.accessKey} placeholder="AKIAIOSFODNN7EXAMPLE"
          class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all font-mono tracking-wide" />
      </div>

      <!-- Secret Key -->
      <!-- svelte-ignore a11y-label-has-associated-control -->
      <div>
        <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Secret Key <span class="text-danger/60">*</span></label>
        <div class="relative">
          {#if showSecret}
            <input bind:value={form.secretKey} type="text" placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLE"
              class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 pr-12 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all font-mono tracking-wide" />
          {:else}
            <input bind:value={form.secretKey} type="password" placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLE"
              class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 pr-12 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all font-mono tracking-wide" />
          {/if}
          <button class="absolute right-3 top-1/2 -translate-y-1/2 w-8 h-8 rounded-lg flex items-center justify-center text-surface-500 hover:text-surface-300 hover:bg-surface-700/50 transition-all"
            on:click={() => showSecret = !showSecret}
            title={showSecret ? 'Hide' : 'Show'}>
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {#if showSecret}
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"/>
              {:else}
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.8" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
              {/if}
            </svg>
          </button>
        </div>
      </div>

      <!-- Divider -->
      <div class="border-t border-surface-700/30 !mt-8 !mb-8"></div>

      <!-- Region + Bucket -->
      <div class="grid grid-cols-2 gap-4">
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <div>
          <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Region</label>
          <input bind:value={form.region} placeholder="us-east-1"
            class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all" />
        </div>
        <!-- svelte-ignore a11y-label-has-associated-control -->
        <div>
          <label class="block text-[13px] font-semibold text-surface-400 uppercase tracking-wider mb-2">Bucket <span class="text-danger/60">*</span></label>
          <input bind:value={form.bucket} placeholder="my-bucket"
            class="w-full bg-surface-900/80 border border-surface-700/60 rounded-xl px-5 py-3.5 text-[15px] text-surface-200 placeholder:text-surface-600 focus:outline-none focus:border-accent/50 focus:ring-2 focus:ring-accent/10 transition-all" />
        </div>
      </div>

      <!-- Toggles -->
      <div class="flex items-center gap-10 pt-3">
        <label class="flex items-center gap-3.5 cursor-pointer group">
          <div class="relative">
            <input type="checkbox" bind:checked={form.useSSL} class="sr-only peer" />
            <div class="w-12 h-7 rounded-full bg-surface-700 peer-checked:bg-accent transition-colors"></div>
            <div class="absolute left-0.5 top-0.5 w-[22px] h-[22px] rounded-full bg-surface-300 peer-checked:translate-x-5 peer-checked:bg-white transition-all shadow-sm"></div>
          </div>
          <span class="text-[15px] text-surface-400 group-hover:text-surface-300 transition-colors">Use SSL</span>
        </label>
        <label class="flex items-center gap-3.5 cursor-pointer group">
          <div class="relative">
            <input type="checkbox" bind:checked={form.pathStyle} class="sr-only peer" />
            <div class="w-12 h-7 rounded-full bg-surface-700 peer-checked:bg-accent transition-colors"></div>
            <div class="absolute left-0.5 top-0.5 w-[22px] h-[22px] rounded-full bg-surface-300 peer-checked:translate-x-5 peer-checked:bg-white transition-all shadow-sm"></div>
          </div>
          <span class="text-[15px] text-surface-400 group-hover:text-surface-300 transition-colors">Path Style</span>
        </label>
      </div>
    </div>

    <!-- Footer -->
    <div class="px-8 py-5 border-t border-surface-700/40 flex items-center gap-3 bg-surface-900/50 shrink-0">
      {#if connection?.id}
        <button on:click={handleDelete}
          class="px-6 py-3 text-[15px] font-medium text-danger/80 hover:text-danger hover:bg-danger/10 rounded-xl transition-all">
          Delete
        </button>
      {/if}
      <div class="flex-1"></div>

      <!-- Test button -->
      <button on:click={handleTest} disabled={testing}
        class="px-6 py-3 text-[15px] font-medium rounded-xl border transition-all disabled:opacity-50 flex items-center gap-2.5
          {testSuccess === true ? 'border-success/30 text-success bg-success/5' :
           testSuccess === false ? 'border-danger/30 text-danger bg-danger/5' :
           'border-surface-700/60 text-surface-300 hover:bg-surface-800 hover:border-surface-600'}">
        {#if testing}
          <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          Testing...
        {:else if testSuccess === true}
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
          </svg>
          Connected
        {:else if testSuccess === false}
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
          </svg>
          Failed
        {:else}
          Test Connection
        {/if}
      </button>

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
        {saving ? 'Saving...' : 'Save Connection'}
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
