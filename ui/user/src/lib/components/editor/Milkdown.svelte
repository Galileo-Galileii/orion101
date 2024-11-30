<script lang="ts">
	import { run } from 'svelte/legacy';

	import type { EditorFile } from '$lib/stores';
	import { Crepe } from '@milkdown/crepe';

	import '@milkdown/crepe/theme/common/style.css';
	import '@milkdown/crepe/theme/nord.css';
	import { listener, listenerCtx } from '@milkdown/kit/plugin/listener';
	import { createEventDispatcher } from 'svelte';
	import { replaceAll } from '@milkdown/utils';

	interface Props {
		file: EditorFile;
	}

	let { file }: Props = $props();

	let setValue: (value: string) => void | undefined = $state();
	let lastSetValue = '';
	let focused = false;
	const dispatcher = createEventDispatcher();

	run(() => {
		setValue?.(file?.contents);
	});

	function editor(node: HTMLElement) {
		lastSetValue = file.contents;

		node.addEventListener('focusin', () => {
			focused = true;
		});
		node.addEventListener('focusout', () => {
			focused = false;
		});

		const crepe = new Crepe({
			root: node,
			defaultValue: file.contents
		});

		crepe.editor
			.config((ctx) => {
				const listener = ctx.get(listenerCtx);
				listener.markdownUpdated((ctx, markdown, prevMarkdown) => {
					if (focused && markdown !== prevMarkdown) {
						dispatcher('changed', {
							name: file.name,
							contents: markdown
						});
					}
				});
			})
			.use(listener);

		crepe.create().then(() => {
			setValue = (value: string) => {
				if (lastSetValue === value) {
					return;
				}

				lastSetValue = value;
				crepe.editor.action(replaceAll(value));
			};
		});

		return {
			destroy() {
				crepe.destroy();
			}
		};
	}
</script>

<div class="milkdown-editor" use:editor></div>
