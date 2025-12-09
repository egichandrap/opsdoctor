package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

type RequestSpec struct {
	Method  string
	URL     string
	Body    []byte
	Header  map[string]string
	Timeout time.Duration
}

type Result struct {
	URL     string        `json:"url"`
	Method  string        `json:"method"`
	Status  int           `json:"status"`
	Latency time.Duration `json:"latency"`
	Error   string        `json:"error,omitempty"`
	Attempt int           `json:"attempt"`
}

type Summary struct {
	TotalRequests int           `json:"total_requests"`
	Success       int           `json:"success"`
	Failed        int           `json:"failed"`
	MinLatency    time.Duration `json:"min_latency"`
	MaxLatency    time.Duration `json:"max_latency"`
	AvgLatency    time.Duration `json:"avg_latency"`
	Duration      time.Duration `json:"duration"`
	Results       []Result      `json:"results,omitempty"`
}

func doRequest(ctx context.Context, client *http.Client, spec RequestSpec) (Result, error) {
	start := time.Now()
	reqBody := bytes.NewReader(spec.Body)
	req, err := http.NewRequestWithContext(ctx, spec.Method, spec.URL, reqBody)
	if err != nil {
		return Result{URL: spec.URL, Method: spec.Method, Latency: 0}, err
	}
	for k, v := range spec.Header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	lat := time.Since(start)
	if err != nil {
		return Result{URL: spec.URL, Method: spec.Method, Latency: lat, Error: err.Error()}, err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return Result{
		URL:     spec.URL,
		Method:  spec.Method,
		Status:  resp.StatusCode,
		Latency: lat,
	}, nil
}

func TestConcurrent(spec RequestSpec, concurrency, requests int) ([]Result, Summary) {
	startAll := time.Now()
	resultsCh := make(chan Result, requests)
	tasks := make(chan int, requests)

	for i := 0; i < requests; i++ {
		tasks <- i
	}
	close(tasks)

	client := &http.Client{}

	wg := sync.WaitGroup{}
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for attempt := range tasks {
				ctx := context.Background()
				if spec.Timeout > 0 {
					var cancel context.CancelFunc
					ctx, cancel = context.WithTimeout(ctx, spec.Timeout)
					cancel()
				}
				res, _ := doRequest(ctx, client, spec)
				res.Attempt = attempt + 1
				resultsCh <- res
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var results []Result
	var sumLatency time.Duration
	var minLatency time.Duration
	var maxLatency time.Duration
	success := 0
	fail := 0
	first := true

	for r := range resultsCh {
		results = append(results, r)
		if r.Error == "" && r.Status >= 200 && r.Status < 300 {
			success++
		} else {
			fail++
		}
		if r.Latency > 0 {
			sumLatency += r.Latency
			if first || r.Latency < minLatency {
				minLatency = r.Latency
			}
			if first || r.Latency > maxLatency {
				maxLatency = r.Latency
			}
			first = false
		}
	}

	duration := time.Since(startAll)
	avg := time.Duration(0)
	if len(results) > 0 && sumLatency > 0 {
		avg = sumLatency / time.Duration(len(results))
	}

	summary := Summary{
		TotalRequests: requests,
		Success:       success,
		Failed:        fail,
		MinLatency:    minLatency,
		MaxLatency:    maxLatency,
		AvgLatency:    avg,
		Duration:      duration,
		Results:       results,
	}
	return results, summary
}

func TestSingle(spec RequestSpec) Result {
	client := &http.Client{Timeout: spec.Timeout}
	ctx := context.Background()
	if spec.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, spec.Timeout)
		cancel()
	}
	res, _ := doRequest(ctx, client, spec)
	return res
}

func NewSpec(method, url string, body []byte, timeout time.Duration, headers map[string]string) RequestSpec {
	if method == "" {
		method = http.MethodGet
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	return RequestSpec{
		Method:  method,
		URL:     url,
		Body:    body,
		Header:  headers,
		Timeout: timeout,
	}
}

func SummaryToJSON(s Summary) ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
