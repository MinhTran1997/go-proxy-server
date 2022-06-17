package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/core-go/config"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
	pb "github.com/core-go/sql/grpc"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strconv"

	"go-service/internal/app"
)

const (
	ProviderGRPC = "grpc"
	ProviderRest = "http"
	ProviderBoth = "both"
)

func main() {
	conf := app.Root{}
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}

	r := mux.NewRouter()
	_, err := log.Initialize(conf.Log)
	if err != nil {
		panic(err)
		return
	}
	r.Use(func(handler http.Handler) http.Handler {
		return mid.BuildContextWithMask(handler, MaskLog)
	})
	logger := mid.NewStructuredLogger()
	if log.IsInfoEnable() {
		r.Use(mid.Logger(conf.MiddleWare, log.InfoFields, logger))
	}
	r.Use(mid.Recover(log.ErrorMsg))
	a, er2 := app.Route(r, context.Background(), conf)
	if er2 != nil {
		panic(er2)
	}
	/*
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
		originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	*/
	handler := cors.AllowAll().Handler(r)
	switch conf.Provider {
	case ProviderGRPC:
		if err := StartGrpc(conf, a); err != nil {
			panic(err)
		}
	case ProviderRest:
		fmt.Println(ServerInfo(conf.Server))
		if err := http.ListenAndServe(Addr(conf.Server.Port), handler); err != nil {
			panic(err)
		}
	case ProviderBoth:
		go func() {
			if err := StartGrpc(conf, a); err != nil {
				panic(err)
			}
		}()
		fmt.Println(ServerInfo(conf.Server))
		if err := http.ListenAndServe(Addr(conf.Server.Port), handler); err != nil {
			panic(err)
		}
	default:
		panic(errors.New("invalid provider"))
	}
}

func StartGrpc(conf app.Root, applicationContext *app.ApplicationContext) error {
	lis, err1 := net.Listen("tcp", Addr(conf.Grpc.Port))
	if err1 != nil {
		return err1
	}
	s := grpc.NewServer()
	if err1 != nil {
		return err1
	}
	pb.RegisterDbProxyServer(s, applicationContext.GrpcHandler)
	fmt.Println(ServerInfo(conf.Grpc))
	err1 = s.Serve(lis)
	return err1
}

func Addr(port *int64) string {
	server := ""
	if port != nil && *port >= 0 {
		server = ":" + strconv.FormatInt(*port, 10)
	}
	return server
}
func ServerInfo(conf app.ServerConfig) string {
	if conf.Port != nil && *conf.Port >= 0 {
		return "Start service: " + conf.Name + " at port " + strconv.FormatInt(*conf.Port, 10)
	} else {
		return "Start service: " + conf.Name
	}
}
func MaskLog(name, s string) string {
	return mid.Mask(s, 1, 6, "x")
}
