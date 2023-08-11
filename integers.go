package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mailgun/mailgun_mjevan93308/util"
)

type Result struct {
	mu     sync.Mutex
	Even   []int   `json:"even"`
	Odd    []int   `json:"odd"`
	Errors []error `json:"errors,omitempty"`
}

func (result *Result) AddToEven(input int) {
	result.mu.Lock()
	defer result.mu.Unlock()
	result.Even = append(result.Even, input)
}

func (result *Result) AddToOdd(input int) {
	result.mu.Lock()
	defer result.mu.Unlock()
	result.Odd = append(result.Odd, input)
}

func main() {
	result := Result{}
	// allocate a buffered chan to control the number of go routines active in parallel
	// use a waitGroup to ensur all goroutines finish before we return the result
	buff := make(chan struct{}, 10)
	var wg sync.WaitGroup
	for i := 0; i <= 100; i++ {
		// add capacity to the waitGroup for the current iteration
		// for every iteration until the buffered chan reaches capacity, add a val to it
		wg.Add(1)
		buff <- struct{}{}
		// re-allocating looping iter var due to unsafe access within go func
		// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/loopclosure
		i := i

		go func() {
			defer wg.Done()
			if err := worker(&result, i); err != nil {
				fmt.Printf("Encountered error when processing %d: %s", i, err)
			}
			// pop val off buffered chan once job is done
			// allowing next job to kick off
			<-buff
		}()
		wg.Wait()
	}

	result_json, err := json.Marshal(&result)
	if err != nil {
		fmt.Println("Could not marshal result to json")
	}
	fmt.Println(string(result_json))
}

// worker:
// - initializes the api
// - makes the api call
// - process the response
func worker(result *Result, iter int) error {
	api := util.NewApi(util.BuildAddr())
	response, err := api.GetInteger(iter)
	if err != nil {
		fmt.Println("error: %s", err)
		result.Errors = append(result.Errors, err)
		return err
	}

	if util.IsEven(response.Value) {
		result.AddToEven(response.Value)
	} else {
		result.AddToOdd(response.Value)
	}
	return nil
}
