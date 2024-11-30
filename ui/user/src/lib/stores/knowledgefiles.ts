import { type KnowledgeFiles } from '$lib/services';
import { writable } from 'svelte/store';

const store = writable<KnowledgeFiles>({
	items: []
});

export default store;
