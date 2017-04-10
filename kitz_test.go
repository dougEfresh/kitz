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
	"net"
	"strings"
	"testing"
	"time"
)

var sent []byte = make([]byte, 1024)

func TestKitz(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:5050")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go acceptConnection(ln)
	l, err := New("123456789", "localhost:5050", log.DefaultTimestampUTC)
	if err != nil {
		t.Fatal(err)
	}

	if l == nil {
		t.Fatal("Logger is nil")
	}
	l.Log("message", "test msg")
	time.Sleep(200 * time.Millisecond)
	msg := string(sent)
	if msg == "" {
		t.Fatal("Message not send")
	}

	if !strings.Contains(msg, "\"message\":") {
		t.Fatal("no message field")
	}

	if !strings.Contains(msg, "\"token\":\"123456789\"") {
		t.Fatal("no token field")
	}

	if !strings.Contains(msg, "\"time\":\"") {
		t.Fatal("no time field")
	}
}

func acceptConnection(listener net.Listener) {
	for {
		conn, _ := listener.Accept()
		conn.Read(sent)
		conn.Close()
	}
}
