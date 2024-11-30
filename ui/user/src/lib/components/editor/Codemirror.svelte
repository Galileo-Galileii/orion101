<script lang="ts">
	import { run } from 'svelte/legacy';

	import {
		lineNumbers,
		highlightActiveLineGutter,
		highlightSpecialChars,
		drawSelection,
		dropCursor,
		keymap
	} from '@codemirror/view';
	import {
		foldGutter,
		indentOnInput,
		syntaxHighlighting,
		defaultHighlightStyle,
		bracketMatching,
		foldKeymap
	} from '@codemirror/language';
	import { history, defaultKeymap, historyKeymap } from '@codemirror/commands';
	import { searchKeymap } from '@codemirror/search';
	import {
		closeBrackets,
		autocompletion,
		closeBracketsKeymap,
		completionKeymap
	} from '@codemirror/autocomplete';
	import { lintKeymap } from '@codemirror/lint';
	import { javascript } from '@codemirror/lang-javascript';
	import { go } from '@codemirror/lang-go';
	import { cpp } from '@codemirror/lang-cpp';
	import { css } from '@codemirror/lang-css';
	import { html } from '@codemirror/lang-html';
	import { sql } from '@codemirror/lang-sql';
	import { vue } from '@codemirror/lang-vue';
	import { sass } from '@codemirror/lang-sass';
	import { rust } from '@codemirror/lang-rust';
	import { java } from '@codemirror/lang-java';
	import { EditorState, Transaction } from '@codemirror/state';
	import { EditorView } from '@codemirror/view';
	import { LanguageSupport } from '@codemirror/language';
	import type { EditorFile } from '$lib/stores';
	import { createEventDispatcher } from 'svelte';
	import type { Tooltip, TooltipView } from '@codemirror/view';
	import { showTooltip } from '@codemirror/view';
	import { StateField } from '@codemirror/state';
	import { oneDark } from '@codemirror/theme-one-dark';
	import { MessageSquareText, CircleHelp } from '$lib/icons';
	import { darkMode } from '$lib/stores';

	const cursorTooltipField = StateField.define<readonly Tooltip[]>({
		create: getExplainTooltip,
		update: updateExplainToolTip,
		provide: (f) => showTooltip.computeN([f], (state) => state.field(f))
	});

	let explain = $state<HTMLElement>();

	function newExplainToolTip(state: EditorState): TooltipView {
		const tooltip = explain?.cloneNode(true) as HTMLElement;
		const explainButton = tooltip.querySelector('.explain-button');
		const improveButton = tooltip.querySelector('.improve-button');
		const improveText = tooltip.querySelector('.improve-text');
		const improveTextarea = tooltip.querySelector('.improve-textarea') as HTMLTextAreaElement;

		explainButton?.addEventListener('click', () => {
			dispatch('explain', {
				explain: {
					filename: file.name,
					selection: state
						.sliceDoc(state.selection.ranges[0].from, state.selection.ranges[0].to)
						.toString()
				}
			});
			tooltip.classList.add('hidden');
		});

		improveButton?.addEventListener('click', (e) => {
			e.preventDefault();
			explainButton?.classList.add('hidden');
			improveButton?.classList.add('hidden');
			improveText?.classList.remove('hidden');
			improveTextarea?.focus();
		});

		improveTextarea?.addEventListener('keydown', (e: KeyboardEvent) => {
			if (e.key !== 'Enter' || e.shiftKey) {
				return;
			}
			e.preventDefault();
			dispatch('improve', {
				prompt: improveTextarea.value,
				improve: {
					filename: file.name,
					selection: state
						.sliceDoc(state.selection.ranges[0].from, state.selection.ranges[0].to)
						.toString()
				}
			});
			tooltip.remove();
		});

		return {
			dom: tooltip,
			offset: {
				x: 0,
				y: 12
			}
		};
	}

	function updateExplainToolTip(tooltips: readonly Tooltip[], tr: Transaction): readonly Tooltip[] {
		if (!focused) {
			return [];
		}
		if (
			tooltips.length == 0 ||
			tr.state.selection.ranges.length != 1 ||
			tr.state.selection.ranges[0].from == tr.state.selection.ranges[0].to
		) {
			return getExplainTooltip(tr.state);
		}

		const obj = tooltips[0] as object;
		if ('state' in obj && obj.state instanceof EditorState) {
			if (
				obj.state.selection.ranges[0]?.from != tr.state.selection.ranges[0].from ||
				obj.state.selection.ranges[0]?.to != tr.state.selection.ranges[0].to
			) {
				return getExplainTooltip(tr.state);
			}
			obj.state = tr.state;
		}
		return tooltips;
	}

	function getExplainTooltip(state: EditorState): readonly Tooltip[] {
		if (
			state.selection.ranges.length == 1 &&
			state.selection.ranges[0].from != state.selection.ranges[0].to
		) {
			const tooltip = {
				pos: state.selection.ranges[0].from,
				above: true,
				strictSide: true,
				create: () => {
					return newExplainToolTip(state);
				},
				state: state,
				selection: state.selection.ranges[0]
			} as Tooltip;
			return [tooltip];
		}
		return [];
	}

	const languages = {
		js: javascript,
		java: java,
		go: go,
		c: cpp,
		h: cpp,
		hpp: cpp,
		cpp: cpp,
		css: css,
		html: html,
		sql: sql,
		vue: vue,
		sass: sass,
		rs: rust
	} as Record<string, () => LanguageSupport>;

	interface Props {
		file: EditorFile;
	}

	let { file }: Props = $props();

	let setValue: (value: string) => void | undefined = $state();
	let lastSetValue = '';
	let focused = false;
	let dispatch = createEventDispatcher();

	run(() => {
		setValue?.(file?.contents);
	});

	const basicSetup = (() => [
		lineNumbers(),
		highlightActiveLineGutter(),
		highlightSpecialChars(),
		history(),
		foldGutter(),
		drawSelection(),
		dropCursor(),
		EditorState.allowMultipleSelections.of(true),
		indentOnInput(),
		syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
		bracketMatching(),
		closeBrackets(),
		autocompletion(),
		keymap.of([
			...closeBracketsKeymap,
			...defaultKeymap,
			...searchKeymap,
			...historyKeymap,
			...foldKeymap,
			...completionKeymap,
			...lintKeymap
		]),
		cursorTooltipField
	])();

	function editor(targetElement: HTMLElement) {
		const ext = file.name.split('.').pop() || 'js';
		const langExtension = languages[ext] ?? javascript;
		lastSetValue = file.contents;

		targetElement.addEventListener('focusin', () => {
			focused = true;
		});
		targetElement.addEventListener('focusout', () => {
			focused = false;
		});

		const updater = EditorView.updateListener.of((update) => {
			if (update.docChanged && focused) {
				dispatch('changed', {
					name: file.name,
					contents: update.state.doc.toString()
				});
			}
		});

		// initial state is never use because it's immediately replaced by the darkMode subscription
		// below
		let state: EditorState = EditorState.create({
			doc: file.contents
		});
		const view = new EditorView({
			parent: targetElement,
			state
		});

		darkMode.subscribe((isDarkMode) => {
			const newState = EditorState.create({
				doc: state.doc,
				extensions: [
					basicSetup,
					isDarkMode ? oneDark : [],
					updater,
					langExtension(),
					cursorTooltipField
				]
			});
			view.setState(newState);
			state = newState;
		});

		setValue = (newContent) => {
			if (lastSetValue === newContent) {
				return;
			}

			lastSetValue = newContent;
			view.dispatch(
				view.state.update({
					changes: { from: 0, to: view.state.doc.length, insert: newContent }
				})
			);
		};
	}
