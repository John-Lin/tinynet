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

	log "github.com/sirupsen/logrus"

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

	readCloser, err = cli.ImagePull(ctx, "docker.io/library/"+imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	} else {
		// because readCloser need to be handle so that image can be download.
		// so we send this output to /dev/null
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
	log.Info("%s\n", resp.ID)
	log.Info("%s\n", cInfo.NetworkSettings.SandboxKey)
	log.Info("%s\n", cInfo.Name)
	return cInfo.Name, cInfo.NetworkSettings.SandboxKey
}
