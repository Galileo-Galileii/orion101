package handlers

import (
	"bytes"
	"encoding/json"

	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/gz"
	v1 "github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1"
)

func convertKnowledgeSource(agentName string, knowledgeSource v1.KnowledgeSource) types.KnowledgeSource {
	var syncDetails []byte
	if len(knowledgeSource.Status.SyncDetails) > 0 {
		_ = gz.Decompress(&syncDetails, knowledgeSource.Status.SyncDetails)
	}
	return types.KnowledgeSource{
		Metadata:                MetadataFrom(&knowledgeSource),
		KnowledgeSourceManifest: knowledgeSource.Spec.Manifest,
		AgentID:                 agentName,
		State:                   knowledgeSource.PublicState(),
		SyncDetails:             syncDetails,
		Status:                  knowledgeSource.Status.Status,
		Error:                   knowledgeSource.Status.Error,
		LastSyncStartTime:       types.NewTime(knowledgeSource.Status.LastSyncStartTime.Time),
		LastSyncEndTime:         types.NewTime(knowledgeSource.Status.LastSyncEndTime.Time),
		LastRunID:               knowledgeSource.Status.RunName,
	}
}

func checkConfigChanged(oldValue, newValue types.KnowledgeSourceInput) bool {
	oldData, _ := json.Marshal(oldValue)
	newData, _ := json.Marshal(newValue)
	return !bytes.Equal(oldData, newData)
}