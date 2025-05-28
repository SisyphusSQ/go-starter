package cron

import (
	"context"
	"errors"
	"time"

	"github.com/SisyphusSQ/golib/models/do/base_do"
	"github.com/bsm/redislock"
	"github.com/robfig/cron/v3"

	"go-starter/config"
	"go-starter/internal/lib/log"
	"go-starter/internal/lib/redis"
	"go-starter/internal/repository/mysql/my_common"
	"go-starter/internal/repository/mysql/my_example"
	"go-starter/utils"
)

type task string

type key string

const (
	flushOwner task = "flush_owner"
)

const (
	owner key = "sample"
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error("Recovered from panic in cron: %v", r)
		}
	}()
}

type Recover struct{}

type Service interface {
	IP() string
}

type ServiceImpl struct {
	ctx      context.Context
	ip       string
	ckConfig *config.Database

	cron *cron.Cron

	cache  *redis.Client
	locker *redislock.Client

	clusterRepo my_example.AuditClusterRepository
	taskRepo    my_common.TaskResultRepository
}

func NewCron(config config.Config,
	cache *redis.Client,
	clusterRepos my_example.AuditClusterRepository,
	taskRepo my_common.TaskResultRepository) (Service, error) {
	if !config.Cron.On {
		return &ServiceImpl{}, nil
	}

	log.Logger.Info("starting cron...")
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	Cron := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(timezone),
		cron.WithLogger(cron.VerbosePrintfLogger(log.Logger)),
		cron.WithChain(cron.Recover(cron.VerbosePrintfLogger(log.Logger))),
	)

	ip, err := utils.GetIP()
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), "gorm:silent", log.SilentLogger{})
	s := &ServiceImpl{
		ctx:         ctx,
		ip:          ip,
		cron:        Cron,
		cache:       cache,
		ckConfig:    &config.Clickhouse,
		clusterRepo: clusterRepos,
		taskRepo:    taskRepo,
	}

	s.locker = redislock.New(s.cache)

	s.cron.AddFunc("@every 30s", s.flushOwner)
	s.cron.Start()
	return s, nil
}

func (s *ServiceImpl) IP() string {
	return s.ip
}

func (s *ServiceImpl) writeTaskResult(res base_do.TaskResult) error {
	var err error
	if res.TaskStatus == base_do.Processing {
		err = s.taskRepo.CreateTaskResult(s.ctx, res)
		if err != nil {
			log.Logger.Errorf("create task result failed, got err: %v", err)
		}
		return err
	}

	err = s.taskRepo.UpdateByUUID(s.ctx, res)
	if err != nil {
		log.Logger.Errorf("update task result failed, got err: %v", err)
	}
	return nil
}

func (s *ServiceImpl) start(name task) (base_do.TaskResult, error) {
	t := base_do.TaskResult{
		UUID:       utils.UUID(),
		TaskName:   string(name),
		ExecIP:     s.ip,
		TaskStatus: base_do.Processing,
	}
	errWrite := s.writeTaskResult(t)
	t.Start = time.Now()
	return t, errWrite
}

func (s *ServiceImpl) end(t base_do.TaskResult, err error) {
	if err == nil {
		if t.ErrMsg != "" {
			t.TaskStatus = base_do.Error
		} else {
			t.TaskStatus = base_do.Finish
		}

		t.TaskCost = time.Since(t.Start).Milliseconds()
		_ = s.writeTaskResult(t)
	}
}

func (s *ServiceImpl) lock(taskName task) (*redislock.Lock, bool) {
	lock, err := s.locker.Obtain(s.ctx, string(taskName), 1*time.Second, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			log.Logger.Infof("Task[%s] %s is not obtained, skip...", taskName, s.ip)
			return nil, true
		}

		log.Logger.Errorf("Task[%s] %s obtain redis lock error: %v", taskName, s.ip, err)
		return nil, true
	}
	return lock, false
}
