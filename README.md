# tdtidy

A command line tool for managing ECS task definitions.

`tdtidy` can deregister and delete old task definitions.

## Prerequisite

The following AWS IAM permissions are required.

- [ecs:DeleteTaskDefinitions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DeleteTaskDefinitions.html)
- [ecs:DeregisterTaskDefinition](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DeregisterTaskDefinition.html)
- [ecs:DescribeTaskDefinition](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DescribeTaskDefinition.html)
- [ecs:ListTaskDefinitions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_ListTaskDefinitions.html)

## Installation

### Homebrew

```
$ brew install manabusakai/tap/tdtidy
```

### go install

```
$ go install github.com/manabusakai/tdtidy/cmd/tdtidy@latest
```

## Usage

```
USAGE:
  tdtidy subcommand [options]

SUBCOMMANDS:
  deregister  Deregisters one or more task definitions.
  delete      Deletes one or more task definitions. You must deregister a task definition before you delete it.

OPTIONS:
  -dry-run
        Turn on dry-run. Output the target task definitions.
  -family-prefix string
        Specify the family name of the task definitions. If specified, filter by family name.
  -retention-period int
        The retention period for task definitions is specified in days. The unit is the number of days, and the default value is zero.
```

## License

MIT License

## Author

Copyright (c) 2023 Manabu Sakai
