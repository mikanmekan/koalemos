# Koalemos Data Model

### Koalemos Metric

Koalemos uses a similar metrics data model to Prometheus & OpenTSDB.

#### Name string

#### Value float64

#### Timestamp int64

#### LabelSet {LabelName: LabelValue...} (map\[string\]string)

-----

### Koalemos Metric Ingestion Format

This is the format expected for the Metrics Payload.

#### Example

#HELP http_requests_total Total number of HTTP requests
#TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 100
http_requests_total{method="POST",status="200"} 50

#HELP http_request_duration_seconds Duration of HTTP requests in seconds
#TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.1"} 5
http_request_duration_seconds_bucket{le="0.5"} 20
http_request_duration_seconds_bucket{le="1.0"} 50
http_request_duration_seconds_bucket{le="5.0"} 100
http_request_duration_seconds_bucket{le="+Inf"} 150
http_request_duration_seconds_sum 45.0
http_request_duration_seconds_count 150

#HELP app_version_info Application version information
#TYPE app_version_info gauge
app_version_info{version="1.0.0",commit="abcdef123",build_time="2023-10-01T12:00:00Z"} 1\EOF