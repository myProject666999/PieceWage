package main

import (
	"fmt"
	"os"

	"piece-wage/internal/config"
	"piece-wage/internal/handler"
	"piece-wage/internal/middleware"
	"piece-wage/pkg/db"
	"piece-wage/pkg/logger"
	"piece-wage/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		fmt.Printf("load config failed: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("init logger failed: %v\n", err)
		os.Exit(1)
	}
	defer logger.Log.Sync()

	if err := db.Init(&cfg.MySQL); err != nil {
		logger.Log.Fatal("init mysql failed", err)
	}
	defer db.Close()
	logger.Log.Info("mysql connected")

	if err := redis.Init(&cfg.Redis); err != nil {
		logger.Log.Warn("init redis failed, cache will be unavailable", err)
	} else {
		defer redis.Close()
		logger.Log.Info("redis connected")
	}

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handler.NewAuthHandler()
	processHandler := handler.NewProcessHandler()
	reportHandler := handler.NewReportHandler()
	wageHandler := handler.NewWageHandler()
	teamHandler := handler.NewTeamHandler()
	userHandler := handler.NewUserHandler()

	r.POST("/api/auth/login", authHandler.Login)

	authGroup := r.Group("/api")
	authGroup.Use(middleware.JWTAuth())
	{
		authGroup.GET("/auth/profile", authHandler.GetProfile)

		authGroup.POST("/products", processHandler.CreateProduct)
		authGroup.GET("/products", processHandler.ListProducts)

		authGroup.POST("/process-steps", processHandler.CreateStep)
		authGroup.GET("/process-steps", processHandler.ListSteps)
		authGroup.GET("/products/:productId/steps", processHandler.ListStepsByProduct)

		authGroup.POST("/process-prices", processHandler.CreatePrice)
		authGroup.GET("/process-prices/effective", processHandler.GetEffectivePrice)
		authGroup.GET("/process-prices/process/:processId", processHandler.ListPricesByProcess)
		authGroup.GET("/process-prices", processHandler.ListAllPrices)

		authGroup.POST("/reports", reportHandler.CreateReport)
		authGroup.GET("/reports/:id", reportHandler.GetReport)
		authGroup.GET("/reports", reportHandler.ListReports)
		authGroup.PUT("/reports/:id/void", reportHandler.VoidReport)

		authGroup.GET("/wage/summary/:workerId/:month", wageHandler.GetMonthlySummary)
		authGroup.GET("/wage/summaries", wageHandler.ListSummaries)
		authGroup.GET("/wage/details", wageHandler.GetWorkerDetails)
		authGroup.GET("/wage/daily/:workerId/:date", wageHandler.GetWorkerDailyDetails)
		authGroup.GET("/wage/realtime/:workerId/:month", wageHandler.GetRealtimeAccumulate)
		authGroup.POST("/wage/settle/:month", wageHandler.SettleMonth)

		authGroup.POST("/teams", teamHandler.CreateTeam)
		authGroup.GET("/teams", teamHandler.ListTeams)
		authGroup.GET("/teams/:teamId/members", teamHandler.GetTeamMembers)
		authGroup.GET("/allocations/:reportId", teamHandler.GetAllocation)

		authGroup.GET("/users", userHandler.ListUsers)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Info(fmt.Sprintf("server starting on %s", addr))
	if err := r.Run(addr); err != nil {
		logger.Log.Fatal("server start failed", err)
	}
}
