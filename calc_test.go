package main

import (
	"testing"
	"reflect"
)

func TestUpdateAverages(t *testing.T) {

	type ValueTime struct {
		float64
		bool
	}
	k := "t"
	tData := []struct {
		name string
		vals []ValueTime
		ea   float64
	}{
		{
			"one_value",
			[]ValueTime{
				{2.5, true},
			},
			2.5,
		},
		{
			"two_values",
			[]ValueTime{
				{2.5, true},
				{3.5, true},
			},
			3,
		},
		{
			"three_values",
			[]ValueTime{
				{2.5, true},
				{3.5, true},
				{3, true},
			},

			3,
		},
		{
			"after_timespan",
			[]ValueTime{
				{2.5, true},
				{3.5, true},
				{3, false},
			},
			3.25,
		},
	}

	for _, test := range tData {
		a := Averages{}
		q := map[string][]float64{}
		t.Run(test.name, func(t *testing.T) {

			var cVals []float64

			for _, v := range test.vals {
				UpdateAverages(k, v.float64, &q, &a, v.bool)
				if v.bool {
					cVals = append(cVals, v.float64)
				} else {
					cVals = cVals[1:]
					cVals = append(cVals, v.float64)
				}
			}

			if a[k] != test.ea {
				t.Fatalf("unexpected average, expected: %v, got: %v", test.ea, a[k])
			}

			if !reflect.DeepEqual(q[k], cVals) {
				t.Fatalf("unexpected queue, expected: %v, got: %v", cVals, q[k])
			}
		})
	}

}
