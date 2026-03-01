import { writable } from 'svelte/store';
import type { Connection } from '../types';

export const connections = writable<Connection[]>([]);
export const activeConnectionId = writable<string | null>(null);
export const showConnectionForm = writable(false);
export const editingConnection = writable<Connection | null>(null);
