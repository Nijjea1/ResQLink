<#
Simple helper script to launch two MeshComm nodes (alice and bob), send a test message
from alice and check bob's /api/messages for receipt.

Usage (PowerShell):
  cd <repo root>\MeshComm
  powershell -ExecutionPolicy Bypass -File .\scripts\run_two_node_test.ps1

This script starts two new PowerShell processes that run the Go servers and redirect
their stdout/stderr to logs in the same folder (alice.log, bob.log). It then sends a
POST to alice's API and GETs bob's messages.
#>

$repoRoot = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$cmdDir = Join-Path $repoRoot 'cmd\meshcomm'

$aliceLog = Join-Path $repoRoot 'alice.log'
$bobLog = Join-Path $repoRoot 'bob.log'

Write-Host "Starting alice (P2P:9001, API:3002) and bob (P2P:9002, API:3003)..."

# Start alice
Start-Process -FilePath pwsh -ArgumentList @('-NoProfile','-Command',"cd '$cmdDir'; go run . -port 9001 -api-port 3002 -nick alice -same_string meshcomm 2>&1 | Out-File -FilePath '$aliceLog' -Encoding utf8 -Append") -WindowStyle Hidden

# Start bob
Start-Process -FilePath pwsh -ArgumentList @('-NoProfile','-Command',"cd '$cmdDir'; go run . -port 9002 -api-port 3003 -nick bob -same_string meshcomm 2>&1 | Out-File -FilePath '$bobLog' -Encoding utf8 -Append") -WindowStyle Hidden

Write-Host "Waiting 4 seconds for nodes to start..."
Start-Sleep -Seconds 4

$postUrl = 'http://localhost:3002/api/send'
$body = '{"content":"hello from alice","category":"GENERAL"}'
try {
    Write-Host "Sending test POST to alice..."
    Invoke-RestMethod -Method POST -ContentType 'application/json' -Body $body $postUrl
    Write-Host "POST succeeded"
} catch {
    Write-Host "POST failed: $_"
}

Start-Sleep -Seconds 2

$getUrl = 'http://localhost:3003/api/messages'
try {
    Write-Host "Fetching messages from bob..."
    $res = Invoke-RestMethod $getUrl
    Write-Host "bob messages:`n"
    $res | ConvertTo-Json -Depth 5
} catch {
    Write-Host "GET failed: $_"
}

Write-Host "--- Tail of alice.log ---"
if (Test-Path $aliceLog) { Get-Content $aliceLog -Tail 40 -Wait:$false } else { Write-Host "alice.log not found" }

Write-Host "--- Tail of bob.log ---"
if (Test-Path $bobLog) { Get-Content $bobLog -Tail 40 -Wait:$false } else { Write-Host "bob.log not found" }

Write-Host "Test script finished. To stop the servers, close their PowerShell windows or kill the 'go' processes."
