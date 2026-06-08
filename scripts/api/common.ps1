$ErrorActionPreference = "Stop"

function Resolve-ApiBaseUrl {
	param(
		[string]$BaseUrl
	)

	if ([string]::IsNullOrWhiteSpace($BaseUrl)) {
		return "http://localhost:8080"
	}

	return $BaseUrl.TrimEnd("/")
}

function ConvertTo-PrettyJson {
	param(
		[AllowNull()]
		[object]$Value
	)

	if ($null -eq $Value) {
		return "null"
	}

	try {
		if ($Value -is [string]) {
			return ($Value | ConvertFrom-Json | ConvertTo-Json -Depth 20)
		}

		return ($Value | ConvertTo-Json -Depth 20)
	} catch {
		return [string]$Value
	}
}

function Read-ResponseBody {
	param(
		[object]$Response
	)

	if ($null -eq $Response) {
		return ""
	}

	try {
		$stream = $Response.GetResponseStream()
		if ($null -eq $stream) {
			return ""
		}

		$reader = New-Object System.IO.StreamReader($stream)
		return $reader.ReadToEnd()
	} catch {
		return ""
	}
}

function Invoke-MiniIssueApi {
	param(
		[Parameter(Mandatory = $true)]
		[string]$Method,

		[Parameter(Mandatory = $true)]
		[string]$Url,

		[hashtable]$Headers = @{},

		[AllowNull()]
		[object]$Body = $null
	)

	Write-Host ""
	Write-Host ">>> REQUEST" -ForegroundColor Cyan
	Write-Host "$Method $Url"

	if ($Headers.Count -gt 0) {
		Write-Host "Headers:"
		foreach ($key in $Headers.Keys) {
			Write-Host "  ${key}: $($Headers[$key])"
		}
	}

	$params = @{
		Method      = $Method
		Uri         = $Url
		Headers     = $Headers
		ContentType = "application/json"
		ErrorAction = "Stop"
	}

	if ($null -ne $Body) {
		$jsonBody = $Body | ConvertTo-Json -Depth 20
		$params.Body = $jsonBody
		Write-Host "Body:"
		Write-Host (ConvertTo-PrettyJson $jsonBody)
	}

	$statusCode = $null
	$content = ""
	$json = $null

	try {
		$response = Invoke-WebRequest @params
		$statusCode = [int]$response.StatusCode
		$content = [string]$response.Content
	} catch {
		$response = $_.Exception.Response
		if ($null -ne $response) {
			$statusCode = [int]$response.StatusCode
			$content = Read-ResponseBody $response
		} else {
			$statusCode = 0
			$content = $_.Exception.Message
		}
	}

	try {
		if (-not [string]::IsNullOrWhiteSpace($content)) {
			$json = $content | ConvertFrom-Json
		}
	} catch {
		$json = $null
	}

	Write-Host ""
	Write-Host "<<< RESPONSE $statusCode" -ForegroundColor Green
	if ([string]::IsNullOrWhiteSpace($content)) {
		Write-Host "<empty>"
	} else {
		Write-Host (ConvertTo-PrettyJson $content)
	}

	return [pscustomobject]@{
		StatusCode = $statusCode
		Content    = $content
		Json       = $json
	}
}
