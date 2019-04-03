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
`gcloud-client --limit=10` adds your key to 10 instances

The existing keys with new key is written to temp file and cleared after added to the instance.

You could customize the flags
- `--limit` total instances to add
- `--filter` regexp to filter the instances while listing. uses gcloud `filter=name~'regex'`
- `--user` username to which your ssh key is added for, defaults to `$USER`
- `--ssh_key` ssh_key file to be uploaded, defaults to `$HOME/.ssh/id_rsa.pub`


```
gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub --filter='.*pg.*' --limit=1 --user username
```

### Add to single instance
You could give `--instance` and `--zone` to add ssh key to single instance, as its faster than listing instances with regexp

```
gcloud-client --instance=some_instance --zone=asia-zone

```


## CheatSheet
globals flags:
timeout

```
gcloud-client ssh grant --key=someone-key.pub --prefix=
gcloud-ciient ssh revoke --name=someone --prefix=

// searching
gcloud-client instances search --prefix vm-prefix-to-search --project=project
gcloud-client instances search --regex some.*regex
gcloud-client instances refresh --timeout
gcloud-client instances list --project=specific-project

// ssh
gcloud-client instances ssh --prefix some-prefix
gcloud-client instances ssh --regex some-prefix
gcloud-client instances ssh --tag some-prefix
 
```


Todo:
SSH ACCESS
* [ ] Display IP after adding the key
* [ ] Adding IP, name mapping to the /etc/hosts file
* [ ] Remove particular ssh key with id
* [ ] Infer id from the sshkey
SEARCH
* [ ] Display progress bar on refresh projects
* [ ] store state (Terminated) information and ignore in search, or show
* [ ] Use knife tags to tag instances than manual
* [ ] Optimize storing in db performance
CODE/FEATURES
* [ ] use logger package
* [ ] Define a list of global flags
* [ ] Add gci command to do ls, switch projects

