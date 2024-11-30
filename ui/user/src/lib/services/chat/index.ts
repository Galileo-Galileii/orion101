import { buildMessagesFromProgress } from './messages';
import {
	getProfile,
	listFiles,
	listAssistants,
	listCredentials,
	deleteCredential,
	getFile,
	listKnowledgeFiles,
	uploadKnowledge,
	deleteKnowledgeFile,
	deleteFile,
	invoke,
	newMessageEventSource,
	listTools,
	enableTool,
	disableTool
} from './operations';
import { baseURL } from './http';

export default {
	progressToMessages: buildMessagesFromProgress,
	listFiles,
	listAssistants,
	listKnowledgeFiles,
	listCredentials,
	deleteCredential,
	uploadKnowledge,
	deleteKnowledgeFile,
	deleteFile,
	getFile,
	invoke,
	newMessageEventSource,
	getProfile,
	listTools,
	enableTool,
	disableTool,
	baseURL
};
