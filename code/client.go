package main

import (
	"kmeansMR/mapReduce"
	"kmeansMR/mapReduce/clusters"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"
	"net"
    "net/http"
    "net/rpc"
	//"gonum.org/v1/plot"
)

func main() {
	k := 10
	rand.Seed(time.Now().UnixNano())

	// set up a random two-dimensional data set (float64 values between 0.0 and 1.0)
	// 1) CREATE OBSERVATIONS
	var d clusters.Observations
	for x := 0; x < 1024*10; x++ {
		d = append(d, clusters.Coordinates{
			rand.Float64(),
			rand.Float64(),
		})
	}
	fmt.Printf("%d data points\n", len(d))

	// Partition the data points into 7 clusters
	// 1) CREATE WORKERS
	

	// 2) CREATE KMEANS STRUCT. CONFIGURABLE
	km := mapReduce.New()
	// 3) INIT CHANNEL

	// 4) CALL SUBROUTINE
	cc, _ := km.Partition(d, k)

	// 5) ANALIZE RESULTS
	var centers [10][2]float64
	for i, c := range cc {
		centers[i][0], centers[i][1] = c.Center[0], c.Center[1]
		fmt.Printf("Cluster: %d\n", i)
		fmt.Printf("Centered at x: %.2f y: %.2f\n", centers[i][0], centers[i][1])

		plotKmeans(centers)

		/*
			for _, p := range c.PointsInDimension(0) {
				fmt.Printf("x: %.2f \n", p)
			}
		*/

		// Plotter
		/*
			xys, err := readData("data.txt")
			if err != nil {
				log.Fatalf("Could not read.txt: %v", err)
			}
		*/
	}
}

func plotKmeans(centers [10][2]float64) {
	p := plot.New()

	xys := make(plotter.XYs, len(centers))
	for i, xy := range centers {
		xys[i].X = xy[0]
		xys[i].Y = xy[1]
	}

	s, err := plotter.NewScatter(xys)
	if err != nil {
		log.Fatalf("could not create scatter: %v", err)
		return
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	wt, err := p.WriterTo(300, 300, "png")
	if err != nil {
		log.Fatalf("Could not create writer: %v", err)
		return
	}

	f, err := os.Create("out.png")
	if err != nil {
		log.Fatalf("could not create out.png: %v", err)
		return
	}
	_, err = wt.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write to out.png: %v", err)
		return
	}

	if err := f.Close(); err != nil {
		log.Fatalf("could not close out.png: %v", err)
		return
	}
}
