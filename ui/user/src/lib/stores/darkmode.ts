import theme from './theme';
import preferredTheme from './preferredtheme';
import { derived } from 'svelte/store';

const store = derived([theme, preferredTheme], ($values) => {
	const [theme, preferredTheme] = $values;
	if (theme === 'system') {
		return preferredTheme === 'dark';
	}
	return theme == 'dark';
});

function set(darkMode: boolean) {
	theme.set(darkMode ? 'dark' : 'light');
}

export default {
	subscribe: store.subscribe,
	set
};
