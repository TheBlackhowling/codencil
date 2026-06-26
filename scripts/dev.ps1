# Codencil dev helpers — Docker Compose only (see agents/BUILD_ORDER.md P0.2).
param(
    [Parameter(Position = 0)]
    [ValidateSet("build", "up", "down", "ps", "api-version", "postgres")]
    [string]$Command = "build"
)

Set-Location (Split-Path $PSScriptRoot -Parent)

switch ($Command) {
    "build"       { docker compose build }
    "up"          { docker compose up -d postgres }
    "down"        { docker compose down }
    "ps"          { docker compose ps }
    "api-version" { docker compose run --rm api go version }
    "postgres"    {
        docker compose up -d postgres
        docker compose ps postgres
    }
}
