# Koalemos
A timeseries database for metrics

## Ingestor
Listens for and stores metrics.

!Periodically writes blocks out to disk || obj. storage

## Blocks
!Blocks contain metrics received during a span of time, by default 2 hours.