locals {  
  // Load team members in variables.yaml
  loaded_team_members = flatten([
    for team in yamldecode(file("variables.yaml")).teams : [
      for tn, t in github_team.all : {
        name    = t.name
        id      = t.id
        slug    = t.slug
        members = team.member
        maintainers = team.maintainer
      } if t.name == team.name
    ]
  ])

  // Filter members having "member" role
  member_users = flatten([ for team in local.loaded_team_members : [
      for name in team.members : {
        identifier = "${team.slug}-${name}"
        team_id = team.id
        username = name
        role = "member"
      }
    ]
  ])

  // filter members having "maintainer" role
  maintainer_users = flatten([ for team in local.loaded_team_members : [
      for name in team.maintainers : {
        identifier = "${team.slug}-${name}"
        team_id = team.id
        username = name
        role = "maintainer"
      }
    ]
  ])

  all_users = concat(local.member_users, local.maintainer_users)

  // Load repositories and associate teams in variables.yaml
  // yasuo repo
  yasuo_teams = flatten([
    for  team in yamldecode(file("variables.yaml")).repositories.yasuo : [
        for tn, t in github_team.all : {
          identifier="yasuo-${tn}-${team.permission}"
          team_id = t.id
          permission = team.permission
        } if t.name == team.name
    ]
  ])
  // yone repo
  yone_teams = flatten([
    for  team in yamldecode(file("variables.yaml")).repositories.yone : [
        for tn, t in github_team.all : {
          identifier="yone-${tn}-${team.permission}"
          team_id = t.id
          permission = team.permission
        } if t.name == team.name
    ]
  ])
}