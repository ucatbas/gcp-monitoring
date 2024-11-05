package gcp

import (
	"context"
	"fmt"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	monitoringpb "cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MetricData struct {
	SubjectID   string       `json:"subject_id,omitempty"`
	SubjectType string       `json:"subject_type,omitempty"`
	MetricType  string       `json:"metric_type"`
	Points      []PointEntry `json:"points"`
}

type PointEntry struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	Count int64 `json:"count"`
}

func NewMetricClient(ctx context.Context, credentialsFile string) *monitoring.MetricClient {
	// Creates a client.
	client, err := monitoring.NewMetricClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close() // Ensure client is closed on exit

	return client
}

// Extract metric data and aggregate counts per minute
func ExtractMetricData(response *monitoringpb.TimeSeries, metricType string) MetricData {
	subjectID := response.Metric.Labels["subject_id"]
	subjectType := response.Metric.Labels["subject_type"]
	var points []PointEntry

	// Check if there are points to process
	if len(response.Points) == 0 {
		return MetricData{SubjectID: subjectID, SubjectType: subjectType, MetricType: metricType, Points: points}
	}

	// Iterate through each point in reverse order
	for i := len(response.Points) - 1; i >= 0; i-- {
		currentPoint := response.Points[i]
		startTime := currentPoint.Interval.StartTime.AsTime()
		endTime := currentPoint.Interval.EndTime.AsTime()
		count := currentPoint.Value.GetDistributionValue().Count

		// If it's the last point, just add it directly
		if i == len(response.Points)-1 {
			points = append(points, PointEntry{
				StartTime: startTime,
				EndTime:   endTime,
				Count:     count,
			})
			continue
		}

		// For subsequent points, check for the same start time
		previousPoint := response.Points[i+1]
		if startTime.Equal(previousPoint.Interval.StartTime.AsTime()) {
			// Calculate the difference in counts
			countDiff := count - previousPoint.Value.GetDistributionValue().Count
			points = append(points, PointEntry{
				StartTime: previousPoint.Interval.EndTime.AsTime(),
				EndTime:   endTime,
				Count:     countDiff,
			})
		} else {
			// Just add the current point as is
			points = append(points, PointEntry{
				StartTime: startTime,
				EndTime:   endTime,
				Count:     count,
			})
		}
	}

	return MetricData{
		SubjectID:   subjectID,
		SubjectType: subjectType,
		MetricType:  metricType,
		Points:      points,
	}
}

// Creates a request to list time series data for the given metric type
func CreateListTimeSeriesRequest(projectId string, metricType string) *monitoringpb.ListTimeSeriesRequest {
	startTime := time.Now().Add(-240 * time.Hour)
	endTime := time.Now()

	return &monitoringpb.ListTimeSeriesRequest{
		Name:   fmt.Sprintf("projects/%s", projectId),
		Filter: fmt.Sprintf(`metric.type="%s"`, metricType),
		Interval: &monitoringpb.TimeInterval{
			StartTime: timestamppb.New(startTime),
			EndTime:   timestamppb.New(endTime),
		},

		View: monitoringpb.ListTimeSeriesRequest_FULL,
	}
}
