package cron

import (
	"go-starter/internal/lib/log"
	"go-starter/utils/timeutil"
)

func cronSample() {
	log.Logger.Info("sample cronTask")
}

func (s *ServiceImpl) flushOwner() {
	log.Logger.Debug("flushing owner data...")

	t, errWrite := s.start(flushOwner)
	defer s.end(t, errWrite)

	err := s.cache.HSet(s.ctx, string(owner), s.ip, timeutil.CSTLayoutString()).Err()
	if err != nil {
		log.Logger.Errorf("%s set owner to %s error: %v", s.ip, s.ip, err)
	}
}
