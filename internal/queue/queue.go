package queue

import (
	"context"
	"log"
	"time"

	"encoding/json"

	"github.com/chinu-anand/crawlerx/internal/models"

	"github.com/redis/go-redis/v9"
)

type JobQueue interface {
	Enqueue(job models.CrawlJob) error
	Dequeue() <-chan models.JobPayload
}

type RedisQueue struct {
	client *redis.Client
	key    string
	ctx    context.Context
}

func NewRedisQueue(redisURL string) *RedisQueue {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("❌ Failed to parse Redis URL:", err)
	}

	client := redis.NewClient(opt)
	return &RedisQueue{
		client: client,
		key:    "crawl_jobs",
		ctx:    context.Background(),
	}
}

func (q *RedisQueue) Enqueue(job models.CrawlJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.RPush(q.ctx, q.key, data).Err()
}

func (q *RedisQueue) Dequeue() <-chan models.JobPayload {
	out := make(chan models.JobPayload)

	go func() {
		for {
			// Blocking pop
			res, err := q.client.BLPop(q.ctx, 0*time.Second, q.key).Result()

			if err != nil || len(res) < 2 {
				continue
			}

			var job models.JobPayload
			err = json.Unmarshal([]byte(res[1]), &job)
			if err != nil {
				log.Println("❌ Failed to parse job:", err)
				continue
			}

			out <- job
		}
	}()

	return out
}
