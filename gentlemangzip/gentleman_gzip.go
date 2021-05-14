package gentlemangzip

import (
	"bytes"
	cgzip "compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/proemergotech/log-go/v3"

	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

const (
	HeaderContentEncoding = "Content-Encoding"
)

// Request return a plugin which will compress body to gzip content and add the `Content-Encoding: gzip` header to the request.
func Request(logger log.Logger) plugin.Plugin {
	return plugin.NewRequestPlugin(func(gCtx *gcontext.Context, h gcontext.Handler) {
		if gCtx.Request == nil || gCtx.Request.Method == "GET" || gCtx.Request.Body == nil || gCtx.Request.Body == http.NoBody || gCtx.Request.ContentLength == 0 {
			h.Next(gCtx)
			return
		}

		var buf bytes.Buffer
		gzw := cgzip.NewWriter(&buf)
		_, err := io.Copy(gzw, gCtx.Request.Body)
		if err != nil {
			logger.Error(gCtx, "cannot compress request body", "error", err)
			h.Next(gCtx)
			return
		}

		if err := gzw.Close(); err != nil {
			logger.Error(gCtx, "cannot compress request body", "error", err)
			h.Next(gCtx)
			return
		}

		gCtx.Request.Body = ioutil.NopCloser(&buf)
		gCtx.Request.Header.Set(HeaderContentEncoding, "gzip")
		h.Next(gCtx)
	})
}
