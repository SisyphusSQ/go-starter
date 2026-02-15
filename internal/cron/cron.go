package cron

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/robfig/cron/v3"

	"go-starter/config"
	"go-starter/internal/lib/log"
	"go-starter/internal/lib/redis"
	"go-starter/utils"
)

const (
	sampleTaskName      = "sample"
	cacheMetricsTask    = "cache_metrics"
	cacheMetricsCounter = "go_starter:cron:cache_metrics:counter"
)

type Service interface {
	IP() string
}

type ServiceImpl struct {
	ctx    context.Context
	ip     string
	cron   *cron.Cron
	cache  *redis.Client
	locker *redislock.Client
}

func NewCron(config config.Config, cache *redis.Client) (Service, error) {
	if !config.Cron.On {
		return &ServiceImpl{}, nil
	}
	if cache == nil {
		return nil, errors.New("redis client is nil")
	}

	log.Logger.Info("starting cron...")
	timezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}

	cronInstance := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(timezone),
		cron.WithLogger(cron.VerbosePrintfLogger(log.Logger)),
		cron.WithChain(cron.Recover(cron.VerbosePrintfLogger(log.Logger))),
	)

	ip, err := utils.GetIP()
	if err != nil {
		return nil, err
	}

	s := &ServiceImpl{
		ctx:   context.Background(),
		ip:    ip,
		cron:  cronInstance,
		cache: cache,
	}

	s.locker = redislock.New(s.cache)

	if _, err = s.cron.AddFunc("@every 30s", s.sample); err != nil {
		return nil, err
	}
	if _, err = s.cron.AddFunc("@every 1m", s.collectCacheMetrics); err != nil {
		return nil, err
	}
	s.cron.Start()
	return s, nil
}

func (s *ServiceImpl) IP() string {
	return s.ip
}

func (s *ServiceImpl) sample() {
	lock, skip := s.lock(sampleTaskName)
	if skip {
		return
	}
	defer func() {
		if err := lock.Release(s.ctx); err != nil {
			log.Logger.Warnf("release redis lock failed: %v", err)
		}
	}()

	log.Logger.Infof("sample cron task executed on %s", s.ip)
}

// collectCacheMetrics demonstrates a cron task that writes lightweight
// operational data to Redis without touching database dependencies.
func (s *ServiceImpl) collectCacheMetrics() {
	lock, skip := s.lock(cacheMetricsTask)
	if skip {
		return
	}
	defer func() {
		if err := lock.Release(s.ctx); err != nil {
			log.Logger.Warnf("release redis lock failed: %v", err)
		}
	}()

	count, err := s.cache.Incr(s.ctx, cacheMetricsCounter).Result()
	if err != nil {
		log.Logger.Errorf("collect cache metrics failed: %v", err)
		return
	}
	if err = s.cache.Expire(s.ctx, cacheMetricsCounter, 24*time.Hour).Err(); err != nil {
		log.Logger.Errorf("set cache metrics ttl failed: %v", err)
		return
	}

	log.Logger.Infof(
		"cache metrics cron task executed on %s, key=%s, count=%d",
		s.ip,
		cacheMetricsCounter,
		count,
	)
}

func (s *ServiceImpl) lock(taskName string) (*redislock.Lock, bool) {
	lock, err := s.locker.Obtain(s.ctx, taskName, 1*time.Second, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			log.Logger.Infof("task[%s] lock not obtained, skip on %s", taskName, s.ip)
			return nil, true
		}

		log.Logger.Errorf("task[%s] obtain redis lock error on %s: %v", taskName, s.ip, err)
		return nil, true
	}
	return lock, false
}
