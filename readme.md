# gcloud client

 scripts to do things which you wish google console client does.

## Add SSH Key to instances
Add your ssh key to all google compute instances. `ssh_key` flag is mandatory.
- `gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub`

The keys are written to tempdir and cleared after added.

You could customize the flags
- `--limit` total instances to add
- `--filter` regexp to filter the instances while listing. uses gcloud `filter=name~'regex'`
- `--user` username to which your ssh key is added for, defaults to `$USER`


`gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub --filter='.*pg.*' --limit=1 --user username`


