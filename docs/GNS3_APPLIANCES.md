StoneWork in GNS3
=================

To import StoneWork into GNS3, all you need is to have stonework docker image installed on your host,
which means that it should be listed in your `docker images` output.
In such state, just import its appliance file, i.e. stonework.gns3a in a standard manner, for more information
about this please follow the Deployment section.


VM Compilation
--------------

To create StoneWork full VM image capable of running in GNS3 run:
```
$ make vm-image
```
from stonework workspace root.

This will create 2 files in build/ directory:
- the gns3 appliance file is a small text file specific to GNS3 which describes requirements and other details about an
  image
- qcow2 image which is standard qemu-kvm image and can be used also in other virtualization software

To specify CNFs for addition to VM, use CNFs spec file as:
```
make vm-image CNFS_SPEC=<path/to/spec-file.yaml>
```
The format of the spec file is described in the [script][add-cnfs-script].


Deployment
----------

This section describes only a key steps of StoneWork deployment process in GNS3, since it is an opensource software
and there is plenty of documentation all over Internet.

Basically we need to import the appliance file.
This can be done as described [here][gns3-import-docs].
Appliance file for StoneWork full VM appears in build directory after compilation and
appliance file for StoneWork mini resides in scripts/vm/ directory.

However the tricky part (on windows host) is that StoneWork full image is a qemu-kvm image and thus requires kvm to
work. On Linux there is no problem, however on Windows, we first need to install custom GNS3 VM server on VMware,
to support images requiring kvm.

This is well described [here][gns3-wizard-docs].
NOTE: If you have troubles connecting to created GNS3 VM server, try port 80 (instead of 3080) and from
drop down menu choose Host binding IP address from the same subnet as shown in GNS3 VM:

![GNS3 VM Screenshot][gns3-vm-screenshot]

Just for imagination, in my case it was: 192.168.41.1. and port 80.
After successful installation and connection, the server should occur with green light in Servers Summary window and
you should be able to import appliances into that server.


Usage / Management
------------------

StoneWork can be managed either from CLI (experienced users) or from config editor (easy-going web UI).
To access the CLI, right click on your running StoneWork mini node and select Auxiliary console.
To access config editor web UI, open your web browser and go to http://localhost:<port>, the exact address together
with port is listed in GNS3 in your "Topology Summary" window, "Console" column.


StonerWork full VM Development
------------------------------

While development, the full VM images were executed in virsh for the sake of development speed.
In this regard, it is useful to list few commands.

Install built image into virt, similarly as GNS3 does:
NOTE: this must be executed from StoneWork workspace root
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

List all VMs:
```
$ virsh list --all
```

Start the StoneWork VM:
```
$ virsh start StoneWork --console
```

Shutdown the VM:
```
$ virsh shutdown StoneWork
```

Remove the VM:
```
$ virsh undefine StoneWork
```


Custom StoneWork container
--------------------------

By default StoneWork docker container images are obtained from
ghcr.io/pantheontech/stonework located [here][stonework-gh-docker-registry].
To use custom built image, you only need to build a local image:
```
$ make images
```

[gns3-vm-screenshot]: img/gns3-vm.png
[add-cnfs-script]: scripts/vm/add-cnfs.py
[gns3-import-docs]: https://docs.gns3.com/docs/using-gns3/beginners/import-gns3-appliance/
[gns3-wizard-docs]: https://docs.gns3.com/docs/getting-started/setup-wizard-gns3-vm/
[stonework-gh-docker-registry]: https://github.com/orgs/PANTHEONtech/packages/container/package/stonework
