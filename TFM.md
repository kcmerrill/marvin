# The Manual

## Some examples

These are some samples plucked from my own personal `marvin.yml` file.

### Github Vanity

I typically will open up my traffic on various github projects. Instead of going to each project individual repo's traffic page, I can open them up pretty easily.

Note in this example I'm using [kronk](https://github.com/kcmerrill/kronk), which is a simple regex tool I wrote that also happens to work awesome with marivn ;)

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
