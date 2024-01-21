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

	tests := []Test{
		{
			desc: "Help Metadata",
			literalInput: `978595200
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total gauge
http_requests_total{method="post",code="200"} 1027`,
			expectedMetrics: &metrics.MetricFamiliesTimeGroup{
				Time: 978595200,
				Families: map[string]*metrics.MetricFamily{
					"http_requests_total": {
						Name: "http_requests_total",
						Metrics: []metrics.MetricPoint{
							{
								Name:     "http_requests_total",
								Value:    1027,
								LabelSet: map[string]string{"code": "200", "method": "post"},
								Time:     0,
							},
						},
						Type: "gauge",
						Help: "The total number of HTTP requests.",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		reader := NewReader()

		bs := []byte(tc.literalInput)
		byteReader := bytes.NewReader(bs)

		res, err := reader.Read(byteReader)
		assert.Equal(t, tc.expectedMetrics, res)
		assert.Equal(t, tc.expectedErr, err)
	}
}
