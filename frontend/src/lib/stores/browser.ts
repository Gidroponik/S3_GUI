import { writable, derived } from 'svelte/store';
import type { S3Object, SortField, SortDir } from '../types';

export const objects = writable<S3Object[]>([]);
export const currentPrefix = writable('');
export const loading = writable(false);
export const searchQuery = writable('');
export const selectedKeys = writable<Set<string>>(new Set());
export const sortField = writable<SortField>('name');
export const sortDir = writable<SortDir>('asc');

export const breadcrumbs = derived(currentPrefix, ($prefix) => {
  if (!$prefix) return [];
  const parts = $prefix.split('/').filter(Boolean);
  return parts.map((part, i) => ({
    name: part,
    prefix: parts.slice(0, i + 1).join('/') + '/',
  }));
});

export const filteredObjects = derived(
  [objects, searchQuery, sortField, sortDir],
  ([$objects, $query, $sortField, $sortDir]) => {
    let filtered = $objects;
    if ($query) {
      const q = $query.toLowerCase();
      filtered = filtered.filter((o) => o.name.toLowerCase().includes(q));
    }
    return [...filtered].sort((a, b) => {
      if (a.isFolder && !b.isFolder) return -1;
      if (!a.isFolder && b.isFolder) return 1;

      let cmp = 0;
      switch ($sortField) {
        case 'name':
          cmp = a.name.localeCompare(b.name);
          break;
        case 'size':
          cmp = a.size - b.size;
          break;
        case 'lastModified':
          cmp = new Date(a.lastModified).getTime() - new Date(b.lastModified).getTime();
          break;
      }
      return $sortDir === 'asc' ? cmp : -cmp;
    });
  }
);
