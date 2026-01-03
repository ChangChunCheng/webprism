#!/bin/bash
# WEBPRISM Platform Detection Script
# Detects the current operating system and sets appropriate variables

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "macos"
            ;;
        Linux*)
            if [ -f /etc/os-release ]; then
                . /etc/os-release
                case "$ID" in
                    ubuntu|debian)
                        echo "ubuntu"
                        ;;
                    centos|rhel|fedora)
                        echo "linux"
                        ;;
                    *)
                        echo "linux"
                        ;;
                esac
            else
                echo "linux"
            fi
            ;;
        CYGWIN*|MINGW*|MSYS*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# Check if Docker is available
check_docker() {
    if command -v docker &> /dev/null && docker ps &> /dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Check if PostgreSQL is installed locally
check_postgres() {
    if command -v psql &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# Main function to display platform info
main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    echo -e "${BLUE}🖥️  Platform Detection${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "Operating System: ${GREEN}$OS${NC}"
    echo -e "Architecture:     ${GREEN}$ARCH${NC}"

    if check_docker; then
        echo -e "Docker:           ${GREEN}✓ Available${NC}"
    else
        echo -e "Docker:           ${YELLOW}✗ Not available${NC}"
    fi

    if check_postgres; then
        echo -e "PostgreSQL:       ${GREEN}✓ Installed${NC}"
    else
        echo -e "PostgreSQL:       ${YELLOW}✗ Not installed${NC}"
    fi

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

    # Export for other scripts
    export WEBPRISM_OS="$OS"
    export WEBPRISM_ARCH="$ARCH"
}

# If sourced, just set variables. If executed, display info
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main
fi
