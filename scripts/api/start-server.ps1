$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path "$PSScriptRoot/../.."

Write-Host "Starting mini-issue server at :8080" -ForegroundColor Yellow
Write-Host "The address is currently configured in cmd/main.go."
Write-Host "Run API scripts from another terminal while this command is running."

Push-Location $repoRoot
try {
	go run ./cmd
} finally {
	Pop-Location
}
