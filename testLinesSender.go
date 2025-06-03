package main

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

type TestLinesReq struct {
	lines []string
}

func (p TestLinesReq) Serialize() ([]byte, error) {
	return []byte(strings.Join(p.lines, "\r\n")), nil
}

func NewTestLinesSender(opts LogSenderOpts) ISender {
	var l *GenericSender
	hdrs := opts.Headers
	opts.Headers = map[string]string{}
	for k, v := range hdrs {
		opts.Headers[k] = v
	}

	l = &GenericSender{
		LogSenderOpts: opts,
		mtx:           sync.Mutex{},
		rnd:           rand.New(rand.NewSource(time.Now().UnixNano())),
		timeout:       time.Second * 15,
		path:          "/test-lines",
	}
	l.generate = func() IRequest {
		resLines := make([]string, 0, opts.LinesPS)

		for i := 0; i < opts.LinesPS; i++ {
			resLines = append(resLines, l.pickRandom(opts.Lines))
		}

		return TestLinesReq{lines: resLines}
	}
	return l
}
