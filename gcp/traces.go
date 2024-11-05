package gcp

import (
	"context"
	"log"

	trace "cloud.google.com/go/trace/apiv1"
	"cloud.google.com/go/trace/apiv1/tracepb"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func NewTraceClient(ctx context.Context, credentialsFile string) *trace.Client {
	client, err := trace.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client
}

func FetchTraces(ctx context.Context, client *trace.Client, projectId string) ([]*tracepb.Trace, error) {

	var traces []*tracepb.Trace

	req := &tracepb.ListTracesRequest{
		ProjectId: projectId,
	}

	iter := client.ListTraces(ctx, req)
	// Fetch the most recent 5 entries.
	for len(traces) < 5 {
		trace, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		traces = append(traces, trace)
	}
	return traces, nil
}
