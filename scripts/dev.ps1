# Codencil dev helpers — Docker Compose only (see agents/BUILD_ORDER.md P0.2).
param(
    [Parameter(Position = 0)]
    [ValidateSet("build", "up", "down", "ps", "api-version", "postgres", "migrate-up", "migrate-down", "migrate-reset")]
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
    "migrate-up"  {
        docker compose up -d postgres
        $db = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { "postgres://codencil:codencil@postgres:5432/codencil?sslmode=disable" }
        docker compose run --rm migrate -path /migrations -database $db up
    }
    "migrate-down" {
        $db = if ($env:DATABASE_URL) { $env:DATABASE_URL } else { "postgres://codencil:codencil@postgres:5432/codencil?sslmode=disable" }
        docker compose run --rm migrate -path /migrations -database $db down 1
    }
    "migrate-reset" {
        & $PSCommandPath migrate-down
        & $PSCommandPath migrate-up
    }
}
