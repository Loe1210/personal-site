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

function Get-ErrorStatusCode($ErrorRecord) {
  $response = $ErrorRecord.Exception.Response
  if ($null -eq $response) {
    throw $ErrorRecord
  }
  return [int]$response.StatusCode
}

function Assert-AuthAnonymous($Url, $Name) {
  try {
    $response = Invoke-WebRequest $Url
    Assert-StatusOk $response $Name
    $json = $response.Content | ConvertFrom-Json
    if ($json.code -eq 0) {
      throw "$Name expected auth error envelope, got success"
    }
    if ([string]::IsNullOrWhiteSpace($json.msg)) {
      throw "$Name expected auth error msg"
    }
  } catch {
    $statusCode = Get-ErrorStatusCode $_
    if ($statusCode -ne 401) {
      throw "$Name expected auth envelope or legacy 401, got $statusCode"
    }
  }
}

function Assert-NotFound($Url, $Name) {
  try {
    $response = Invoke-WebRequest $Url
    throw "$Name expected 404, got $($response.StatusCode)"
  } catch {
    $statusCode = Get-ErrorStatusCode $_
    if ($statusCode -ne 404) {
      throw "$Name expected 404, got $statusCode"
    }
  }
}

Write-Host "Checking gateway health..."
$gateway = Invoke-WebRequest "http://127.0.0.1:8888/healthz"
Assert-StatusOk $gateway "gateway health"

Write-Host "Checking auth anonymous /me..."
Assert-AuthAnonymous "http://127.0.0.1:9001/me" "auth anonymous /me"

Write-Host "Checking content article list..."
$contentList = Invoke-WebRequest "http://127.0.0.1:9003/articles?page=1&page_size=1"
Assert-StatusOk $contentList "content article list"

Write-Host "Checking gateway content article list..."
$gatewayContentList = Invoke-WebRequest "http://127.0.0.1:8888/api/content/articles?page=1&page_size=1"
Assert-StatusOk $gatewayContentList "gateway content article list"

Write-Host "Checking deprecated gateway /api/articles 404..."
Assert-NotFound "http://127.0.0.1:8888/api/articles?page=1&page_size=1" "deprecated /api/articles"

Write-Host "Checking login cookie flow..."
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$body = @{ username = "admin"; password = "admin" }
$login = Invoke-JsonPost "http://127.0.0.1:9001/login" $body $session
Assert-StatusOk $login "auth login"
$me = Invoke-WebRequest "http://127.0.0.1:9001/me" -WebSession $session
Assert-StatusOk $me "auth /me"

Write-Host "Microservice smoke verification passed."
