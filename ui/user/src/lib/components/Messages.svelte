<script lang="ts">
	import Input from '$lib/components/messages/Input.svelte';
	import type { Input as InputType } from '$lib/components/messages/Input.svelte';
	import MessageSource from '$lib/components/messages/MessageSource.svelte';
	import type { Messages, Progress } from '$lib/services';
	import { ChatService } from '$lib/services';
	import Message from '$lib/components/messages/Message.svelte';
	import { fade } from 'svelte/transition';

	interface Props {
		assistant: string;
		onMessages: (messages: Messages) => void;
		onError: (err: Error) => void;
		onFocus?: () => void;
		onLoadFile: (filename: string) => void;
	}

	let { assistant, onMessages, onError, onFocus, onLoadFile }: Props = $props();

	let progressEvents: Progress[] = [];
	let replayComplete = false;
	let messages: Messages = $state({ messages: [], inProgress: false });

	function handleMessage(progress: Progress) {
		progressEvents = [...progressEvents, progress];
		if (!replayComplete) {
			replayComplete = progressEvents.find((e) => e.replayComplete) !== undefined;
		}

		if (!replayComplete) {
			return;
		}

		messages = ChatService.progressToMessages(progressEvents);

		// forward the messages to the parent component
		onMessages(messages);
	}

	let inputBox: ReturnType<typeof Input>;

	export async function submit(input: InputType) {
		return inputBox.submit(input);
	}
</script>

<MessageSource {assistant} onMessage={handleMessage} {onError} />

<div transition:fade|global class="flex flex-col gap-8">
	{#each messages.messages as msg}
		{#if !msg.ignore}
			<Message {msg} {onLoadFile} />
		{/if}
	{/each}
	<Input {assistant} bind:this={inputBox} readonly={messages.inProgress} {onError} {onFocus} />
</div>
