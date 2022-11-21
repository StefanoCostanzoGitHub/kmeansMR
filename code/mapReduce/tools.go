package mapReduce

import "kmeansMR/mapReduce/clusters"

type keyValue struct {
	center int
	obs    clusters.Observation
}

/*


func (c *Cluster) Append(point Observation) {
	c.Observations = append(c.Observations, point)
}

*/
