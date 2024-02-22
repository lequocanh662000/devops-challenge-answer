# Create infrastructure repository
# Add memberships for infrastructure repository (rm stands for repository member)
// yasuo
resource "github_repository" "yasuo" {
  name = "yasuo-deadth-is-like-wind"
}

resource "github_team_repository" "yasuo" {
  for_each = { for rm in local.yasuo_teams : rm.identifier => rm}

  repository = github_repository.yasuo.id
  team_id    = each.value.team_id
  permission = each.value.permission
}

// yone
resource "github_repository" "yone" {
  name = "yone-yasuo-bros"
}

resource "github_team_repository" "yone" {
  for_each = { for rm in local.yone_teams : rm.identifier => rm}

  repository = github_repository.yone.id
  team_id    = each.value.team_id
  permission = each.value.permission
}
