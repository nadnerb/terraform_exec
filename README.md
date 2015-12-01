terraform_exec
=============

Terraform wrapper that allows terraform projects to have multiple environments, synced to S3.

For example terraform_exec allows a terraform elasticsearch project to have staging and production environments in multiple aws regions.

## Usage:

First you need to install [terraform](https://terraform.io)

Next install terraform_exec ([golang required](https://golang.org/) and setup)

```
go get github.com/nadnerb/terraform_exec
go install github.com/nadnerb/terraform_exec
terraform_exec --help
```

### terraform_exec run

By default you will run `terraform_exec` within an existing terraform project. It will sync local state with s3 similar to atlas, additionally supporting
multiple 'environments'.

The run command allows running normal terraform commands such as `plan`, `apply`, `refresh` and `destroy`.

e.g `terraform_exec run plan staging`

### Configuration

All `terraform_exec` commands will look in the ./config directory for a `staging.tfvars` file. At a minimum it will need the following variables to
save state to s3:

```
aws_region="ap-southeast-2"
s3_bucket="a-bucket"
s3_key="an-s3-key"
```

#### AWS security

Out of the box, `terraform_exec` will look for AWS credentials set in environment variables. If running on an ec2 box in AWS, retrieving credentials
via the machines IAM role are supported:

```
terraform_exec run plan staging --security=aws-internal --security-role=your-iam-role
```

Use `terraform_exec run --help` for more details.

#### S3 sync

If for some reason you need to skip the inital sync with s3, the `--no-sync=true` flag can be used.

### terraform_exec upload

Upload existing environment state to s3

See `terraform_exec upload --help` for more details

### terraform_exec download

Download existing environment state from s3

See `terraform_exec download --help` for more details

### Testing terraform_exec

```shell
$ go test ./...
```

### Issues

When testing

`cannot find package "github.com/stretchr/testify/assert" in any of: ...`

You will need to

```shell
$ go get github.com/stretchr/testify/assert
```

### TODO

* use terraform go project (most likely need to change to its cli)
* improve this readme!!
* remove unnessessary s3_Key variable
* only warn for inital sync if new project, fail otherwise
