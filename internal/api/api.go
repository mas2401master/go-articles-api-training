package api

import (
	"fmt"

	_ "github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/internal/api/router"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/config"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/db"
)

func Init(configPath string) {
	if configPath == "" {
		configPath = "data/config.yml"
	}
	config.Setup(configPath)
	db.SetupDB()
	conf := config.GetConfig()
	webapi := router.Setup()
	fmt.Println("Go API REST Running on port " + conf.Server.Port)
	_ = webapi.Run(":" + conf.Server.Port)
}
