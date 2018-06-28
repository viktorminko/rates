package main

/*
When new rates are received before time frame for moving averages
we just collect them.
If later, then we pop first from queue before adding new rate and
updating average according to formula.
*/
func UpdateAverages(
	i string,
	v float64,
	q *map[string][]float64,
	average *Averages,
	timeDiff bool,
) {

	a := average

	n := len((*q)[i])
	fn := float64(n)

	if _, ok := (*q)[i]; !ok {
		(*q)[i] = []float64{}
	}

	if _, ok := a.Load(i); !ok {
		a.Store(i,v)
	}

	if timeDiff || n == 0 {
		(*q)[i] = append((*q)[i], v)
		t,_ :=  a.Load(i)
		a.Store(i, (t.(float64)*fn + v) / (fn + 1))

		return
	}
	f := (*q)[i][0]
	(*q)[i] = (*q)[i][1:]
	(*q)[i] = append((*q)[i], v)
	t,_ :=  a.Load(i)
	a.Store(i, t.(float64) - f/(fn) + v/(fn))

	return
}
