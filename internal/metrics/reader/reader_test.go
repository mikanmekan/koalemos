package reader

import (
	"bytes"
	"testing"

	"github.com/mikanmekan/koalemos/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_Read(t *testing.T) {
	type Test struct {
		desc            string
		literalInput    string
		expectedMetrics *metrics.MetricFamiliesTimeGroup
		expectedErr     error
	}

	mp1 := metrics.MetricPoint{
		Name:     "http_requests_total",
		Value:    1027,
		LabelSet: map[string]string{"code": "200", "method": "post"},
		Time:     0,
		Hash:     15123803854892908114,
	}

	mp2 := metrics.MetricPoint{
		Name:     "http_requests_total",
		Value:    1,
		LabelSet: map[string]string{"code": "422", "method": "post"},
		Time:     0,
		Hash:     15123803854892908114,
	}

	mp3 := metrics.MetricPoint{
		Name:     "http_requests_latency_ms",
		Value:    70000,
		LabelSet: map[string]string{"code": "200", "method": "post"},
		Time:     0,
		Hash:     2862631026593619755,
	}

	tests := []Test{
		{
			desc: "[POSITIVE] successfully read two metric points from two families",
			literalInput: `978595200
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total gauge
# HELP http_requests_latency_ms The total latency of HTTP requests.
# TYPE http_requests_latency_ms gauge
http_requests_total{method="post",code="200"} 1027
http_requests_total{method="post",code="422"} 1
http_requests_latency_ms{method="post",code="200"} 70000`,
			expectedMetrics: &metrics.MetricFamiliesTimeGroup{
				Time: 978595200,
				Families: map[string]*metrics.MetricFamily{
					"http_requests_total": {
						Definition: metrics.MetricDefinition{
							Name: "http_requests_total",
							Type: "gauge",
							Help: "The total number of HTTP requests.",
						},
						HashedMetrics: map[uint64][]*metrics.MetricPoint{mp1.Hash: {&mp1, &mp2}},
					},
					"http_requests_latency_ms": {
						Definition: metrics.MetricDefinition{
							Name: "http_requests_latency_ms",
							Type: "gauge",
							Help: "The total latency of HTTP requests.",
						},
						HashedMetrics: map[uint64][]*metrics.MetricPoint{mp3.Hash: {&mp3}},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "[POSITIVE] successfully read two metric points",
			literalInput: `978595200
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total gauge
http_requests_total{method="post",code="200"} 1027
http_requests_total{method="post",code="422"} 1`,
			expectedMetrics: &metrics.MetricFamiliesTimeGroup{
				Time: 978595200,
				Families: map[string]*metrics.MetricFamily{
					"http_requests_total": {
						Definition: metrics.MetricDefinition{
							Name: "http_requests_total",
							Type: "gauge",
							Help: "The total number of HTTP requests.",
						},
						HashedMetrics: map[uint64][]*metrics.MetricPoint{mp1.Hash: {&mp1, &mp2}},
					},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "[NEGATIVE] input metrics payload has duplicated label set, i.e. two method=post & code=200 entries",
			literalInput: `978595200
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total gauge
http_requests_total{method="post",code="200"} 1027
http_requests_total{method="post",code="200"} 1`,
			expectedErr: metrics.ErrDuplicateMetricLabelSet,
		},
	}

	for _, tc := range tests {
		reader := NewReader()

		bs := []byte(tc.literalInput)
		byteReader := bytes.NewReader(bs)

		res, err := reader.Read(byteReader)

		assert.ErrorIs(t, err, tc.expectedErr)
		// We don't care for the results if we encounter an error while reading.
		// Read errors will be fully forgotten about.
		if err == nil {
			assert.Equal(t, tc.expectedMetrics, res)
		}
	}
}
