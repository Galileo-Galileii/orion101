import { writable } from 'svelte/store';
import { storeWithInit } from './storeinit';
import assistants from './assistants';
import type { Assistant } from '$lib/services';

const def: Assistant = {
	id: '',
	icons: {},
	current: false
};

const store = writable<Assistant>(def);

export default storeWithInit(store, async () => {
	assistants.subscribe(async (assistants) => {
		for (const assistant of assistants) {
			if (assistant.current) {
				store.set(assistant);
				return;
			}
		}
		store.set(def);
	});
});
