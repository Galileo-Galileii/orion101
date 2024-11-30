<script lang="ts">
	import { Brain } from '$lib/icons';
	import { knowledgeFiles, currentAssistant } from '$lib/stores';
	import { ChatService } from '$lib/services';
	import { popover } from '$lib/actions';
	import Modal from '$lib/components/Modal.svelte';
	import KnowledgeFile from './KnowledgeFile.svelte';
	import KnowledgeUpload from '$lib/components/navbar/KnowledgeUpload.svelte';

	const { ref, tooltip, toggle } = popover({
		placement: 'bottom'
	});

	async function loadFiles() {
		knowledgeFiles.set(await ChatService.listKnowledgeFiles($currentAssistant.id));
	}

	let fileToDelete = $state<string | undefined>();

	async function deleteFile() {
		if (!fileToDelete) {
			return;
		}
		await ChatService.deleteKnowledgeFile($currentAssistant.id, fileToDelete);
		await loadFiles();
		fileToDelete = undefined;
	}
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
	<Brain class="h-5 w-5" />
</button>

<!-- Dropdown menu -->
<div
	use:tooltip
	class="w-96 divide-y divide-gray-100 rounded-lg bg-white shadow dark:divide-gray-600 dark:bg-gray-700"
>
	<div class="px-4 py-3 text-sm text-gray-900 dark:text-white">
		Knowledge Files
		<p class="mt-1 text-xs font-normal text-gray-500 dark:text-gray-300">
			Additional knowledge the AI can search.
		</p>
	</div>
	{#if $knowledgeFiles.items.length === 0}
		<p class="p-3 text-sm text-gray-500 dark:text-gray-300">No files</p>
	{:else}
		<ul class="space-y-3 p-3 text-sm" aria-labelledby="dropdownCheckboxButton">
			{#each $knowledgeFiles.items as file}
				<li>
					<KnowledgeFile
						{file}
						onDelete={() => {
							fileToDelete = file.fileName;
						}}
					/>
				</li>
			{/each}
		</ul>
	{/if}
	<KnowledgeUpload onUpload={loadFiles} />
</div>

<Modal
	show={fileToDelete !== undefined}
	msg={`Are you sure you want to delete ${fileToDelete}?`}
	onsuccess={deleteFile}
	oncancel={() => (fileToDelete = undefined)}
/>
