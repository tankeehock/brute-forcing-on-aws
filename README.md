# brute-forcing-on-aws
A proof-of-concept project that automatically provisions EC2 instances and perform specific tasks (cracking access keys, brute-forcing S3 buckets, etc).

The codes in this repository is focused on cracking AWS access keys. The full narrative for this demonstration can be found in this [article](https://medium.com/csg-govtech/lets-try-to-crack-some-aws-credentials-e2fb5359ef3b).

## Cracker
The cracker can be found in `/app` folder with the compiled executable used in my experiment. However, do note that the folder is working directory (may not be the most updated codes from the owner). The repository of the cracker can be found [here](https://github.com/violenttestpen/aws-key).

## Deployment

### Validation

``` bash
terraform init
terraform fmt
terraform validate
```

### Setup and Tear-down

Please ensure that all the configurations in `configurations.tfvars` are accurate.

``` bash
terraform apply -var-file="configurations.tfvars" -auto-approve
terraform output ssh_private_key
terraform destroy -var-file="configurations.tfvars" -auto-approve
```

## Credits

Thanks [@violenttesten](https://github.com/violenttestpen) for his cracker.
