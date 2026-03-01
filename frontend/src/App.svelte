<script lang="ts">
  import { onMount } from 'svelte';
  import Sidebar from './components/Sidebar.svelte';
  import TopBar from './components/TopBar.svelte';
  import FileBrowser from './components/FileBrowser.svelte';
  import TransferPanel from './components/TransferPanel.svelte';
  import ConnectionForm from './components/ConnectionForm.svelte';
  import ConfirmDialog from './components/ConfirmDialog.svelte';
  import CreateFolderDialog from './components/CreateFolderDialog.svelte';
  import ContextMenu from './components/ContextMenu.svelte';
  import ToastContainer from './components/ToastContainer.svelte';
  import SettingsDialog from './components/SettingsDialog.svelte';
  import { connections, activeConnectionId, showConnectionForm, editingConnection } from './lib/stores/connections';
  import { objects, currentPrefix, loading, selectedKeys } from './lib/stores/browser';
  import { updateTransfer } from './lib/stores/transfers';
  import { addToast } from './lib/stores/toasts';
  import type { Connection, S3Object } from './lib/types';
  import { GetConnections, Connect, Disconnect, ListObjects, DeleteObjects, DeletePrefix, CreateFolder, UploadFiles, DownloadFiles, UploadFolder, DownloadPrefix, CancelTransfer, ClearCompletedTransfers, GetTransfers, UploadDroppedFiles, GetPresignedURL, GetDirectURL } from '../wailsjs/go/main/App';
  import { EventsOn, OnFileDrop, OnFileDropOff } from '../wailsjs/runtime/runtime';

  let connected = false;
  let showSettings = false;
  let isDragOver = false;
  let confirmDialog: { show: boolean; title: string; message: string; onConfirm: () => void } = {
    show: false, title: '', message: '', onConfirm: () => {}
  };
  let contextMenu: { show: boolean; x: number; y: number; items: Array<{ label: string; icon?: string; danger?: boolean; action: () => void }> } = {
    show: false, x: 0, y: 0, items: []
  };

  let createFolderDialog: { show: boolean; parentPath: string } = {
    show: false, parentPath: ''
  };

  let dragCounter = 0;

  onMount(async () => {
    await loadConnections();

    EventsOn('transfer:progress', (t: any) => {
      updateTransfer(t);
    });
    EventsOn('transfer:complete', (t: any) => {
      updateTransfer(t);
      addToast('success', `${t.type === 'upload' ? 'Upload' : 'Download'} complete: ${t.fileName}`);
      if (connected) refresh();
    });
    EventsOn('transfer:cleared', () => {
      ClearCompletedTransfers();
    });

    // Drag-and-drop file handling
    OnFileDrop((x: number, y: number, paths: string[]) => {
      isDragOver = false;
      dragCounter = 0;
      if (!connected || !paths || paths.length === 0) return;
      handleDroppedFiles(paths);
    }, true);

    document.addEventListener('keydown', handleKeydown);
    document.addEventListener('dragenter', handleDragEnter);
    document.addEventListener('dragleave', handleDragLeave);
    document.addEventListener('dragover', handleDragOver);
    document.addEventListener('drop', handleDrop);

    return () => {
      document.removeEventListener('keydown', handleKeydown);
      document.removeEventListener('dragenter', handleDragEnter);
      document.removeEventListener('dragleave', handleDragLeave);
      document.removeEventListener('dragover', handleDragOver);
      document.removeEventListener('drop', handleDrop);
      OnFileDropOff();
    };
  });

  function handleDragEnter(e: DragEvent) {
    e.preventDefault();
    dragCounter++;
    if (connected) isDragOver = true;
  }

  function handleDragLeave(e: DragEvent) {
    e.preventDefault();
    dragCounter--;
    if (dragCounter <= 0) {
      isDragOver = false;
      dragCounter = 0;
    }
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragOver = false;
    dragCounter = 0;
  }

  async function handleDroppedFiles(paths: string[]) {
    try {
      await UploadDroppedFiles(paths, $currentPrefix);
      addToast('info', `Uploading ${paths.length} item(s)...`);
    } catch (e: any) {
      addToast('error', 'Drop upload failed: ' + e);
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      contextMenu.show = false;
      // Do NOT close the connection form on Escape — only close via X or Cancel
      confirmDialog.show = false;
      createFolderDialog.show = false;
      showSettings = false;
    }
    if (!connected) return;
    if (e.ctrlKey && e.key === 'r') { e.preventDefault(); refresh(); }
    if (e.ctrlKey && e.key === 'a') { e.preventDefault(); selectAll(); }
    if (e.key === 'Delete' && $selectedKeys.size > 0) { deleteSelected(); }
    if (e.key === 'Backspace' && $currentPrefix) { navigateUp(); }
  }

  async function loadConnections() {
    try {
      const conns = await GetConnections();
      connections.set(conns || []);
    } catch (e: any) {
      addToast('error', 'Failed to load connections: ' + e);
    }
  }

  async function handleConnect(id: string) {
    try {
      loading.set(true);
      await Connect(id);
      activeConnectionId.set(id);
      connected = true;
      currentPrefix.set('');
      selectedKeys.set(new Set());
      await refresh();
      addToast('success', 'Connected');
    } catch (e: any) {
      addToast('error', 'Connection failed: ' + e);
    } finally {
      loading.set(false);
    }
  }

  async function handleDisconnect() {
    await Disconnect();
    activeConnectionId.set(null);
    connected = false;
    objects.set([]);
    currentPrefix.set('');
    selectedKeys.set(new Set());
  }

  async function refresh() {
    if (!connected) return;
    try {
      loading.set(true);
      const result = await ListObjects($currentPrefix);
      objects.set(result || []);
      selectedKeys.set(new Set());
    } catch (e: any) {
      addToast('error', 'Failed to list objects: ' + e);
    } finally {
      loading.set(false);
    }
  }

  async function navigateTo(prefix: string) {
    currentPrefix.set(prefix);
    await refresh();
  }

  function navigateUp() {
    const parts = $currentPrefix.split('/').filter(Boolean);
    parts.pop();
    navigateTo(parts.length ? parts.join('/') + '/' : '');
  }

  function selectAll() {
    const allKeys = new Set($objects.map(o => o.key));
    selectedKeys.set(allKeys);
  }

  async function handleUpload() {
    try {
      await UploadFiles($currentPrefix);
    } catch (e: any) {
      addToast('error', 'Upload failed: ' + e);
    }
  }

  async function handleUploadFolder() {
    try {
      await UploadFolder($currentPrefix);
    } catch (e: any) {
      addToast('error', 'Folder upload failed: ' + e);
    }
  }

  async function handleDownload(keys: string[]) {
    try {
      await DownloadFiles(keys);
    } catch (e: any) {
      addToast('error', 'Download failed: ' + e);
    }
  }

  async function handleDownloadPrefix(prefix: string) {
    try {
      await DownloadPrefix(prefix);
    } catch (e: any) {
      addToast('error', 'Folder download failed: ' + e);
    }
  }

  function deleteSelected() {
    const keys = Array.from($selectedKeys);
    if (keys.length === 0) return;

    // Separate folders and files
    const folderKeys = keys.filter(k => k.endsWith('/'));
    const fileKeys = keys.filter(k => !k.endsWith('/'));
    const hasFolders = folderKeys.length > 0;

    // Build a list of item names for the confirmation message
    const names = keys.map(k => {
      const parts = k.replace(/\/$/, '').split('/');
      const name = parts[parts.length - 1];
      return k.endsWith('/') ? `📁 ${name}/` : name;
    });

    const namesList = names.length <= 8
      ? names.join('\n')
      : names.slice(0, 7).join('\n') + `\n...and ${names.length - 7} more`;

    const warning = hasFolders
      ? '\n\nFolders and all their contents will be permanently removed.'
      : '\n\nThis action cannot be undone.';

    const message = `Delete ${keys.length} item(s)?\n\n${namesList}${warning}`;

    confirmDialog = {
      show: true,
      title: 'Delete Objects',
      message,
      onConfirm: async () => {
        try {
          // Delete folders recursively
          for (const folder of folderKeys) {
            await DeletePrefix(folder);
          }
          // Delete files in batch
          if (fileKeys.length > 0) {
            await DeleteObjects(fileKeys);
          }
          addToast('success', `Deleted ${keys.length} item(s)`);
          await refresh();
        } catch (e: any) {
          addToast('error', 'Delete failed: ' + e);
        }
        confirmDialog.show = false;
      }
    };
  }

  async function handleCreateFolder(name: string, prefix?: string) {
    try {
      const targetPrefix = prefix !== undefined ? prefix : $currentPrefix;
      await CreateFolder(targetPrefix + name);
      addToast('success', `Folder "${name}" created`);
      await refresh();
    } catch (e: any) {
      addToast('error', 'Failed to create folder: ' + e);
    }
  }

  function handleContextMenu(e: CustomEvent<{ x: number; y: number; object: S3Object }>) {
    const { x, y, object } = e.detail;
    const items = [];
    if (object.isFolder) {
      items.push({ label: 'Open', action: () => navigateTo(object.key) });
      items.push({ label: 'Download Folder', action: () => handleDownloadPrefix(object.key) });
      items.push({ label: 'New Subfolder', action: () => promptCreateFolder(object.key) });
    } else {
      items.push({ label: 'Download', action: () => handleDownload([object.key]) });
    }
    items.push({ label: 'Copy Presigned URL', action: async () => {
      try {
        const url = await GetPresignedURL(object.key);
        await navigator.clipboard.writeText(url);
        addToast('success', 'Presigned URL copied (valid 1 hour)');
      } catch (e: any) {
        addToast('error', 'Failed to generate URL: ' + e);
      }
    }});
    items.push({ label: 'Copy Direct URL', action: async () => {
      try {
        const url = await GetDirectURL(object.key);
        await navigator.clipboard.writeText(url);
        addToast('info', 'Direct URL copied');
      } catch (e: any) {
        addToast('error', 'Failed to get URL: ' + e);
      }
    }});
    items.push({ label: 'Copy Key', action: () => navigator.clipboard.writeText(object.key) });
    items.push({ label: 'Delete', danger: true, action: () => {
      selectedKeys.set(new Set([object.key]));
      deleteSelected();
    }});
    contextMenu = { show: true, x, y, items };
  }

  function handleBgContextMenu(e: CustomEvent<{ x: number; y: number }>) {
    const { x, y } = e.detail;
    contextMenu = {
      show: true, x, y,
      items: [
        { label: 'New Folder', action: () => promptCreateFolder($currentPrefix) },
        { label: 'Upload Files', action: () => handleUpload() },
        { label: 'Upload Folder', action: () => handleUploadFolder() },
      ]
    };
  }

  function promptCreateFolder(prefix: string) {
    createFolderDialog = { show: true, parentPath: prefix };
  }

  function handleConnectionSaved() {
    loadConnections();
    showConnectionForm.set(false);
    editingConnection.set(null);
  }
