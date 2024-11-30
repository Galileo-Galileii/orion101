<script module lang="ts">
	import type { Explain } from '$lib/services';

	export interface Input {
		prompt: string;
		explain?: Explain;
		improve?: Explain;
		changedFiles: Record<string, string>;
	}
</script>

<script lang="ts">
	import { ChatService } from '$lib/services';
	import { onMount } from 'svelte';
	import { editorFiles } from '$lib/stores';
	import { SendHorizontal } from '$lib/icons';

	interface Props {
		onError?: (err: Error) => void;
		onFocus?: () => void;
		readonly?: boolean;
		assistant: string;
	}

	let { onError = () => {}, onFocus = () => {}, assistant, readonly }: Props = $props();

	let value = $state('');
	let chat: HTMLTextAreaElement;

	function toInvokeInput(input: Input): Input | string {
		// This is just to make it pretty and send simple prompts if we can
		if (input.explain || input.improve) {
			return input;
		}
		if (input.changedFiles && Object.keys(input.changedFiles).length !== 0) {
			return input;
		}
		if (input.prompt) {
			return input.prompt;
		}
		return input;
	}

	export async function submit(override?: Input) {
		let input: Input = {
			prompt: value,
			changedFiles: {}
		};

		if (override) {
			input = override;
			if (chat) {
				chat.focus();
			}
		}

		for (const file of $editorFiles) {
			if (file.modified) {
				input.changedFiles[file.name] = file.buffer;
			}
		}

		try {
			await ChatService.invoke(assistant, toInvokeInput(input));
		} catch (err) {
			if (err instanceof Error) {
				onError(err);
			} else {
				onError(new Error(String(err)));
			}
			return;
		}

		if (input.changedFiles) {
			editorFiles.update((files) => {
				for (const file of files) {
					if (input.changedFiles[file.name]) {
						file.contents = input.changedFiles[file.name];
						file.modified = false;
						file.buffer = '';
					}
				}
				return files;
			});
		}

		value = '';
	}

	async function onKey(e: KeyboardEvent) {
		if (e.key !== 'Enter' || e.shiftKey) {
			return;
		}
		e.preventDefault();
		await submit();
		resize();
	}

	function resize() {
		chat.style.height = 'auto';
		// not totally sure why 4 is needed here, otherwise the textarea is too small and we
		// get a scrollbar
		chat.style.height = chat.scrollHeight + 4 + 'px';
	}

	function focusChat(e: KeyboardEvent) {
		if (!chat) {
			return;
		}

		const alphanumericRegex = /^[a-zA-Z0-9]$/;
		if (alphanumericRegex.test(e.key) && !e.ctrlKey && !e.altKey && !e.metaKey) {
			chat.focus();
		} else if ((e.key == 'Backspace' || e.key == 'ArrowLeft') && chat.value !== '') {
			chat.focus();
		}
	}

	onMount(() => {
		document.addEventListener('keydown', focusChat);
	});
</script>

<div
	class="absolute inset-x-0 bottom-0 flex justify-center bg-gradient-to-t from-white px-12 pb-8 pt-10 dark:from-black"
>
	<div class="w-full max-w-[700px]">
		<label for="chat" class="sr-only">Your message</label>
		<div class="flex items-center rounded-lg px-3 py-2 dark:bg-black">
			<textarea
				id="chat"
				rows="1"
				bind:value
				{readonly}
				onkeydown={onKey}
				bind:this={chat}
				oninput={resize}
				onfocus={onFocus}
				class="ml-4 mr-2 block w-full resize-none rounded-lg border-2 border-gray-300 bg-white p-2.5 text-sm text-gray-900
								outline-none focus:border-blue-500 focus:ring-blue-500 dark:border-gray-300 dark:bg-black
								 dark:text-white dark:placeholder-gray-300 dark:focus:border-blue-500 dark:focus:ring-blue-500"
				placeholder="Your message..."
			></textarea>
			<button
				type="submit"
				onclick={() => submit()}
				class="inline-flex cursor-pointer justify-center rounded-full p-2 text-blue-600
							 hover:bg-blue-100 dark:text-blue-500 dark:hover:bg-gray-800"
			>
				<SendHorizontal class="text-blue-500 dark:text-blue-500" />
				<span class="sr-only">Send message</span>
			</button>
		</div>
	</div>
</div>
