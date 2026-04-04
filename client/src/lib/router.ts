import { writable } from 'svelte/store';

export type Page =
  | { name: 'garden' }
  | { name: 'seeds' }
  | { name: 'herbarium' }
  | { name: 'feed' }
  | { name: 'leaderboard' }
  | { name: 'seed-detail'; id: number }
  | { name: 'herbarium-detail'; id: number }
  | { name: 'flower-detail'; userId: number };

export type TabName = 'garden' | 'feed' | 'leaderboard';
export const gardenTabs = ['garden', 'seeds', 'herbarium'] as const;

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

/** Switch between the 3 garden sub-screens, keeping feed/leaderboard untouched */
export function switchGardenTab(tab: 'garden' | 'seeds' | 'herbarium') {
  pageStack.set([{ name: tab }]);
}
