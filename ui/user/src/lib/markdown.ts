import { micromark } from 'micromark';
import { gfmHtml, gfm } from 'micromark-extension-gfm';

export function toHTMLFromMarkdown(markdown: string): string {
	return micromark(markdown, {
		extensions: [gfm()],
		htmlExtensions: [gfmHtml()],
		allowDangerousHtml: true
	});
}
