# StoneWork update procedure

This document describes steps needed to update StoneWork for newer VPP release.

First to mention is that StoneWork is based on vpp-agent and thus supports
stable VPP versions supported by vpp-agent.

# 1 Update vpp-agent to support desired VPP version

## 1.1 Enable build of new VPP version in vpp-base

Check whether ligato/vpp-base with your desired version already exists.
To do so, just look at ligato/vpp-base tags on [dockerhub][dockerhub-tags].
Or directly by docker pull command, for example
`docker pull ligato/vpp-base:21.06`.
If there is no such tagged version, you need to create it.
To do so, inspire yourself by following [pull request][inspiration-pr].

## 1.2 Add support for new version into vpp-agent

Now continue to vpp-agent update. Complete instructions can be found on this
[wiki][agent-instructions].


# 2 Update StoneWork codebase

## 2.1 Update custom VPP plugins

Currently it is needed to update couple of VPP plugins - abx,isisx. It can be compiled
separately from the rest of StoneWork, take a look at vpp/abx/README.md, vpp/isisx/README.md for
details. Successful build will produce .so and .api.json files.
You can test it easily, as any externally built VPP plugin: just copy build
artifacts into standard places as /usr/lib/x86_64-linux-gnu/vpp_plugins/
and others as seen in docker/vpp.Dockerfile (assuming you have the same version
of VPP installed from packages on your host). Then verify its presence by
`sudo vppctl sh plugins`.

## 2.2 Update StoneWork-specific vpp-agent plugins

StoneWork contains a few additional vpp-agent plugins that are not present in
open source vpp-agent codebase. These live in plugins/ subdirectory.
Their update follows the same steps and principles as update of vpp-agent
plugins you are already familiar with.

## 2.3 Update StoneWork codebase to use newer VPP and vpp-agent

Now just update all occurrences of old version to new one. Easy step.
Just note that if your
vpp-agent was not yet merged into upstream repository, but you want to use it
in StoneWork anyway, you will need to set the custom version to be used
temporarily.
To do so, examine the go.mod and use the replace clause prepared on the bottom.
Now verify that `ls submodule/vpp-agent` is not empty,
if it is, download StoneWork submodules by:
`git submodule update --init --recursive`
After these steps, StoneWork will use vpp-agent from submodule/vpp-agent as its
build dependency. Good for development purposes. Now you just need to apply your
patches (or track remote branch) to this vpp-agent workspace.

## 2.4 Run StoneWork tests

Run `make test`. If tests fail, most probably you have done some mistake in 2.2
and most likely not all containers started successfully, especially those
containing vpp and vpp-agent.
In that case do
`cd examples/testing/100-mock-cnf-module && make start-example`
and check whether all containers are running or some of them exited prematurely.
In latter case, take a look at docker logs of particular container and ensure
its vpp-agent and all of its plugins support the newly added version, i.e. it
found a compatible API.

## 2.5 Upload new tagged version of the StoneWork to image repository

Create tag in repository to trigger update of the image in repository.

[dockerhub-tags]: https://hub.docker.com/r/ligato/vpp-base/tags?page=1&ordering=last_updated&name=21.06
[inspiration-pr]: https://github.com/ligato/vpp-base/pull/18
[agent-instructions]: https://github.com/ligato/vpp-agent/wiki/Guide-for-adding-new-VPP-version

