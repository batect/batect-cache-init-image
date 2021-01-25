# batect-cache-init-image

[![Build Status](https://img.shields.io/github/workflow/status/batect/batect-cache-init-image/Pipeline/main)](https://github.com/batect/batect-cache-init-image/actions?query=workflow%3APipeline+branch%3Amain)
[![License](https://img.shields.io/github/license/batect/batect-cache-init-image.svg)](https://opensource.org/licenses/Apache-2.0)
[![Chat](https://img.shields.io/badge/chat-on%20spectrum-brightgreen.svg)](https://spectrum.chat/batect)

A Docker image that initialises cache volumes for [Batect](https://batect.dev).

## Why is this necessary?

Docker volumes have two major drawbacks when used for caches:

* If they are empty (including if they've just been created), then everything in the target directory of the next container they're mounted into is copied into the volume - this is potentially time consuming and can lead to unexpected behaviour.

* By default, they're mounted into the container with `root` as the owner, so if the container is running as a non-root user, it can't use the directory - which presents problems for containers running with Batect's ['run as current user' mode](https://batect.dev/docs/concepts/run-as-current-user-mode).

Batect uses this image to initialise cache volumes before they're used by a user's containers:

* If the volume is empty, it creates a dummy file (`.cache-init`) in the root of the volume to stop Docker copying the contents of the target container's directory into the volume.

* If the target container is running with 'run as current user' mode, it sets the owner and group of the volume to the desired user so that the volume can be used successfully.
