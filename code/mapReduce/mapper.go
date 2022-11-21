package mapReduce

import (
	"kmeansMR/mapReduce/clusters"
	"fmt"
	"net"
    "net/http"
    "net/rpc"
)

func Mapper(points []int, obs clusters.Observations, chInMap chan clusters.Clusters, chInRed chan InRed) {
	// dataset clusters.Observations, points []int, cc clusters.Clusters, ch_mapToMast chan int, ch_mapToRed chan clusters.Clusters
	////////////////////////////// Implementation single reducer without combiner works
	// 0) input: observations, points, clusters
	if len(points) != len(obs) {
		fmt.Printf("Error in dimensions array!\n")
		return
	}
	for {
		cc := <-chInMap
		// If nil the process is finished (DEV)
		if cc == nil {
			fmt.Printf("Mapper finished work correctly\n")
			return
		}
		// 1) Exec mapper work, save changes
		var kvs []keyValue
		changes := 0
		for p, point := range obs {
			ci := cc.Nearest(point)
			kvs = append(kvs, keyValue{ci, point})
			if points[p] != ci {
				points[p] = ci
				changes++
			}
		}

		// 2) Return num changes
		chInRed <- InRed{kvs, changes}
	}
}
