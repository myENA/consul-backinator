# Why dockerize this utility?

So glad you asked! The use cases for dockerizing this utility are:

* Another distribution mechanism
* Ability to run wherever you can run docker (mesos, kubernetes, docker swarm,
docker engine, etc.)

You may be thinking, "Well, duh! Of course you can run on docker engine now,
but _why_ would I ever want to do that?" The ability to run on anything that can
run docker makes it quite powerful to offload this task to different services
which may be hosted. Consider if you had a service running on docker 1.12 that
had its own private overlay network and a consul cluster in that network. One
could then use this docker image as a docker service that connects to this
network, points at the consul service in the docker network, and creates a
backup to a DFS that it has volume mounted in. This consul backup docker
service could then have a `--restart-delay` of `10m` which would essentially
take backups of your consul cluster every 10 minutes, dumping it to the DFS of
your choice.

There are probably many more, use your imagination and see what you come up
with. If you find it beneficial for others to know, submit a PR updating this
doc.

Thanks.
