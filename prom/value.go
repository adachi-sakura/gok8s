package prom

import "github.com/prometheus/common/model"

func GetLatestSampleStream(m model.Matrix) *model.SampleStream{
	if m.Len() == 0 {
		return nil
	}
	time := model.Time(0)
	pos := -1
	for n, sampleStream := range m {
		samplePair := sampleStream.Values[len(sampleStream.Values)-1]
		if samplePair.Timestamp.After(time) {
			time = samplePair.Timestamp
			pos = n
		}
	}
	if pos == -1 {
		panic("invalid index")
	}
	return m[pos]
}