#Requires -Version 5.1
<#
.SYNOPSIS
    Prompt Optimizer Skill - Installation Script for Windows

.DESCRIPTION
    Install, update, or uninstall the prompt-optimizer skill for Claude Code.

.PARAMETER Action
    The action to perform: install, update, or uninstall

.EXAMPLE
    .\install.ps1
    .\install.ps1 -Action install
    .\install.ps1 -Action update
    .\install.ps1 -Action uninstall

.LINK
    https://github.com/geq1fan/prompt-optimizer-skill
#>

param(
    [Parameter(Position = 0)]
    [ValidateSet("install", "update", "uninstall", "help")]
    [string]$Action = "install"
)

$ErrorActionPreference = "Stop"

$SkillName = "prompt-optimizer-skill"
$ClaudeSkillsDir = Join-Path $env:USERPROFILE ".claude\skills"
$Repo = "geq1fan/prompt-optimizer-skill"
$GitHubApi = "https://api.github.com/repos/$Repo/releases/latest"
$GitHubReleases = "https://github.com/$Repo/releases/download"

function Write-Step {
    param([string]$Message)
    Write-Host "==> " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Warning {
    param([string]$Message)
    Write-Host "Warning: " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

function Write-Error {
    param([string]$Message)
    Write-Host "Error: " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

function Get-LatestVersion {
    Write-Step "Fetching latest version..."

    try {
        $response = Invoke-RestMethod -Uri $GitHubApi -UseBasicParsing
        $script:Version = $response.tag_name
        Write-Host "Latest version: $script:Version"
    }
    catch {
        Write-Error "Failed to get latest version. Check your network connection."
        exit 1
    }
}

function Get-PackageName {
    $script:PackageName = "prompt-optimizer-skill-windows-amd64.zip"
}

function Download-Package {
    $url = "$GitHubReleases/$script:Version/$script:PackageName"
    $tempDir = Join-Path $env:TEMP "prompt-optimizer-install-$(Get-Date -Format 'yyyyMMddHHmmss')"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    $packagePath = Join-Path $tempDir $script:PackageName

    Write-Step "Downloading $script:PackageName..."

    try {
        Invoke-WebRequest -Uri $url -OutFile $packagePath -UseBasicParsing
    }
    catch {
        Write-Error "Download failed: $_"
        exit 1
    }

    Write-Step "Extracting..."

    Expand-Archive -Path $packagePath -DestinationPath $tempDir -Force

    $script:TempSkillsDir = Join-Path $tempDir "skills"
}

function Ensure-SkillsDirectory {
    if (-not (Test-Path $ClaudeSkillsDir)) {
        Write-Step "Creating skills directory..."
        New-Item -ItemType Directory -Path $ClaudeSkillsDir -Force | Out-Null
    }
}

function Install-Skill {
    $TargetDir = Join-Path $ClaudeSkillsDir $SkillName

    if (Test-Path $TargetDir) {
        Write-Warning "Skill already installed at $TargetDir"
        Write-Host "Use -Action update to update the existing installation."
        exit 1
    }

    Get-PackageName
    Get-LatestVersion
    Download-Package

    Write-Step "Installing $SkillName..."
    Move-Item -Path $script:TempSkillsDir -Destination $TargetDir

    Write-Host ""
    Write-Step "Installation complete! (version: $script:Version)"
    Write-Host ""
    Write-Host "Usage:"
    Write-Host "  /optimize-prompt <your prompt>"
    Write-Host "  /optimize-prompt iterate <improvement instructions>"
    Write-Host ""
    Write-Host "Restart Claude Code to load the new skill."
}

function Update-Skill {
    $TargetDir = Join-Path $ClaudeSkillsDir $SkillName

    if (-not (Test-Path $TargetDir)) {
        Write-Error "Skill not installed. Use install action first."
        exit 1
    }

    Get-PackageName
    Get-LatestVersion
    Download-Package

    Write-Step "Updating $SkillName..."

    # Remove existing installation
    Remove-Item -Path $TargetDir -Recurse -Force

    # Move new version
    Move-Item -Path $script:TempSkillsDir -Destination $TargetDir

    Write-Host ""
    Write-Step "Update complete! (version: $script:Version)"
    Write-Host ""
    Write-Host "Restart Claude Code to load the updated skill."
}

function Uninstall-Skill {
    $TargetDir = Join-Path $ClaudeSkillsDir $SkillName

    if (-not (Test-Path $TargetDir)) {
        Write-Error "Skill not installed."
        exit 1
    }

    Write-Step "Uninstalling $SkillName..."
    Remove-Item -Path $TargetDir -Recurse -Force

    Write-Step "Uninstall complete!"
}

function Show-Help {
    Write-Host "Prompt Optimizer Skill - Installation Script"
    Write-Host ""
    Write-Host "Usage: .\install.ps1 [-Action <action>]"
    Write-Host ""
    Write-Host "Actions:"
    Write-Host "  install     Install the skill (default)"
    Write-Host "  update      Update to latest version"
    Write-Host "  uninstall   Remove the skill"
    Write-Host "  help        Show this help message"
    Write-Host ""
}

# Main
switch ($Action) {
    "install" {
        Ensure-SkillsDirectory
        Install-Skill
    }
    "update" {
        Update-Skill
    }
    "uninstall" {
        Uninstall-Skill
    }
    "help" {
        Show-Help
    }
}
