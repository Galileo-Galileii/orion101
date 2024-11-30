<script lang="ts">
	import { CircleX, FileText, Trash } from '$lib/icons';
	import { type KnowledgeFile } from '$lib/services';
	import { popover } from '$lib/actions';
	import Loading from '$lib/icons/Loading.svelte';

	interface Props {
		onDelete?: () => void;
		file: KnowledgeFile;
	}

	const { onDelete, file }: Props = $props();

	const tt = popover({
		placement: 'top-start',
		hover: true
	});
</script>

<div class="group flex" use:tt.ref>
	<button class="flex flex-1 items-center">
		<FileText class="h-8 w-8 rounded-md bg-gray-100 p-1 text-black" />
		<span class="ms-3 text-sm font-medium dark:text-gray-100">{file.fileName.slice(0, 32)}</span>
		{#if file.state === 'error' || file.state === 'failed'}
			<CircleX class="h-4 text-red-500" />
			<div class="rounded-md bg-red-500 p-2 text-white" use:tt.tooltip>
				{file.error ? file.error : 'Failed'}
			</div>
		{:else if file.state === 'pending' || file.state === 'ingesting'}
			<Loading />
		{/if}
	</button>
	<button
		class="hidden group-hover:block"
		onclick={() => {
			if (file.state === 'ingested') {
				onDelete?.();
			}
		}}
	>
		<Trash class="h-5 w-5 text-gray-400" />
	</button>
</div>
