# Wiki
The `local.team_members` value depends on the `github_team.all` resource, so you must create your teams before you can add members to them.

### Pre-requirements: 
1. Set 2 env variables `GITHUB_TOKEN`(personal access token), `GITHUB_OWNER`(organization)
```
export GITHUB_TOKEN=<sensitive>
export GITHUB_OWNER=<sensitive>
```

2. Define teams & members in `variables.yaml`

### Apply configuration
Apply this terraform manifest with following orders:
1. Define member in the organization
```
terraform apply -target github_team.all
```

2. Assign members into teams
```
terraform apply
```