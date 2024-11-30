<script lang="ts">
	import { Check, ChevronDown } from '$lib/icons';
	import { assistants } from '$lib/stores';
	import { darkMode } from '$lib/stores';
	import type { Assistant } from '$lib/services';
	import { popover } from '$lib/actions';
	import { fade } from 'svelte/transition';

	const selected = $derived($assistants.find((a) => a.current));

	const { ref, tooltip, toggle } = popover({
		placement: 'bottom-start'
	});

	function icon(a: Assistant | undefined): string {
		if (!a) {
			return '';
		}

		if ($darkMode) {
			return (a.icons.iconDark ? a.icons.iconDark : a.icons.icon) ?? '';
		}
		return a.icons.icon ?? '';
	}

	function collapsedIcon(a: Assistant | undefined): string {
		if (!a) {
			return '';
		}

		if ($darkMode) {
			return a.icons.collapsedDark ?? a.icons.iconDark ?? a.icons.collapsed ?? a.icons.icon ?? '';
		}
		return a.icons.collapsed ?? a.icons.icon ?? '';
	}
</script>

<div class="flex items-center justify-start" transition:fade|global>
	<a
		use:ref
		href={`/${selected?.id ?? ''}`}
		class="flex items-center gap-2"
		onclick={() => {
			if ($assistants.length > 1) {
				toggle();
			} else {
				window.location.href = `/${selected?.id ?? ''}`;
			}
		}}
	>
		{#if collapsedIcon(selected)}
			<img src={collapsedIcon(selected)} alt="assistant icon" class="ml-3 h-8" />
		{:else if selected?.name}
			<div
				class="flex h-8 w-8 items-center justify-center rounded-full bg-gray-200 dark:bg-gray-500"
			>
				{selected?.name ? selected?.name[0].toUpperCase() : '?'}
			</div>
			<span class="font-semibold dark:text-gray-100">{selected?.name ?? ''}</span>
		{/if}
	</a>

	<!-- Dropdown menu -->
	<div
		use:tooltip
		class="mt-4 w-60 divide-y divide-gray-100 rounded-lg bg-white shadow dark:divide-gray-600 dark:bg-gray-700"
	>
		<ul
			class="space-y-1 p-3 text-sm text-gray-700 dark:text-gray-200"
			aria-labelledby="dropdownHelperButton"
		>
			{#each $assistants as assistant}
				<li>
					<a
						href={'/' + assistant.id}
						data-sveltekit-reload
						class="flex rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-600"
					>
						<div class="flex h-5 items-center">
							{#if icon(assistant)}
								<img src={icon(assistant)} alt="assistant icon" class="h-5 w-5 rounded-full" />
							{:else}
								<div
									class="flex h-5 w-5 items-center justify-center rounded-full bg-gray-200 dark:bg-gray-500"
								>
									{assistant.name ? assistant.name[0].toUpperCase() : '?'}
								</div>
							{/if}
						</div>
						<div class="ms-2 text-sm">
							<label
								for="helper-checkbox-1"
								class="flex items-center gap-1 font-medium text-gray-900 dark:text-white"
							>
								{assistant.name}
							</label>
							<p class="text-xs font-normal text-gray-500 dark:text-gray-300">
								{assistant.description}
							</p>
						</div>
						<div class="flex flex-1 justify-end">
							{#if assistant.current}
								<Check class="h-4 w-4 self-center" />
							{/if}
						</div>
					</a>
				</li>
			{/each}
		</ul>
	</div>
	{#if $assistants.length > 1}
		<button class="h-full" onclick={toggle}>
			<ChevronDown
				class="ms-2 h-5 w-5 rounded text-gray-200 hover:bg-gray-100 hover:text-black dark:text-gray-700 hover:dark:bg-gray-700 hover:dark:text-white"
			/>
		</button>
	{/if}
</div>
