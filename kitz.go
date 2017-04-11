// Copyright © 2017 Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kitz

import (
	"errors"
	"github.com/cenkalti/backoff"
	"github.com/go-kit/kit/log"
	"io"
	"net"
)

const (
	Endpoint = "listener.logz.io:5050"
	proto    = "tcp"
)

var ErrorConnection = errors.New("No connection defined to logz")
var ErrorInvalidToken = errors.New("Invalid token")

type Logger struct {
	conn   io.Writer
	ep     string
	ts     log.Valuer
	logger log.Logger
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if l.conn == nil {
		return 0, ErrorConnection
	}
	err = backoff.Retry(func() error {
		_, err := l.conn.Write(p)
		return err
	}, backoff.NewExponentialBackOff())
	if err != nil {
		return 0, err
	}
	return len(p), err
}

// WithTimestamp overrides DefaultTimestampUTC
func (l *Logger) WithTimestamp(ts log.Valuer) *Logger {
	l.logger = log.With(l.logger, "time", ts)
	return l
}

// WithEndpoint overrides default endpoint: listener.logz.io:5050 endpoint
func (l *Logger) WithEndpoint(ep string) (*Logger, error) {
	conn, err := net.Dial(proto, ep)
	if err != nil {
		return l, err
	}
	l.conn = conn
	return l, nil
}

// Creates a new kitz logger with defaults listener.logz.io:5050 and DefaultTimestampUTC
func New(token string) (*Logger, error) {
	l := Logger{
		ep: Endpoint,
		ts: log.DefaultTimestampUTC,
	}
	if token == "" {
		return nil, ErrorInvalidToken
	}
	klogger := log.NewJSONLogger(&l)
	klogger = log.With(klogger, "token", token)
	klogger = log.With(klogger, "time", l.ts)
	l.logger = klogger
	conn, err := net.Dial(proto, Endpoint)
	if err != nil {
		return &l, err
	}
	l.conn = conn
	return &l, nil
}

// Build returns the (configured) go-kit logger
func (l *Logger) Build() log.Logger {
	return l.logger
}
