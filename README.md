# release-instructions

## Configuration

```YAML
HTTPListen: :8081
jira:
  host: "https://jira.my.net"
  user: "user"
  password: "pass"
bitbucket:
  user: "user"
  password: "pass"
  owner: "owner"
  repos:
    - name: myprj
      weight: 2
      branch: staging
      excludeBranch: live
    - name: myanotherprj
      weight: 10
      branch: staging
      excludeBranch: live
```