import { writable } from 'svelte/store';
import { ChatService, type Profile } from '$lib/services';
import { storeWithInit } from './storeinit';

const store = writable<Profile>({
	email: '',
	iconURL: '',
	role: 0
});

async function init() {
	try {
		store.set(await ChatService.getProfile());
	} catch (e) {
		if (e instanceof Error && e.message.startsWith('403')) {
			store.set({
				email: '',
				iconURL: '',
				role: 0,
				unauthorized: true
			});
		} else {
			setTimeout(init, 5000);
		}
	}
}
export default storeWithInit(store, init);
