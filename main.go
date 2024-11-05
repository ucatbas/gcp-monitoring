package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"monitoring/gcp"

	"github.com/joho/godotenv"
)

func getMetrics(ctx context.Context, projectId string, credentialsFile string, metrics []string) {
	client := gcp.NewMetricClient(ctx, credentialsFile)
	for _, metricType := range metrics {
		req := gcp.CreateListTimeSeriesRequest(projectId, metricType)
		it := client.ListTimeSeries(ctx, req)
		fmt.Printf("Fetching data for metric type: %s\n", metricType)
		for {
			resp, err := it.Next()
			if err != nil {
				if err.Error() == "no more items in iterator" {
					break
				}
				log.Fatalf("Failed to retrieve time series data for %s: %v", metricType, err)
			}
			data := gcp.ExtractMetricData(resp, metricType)
			fmt.Println("-------------------")
			fmt.Printf("Metric Type: %s\n", data.MetricType)
			if data.SubjectID != "" {
				fmt.Printf("Subject ID: %s, Subject Type: %s\n", data.SubjectID, data.SubjectType)
			}
			for _, point := range data.Points {
				if point.Count > 0 {
					fmt.Printf("Start Time: %s, End Time: %s, Count: %d\n", point.StartTime.Format(time.RFC3339), point.EndTime.Format(time.RFC3339), point.Count)
				}
			}
			fmt.Println("-------------------")
		}
	}

	fmt.Println("Done retrieving time series data.")
}

func getLogs(ctx context.Context, projectId string, credentialsFile string, logName string) {
	client := gcp.NewLogClient(ctx, projectId, credentialsFile)
	defer client.Close() // Close client here after we're done with it

	entries, err := gcp.GetEntries(ctx, client, logName)
	if err != nil {
		log.Fatalf("Failed to retrieve entries: %v", err)
	}
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

func getTraces(ctx context.Context, projectId string, credentialsFile string) {
	client := gcp.NewTraceClient(ctx, credentialsFile)
	defer client.Close() // Close client here after we're done with it

	traces, err := gcp.FetchTraces(ctx, client, projectId)
	if err != nil {
		log.Fatalf("Failed to retrieve traces: %v", err)
	}
	for _, trace := range traces {
		fmt.Println(trace)
	}

}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Access environment variables
	projectId := os.Getenv("PROJECT_ID")
	credentialsFile := os.Getenv("PROJECT_CREDENTIALS")

	ctx := context.Background()

	// retrieve logs
	logName := "log-name"
	getLogs(ctx, projectId, credentialsFile, logName)

	// retrieve metrics
	metrics := []string{
		"metric1",
		"metric2",
	}
	getMetrics(ctx, projectId, credentialsFile, metrics)

	// retrieve traces
	getTraces(ctx, projectId, credentialsFile)
}
