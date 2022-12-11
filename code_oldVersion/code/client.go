package main

import (
	"fmt"
	"image/color"
	utils "kmeansMR/cluster"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

var PORT string = "5000"
var NMAPS [1]int = [1]int{5}
var THRSHOLD float64 = 0.01
var MAXITER int = 5

var KS [1]int = [1]int{10}
var PATHS [1]string = [1]string{"./points/rand10000.txt"}

//var PATHS [2]string = [2]string{"./points/rand1000.txt", "./points/rand10000.txt"}

type result struct {
	path string
	k    int
	nmap int
	time time.Duration
}

func main() {
	client, err := rpc.Dial("tcp", ":"+PORT)
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	var out1 int
	var out2 utils.Output
	// default test

	var results []result

	fmt.Println("\ndefault config: THRSHOLD = ", THRSHOLD, ", MAXITER = ", MAXITER, "\n")
	// Loop for path
	for _, p := range PATHS {
		fmt.Println("Path: ", p)
		// loop for k
		for _, k := range KS {
			fmt.Println("-- Number k: ", k)
			// loop for nmap
			for _, nmap := range NMAPS {
				err = client.Call("API.SetKmeans", utils.TestInput{NumMap: nmap, MaxIter: MAXITER, ThrShold: THRSHOLD}, &out1)
				if err != nil || out1 != 0 {
					log.Fatal("[PATHS: ", p, " nmap: ", nmap, " ks: ", k, "] Error in API.SetKmeans: ", err)
				}
				start := time.Now()
				err = client.Call("API.MapReduce", utils.Input{K: k, NumMap: nmap, File: p}, &out2)
				if err != nil {
					log.Fatal("[PATHS: ", p, " nmap: ", nmap, " ks: ", k, "] Error in API.MapReduce: ", err)
				}
				elapsed := time.Since(start)
				results = append(results, result{p, k, nmap, elapsed})
				fmt.Println("---- nmap = ", nmap, " -> time: ", elapsed)
			}
			fmt.Println("")
		}
	}

	graphResult(PATHS[0], results)
}

func graphResult(path string, results []result) {
	testChart := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	testChart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeInfographic}),
		charts.WithTitleOpts(opts.Title{Title: "Test result", Subtitle: "Path = " + path + ". x = num map,  y = time(ms)"}))

	// Put data into instance
	for _, k := range KS {
		items := make([]opts.LineData, 0)
		for _, res := range results {
			if path == res.path && k == res.k {
				items = append(items, opts.LineData{Value: res.time / 1000000})

			}
		}
		testChart.SetXAxis([]string{"1", "2", "3", "4", "5"}).
			AddSeries("Category "+strconv.Itoa(k), items)
	}

	f, _ := os.Create("line.html")
	_ = testChart.Render(f)
}

/*
// Sort by age, keeping original order or equal elements.
sort.SliceStable(family, func(i, j int) bool {
	return family[i].Age < family[j].Age
})
*/

/*
	fmt.Printf("Distributed Map Reduce Client\n")

	// Analize input data
	if len(os.Args) != 4 {
		fmt.Printf("Usage: go run client.go [k] [fileObj] [port server]\n")
		os.Exit(1)
	}
	file := os.Args[2]
	k, e := strconv.Atoi(os.Args[1])
	if e != nil {
		fmt.Printf("k values is not correct!")
		os.Exit(1)
	}

	// Open Rpc connection
	client, err := rpc.Dial("tcp", ":"+os.Args[3])
	defer client.Close()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	var out Output
	// Call function Partition
	err = client.Call("API.MapReduce", Input{k, file}, &out)
	if err != nil {
		log.Fatal("Error in API.MapReduce: ", err)
	}

	// Analize output data
	for i, c := range out.Cc {
		fmt.Printf("Cluster: %d\n", i)
		fmt.Printf("Centered at x: %.2f y: %.2f\n", c.Center[0], c.Center[1])
	}
	plotKmeans(out.Cc, k, out.NPoints)
}
*/
func plotKmeans(cc utils.Clusters, k int, nPoints int) {
	/** Plot centers in Scatter plot **/
	p := plot.New()

	xysCenters := make(plotter.XYs, k)
	for i, c := range cc {
		xysCenters[i].X = c.Center[0]
		xysCenters[i].Y = c.Center[1]
	}

	s, err := plotter.NewScatter(xysCenters)

	if err != nil {
		log.Fatalf("could not create scatter: %v", err)
		return
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)

	wt, err := p.WriterTo(500, 500, "png")
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
