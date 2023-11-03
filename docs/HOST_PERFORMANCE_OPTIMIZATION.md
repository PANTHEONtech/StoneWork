# Optimizing Host Performance

To enhance the performance of your host system, you can apply the following kernel parameter settings. These settings are especially beneficial for applications that require substantial memory resources.

## Temporary Kernel Parameter Modification

You can modify kernel parameters in runtime using the `sysctl` command with the `-w` prefix. However, please note that these changes will only be effective until the next system restart, as the parameters will revert to their default values.

To temporarily modify kernel parameters, use the following command format:

```bash
sysctl -w parameter_name=value
```

## Permanent Kernel Parameter Configuration

To ensure that your kernel parameters persist across system reboots, you should write them to the /etc/sysctl.conf configuration file. This will make your system use the specified values as the defaults at each boot.

Recommended Kernel Parameters
Here are some recommended kernel parameters to optimize performance:

1. Number of 2MB Hugepages
This parameter controls the number of 2MB hugepages allocated by the system. Applications that benefit from large memory pages can utilize this setting for improved performance.

    ```bash
    vm.nr_hugepages=1024
    ```
2. Maximum Map Count
The vm.max_map_count parameter defines the maximum number of memory map areas (vmas) for a process. It is advised to set this value to be greater than or equal to twice the number of hugepages allocated.

    ```bash
    vm.max_map_count=3096
    ```
3. Access to Hugepages
This setting determines which user groups are allowed to access hugepages. By setting it to 0, you grant access to all groups.

    ```bash
    vm.hugetlb_shm_group=0
    ```
4. Shared Memory Maximum
The kernel.shmmax parameter configures the maximum shared memory size and should be set to at least the total size of the allocated hugepages. For 2MB pages, calculate the TotalHugepageSize as follows:

    TotalHugepageSize = vm.nr_hugepages * 2 * 1024 * 1024
    If the existing kernel.shmmax value (check with cat /proc/sys/kernel/shmmax) is greater than TotalHugepageSize, it's recommended to set kernel.shmmax to the current shmmax value.

    ```bash
    kernel.shmmax=2147483648
    ```

## Applying Configuration Changes
After making changes to the /etc/sysctl.conf file, you can apply the new settings with the following command:

```bash
sysctl -p
```

These adjustments will ensure that your system consistently benefits from the specified kernel parameter settings.