package g

import (
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	RID = "rid"
)

var (
	Conf  *viper.Viper
	DB    *gorm.DB
	Pool  *ants.Pool
	Redis *redis.Client
)
