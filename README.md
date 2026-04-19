# HBT Beerienteering

Terraform stack and UI

## Tasks

### apply

directory: tf
requires: build-cmd, upload-cmd
environment: AWS_PROFILE=beerienteering

```shell
terraform apply -auto-approve -input=false -var-file=tfvars/pro.auto.tfvars
```

### build-cmd

requires: clean

```shell
go run ./scripts/build-cmd --zip
```

### build-ui

directory: ui

```shell
npm run build
```

### clean

```shell
rm -rf build/
```

### format

requires: format-ui, format-tf

### format-tf

directory: tf

```shell
terraform fmt -recursive -write .
```

### format-ui

directory: ui

```shell
npx prettier -w .
```

### init

directory: tf
environment: AWS_PROFILE=beerienteering

```shell
terraform init -backend-config=tfvars/pro.backend.tfvars
```

### plan

directory: tf
requires: build-cmd, upload-cmd
environment: AWS_PROFILE=beerienteering

```shell
terraform plan -var-file=tfvars/pro.auto.tfvars -input=false
```

### upload-cmd

environment: AWS_PROFILE=beerienteering

```shell
#!/bin/bash
source tf/tfvars/.env
BINARY_BUCKET=$BINARY_BUCKET go run ./scripts/upload-binaries
```

### upload-ui

requires: build-ui
environment: AWS_PROFILE=beerienteering
directory: ui

```shell
#!/bin/bash
source ../tf/tfvars/.env
aws s3 sync dist/ s3://${WEB_BUCKET}
```
