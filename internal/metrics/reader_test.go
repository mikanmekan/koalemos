package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Read(t *testing.T) {
	type Test struct {
		desc            string
		literalInput    string
		expectedMetrics map[string]*MetricFamily
		expectedErr     error
	}

	tests := []Test{
		{
			desc: "Help Metadata",
			literalInput: `# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total gauge`,
			expectedMetrics: map[string]*MetricFamily{
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
