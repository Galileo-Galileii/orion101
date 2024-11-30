<script lang="ts">
	import { FileText, Trash } from '$lib/icons';
	import { files, currentAssistant } from '$lib/stores';
	import { ChatService } from '$lib/services';
	import { popover } from '$lib/actions';
	import Modal from '$lib/components/Modal.svelte';
	import { getContext } from 'svelte';
	import type { Editor } from '$lib/components/Editor.svelte';

	const { ref, tooltip, toggle } = popover({
		placement: 'bottom'
	});

	async function loadFiles() {
		files.set(await ChatService.listFiles($currentAssistant.id));
	}

	async function deleteFile() {
		if (!fileToDelete) {
			return;
		}
		await ChatService.deleteFile($currentAssistant.id, fileToDelete);
		await loadFiles();
		fileToDelete = undefined;
	}

	const editor: Editor = getContext('editor');
	let fileToDelete = $state<string | undefined>();
</script>

<button
	use:ref
	class="icon-button z-20"
	onclick={async () => {
		await loadFiles();
		toggle();
	}}
	type="button"
>
	<FileText class="h-5 w-5" />
</button>

<!-- Dropdown menu -->
<div
	use:tooltip
	class="w-96 divide-y divide-gray-100 rounded-lg bg-white shadow dark:divide-gray-600 dark:bg-gray-700"
>
	<div class="px-4 py-3 text-sm text-gray-900 dark:text-white">
		Files
		<p class="mt-1 text-xs font-normal text-gray-500 dark:text-gray-300">
			Click to view or edit files
		</p>
	</div>
	{#if $files.items.length === 0}
		<p class="p-3 text-sm text-gray-500 dark:text-gray-300">No files</p>
	{:else}
		<ul class="space-y-3 p-3 text-sm" aria-labelledby="dropdownCheckboxButton">
			{#each $files.items as file}
				<li class="group">
					<div class="flex">
						<button
							class="flex flex-1 items-center"
							onclick={() => {
								editor.loadFile(file.name);
							}}
						>
							<FileText class="h-8 w-8 rounded-md bg-gray-100 p-1 text-black" />
							<span class="ms-3 text-sm font-medium dark:text-gray-100">{file.name}</span>
						</button>
						<button
							class="hidden group-hover:block"
							onclick={() => {
								fileToDelete = file.name;
							}}
						>
							<Trash class="h-5 w-5 text-gray-400" />
						</button>
					</div>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<Modal
	show={fileToDelete !== undefined}
	msg={`Are you sure you want to delete ${fileToDelete}?`}
	onsuccess={deleteFile}
	oncancel={() => (fileToDelete = undefined)}
/>
