package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"labs/charting"
)

// Consider rewriting storing logic, as some variable can have too maky of
// commbinations, and we can end up with a very large cache. Maybe we can store only some variables,
// or use a more efficient way to generate the key
type ResponseCache struct {
	outcoming map[string]*charting.RenderResponse
	incoming  map[string]*charting.RenderRequest
}

func NewResponseCache() *ResponseCache {
	return &ResponseCache{
		outcoming: make(map[string]*charting.RenderResponse),
		incoming:  make(map[string]*charting.RenderRequest),
	}
}

// generateCacheKey creates a hash of the entire request to account for all variables
func (rc *ResponseCache) generateCacheKey(req *charting.RenderRequest) string {
	jsonData, _ := json.Marshal(req)
	hash := md5.Sum(jsonData)
	return fmt.Sprintf("%x", hash)
}

func (rc *ResponseCache) GetResponse(req *charting.RenderRequest) (*charting.RenderResponse, bool) {
	if req == nil {
		return nil, false
	}
	key := rc.generateCacheKey(req)
	if res, ok := rc.outcoming[key]; ok && !res.IsExpired() {
		return res, true
	}
	return nil, false
}

func (rc *ResponseCache) StoreResponse(req *charting.RenderRequest, res *charting.RenderResponse) {
	if req == nil || res == nil || res.CachePolicy == charting.CachePolicyDontCache {
		return
	}
	key := rc.generateCacheKey(req)
	rc.outcoming[key] = res
	rc.incoming[key] = req
}

func (rc *ResponseCache) Clear() {
	rc.outcoming = make(map[string]*charting.RenderResponse)
	rc.incoming = make(map[string]*charting.RenderRequest)
}
