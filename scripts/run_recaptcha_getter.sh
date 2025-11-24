#!/bin/bash

# Wrapper script to run the automated reCAPTCHA token getter
# 自動化 reCAPTCHA token 獲取器的包裝腳本

# 顏色輸出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}Automated reCAPTCHA Token Getter${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

# 檢查 Python 是否安裝（支持 Windows Git Bash，跳過 Windows Store 快捷方式）
PYTHON_CMD=""
PYTHON_VERSION=""

# 函數：測試 Python 命令是否真的可用（不是 Windows Store 快捷方式）
test_python() {
    local cmd="$1"
    local test_output
    test_output=$($cmd --version 2>&1)
    
    # 檢查是否包含 Windows Store 錯誤訊息
    if echo "$test_output" | grep -qi "Microsoft Store\|App execution aliases"; then
        return 1
    fi
    
    # 檢查是否是有效的 Python 版本輸出
    if echo "$test_output" | grep -qiE "Python [0-9]+\.[0-9]+"; then
        PYTHON_VERSION="$test_output"
        return 0
    fi
    
    return 1
}

# 優先查找實際的 Python 安裝路徑（跳過 Windows Store 快捷方式）
# 1. 查找常見的 Python 安裝路徑
if [ -d "/c/Users/$USER/AppData/Local/Programs/Python" ]; then
    PYTHON_CANDIDATE=$(find "/c/Users/$USER/AppData/Local/Programs/Python" -name "python.exe" 2>/dev/null | head -1)
    if [ -n "$PYTHON_CANDIDATE" ] && test_python "$PYTHON_CANDIDATE"; then
        PYTHON_CMD="$PYTHON_CANDIDATE"
    fi
fi

# 2. 查找 C:\Python* 路徑
if [ -z "$PYTHON_CMD" ] && [ -d "/c/Python" ]; then
    PYTHON_CANDIDATE=$(find /c/Python* -name "python.exe" 2>/dev/null | head -1)
    if [ -n "$PYTHON_CANDIDATE" ] && test_python "$PYTHON_CANDIDATE"; then
        PYTHON_CMD="$PYTHON_CANDIDATE"
    fi
fi

# 3. 嘗試 python.exe（Windows 上通常更可靠）
if [ -z "$PYTHON_CMD" ] && command -v python.exe &> /dev/null; then
    if test_python "python.exe"; then
        PYTHON_CMD="python.exe"
    fi
fi

# 4. 嘗試 python（不帶 .exe）
if [ -z "$PYTHON_CMD" ] && command -v python &> /dev/null; then
    if test_python "python"; then
        PYTHON_CMD="python"
    fi
fi

# 5. 最後嘗試 python3（可能被 Windows Store 快捷方式覆蓋）
if [ -z "$PYTHON_CMD" ] && command -v python3 &> /dev/null; then
    if test_python "python3"; then
        PYTHON_CMD="python3"
    fi
fi

# 如果還是找不到，嘗試查找其他常見位置
if [ -z "$PYTHON_CMD" ]; then
    # 查找 AppData\Roaming\Python
    if [ -d "/c/Users/$USER/AppData/Roaming/Python" ]; then
        PYTHON_CANDIDATE=$(find "/c/Users/$USER/AppData/Roaming/Python" -name "python.exe" 2>/dev/null | head -1)
        if [ -n "$PYTHON_CANDIDATE" ] && test_python "$PYTHON_CANDIDATE"; then
            PYTHON_CMD="$PYTHON_CANDIDATE"
        fi
    fi
fi

if [ -z "$PYTHON_CMD" ]; then
    echo -e "${RED}✗ Python is not installed or not found.${NC}"
    echo -e "${YELLOW}The 'python3' command may be a Windows Store shortcut.${NC}"
    echo -e "${YELLOW}Please:${NC}"
    echo -e "${YELLOW}  1. Install Python from https://www.python.org/downloads/${NC}"
    echo -e "${YELLOW}  2. During installation, check 'Add Python to PATH'${NC}"
    echo -e "${YELLOW}  3. Or disable Windows Store aliases in Settings > Apps > Advanced app settings${NC}"
    exit 1
fi

