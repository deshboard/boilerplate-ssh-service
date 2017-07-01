# Documentation

This folder contains documentation for the service. It mainly includes information about setting up and running the development environment.


## Table of Contents

- [Requirements](#requirements)
- [Framework overview](#framework-overview)
- [Setup](setup.md)
- [Testing](testing.md)


## Requirements

- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/)
- [Go](https://golang.org/) (1.8 or above)
- [Glide](http://glide.sh/)


### Optional

These are not hard requirements of the project.

- [Godotenv](https://github.com/joho/godotenv) (recommended for easier env setup)
- [Reflex](https://github.com/cespare/reflex) (required for watching code changes)
- [Docker Compose](https://docs.docker.com/compose/)


Please make sure that you have the latest versions installed.


## Framework overview

This project does not use any third-party framework (except ones required by the application logic), but relies heavily on the standard library and separate third-party components. The integration layer for these components and the main execution logic can be found in the [main/](../main/) directory.
