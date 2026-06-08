param(
	[string]$BaseUrl = "http://localhost:8080",
	[string]$Username = "alice",
	[string]$Password = "123456"
)

. "$PSScriptRoot/common.ps1"

$base = Resolve-ApiBaseUrl $BaseUrl

Invoke-MiniIssueApi `
	-Method "POST" `
	-Url "$base/register" `
	-Body @{
		username = $Username
		password = $Password
	}

