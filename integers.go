package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mailgun/mailgun_mjevan93308/util"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

const (
	MaxRoutines = 10
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
	api := util.Api

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelFn()

	// using `sync/errgroup` due to better error handling and easier goroutine limit enforcement
	// thus eliminating need for waitgroup and buffered channel
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(MaxRoutines)
	for i := 0; i < 100; i++ {
		// re-allocating looping iter var due to unsafe access within go func
		// https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/loopclosure
		i := i

		eg.Go(func() error {
			// check for context err, don't proceed if found
			if ctx.Err() != nil {
				fmt.Println("encountered ctx err")
				return ctx.Err()
			}
			if err := worker(ctx, &result, api, i); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("encountered error when processing %d: %s", i, err))
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		fmt.Printf("errgroup returned err: %s\n", err)
		return
	}

	resultJson, err := json.Marshal(&result)
	if err != nil {
		fmt.Println("Could not marshal result to json")
	} else {
		fmt.Println(string(resultJson))
	}
}

// worker initiates the api call and process the response
// - takes a context, pointer to result, pointer to shared api client, and the current iteration var
// - returns an err if encountered
func worker(ctx context.Context, result *Result, api *util.API, iter int) error {
	response, err := api.GetInteger(ctx, iter)
	if err != nil {
		return err
	}

	if util.IsEven(response.Value) {
		result.AddToEven(response.Value)
	} else {
		result.AddToOdd(response.Value)
	}
	return nil
}
