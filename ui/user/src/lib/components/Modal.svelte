<script lang="ts">
	import { CircleAlert, X } from '$lib/icons';

	interface Props {
		show: boolean;
		msg?: string;
		onsuccess: () => void;
		oncancel: () => void;
	}

	let { show = false, msg = 'OK?', onsuccess, oncancel }: Props = $props();
</script>

<div
	tabindex="-1"
	class:hidden={!show}
	class:flex={show}
	class="fixed left-0 right-0 top-0 z-50 h-[calc(100%-1rem)] max-h-full w-full items-center
 justify-center overflow-y-auto overflow-x-hidden bg-black bg-opacity-5 md:inset-0"
>
	<div class="relative max-h-full w-full max-w-md p-4">
		<div class="relative rounded-lg bg-white shadow dark:bg-gray-700">
			<button
				type="button"
				onclick={oncancel}
				class="absolute end-2.5 top-3 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-transparent text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900 dark:hover:bg-gray-600 dark:hover:text-white"
				data-modal-hide="popup-modal"
			>
				<X class="h-5 w-5" />
				<span class="sr-only">Close modal</span>
			</button>
			<div class="p-4 text-center md:p-5">
				<CircleAlert class="mx-auto mb-4 h-12 w-12 text-gray-400 dark:text-gray-100" />
				<h3 class="mb-5 text-lg font-normal text-black dark:text-gray-100">{msg}</h3>
				<button
					onclick={onsuccess}
					type="button"
					class="inline-flex items-center rounded-lg bg-red-600 px-5 py-2.5 text-center text-sm font-medium text-white hover:bg-red-800 focus:outline-none focus:ring-4 focus:ring-red-300 dark:focus:ring-red-800"
				>
					Yes, I'm sure
				</button>
				<button
					onclick={oncancel}
					type="button"
					class="ms-3 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-medium text-gray-900 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white dark:focus:ring-gray-700"
					>No, cancel</button
				>
			</div>
		</div>
	</div>
</div>
