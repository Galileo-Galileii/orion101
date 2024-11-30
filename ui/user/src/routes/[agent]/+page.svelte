<script lang="ts">
	import { page } from '$app/stores';
	import { profile } from '$lib/stores';
	import Navbar from '$lib/components/Navbar.svelte';
	import Messages from '$lib/components/Messages.svelte';
	import Editor from '$lib/components/Editor.svelte';
	import type { Editor as EditorType } from '$lib/components/Editor.svelte';
	import { loadFile } from '$lib/components/Editor.svelte';
	import Notifications from '$lib/components/Notifications.svelte';
	import { NotificationMessage } from '$lib/components/Notifications.svelte';
	import type { Input } from '$lib/components/messages/Input.svelte';
	import type { Messages as MessagesType } from '$lib/services';
	import { setContext } from 'svelte';

	let editorVisible = $state(false);
	let assistant = $page.params.agent;
	let notification: ReturnType<typeof Notifications>;
	let messageDiv: HTMLDivElement | undefined;
	let messages: ReturnType<typeof Messages>;

	function handleError(event: Error) {
		notification.addNotification(new NotificationMessage(event));
	}

	async function submit(e: CustomEvent<Input>) {
		return messages.submit(e.detail);
	}

	function handleLoadFile(e: string) {
		loadFile(assistant, e);
		editorVisible = true;
	}

	function handleMessages(e: MessagesType) {
		if (!messageDiv) {
			return;
		}

		// Check if messageDiv is already scrolled to the bottom
		let isScrolledToBottom =
			messageDiv.scrollHeight - messageDiv.clientHeight <= messageDiv.scrollTop + 1;

		const messages = e;
		if (messages.messages.length > 0 && messages.messages[messages.messages.length - 1].sent) {
			// If the last message is a sent (user input) message, scroll to the bottom always
			isScrolledToBottom = true;
		}

		if (isScrolledToBottom) {
			setTimeout(() => {
				if (messageDiv) {
					messageDiv.scrollTop = messageDiv.scrollHeight - messageDiv.clientHeight;
				}
			}, 100);
		}
	}

	setContext('editor', {
		loadFile: handleLoadFile
	} as EditorType);

	$effect(() => {
		if ($profile.unauthorized) {
			window.location.href = '/oauth2/start?rd=' + window.location.pathname;
		}
	});
</script>

<Navbar />

<main id="main-content" class="flex h-screen justify-center">
	<div class="relative flex w-1/2 flex-1 justify-center">
		<div bind:this={messageDiv} class="w-full overflow-auto px-8 pb-32 pt-16 scrollbar-none">
			<div class="mx-auto max-w-[1000px]">
				<Messages
					bind:this={messages}
					{assistant}
					onError={handleError}
					onMessages={handleMessages}
					onLoadFile={handleLoadFile}
				/>
			</div>
		</div>
	</div>

	{#if editorVisible}
		<div class="w-1/2 overflow-auto pb-16 pt-16 scrollbar-none">
			<Editor
				onEditorClose={() => {
					editorVisible = false;
				}}
				on:explain={submit}
				on:improve={submit}
			/>
		</div>
	{/if}

	<Notifications bind:this={notification} />
</main>
