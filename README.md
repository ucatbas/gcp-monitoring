# gcp-monitoring
Sample GCP Monitoring Client

This repo is a benchmark of fetching and processing metrics, logs, and traces from Google Cloud, making it easier to integrate Google Cloud Traces, Metrics and Logging into your application.

## Features
- **Logging**: Fetch and filter log entries from Google Cloud Logging.
- **Metrics**: Retrieve time series data from Google Cloud Monitoring with customizable filters.
- **Tracing**: Retrieve traces from Google Cloud Trace for analysis and debugging.

## Requirements
- Go 1.18+
- Google Cloud SDK (if you need to authenticate locally)
- Credentials JSON file for a Google Cloud service account with logging, monitoring, and tracing permissions.