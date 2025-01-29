package main

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func Test_FullWorkflow_Create_S3_Binary_Text(t *testing.T) {
	mocks, ctx, event := prepareHandler(t, "events/event_create_s3_binary_text.json")

	phys, data, err := mocks.syncSopsToSecretsmanager(ctx, event)
	check(err)
	snaps.MatchSnapshot(t, ">>>syncSopsToSecretsmanager", phys, data, err)
}

func Test_FullWorkflow_Create_S3_Binary_Binary(t *testing.T) {
	mocks, ctx, event := prepareHandler(t, "events/event_create_s3_binary_binary.json")

	phys, data, err := mocks.syncSopsToSecretsmanager(ctx, event)
	check(err)
	snaps.MatchSnapshot(t, ">>>syncSopsToSecretsmanager", phys, data, err)
}
