resource "github_membership" "all" {
  for_each = {
    for user in yamldecode(file("variables.yaml")).users:
    user.username => user
  }

  username = each.value.username
  role     = each.value.role
}