StoneWork in GNS3
=================

To import StoneWork into GNS3, all you need is to have the StoneWork Docker image installed on your host,
which means that it should be listed in your `docker images` output.

In such state, just import its appliance file, `stonework.gns3a` in a standard manner, for more information
about this please follow the Deployment section.

VM Compilation
--------------

To create a full StoneWork VM image, capable of running in GNS3, run:
```
$ make vm-image
```
from stonework workspace root.

This will create 2 files in the `build/` directory:

- The GNS3 appliance file is a small text file, specific to GNS3 which describes requirements and other details about an image
- The qcow2 image, which is a standard qemu-kvm image and can be used also in other virtualization software

To specify CNFs for addition to VM, use the CNFs spec file as:
```
make vm-image CNFS_SPEC=<path/to/spec-file.yaml>
```
The format of the spec file is described in the [script][add-cnfs-script].


Deployment
----------

This section only describes key steps of StoneWork deployment process in GNS3, since it is an open-source software
and there is plenty of documentation all over the internet.

Basically, we need to import the appliance file. This can be done as described [here][gns3-import-docs].

An appliance file for the full StoneWork VM appears in the build directory after compilation and the appliance file for StoneWork mini resides in the `scripts/vm/` directory.

However, the tricky part (on a  Windows host), is that the full StoneWork image is a qemu-kvm image and thus requires kvm to
work. On Linux, there is no problem, however on Windows, we first need to install a custom GNS3 VM server on VMware, to support images requiring kvm.

This is well described [here][gns3-wizard-docs].

**Note:** If you have troubles connecting to the created GNS3 VM server, try port 80 (instead of 3080) and from the drop down menu, choose **Host binding IP address from the same subnet** as shown in the GNS3 VM:

![GNS3 VM Screenshot][gns3-vm-screenshot]

For better imagination, in my case it was: 192.168.41.1. and port 80.

After a successful installation and connection, the server should appear with a green light in the **Servers Summary** window and
you should be able to import appliances into that server.


Usage / Management
------------------

StoneWork can be managed either from CLI (experienced users) or from config editor (easy-going web UI).

- To access the CLI, right click on your running StoneWork mini node and select *Auxiliary console*.
- To access config editor web UI, open your web browser and go to http://localhost:<port>.

The exact address, together with the port, is listed in GNS3 in your **Topology Summary** window, in the **Console** column.


Full StonerWork VM Development
------------------------------

While in development, the full VM images were executed in *virsh*, for the sake of development speed.

In this regard, it is useful to list few commands.

**Note**: this must be executed from StoneWork workspace root

1. Install the built image into virt, similarly as GNS3 does:
```
$ virt-install \
--name StoneWork \
--memory 4096 \
--vcpus 2 \
--cpu host \
--disk build/stonework.qcow2,bus=sata \
--import \
--os-type linux \
--os-variant ubuntu20.04 \
--network default,model=e1000 \
--graphics none \
--console pty,target_type=serial
```

2. List all VMs:
```
$ virsh list --all
```

3. Start the StoneWork VM:
```
$ virsh start StoneWork --console
```

- Shutdown the VM:
```
$ virsh shutdown StoneWork
```

- Remove the VM:
```
$ virsh undefine StoneWork
```
Custom StoneWork Container
--------------------------

By default, the StoneWork Docker container images are obtained from
`ghcr.io/pantheontech/stonework`, located [here][stonework-gh-docker-registry].

To use a custom-built image, you only need to build a local image:
```
$ make images
```

[gns3-vm-screenshot]: img/gns3-vm.png
[add-cnfs-script]: scripts/vm/add-cnfs.py
[gns3-import-docs]: https://docs.gns3.com/docs/using-gns3/beginners/import-gns3-appliance/
[gns3-wizard-docs]: https://docs.gns3.com/docs/getting-started/setup-wizard-gns3-vm/
[stonework-gh-docker-registry]: https://github.com/orgs/PANTHEONtech/packages/container/package/stonework
