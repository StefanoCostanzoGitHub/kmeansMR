// Package mapReduce implements the k-means clustering algorithm
package mapReduce

import (
	"kmeansMR/mapReduce/clusters"
	"fmt"
)

// Kmeans configuration/option struct
type Kmeans struct {
	int numMappers
	int numReducers
	string[] mappers
	string[] reducers
	// when a plotter is set, Plot gets called after each iteration
	plotter Plotter
	// deltaThreshold (in percent between 0.0 and 0.1) aborts processing if
	// less than n% of data points shifted clusters in the last iteration
	deltaThreshold float64
	// iterationThreshold aborts processing when the specified amount of
	// algorithm iterations was reached
	iterationThreshold int
}

// The Plotter interface lets you implement your own plotters
type Plotter interface {
	Plot(cc clusters.Clusters, iteration int) error
}

// NewWithOptions returns a Kmeans configuration struct with custom settings
func NewWithOptions(numMappers int, numReducers int, deltaThreshold float64, plotter Plotter) (Kmeans, error) {
	if deltaThreshold <= 0.0 || deltaThreshold >= 1.0 {
		return Kmeans{}, fmt.Errorf("threshold is out of bounds (must be >0.0 and <1.0, in percent)")
	}
	kmeans = Kmeans{
		numMappers: 		numMappers,
		numReducers:		numReducers,
		plotter:            plotter,
		deltaThreshold:     deltaThreshold,
		iterationThreshold: 96,
	}


}

/*
type InRed struct {
	cc      clusters.Clusters
	changes int
	//arrSums [][]int
}

type OutRed struct {
	cc      clusters.Clusters
	changes int
}

type InMap struct {
	chunk  clusters.Observations
	points []int
	cc     clusters.Clusters
}
*/

type InRed struct {
	kvs     []keyValue
	changes int
	//arrSums [][]int
}

// New returns a Kmeans configuration struct with default settings
func New() Kmeans {
	m, _ := NewWithOptions(0.01, nil)
	return m
}

// Partition executes the k-means algorithm on the given dataset and
// partitions it into k clusters
func (m Kmeans) Partition(dataset clusters.Observations, k int) (clusters.Clusters, error) {
	//////////////////////////////////////////////// MASTER START (1) ///////////////////////////////////////

	// 	0) VERIFICA VALORI INGRESSO
	if k > len(dataset) {
		return clusters.Clusters{}, fmt.Errorf("the size of the data set must at least equal k")
	}

	// 1) INITIALIZE CLUSTERS, associate punti a clusters random e var changes
	cc, err := clusters.New(k, dataset)
	if err != nil {
		return cc, err
	}

	points := make([]int, len(dataset))
	/*
		chInpRed := make(chan InRed)
		chOutRed := make(chan OutRed)
		chMap := make(chan InMap)
	*/
	// 1.0) Create channels
	chInMap := make(chan clusters.Clusters)
	chInRed := make(chan InRed)
	chOutRed := make(chan clusters.Clusters)
	chMastRed := make(chan bool)

	// 1.1) start threads
	numMapper := 64
	pos, lenChunk, resChunk := 0, len(dataset)/numMapper, len(dataset)%numMapper
	pos_fin := lenChunk + resChunk
	for i := 0; i < numMapper; i++ {
		go Mapper(points[pos:pos_fin], dataset[pos:pos_fin], chInMap, chInRed)		
		pos = pos_fin
		pos_fin += lenChunk
	}
	go Reducer(cc, len(dataset), m.deltaThreshold, chInRed, chOutRed, chMastRed)


	for i := 0; i < m.iterationThreshold; i++ {
		// 1.2) Reset clusters points
		cc.Reset()

		// 2) send values threads
		chMastRed <- false
		for i := 0; i < numMapper; i++ {
			chInMap <- cc
		}
		// 3) wait cluster results (Channel2)
		cc = <-chOutRed
		if end := <-chMastRed; end {
			chInMap <- nil
			break
		}

		// 4) Plotter (DEV) servira?
		if m.plotter != nil {
			err := m.plotter.Plot(cc, i) // (DEV) attenzione la i è cambiata
			if err != nil {
				return nil, fmt.Errorf("failed to plot chart: %s", err)
			}
		}

	}

	// Return clusters with Channel1
	return cc, nil
	/////////////////////////////////////////////// MASTER STOP (2) ////////////////////////////////
}
