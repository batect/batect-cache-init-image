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
	"strings"
)

type config struct {
	Caches []cache
}

type cache struct {
	Path string
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

	arg := os.Args[1]
	spec := config{}
	decoder := json.NewDecoder(strings.NewReader(arg))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&spec); err != nil {
		exitWithError("Input is invalid", err)
	}
}

func exitWithError(context string, err error) {
	fmt.Printf("%v: %v\n", context, err)
	os.Exit(1)
}
