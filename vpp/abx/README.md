This subdirectory contains ABX plugin, which can be compiled externally of VPP,
that means, you don't need to compile whole VPP, which takes much longer.
However it is needed to have compatible VPP version installed in your system
(because of headers) and also to clone or download VPP sources of the same
version.

To build the abx, use:
cd vpp2101 && ./build.sh /path/to/vpp/workspace
