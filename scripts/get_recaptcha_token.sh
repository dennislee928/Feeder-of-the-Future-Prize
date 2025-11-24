#!/bin/bash

# Script to help obtain a fresh reCAPTCHA token for University of Glasgow inquiry form
# 獲取格拉斯哥大學諮詢表單的最新 reCAPTCHA token 助手腳本

# 顏色輸出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}reCAPTCHA Token Getter${NC}"
echo -e "${GREEN}================================${NC}"
echo ""

echo -e "${YELLOW}This script will help you obtain a fresh reCAPTCHA token.${NC}"
echo -e "${YELLOW}Since reCAPTCHA requires human interaction, you'll need to complete the challenge manually.${NC}"
echo ""

# 檢查是否安裝了必要的工具
check_dependencies() {
    echo -e "${BLUE}Checking dependencies...${NC}"

    if ! command -v curl &> /dev/null; then
        echo -e "${RED}✗ curl is not installed. Please install curl first.${NC}"
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}⚠ jq is not installed. JSON parsing will be limited.${NC}"
        echo -e "${YELLOW}  Install jq for better output: https://stedolan.github.io/jq/download/${NC}"
        HAS_JQ=false
    else
        HAS_JQ=true
    fi

    echo -e "${GREEN}✓ Dependencies check completed${NC}"
    echo ""
}

# 生成隨機用戶代理
get_random_user_agent() {
    local agents=(
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0"
    )
    echo "${agents[$RANDOM % ${#agents[@]}]}"
}

# 獲取表單頁面
get_form_page() {
    echo -e "${BLUE}Fetching the inquiry form page...${NC}"

    local user_agent=$(get_random_user_agent)
    local response_file=$(mktemp)

    curl -s -A "$user_agent" \
         -H "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" \
         -H "Accept-Language: en-US,en;q=0.5" \
         -H "Accept-Encoding: gzip, deflate" \
         -H "Connection: keep-alive" \
         -H "Upgrade-Insecure-Requests: 1" \
         -o "$response_file" \
         "https://www.gla.ac.uk/study/enquire/send/index.html"

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Form page fetched successfully${NC}"
        echo -e "${YELLOW}Page saved to: $response_file${NC}"

        # 檢查頁面是否包含 reCAPTCHA
        if grep -q "recaptcha" "$response_file"; then
            echo -e "${GREEN}✓ reCAPTCHA detected on the page${NC}"
        else
            echo -e "${YELLOW}⚠ reCAPTCHA not found on the page${NC}"
        fi

        echo ""
        echo -e "${YELLOW}Next steps:${NC}"
        echo -e "1. Open the page in your browser: https://www.gla.ac.uk/study/enquire/send/index.html"
        echo -e "2. Open Developer Tools (F12) and go to Network tab"
        echo -e "3. Fill out the form and complete the reCAPTCHA challenge"
        echo -e "4. Click submit"
        echo -e "5. Find the POST request in Network tab"
        echo -e "6. Look for 'g-recaptcha-response' in the request payload"
        echo -e "7. Copy the token value"
        echo ""
    else
        echo -e "${RED}✗ Failed to fetch form page${NC}"
        rm -f "$response_file"
        exit 1
    fi
}

# 測試 token
test_token() {
    local token="$1"

    if [ -z "$token" ]; then
        echo -e "${RED}✗ No token provided${NC}"
        return 1
    fi

    echo -e "${BLUE}Testing reCAPTCHA token...${NC}"

    # 這裡可以添加一個簡單的測試來驗證 token 格式
    if [[ "$token" =~ ^0\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$ ]]; then
        echo -e "${GREEN}✓ Token format appears valid${NC}"
        return 0
    else
        echo -e "${RED}✗ Token format appears invalid${NC}"
        return 1
    fi
}

# 顯示使用說明
show_instructions() {
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}Manual reCAPTCHA Token Retrieval${NC}"
    echo -e "${GREEN}================================${NC}"
    echo ""
    echo -e "${YELLOW}Step-by-step instructions:${NC}"
    echo ""
    echo -e "1. ${BLUE}Open your browser and visit:${NC}"
    echo -e "   https://www.gla.ac.uk/study/enquire/send/index.html"
    echo ""
    echo -e "2. ${BLUE}Open Developer Tools:${NC}"
    echo -e "   - Press F12"
    echo -e "   - Click on the 'Network' tab"
    echo ""
    echo -e "3. ${BLUE}Fill out the form with your information:${NC}"
    echo -e "   - Name: Pei Chen, Lee"
    echo -e "   - Email: 17-DEC-96"
    echo -e "   - Category: PGT"
    echo -e "   - Service: Application Self Service Issues"
    echo -e "   - Subcategory: Uploading documents"
    echo -e "   - Summary: PGT - Application Self Service Issues - Uploading documents"
    echo -e "   - Symptom: [your message content]"
    echo ""
    echo -e "4. ${BLUE}Complete the reCAPTCHA challenge:${NC}"
    echo -e "   - Click the reCAPTCHA checkbox"
    echo -e "   - Solve any additional challenges if presented"
    echo ""
    echo -e "5. ${BLUE}Click the submit button${NC}"
    echo ""
    echo -e "6. ${BLUE}Find the network request:${NC}"
    echo -e "   - In Network tab, look for a POST request to '/study/enquire/send/index.html'"
    echo -e "   - Click on it to see details"
    echo ""
    echo -e "7. ${BLUE}Extract the reCAPTCHA token:${NC}"
    echo -e "   - Look under 'Form Data' or 'Payload'"
    echo -e "   - Find the 'g-recaptcha-response' field"
    echo -e "   - Copy the long token value (starts with '0.')"
    echo ""
    echo -e "8. ${BLUE}Set the environment variable:${NC}"
    echo -e "   export RECAPTCHA_TOKEN='your_copied_token_here'"
    echo ""
    echo -e "9. ${BLUE}Run the submission script:${NC}"
    echo -e "   ./scripts/submit-inquiry.sh"
    echo ""
}

# 主函數
main() {
    check_dependencies
    get_form_page
    show_instructions

    echo -e "${YELLOW}Would you like to test a token? Paste it below (or press Enter to skip):${NC}"
    echo -e "${YELLOW}Token: ${NC}"
    read -r test_token_input

    if [ -n "$test_token_input" ]; then
        if test_token "$test_token_input"; then
            echo ""
            echo -e "${GREEN}✓ Token validation passed!${NC}"
            echo -e "${YELLOW}Set it as environment variable:${NC}"
            echo -e "export RECAPTCHA_TOKEN='$test_token_input'"
            echo ""
            echo -e "${GREEN}Then run: ./scripts/submit-inquiry.sh${NC}"
        fi
    else
        echo -e "${YELLOW}Skipping token test. Follow the manual instructions above.${NC}"
    fi
}

# 檢查命令行參數
case "$1" in
    --help|-h)
        show_instructions
        exit 0
        ;;
    --test)
        if [ -n "$2" ]; then
            test_token "$2"
        else
            echo -e "${RED}Usage: $0 --test 'your_token_here'${NC}"
            exit 1
        fi
        ;;
    *)
        main
        ;;
esac
