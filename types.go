package main

import "time"

// FlagEvaluateResponse ...
type FlagEvaluateResponse struct {
	RequestID      string `json:"requestId"`
	EntityID       string `json:"entityId"`
	RequestContext struct {
		Fingerprint string `json:"fingerprint"`
	} `json:"requestContext"`
	Match                 bool      `json:"match"`
	FlagKey               string    `json:"flagKey"`
	SegmentKey            string    `json:"segmentKey"`
	Timestamp             time.Time `json:"timestamp"`
	Value                 string    `json:"value"`
	RequestDurationMillis float64   `json:"requestDurationMillis"`
}
