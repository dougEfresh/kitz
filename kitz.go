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

	"github.com/dougEfresh/logzio-go"
	"github.com/go-kit/kit/log"
)

// ErrorInvalidToken Error for tokens
var ErrorInvalidToken = errors.New("Invalid token")

// ClientOptionFunc to set options
type ClientOptionFunc func(*Logger) error

// Logger structure
type Logger struct {
	ts     log.Valuer
	logger log.Logger
	logz   *logzio.LogzioSender
}

// Log Statisfies go-kit log interface
func (l *Logger) Log(keyvals ...interface{}) error {
	return l.logger.Log(keyvals...)
}

// Write Statisfies go-kit log interface
func (l *Logger) Write(p []byte) (n int, err error) {
	err = l.logz.Send(p)
	n = len(p)
	return
}

// SetTimestamp overrides DefaultTimestampUTC
func SetTimestamp(ts log.Valuer) ClientOptionFunc {
	return func(l *Logger) error {
		l.logger = log.With(l.logger, "time", ts)
		return nil
	}
}

// SetUrl overrides default endpoint
func SetUrl(url string) ClientOptionFunc {
	return func(l *Logger) error {
		return logzio.SetUrl(url)(l.logz)
	}
}

// New Creates a new kitz logger with DefaultTimestampUTC
//  l, err := kitz.New("api-key")
func New(token string, options ...ClientOptionFunc) (*Logger, error) {
	l := Logger{
		ts: log.DefaultTimestampUTC,
	}
	if token == "" {
		return nil, ErrorInvalidToken
	}
	klogger := log.NewJSONLogger(&l)
	klogger = log.With(klogger, "token", token)
	klogger = log.With(klogger, "time", l.ts)
	l.logger = klogger

	logz, e := logzio.New(token)
	if e != nil {
		return nil, e
	}
	l.logz = logz
	for _, option := range options {
		if err := option(&l); err != nil {
			return nil, err
		}
	}
	return &l, nil
}

// Stop the logger and drain local values to Logz.io
func (l *Logger) Stop() {
	l.logz.Stop()
}

func (l *Logger) Drain() {
	l.logz.Drain()
}
