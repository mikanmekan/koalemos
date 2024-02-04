package metrics

type Block struct {
	// metrics stored as a hash of metric definition + labelset
	// hash(metricName, labelset) -> []timeseries.
	// hash returns slice of timeseries because we are hash collision cognizant.
	metrics map[uint64][]MetricFamilyTimeSeries
}

func NewBlock() *Block {
	return &Block{metrics: make(map[uint64][]MetricFamilyTimeSeries)}
}

func (b *Block) AddMetricFamily(metricFamily *MetricFamily, hashed uint64) {
	timeSeries := b.metrics[hashed]

	for _, metrics := range metricFamily.HashedMetrics {
		for _, mp := range metrics {
			if len(timeSeries) == 0 {
				b.addNewTimeSeries(metricFamily, mp, hashed)
			} else {
				foundTs := false
				for _, ts := range timeSeries {
					if mp.MetadataEquals(ts.ToMetricPoint()) {
						ts.metrics = append(ts.metrics, *mp)
						break
					}
				}
				if !foundTs {
					b.addNewTimeSeries(metricFamily, mp, hashed)
				}
			}
		}
	}
}

func (b *Block) addNewTimeSeries(metricFamily *MetricFamily, mp *MetricPoint, hashed uint64) {
	ts := *metricFamily.ToMetricTimeSeries()
	ts.LabelSet = mp.LabelSet
	ts.metrics = append(ts.metrics, *mp)
	b.metrics[hashed] = append(b.metrics[hashed], ts)
}

func (b *Block) GetTimeSeries(timeSeries MetricFamilyTimeSeries) (*MetricFamilyTimeSeries, error) {
	timeSeriesHash, err := timeSeries.Hash()
	if err != nil {
		panic("unexpected err hashing timeseries")
	}

	hashedTimeSeries := b.metrics[timeSeriesHash]
	mp := MetricPoint{
		Name:     timeSeries.Def.Name,
		LabelSet: timeSeries.LabelSet,
	}

	for _, ts := range hashedTimeSeries {
		tsMp := MetricPoint{
			Name:     ts.Def.Name,
			LabelSet: ts.LabelSet,
		}
		if mp.MetadataEquals(&tsMp) {
			return &ts, nil
		}
	}
	return nil, ErrTimeSeriesNotFound
}