echo -e "${BLUE}Using Python: $PYTHON_VERSION${NC}"
echo -e "${BLUE}Python path: $PYTHON_CMD${NC}"

# 檢查 pip 是否安裝（使用找到的 Python 命令）
PIP_CMD=""
if $PYTHON_CMD -m pip --version &> /dev/null; then
    # 使用 python -m pip（最可靠的方式）
    PIP_CMD="$PYTHON_CMD -m pip"
elif command -v pip3 &> /dev/null; then
    # 測試 pip3 是否真的可用
    if pip3 --version 2>&1 | grep -qiE "pip [0-9]+\.[0-9]+"; then
        PIP_CMD="pip3"
    fi
elif command -v pip &> /dev/null; then
    # 測試 pip 是否真的可用
    if pip --version 2>&1 | grep -qiE "pip [0-9]+\.[0-9]+"; then
        PIP_CMD="pip"
    fi
fi

if [ -z "$PIP_CMD" ]; then
    echo -e "${RED}✗ pip is not installed or not found.${NC}"
    echo -e "${YELLOW}Please install pip: $PYTHON_CMD -m ensurepip --upgrade${NC}"
    exit 1
fi

echo -e "${BLUE}Using pip: $PIP_CMD${NC}"

# 安裝依賴項
echo -e "${BLUE}Checking/installing dependencies...${NC}"
if [ -f "requirements_recaptcha.txt" ]; then
    $PIP_CMD install -r requirements_recaptcha.txt
    if [ $? -ne 0 ]; then
        echo -e "${RED}✗ Failed to install dependencies${NC}"
        exit 1
    fi
else
    echo -e "${RED}✗ requirements_recaptcha.txt not found${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Dependencies installed${NC}"
echo ""

# 檢查腳本是否存在
if [ ! -f "get_recaptcha_automated.py" ]; then
    echo -e "${RED}✗ get_recaptcha_automated.py not found${NC}"
    exit 1
fi

# 解析命令行參數
BROWSER="chrome"
HEADLESS=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --firefox)
            BROWSER="firefox"
            shift
            ;;
        --edge)
            BROWSER="edge"
            shift
            ;;
        --headless)
            HEADLESS="--headless"
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [--firefox|--edge] [--headless]"
            echo ""
            echo "Options:"
            echo "  --firefox    Use Firefox instead of Chrome"
            echo "  --edge       Use Microsoft Edge (Windows)"
            echo "  --headless   Run in headless mode (no GUI)"
            echo "  --help, -h    Show this help message"
            echo ""
            echo "Requirements:"
            echo "- Python 3.x with pip"
            echo "- Chrome, Firefox, or Edge browser"
            echo "- ChromeDriver (for Chrome/Edge) or GeckoDriver (for Firefox)"
            echo ""
            echo "The script will automatically download the appropriate WebDriver."
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

echo -e "${BLUE}Starting reCAPTCHA token getter with ${BROWSER}...${NC}"
echo -e "${YELLOW}Note: You will still need to manually complete the reCAPTCHA challenge.${NC}"
echo ""

# 運行 Python 腳本
$PYTHON_CMD get_recaptcha_automated.py --browser $BROWSER $HEADLESS

# 檢查是否成功獲取 token
if [ -f "recaptcha_token.txt" ]; then
    echo ""
    echo -e "${GREEN}✓ Token retrieved successfully!${NC}"
    echo -e "${YELLOW}Token saved to: recaptcha_token.txt${NC}"

    # 顯示 token 並提供使用說明
    TOKEN=$(cat recaptcha_token.txt)
    echo -e "${BLUE}Token preview: ${TOKEN:0:50}...${NC}"
    echo ""
    echo -e "${GREEN}To use this token, run:${NC}"
    echo -e "export RECAPTCHA_TOKEN='$TOKEN'"
    echo -e "./scripts/submit-inquiry.sh"
    echo ""
    echo -e "${YELLOW}Or copy the token from recaptcha_token.txt${NC}"
else
    echo ""
    echo -e "${RED}✗ Token retrieval failed${NC}"
    echo -e "${YELLOW}Please check the error messages above${NC}"
    exit 1
fi
