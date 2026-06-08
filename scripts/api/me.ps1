param(
	[string]$BaseUrl = "",
	[string]$Token = "",
	[string]$TokenFile = ".mini-issue-token.json"
)

. "$PSScriptRoot/common.ps1"

if ([string]::IsNullOrWhiteSpace($Token) -and (Test-Path $TokenFile)) {
	$tokenState = Get-Content -Path $TokenFile -Raw -Encoding UTF8 | ConvertFrom-Json
	$Token = $tokenState.token

	if ([string]::IsNullOrWhiteSpace($BaseUrl)) {
		$BaseUrl = $tokenState.base_url
	}
}

if ([string]::IsNullOrWhiteSpace($Token)) {
	throw "No token provided. Run scripts/api/login.ps1 first, or pass -Token."
}

$base = Resolve-ApiBaseUrl $BaseUrl

Invoke-MiniIssueApi `
	-Method "GET" `
	-Url "$base/me" `
	-Headers @{
		Authorization = "Bearer $Token"
	}

