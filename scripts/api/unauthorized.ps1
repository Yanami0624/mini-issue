param(
	[string]$BaseUrl = "http://localhost:8080"
)

. "$PSScriptRoot/common.ps1"

$base = Resolve-ApiBaseUrl $BaseUrl

Write-Host "Case 1: GET /me without Authorization header" -ForegroundColor Yellow
Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/me"

Write-Host ""
Write-Host "Case 2: GET /me with invalid Bearer token" -ForegroundColor Yellow
Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/me" `
	-Headers @{
		Authorization = "Bearer invalid-token"
	}

