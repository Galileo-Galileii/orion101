import type { Config } from 'tailwindcss';
import colors from 'tailwindcss/colors';
import scrollbars from 'tailwind-scrollbar';

const apurple = {
	'50': '#380067',
	'100': '#380067',
	'200': '#380067',
	'300': '#380067',
	'400': '#380067',
	'500': '#380067',
	'600': '#380067',
	'700': '#380067',
	'800': '#380067',
	'900': '#380067',
	'950': '#380067'
};

const apurple2 = {
	'50': '#faf4ff',
	'100': '#f3e6ff',
	'200': '#e8d2ff',
	'300': '#d7aeff',
	'400': '#be7bff',
	'500': '#a549ff',
	'600': '#9125f8',
	'700': '#7c15db',
	'800': '#6a17b2',
	'900': '#57148f',
	'950': '#380067'
};

const black = {
	'50': '#f6f6f6',
	'100': '#e7e7e7',
	'200': '#d1d1d1',
	'300': '#b0b0b0',
	'400': '#888888',
	'500': '#6d6d6d',
	'600': '#5d5d5d',
	'700': '#4f4f4f',
	'800': '#454545',
	'900': '#3d3d3d',
	'950': '#000000'
};

const ablue2 = {
	'50': '#eff5ff',
	'100': '#dce7fd',
	'200': '#c0d5fd',
	'300': '#95bcfb',
	'400': '#6397f7',
	'500': '#4f7ef3',
	'600': '#2953e7',
	'700': '#213fd4',
	'800': '#2135ac',
	'900': '#203188',
	'950': '#182153'
};

const ablue = {
	'50': '#4f7ef3',
	'100': '#4f7ef3',
	'200': '#4f7ef3',
	'300': '#4f7ef3',
	'400': '#4f7ef3',
	'500': '#4f7ef3',
	'600': '#4f7ef3',
	'700': '#4f7ef3',
	'800': '#4f7ef3',
	'900': '#4f7ef3',
	'950': '#4f7ef3'
};

export default {
	content: [
		'./src/**/*.{html,js,svelte,ts}',
		'./node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'
	],

	darkMode: 'class',

	plugins: [scrollbars],

	theme: {
		fontFamily: {
			sans: [
				'Poppins',
				'ui-sans-serif',
				'system-ui',
				'-apple-system',
				'system-ui',
				'Segoe UI',
				'Roboto',
				'Helvetica Neue',
				'Arial',
				'Noto Sans',
				'sans-serif',
				'Apple Color Emoji',
				'Segoe UI Emoji',
				'Segoe UI Symbol',
				'Noto Color Emoji'
			],
			body: [
				'Poppins',
				'ui-sans-serif',
				'system-ui',
				'-apple-system',
				'system-ui',
				'Segoe UI',
				'Roboto',
				'Helvetica Neue',
				'Arial',
				'Noto Sans',
				'sans-serif',
				'Apple Color Emoji',
				'Segoe UI Emoji',
				'Segoe UI Symbol',
				'Noto Color Emoji'
			],
			mono: [
				'ui-monospace',
				'SFMono-Regular',
				'Menlo',
				'Monaco',
				'Consolas',
				'Liberation Mono',
				'Courier New',
				'monospace'
			]
		},

		colors: {
			transparent: 'transparent',
			current: 'currentColor',
			white: colors.white,
			black: colors.black,
			red: colors.red,
			gray: black,
			blue: ablue,
			apurple: apurple,
			apurple2: apurple2,
			ablue: ablue,
			ablue2: ablue2
		}
	}
} as Config;
