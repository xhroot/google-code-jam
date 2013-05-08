// Run as Node:
//    go run prog.go
// Run as Master:
//    go run prog.go -mode=2 < prog.in > prog.out

package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"runtime"
	"time"
)

const port = ":4000"

var nodeUrls = []string{
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
	"ec2-111-111-111-111.compute-1.amazonaws.com",
}

// This is the basic algorithm to be executed on each case.
func solver(input *Input, c chan<- *Result) {
	paint_used := uint64(0)
	rings_drawn := uint64(0)

	for {
		current_r := uint64(2)*rings_drawn + input.R
		paint_needed := uint64(2)*current_r + uint64(1)
		if paint_needed+paint_used > input.T {
			break
		}
		paint_used += paint_needed
		rings_drawn += uint64(1)
	}

	c <- &Result{input.Id, rings_drawn}
}

// Define the case inputs.
type Input struct {
	Id int
	R  uint64
	T  uint64
}

// Case number and answer.
type Result struct {
	Id     int
	Answer uint64
}

func main() {
	// Use all available processors.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Check mode flag. `Node` (default) will receive a batch of inputs and return
	// the results. `Master` sends batches. `Single` processes all inputs 
	// locally.
	ip := flag.Int("mode", 1, "1=Node, 2=Master, 3=Single")
	flag.Parse()
	if *ip == 1 {
		// Run as `Node`.
		batchSolverNode()
		return
	}

	// Get case count.
	var T int
	fmt.Scan(&T)

	// Read all case inputs.
	inputs := make([]*Input, T)
	for i := 0; i < T; i++ {
		var r, t uint64
		fmt.Scan(&r)
		fmt.Scan(&t)
		inputs[i] = &Input{i, r, t}
	}

	// Create the channel results will be reported on.
	c := make(chan *Result, T)

	// The clock is ticking.
	startTime := time.Now()

	if *ip == 2 {
		// Run as `Master`.
		// Create one input batch per node.
		nodeCount := len(nodeUrls)
		inputBatches := make([][]*Input, nodeCount)
		for i := 0; i < nodeCount; i++ {
			inputBatches[i] = make([]*Input, 0)
		}
		// Poor man's load balancing: round-robin distribution.
		for index, input := range inputs {
			x := index % nodeCount
			inputBatches[x] = append(inputBatches[x], input)
		}
		// Send batches async.
		for index, inputBatch := range inputBatches {
			go sendBatchToNode(inputBatch, nodeUrls[index], c)
		}
	} else {
		// Run as `Single`.
		spawnSolvers(inputs, c)
	}

	// Listen on the channel. Gather answers into answer slice in order.
	answers := make([]uint64, T)
	for i := 0; i < T; i++ {
		res := <-c
		answers[res.Id] = res.Answer
	}

	// All results have returned at this point; end timer.
	done := time.Since(startTime)

	// Print answers.
	for index, answer := range answers {
		fmt.Printf("Case #%v: %v\n", index+1, answer)
	}

	fmt.Printf("t=`%v`\n", done)
}

func sendBatchToNode(inputs []*Input, url string, c chan<- *Result) {
	// Establish TCP connection to node. 
	conn, _ := net.Dial("tcp", url+port)
	// Send batch of inputs.
	gob.NewEncoder(conn).Encode(inputs)
	results := make([]*Result, len(inputs))
	// Wait for results and unpack.
	gob.NewDecoder(conn).Decode(&results)
	// Send results 1 at a time through the channel.
	for _, res := range results {
		c <- res
	}
}

func batchSolverNode() {
	// Wait for connection.
	l, _ := net.Listen("tcp", port)
	conn, _ := l.Accept()
	defer conn.Close()
	// Receive input batch.
	var inputs []*Input
	gob.NewDecoder(conn).Decode(&inputs)
	batchSize := len(inputs)
	// Create a channel for results to be returned to.
	c := make(chan *Result, batchSize)
	// Execute solution on the inputs.
	spawnSolvers(inputs, c)
	// Gather channel results into Result slice. Ignore order for now.
	results := make([]*Result, batchSize)
	for i := 0; i < batchSize; i++ {
		res := <-c
		results[i] = res
	}
	// Send back the result batch.
	gob.NewEncoder(conn).Encode(results)
}

// Spawn solvers; 1 thread per input. Answer will be returned back via channel.
func spawnSolvers(inputs []*Input, c chan<- *Result) {
	for _, input := range inputs {
		go solver(input, c)
	}
}
