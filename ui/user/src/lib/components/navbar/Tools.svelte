<script lang="ts">
	import { Wrench } from '$lib/icons';
	import { tools } from '$lib/stores';
	import { currentAssistant } from '$lib/stores';
	import { ChatService } from '$lib/services';
	import { popover } from '$lib/actions';

	const { ref, tooltip, toggle } = popover({
		placement: 'bottom'
	});

	async function updateTool(enabled: boolean, tool: string | undefined) {
		if (!tool) {
			return;
		}
		if (enabled) {
			tools.set(await ChatService.enableTool($currentAssistant.id, tool));
		} else {
			tools.set(await ChatService.disableTool($currentAssistant.id, tool));
		}
	}
</script>

<button use:ref class="icon-button" onclick={toggle} type="button">
	<Wrench class="h-5 w-5" />
</button>

<!-- Dropdown menu -->
<div
	use:tooltip
	class="w-96 divide-y divide-gray-100 rounded-lg bg-white shadow dark:divide-gray-600 dark:bg-gray-700"
>
	<div class="px-4 py-3 text-sm text-gray-900 dark:text-white">
		Tools
		<p class="mt-1 text-xs font-normal text-gray-500 dark:text-gray-300">
			Enable or disable available tools
		</p>
	</div>
	<ul class="space-y-3 p-3 text-sm" aria-labelledby="dropdownCheckboxButton">
		{#each $tools.items as tool, i}
			{#if !tool.builtin}
				<li>
					<div class="flex">
						{#if tool.icon}
							<img class="h-8 w-8 rounded-md bg-gray-100 p-1" src={tool.icon} alt="message icon" />
						{:else}
							<Wrench class="h-8 w-8 rounded-md bg-gray-100 p-1 text-black" />
						{/if}
						<div class="flex flex-1 px-2">
							<div>
								<label for="checkbox-item-{i}" class="text-sm font-medium dark:text-gray-100"
									>{tool.name}</label
								>
								<p class="text-xs font-normal text-gray-500 dark:text-gray-300">
									{tool.description}
								</p>
							</div>
						</div>
						<input
							id="checkbox-item-{i}"
							type="checkbox"
							checked={tool.enabled}
							onchange={(e) => updateTool(e.currentTarget.checked, tool.id)}
							disabled={tool.builtin || $tools.readonly}
							class="h-4 w-4 self-center rounded border-gray-300 bg-gray-100 focus:ring-blue-500
								  dark:border-gray-500 dark:bg-gray-600 dark:focus:ring-blue-600"
						/>
					</div>
				</li>
			{/if}
		{/each}
	</ul>
</div>
