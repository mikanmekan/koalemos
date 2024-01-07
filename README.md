# Koalemos
A timeseries database for metrics

## How does it work?

                             0..N
                             ┌──────────────┐
                      ┌──────►    Client    │
                      │      │              │
                      │      └──────┬───────┘
                Query metrics       │
                      │             │
                      │             │Metrics payload
                      │             │
                      │      1      │   Inspects metrics payloads, and distributes metrics load across ingestors.
                      │      ┌──────▼──────┐ Remembers where timeseries land for smart timeseries/query placement
                      └──────►             │
                             │   Doorman   │
                    ┌────────┤             ├────────┐
                    │        └─────────────┘        │
                    │                               │
     Ingestors 1..M │                               │ Metrics payload
    ── ── ── ── ── ─┤ ── ── ── ── ── ── ── ── ── ── ├─ ── ── ── ── ── ── ── ── ── ── ── ── ── ── ── ── ── ── ── ── ──
                    │                               │
               ┌────▼───────┐                ┌──────▼─────┐
               │  Ingestor  │                │  Ingestor  │
               │            │                │            │
               └────────────┘                └────────────┘





      ┌─Ingestor────────────────────────────────┬─────────────────────────────────┐
      │                                         │ 
      │                                         │  Serialise blocks of timeseries │
      │           ┌───────────┐                                                   │
      │           │           │                 │  Write to disk                  │
      │           │   Parse   │                 │                                 │
      │           │           │                    Write to remote (obj storage)  │
      │           └───────────┘                                                   │
      │                                         │  Ingest logs?                   │
      │                                         │                                 │
      │    ┌───────────────────────────┐        │                                 │
      │    │                           │                                          │
      │    │ Maintain chronological    │        │                                 │
      │    │ block for each timeseries │        │                                 │
      │    │                           │        │                                 │
      │    └───────────────────────────┘                                          │
      │                                         │                                 │
      │                                         │                                 │
      │                                                                           │
      │                                                                           │
      │                                         │                                 │
      │                                         │                                 │
      │                                                                           │
      └───────────────────────────────────────────────────────────────────────────┘


## Doorman
Forwards metrics payloads to recipient ingestors.

Queries ingestors(?) Would make sense to have a separate querier component which will eventually deal with reading out blocks from disk & obj storage.

## Ingestor
Receives and parses metrics, and stores them in memory. These need to be fast to retrieve chronologically to support range queries effectively.

Periodically writes blocks out to disk || obj. storage

## Blocks
Blocks contain metrics received during a span of time, by default 2 hours.