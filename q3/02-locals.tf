locals {  
  team_members = flatten([
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


  members = flatten([ for team in local.team_members : [
      for name in team.members : {
        identifier = "${team.slug}-${name}"
        team_id = team.id
        username = name
        role = "member"
      }
    ]
  ])

  maintainers = flatten([ for team in local.team_members : [
      for name in team.maintainers : {
        identifier = "${team.slug}-${name}"
        team_id = team.id
        username = name
        role = "maintainer"
      }
    ]
  ])

  all_users = concat(local.members, local.maintainers)
}