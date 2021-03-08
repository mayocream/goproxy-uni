package handler

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aofei/air"
	"github.com/goproxy/goproxy"
	"github.com/mayocream/goproxy-uni/base"
)

var (
	// goproxyViper is used to get the configuration items of the Goproxy.
	goproxyViper = base.Viper.Sub("goproxy")

	// hhGoproxy is an instance of the `goproxy.Goproxy`.
	hhGoproxy = goproxy.New()

	// goproxyAutoRedirect indicates whether the automatic redirection
	// feature is enabled for Goproxy.
	goproxyAutoRedirect = goproxyViper.GetBool("auto_redirect")

	// goproxyAutoRedirectMinSize is the minimum size of the Goproxy used to
	// limit at least how big Goproxy cache can be automatically redirected.
	goproxyAutoRedirectMinSize = goproxyViper.GetInt64("auto_redirect_min_size")

	// goproxyReplaceMap ...
	goproxyReplaceMap = goproxyViper.GetStringMapString("replace")

	// goproxyReplacer
	goproxyReplacer = newGoproxyPathReplacer(goproxyReplaceMap)
)

func init() {
	if err := goproxyViper.Unmarshal(hhGoproxy); err != nil {
		base.Logger.Fatal().Err(err).
			Msg("failed to unmarshal goproxy configuration items")
	}

	goproxyLocalCacheRoot, err := os.MkdirTemp(
		goproxyViper.GetString("local_cache_root"),
		"goproxy-local-caches",
	)
	if err != nil {
		base.Logger.Fatal().Err(err).
			Msg("failed to create goproxy local cache root")
	}
	base.Air.AddShutdownJob(func() {
		for i := 0; i < 60; i++ {
			time.Sleep(time.Second)
			err := os.RemoveAll(goproxyLocalCacheRoot)
			if err == nil {
				break
			}
		}
	})

	hhGoproxy.Cacher = diskCacher

	hhGoproxy.ErrorLogger = log.New(base.Logger, "", 0)

	base.Air.BATCH(nil, "/*", hGoproxy)
}

// hGoproxy handles requests to play with Go module proxy.
func hGoproxy(req *air.Request, res *air.Response) error {
 	req.Path = goproxyReplacer.Replace(req.Path)

	hhGoproxy.ServeHTTP(res.HTTPResponseWriter(), req.HTTPRequest())
	return nil
}

func newGoproxyPathReplacer(s map[string]string) *strings.Replacer {
	var oldNew []string
	for k, v := range s {
		oldNew = append(oldNew, k, v)
	}
	return strings.NewReplacer(oldNew...)
}
