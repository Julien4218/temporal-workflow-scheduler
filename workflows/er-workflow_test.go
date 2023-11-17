package workflows

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const incidentReaction = "one"
const notAnIncidentReaction = "two"

//todo tests

func Test_hasIsIncidentReactionOnMessage_shouldPassHappyPath(t *testing.T) {
	reactionKeysMap := make(map[string]bool)
	reactionKeysMap[incidentReaction] = true
	reactionCountsMap := make(map[string]int)
	reactionCountsMap[incidentReaction] = 2
	got := hasIsIncidentReactionOnMessage(reactionKeysMap, reactionCountsMap)
	assert.True(t, got)
}

func Test_hasIsIncidentReactionOnMessage_shouldFailIfOnlyHasStarterReaction(t *testing.T) {
	reactionKeysMap := make(map[string]bool)
	reactionKeysMap[incidentReaction] = true
	reactionCountsMap := make(map[string]int)
	reactionCountsMap[incidentReaction] = 1
	got := hasIsIncidentReactionOnMessage(reactionKeysMap, reactionCountsMap)
	assert.False(t, got)
}

func Test_hasIsIncidentReactionOnMessage_shouldFailIfMissingIncidentReaction(t *testing.T) {
	reactionKeysMap := make(map[string]bool)
	reactionCountsMap := make(map[string]int)
	reactionCountsMap[incidentReaction] = 2
	got := hasIsIncidentReactionOnMessage(reactionKeysMap, reactionCountsMap)
	assert.False(t, got)
}

func Test_hasNotAnIncidentReactionOnMessage_shouldPassHappyPath(t *testing.T) {
	reactionKeysMap := make(map[string]bool)
	reactionKeysMap[notAnIncidentReaction] = true
	reactionCountsMap := make(map[string]int)
	reactionCountsMap[notAnIncidentReaction] = 2
	got := hasNotAnIncidentReactionOnMessage(reactionKeysMap, reactionCountsMap)
	assert.True(t, got)
}

func Test_hasNotAnIncidentReactionOnMessage_shouldFailIfOnlyHasStarterReaction(t *testing.T) {
	reactionKeysMap := make(map[string]bool)
	reactionKeysMap[notAnIncidentReaction] = true
	reactionCountsMap := make(map[string]int)
	reactionCountsMap[notAnIncidentReaction] = 1
	got := hasNotAnIncidentReactionOnMessage(reactionKeysMap, reactionCountsMap)
	assert.False(t, got)
}

func Test_hasNotAnIncidentReactionOnMessage_shouldFailIfMissingIncidentReaction(t *testing.T) {
	reactionKeysMap := make(map[string]bool)
	reactionCountsMap := make(map[string]int)
	reactionCountsMap[notAnIncidentReaction] = 2
	got := hasNotAnIncidentReactionOnMessage(reactionKeysMap, reactionCountsMap)
	assert.False(t, got)
}

func Test_updateWorkflowContextOptions_shouldSetStartCloseTimeoutToTenSeconds(t *testing.T) {
	//got := updateWorkflowContextOptions()
	//	todo
}
