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
	"testing"
)

func TestNewOVSSwitch(t *testing.T) {
	bridgeName := "tinynet"
	ovsSwitch, err := NewOVSSwitch(bridgeName)
	assert.NoError(t, err)

	exist := ovsSwitch.ovsdb.IsBridgePresent(bridgeName)
	assert.True(t, exist)
}
