import { writable } from 'svelte/store';

export interface EditorFile {
	name: string;
	contents: string;
	buffer: string;
	modified?: boolean;
	selected?: boolean;
}

export default writable([] as EditorFile[]);
