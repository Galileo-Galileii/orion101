<script lang="ts">
	import type { Message } from '$lib/services';
	import Loading from '$lib/icons/Loading.svelte';
	import highlight from 'highlight.js';
	import MessageIcon from '$lib/components/messages/MessageIcon.svelte';
	import { FileText } from '$lib/icons';
	import { toHTMLFromMarkdown } from '$lib/markdown.js';

	interface Props {
		msg: Message;
		onLoadFile?: (filename: string) => void;
	}

	let { msg, onLoadFile = () => {} }: Props = $props();

	let content = $derived(msg.message ? msg.message.join('') : '');
	let fullWidth = !msg.sent && !msg.oauthURL && !msg.tool;
	let showBubble = msg.sent;
	let renderMarkdown = !msg.sent && !msg.oauthURL && !msg.tool;

	$effect(() => {
		// this is a hack to make sure this effect is run after the content is updated
		if (content.length == 0) {
			return;
		}

		const blocks = document.querySelectorAll('.message-content pre > code');
		blocks.forEach((block) => {
			if (block instanceof HTMLElement && block.dataset.highlighted !== 'yes') {
				highlight.highlightElement(block);
			}
		});

		const links = document.querySelectorAll('.message-content a');
		links.forEach((link) => {
			if (link instanceof HTMLAnchorElement && link.target == '') {
				link.target = '_blank';
			}
		});
	});
</script>

<div class="flex items-start gap-2.5" class:justify-end={msg.sent}>
	<MessageIcon {msg} />

	<div class="leading-1.5 flex w-full flex-col" class:w-full={fullWidth}>
		<div class="mb-2 flex items-center space-x-2 rtl:space-x-reverse">
			{#if msg.sourceName}
				<span class="text-sm font-semibold text-gray-900 dark:text-white">{msg.sourceName}</span>
			{/if}
			{#if msg.time}
				<span class="text-sm font-normal text-gray-500 dark:text-gray-400"
					>{msg.time.toLocaleDateString(undefined, {
						year: 'numeric',
						month: 'short',
						day: 'numeric',
						hour: 'numeric',
						minute: 'numeric'
					})}</span
				>
			{/if}
		</div>
		<div
			style:display={showBubble ? 'flex' : 'contents'}
			class:message-content={renderMarkdown}
			class="leading-1.5 flex flex-col rounded-e-xl rounded-es-xl border-gray-200 bg-gray-900 p-4 text-white dark:bg-gray-700"
		>
			{#if msg.oauthURL}
				<a
					href={msg.oauthURL}
					class="rounded-xl bg-ablue-900 p-4
						text-white
						hover:bg-ablue2-600
					  hover:dark:bg-ablue2-600"
					target="_blank"
					>Authentication is required
					<span class="underline">click here</span> to log-in using OAuth
				</a>
			{:else if content}
				{#if msg.sent}
					{#if msg.explain}
						<div class="flex items-center gap-1 pb-3 pl-1">
							<FileText class="h-4 w-4 text-white" />
							<span>{msg.explain.filename}</span>
						</div>
						<pre
							class="mb-4 overflow-x-auto rounded border-white bg-gray-100 p-4 text-black shadow-white">{msg
								.explain.selection}</pre>
					{/if}
					{content}
				{:else}
					{@html toHTMLFromMarkdown(content)}
				{/if}
			{/if}
			{#if msg.file?.filename}
				<div class="flex items-center">
					<button
						onclick={() => {
							if (msg.file?.filename) {
								onLoadFile(msg.file?.filename);
							}
						}}
						class="flex items-center gap-2 rounded border border-gray-200 p-2 px-4 text-black shadow hover:bg-gray-100 dark:bg-gray-900 dark:text-white hover:dark:bg-gray-700"
					>
						<FileText class="text-black" />
						<span>{msg.file.filename}</span>
					</button>
				</div>
			{/if}
			{#if !msg.sent}
				<div class="mt-3 flex h-6 items-end justify-end">
					{#if !msg.done}
						<div class="mx-1.5" role="status">
							<Loading />
							<span class="sr-only">Loading...</span>
						</div>
						<span class="text-sm font-normal text-gray-500 dark:text-gray-400">Loading...</span>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</div>
