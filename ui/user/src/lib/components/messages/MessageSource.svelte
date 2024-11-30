<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { ChatService } from '$lib/services';
	import type { Progress } from '$lib/services';

	interface Props {
		assistant: string;
		onMessage?: (event: Progress) => void;
		onError?: (event: Error) => void;
	}

	let { assistant, onMessage = () => {}, onError = () => {} }: Props = $props();

	let es: EventSource;
	let replayComplete = false;

	function disconnect() {
		if (es) {
			es.close();
		}
	}

	function connect() {
		disconnect();
		es = ChatService.newMessageEventSource(assistant);
		es.onmessage = handleMessage;
		es.onopen = () => {
			console.log('Message EventSource opened');
		};
		es.onerror = (e: Event) => {
			if (e.eventPhase === EventSource.CLOSED) {
				console.log('Message EventSource closed');
			}
		};
	}

	function handleMessage(event: MessageEvent) {
		const message = JSON.parse(event.data) as Progress;
		if (message.replayComplete) {
			replayComplete = true;
		}
		if (message.error) {
			if (replayComplete) {
				onError(new Error(message.error));
			}
		} else {
			onMessage(message);
		}
	}

	onMount(connect);
	onDestroy(disconnect);
</script>
