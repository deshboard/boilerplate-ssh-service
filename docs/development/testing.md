# Testing

Go comes with a great set of tools in it's standard library and fortunately a [testing](https://golang.org/pkg/testing/)
"framework" is one of them.


## Running tests

As most of the other development related tasks, test running is also included in the `Makefile`.

```bash
$ make check
```

Besides running tests, the above code also runs code style checks. You can fix any code style breaks by executing:

```bash
$ gofmt -l -w -s path/to/files
```

You can run the tests only with the following command:

```bash
$ make test
```

You can also pass arguments to the test command using the `ARGS` variable.
For example the following command will make the tests to run in verbose mode.

```bash
$ make test ARGS="-v"
```

Ultimately you can fall back to the builtin go test command:

```bash
$ go test
```

In this case you either have to pass the package to be tested as an argument or you have to change to the directory of
the package you want to test. The above make commands run tests for all the packages in the project.


## Structure of tests

The builtin testing package provides a flexible way to write unit/integration tests,
benchmarks and even so called examples which then gets built into the documentation.

In order to add support for BDD style acceptance tests, this project depends on
[godog](https://github.com/DATA-DOG/godog). It allows you to write user stories using
[Gherkin](https://github.com/cucumber/cucumber/wiki/Gherkin) language.

Normally you would want to run unit level tests locally as in most of the cases as acceptance and integration
tests can take a long time to run. But of course from time to time you have to run those as well.
To separate different tests this project uses go [build tags](https://golang.org/pkg/go/build/#hdr-Build_Constraints).

In order to run those tests locally you need to execute the following commands:

```bash
$ make test TAGS="integration"
$ make test TAGS="acceptance"
```
