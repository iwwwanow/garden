import { writable } from 'svelte/store';

export type Page =
  | { name: 'garden' }
  | { name: 'feed' }
  | { name: 'leaderboard' }
  | { name: 'seed-detail'; id: number }
  | { name: 'herbarium-detail'; id: number }
  | { name: 'flower-detail'; userId: number };

export type TabName = 'garden' | 'feed' | 'leaderboard';

export const pageStack = writable<Page[]>([{ name: 'garden' }]);

export function navigate(page: Page) {
  pageStack.update(s => [...s, page]);
}

export function back() {
  pageStack.update(s => (s.length > 1 ? s.slice(0, -1) : s));
}

export function switchTab(tab: TabName) {
  pageStack.set([{ name: tab }]);
}
