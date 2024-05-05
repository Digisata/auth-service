// Package gateway is described reusable package for create gateway server
package gateway

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/digisata/auth-service/api"
	"github.com/digisata/auth-service/bootstrap"
	"github.com/digisata/auth-service/docs"
	"github.com/digisata/auth-service/pkg/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	MaxHeaderBytes = 1 << 20
	ReadTimeOut    = 10 * time.Second
	WriteTimeOut   = 10 * time.Second
)

type Gateway struct {
	*runtime.ServeMux
	Addr           string
	MaxHeaderBytes int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func NewGateway(addr string, opts ...runtime.ServeMuxOption) *Gateway {
	opts = append(opts,
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	gwMux := runtime.NewServeMux(opts...)

	return &Gateway{
		ServeMux:       gwMux,
		Addr:           addr,
		MaxHeaderBytes: MaxHeaderBytes,
		ReadTimeout:    ReadTimeOut,
		WriteTimeout:   WriteTimeOut,
	}
}

func (gw *Gateway) swaggerUIHandler() (http.Handler, error) {
	err := mime.AddExtensionType(".svg", "image/svg+xml")
	if err != nil {
		return nil, errors.Wrap(err, "mime.AddExtensionType")
	}
	subFS, err := fs.Sub(docs.SwaggerUI, "swagger-ui")
	if err != nil {
		return nil, errors.Wrap(err, "fs.Sub")
	}
	return http.FileServer(http.FS(subFS)), nil
}

func (gw *Gateway) Run(ctx context.Context, cfg *bootstrap.Config) error {
	sw, err := gw.swaggerUIHandler()
	if err != nil {
		return errors.Wrap(err, "gw.swaggerUIHandler")
	}

	fileServer := http.FileServer(http.FS(api.FS))
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", sw)

	gwServer := &http.Server{
		Addr: fmt.Sprintf(":%v", gw.Addr),
		Handler: middleware.CORS(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if strings.HasPrefix(request.URL.Path, "/api/v1") {
				gw.ServeMux.ServeHTTP(writer, request)
				return
			}
			mux.ServeHTTP(writer, request)
		}), cfg.Port),
		ReadTimeout:    gw.ReadTimeout,
		WriteTimeout:   gw.WriteTimeout,
		MaxHeaderBytes: gw.MaxHeaderBytes,
	}

	go func() {
		<-ctx.Done()
		gwServer.Shutdown(ctx)
	}()

	err = gwServer.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "gwServer.ListenAndServe")
	}

	return nil
}
