package mapReduce

import (
	"kmeansMR/mapReduce/clusters"
	"fmt"
	"net"
    "net/http"
    "net/rpc"
)

func Reducer(cc clusters.Clusters, numPoints int, deltaThreshold float64, chInRed chan InRed, chOutRed chan clusters.Clusters, chMastRed chan bool) {
	numop := 0
	for {
		// AUDIT
		changes := 0
		n := 0

		// 1) Input: clusters, changes
		if flagEnd := <-chMastRed; flagEnd == true {
			fmt.Printf("Reducer finished correctly the execution\n")
			return
		}

		for n < numPoints {
			in := <-chInRed
			changes += in.changes
			// 2) Join clusters
			for _, kv := range in.kvs {
				cc[kv.center].Append(kv.obs)
				n++
			}
			// 3) Update array sums

			// Wait all clusters

			// 3.1) Resolve empty clusters

		}
		numop += 1
		fmt.Printf("numop %d\n", numop)
		// 4) Recenter (depending on reducers or not)
		cc.Recenter()
		chOutRed <- cc
		if changes > 0 || changes < int(float64(numPoints)*deltaThreshold) {
			chMastRed <- false
		} else {
			// Send finish process message to master
			fmt.Printf("Reducer finished correctly the execution\n")
			chMastRed <- true
			return
		}
	}
}

/**
for ci := 0; ci < len(cc); ci++ {
	if len(cc[ci].Observations) == 0 {
		// During the iterations, if any of the cluster centers has no
		// data points associated with it, assign a random data point
		// to it.
		// Also see: http://user.ceng.metu.edu.tr/~tcan/ceng465_f1314/Schedule/KMeansEmpty.html
		var ri int
		for {
			// find a cluster with at least two data points, otherwise
			// we're just emptying one cluster to fill another
			ri = rand.Intn(len(dataset)) //nolint:gosec // rand.Intn is good enough for this
			if len(cc[points[ri]].Observations) > 1 {
				break
			}
		}
		cc[ci].Append(dataset[ri])
		points[ri] = ci

		// Ensure that we always see at least one more iteration after
		// randomly assigning a data point to a cluster
		changes = len(dataset)
	}
}
**/
