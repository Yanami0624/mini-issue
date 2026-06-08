param(
	[string]$BaseUrl = "http://localhost:8080",
	[string]$Username = "",
	[string]$Password = "123456"
)

. "$PSScriptRoot/common.ps1"

$base = Resolve-ApiBaseUrl $BaseUrl
if ([string]::IsNullOrWhiteSpace($Username)) {
	$Username = "user_" + (Get-Date -Format "yyyyMMddHHmmss")
}

Write-Host "Step 1: register user $Username" -ForegroundColor Yellow
Invoke-MiniIssueApi `
	-Method "POST" `
	-Url "$base/register" `
	-Body @{
		username = $Username
		password = $Password
	}

Write-Host ""
Write-Host "Step 2: login user $Username" -ForegroundColor Yellow
$loginResponse = Invoke-MiniIssueApi `
	-Method "POST" `
	-Url "$base/login" `
	-Body @{
		username = $Username
		password = $Password
	}

if ($loginResponse.StatusCode -ne 200 -or $null -eq $loginResponse.Json.data.token) {
	throw "Login did not return a token. Stop before GET /me."
}

$token = $loginResponse.Json.data.token

Write-Host ""
Write-Host "Step 3: GET /me with JWT" -ForegroundColor Yellow
Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/me" `
	-Headers @{
		Authorization = "Bearer $token"
	}

