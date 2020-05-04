# Ocilot

Ocilot is a small tool to build applicative container images "the right way"

Instead of spawning a container and snapshotting the filesystem of a running container, 
it simply appends filesystem overlay layers over a base image.

Since most image builds have a huge constant part and a small variable application artefact,
downloading the whole image, spawning a container (with the security issues associated) 
and executing commands is overkill when all you want to do is to add your binary or source files.

## Usage

1. Learn Lua
2. Use ocilot ^^

## Why

Building images quickly and on a kubernetes environment for CI/CD is a pain with Dockerbuild, 
You either have to elevate the pod to what amounts to root on the node, or have to work with compatible but unflexible 
and resource hungry systems like kaniko.

It shouldn’t take more than a few seconds to just add your last python package to a base image, any solution that takes minutes is wrong.

Ocilot provides a small set of functions to help build images YOUR WAY, 
handling layer caching and pure remote image manipulation to keep most CI/CD build scripts
on a very short run time

## But Why, bazel docker_rules already does that !

Yes, I suppose it does, and if you already have an investment in bazel, it’s probably the best solution, but if you don’t the entry cost is prohibitive.

## Can’t you do that with skopeo and umoci and a bunch of python scripts ?

Yep, that was the original plan, it’s … not for the faint of heart since skopeo does the right thing, but umoci needs the full OCI bundle to work, which means downloading all layers.
 