# The Manual

## Some examples

These are some samples plucked from my own personal `marvin.yml` file.

#### Github Vanity

Previously, if I wanted to check my repository stats, I'd manually click on all my repos, then go to insights then traffic. Another example of how I am using marvin to check my repo stats painlessly. 

Note in this example I'm using [kronk](https://github.com/kcmerrill/kronk), which is a simple regex tool I wrote that also happens to work awesome with marvin!

```
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
        repos: |
            curl -s https://api.github.com/users/kcmerrill/repos?per_page=100 | kronk 'repo:"full_name": "(.*?)"'
```
