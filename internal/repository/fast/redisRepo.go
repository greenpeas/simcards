package fast

import (
	"app/internal/dto"
	"app/internal/services/logger"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
)

const ICCIDs_LIST = "ICCIDs:list"

type RedisRepo struct {
	rdb    *redis.Client
	logger *logger.Logger
	ctx    context.Context
}

func NewRedisRepo(rdb *redis.Client, logger *logger.Logger, ctx context.Context) *RedisRepo {
	return &RedisRepo{
		rdb,
		logger,
		ctx,
	}
}

// Получение сим-карты из кеша
func (r *RedisRepo) Get(iccid string) *dto.SimCardMTS {

	simJsonString, prevErr := r.rdb.HGet(r.ctx, ICCIDs_LIST, iccid).Result()

	if prevErr != nil && prevErr != redis.Nil {
		r.logger.Error.Println(prevErr)
	}

	var simCard *dto.SimCardMTS
	json.Unmarshal([]byte(simJsonString), &simCard)

	return simCard
}

// Добавление сим-карты в кеш
func (r *RedisRepo) Add(sim *dto.SimCardMTS) {

	simJson, err := json.Marshal(sim)

	if err != nil {
		r.logger.Error.Println(err)
	}

	addSimErr := r.rdb.HSet(r.ctx, ICCIDs_LIST, sim.Sim.Iccid, simJson).Err()
	if addSimErr != nil {
		r.logger.Error.Println(addSimErr)
	}
}

// Удаление сим-карты из кэша
func (r *RedisRepo) Del(iccid string) {
	r.rdb.HDel(r.ctx, ICCIDs_LIST, iccid)
}
