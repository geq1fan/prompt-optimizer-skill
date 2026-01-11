#!/bin/bash
#
# Prompt Optimizer Skill - Installation Script for macOS/Linux
# https://github.com/geq1fan/prompt-optimizer-skill
#

set -e

SKILL_NAME="prompt-optimizer-skill"
CLAUDE_SKILLS_DIR="$HOME/.claude/skills"
REPO_URL="https://github.com/geq1fan/prompt-optimizer-skill"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${GREEN}==>${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

print_error() {
    echo -e "${RED}Error:${NC} $1"
}

# Check if git is installed
check_git() {
    if ! command -v git &> /dev/null; then
        print_error "git is not installed. Please install git first."
        exit 1
    fi
}

# Create skills directory if not exists
ensure_skills_dir() {
    if [ ! -d "$CLAUDE_SKILLS_DIR" ]; then
        print_step "Creating skills directory..."
        mkdir -p "$CLAUDE_SKILLS_DIR"
    fi
}

# Install skill
install_skill() {
    local target_dir="$CLAUDE_SKILLS_DIR/$SKILL_NAME"

    if [ -d "$target_dir" ]; then
        print_warning "Skill already installed at $target_dir"
        echo "Use --update to update the existing installation."
        exit 1
    fi

    print_step "Installing $SKILL_NAME..."
    git clone "$REPO_URL" "$target_dir"

    print_step "Cleaning up..."
    rm -rf "$target_dir/.git"
    rm -f "$target_dir/install.sh"
    rm -f "$target_dir/install.ps1"
    rm -rf "$target_dir/task_plan.md"
    rm -rf "$target_dir/findings.md"
    rm -rf "$target_dir/progress.md"

    echo ""
    print_step "Installation complete!"
    echo ""
    echo "Usage:"
    echo "  /optimize-prompt <your prompt>"
    echo "  /optimize-prompt iterate <improvement instructions>"
    echo ""
    echo "Restart Claude Code to load the new skill."
}

# Update skill
update_skill() {
    local target_dir="$CLAUDE_SKILLS_DIR/$SKILL_NAME"
    local temp_dir=$(mktemp -d)

    if [ ! -d "$target_dir" ]; then
        print_error "Skill not installed. Use install mode first."
        exit 1
    fi

    print_step "Updating $SKILL_NAME..."

    # Clone to temp directory
    git clone "$REPO_URL" "$temp_dir/$SKILL_NAME"

    # Backup existing installation
    local backup_dir="$target_dir.backup.$(date +%Y%m%d%H%M%S)"
    mv "$target_dir" "$backup_dir"

    # Move new version
    mv "$temp_dir/$SKILL_NAME" "$target_dir"

    # Cleanup
    rm -rf "$target_dir/.git"
    rm -f "$target_dir/install.sh"
    rm -f "$target_dir/install.ps1"
    rm -rf "$target_dir/task_plan.md"
    rm -rf "$target_dir/findings.md"
    rm -rf "$target_dir/progress.md"
    rm -rf "$temp_dir"

    echo ""
    print_step "Update complete!"
    echo "Backup saved to: $backup_dir"
    echo ""
    echo "Restart Claude Code to load the updated skill."
}

# Uninstall skill
uninstall_skill() {
    local target_dir="$CLAUDE_SKILLS_DIR/$SKILL_NAME"

    if [ ! -d "$target_dir" ]; then
        print_error "Skill not installed."
        exit 1
    fi

    print_step "Uninstalling $SKILL_NAME..."
    rm -rf "$target_dir"

    print_step "Uninstall complete!"
}

# Show help
show_help() {
    echo "Prompt Optimizer Skill - Installation Script"
    echo ""
    echo "Usage: ./install.sh [command]"
    echo ""
    echo "Commands:"
    echo "  install     Install the skill (default)"
    echo "  update      Update to latest version"
    echo "  uninstall   Remove the skill"
    echo "  help        Show this help message"
    echo ""
}

# Main
main() {
    local command="${1:-install}"

    case "$command" in
        install)
            check_git
            ensure_skills_dir
            install_skill
            ;;
        update|--update|-u)
            check_git
            update_skill
            ;;
        uninstall|--uninstall|remove)
            uninstall_skill
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
