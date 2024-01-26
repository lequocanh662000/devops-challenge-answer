provider "github" {
  token = vault_generic_secret.github_credentials["token"]
  organization = "your-organization"
}