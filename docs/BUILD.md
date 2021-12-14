Building StoneWork
==================

StoneWork is packaged inside a Docker container, as an image: `ghcr.io/pantheontech/stonework:<VPP-version>`.

In order to build the image, **only Docker is required** to be installed. Every other dependency is either downloaded
or compiled inside the image build process (including Golang, for example).

## Building The Image

Image build is split into two stages, first a development image is built containing all the sources of StoneWork
as well as tools for development and debugging. From this a minimalistic production-ready image is extracted,
containing only those binaries and libraries that are really needed to run StoneWork.

- To build development image, execute Makefile target:

```
$ make dev-image
```

- To select non-default VPP version, specify also `VPP_VERSION` variable as follows:

```
$ VPP_VERSION=x.y make dev-image
```

- Once the development image is built, you can proceed with the build of the production-ready image:

```
$ make prod-image
```

- If a non-default VPP version was used, `VPP_VERSION` variable needs to be defined at this stage as well:

```
$ VPP_VERSION=x.y make prod-image
```
