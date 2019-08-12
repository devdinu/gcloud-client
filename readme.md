# gcloud client
[![Go Report Card](https://goreportcard.com/badge/github.com/devdinu/gcloud-client)](https://goreportcard.com/report/github.com/devdinu/gcloud-client)

 scripts to do things which you wish google console client does.

## Installation

### Required libs
- [gcloud](https://cloud.google.com/sdk/gcloud) authenticated with projects and config set
- [tmux](https://github.com/tmux/tmux)
- [tmuxinator](https://github.com/tmuxinator/tmuxinator) gem

### Go
`go get -u github.com/devdinu/gcloud-client`

### Homebrew
`brew install devdinu/devlife/gcloud-client`

If you've installed with brew, command will be `gcl`, might need to remove alias from git plugin of oh-my-zsh `~/.oh-my-zsh/plugins/git/git.plugin.zsh`

## Usage

### Refresh
all instances, which's stored in boltd
- `gcloud-client instances refresh`

optionally you can pass flag `--projects proj1,proj2` to refresh specific projects
You could add the refresh as cron

![Demo](https://raw.githubusercontent.com/devdinu/gcloud-client/master/demo/refresh-demo.gif)

### Search

- `gcloud-client --help` to show the help with flag information


### Add SSH key to compute instances

You could add your ssh key to compute instance[s], so you could ssh directly. default [gcloud compute add ssh](https://cloud.google.com/compute/docs/instances/adding-removing-ssh-keys) overrides the existing keys.
you can add your key along with existing ones in instances with this command.

```
# adds your ~/.ssh/id_rsa.pub key to 10 instances
gcloud-client ssh_access --limit=10
```

The existing keys with new key is written to temp file and cleared after added to the instance.

You could customize the flags
- `--limit` total instances to add
- `--filter` regexp to filter the instances while listing. uses gcloud `filter=name~'regex'`
- `--user` username to which your ssh key is added for, defaults to `$USER`
- `--ssh_key` ssh_key file to be uploaded, defaults to `$HOME/.ssh/id_rsa.pub`
- `--dbfile` file to store the instances and search, defaults to `$HOME/hosts.db`
- `--projects` list of project-ids to search for while login, or to refresh


```
# customize ssh file, username
gcloud-client ssh_access --ssh_key=$HOME/.ssh/gcp_id_rsa.pub --filter='.*pg.*' --limit=10 --user username
```

#### Add to single instance
You could give `--instance` and `--zone` to add ssh key to single instance, as its faster than listing instances with regexp

```
gcloud-client --instance=some_instance --zone=asia-zone
```

![Demo](https://raw.githubusercontent.com/devdinu/gcloud-client/master/demo/ssh-access.gif)

### Login to instances

ssh into the instances which you searched, open each in tmux pane.

```
gcloud-client instances login --regex=some-prefix
gcloud-client instances login --prefix='some-.*db.*' --user username --session session_name
gcl instances search --regex=".*kafka.*" --projects=proj-integration,proj-production
```
* prefix search is much faster than regex search

Customize flags
- `--user` username for ssh
- `--session` tmux session name, so multiple ssh sessions can be done at same time
- `--projects` if you've many gcloud projects, mention the projects to search

![Demo](https://raw.githubusercontent.com/devdinu/gcloud-client/master/demo/login.gif)

## Additional information

### Homebrew
if you've installed via homebrew, use `gcl` instead of `gcloud-client`
`/usr/local/Cellar/gcloud-client/` have the template file and bin directory with the executable

### Tmux
enable syncronized panes, so most cases require you to run command on all machines
```
:set synchronized-panes on
```

### TODO:
* customize which `cmd` to run once ssh to all instances
* customize login via external ip / internal ip via flags (currently only internal ip is supported)
* output machine ip mapping so can be added to `/etc/hosts`, or add to file
* Infer id from ssh key
* revoking ssh key for user from machine
* ssh customize project to use via flag
* TAG: search list could be tagged, so can login via tags
