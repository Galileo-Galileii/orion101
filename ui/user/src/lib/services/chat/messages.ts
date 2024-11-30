import type { Progress, InputMessage, Message, Messages, Explain } from './types';
import editorFiles from '$lib/stores/editorfiles';

const orion101AIIcon = 'Orion101';
const profileIcon = 'Profile';

function toMessageFromInput(s: string): string {
	try {
		const input = JSON.parse(s) as InputMessage;
		if (input.type === 'orion101-prompt') {
			return input.prompt;
		}
	} catch {
		// ignore error
	}
	return s;
}

function setFileContent(name: string, content: string, full: boolean = false) {
	editorFiles.update((files) => {
		const existing = files.find((f) => f.name === name);
		if (existing) {
			if (full) {
				existing.contents = content;
			} else if (content.length < existing.contents.length) {
				existing.contents = content + existing.contents.slice(content.length);
			} else {
				existing.contents = content;
			}
		} else {
			files.push({
				name: name,
				contents: content,
				buffer: ''
			});
		}

		// select the file
		files.forEach((f) => {
			f.selected = f.name === name;
		});

		return files;
	});
}

function reformatInputMessage(msg: Message) {
	try {
		const input = JSON.parse(msg.message.join('')) as {
			prompt?: string;
			explain?: Explain;
			improve?: Explain;
		};
		if (input.prompt) {
			msg.message = [input.prompt];
		}
		if (input.explain) {
			msg.explain = input.explain;
			msg.message = ['Explain'];
		} else if (input.improve) {
			msg.explain = input.improve;
		}
	} catch {
		// ignore error
	}
}

function reformatWriteMessage(msg: Message, last: boolean) {
	msg.icon = 'stock:Pencil';
	msg.done = !last || msg.toolCall;
	msg.sourceName = msg.done ? 'Wrote to Workspace' : 'Writing to Workspace';
	let content = msg.message.join('').trim();
	if (!content.endsWith('"}')) {
		content += '"}';
	}
	try {
		const obj = JSON.parse(content);
		if (obj.filename) {
			msg.file = {
				filename: obj.filename,
				content: ''
			};
			msg.file.filename = obj.filename;
			if (obj.content) {
				msg.file.content = obj.content;
			}
		}
		msg.message = [];
	} catch {
		// ignore error
	}

	if (last && msg.file?.filename && msg.file?.content) {
		setFileContent(msg.file.filename, msg.file.content, msg.toolCall);
	}
}

export function buildMessagesFromProgress(progresses: Progress[]): Messages {
	const messages = toMessages(progresses);

	// Post Process for much more better-ness
	messages.messages.forEach((item, i) => {
		if (item.tool && item.sourceName == 'workspace_write') {
			reformatWriteMessage(item, i == messages.messages.length - 1);
			return;
		} else if (item.sent) {
			reformatInputMessage(item);
		}

		if (item.toolInput) {
			item.message = ['Preparing information...'];
		} else if (item.toolCall) {
			item.message = ['Calling...'];
		}

		// For all but last message
		if (i < messages.messages.length - 1) {
			if (item.oauthURL) {
				item.ignore = true;
			}
			if (item.tool) {
				item.done = true;
				item.message = [];
			}
		}
	});

	return messages;
}

