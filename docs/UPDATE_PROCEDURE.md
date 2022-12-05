# StoneWork Update Procedure

This document describes the steps needed to update StoneWork for a newer VPP release.

StoneWork is based on *vpp-agent* and thus supports stable VPP versions, supported by vpp-agent.

# 1 Update vpp-agent to support the desired VPP version

## 1.1 Enable build of new VPP version in vpp-base

Check whether `ligato/vpp-base`, with your desired version, already exists.

To do so, just look at `ligato/vpp-base` tags on [DockerHub][dockerhub-tags]. Or directly by docker pull command, for example:
`docker pull ligato/vpp-base:22.10`

If there is no such tagged version, you need to create it.

To do so, inspire yourself by following this [pull request][inspiration-pr].

## 1.2 Add support for new version into vpp-agent

Now, continue the *vpp-agent* update. Complete instructions can be found on this [wiki][agent-instructions].
# 2 Update StoneWork codebase
## 2.1 Update custom VPP plugins

Currently, it is required to update a couple of VPP plugins: 
- abx
- isisx

They can be compiled separately from the rest of StoneWork - take a look at `vpp/abx/README.md`, `vpp/isisx/README.md` for
details. Successful build will produce .so and .api.json files.

You can test it easily, as any externally built VPP plugin. Copy the build
artifacts into standard places, such as `/usr/lib/x86_64-linux-gnu/vpp_plugins/`
and others as seen in `docker/vpp.Dockerfile` (assuming you have the same version
of VPP installed from packages on your host). 

Then verify its presence with `sudo vppctl sh plugins`.

## 2.2 Update StoneWork-specific vpp-agent plugins

StoneWork contains a few additional vpp-agent plugins that are not present in the
*vpp-agent* codebase. These reside in the `plugins/` subdirectory.
Their update follows the same steps and principles as the update of *vpp-agent*
plugins you are already familiar with.

## 2.3 Update StoneWork codebase to use newer VPP and vpp-agent

Now, just update all occurrences of old version to new one.

**Note:** If your *vpp-agent* was not yet merged into upstream repository, but you want to use it
in StoneWork anyway, you will need to set the custom version to be used temporarily.

To do so, use a replace clause in the `go.mod` file.

StoneWork will then use vpp-agent from the folder specified in the replace clause as its build dependency. 

## 2.4 Run StoneWork tests
Run `make test-vpp-plugins`. This builds VPP inside a test image specified in `docker/vpp-test.Dockerfile`, runs
custom VPP plugin tests and shows the test results at the end of the output.

If you make further changes to the custom VPP plugins (for example to fix a failing test case), you can 
then run `make test-vpp-plugins-prebuilt` to skip the lengthy VPP build process.

After that, run `make test`. If tests fail, you have most probably done some mistake in 2.2
and not all containers started successfully, especially those containing VPP and vpp-agent.

In that case, execute the command
`cd examples/testing/100-mock-cnf-module && STONEWORK_IMAGE=stonework:<version> make start-example`
and check whether all containers are running, or some of them exited prematurely.

In latter case, take a look at Docker logs of a particular container and ensure
its *vpp-agent* and all of its plugins support the newly added version, i.e. it
found a compatible API.

## 2.5 Upload new tag for the StoneWork to image repository

StoneWork docker images are present on [GitHub Container Registry][ghcr].

To update images, create and push a git tag into image repository according to
the following convention:
`v<VPP-major>.<VPP-minor>.<patch><optional-identifier>` (for example
`v22.10.0`), where `<patch>` may increase if VPP is updated by its patch version
or if some change is submitted into the control-plane.

This triggers build of the images in repository and tags StoneWork production
image as `ghcr.io/pantheontech/stonework` with the following version tags:
1. Full git tag as-it-is, with trimmed leading '`v`', for example
   `ghcr.io/pantheontech/stonework:22.10.0`. This tag is fixed
   and should never be changed.
2. `<VPP-major>.<VPP-minor>`, for example
   `ghcr.io/pantheontech/stonework:22.10`. This tag points to the
   latest version with the same major and minor version number.
3. `latest`, for example `ghcr.io/pantheontech/stonework:latest`.

This convention allows to keep track of all the patches of the StoneWork
while keeping the `latest` and `<VPP-major>.<VPP-minor>` tags still updated.

After that, all three tagged images are automatically pushed into GitHub
Container Registry.

[dockerhub-tags]: https://hub.docker.com/r/ligato/vpp-base/tags?page=1&ordering=last_updated&name=22.10
[inspiration-pr]: https://github.com/ligato/vpp-base/pull/18
[agent-instructions]: https://github.com/ligato/vpp-agent/wiki/Guide-for-adding-new-VPP-version
[ghcr]: https://github.com/orgs/PANTHEONtech/packages/container/package/stonework
