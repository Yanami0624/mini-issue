param(
	[string]$BaseUrl = "http://localhost:8080",
	[string]$Username = "",
	[string]$Password = "123456",
	[Nullable[int64]]$IssueId = $null
)

$ErrorActionPreference = "Stop"

. "$PSScriptRoot/../../scripts/api/common.ps1"

$base = Resolve-ApiBaseUrl $BaseUrl
if ([string]::IsNullOrWhiteSpace($Username)) {
	$Username = "user_" + (Get-Date -Format "yyyyMMddHHmmss")
}

Write-Host "Step 1: register user $Username" -ForegroundColor Yellow
$null = Invoke-MiniIssueApi `
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
	throw "Login did not return a token. Stop before authenticated requests."
}

$token = $loginResponse.Json.data.token
$authHeaders = @{
	Authorization = "Bearer $token"
}

Write-Host ""
Write-Host "Step 3: GET /me" -ForegroundColor Yellow
$null = Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/me" `
	-Headers $authHeaders

Write-Host ""
Write-Host "Step 4: POST /issues" -ForegroundColor Yellow
$null = Invoke-MiniIssueApi `
	-Method "POST" `
	-Url "$base/issues" `
	-Headers $authHeaders `
	-Body @{
		title = "First issue from CLI"
		content = "Created by examples/cli/powershell_issue_flow.ps1"
		status = "OPEN"
		priority = 1
	}

Write-Host ""
Write-Host "Step 5: GET /issues" -ForegroundColor Yellow
$listResponse = Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/issues?page=1&page_size=10" `
	-Headers $authHeaders

if ($null -eq $IssueId) {
	if ($listResponse.StatusCode -ne 200 -or $null -eq $listResponse.Json.data.list -or $listResponse.Json.data.list.Count -eq 0) {
		throw "Could not find an issue id from GET /issues. Pass -IssueId explicitly after creating an issue."
	}

	$IssueId = [int64]$listResponse.Json.data.list[0].id
	Write-Host "Using issue id $IssueId from GET /issues." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Step 6: GET /issues/$IssueId" -ForegroundColor Yellow
$null = Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/issues/$IssueId" `
	-Headers $authHeaders

Write-Host ""
Write-Host "Step 7: PATCH /issues/$IssueId" -ForegroundColor Yellow
$null = Invoke-MiniIssueApi `
	-Method "PATCH" `
	-Url "$base/issues/$IssueId" `
	-Headers $authHeaders `
	-Body @{
		title = "Updated issue from CLI"
		content = "Updated by examples/cli/powershell_issue_flow.ps1"
		status = "IN_PROGRESS"
		priority = 2
	}

Write-Host ""
Write-Host "Step 8: DELETE /issues/$IssueId" -ForegroundColor Yellow
$null = Invoke-MiniIssueApi `
	-Method "DELETE" `
	-Url "$base/issues/$IssueId" `
	-Headers $authHeaders
