# Slurm Docker Cluster

This is a multi-container Slurm cluster using docker-compose.
The compose file creates named volumes for persistent storage
of MySQL data files as well as Slurm state and log directories.

For a quick end to end demo, including rebuilding the container
image, please run:
```console
$ ./demo.sh
```

Many ideas on how to package Slurm came from here:
https://github.com/giovtorres/slurm-docker-cluster

## Containers and Volumes

The compose file will run the following containers:

* mysql
* slurmdbd
* slurmctld
* c1 (slurmd)
* c2 (slurmd)

The compose file will create the following named volumes:

* etc_munge         ( -> /etc/munge     )
* etc_slurm         ( -> /etc/slurm     )
* slurm_jobdir      ( -> /data          )
* var_lib_mysql     ( -> /var/lib/mysql )
* var_log_slurm     ( -> /var/log/slurm )

## Building the Docker Image

Build the image locally:

```console
$ docker-compose build
```

## Starting the Cluster

Run `docker-compose` to instantiate the cluster:

```console
$ docker-compose up -d
```

## Register the Cluster with SlurmDBD

To register the cluster to the slurmdbd daemon, run the `register_cluster.sh`
script:

```console
$ ./register_cluster.sh
```

> Note: You may have to wait a few seconds for the cluster daemons to become
> ready before registering the cluster.  Otherwise, you may get an error such
> as **sacctmgr: error: Problem talking to the database: Connection refused**.
>
> You can check the status of the cluster by viewing the logs: `docker-compose
> logs -f`

## Accessing the Cluster

Use `docker exec` to run a bash shell on the controller container:

```console
$ docker exec -it slurmctld bash
```

From the shell, execute slurm commands, for example:

```console
[root@slurmctld /]# sinfo
PARTITION AVAIL  TIMELIMIT  NODES  STATE NODELIST
normal*      up 5-00:00:00      2   idle c[1-2]
```

You can check the burst buffer is reporting correctly:

```console
[root@slurmctld /]# scontrol show burstbuffer
Name=cray DefaultPool=dwcache Granularity=16MiB TotalSpace=32GiB FreeSpace=32GiB UsedSpace=0
  AltPoolName[0]=test_pool Granularity=16MiB TotalSpace=32GiB FreeSpace=32GiB UsedSpace=0
  Flags=EnablePersistent
  StageInTimeout=30 StageOutTimeout=30 ValidateTimeout=5 OtherTimeout=300
  AllowUsers=root,slurm
  GetSysState=/opt/cray/dw_wlm/default/bin/dw_wlm_cli
```

## Submitting Jobs

The `slurm_jobdir` named volume is mounted on each Slurm container as `/data`.
Therefore, in order to see job output files while on the controller, change to
the `/data` directory when on the **slurmctld** container and then submit a job:

```console
[root@slurmctld /]# cd /data/
[root@slurmctld data]# sbatch --wrap="uptime"
Submitted batch job 2
[root@slurmctld data]# ls
slurm-2.out
[root@slurmctld data]# srun -n2 hostname
c1
c2
```

To create a burst buffer you need to be the slurm user, not root:

```console
su slurm
srun --bb="capacity=1G" hostname
```

To update the burst buffer python code and run a test job run:

```console
./update_burstbuffer.sh
docker exec slurmctld bash -c "cd /data && su slurm -c 'srun --bb=\"capacity=1G\" bash -c \"set\"'"
```

## Stopping and Restarting the Cluster

```console
$ docker-compose stop
```

```console
$ docker-compose start
```

## Deleting the Cluster

To remove all containers and volumes, run:

```console
$ docker-compose rm -sf
$ docker volume rm slurmdockercluster_etc_munge slurmdockercluster_etc_slurm slurmdockercluster_slurm_jobdir slurmdockercluster_var_lib_mysql slurmdockercluster_var_log_slurm
```
