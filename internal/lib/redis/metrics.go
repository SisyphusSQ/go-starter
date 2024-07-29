package redis

import (
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const namespace = "redis_v8_client"

func (o OpenTracingHook) report(pipe bool, elapsed time.Duration, cmds ...redis.Cmder) {
	//address := o.cfg.Addr
	//name := o.cfg.Name
	errStr := ""
	cmdStr := ""
	//pipeStr := fmt.Sprintf("%t", pipe)

	for _, cmd := range cmds {
		cmdStr += cmd.Name() + ";"

		if err := cmd.Err(); err != nil && !errors.Is(err, redis.Nil) {
			errStr += err.Error() + ";"
		}

		//if cmd.Err() == redis.Nil {
		//	_metricMisses.Inc(name, address) //未命中
		//}
		//
		//if cmd.Err() == nil {
		//	_metricHits.Inc(name, address) //命中缓存
		//}
	}
	cmdStr = strings.TrimSuffix(cmdStr, ";")

	//if len(errStr) > 0 {
	//	_metricReqErr.Inc(name, address, cmdStr, errStr)
	//}
	//
	//_metricReqDur.Observe(int64(elapsed.Seconds()), name, address, cmdStr)
	//
	//_metricConnTotal.Add(float64(o.status.Hits), name, address, "total")
	//
	//_metricConnCurrent.Set(float64(o.status.TotalConns), name, address, "total")
}
