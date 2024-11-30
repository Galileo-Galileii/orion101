<!-- @migration-task Error while migrating Svelte code: `<button>` is invalid inside `<button>` -->
<script lang="ts" module>
	import { editorFiles } from '$lib/stores';
	import { ChatService } from '$lib/services';
	import type { EditorFile } from '$lib/stores/editorfiles';

	export interface Editor {
		loadFile: (file: string) => void;
	}

	function hasFile(name: string): boolean {
		let files: EditorFile[] | undefined;
		editorFiles.subscribe((value) => (files = value))();
		return files?.find((file) => file.name === name) !== undefined;
	}

	export async function loadFile(assistant: string, file: string) {
		if (hasFile(file)) {
			selectFile(file);
			return;
		}

		try {
			const contents = await ChatService.getFile(assistant, file);
			const targetFile = {
				name: file,
				contents,
				buffer: '',
				modified: false,
				selected: true
			};

			editorFiles.update((files) => {
				const index = files.findIndex((f) => f.name === targetFile.name);
				if (index === -1) {
					files.push(targetFile);
				}
				return files;
			});
			selectFile(targetFile.name);
		} catch {
			// ignore error
		}
	}

	function selectFile(name: string) {
		editorFiles.update((files) => {
			let found = false;
			files.forEach((file) => {
				file.selected = file.name === name;
				found = found || file.selected;
			});
			if (!found && files.length > 0) {
				files[0].selected = true;
			}
			return files;
		});
	}
</script>

<script lang="ts">
	import { FileText, X } from '$lib/icons';
	import { onMount } from 'svelte';
	import Milkdown from '$lib/components/editor/Milkdown.svelte';
	import Codemirror from '$lib/components/editor/Codemirror.svelte';

	interface Props {
		onEditorClose?: () => void;
	}

	const { onEditorClose }: Props = $props();

	function fileChanged(e: CustomEvent<{ name: string; contents: string }>) {
		editorFiles.update((files) => {
			files.forEach((file) => {
				if (file.name === e.detail.name) {
					file.buffer = e.detail.contents;
					file.modified = true;
				}
			});
			return files;
		});
	}

	function remove(i: number) {
		editorFiles.update((files) => {
			files.splice(i, 1);
			if (files.length == 0) {
				onEditorClose?.();
			}
			return files;
		});
		selectFile(i == 0 ? '' : $editorFiles[i - 1].name);
	}

	function isSelected(name: string): boolean {
		if ($editorFiles.length <= 1) {
			return true;
		}
		return $editorFiles.find((file) => file.name === name)?.selected || false;
	}

	onMount(() => {
		if (window.location.href.indexOf('#editor:') > -1) {
			selectFile(window.location.href.split('#editor:')[1]);
		} else {
			selectFile('');
		}
	});
</script>

<div>
	<div class="flex items-center justify-between">
		<ul class="mb-4 flex flex-wrap text-center text-sm font-medium">
			{#each $editorFiles as file, i}
				<li class="me-2">
					<a
						href={`#editor:${$editorFiles[i].name}`}
						class:selected={isSelected(file.name)}
						onclick={() => {
							selectFile($editorFiles[i].name);
							window.location.href = `#editor:${$editorFiles[i].name}`;
						}}
						class="selected active group flex items-center justify-center gap-2 rounded-t-lg p-4 text-black dark:border-blue-500 dark:text-white"
						aria-current="page"
					>
						<FileText />
						<span>{file.name}</span>
						<button
							class="ml-2"
							onclick={() => {
								remove(i);
							}}
						>
							<X />
						</button>
					</a>
				</li>
			{/each}
		</ul>
		<button
			class="icon-button"
			onclick={() => {
				onEditorClose?.();
			}}
		>
			<X />
		</button>
	</div>

	<div
		id="editor"
		onkeydown={(e) => {
			e.stopPropagation();
		}}
		role="none"
	>
		{#each $editorFiles as file}
			<div class:hidden={!isSelected(file.name)}>
				{#if file.name.toLowerCase().endsWith('.md')}
					<Milkdown {file} on:changed={fileChanged} />
				{:else}
					<Codemirror {file} on:changed={fileChanged} on:explain on:improve />
				{/if}
			</div>
		{/each}
	</div>
</div>

<style lang="postcss">
	.selected {
		@apply border-b-2 border-blue-600;
	}
</style>
