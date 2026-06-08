param(
	[string]$BaseUrl = "http://localhost:8080",
	[string]$Username = "alice",
	[string]$Password = "123456",
	[string]$TokenFile = ".mini-issue-token.json"
)

. "$PSScriptRoot/common.ps1"

$base = Resolve-ApiBaseUrl $BaseUrl

$response = Invoke-MiniIssueApi `
	-Method "POST" `
	-Url "$base/login" `
	-Body @{
		username = $Username
		password = $Password
	}

if ($response.StatusCode -eq 200 -and $null -ne $response.Json.data.token) {
	$tokenState = @{
		base_url = $base
		username = $Username
		token    = $response.Json.data.token
	}

	$tokenState | ConvertTo-Json -Depth 20 | Set-Content -Path $TokenFile -Encoding UTF8
	Write-Host ""
	Write-Host "Saved token to $TokenFile" -ForegroundColor Yellow
}

