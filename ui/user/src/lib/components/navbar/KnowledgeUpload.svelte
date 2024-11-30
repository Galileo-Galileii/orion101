<script lang="ts">
	import { Upload } from '$lib/icons';
	import { ChatService } from '$lib/services';
	import type { KnowledgeFile } from '$lib/services';
	import Loading from '$lib/icons/Loading.svelte';
	import Error from '$lib/components/Error.svelte';
	import { currentAssistant } from '$lib/stores';

	interface Props {
		onUpload?: () => void | Promise<void>;
	}

	let { onUpload }: Props = $props();

	let files = $state<FileList>();
	let uploadInProgress = $state<Promise<KnowledgeFile>>();

	$effect(() => {
		if (!files?.length) {
			return;
		}

		uploadInProgress = ChatService.uploadKnowledge($currentAssistant.id, files[0]);
		uploadInProgress
			.then(() => {
				onUpload?.();
			})
			.catch((error) => {
				console.error(error);
			})
			.finally(() => {
				uploadInProgress = undefined;
			});

		files = undefined;
	});
</script>

<label class="flex cursor-pointer justify-end gap-2 p-4 text-gray-300 hover:text-gray-100">
	Upload (.pdf, .pptx, .txt)
	<input bind:files type="file" class="hidden" accept=".pdf, .txt, .docs" />
	{#await uploadInProgress}
		<Loading class="h-5 w-5" />
	{:catch error}
		<Error {error} />
	{/await}
	{#if !uploadInProgress}
		<Upload class="h-5 w-5" />
	{/if}
</label>
