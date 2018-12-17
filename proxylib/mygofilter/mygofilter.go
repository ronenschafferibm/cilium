// Copyright 2018 Authors of Cilium
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

package mygofilter

import (
	"bytes"

	"github.com/cilium/cilium/proxylib/proxylib"

	log "github.com/sirupsen/logrus"
)

type factory struct{}

func init() {
	log.Info("init(): Registering mygofilterFactory")
	proxylib.RegisterParserFactory("mygofilter", &factory{})
}

type parser struct {
	connection *proxylib.Connection
	inserted   bool
}

func (f *factory) Create(connection *proxylib.Connection) proxylib.Parser {
	log.Debugf("mygofilterFactory: Create: %v", connection)

	return &parser{connection: connection}
}

func (p *parser) OnData(reply, endStream bool, dataArray [][]byte) (proxylib.OpType, int) {

	// inefficient, but simple for now
	myData := bytes.Join(dataArray, []byte{})
	log.Debugf("*-*-*-OnData-len(myData): %v", len(myData))
	if len(myData) == 0 {
	return proxylib.NOP, 0
	}
	return proxylib.PASS, len(myData)
}
