package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/coderollers/go-logger"
	"github.com/coderollers/go-stats/concurrency"
	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"my-microservice/configuration"
)

func StartHttpServer(ctx context.Context, ginRouter *gin.Engine, grpcServer *grpc.Server, grpcWebWrapper *grpcweb.WrappedGrpcServer) {
	concurrency.GlobalWaitGroup.Add(1)
	defer concurrency.GlobalWaitGroup.Done()

	conf := configuration.AppConfig()
	log := logger.SugaredLogger()

	if grpcWebWrapper == nil {
		// No GRPC-Web, Gin-only server
		httpSrv := &http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HttpPort),
			Handler: ginRouter,
		}

		// Start the HTTP Server
		go func() {
			log.Infof("Listening on port %d", conf.HttpPort)
			if err := httpSrv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalf("Unrecoverable HTTP Server failure: %s", err.Error())
				}
			}
		}()

		// Block until SIGTERM/SIGINT
		<-ctx.Done()

		// Clean up and shutdown the HTTP server
		cleanCtx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.CleanupTimeoutSec)*time.Second)
		defer cancel()
		log.Infof("Attempting to shutdown the HTTP server with a timeout of %d seconds", conf.CleanupTimeoutSec)
		if err := httpSrv.Shutdown(cleanCtx); err != nil {
			log.Errorf("HTTP server failed to shutdown gracefully: %s", err.Error())
		} else {
			log.Infof("HTTP Server was shutdown successfully")
		}
	} else {
		// GRPC native + GRPC-Web + Gin on same port

		var mixedHandler, http1Handler http.Handler

		if conf.Development {
			http1Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if grpcWebWrapper.IsGrpcWebRequest(request) {
					log.Debugf("Request handled by GRPC-Web")
					grpcWebWrapper.ServeHTTP(writer, request)
					return
				}
				log.Debugf("Request handled by Gin")
				ginRouter.ServeHTTP(writer, request)
			})
			mixedHandler = newHttpAndGrpcMuxWithDebug(http1Handler, grpcServer)
		} else {
			// Handle GRPC-Web and Gin multiplexing
			http1Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				if grpcWebWrapper.IsGrpcWebRequest(request) {
					grpcWebWrapper.ServeHTTP(writer, request)
					return
				}
				ginRouter.ServeHTTP(writer, request)
			})

			// Add GRPC Native to the multiplexer
			mixedHandler = newHttpAndGrpcMux(http1Handler, grpcServer)
		}
		http2Srv := &http2.Server{}
		http1Srv := &http.Server{
			Handler: h2c.NewHandler(mixedHandler, http2Srv),
		}

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.HttpPort))
		if err != nil {
			log.Fatalf("Cannot open listener socket: %s", err.Error())
		}

		// Start the HTTP Server
		go func() {
			concurrency.GlobalWaitGroup.Add(1)
			defer concurrency.GlobalWaitGroup.Done()

			if err = http1Srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Failed to serve multiplexed endpoint: %s", err.Error())
			}
		}()

		// Block until SIGTERM/SIGINT
		<-ctx.Done()

		// Clean up and shutdown the HTTP server
		cleanCtx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.CleanupTimeoutSec)*time.Second)
		defer cancel()
		log.Infof("Attempting to shutdown the HTTP server with a timeout of %d seconds", conf.CleanupTimeoutSec)
		if err := http1Srv.Shutdown(cleanCtx); err != nil {
			log.Errorf("HTTP server failed to shutdown gracefully: %s", err.Error())
		} else {
			log.Infof("HTTP Server was shutdown successfully")
		}
	}
}

func newHttpAndGrpcMuxWithDebug(httpHandler http.Handler, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log := logger.SugaredLogger()
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("content-type"), "application/grpc") {
			log.Debugf("HTTP/2 GRPC")
			grpcHandler.ServeHTTP(writer, request)
			return
		}
		log.Debugf("Received non-GRPC HTTP/2 request or HTTP/1.1 request")
		httpHandler.ServeHTTP(writer, request)
	})
}

func newHttpAndGrpcMux(httpHandler http.Handler, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.ProtoMajor == 2 && strings.HasPrefix(request.Header.Get("content-type"), "application/grpc") {
			grpcHandler.ServeHTTP(writer, request)
			return
		}
		httpHandler.ServeHTTP(writer, request)
	})
}
