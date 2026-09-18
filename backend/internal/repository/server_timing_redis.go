package repository

import (
	"context"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
	"github.com/redis/go-redis/v9"
)

type serverTimingRedisHook struct {
	metric string
}

func (serverTimingRedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h serverTimingRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if !servertiming.Active(ctx) {
			return next(ctx, cmd)
		}
		startedAt := time.Now()
		err := next(ctx, cmd)
		h.record(ctx, startedAt, time.Now(), 1)
		return err
	}
}

func (h serverTimingRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if !servertiming.Active(ctx) {
			return next(ctx, cmds)
		}
		startedAt := time.Now()
		err := next(ctx, cmds)
		h.record(ctx, startedAt, time.Now(), len(cmds))
		return err
	}
}

func (h serverTimingRedisHook) record(ctx context.Context, startedAt, endedAt time.Time, count int) {
	servertiming.Record(ctx, servertiming.MetricRedis, startedAt, endedAt, count)
	if h.metric != "" && h.metric != servertiming.MetricRedis {
		servertiming.Record(ctx, h.metric, startedAt, endedAt, count)
	}
}
