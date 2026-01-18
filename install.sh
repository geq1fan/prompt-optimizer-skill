#!/bin/bash
#
# Prompt Optimizer Skill - Installation Script for macOS/Linux
# https://github.com/geq1fan/prompt-optimizer-skill
#

set -e

SKILL_NAME="prompt-optimizer-skill"
CLAUDE_SKILLS_DIR="$HOME/.claude/skills"
REPO="geq1fan/prompt-optimizer-skill"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"
GITHUB_RELEASES="https://github.com/$REPO/releases/download"

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

# Detect platform
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        linux*)
            PLATFORM="linux"
            ;;
        darwin*)
            PLATFORM="darwin"
            ;;
        *)
            print_error "Unsupported OS: $os"
            exit 1
            ;;
    esac

    # For now, we use universal binary for macOS and amd64 for Linux
    if [ "$PLATFORM" = "darwin" ]; then
        ARCH="universal"
        EXT="tar.gz"
    else
        ARCH="amd64"
        EXT="tar.gz"
    fi

    PACKAGE_NAME="prompt-optimizer-skill-${PLATFORM}-${ARCH}.${EXT}"
}

# Get latest version from GitHub
get_latest_version() {
    print_step "Fetching latest version..."

    if command -v curl &> /dev/null; then
        VERSION=$(curl -s "$GITHUB_API" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        VERSION=$(wget -qO- "$GITHUB_API" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
    else
        print_error "curl or wget is required"
        exit 1
    fi

    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version. Check your network connection."
        exit 1
    fi

    echo "Latest version: $VERSION"
}

# Download and extract package
download_package() {
    local url="$GITHUB_RELEASES/$VERSION/$PACKAGE_NAME"
    local temp_dir=$(mktemp -d)
    local package_path="$temp_dir/$PACKAGE_NAME"

    print_step "Downloading $PACKAGE_NAME..."

    if command -v curl &> /dev/null; then
        curl -L -o "$package_path" "$url"
    else
        wget -O "$package_path" "$url"
    fi

    if [ ! -f "$package_path" ]; then
        print_error "Download failed"
        exit 1
    fi

    print_step "Extracting..."

    tar -xzf "$package_path" -C "$temp_dir"

    TEMP_SKILLS_DIR="$temp_dir/skills"
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
        echo "Use 'update' to update the existing installation."
        exit 1
    fi

    detect_platform
    get_latest_version
    download_package

    print_step "Installing $SKILL_NAME..."
    mv "$TEMP_SKILLS_DIR" "$target_dir"

    echo ""
    print_step "Installation complete! (version: $VERSION)"
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

    if [ ! -d "$target_dir" ]; then
        print_error "Skill not installed. Use install mode first."
        exit 1
    fi

    detect_platform
    get_latest_version
    download_package

    print_step "Updating $SKILL_NAME..."

    # Remove existing installation
    rm -rf "$target_dir"

    # Move new version
    mv "$TEMP_SKILLS_DIR" "$target_dir"

    echo ""
    print_step "Update complete! (version: $VERSION)"
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
            ensure_skills_dir
            install_skill
            ;;
        update|--update|-u)
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
