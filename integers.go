package main

import (
	"encoding/json"
	"fmt"

	"github.com/mailgun/mailgun_mjevan93308/util"
)

type Result struct {
	Even   []int   `json:"even"`
	Odd    []int   `json:"odd"`
	Errors []error `json:"errors,omitempty"`
}

func main() {
	result := Result{}

	// create buffered bool chan of size 10
	s := make(chan bool, 10)
	for i := 0; i <= 100; i++ {
		// for every iteration until the chan reaches capacity, add a val to it
		s <- true

		// re-allocating looping iter var due to unsafe access within go func
		// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/loopclosure
		iter := i
		go func() {
			if err := worker(&result, iter); err != nil {
				fmt.Printf("Encountered error when processing %d: %s", iter, err)
			}

			// pop bool val off buffered chan once job is done
			// allowing next job to kick off
			<-s
		}()
	}
	result_json, err := json.Marshal(result)
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
		result.Errors = append(result.Errors, err)
		return err
	}

	if util.IsEven(response.Value) {
		result.Even = append(result.Even, response.Value)
	} else {
		result.Odd = append(result.Odd, response.Value)
	}
	return nil
}