</script>

<div use:editor class="mt-8"></div>

<template>
	<div bind:this={explain}>
		<!-- Two buttons, one for explain and one for enhance -->
		<div class="flex justify-center">
			<button
				class="explain-button flex h-12 items-center justify-center gap-2 rounded-s bg-blue-500 px-4 py-2 text-white hover:bg-ablue2-400"
			>
				<span class="font-bold">Explain</span>
				<CircleHelp class="h-4 w-4 text-white dark:text-white" />
			</button>
			<button
				class="improve-button flex h-12 items-center justify-center gap-2 rounded-e border-l border-ablue2-400 bg-blue-500 px-4 py-2 text-white hover:bg-ablue2-400"
			>
				<span class="font-bold">Improve</span>
				<MessageSquareText class="h-4 w-4 text-white dark:text-white" />
			</button>
			<button
				aria-label="Improve"
				class="improve-text hidden h-12 items-center justify-center gap-2 rounded-e border-l border-ablue2-400 bg-blue-500 px-2 py-2 text-white hover:bg-ablue2-400"
			>
				<textarea
					placeholder="How would like to improve this?"
					class="improve-textarea h-8 min-h-8 w-[500px] px-2 pt-1 text-black focus:outline-none"
				></textarea>
			</button>
		</div>
	</div>
</template>
