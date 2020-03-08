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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Caches []cache
}

type cache struct {
	Path string
	Uid  int
	Gid  int
}

func main() {
	if len(os.Args) <= 1 {
		println("No arguments provided.")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		println("Too many arguments provided.")
		os.Exit(1)
	}

	config := loadConfig()

	for _, cache := range config.Caches {
		if err := process(cache); err != nil {
			exitWithError(fmt.Sprintf("Could not process %v", cache.Path), err)
		}

		fmt.Printf("Processed %v.\n", cache.Path)
	}

	println("Done.")
}

func process(cache cache) error {
	initFilePath := filepath.Join(cache.Path, ".cache-init")
	f, err := os.Create(initFilePath)

	if err != nil {
		return fmt.Errorf("could not create %v: %w", initFilePath, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %v: %w", initFilePath, err)
	}

	if err := os.Chown(initFilePath, cache.Uid, cache.Gid); err != nil {
		return fmt.Errorf("could not set owner and group for %v: %w", initFilePath, err)
	}

	if err := os.Chown(cache.Path, cache.Uid, cache.Gid); err != nil {
		return fmt.Errorf("could not set owner and group for %v: %w", cache.Path, err)
	}

	return nil
}

func loadConfig() config {
	arg := os.Args[1]
	config := config{}
	decoder := json.NewDecoder(strings.NewReader(arg))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&config); err != nil {
		exitWithError("Input is invalid", err)
	}

	return config
}

func exitWithError(context string, err error) {
	fmt.Printf("%v: %v\n", context, err)
	os.Exit(1)
}
