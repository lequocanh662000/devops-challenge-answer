resource "github_membership" "all" {
  for_each = {
    for member in yamldecode(file("variables.yaml")).members:
    member.username => member
  }

  username = each.value.username
  role     = each.value.role
}