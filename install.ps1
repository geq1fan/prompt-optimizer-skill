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
$RepoUrl = "https://github.com/geq1fan/prompt-optimizer-skill"

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

function Test-GitInstalled {
    try {
        $null = git --version
        return $true
    }
    catch {
        return $false
    }
}

function Ensure-SkillsDirectory {
    if (-not (Test-Path $ClaudeSkillsDir)) {
        Write-Step "Creating skills directory..."
        New-Item -ItemType Directory -Path $ClaudeSkillsDir -Force | Out-Null
    }
}

function Remove-InstallationFiles {
    param([string]$TargetDir)

    $filesToRemove = @(
        ".git",
        "install.sh",
        "install.ps1",
        "task_plan.md",
        "findings.md",
        "progress.md"
    )

    foreach ($file in $filesToRemove) {
        $path = Join-Path $TargetDir $file
        if (Test-Path $path) {
            Remove-Item -Path $path -Recurse -Force
        }
    }
}

function Install-Skill {
    $TargetDir = Join-Path $ClaudeSkillsDir $SkillName

    if (Test-Path $TargetDir) {
        Write-Warning "Skill already installed at $TargetDir"
        Write-Host "Use -Action update to update the existing installation."
        exit 1
    }

    Write-Step "Installing $SkillName..."
    git clone $RepoUrl $TargetDir

    Write-Step "Cleaning up..."
    Remove-InstallationFiles -TargetDir $TargetDir

    Write-Host ""
    Write-Step "Installation complete!"
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

    Write-Step "Updating $SkillName..."

    # Create temp directory
    $TempDir = Join-Path $env:TEMP "prompt-optimizer-update-$(Get-Date -Format 'yyyyMMddHHmmss')"

    # Clone to temp directory
    git clone $RepoUrl $TempDir

    # Backup existing installation
    $BackupDir = "$TargetDir.backup.$(Get-Date -Format 'yyyyMMddHHmmss')"
    Move-Item -Path $TargetDir -Destination $BackupDir

    # Move new version
    Move-Item -Path $TempDir -Destination $TargetDir

    # Cleanup
    Remove-InstallationFiles -TargetDir $TargetDir

    Write-Host ""
    Write-Step "Update complete!"
    Write-Host "Backup saved to: $BackupDir"
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
if (-not (Test-GitInstalled)) {
    Write-Error "git is not installed. Please install git first."
    exit 1
}

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
