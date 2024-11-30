import { type Assistant, ChatService } from '$lib/services';
import { writable } from 'svelte/store';
import { storeWithInit } from '$lib/stores/storeinit';
import { page } from '$app/stores';

function assignSelected(assistants: Assistant[], selectedName: string): Assistant[] {
	const result: Assistant[] = [];

	for (const assistant of assistants) {
		assistant.current = selectedName !== '' && assistant.id === selectedName;
		result.push(assistant);
	}

	return result;
}

const store = writable<Assistant[]>(assignSelected([], ''));

export default storeWithInit(store, async () => {
	page.subscribe(async (value) => {
		const selectedName = value.params?.agent ?? '';
		try {
			const assistants = await ChatService.listAssistants();
			store.set(assignSelected(assistants.items, selectedName));
		} catch {
			// just ignore
		}
	});
});
