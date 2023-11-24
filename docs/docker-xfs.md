# XFS filesystem for Docker

1. Install the necessary packages:

```bash
sudo apt-get update
sudo apt-get install xfsprogs
```

2. Identify the disk or partition you want to use for the new filesystem. You can use the `lsblk` or `fdisk -l` command to list the available disks and partitions. Make sure you select the correct one as this process will erase all data on it.

3. If the partition you want to use is not formatted, format it with an XFS filesystem. Replace `/dev/sdX` with the appropriate device identifier:

```bash
sudo mkfs.xfs /dev/sdX
```

4. Once the partition is formatted, create a mount point. This will be the directory where the new filesystem will be mounted:

```bash
sudo mkdir /mnt/docker
```

5. Update the `/etc/fstab` file to automatically mount the new filesystem at boot:

```
/dev/sdX  /mnt/docker  xfs  defaults,nofail,discard,noatime,quota,prjquota,pquota,gquota  0 2
```

6. Mount the new filesystem and verify that it is working:

```bash
sudo mount -a
df -h
```

The output of `df -h` should show the new filesystem mounted at `/mnt/docker`.

7. Stop the docker daemon:

```bash
systemctl stop docker
```

8. Update the `/etc/docker/daemon.json` file to point docker to the new mount point:

```json
{
    "data-root": "/mnt/docker"
}
```

9. Start the docker daemon:

```bash
systemctl start docker
```

10. Build the images:

```bash
su - codapi
make images
```

11. Verify that docker can now limit the storage size:

```bash
docker run -it --storage-opt size=16m codapi/alpine /bin/df -h | grep overlay
```
