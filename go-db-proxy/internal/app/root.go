package app

import (
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
	"github.com/core-go/sql"
)

type Root struct {
	Provider   string        `mapstructure:"provider"`
	Server     ServerConfig  `mapstructure:"server"`
	Grpc       ServerConfig  `mapstructure:"grpc"`
	Sql        sql.Config    `mapstructure:"sql"`
	Log        log.Config    `mapstructure:"log"`
	MiddleWare mid.LogConfig `mapstructure:"middleware"`
}
type ServerConfig struct {
	Name string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Port *int64 `mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
}
