// Copyright 2020 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main_test

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type executionResult struct {
	Output   string
	ExitCode int
}

func runImage(volumesToMount []string, input ...string) executionResult {
	imageTag := getImageTag()

	args := []string{"run", "--rm", "-t"}

	for i, volumeToMount := range volumesToMount {
		args = append(args, "-v")
		args = append(args, fmt.Sprintf("%v:/caches/%d", volumeToMount, i+1))
	}

	args = append(args, imageTag)
	args = append(args, input...)
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()

	result := executionResult{
		Output: strings.TrimSpace(string(output)),
	}

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exiterr.ExitCode()
		} else {
			panic(fmt.Sprintf("Could not run application: %v", err))
		}
	}

	return result
}

func getImageTag() string {
	imageTag, haveImageTag := os.LookupEnv("IMAGE_TAG")

	if !haveImageTag {
		panic("IMAGE_TAG environment variable not set.")
	}

	return imageTag
}

func createEmptyVolume() string {
	volumeName := fmt.Sprintf("batect-cache-init-image-tests-%v", randomIdentifier())

	cmd := exec.Command("docker", "volume", "create", volumeName)

	if output, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("Could not create volume, error was: %v, output was: %v", err, output))
	}

	return volumeName
}

func addCacheInitFile(volumeName string) {
	// #nosec G204
	cmd := exec.Command("docker", "run", "--rm", "-t", "-v", fmt.Sprintf("%v:/cache", volumeName), "alpine:3.11.3", "touch", "/cache/.cache-init")

	if output, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("Could not add .cache-init file to volume, error was: %v, output was: %v", err, string(output)))
	}
}

func getVolumeContents(volumeName string) volume {
	// #nosec G204
	cmd := exec.Command("docker", "run", "--rm", "-t", "-v", fmt.Sprintf("%v:/cache", volumeName), "alpine:3.11.3", "ls", "-an", "--color=never", "/cache")
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(fmt.Sprintf("Could not get volume contents, error was: %v, output was: %v", err, string(output)))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\r\n")
	volume := volume{
		Contents: []volumeEntry{},
	}

	splitter := regexp.MustCompile(`\s+`)

	for _, line := range lines[1:] {
		parts := splitter.Split(line, -1)
		uid := parseInt(parts[2])
		gid := parseInt(parts[3])
		name := parts[8]

		if name == "." {
			volume.UID = uid
			volume.GID = gid
		} else if name != ".." {
			volume.Contents = append(volume.Contents, volumeEntry{
				UID:  uid,
				GID:  gid,
				Name: name,
			})
		}
	}

	return volume
}

type volume struct {
	UID      int
	GID      int
	Contents []volumeEntry
}

type volumeEntry struct {
	UID  int
	GID  int
	Name string
}

func parseInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)

	if err != nil {
		panic(fmt.Sprintf("Could not convert %v to an integer, error was: %v", s, err))
	}

	return int(i)
}

func deleteVolume(volumeName string) {
	cmd := exec.Command("docker", "volume", "rm", volumeName)

	if output, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("Could not delete volume, error was: %v, output was: %v", err, output))
	}
}

func randomIdentifier() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
