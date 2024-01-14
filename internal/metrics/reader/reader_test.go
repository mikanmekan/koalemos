package reader

import (
	"testing"

	"github.com/mikanmekan/koalemos/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_Read(t *testing.T) {
	type Test struct {
		desc            string
		literalInput    string
		expectedMetrics map[string]*metrics.MetricFamily
		expectedErr     error
	}

	tests := []Test{
		{
			desc: "Help Metadata",
			literalInput: `# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total gauge
http_requests_total{method="post",code="200"} 1027`,
			expectedMetrics: map[string]*metrics.MetricFamily{
				"http_requests_total": {
					Name:    "http_requests_total",
					Metrics: nil,
					Type:    "gauge",
					Help:    "The total number of HTTP requests.",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		reader := NewReader()

		bs := []byte(tc.literalInput)

		res, err := reader.Read(bs)
		assert.Equal(t, tc.expectedMetrics, res)
		assert.Equal(t, tc.expectedErr, err)
	}
}
