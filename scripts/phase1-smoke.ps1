# Phase 1 end-to-end smoke test (Docker stack must be running).
param(
    [string]$ApiUrl = "http://localhost:8080"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

Write-Host "Phase 1 smoke against $ApiUrl"

$health = Invoke-RestMethod -Uri "$ApiUrl/health" -Method Get
if ($health.status -ne "ok") { throw "health check failed" }

$created = Invoke-RestMethod -Uri "$ApiUrl/documents" -Method Post -ContentType "application/json" -Body (@{
    title = "Smoke test"
    draft_markdown = "# Hello Codencil"
} | ConvertTo-Json)

$id = $created.id
Write-Host "Created document $id"

Invoke-RestMethod -Uri "$ApiUrl/documents/$id" -Method Patch -ContentType "application/json" -Body (@{
    draft_markdown = "# Hello Codencil`n`nPublished snapshot."
} | ConvertTo-Json) | Out-Null

$v1 = Invoke-RestMethod -Uri "$ApiUrl/documents/$id/publish" -Method Post -ContentType "application/json" -Body (@{
    published_by = "smoke-script"
} | ConvertTo-Json)

if ($v1.version -ne 1) { throw "expected version 1" }

$snapshot = Invoke-RestMethod -Uri "$ApiUrl/documents/$id/versions/1" -Method Get
if ($snapshot.markdown -notmatch "Published snapshot") { throw "snapshot markdown mismatch" }

Write-Host "OK: published v1 for document $id"
Write-Host "Preview: http://localhost:3000/documents/$id/versions/1"
