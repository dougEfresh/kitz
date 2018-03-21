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
	"github.com/go-kit/kit/log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var sent []byte = make([]byte, 1024)

func TestKitz(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	defer ts.Close()
	logger, _ := New("123456789", SetUrl(ts.URL), SetTimestamp(log.DefaultTimestampUTC))
	if logger == nil {
		t.Fatal("Logger is nil")
	}
	logger.Log("message", "test msg")
	logger.Stop()
	time.Sleep(500 * time.Millisecond)

	msg := string(sent)
	if msg == "" {
		t.Fatal("Message not send")
	}
	if !strings.Contains(msg, "\"message\":") {
		t.Fatalf("no message field %s", msg)
	}
	if !strings.Contains(msg, "\"token\":\"123456789\"") {
		t.Fatal("no token field")
	}
	if !strings.Contains(msg, "\"time\":\"") {
		t.Fatal("no time field")
	}
}

func BenchmarkNew(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	logger, _ := New("123456789", SetUrl(ts.URL), SetTimestamp(log.DefaultTimestampUTC))
	if logger == nil {
		b.Fatal("Logger is nil")
	}

	for i := 0; i < b.N; i++ {
		logger.Log("message", "test msg")
	}
}
