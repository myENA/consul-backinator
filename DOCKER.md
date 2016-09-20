# Why Dockerize this utility?

So glad you asked! The use cases for Dockerizing this utility are:

* Another distribution mechanism
* Ability to run wherever you can run docker (Mesos, Kubernetes, Docker Swarm,
Docker Engine, etc.)

You may be thinking, "Well, duh! Of course you can run on Docker Engine now,
but _why_ would I ever want to do that?" The ability to run on anything that can
run Docker makes it quite powerful to offload this task to different services
which may be hosted. Consider if you had a service running on Docker 1.12+ that
had its own private overlay network and a Consul cluster in that network. One
could then use this Docker image as a Docker service that connects to this
network, points at the Consul service in the Docker network, and creates a
backup to a DFS volume mounted in the container. This Consul backup Docker
service could then have a `--restart-delay` of `10m` which would essentially
take backups of your Consul cluster every 10 minutes, dumping it to the DFS of
your choice.

There are probably many more, use your imagination and see what you come up
with. If you find it beneficial for others to know, submit a PR updating this
doc.

Thanks.
