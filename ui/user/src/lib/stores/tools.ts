import { type AssistantToolList, ChatService } from '$lib/services';
import { writable } from 'svelte/store';
import { storeWithInit } from '$lib/stores/storeinit';
import assistants from '$lib/stores/assistants';

const store = writable<AssistantToolList>({
	readonly: true,
	items: []
});

export default storeWithInit(store, async () => {
	assistants.subscribe(async (assistants) => {
		for (const assistant of assistants) {
			if (assistant.current && assistant.id) {
				store.set(await ChatService.listTools(assistant.id));
				break;
			}
		}
	});
});
