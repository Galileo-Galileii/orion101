<script>
	import { Moon, Sun } from '$lib/icons';
	import { darkMode } from '$lib/stores';
</script>

<button
	type="button"
	onclick={() => {
		darkMode.set(!$darkMode);
	}}
	class="icon-button"
>
	{#if $darkMode}
		<Sun class="h-5 w-5" />
	{:else}
		<Moon class="h-5 w-5" />
	{/if}
</button>
