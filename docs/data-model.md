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