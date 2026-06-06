package service

import (
	"io"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

func gatewayErrorBodyReadLimitForConfig(cfg *config.Config) int64 {
	limit := gatewayUpstreamErrorBodyReadLimit
	if cfg != nil && cfg.Gateway.LogUpstreamErrorBody && cfg.Gateway.LogUpstreamErrorBodyMaxBytes > int(limit) {
		limit = int64(cfg.Gateway.LogUpstreamErrorBodyMaxBytes)
	}
	return limit
}

func (s *GatewayService) readUpstreamErrorBody(resp *http.Response) []byte {
	if resp == nil || resp.Body == nil {
		return nil
	}
	limit := gatewayErrorBodyReadLimitForConfig(s.cfg)
	body, _ := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if int64(len(body)) > limit {
		return body[:limit]
	}
	return body
}
