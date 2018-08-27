# gcloud client

 scripts to do things which you wish google console client does.

## Installation
- required `gcloud` installed, authenticated with project config set
- install lib with `go get -u github.com/devdinu/gcloud-client`

## Usage
- `gcloud-client --help` to show the help with flag information

## Add SSH key to compute instances

You could add your ssh key to compute instance[s], so you could ssh directly. default [gcloud compute add ssh](https://cloud.google.com/compute/docs/instances/adding-removing-ssh-keys) overrides the existing keys.
you can add your key along with existing ones in instances with this command.

### Add your ssh key to all google compute instances.
`ssh_key` flag is mandatory.
- `gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub`

The existing keys with new key is written to temp file and cleared after added to the instance.

You could customize the flags
- `--limit` total instances to add
- `--filter` regexp to filter the instances while listing. uses gcloud `filter=name~'regex'`
- `--user` username to which your ssh key is added for, defaults to `$USER`


```
gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub --filter='.*pg.*' --limit=1 --user username
```

### Add to single instance
You could give `--instance` and `--zone` to add ssh key to single instance, as its faster than listing instances with regexp

```
gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub --instance=some_instance --zone=asia-zone
```