</script>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
<div class="flex h-full w-full bg-surface-900" on:click={() => contextMenu.show = false}>
  <Sidebar
    on:connect={(e) => handleConnect(e.detail)}
    on:disconnect={handleDisconnect}
    on:edit={(e) => { editingConnection.set(e.detail); showConnectionForm.set(true); }}
    on:settings={() => showSettings = true}
  />

  <div class="flex flex-col flex-1 min-w-0">
    {#if connected}
      <TopBar
        on:refresh={refresh}
        on:upload={handleUpload}
        on:uploadFolder={handleUploadFolder}
        on:createFolder={() => promptCreateFolder($currentPrefix)}
        on:navigateTo={(e) => navigateTo(e.detail)}
        on:navigateUp={navigateUp}
        on:deleteSelected={deleteSelected}
        on:downloadSelected={() => handleDownload(Array.from($selectedKeys))}
      />

      <FileBrowser
        {isDragOver}
        on:navigate={(e) => navigateTo(e.detail)}
        on:contextmenu={handleContextMenu}
        on:bgContextmenu={handleBgContextMenu}
        on:download={(e) => handleDownload([e.detail])}
        on:deleteSelected={deleteSelected}
        on:upload={handleUpload}
        on:createFolder={() => promptCreateFolder($currentPrefix)}
      />
    {:else}
      <div class="flex-1 flex items-center justify-center">
        <div class="flex flex-col items-center">
          <div class="w-32 h-32 rounded-[32px] bg-surface-850 flex items-center justify-center mb-8 ring-1 ring-surface-750/50 shadow-2xl shadow-surface-950/60">
            <svg class="w-16 h-16 text-surface-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="0.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
            </svg>
          </div>
          <p class="text-[24px] font-bold text-surface-300 mb-3">No Connection</p>
          <p class="text-[16px] text-surface-500 leading-relaxed mb-8 text-center">Select a connection from the sidebar<br/>or create a new one to get started</p>
          <button class="px-10 py-4 text-[16px] font-bold rounded-2xl transition-all duration-200 inline-flex items-center gap-3 bg-gradient-to-r from-accent to-accent-dim hover:from-accent-hover hover:to-accent text-white shadow-xl shadow-accent/25 hover:shadow-2xl hover:shadow-accent/40 active:scale-[0.97] hover:-translate-y-0.5"
            on:click={() => { editingConnection.set(null); showConnectionForm.set(true); }}>
            <svg class="w-5.5 h-5.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
            </svg>
            Add Connection
          </button>
        </div>
      </div>
    {/if}

    <TransferPanel
      on:cancel={(e) => CancelTransfer(e.detail)}
      on:clear={() => ClearCompletedTransfers().then(() => {
        GetTransfers().then(t => { /* transfers updated via events */ });
      })}
    />
  </div>
</div>

{#if $showConnectionForm}
  <ConnectionForm
    connection={$editingConnection}
    on:save={handleConnectionSaved}
    on:close={() => { showConnectionForm.set(false); editingConnection.set(null); }}
  />
{/if}

{#if confirmDialog.show}
  <ConfirmDialog
    title={confirmDialog.title}
    message={confirmDialog.message}
    on:confirm={confirmDialog.onConfirm}
    on:cancel={() => confirmDialog.show = false}
  />
{/if}

{#if createFolderDialog.show}
  <CreateFolderDialog
    parentPath={createFolderDialog.parentPath}
    on:create={(e) => { handleCreateFolder(e.detail, createFolderDialog.parentPath); createFolderDialog.show = false; }}
    on:cancel={() => createFolderDialog.show = false}
  />
{/if}

{#if contextMenu.show}
  <ContextMenu x={contextMenu.x} y={contextMenu.y} items={contextMenu.items} on:close={() => contextMenu.show = false} />
{/if}

{#if showSettings}
  <SettingsDialog
    on:close={() => showSettings = false}
    on:imported={loadConnections}
  />
{/if}

<ToastContainer />
