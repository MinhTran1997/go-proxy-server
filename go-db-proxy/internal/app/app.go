package app

import (
	"context"
	. "github.com/core-go/health"
	"github.com/core-go/log/zap"
	s "github.com/core-go/sql"
	c "github.com/core-go/sql/cache"
	g "github.com/core-go/sql/grpc-server"
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/godror/godror"
	"github.com/teris-io/shortid"
)

type ApplicationContext struct {
	HealthHandler *Handler
	Handler       *s.Handler
	GrpcHandler   *g.GRPCHandler
}

func NewApp(ctx context.Context, conf Root) (*ApplicationContext, error) {
	db, er1 := s.Open(conf.Sql)
	if er1 != nil {
		return nil, er1
	}

	sqlHealthChecker := s.NewHealthChecker(db)
	healthHandler := NewHandler(sqlHealthChecker)
	cache, err := c.NewMemoryCacheServiceByConfig(c.CacheConfig{
		Size:             1024,
		CleaningEnable:   true,
		CleaningInterval: 1600,
	})
	grpcCache, err := c.NewMemoryCacheServiceByConfig(c.CacheConfig{
		Size:             1024,
		CleaningEnable:   true,
		CleaningInterval: 1600,
	})
	if err != nil {
		return nil, err
	}
	handler := s.NewHandler(db, s.ToCamelCase, cache, Generate, log.ErrorMsg)
	grpcHandler := g.NewHandler(db, s.ToCamelCase, grpcCache, Generate, log.ErrorMsg)
	app := &ApplicationContext{
		HealthHandler: healthHandler,
		Handler:       handler,
		GrpcHandler:   grpcHandler,
	}
	return app, nil
}

var sid *shortid.Shortid

func Generate(ctx context.Context) (string, error) {
	if sid == nil {
		s2, err := shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			return "", err
		}
		sid = s2
	}
	return sid.Generate()
}
