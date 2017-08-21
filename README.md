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
    mysql -u :user -p :password -h :host -e ":args"
inventory:
    static: |
      db:master host:master.db.kcmerrill.com user:db_user password:$PASSWORD
      db:replica host:replica.db.kcmerrill.com user:db_user password:$PASSWORD
      db:manual-query host:manual-query.db.kcmerrill.com user:db_user password:$PASSWORD
```
