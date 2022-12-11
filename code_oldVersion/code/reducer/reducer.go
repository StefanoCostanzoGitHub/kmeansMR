package main

import (
	"fmt"
	utils "kmeansMR/cluster"
	"log"
	"net"
	"net/rpc"
	"time"
)

type API int

var Cc utils.Clusters

// func Reducer(cc clusters.Clusters, numPoints int, deltaThreshold float64, chInRed chan InRed, chOutRed chan clusters.Clusters, chMastRed chan bool) {
func (a *API) Reducer(input utils.InRed, reply *utils.OutRed) error {
	start := time.Now()
	for _, kv := range input { // Join clusters
		Cc[kv.Center].Append(kv.Obs)
	}
	fmt.Println("esalpsed join (high): ", time.Since(start))

	Cc.Recenter() // Recenter (depending on reducers or not)
	fmt.Println("esalpsed recenter (medium): ", time.Since(start))

	Cc.Reset() // Empty cluster from observations
	fmt.Println("esalpsed reset (low): ", time.Since(start))

	*reply = utils.OutRed(Cc) // Send cluster result
	return nil
}

// Set initial cluster
func (a *API) InitReducer(input utils.Clusters, reply *utils.OutRed) error {
	Cc = input
	return nil
}

func main() {
	// Open Rpc connection
	api := new(API)
	server := rpc.NewServer()
	err := server.RegisterName("API", api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	// Rpc lister
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal("Listener error", err)
	}

	log.Printf("Reducer is listening")
	server.Accept(listener)
}
