#!/bin/bash

# Python 查找和診斷腳本
# 幫助在 Windows Git Bash 中找到正確的 Python 安裝

echo "=========================================="
echo "Python Detection and Diagnostic Tool"
echo "=========================================="
echo ""

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

# 測試函數
test_python() {
    local cmd="$1"
    local test_output
    test_output=$($cmd --version 2>&1)
    
    if echo "$test_output" | grep -qi "Microsoft Store\|App execution aliases"; then
        echo -e "${RED}✗ Windows Store shortcut (not real Python)${NC}"
        return 1
    fi
    
    if echo "$test_output" | grep -qiE "Python [0-9]+\.[0-9]+"; then
        echo -e "${GREEN}✓ Valid Python: $test_output${NC}"
        return 0
    fi
    
    echo -e "${RED}✗ Invalid output: $test_output${NC}"
    return 1
}

echo -e "${BLUE}Testing common Python commands...${NC}"
echo ""

# 測試 python3
echo -n "python3: "
if command -v python3 &> /dev/null; then
    test_python "python3"
else
    echo -e "${YELLOW}Not found${NC}"
fi

# 測試 python
echo -n "python: "
if command -v python &> /dev/null; then
    test_python "python"
else
    echo -e "${YELLOW}Not found${NC}"
fi

# 測試 python.exe
echo -n "python.exe: "
if command -v python.exe &> /dev/null; then
    test_python "python.exe"
else
    echo -e "${YELLOW}Not found${NC}"
fi

echo ""
echo -e "${BLUE}Searching for Python installations...${NC}"
echo ""

# 搜索常見安裝路徑
SEARCH_PATHS=(
    "/c/Users/$USER/AppData/Local/Programs/Python"
    "/c/Python*"
    "/c/Users/$USER/AppData/Roaming/Python"
    "/c/Program Files/Python*"
    "/c/Program Files (x86)/Python*"
)

FOUND_PYTHON=0

for search_path in "${SEARCH_PATHS[@]}"; do
    if [ -d "$search_path" ] 2>/dev/null; then
        echo -e "${BLUE}Searching: $search_path${NC}"
        PYTHON_FOUND=$(find "$search_path" -name "python.exe" 2>/dev/null | head -1)
        if [ -n "$PYTHON_FOUND" ]; then
            echo -e "${GREEN}  Found: $PYTHON_FOUND${NC}"
            if test_python "$PYTHON_FOUND"; then
                echo -e "${GREEN}  ✓ This Python works!${NC}"
                echo ""
                echo -e "${GREEN}You can use this Python with:${NC}"
                echo -e "${YELLOW}  export PYTHON_CMD=\"$PYTHON_FOUND\"${NC}"
                echo -e "${YELLOW}  Or add to PATH${NC}"
                FOUND_PYTHON=1
            fi
        fi
    fi
done

echo ""
if [ $FOUND_PYTHON -eq 0 ]; then
    echo -e "${RED}No working Python installation found.${NC}"
    echo ""
    echo -e "${YELLOW}Solutions:${NC}"
    echo "1. Install Python from https://www.python.org/downloads/"
    echo "2. During installation, check 'Add Python to PATH'"
    echo "3. Or disable Windows Store aliases:"
    echo "   Settings > Apps > Advanced app settings > App execution aliases"
    echo "   Turn off 'python.exe' and 'python3.exe'"
else
    echo -e "${GREEN}Python installation found!${NC}"
fi

echo ""
echo "=========================================="
