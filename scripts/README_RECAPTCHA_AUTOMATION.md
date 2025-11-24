# Automated reCAPTCHA Token Retrieval

This guide explains how to use automated scripts to obtain fresh reCAPTCHA tokens for the University of Glasgow inquiry form.

## ⚠️ Important Limitations

**reCAPTCHA cannot be fully automated** because:
- reCAPTCHA v2/v3 requires human interaction to solve challenges
- Google actively detects and blocks automated solving attempts
- Completely automated solving violates Google's Terms of Service

The scripts provided here **semi-automate** the process by handling browser setup and form filling, but you still need to **manually complete the reCAPTCHA challenge**.

## 📋 Available Scripts

### 1. `get_recaptcha_token.sh` - Manual Helper
Simple script that provides instructions and basic form page fetching.

### 2. `run_recaptcha_getter.sh` - Semi-Automated Browser Script
Uses Selenium to automatically open browser, fill form, and extract token after manual reCAPTCHA completion.

### 3. `get_recaptcha_automated.py` - Python Selenium Script
The core automation script (called by `run_recaptcha_getter.sh`).

## 🚀 Quick Start

### Prerequisites

1. **Python 3.x** with pip
2. **Chrome** or **Firefox** browser
3. **WebDriver** for your chosen browser

### Installation

```bash
cd scripts

# Install Python dependencies
pip install -r requirements_recaptcha.txt

# For Chrome: Download ChromeDriver
# For Firefox: Download GeckoDriver
# Add the driver to your PATH or place in the scripts directory
```

### Run the Automated Script

```bash
# Use Chrome (default)
./run_recaptcha_getter.sh

# Use Firefox
./run_recaptcha_getter.sh --firefox

# Use headless mode (no GUI)
./run_recaptcha_getter.sh --headless
```

## 📖 Detailed Instructions

### Step 1: Run the Script

```bash
cd scripts
./run_recaptcha_getter.sh
```

The script will:
1. Check dependencies
2. Install required Python packages
3. Open your browser
4. Navigate to the University of Glasgow form
5. Automatically fill out the form fields

### Step 2: Complete reCAPTCHA Manually

When the script shows:
```
================================
FORM IS READY!
================================

Please complete the following steps:
1. Complete the reCAPTCHA challenge
2. Click the submit button

The script will automatically detect when you submit...
```

**You need to:**
1. ✅ Complete the reCAPTCHA challenge (checkbox + any additional puzzles)
2. 🖱️ Click the **Submit** button

### Step 3: Token Extraction

The script automatically:
- Monitors the network requests
- Detects when you submit the form
- Extracts the `g-recaptcha-response` token
- Saves it to `recaptcha_token.txt`

### Step 4: Use the Token

```bash
# Option 1: Set environment variable
export RECAPTCHA_TOKEN="$(cat recaptcha_token.txt)"
./scripts/submit-inquiry.sh

# Option 2: Direct usage
RECAPTCHA_TOKEN="$(cat recaptcha_token.txt)" ./scripts/submit-inquiry.sh
```

## 🔧 Troubleshooting

### Common Issues

#### 1. WebDriver Not Found
```
Message: 'chromedriver' executable needs to be in PATH
```

**Solution:**
```bash
# Download ChromeDriver
# For Chrome: https://chromedriver.chromium.org/
# For Firefox: https://github.com/mozilla/geckodriver/releases

# Add to PATH or place in scripts directory
export PATH=$PATH:/path/to/driver
```

#### 2. Browser Not Detected
```
selenium.common.exceptions.WebDriverException: Message: chrome not reachable
```

**Solutions:**
- Update your browser to the latest version
- Try Firefox instead: `./run_recaptcha_getter.sh --firefox`
- Use headless mode: `./run_recaptcha_getter.sh --headless`

#### 3. Form Fields Not Found
```
Failed to fill field_name: element not found
```

**Solution:**
The website structure may have changed. Use the manual method instead.

#### 4. Token Not Extracted
```
g-recaptcha-response not found in form data
```

**Solution:**
- Ensure you complete reCAPTCHA before submitting
- Check that the form actually submits (look for success/error messages)

### Dependencies Installation Issues

```bash
# If pip fails
python -m ensurepip --upgrade
pip install --user -r requirements_recaptcha.txt

# Or use conda
conda install selenium
```

## 🔄 Alternative Methods

### Manual Method (Most Reliable)

1. Open browser manually: https://www.gla.ac.uk/study/enquire/send/index.html
2. Open Developer Tools (F12) → Network tab
3. Fill form and complete reCAPTCHA
4. Submit form
5. Find POST request, extract `g-recaptcha-response`

### Browser Extensions

Use extensions that help with form filling and reCAPTCHA monitoring:
- **reCAPTCHA Solver** (helps with some challenges)
- **Form Filler** (auto-fill form data)
- **Network Monitor** (inspect requests)

## ⚖️ Legal and Ethical Considerations

- **Allowed**: Automating form filling and token extraction for personal use
- **Borderline**: Using reCAPTCHA solving services (may violate terms)
- **Not Allowed**: Completely automated reCAPTCHA solving without human interaction
- **Violation**: Using automated tools to bypass security measures for malicious purposes

## 📁 File Structure

```
scripts/
├── get_recaptcha_token.sh           # Manual helper script
├── run_recaptcha_getter.sh          # Semi-automated wrapper
├── get_recaptcha_automated.py       # Core Selenium script
├── requirements_recaptcha.txt       # Python dependencies
├── README_RECAPTCHA_AUTOMATION.md   # This documentation
└── recaptcha_token.txt              # Generated token (after successful run)
```

## 🎯 Success Indicators

### Script Output
```
🎉 SUCCESS! reCAPTCHA token extracted!
Token: 0.ABC123...[truncated]
✓ Token also saved to recaptcha_token.txt
```

### Token Format
Valid tokens look like:
```
0.ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.
ABCDEF... (very long string, usually 1000+ characters)
```

## 🆘 Getting Help

If you encounter issues:

1. **Check the error messages** in the script output
2. **Try different browsers** (Chrome vs Firefox)
3. **Use manual method** if automation fails
4. **Update dependencies** to latest versions
5. **Check browser compatibility**

The automated scripts make the process easier, but the **manual method is always available as a fallback**.
