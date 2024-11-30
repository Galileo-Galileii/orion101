<script lang="ts">
	import ProfileIcon from '$lib/components/profile/ProfileIcon.svelte';
	import { profile, currentAssistant } from '$lib/stores';
	import { popover } from '$lib/actions';
	import { ChatService, type CredentialList } from '$lib/services';
	import { Trash } from '$lib/icons';
	import Loading from '$lib/icons/Loading.svelte';

	const { ref, tooltip, toggle } = popover({
		placement: 'bottom-end'
	});

	let credPromise = $state<Promise<CredentialList>>();

	async function deleteCred(name: string) {
		await ChatService.deleteCredential($currentAssistant.id, name);
		credPromise = ChatService.listCredentials($currentAssistant.id);
	}

	function loadCredsAndToggle() {
		if ($currentAssistant.id) {
			credPromise = ChatService.listCredentials($currentAssistant.id);
		}
		toggle();
	}
</script>

<!-- Profile -->
<div class="ml-1 flex items-center" use:ref>
	<button
		onclick={loadCredsAndToggle}
		type="button"
		class="flex rounded-full bg-gray-800 text-sm focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600"
	>
		<span class="sr-only">Open user menu</span>
		<ProfileIcon />
	</button>
	<!-- Dropdown menu -->
	<div
		use:tooltip
		class="mt-2 list-none divide-y divide-gray-100 rounded bg-white text-base shadow dark:divide-gray-600 dark:bg-gray-700"
	>
		<div class="px-4 py-3" role="none">
			<p class="truncate text-sm font-medium text-gray-900 dark:text-white" role="none">
				{$profile.email || 'Anonymous'}
			</p>
		</div>
		<div class="px-4 py-3" role="none">
			{#if credPromise !== undefined}
				{#await credPromise}
					<p class="mb-1 truncate text-sm text-gray-900 dark:text-white" role="none">
						Credentials <Loading class="mb-0.5 ms-1 h-3 w-3" />
					</p>
				{:then credentials}
					<p class="mb-1 truncate text-sm text-gray-900 dark:text-white" role="none">Credentials</p>
					{#if credentials?.items.length > 0}
						<ul class="py-1" role="none">
							{#each credentials.items as cred}
								<li class="flex">
									<span class="flex-1 py-2 text-sm text-black dark:text-white">{cred.name}</span>
									<button>
										<Trash
											class="h-5 w-5 text-gray-400"
											onclick={() => {
												deleteCred(cred.name);
											}}
										/>
									</button>
								</li>
							{/each}
						</ul>
					{:else}
						<span class="flex-1 py-2 text-sm text-black dark:text-white">No credentials</span>
					{/if}
				{/await}
			{/if}
		</div>
		<ul class="py-1" role="none">
			{#if $profile.role === 1}
				<li>
					<a
						href="/admin/"
						rel="external"
						class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-white dark:hover:bg-gray-600 dark:hover:text-white"
						role="menuitem">Settings</a
					>
				</li>
			{/if}
			{#if $profile.email}
				<li>
					<a
						href="/oauth2/sign_out?rd=/"
						rel="external"
						class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-white dark:hover:bg-gray-600 dark:hover:text-white"
						role="menuitem">Sign out</a
					>
				</li>
			{/if}
		</ul>
	</div>
</div>
