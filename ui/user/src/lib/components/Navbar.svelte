<script lang="ts">
	import Logo from '$lib/components/navbar/Logo.svelte';
	import DarkModeToggle from '$lib/components/navbar/DarkModeToggle.svelte';
	import Profile from '$lib/components/navbar/Profile.svelte';
	import Tools from '$lib/components/navbar/Tools.svelte';
	import Files from '$lib/components/navbar/Files.svelte';
	import KnowledgeFile from '$lib/components/navbar/KnowledgeFiles.svelte';
</script>

<nav class="fixed z-30 w-full bg-white dark:bg-transparent dark:bg-gradient-to-b dark:from-black">
	<div class="px-3 py-3 lg:px-5 lg:pl-3">
		<div class="flex items-center justify-between">
			<Logo />
			<div class="flex items-center gap-1 pr-2">
				<KnowledgeFile />
				<Files />
				<Tools />
				<DarkModeToggle />
				<Profile />
			</div>
		</div>
	</div>
</nav>
