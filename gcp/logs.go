package gcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func NewLogClient(ctx context.Context, projectId string, credentialsFile string) *logadmin.Client {
	client, err := logadmin.NewClient(ctx, projectId, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func GetEntries(ctx context.Context, client *logadmin.Client, logName string) ([]*logging.Entry, error) {
	var entries []*logging.Entry
	lastHour := time.Now().Add(-72 * time.Hour).Format(time.RFC3339)

	iter := client.Entries(ctx,
		logadmin.Filter(fmt.Sprintf(`logName = "%s" AND timestamp > "%s"`, logName, lastHour)),
		logadmin.NewestFirst(),
	)

	// Fetch the most recent 5 entries.
	for len(entries) < 5 {
		entry, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
