import { writable } from 'svelte/store';
import type { Transfer } from '../types';

export const transfers = writable<Transfer[]>([]);
export const transferPanelOpen = writable(false);

export function updateTransfer(t: Transfer) {
  transfers.update((list) => {
    const idx = list.findIndex((x) => x.id === t.id);
    if (idx >= 0) {
      list[idx] = t;
      return [...list];
    }
    return [...list, t];
  });
}

export function clearTransfers() {
  transfers.update((list) =>
    list.filter((t) => t.status === 'pending' || t.status === 'in_progress')
  );
}
