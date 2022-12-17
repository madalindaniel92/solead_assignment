# Setup elasticsearch with kibana for development and testing

In order to install elasticsearch with kibana on Linux Manjaro, I used [snapcraft](https://snapcraft.io/install/docker/manjaro)

```
sudo pacman -S snapd

sudo systemctl enable --now snapd.socket

sudo ln -s /var/lib/snapd/snap /snap

sudo snap install docker
```

## Issues encountered

While starting the local cluster using docker-compose I encountered the
following issue:
```
bootstrap check failure [1] of [1]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
```
log obtained by running `docker-compose logs es01`

https://stackoverflow.com/questions/51445846/elasticsearch-max-virtual-memory-areas-vm-max-map-count-65530-is-too-low-inc#answer-51448773

```
sysctl -w vm.max_map_count=262144
```

If you want to set this permanently, you need to edit `/etc/sysctl.conf` and set vm.max_map_count to 262144.

## Resources and links
- https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html

- https://www.elastic.co/guide/en/kibana/current/docker.html

