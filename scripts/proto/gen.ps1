param(
    [string]$Module = "github.com/Loe1210/personal-site"
)

$ErrorActionPreference = "Stop"

$kitex = Get-Command kitex -ErrorAction SilentlyContinue
if (-not $kitex) {
    Write-Error "kitex is not installed or not in PATH. Install Kitex tooling first, then rerun make proto-gen."
}

$protoFiles = @(
    "idl/auth/auth.proto",
    "idl/content/content.proto",
    "idl/media/media.proto"
)

foreach ($proto in $protoFiles) {
    if (-not (Test-Path $proto)) {
        Write-Error "missing proto file: $proto"
    }
    & $kitex.Source -module $Module -type protobuf -I idl $proto
}
