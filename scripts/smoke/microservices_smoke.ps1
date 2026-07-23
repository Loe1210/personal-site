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
$gatewayContentJson = $gatewayContentList.Content | ConvertFrom-Json
if ($gatewayContentJson.code -ne 0) {
  throw "gateway content article list expected success envelope, got code $($gatewayContentJson.code)"
}
if ([string]::IsNullOrWhiteSpace($gatewayContentJson.msg)) {
  throw "gateway content article list expected envelope msg"
}

Write-Host "Checking deprecated gateway /api/articles 404..."
Assert-NotFound "http://127.0.0.1:8888/api/articles?page=1&page_size=1" "deprecated /api/articles"

Write-Host "Checking media upload through gateway..."
Add-Type -AssemblyName System.Net.Http
$pngBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/p9sAAAAASUVORK5CYII="
$pngBytes = [Convert]::FromBase64String($pngBase64)
$httpClient = New-Object System.Net.Http.HttpClient
$multipart = New-Object System.Net.Http.MultipartFormDataContent
$fileContent = New-Object System.Net.Http.ByteArrayContent -ArgumentList @(,$pngBytes)
$fileContent.Headers.ContentType = [System.Net.Http.Headers.MediaTypeHeaderValue]::Parse("image/png")
$multipart.Add($fileContent, "file", "micro-smoke.png")
try {
  $mediaResponse = $httpClient.PostAsync("http://127.0.0.1:8888/api/media/upload?user_id=1&biz_type=smoke&biz_id=micro-smoke", $multipart).GetAwaiter().GetResult()
  if ([int]$mediaResponse.StatusCode -lt 200 -or [int]$mediaResponse.StatusCode -ge 300) {
    throw "gateway media upload failed with status $([int]$mediaResponse.StatusCode)"
  }
  $mediaContent = $mediaResponse.Content.ReadAsStringAsync().GetAwaiter().GetResult()
  $mediaJson = $mediaContent | ConvertFrom-Json
  if ($mediaJson.code -ne 0) {
    throw "gateway media upload expected success code, got $($mediaJson.code)"
  }
  if ([string]::IsNullOrWhiteSpace($mediaJson.data.url)) {
    throw "gateway media upload expected media url"
  }
} finally {
  if ($null -ne $mediaResponse) { $mediaResponse.Dispose() }
  $multipart.Dispose()
  $httpClient.Dispose()
}

Write-Host "Checking login cookie flow..."
$session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$body = @{ username = "admin"; password = "admin" }
$login = Invoke-JsonPost "http://127.0.0.1:9001/login" $body $session
Assert-StatusOk $login "auth login"
$me = Invoke-WebRequest "http://127.0.0.1:9001/me" -WebSession $session
Assert-StatusOk $me "auth /me"

Write-Host "Microservice smoke verification passed."

