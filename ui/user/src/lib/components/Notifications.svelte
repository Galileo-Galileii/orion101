<script lang="ts" module>
	export class NotificationMessage {
		level: 'info' | 'error';
		message: string;

		constructor(messageOrError: string | Error, level?: 'info' | 'error') {
			if (messageOrError instanceof Error) {
				this.level = 'error';
				this.message = messageOrError.message;
			} else {
				this.level = level || 'info';
				this.message = messageOrError;
			}
		}
	}
</script>

<script lang="ts">
	import { CircleCheck } from '$lib/icons';
	import { CircleX } from '$lib/icons';
	import { X } from '$lib/icons';
	import { errors, profile } from '$lib/stores';

	let notifications: NotificationMessage[] = $state([]);
	let div: HTMLElement;

	export function addNotification(notification: NotificationMessage) {
		notifications = [...notifications, notification];
		setTimeout(() => {
			notifications = notifications.slice(1);
		}, 5000);
	}

	$effect(() => {
		if ($profile.loaded && div.classList.contains('hidden')) {
			div.classList.remove('hidden');
			div.classList.add('flex');
		}
	});
</script>

<div bind:this={div} class="absolute bottom-0 right-0 z-50 mb-20 mr-4 hidden flex-col">
	{#each $errors as error, i}
		<div
			class="mb-4 flex w-full max-w-xs items-center rounded-lg bg-white p-4 text-gray-500 shadow dark:bg-gray-800 dark:text-white"
		>
			<div
				class="text-green-500 bg-green-100 dark:bg-green-800 dark:text-green-200 inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg"
			>
				<CircleX />
				<span class="sr-only">Check icon</span>
			</div>
			<div class="ms-3 text-sm font-normal">{error.message}</div>
			<button
				type="button"
				onclick={() => errors.remove(i)}
				class="-mx-1.5 -my-1.5 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-white p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-900 focus:ring-2 focus:ring-gray-300 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-700 dark:hover:text-white"
			>
				<span class="sr-only">Close</span>
				<X />
			</button>
		</div>
	{/each}
	{#each notifications as notification, i}
		<div
			class="mb-4 flex w-full max-w-xs items-center rounded-lg bg-white p-4 text-gray-500 shadow dark:bg-gray-800 dark:text-white"
		>
			<div
				class="text-green-500 bg-green-100 dark:bg-green-800 dark:text-green-200 inline-flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg"
			>
				{#if notification.level === 'error'}
					<CircleX />
				{:else}
					<CircleCheck />
				{/if}
				<span class="sr-only">Check icon</span>
			</div>
			<div class="ms-3 text-sm font-normal">{notification.message}</div>
			<button
				type="button"
				onclick={() => (notifications = notifications.filter((_, index) => index !== i))}
				class="-mx-1.5 -my-1.5 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-white p-1.5 text-gray-400 hover:bg-gray-100 hover:text-gray-900 focus:ring-2 focus:ring-gray-300 dark:bg-gray-800 dark:text-white dark:hover:bg-gray-700 dark:hover:text-white"
			>
				<span class="sr-only">Close</span>
				<X />
			</button>
		</div>
	{/each}
</div>
