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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the cache-init image", func() {
	var noVolumes []string

	Context("when the application is given an empty configuration", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage(noVolumes, "{}")
		})

		It("returns a zero exit code and no error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "Done.",
				ExitCode: 0,
			}))
		})
	})

	Context("when the application is given a configuration with a single volume", func() {
		var volumeName string
		var volumes []string

		BeforeEach(func() {
			volumeName = createEmptyVolume()
			volumes = []string{volumeName}
		})

		AfterEach(func() {
			deleteVolume(volumeName)
		})

		singleVolumeConfig := `{
								"caches": [
									{ "path": "/caches/1" }
								]
							}`

		singleVolumeWithUIDAndGIDConfig := `{
												"caches": [
													{ "path": "/caches/1", "uid": 123, "gid": 456 }
												]
											}`

		itReturnsAZeroExitCodeAndNoError := func(result *executionResult) {
			It("returns a zero exit code and no error", func() {
				Expect(*result).To(Equal(executionResult{
					Output:   "Processed /caches/1.\r\nDone.",
					ExitCode: 0,
				}))
			})
		}

		itEnsuresVolumeContainsCacheInitFile := func(volumeContents *volume, uid int, gid int) {
			It("ensures that the volume contains a .cache-init file with the provided user and group", func() {
				Expect(volumeContents.Contents).To(ContainElement(volumeEntry{
					UID:  uid,
					GID:  gid,
					Name: ".cache-init",
				}))
			})
		}

		itSetsOwnerAndGroupOfVolume := func(volumeContents *volume, uid int, gid int) {
			It("sets the owner of the volume to the provided user", func() {
				Expect(volumeContents.UID).To(Equal(uid))
			})

			It("sets the group of the volume to the provided group", func() {
				Expect(volumeContents.GID).To(Equal(gid))
			})
		}

		Context("when the volume is empty", func() {
			Context("when no user or group is provided", func() {
				var result executionResult
				var volumeContents volume

				BeforeEach(func() {
					result = runImage(volumes, singleVolumeConfig)
					volumeContents = getVolumeContents(volumeName)
				})

				itReturnsAZeroExitCodeAndNoError(&result)
				itEnsuresVolumeContainsCacheInitFile(&volumeContents, 0, 0)
			})

			Context("when a user and group are provided", func() {
				var result executionResult
				var volumeContents volume

				BeforeEach(func() {
					result = runImage(volumes, singleVolumeWithUIDAndGIDConfig)
					volumeContents = getVolumeContents(volumeName)
				})

				itReturnsAZeroExitCodeAndNoError(&result)
				itSetsOwnerAndGroupOfVolume(&volumeContents, 123, 456)
				itEnsuresVolumeContainsCacheInitFile(&volumeContents, 123, 456)
			})
		})

		Context("when the volume already contains a .cache-init file", func() {
			BeforeEach(func() {
				addCacheInitFile(volumeName)
			})

			Context("when no user or group is provided", func() {
				var result executionResult
				var volumeContents volume

				BeforeEach(func() {
					result = runImage(volumes, singleVolumeConfig)
					volumeContents = getVolumeContents(volumeName)
				})

				itReturnsAZeroExitCodeAndNoError(&result)
				itEnsuresVolumeContainsCacheInitFile(&volumeContents, 0, 0)
			})

			Context("when a user and group are provided", func() {
				var result executionResult
				var volumeContents volume

				BeforeEach(func() {
					result = runImage(volumes, singleVolumeWithUIDAndGIDConfig)
					volumeContents = getVolumeContents(volumeName)
				})

				itReturnsAZeroExitCodeAndNoError(&result)
				itSetsOwnerAndGroupOfVolume(&volumeContents, 123, 456)
				itEnsuresVolumeContainsCacheInitFile(&volumeContents, 123, 456)
			})
		})
	})

	Context("when the application is given a configuration with multiple volumes", func() {
		var volume1Name string
		var volume2Name string
		var result executionResult
		var volume1Contents volume
		var volume2Contents volume

		BeforeEach(func() {
			volume1Name = createEmptyVolume()
			volume2Name = createEmptyVolume()
			volumes := []string{volume1Name, volume2Name}

			config := `{
				"caches": [
					{ "path": "/caches/1" },
					{ "path": "/caches/2" }
				]
			}`

			result = runImage(volumes, config)
			volume1Contents = getVolumeContents(volume1Name)
			volume2Contents = getVolumeContents(volume2Name)
		})

		AfterEach(func() {
			deleteVolume(volume1Name)
			deleteVolume(volume2Name)
		})

		It("returns a zero exit code and no error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "Processed /caches/1.\r\nProcessed /caches/2.\r\nDone.",
				ExitCode: 0,
			}))
		})

		It("ensures that both volumes contain a .cache-init file", func() {
			Expect(volume1Contents.Contents).To(ContainElement(volumeEntry{
				Name: ".cache-init",
			}))

			Expect(volume2Contents.Contents).To(ContainElement(volumeEntry{
				Name: ".cache-init",
			}))
		})
	})

	Context("when the application is given no arguments", func() {
		var result executionResult

		BeforeEach(func() {
			result = runImage(noVolumes)
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
			result = runImage(noVolumes, "{}", "{}")
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
			result = runImage(noVolumes, "{")
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
			result = runImage(noVolumes, `{"thing":2}`)
		})

		It("returns a non-zero exit code and an error", func() {
			Expect(result).To(Equal(executionResult{
				Output:   "Input is invalid: json: unknown field \"thing\"",
				ExitCode: 1,
			}))
		})
	})
})
