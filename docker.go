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
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func initDocker(imageName string) (string, string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// docker pull busybox
	readCloser, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	} else {
		// because readCloser need to be handle so that image can be download.
		// we don't need output so send this to /dev/null
		io.Copy(ioutil.Discard, readCloser)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"sleep", "3600"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	var cInfo types.ContainerJSON
	cInfo, err = cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	return resp.ID, cInfo.NetworkSettings.SandboxKey
}
