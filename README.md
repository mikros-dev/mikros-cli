# mikros CLI

## About

`mikros` is a CLI utility that helps a user managing its mikros environment
from the command line.

## Installing

In order to install the CLI locally, one can execute the following command:

```bash
go install github.com/mikros-dev/mikros-cli/cmd/mikros@latest
```

## Creating a new project

After installing the `mikros` command locally, it is possible to create
new project templates with it by executing in the following way:

```bash
mikros new
```

This will execute a little survey, where mandatory information must be
entered to successfully generate the desired templates.

And, if everything executed the way it should, you should have a new folder,
with the project selected at the survey and with some source files in it.

## Roadmap

* ~~Change main command to `new`~~
* Full support for rust services
* ~~Support for creating protobuf projects~~
* Support for creating services monorepo projects
* ~~Command for creating protobuf file from templates~~
* ~~Use [bubbletea](https://github.com/charmbracelet/bubbletea) for UI surveys.~~
* ~~Add command to generate default config file~~

## License

[Mozilla Public License 2.0](LICENSE)
