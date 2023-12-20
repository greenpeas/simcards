package app

import (
	"app/internal/config"
	"app/internal/repository/fast"
	"app/internal/repository/spo"
	"app/internal/services/logger"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg      config.Config
	ctx      context.Context
	dbPool   *pgxpool.Pool
	logger   *logger.Logger
	spoRepo  SpoRepo
	fastRepo *fast.RedisRepo
}

func NewApp(cfg config.Config, ctx context.Context, dbPool *pgxpool.Pool, logger *logger.Logger, spoRepo SpoRepo, fastRepo *fast.RedisRepo) *App {

	return &App{
		cfg,
		ctx,
		dbPool,
		logger,
		spoRepo,
		fastRepo,
	}
}

func InitAndRun(configPath string) {

	log.Println("Init application")

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel

	cfg := config.Init(configPath)
	lgr := logger.NewLogger()

	if !cfg.EnableService {
		log.Println("Service is disabled by config")
		return
	}

	log.Println("Open DB connection pool")
	dbpool, err := connectPg(cfg.Database.Psql.Url, 3, ctx) // 3 попытки

	if err != nil {
		log.Println(err)
		return
	}

	spoRepo := spo.NewSpoRepo(dbpool, lgr, ctx)

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.Database.Redis.Addr,
		Password:     cfg.Database.Redis.Password,
		DB:           cfg.Database.Redis.Db,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	fastRepo := fast.NewRedisRepo(redisClient, lgr, ctx)

	a := NewApp(cfg, ctx, dbpool, lgr, spoRepo, fastRepo)
	a.Run()

}
