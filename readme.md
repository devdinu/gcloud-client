# gcloud client

 scripts to do things which you wish google console client does.

## Installation

### Required libs
- [gcloud](https://cloud.google.com/sdk/gcloud) authenticated with projects and config set
- [tmux](https://github.com/tmux/tmux)
- [tmuxinator](https://github.com/tmuxinator/tmuxinator)

### Go
`go get -u github.com/devdinu/gcloud-client`

### Homebrew
`brew install devdinu/devlife/gcloud-client`


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
gcloud-client --ssh_key=$HOME/.ssh/id_rsa.pub --filter='.*pg.*' --limit=10 --user username
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
gcloud-ciient ssh revoke --name=someone --prefix= # Unimplemented

// searching
gcloud-client instances search --prefix vm-prefix-to-search --project=project
gcloud-client instances search --regex some.*regex
gcloud-client instances refresh --timeout
gcloud-client instances list --project=specific-project  # WIP

// ssh
gcloud-client instances ssh --prefix some-prefix --user username --session session_name
gcloud-client instances ssh --regex some-prefix # WIP
gcloud-client instances ssh --tag some-prefix   # WIP
 
```


Todo:
* [ ] Use absolute path (sensible default) for db file & configs if required
* [ ] Setup script to install gcl and tmuxinator template and tmuxinator if needed (2.7)
* [ ] CI to build binary
* [ ] enable os.Stdin in cmd.execute
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
* [X] use logger package
* [X] Define a list of global flags
* [ ] Add gci command to do ls, switch projects

