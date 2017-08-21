[![Build Status](https://travis-ci.org/kcmerrill/marvin.svg?branch=master)](https://travis-ci.org/kcmerrill/marvin) [![Go Report Card](https://goreportcard.com/badge/github.com/kcmerrill/marvin)](https://goreportcard.com/report/github.com/kcmerrill/marvin)
![marvin](assets/marvin.jpg "marvin")

# Marvin

An inventory based task runner. Inspired by [Knife](https://docs.chef.io/knife.html) and [Ansible](https://www.ansible.com/). Marvin allows you to define the tasks you want to run, how you want to run them and where depending on how you define and setup your tasks.

```bash
# marvin usage
$> marvin <inventory:filter> <taskname> <any additional args you might want to use>
```

Here is an example how you can use marvin to run queries across multiple databases.

```bash
# run a query on all databases in the inventory
$> marvin db:* query select count(*) from tablename

# run a query on just the master database
$> marvin db:master query select count(*) from tablename
```

The previous commands are made possible given the following `marvin.yml` configuration file.

```yaml
tasks:
  query: |
    mysql -u {{ .user }} -p {{ .password }} -h {{ .host }} -e "{{ .args }} "
inventory:
    dynamic: 
      files: ls -1
    static: |
      db:master host:master.db.kcmerrill.com user:db_user password:$PASSWORD
      db:replica host:replica.db.kcmerrill.com user:db_user password:$PASSWORD
      db:manual-query host:manual-query.db.kcmerrill.com user:db_user password:$PASSWORD
```

## Built in tasks

By default, there are a few built in tasks. You can overide these if you'd like, but by default you get `ls` and `ssh` as described below.

1. `ls` will show all available inventory, both dynamic and static
1. `ssh` will allow you to run a command with `{{ .args }}` on remote hosts
1. Need more? Submit a PR ...

## Dyanmic Inventory

You can specify commands, that when filtered by, will run commands to generate dynamic inventory. A great example would be, generating ec2 instances inventory on the fly, then running commands on said instances.

```bash
$> marvin : ls #display all available inventory
$> marvin *: ls #display all available inventory
$> marvin *:* ls #display all available inventory

$> marvin env:prod ssh whoami #ssh {{ .host }} {{ .args }}
```

## Binaries && Installation

Feel free and use `go get` if you already have golang instaled. If not, feel free and download a compiled binary just for you and your OS: 

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")](http://go-dist.kcmerrill.com/kcmerrill/marvin/mac/amd64) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")](http://go-dist.kcmerrill.com/kcmerrill/marvin/linux/amd64)