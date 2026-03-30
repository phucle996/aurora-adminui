package victoria

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

type Point struct {
	Timestamp time.Time
	Value     float64
}

type Series struct {
	Labels map[string]string
	Points []Point
}

func NewClient(baseURL string) *Client {
	trimmed := strings.TrimSpace(baseURL)
	trimmed = strings.TrimRight(trimmed, "/")
	return &Client{
		baseURL: trimmed,
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Configured() bool {
	return c != nil && c.baseURL != ""
}

func (c *Client) Query(ctx context.Context, query string, at time.Time) (float64, error) {
	if !c.Configured() {
		return 0, fmt.Errorf("victoria query base url is not configured")
	}

	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, fmt.Errorf("parse victoria query base url: %w", err)
	}
	endpoint.Path = path.Join(endpoint.Path, "/api/v1/query")

	params := endpoint.Query()
	params.Set("query", query)
	params.Set("time", strconv.FormatInt(at.Unix(), 10))
	endpoint.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return 0, fmt.Errorf("victoria query returned status %d", resp.StatusCode)
	}

	var payload instantQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("decode victoria query response: %w", err)
	}
	if payload.Status != "success" {
		if payload.Error != "" {
			return 0, fmt.Errorf("victoria query failed: %s", payload.Error)
		}
		return 0, fmt.Errorf("victoria query returned status %q", payload.Status)
	}

	total := 0.0
	for _, series := range payload.Data.Result {
		if len(series.Value) != 2 {
			continue
		}
		value, err := parseJSONStringFloat(string(series.Value[1]))
		if err != nil {
			continue
		}
		total += value
	}

	return total, nil
}

func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]Point, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("victoria query base url is not configured")
	}

	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse victoria query base url: %w", err)
	}
	endpoint.Path = path.Join(endpoint.Path, "/api/v1/query_range")

	params := endpoint.Query()
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.Unix(), 10))
	params.Set("end", strconv.FormatInt(end.Unix(), 10))
	params.Set("step", strconv.FormatFloat(step.Seconds(), 'f', -1, 64))
	endpoint.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("victoria query_range returned status %d", resp.StatusCode)
	}

	var payload queryRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode victoria query_range response: %w", err)
	}
	if payload.Status != "success" {
		if payload.Error != "" {
			return nil, fmt.Errorf("victoria query_range failed: %s", payload.Error)
		}
		return nil, fmt.Errorf("victoria query_range returned status %q", payload.Status)
	}

	aggregated := map[int64]float64{}
	for _, series := range payload.Data.Result {
		for _, sample := range series.Values {
			if len(sample) != 2 {
				continue
			}
			ts, err := parseUnixSeconds(string(sample[0]))
			if err != nil {
				continue
			}
			value, err := parseJSONStringFloat(string(sample[1]))
			if err != nil {
				continue
			}
			aggregated[ts] += value
		}
	}

	points := make([]Point, 0, len(aggregated))
	for ts, value := range aggregated {
		points = append(points, Point{
			Timestamp: time.Unix(ts, 0).UTC(),
			Value:     value,
		})
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Timestamp.Before(points[j].Timestamp)
	})
	return points, nil
}

func (c *Client) QueryRangeSeries(ctx context.Context, query string, start, end time.Time, step time.Duration) ([]Series, error) {
	if !c.Configured() {
		return nil, fmt.Errorf("victoria query base url is not configured")
	}

	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse victoria query base url: %w", err)
	}
	endpoint.Path = path.Join(endpoint.Path, "/api/v1/query_range")

	params := endpoint.Query()
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.Unix(), 10))
	params.Set("end", strconv.FormatInt(end.Unix(), 10))
	params.Set("step", strconv.FormatFloat(step.Seconds(), 'f', -1, 64))
	endpoint.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("victoria query_range returned status %d", resp.StatusCode)
	}

	var payload queryRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode victoria query_range response: %w", err)
	}
	if payload.Status != "success" {
		if payload.Error != "" {
			return nil, fmt.Errorf("victoria query_range failed: %s", payload.Error)
		}
		return nil, fmt.Errorf("victoria query_range returned status %q", payload.Status)
	}

	out := make([]Series, 0, len(payload.Data.Result))
	for _, matrixSeries := range payload.Data.Result {
		points := make([]Point, 0, len(matrixSeries.Values))
		for _, sample := range matrixSeries.Values {
			if len(sample) != 2 {
				continue
			}
			ts, err := parseUnixSeconds(string(sample[0]))
			if err != nil {
				continue
			}
			value, err := parseJSONStringFloat(string(sample[1]))
			if err != nil {
				continue
			}
			points = append(points, Point{
				Timestamp: time.Unix(ts, 0).UTC(),
				Value:     value,
			})
		}
		sort.Slice(points, func(i, j int) bool {
			return points[i].Timestamp.Before(points[j].Timestamp)
		})
		out = append(out, Series{
			Labels: matrixSeries.Metric,
			Points: points,
		})
	}
	return out, nil
}

type queryRangeResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Data   struct {
		Result []queryRangeMatrixTS `json:"result"`
	} `json:"data"`
}

type instantQueryResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Data   struct {
		Result []instantQueryVectorTS `json:"result"`
	} `json:"data"`
}

type queryRangeMatrixTS struct {
	Metric map[string]string   `json:"metric"`
	Values [][]json.RawMessage `json:"values"`
}

type instantQueryVectorTS struct {
	Metric map[string]string `json:"metric"`
	Value  []json.RawMessage `json:"value"`
}

func parseUnixSeconds(raw string) (int64, error) {
	raw = strings.Trim(raw, "\"")
	if strings.Contains(raw, ".") {
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return 0, err
		}
		return int64(parsed), nil
	}
	return strconv.ParseInt(raw, 10, 64)
}

func parseJSONStringFloat(raw string) (float64, error) {
	return strconv.ParseFloat(strings.Trim(raw, "\""), 64)
}
