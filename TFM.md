# The Manual

## Some examples

These are some samples plucked from my own personal `marvin.yml` file.

#### Github Vanity

Previously, if I wanted to check my repository stats, I'd manually click on all my repos, then go to insights then traffic. Another example of how I am using marvin to check my repo stats painlessly. 

Note in this example I'm using [kronk](https://github.com/kcmerrill/kronk), which is a simple regex tool I wrote that also happens to work awesome with marvin!

```bash
# look at all repository traffic
$ marvin repos:* traffic

# only look at my go projects
$ marvin repos:*.go traffic

# only look at marvin's traffic
$ marvin repos:marvin traffic
```

```yaml
tasks:
  clone: |
    git clone git@github.com:{{ .id }}.git
  traffic: |
      open https://github.com/{{ .id }}/graphs/traffic
inventory:
    dynamic:
        myorg: |
            curl -s https://api.github.com/orgs/myorg/repos?per_page=100 | kronk 'repo:"full_name": "(.*?)"'
        repos: |
            curl -s https://api.github.com/users/kcmerrill/repos?per_page=100 | kronk 'repo:"full_name": "(.*?)"'
```

#### EC2 instance runner

Manage a bunch of AWS EC2 instances? Need to do a quick ssh command on each of them? By default, `marvin` gives you a handy `ec2` dyanmic inventory baked in. In doing so, it will return a list of ec2 instance names that are running. 

At this point, it's just a matter of hooking up a task with the generated inventory. In this case, we will override the default `ssh` task, and configure it so it's magic for our specific system. 

```yaml
tasks:
    ssh: |
        ssh -q -i ~/.ssh/user -o StrictHostKeyChecking=no user@{{ .id }}.mydomain.com {{ .args }}
    query: |
        mysql --login-path=master database -e '{{ .args }}'
inventory:
    dynamic:
        myorg: |
            curl -s https://api.github.com/orgs/myorg/repos?per_page=100 | kronk 'repo:"full_name": "(.*?)"'
        repo: |
            curl -s https://api.github.com/users/kcmerrill/repos?per_page=100 | kronk 'repo:"full_name": "(.*?)"'
    static: |
        db:master
        db:read1
        db:reporting
        db:manual-queries
```

```bash
# restart a service on all of the machines
$ marvin ec2:* ssh service httpd restart

# shutdown autoscale machines 1 && 2
$ marvin "ec2:autoscale(1|2)" ssh shutdown now
# quotes are needed because bash is no fan of parens
```
