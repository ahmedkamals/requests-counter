package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"requests-counter/config"
	"requests-counter/utils/logger"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

const (
	LETTER_BYTES    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LETTER_IDX_BITS = 6                      // 6 bits to represent a letter index.
	LETTER_IDX_MASK = 1<<LETTER_IDX_BITS - 1 // All 1-bits, as many as LETTER_IDX_BITS.
	LETTER_IDX_MAX  = 63 / LETTER_IDX_BITS   // # of letter indices fitting in 63 bits.
)

type (
	Locator struct {
	}
)

func NewLocator() *Locator {
	return &Locator{}
}

func (l *Locator) LoadConfig(data []byte) (*config.Config, error) {
	parsedConfig := config.Config{}

	err := json.Unmarshal(data, &parsedConfig)

	if nil != err {
		return nil, err
	}

	return &parsedConfig, nil
}

func (*Locator) Logger() *logger.Logger {
	return logger.NewLogger()
}

func (*Locator) BlockIndefinitely() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-sigc
}

func (l *Locator) Stats(duration time.Duration) {
	logger := l.Logger()

	for {
		logger.Info("Number of Go routines:", runtime.NumGoroutine())
		time.Sleep(duration * time.Second)
	}
}

func (l *Locator) catchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		l.Logger().Log(goRoutine, functionName, fmt.Sprintf("PANIC Defered [%v] : Stack Trace : %v", r, string(buf)))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func (l *Locator) GetAsTimestamp(nanoseconds int64) time.Time {
	return time.Unix(0, nanoseconds)
}

func (l *Locator) RandString(prefix string, n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), LETTER_IDX_MAX; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), LETTER_IDX_MAX
		}

		if idx := int(cache & LETTER_IDX_MASK); idx < len(LETTER_BYTES) {
			b[i] = LETTER_BYTES[idx]
			i--
		}
		cache >>= LETTER_IDX_BITS
		remain--
	}

	return prefix + string(b)
}

func (l *Locator) FloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', 6, 64)
}
