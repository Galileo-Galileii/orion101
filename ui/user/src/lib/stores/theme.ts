import { writable } from 'svelte/store';
import { storeWithInit } from '$lib/stores/storeinit';

export type Theme = 'light' | 'dark' | 'system';
export const DefaultTheme: Theme = 'system';

const store = writable(DefaultTheme as Theme);

function init() {
	const initialTheme = localStorage.getItem('theme') as Theme | null;
	if (initialTheme) {
		store.set(initialTheme);
	}

	store.subscribe((value) => {
		if (value === 'system' && localStorage.getItem('theme') === null) {
			// no point in setting it
			return;
		}
		localStorage.setItem('theme', value);
	});
}

export default {
	...storeWithInit(store, init),
	init
};
