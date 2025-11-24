#!/bin/bash

# University of Glasgow Inquiry Management Script
# 格拉斯哥大學諮詢管理腳本 - 統一管理所有相關功能

# 顏色輸出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 腳本路徑
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GET_TOKEN_SCRIPT="$SCRIPT_DIR/get_recaptcha_token.sh"
RUN_GETTER_SCRIPT="$SCRIPT_DIR/run_recaptcha_getter.sh"
SUBMIT_SCRIPT="$SCRIPT_DIR/submit-inquiry.sh"

# 顯示橫幅
show_banner() {
    echo -e "${GREEN}====================================================${NC}"
    echo -e "${GREEN}  University of Glasgow Inquiry Management System${NC}"
    echo -e "${GREEN}====================================================${NC}"
    echo ""
}

# 顯示幫助信息
show_help() {
    show_banner
    echo -e "${YELLOW}Usage: $0 <command> [options]${NC}"
    echo ""
    echo -e "${CYAN}Commands:${NC}"
    echo -e "  ${GREEN}help${NC}                    Show this help message"
    echo -e "  ${GREEN}manual${NC}                  Manual reCAPTCHA token retrieval (with instructions)"
    echo -e "  ${GREEN}auto${NC}                    Automated reCAPTCHA token retrieval (Selenium)"
    echo -e "  ${GREEN}submit${NC}                  Start automated inquiry submission (every 4 hours)"
    echo -e "  ${GREEN}status${NC}                  Check current system status"
    echo -e "  ${GREEN}logs${NC}                    Show recent log entries"
    echo ""
    echo -e "${CYAN}Auto Command Options:${NC}"
    echo -e "  --firefox                 Use Firefox instead of Chrome"
    echo -e "  --headless                Run in headless mode"
    ""
    echo -e "${CYAN}Examples:${NC}"
    echo -e "  $0 manual                 # Get manual instructions for reCAPTCHA token"
    echo -e "  $0 auto                   # Run automated token getter with Chrome"
    echo -e "  $0 auto --firefox         # Run automated token getter with Firefox"
    echo -e "  $0 submit                 # Start automated submission every 4 hours"
    echo -e "  $0 status                 # Check if token is set and system is ready"
    echo -e "  $0 logs                   # Show last 10 log entries"
    echo ""
}