function toMessages(progresses: Progress[]): Messages {
	let messages: Message[] = [];
	let lastRunID: string | undefined;
	let inProgress = false;

	for (const [i, progress] of progresses.entries()) {
		if (progress.error) {
			// Errors are handled as events, so we can just ignore them here
			continue;
		}

		if (progress.runID) {
			lastRunID = progress.runID;
			inProgress = true;
		} else {
			// if it doesn't have a runID we don't know what do to with it, so ignore
			continue;
		}

		if (progress.runComplete) {
			lastRunID = progress.runID;
			inProgress = false;
			for (const message of messages) {
				if (message.runID === progress.runID) {
					message.done = true;
				}
			}
		} else {
			inProgress = true;
		}

		if (progress.error) {
			// Errors are handled as events, so we can just ignore them here
			continue;
		} else if (progress.waitingOnModel) {
			if (i === progresses.length - 1) {
				// Only add if it's the last one
				messages.push(newWaitingOnModelMessage(progress));
			}
		} else if (progress.prompt && progress.prompt.metadata?.authType === 'oauth') {
			messages.push(newOAuthMessage(progress));
		} else if (progress.input) {
			// delete the current runID, this is to avoid duplicate messages
			messages = messages.filter((m) => m.runID !== progress.runID);
			messages.push(newInputMessage(progress));
		} else if (progress.content) {
			const found = messages.findLast(
				(m) => m.contentID === progress.contentID && progress.contentID
			);
			if (found) {
				found.message.push(progress.content);
				found.time = new Date(progress.time);
			} else {
				messages.push(newContentMessage(progress));
			}
		} else if (progress.toolInput) {
			const found = messages.findLast(
				(m) => m.contentID === progress.contentID && progress.contentID
			);
			if (found) {
				if (progress.toolInput.input) {
					found.message.push(progress.toolInput.input);
				}
			} else {
				messages.push(newContentMessage(progress));
			}
		} else if (progress.toolCall) {
			// once we see a toolCall ignore all previous toolInputs
			for (const msg of messages) {
				if (msg.runID === progress.runID && msg.toolInput) {
					msg.ignore = true;
				}
			}
			messages.push(newContentMessage(progress));
		}
	}

	return {
		lastRunID,
		messages,
		inProgress
	};
}

function newInputMessage(progress: Progress): Message {
	return {
		runID: progress.runID || '',
		time: new Date(progress.time),
		icon: profileIcon,
		sourceName: 'You',
		sent: true,
		message: [toMessageFromInput(progress.input || '')],
		done: true
	};
}

function newOAuthMessage(progress: Progress): Message {
	// prompt will not be undefined at this point
	return {
		runID: progress.runID || '',
		time: new Date(progress.time),
		icon: progress.prompt?.metadata?.icon || orion101AIIcon,
		sourceName: progress.prompt?.name || 'Orion101',
		sourceDescription: progress.prompt?.description,
		oauthURL: progress.prompt?.metadata?.authURL || '',
		message: progress.prompt?.message ? [progress.prompt?.message] : []
	};
}

function newWaitingOnModelMessage(progress: Progress): Message {
	return {
		runID: progress.runID || '',
		time: new Date(progress.time),
		icon: orion101AIIcon,
		sourceName: 'Orion101',
		message: ['Thinking really hard...']
	};
}

function newContentMessage(progress: Progress): Message {
	const result: Message = {
		time: new Date(progress.time),
		runID: progress.runID || '',
		icon: orion101AIIcon,
		sourceName: 'Orion101',
		message: [progress.toolInput?.input ?? progress.content],
		contentID: progress.contentID
	};

	if (progress.toolInput) {
		if (progress.toolInput.name) {
			result.sourceName = progress.toolInput.name;
		}
		result.sourceDescription = progress.toolInput.description;
		if (progress.toolInput.metadata?.icon) {
			result.icon = progress.toolInput.metadata.icon;
		}
		result.message = progress.toolInput.input ? [progress.toolInput.input] : [];
		result.toolInput = true;
		result.tool = true;
	}

	if (progress.toolCall) {
		if (progress.toolCall.name) {
			result.sourceName = progress.toolCall.name;
		}
		result.sourceDescription = progress.toolCall.description;
		if (progress.toolCall.metadata?.icon) {
			result.icon = progress.toolCall.metadata.icon;
		}
		result.message = progress.toolCall.input ? [progress.toolCall.input] : [];
		result.toolCall = true;
		result.tool = true;
	}

	return result;
}
