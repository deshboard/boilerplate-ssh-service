# Setup

## Automatic

The project comes with a Makefile which includes everything for starting the project.
The easiest way to go is simply executing the following command:

``` bash
$ make setup
```

It sets up the project with decent defaults (installs all dependencies, sets the default configuration, etc).

If you want to clean the project and give it a fresh start:

``` bash
$ make clean
```


## Manual

Although in most cases the automatic installation should be just fine, it's possible to go through the same steps manually.
Consider this part an explanation of the section above.

As the minus one step you might want to check if every dependency is installed:

``` bash
$ make envcheck
```

First thing you need to do is installing the dependencies. This project uses [Glide](http://glide.sh/) for dependency management.

``` bash
$ glide install # Or make install
```

Next you need to setup the environment configuration. You can use `.env.example` as a base, it contains the default values.
For testing there is a separate environment which usually can be the same as the development environment.

``` bash
$ cp .env.example .env # or make .env
$ cp .env.example .env.test # or make .env.test
```

Cleaning up is as easy as deleting the files created above:

``` bash
$ rm -rf vendor/ .env .env.test
```


## Docker

Some projects may come with a Docker environment as well in form of a [Docker Compose](https://docs.docker.com/compose/) config.
Using Docker to create a disposable development environment is optional, but it's usually easier,
especially when the project has a lot of dependencies.

Since it's optional, if the project contains a Docker setup as well you need to start it manually:

``` bash
$ docker-compose up -d
```

There are usually custom configurations you can make (port mappings, volumes for persistence).
By default the project should just work OOTB without those, but if you want to access for example the database
or want to persist it during the development, so your data is not lost between restarts, it's recommended to make those configurations.

As always, sane defaults should be provided in a `docker-compose.override.yml.example` file.

``` bash
$ cp docker-compose.override.yml.example docker-compose.override.yml
```

**Note:** Although custom Docker config is optional you might need it to run integration or acceptance tests.

If you develop multiple projects simultaneously it makes sense to stop the running containers when switching projects:

``` bash
$ docker-compose stop
```

The Docker environment is disposable, you can easily destroy it and give it a fresh start:

``` bash
$ docker-compose down
$ rm -rf .docker/
```
