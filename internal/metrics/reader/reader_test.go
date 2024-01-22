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
		Hash:     7274857175809454558,
	}

	mp2 := metrics.MetricPoint{
		Name:     "http_requests_total",
		Value:    1,
		LabelSet: map[string]string{"code": "422", "method": "post"},
		Time:     0,
		Hash:     14315283831771060870,
	}

	tests := []Test{
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
						Name: "http_requests_total",
						Metrics: []metrics.MetricPoint{
							mp1, mp2,
						},
						Type:   "gauge",
						Help:   "The total number of HTTP requests.",
						Hashes: map[uint64][]*metrics.MetricPoint{mp1.Hash: {&mp1}, mp2.Hash: {&mp2}},
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
