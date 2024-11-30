import {
	autoUpdate,
	computePosition,
	flip,
	shift,
	offset,
	type ComputePositionConfig
} from '@floating-ui/dom';
import type { Action, ActionReturn } from 'svelte/action';
import { type Readable, writable, type Writable } from 'svelte/store';

interface Popover extends Readable<boolean> {
	ref: Action;
	tooltip: Action;
	open: Writable<boolean>;
	toggle: () => void;
}

interface PopoverOptions extends Partial<ComputePositionConfig> {
	hover?: boolean;
}

let id = 0;

export default function popover(opts?: PopoverOptions): Popover {
	let ref: HTMLElement;
	let tooltip: HTMLElement;
	const open = writable(false);

	function build(): ActionReturn | void {
		if (!ref || !tooltip) return;

		const selfId = id++;
		document.addEventListener('toolOpen', (e: CustomEvent<string>) => {
			if (e.detail !== selfId.toString()) {
				open.set(false);
			}
		});

		function updatePosition() {
			computePosition(ref, tooltip, {
				placement: 'bottom-end',
				middleware: [
					flip(),
					shift({
						padding: 8
					}),
					offset(8)
				],
				...opts
			}).then(({ x, y }) => {
				Object.assign(tooltip.style, {
					left: `${x}px`,
					top: `${y}px`
				});
			});
		}

		open.subscribe((value) => {
			if (!value) {
				return;
			}

			if (!opts?.hover) {
				const div = document.createElement('div');
				div.classList.add('fixed', 'inset-0', 'z-20', 'cursor-default');
				div.onclick = () => {
					open.set(false);
					div.remove();
				};
				document.body.append(div);
				open.subscribe((value) => {
					if (!value) {
						div.remove();
					}
				});
			}
		});

		tooltip.classList.add('hidden');
		tooltip.classList.add('absolute');

		if (opts?.hover) {
			ref.addEventListener('mouseenter', () => {
				open.set(true);
			});
			ref.addEventListener('mouseleave', () => {
				open.set(false);
			});
		}

		let close: (() => void) | null;
		open.subscribe((value) => {
			if (value) {
				tooltip.classList.remove('hidden');
				updatePosition();
				close = autoUpdate(ref, tooltip, updatePosition);
			} else {
				if (close) {
					close();
				}
				tooltip.classList.add('hidden');
				close = null;
			}
		});

		return {
			destroy: function () {
				if (close) {
					close();
				}
			}
		};
	}

	return {
		ref: (node: HTMLElement) => {
			ref = node;
			return build();
		},
		tooltip: (node: HTMLElement) => {
			tooltip = node;
			return build();
		},
		open,
		subscribe: open.subscribe,
		toggle: () => {
			open.update((value) => {
				if (!value && !opts?.hover) {
					document.dispatchEvent(
						new CustomEvent('toolOpen', {
							detail: id.toString()
						})
					);
				}
				return !value;
			});
		}
	};
}
