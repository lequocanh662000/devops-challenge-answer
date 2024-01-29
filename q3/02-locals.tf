locals {  
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


  member_users = flatten([ for team in local.loaded_team_members : [
      for name in team.members : {
        identifier = "${team.slug}-${name}"
        team_id = team.id
        username = name
        role = "member"
      }
    ]
  ])

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
}