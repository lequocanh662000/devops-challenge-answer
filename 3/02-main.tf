resource "github_repository" "example_repo" {
  name = "example-repo"
}

resource "vault_generic_secret" "github_credentials" {
  path = "secret/github"
  data = {
    token = "your-github-token"
  }
}