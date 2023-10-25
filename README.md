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
$ tdtidy -help
Usage of tdtidy:
  -dry-run
        Turn on dry-run. List the target task definitions.
  -family-prefix string
        Family name of task definitions. If specified, filter by family name.
  -retention-period int
        Retention period for task definitions. Unit is number of days. The default value is zero.
```

## License

MIT License

## Author

Copyright (c) 2023 Manabu Sakai
