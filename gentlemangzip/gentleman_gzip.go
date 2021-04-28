package gentlemangzip

import (
	"bytes"
	cgzip "compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.com/proemergotech/errors"
	"gitlab.com/proemergotech/log-go/v3"

	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

const (
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderContentEncoding = "Content-Encoding"
)

// Response return a plugin which will add `Accept-Encoding: gzip` header to the request and decompress the gzip content of the response, if Content-Encoding header is exist and value is gzip.
func Response(logger log.Logger) plugin.Plugin {
	handlers := plugin.Handlers{}

	handlers["before dial"] = func(gCtx *gcontext.Context, h gcontext.Handler) {
		gCtx.Request.Header.Set(HeaderAcceptEncoding, "gzip")
		h.Next(gCtx)
	}

	handlers["after dial"] = func(gCtx *gcontext.Context, h gcontext.Handler) {
		if !strings.Contains(gCtx.Response.Header.Get(HeaderContentEncoding), "gzip") {
			logger.Debug(gCtx, "no need for decompress, there is no content-encoding header with gzip value")
			h.Next(gCtx)
			return
		}

		var buf bytes.Buffer
		gzr, err := cgzip.NewReader(gCtx.Response.Body)
		if err != nil {
			err = errors.Wrap(err, "cannot create response body reader")
			logger.Error(gCtx, "decompression failed", "error", err)
			h.Error(gCtx, err)
			return
		}
		_, err = io.Copy(&buf, gzr) // #nosec
		if err != nil {
			err = errors.Wrap(err, "cannot decompress response body (copy)")
			logger.Error(gCtx, "decompression failed", "error", err)
			h.Error(gCtx, err)
			return
		}

		if err := gzr.Close(); err != nil {
			err = errors.Wrap(err, "cannot decompress response body (close)")
			logger.Error(gCtx, "decompression failed", "error", err)
			h.Error(gCtx, err)
			return
		}

		gCtx.Response.Body = ioutil.NopCloser(&buf)
		h.Next(gCtx)
	}

	return &plugin.Layer{Handlers: handlers}
}

// Request return a plugin which will compress body to gzip content and add the `Content-Encoding: gzip` header to the request.
func Request(logger log.Logger) plugin.Plugin {
	return plugin.NewRequestPlugin(func(gCtx *gcontext.Context, h gcontext.Handler) {

		if gCtx.Request == nil || gCtx.Request.Method == "GET" || gCtx.Request.Body == nil || gCtx.Request.Body == http.NoBody || gCtx.Request.ContentLength <= 0 {
			logger.Debug(gCtx, "calling next because no need to compress")
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
