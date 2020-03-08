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
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the cache-init image", func() {
	Context("when the application is given an empty specification", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage("{}")
		})

		It("returns a zero exit code and no output", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "",
				ExitCode: 0,
			}))
		})
	})

	Context("when the application is given no arguments", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage()
		})

		It("returns a non-zero exit code and an error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "No arguments provided.",
				ExitCode: 1,
			}))
		})
	})

	Context("when the application is given multiple arguments", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage("{}", "{}")
		})

		It("returns a non-zero exit code and an error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "Too many arguments provided.",
				ExitCode: 1,
			}))
		})
	})

	Context("when the application is given malformed JSON", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage("{")
		})

		It("returns a non-zero exit code and an error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "Input is invalid: unexpected EOF",
				ExitCode: 1,
			}))
		})
	})

	Context("when the application is given JSON with an unknown field", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage(`{"thing":2}`)
		})

		It("returns a non-zero exit code and an error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "Input is invalid: json: unknown field \"thing\"",
				ExitCode: 1,
			}))
		})
	})
})

type executionResult struct {
	Output   string
	ExitCode int
}

func runImage(input ...string) executionResult {
	imageTag, haveImageTag := os.LookupEnv("IMAGE_TAG")

	if !haveImageTag {
		panic("IMAGE_TAG environment variable not set.")
	}

	args := append([]string{"run", "--rm", "-t", imageTag}, input...)
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
