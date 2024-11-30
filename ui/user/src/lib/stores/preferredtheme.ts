import { storeWithInit } from './storeinit';
import { writable } from 'svelte/store';

type PreferredTheme = 'light' | 'dark';

const store = storeWithInit(writable('light' as PreferredTheme), init);

function init() {
	const mm = window.matchMedia('(prefers-color-scheme: dark)');
	mm.addEventListener('change', (e) => {
		store.set(e.matches ? 'dark' : 'light');
	});
	store.set(mm.matches ? 'dark' : 'light');
}

// mask writable as readable
export default {
	subscribe: store.subscribe,
	init
};
