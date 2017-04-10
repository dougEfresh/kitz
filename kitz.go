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

import "io"
import (
	"github.com/go-kit/kit/log"
	"net"
)

const (
	Endpoint = "listener.logz.io:5050"
	proto    = "tcp"
)

type logger struct {
	conn io.Writer
}

func (l logger) Write(p []byte) (n int, err error) {
	return l.conn.Write(p)
}
// WithDefaults creates a new go-kit logger
// This uses DefaultTimestampUTC as a time format and listener.logz.io:5050 endpoint
func WithDefaults(token string) (log.Logger, error) {
	return New(token, Endpoint, log.DefaultTimestampUTC)
}

// Creates a new go-kit logger
func New(token, ep string, ts log.Valuer) (log.Logger, error) {
	conn, err := net.Dial(proto, ep)
	if err != nil {
		return nil, err
	}
	klogger := log.NewJSONLogger(logger{conn})
	klogger = log.With(klogger, "token", token)
	klogger = log.With(klogger, "time", ts)
	return klogger, nil
}
