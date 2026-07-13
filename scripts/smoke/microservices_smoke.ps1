$ErrorActionPreference = "Stop"

function Assert-StatusOk($Response, $Name) {
  if ($Response.StatusCode -lt 200 -or $Response.StatusCode -ge 300) {
    throw "$Name failed with status $($Response.StatusCode)"
  }
}

function Invoke-JsonPost($Url, $Body, $Session) {
  $json = $Body | ConvertTo-Json
  return Invoke-WebRequest $Url -Method Post -Body $json -ContentType "application/json" -WebSession $Session
}

function Assert-Unauthorized($Url, $Name) {
  try {
    $response = Invoke-WebRequest $Url
    throw "$Name expected 401, got $($response.StatusCode)"
  } catch {
    $response = $_.Exception.Response
    if ($null -eq $response) {
      throw
    }

    $statusCode = [int]$response.StatusCode
    if ($statusCode -ne 401) {
      throw "$Name expected 401, got $statusCode"
    }
  }
}

Write-Host "Checking gateway health..."
$gateway = Invoke-WebRequest "http://127.0.0.1:8888/healthz"
Assert-StatusOk $gateway "gateway health"

Write-Host "Checking auth anonymous /me..."
Assert-Unauthorized "http://127.0.0.1:9001/me" "auth anonymous /me"

Write-Host "Checking content article list..."
$contentList = Invoke-WebRequest "http://127.0.0.1:9003/articles?page=1&page_size=1"
Assert-StatusOk $contentList "content article list"

Write-Host "Checking login cookie flow..."
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$body = @{ username = "admin"; password = "admin" }
$login = Invoke-JsonPost "http://127.0.0.1:9001/login" $body $session
Assert-StatusOk $login "auth login"
$me = Invoke-WebRequest "http://127.0.0.1:9001/me" -WebSession $session
Assert-StatusOk $me "auth /me"

Write-Host "Microservice smoke verification passed."