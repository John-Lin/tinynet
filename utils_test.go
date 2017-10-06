// Copyright (c) 2017 Che Wei, Lin
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

package tinynet

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestGetAllIPsfromCIDR(t *testing.T) {

	expected := []string{"140.113.234.121", "140.113.234.122", "140.113.234.123", "140.113.234.124", "140.113.234.125", "140.113.234.126"}
	addr, _ := getAllIPsfromCIDR("140.113.234.123/29")

	for i, _ := range addr {
		assert.Equal(t, addr[i], expected[i], "Those two address should be same")
	}
}

func TestInc(t *testing.T) {
	ip, _, _ := net.ParseCIDR("140.113.234.123/30")
	expected := []string{"140.113.234.123", "140.113.234.124", "140.113.234.125"}

	for i, _ := range expected {
		assert.Equal(t, ip.String(), expected[i], "Those two address should be same")
		inc(ip)
	}
}
