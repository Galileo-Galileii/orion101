import { type Readable, type Subscriber, type Unsubscriber, type Writable } from 'svelte/store';

import errorStore from './errors';

function toError(error: unknown): Error {
	if (error instanceof Error) {
		return error;
	}
	return new Error(String(error));
}

export type Initializer = () => void | Promise<void>;
export type ErrorCallback = (e: Error) => void;

export function storeWithInit<V, T extends Readable<V> | Writable<V>>(
	target: T,
	init: Initializer,
	errorCallback?: ErrorCallback
): T {
	let initialized = false;
	let inflight = false;

	function handleError(e: unknown) {
		if (errorCallback) {
			errorCallback(toError(e));
		}
		errorStore.append(toError(e));
	}

	function initialize(): void {
		if (initialized || inflight || typeof window === 'undefined') {
			return;
		}

		try {
			const promise = init();
			if (promise instanceof Promise) {
				inflight = true;
				promise
					.then(() => {
						initialized = true;
					})
					.catch(handleError)
					.finally(() => {
						inflight = false;
					});
			} else {
				initialized = true;
			}
		} catch (e) {
			handleError(e);
		}
	}

	return {
		...target,
		subscribe(run: Subscriber<V>, invalidate?: () => void): Unsubscriber {
			initialize();
			return target.subscribe(run, invalidate);
		}
	};
}
