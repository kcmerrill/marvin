[![Build Status](https://travis-ci.org/kcmerrill/marvin.svg?branch=master)](https://travis-ci.org/kcmerrill/marvin) [![Go Report Card](https://goreportcard.com/badge/github.com/kcmerrill/marvin)](https://goreportcard.com/report/github.com/kcmerrill/marvin)
![marvin](assets/marvin.jpg "marvin")

# Marvin

An inventory based task runner. Inspired by [Knife](https://docs.chef.io/knife.html) and [Ansible](https://www.ansible.com/). Marvin allows you to define the tasks you want to run, how you want to run them and where depending on how you define and setup your tasks.

```bash
# run a query on all databases in the inventory
$> marvin db:* query select count(*) from tablename

# run a query on just the master database
$> marvin db:master query select count(*) from tablename
```

The previous commands are possible given the following `marvin.yml` configuration file.

```yaml
tasks:
  query: |
    mysql -u :user -p :password -h :db -e ":args"
inventory: 
  db:master host:master.db.kcmerrill.com user:db_user password:$PASSWORD 
  db:replica host:replica.db.kcmerrill.com user:db_user password:$PASSWORD
  db:manual-query host:manual-query.db.kcmerrill.com user:db_user password:$PASSWORD
```


## Inventory

Inventory can come from 3 ways `Dynamic`, `Static` and `stdin`. 

### `Static` Inventory

Simple strings that are `key:value` pairs on a single line. An example of mysql 