# mikros-cli

## About

`mikros-cli` is a CLI utility that helps a user creating or editing services
from the command line.

## Installing

In order to install the CLI locally, one can execute the following command:

```bash
go install github.com/somatech1/mikros-cli@latest
```

## Creating a new service template

After installing the `mikros-cli` command locally, is possible to create
new service templates with it by executing in the following way:

```bash
mikros-cli service init
```

This will execute a little survey, where mandatory information must be
entered to successfully generate the templates.

And, if everything executed the way it should, you should have a new folder,
with the service name entered at the survey and with some source files in
it.

## Roadmap

WIP

## License

[Mozilla Public License 2.0](LICENSE)
