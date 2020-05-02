package main

import (
	optly "github.com/optimizely/go-sdk"
	"github.com/optimizely/go-sdk/pkg/client"
	"github.com/optimizely/go-sdk/pkg/entities"
	"github.com/valyala/fasthttp"
)

var (
	reverseProxyV1, reverseProxyV2 *ReverseProxy
	optlyClient *client.OptimizelyClient
)

func init() {
	reverseProxyV1 = newReverseProxy(config.TargetSite1)
	reverseProxyV2 = newReverseProxy(config.TargetSite2)
	optlyClient, _ = optly.Client(config.OptimizelySDKKey)
}

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, ctx *fasthttp.RequestCtx) {
	if target == config.TargetSite2 {
		reverseProxyV2.serveHTTP(ctx)
		return
	}

	reverseProxyV1.serveHTTP(ctx)
}

// GetFlags ...
func GetFlags(flagKey string, entityID string) bool {
	if !config.FlipEnabled {
		log.Infof("Flipt is not enabled")
		return false
	}
	enabled, err := optlyClient.IsFeatureEnabled(flagKey, entities.UserContext{ID: entityID})
	log.Infof("%s isEnable %v", flagKey, enabled)
	if err != nil {
		log.Infof("Optimizely err", err)
	}

	return enabled
}

// HandleRequestAndRedirect given a request send it to the appropriate url
func HandleRequestAndRedirect(ctx *fasthttp.RequestCtx) {
	clientGAId := ctx.Request.Header.Cookie("_ga")
	// default will be the old xe web site
	urlProxy := config.TargetSite1

	if clientGAId == nil {
		serveReverseProxy(urlProxy, ctx)
		return
	}

	entityID := string(clientGAId)
	log.Infof("entityID: %s", entityID)
	enable := GetFlags(config.FlagKey, entityID)

	if enable {
		urlProxy = config.TargetSite2
	}
	serveReverseProxy(urlProxy, ctx)
}

func healthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("OK"))
}