# 檢查腳本是否存在
check_scripts() {
    local missing_scripts=()

    if [ ! -f "$GET_TOKEN_SCRIPT" ]; then
        missing_scripts+=("get_recaptcha_token.sh")
    fi

    if [ ! -f "$RUN_GETTER_SCRIPT" ]; then
        missing_scripts+=("run_recaptcha_getter.sh")
    fi

    if [ ! -f "$SUBMIT_SCRIPT" ]; then
        missing_scripts+=("submit-inquiry.sh")
    fi

    if [ ${#missing_scripts[@]} -ne 0 ]; then
        echo -e "${RED}✗ Missing required scripts:${NC}"
        for script in "${missing_scripts[@]}"; do
            echo -e "${RED}  - $script${NC}"
        done
        echo ""
        echo -e "${YELLOW}Please ensure all scripts are in the same directory.${NC}"
        exit 1
    fi
}

# 手動 reCAPTCHA token 獲取
run_manual_mode() {
    echo -e "${BLUE}Starting manual reCAPTCHA token retrieval...${NC}"
    echo ""

    if [ -f "$GET_TOKEN_SCRIPT" ]; then
        bash "$GET_TOKEN_SCRIPT"
    else
        echo -e "${RED}✗ Manual token getter script not found: $GET_TOKEN_SCRIPT${NC}"
        exit 1
    fi
}

# 自動化 reCAPTCHA token 獲取
run_auto_mode() {
    echo -e "${BLUE}Starting automated reCAPTCHA token retrieval...${NC}"
    echo ""

    if [ -f "$RUN_GETTER_SCRIPT" ]; then
        # 傳遞所有額外參數給自動腳本
        bash "$RUN_GETTER_SCRIPT" "$@"
    else
        echo -e "${RED}✗ Automated token getter script not found: $RUN_GETTER_SCRIPT${NC}"
        exit 1
    fi
}

# 自動提交模式
run_submit_mode() {
    echo -e "${BLUE}Starting automated inquiry submission...${NC}"
    echo ""

    # 檢查 reCAPTCHA token 是否設置
    if [ -z "$RECAPTCHA_TOKEN" ]; then
        echo -e "${YELLOW}⚠ Warning: RECAPTCHA_TOKEN environment variable is not set${NC}"

        # 檢查是否有 token 文件
        if [ -f "recaptcha_token.txt" ]; then
            echo -e "${YELLOW}Found token file: recaptcha_token.txt${NC}"
            echo -e "${YELLOW}Loading token from file...${NC}"
            export RECAPTCHA_TOKEN="$(cat recaptcha_token.txt)"
            echo -e "${GREEN}✓ Token loaded from file${NC}"
        else
            echo -e "${RED}✗ No RECAPTCHA_TOKEN set and no token file found${NC}"
            echo ""
            echo -e "${YELLOW}Please get a reCAPTCHA token first:${NC}"
            echo -e "  $0 manual    # Manual instructions"
            echo -e "  $0 auto      # Automated retrieval"
            echo ""
            exit 1
        fi
    else
        echo -e "${GREEN}✓ RECAPTCHA_TOKEN is set${NC}"
    fi

    echo ""

    if [ -f "$SUBMIT_SCRIPT" ]; then
        bash "$SUBMIT_SCRIPT"
    else
        echo -e "${RED}✗ Submit script not found: $SUBMIT_SCRIPT${NC}"
        exit 1
    fi
}

# 檢查系統狀態
check_status() {
    echo -e "${BLUE}System Status Check${NC}"
    echo -e "${BLUE}==================${NC}"
    echo ""

    # 檢查腳本
    echo -e "${CYAN}Scripts:${NC}"
    if [ -f "$GET_TOKEN_SCRIPT" ]; then
        echo -e "${GREEN}✓${NC} get_recaptcha_token.sh (manual mode)"
    else
        echo -e "${RED}✗${NC} get_recaptcha_token.sh (manual mode)"
    fi

    if [ -f "$RUN_GETTER_SCRIPT" ]; then
        echo -e "${GREEN}✓${NC} run_recaptcha_getter.sh (auto mode)"
    else
        echo -e "${RED}✗${NC} run_recaptcha_getter.sh (auto mode)"
    fi

    if [ -f "$SUBMIT_SCRIPT" ]; then
        echo -e "${GREEN}✓${NC} submit-inquiry.sh (submit mode)"
    else
        echo -e "${RED}✗${NC} submit-inquiry.sh (submit mode)"
    fi

    echo ""

    # 檢查環境變數
    echo -e "${CYAN}Environment:${NC}"
    if [ -n "$RECAPTCHA_TOKEN" ]; then
        echo -e "${GREEN}✓${NC} RECAPTCHA_TOKEN is set (${#RECAPTCHA_TOKEN} characters)"
    else
        echo -e "${RED}✗${NC} RECAPTCHA_TOKEN is not set"
    fi

    if [ -n "$GLA_REQUEST_LOGS" ]; then
        echo -e "${GREEN}✓${NC} GLA_REQUEST_LOGS is set to: $GLA_REQUEST_LOGS"
    else
        echo -e "${YELLOW}⚠${NC} GLA_REQUEST_LOGS not set, using default: GLA_REQUEST_LOGS"
    fi

    echo ""

    # 檢查文件
    echo -e "${CYAN}Files:${NC}"
    if [ -f "recaptcha_token.txt" ]; then
        local token_length=$(wc -c < "recaptcha_token.txt" 2>/dev/null || echo "0")
        echo -e "${GREEN}✓${NC} recaptcha_token.txt exists (${token_length} characters)"
    else
        echo -e "${RED}✗${NC} recaptcha_token.txt not found"
    fi

    if [ -d "GLA_REQUEST_LOGS" ]; then
        local log_count=$(find GLA_REQUEST_LOGS -name "*.log" 2>/dev/null | wc -l)
        echo -e "${GREEN}✓${NC} GLA_REQUEST_LOGS directory exists (${log_count} log files)"
    else
        echo -e "${RED}✗${NC} GLA_REQUEST_LOGS directory not found"
    fi

    echo ""

    # 檢查依賴項
    echo -e "${CYAN}Dependencies:${NC}"
    if command -v python3 &> /dev/null || command -v python &> /dev/null; then
        local python_cmd="python3"
        if ! command -v python3 &> /dev/null; then
            python_cmd="python"
        fi
        local python_version=$($python_cmd --version 2>&1 | head -1)
        echo -e "${GREEN}✓${NC} Python available: $python_version"
    else
        echo -e "${RED}✗${NC} Python not found"
    fi

    if command -v curl &> /dev/null; then
        echo -e "${GREEN}✓${NC} curl available"
    else
        echo -e "${RED}✗${NC} curl not found"
    fi

    echo ""

    # 總結
    echo -e "${CYAN}Summary:${NC}"
    if [ -n "$RECAPTCHA_TOKEN" ] || [ -f "recaptcha_token.txt" ]; then
        echo -e "${GREEN}✓ System is ready for submission${NC}"
    else
        echo -e "${YELLOW}⚠ System needs reCAPTCHA token${NC}"
        echo -e "${YELLOW}  Run: $0 manual  or  $0 auto${NC}"
    fi
}

# 顯示日誌
show_logs() {
    echo -e "${BLUE}Recent Log Entries${NC}"
    echo -e "${BLUE}==================${NC}"
    echo ""

    local log_dir="${GLA_REQUEST_LOGS:-GLA_REQUEST_LOGS}"

    if [ ! -d "$log_dir" ]; then
        echo -e "${RED}✗ Log directory not found: $log_dir${NC}"
        echo -e "${YELLOW}No logs available yet. Run the submit command first.${NC}"
        return
    fi

    # 找到最新的日誌文件
    local latest_log=$(find "$log_dir" -name "*.log" -type f -printf '%T@ %p\n' 2>/dev/null | sort -n | tail -1 | cut -d' ' -f2-)

    if [ -z "$latest_log" ]; then
        echo -e "${YELLOW}No log files found in $log_dir${NC}"
        return
    fi

    echo -e "${CYAN}Latest log file: $latest_log${NC}"
    echo ""

    # 顯示最後10行
    if [ -f "$latest_log" ]; then
        echo -e "${CYAN}Last 10 entries:${NC}"
        echo ""
        tail -10 "$latest_log"
    else
        echo -e "${RED}✗ Cannot read log file: $latest_log${NC}"
    fi
}

# 主函數
main() {
    check_scripts

    case "${1:-help}" in
        "help"|"-h"|"--help")
            show_help
            ;;
        "manual"|"m")
            shift
            run_manual_mode "$@"
            ;;
        "auto"|"a")
            shift
            run_auto_mode "$@"
            ;;
        "submit"|"s")
            shift
            run_submit_mode "$@"
            ;;
        "status"|"st")
            shift
            check_status "$@"
            ;;
        "logs"|"l")
            shift
            show_logs "$@"
            ;;
        *)
            echo -e "${RED}✗ Unknown command: $1${NC}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# 執行主函數
main "$@"
