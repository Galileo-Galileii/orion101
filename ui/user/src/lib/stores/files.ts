import { type Files } from '$lib/services';
import { writable } from 'svelte/store';

const store = writable<Files>({
	items: []
});

export default store;
