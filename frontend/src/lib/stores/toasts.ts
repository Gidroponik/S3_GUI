import { writable } from 'svelte/store';
import type { Toast } from '../types';
import { nextToastId } from '../utils';

export const toasts = writable<Toast[]>([]);

export function addToast(type: Toast['type'], message: string, duration = 4000) {
  const id = nextToastId();
  toasts.update((t) => [...t, { id, type, message }]);
  setTimeout(() => {
    toasts.update((t) => t.filter((x) => x.id !== id));
  }, duration);
}
