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

package dns

import (
	"bytes"
	"encoding/binary"

	"github.com/cilium/cilium/proxylib/proxylib"

	log "github.com/sirupsen/logrus"
	"github.com/miekg/dns"
)

type factory struct{}

func init() {
	log.Info("init(): Registering dnsParserFactory")
	proxylib.RegisterParserFactory("dns", &factory{})
}

type parser struct {
	connection *proxylib.Connection
	inserted   bool
}

func (f *factory) Create(connection *proxylib.Connection) proxylib.Parser {
	log.Debugf("dnsParserFactory: Create: %v", connection)

	return &parser{connection: connection}
}

func (p *parser) OnData(reply, endStream bool, dataArray [][]byte) (proxylib.OpType, int) {
	data := bytes.Join(dataArray, []byte{})

	if len(data) == 0 {
		return proxylib.NOP, 0
	}

	dnsLen := binary.BigEndian.Uint16(data[0:2])
	dnsMsg := data[2:]

	if len(dnsMsg) < int(dnsLen) {
		bytesToRequest := int(dnsLen) - len(dnsMsg)
		log.Debugf("*-*-*-OnData-bytesToRequest: %v", bytesToRequest)
		return proxylib.MORE, bytesToRequest
	}

	log.Debugf("*-*-*-OnData-dnsMsg: %v", dnsMsg)
	log.Debugf("*-*-*-OnData-data: %v", data)

	m := new(dns.Msg)
	if err := m.Unpack(dnsMsg); err != nil {
		log.Debugf("*-*-*-OnData-err: %v", err)
	}
	log.Debugf("*-*-*-OnData-m: %v", m)
	log.Debugf("*-*-*-OnData-m.Question: %v", m.Question)


	// Prevent the message from sending to the dns server
	if m.Question[0].Name == "google.com." {
		log.Debugf("*-*-*-========= GOOGLE.COM \n")
		return proxylib.DROP, len(data)
	}else{
		log.Debugf("*-*-*-========= *NOT* GOOGLE.COM \n")
	}


	if len(m.Answer) > 0 {
		ans := m.Answer[0]
		switch rrType := ans.(type) {
		default:
			log.Debugf("*-*-*-OnData-m.Answer: %T", rrType)
		case *dns.A:
			x := ans.(*dns.A)
			log.Debugf("*-*-*-OnData-m.Answer: %#v", x.A)
		}
		log.Debugf("*-*-*-OnData-m.Answer: %#v", m.Answer[0])
	}

	log.Debugf("*-*-*-OnData-len(data): %v", len(data))
	log.Debugf("*-*-*-OnData-reply: %v", reply)

	return proxylib.PASS, len(data)
	//return proxylib.MORE, 1
}
