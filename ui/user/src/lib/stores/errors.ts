import { writable } from 'svelte/store';

const store = writable([] as Error[]);

function append(e: Error) {
	store.update((v) => {
		return [...v, e];
	});
	setTimeout(() => {
		store.update((v) => {
			return v.slice(1);
		});
	}, 2000);
}

function remove(index: number) {
	store.update((v) => {
		return v.filter((v, i) => i !== index);
	});
}

export default {
	subscribe: store.subscribe,
	append,
	remove
};
