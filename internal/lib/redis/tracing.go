package redis

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type startKey struct{}
type ctxKey struct{}

// OpenTracingHook opentracing hook
type OpenTracingHook struct {
	cfg    *Config
	status *redis.PoolStats
}

var _ redis.Hook = OpenTracingHook{}

func (o OpenTracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {

	spanCtx, span := otel.Tracer("Redis-V8").Start(ctx, "Redis-V8")
	span.SetAttributes(attribute.Key("db.system").String("redis"))
	span.SetAttributes(attribute.Key("redis.cmd").String(cmd.String()))

	spanCtx = context.WithValue(spanCtx, ctxKey{}, span)
	return context.WithValue(spanCtx, startKey{}, time.Now()), nil
}

func (o OpenTracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start := ctx.Value(startKey{}).(time.Time)
	span := ctx.Value(ctxKey{}).(trace.Span)
	elapsed := time.Since(start)
	if err := cmd.Err(); err != nil {
		recordError(span, err, 0)
	}
	span.End()
	o.report(false, elapsed, cmd)
	return nil
}

func (o OpenTracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	content := ""
	for _, cmd := range cmds {
		content += cmd.String() + "\n"
	}

	spanCtx, span := otel.Tracer("Redis-V8").Start(ctx, o.cfg.Name)
	span.SetAttributes(attribute.Key("db.system").String("redis"))
	span.SetAttributes(attribute.Key("redis.num_cmd").Int(len(cmds)))
	span.SetAttributes(attribute.Key("redis.cmds").String(content))

	spanCtx = context.WithValue(spanCtx, ctxKey{}, span)
	return context.WithValue(spanCtx, startKey{}, time.Now()), nil
}

func (o OpenTracingHook) AfterProcessPipeline(ctx context.Context, cmdList []redis.Cmder) error {
	start := ctx.Value(startKey{}).(time.Time)
	span := ctx.Value(ctxKey{}).(trace.Span)
	elapsed := time.Since(start)
	for i, cmd := range cmdList {
		if err := cmd.Err(); err != nil {
			recordError(span, err, i)
		}
	}
	span.End()
	o.report(true, elapsed, cmdList...)
	return nil
}

func recordError(span trace.Span, err error, index int) {
	if index == 0 {
		span.SetAttributes(attribute.Key("redis.cmd.error").String(fmt.Sprintf("err: %v", err.Error())))
	} else {
		span.SetAttributes(attribute.Key("redis.cmd."+strconv.Itoa(index)+".error").
			String(fmt.Sprintf("err: %v", err.Error())))
	}
}

func appendCmd(b []byte, cmd redis.Cmder) []byte {
	const lenLimit = 64

	for i, arg := range cmd.Args() {
		if i > 0 {
			b = append(b, ' ')
		}

		start := len(b)
		b = appendArg(b, arg)
		if len(b)-start > lenLimit {
			b = append(b[:start+lenLimit], "..."...)
		}
	}

	if err := cmd.Err(); err != nil {
		b = append(b, ": "...)
		b = append(b, err.Error()...)
	}

	return b
}

func appendArg(b []byte, v interface{}) []byte {
	switch v := v.(type) {
	case nil:
		return append(b, "<nil>"...)
	case string:
		return appendUTF8String(b, BytesFormat(v))
	case []byte:
		return appendUTF8String(b, v)
	case int:
		return strconv.AppendInt(b, int64(v), 10)
	case int8:
		return strconv.AppendInt(b, int64(v), 10)
	case int16:
		return strconv.AppendInt(b, int64(v), 10)
	case int32:
		return strconv.AppendInt(b, int64(v), 10)
	case int64:
		return strconv.AppendInt(b, v, 10)
	case uint:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint8:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint16:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint32:
		return strconv.AppendUint(b, uint64(v), 10)
	case uint64:
		return strconv.AppendUint(b, v, 10)
	case float32:
		return strconv.AppendFloat(b, float64(v), 'f', -1, 64)
	case float64:
		return strconv.AppendFloat(b, v, 'f', -1, 64)
	case bool:
		if v {
			return append(b, "true"...)
		}
		return append(b, "false"...)
	case time.Time:
		return v.AppendFormat(b, time.RFC3339Nano)
	default:
		return append(b, fmt.Sprint(v)...)
	}
}

func appendUTF8String(dst []byte, src []byte) []byte {
	if isSimple(src) {
		dst = append(dst, src...)
		return dst
	}

	s := len(dst)
	dst = append(dst, make([]byte, hex.EncodedLen(len(src)))...)
	hex.Encode(dst[s:], src)
	return dst
}

func isSimple(b []byte) bool {
	for _, c := range b {
		if !isSimpleByte(c) {
			return false
		}
	}
	return true
}

func isSimpleByte(c byte) bool {
	return simple[c]
}

var simple = [256]bool{
	'-': true,
	'_': true,

	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,

	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'g': true,
	'h': true,
	'i': true,
	'j': true,
	'k': true,
	'l': true,
	'm': true,
	'n': true,
	'o': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'u': true,
	'v': true,
	'w': true,
	'x': true,
	'y': true,
	'z': true,

	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
	'G': true,
	'H': true,
	'I': true,
	'J': true,
	'K': true,
	'L': true,
	'M': true,
	'N': true,
	'O': true,
	'P': true,
	'Q': true,
	'R': true,
	'S': true,
	'T': true,
	'U': true,
	'V': true,
	'W': true,
	'X': true,
	'Y': true,
	'Z': true,
}

func StringFormat(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Bytes converts string to byte slice.
func BytesFormat(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
