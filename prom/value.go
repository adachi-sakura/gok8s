package prom

import (
	"github.com/prometheus/common/model"
	"math"
)

func getLatestSampleStream(m model.Matrix) *model.SampleStream{
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

type PromRangeValue []model.SamplePair

func GetMatrixValues(val model.Value) []PromRangeValue {
	if val.Type() != model.ValMatrix {
		panic("not matrix value...")
	}
	mat := val.(model.Matrix)
	rangeValues := []PromRangeValue{}
	for _, sampleStream := range mat {
		rangeValues = append(rangeValues, sampleStream.Values)
	}
	return rangeValues
}

func GetVectorValues(val model.Value) []model.SampleValue {
	if val.Type() != model.ValVector {
		panic("not vector value...")
	}
	ret := []model.SampleValue{}
	vec := val.(model.Vector)
	for _, sample := range vec {
		ret = append(ret, sample.Value)
	}
	return ret
}

func Max(values ...model.SampleValue) model.SampleValue {
	ret := 0.
	for _, value := range values {
		ret = math.Max(ret, float64(value))
	}
	return model.SampleValue(ret)
}

func (values PromRangeValue) Increment() float64 {
	return float64(values[len(values)-1].Value - values[0].Value)
}

func SumIncrement(values ...PromRangeValue) float64 {
	sum := 0.
	for _, value := range values{
		sum += value.Increment()
	}
	return sum
}

func (values PromRangeValue) ElapsedTime() float64 {
	duration := values[len(values)-1].Timestamp.Sub(values[0].Timestamp)
	return duration.Seconds()
}

func SumElapsedTime(values ...PromRangeValue) float64 {
	sum := 0.
	for _, value := range values {
		sum += value.ElapsedTime()
	}
	return sum
}