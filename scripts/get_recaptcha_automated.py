#!/usr/bin/env python3
"""
Automated reCAPTCHA Token Retrieval for University of Glasgow Inquiry Form

This script automates the browser interaction but still requires manual reCAPTCHA solving.
"""

import time
import json
import sys
import os
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service as ChromeService
from selenium.webdriver.edge.options import Options as EdgeOptions
from selenium.webdriver.edge.service import Service as EdgeService
from selenium.common.exceptions import TimeoutException, WebDriverException

# Firefox imports - 使用 Selenium 4.x 兼容方式
try:
    from selenium.webdriver.firefox.service import Service as FirefoxService
except ImportError:
    FirefoxService = None

# WebDriver Manager imports
try:
    from webdriver_manager.chrome import ChromeDriverManager
    from webdriver_manager.microsoft import EdgeChromiumDriverManager
    from webdriver_manager.firefox import GeckoDriverManager
    USE_WEBDRIVER_MANAGER = True
except ImportError:
    USE_WEBDRIVER_MANAGER = False
    ChromeDriverManager = None
    EdgeChromiumDriverManager = None
    GeckoDriverManager = None


class RecaptchaTokenGetter:
    def __init__(self, browser='chrome', headless=False):
        self.browser = browser.lower()
        self.headless = headless
        self.driver = None
        self.wait = None

        # Form data
        self.form_data = {
            'name': 'Pei Chen, Lee',
            'email': '17-DEC-96',
            'category': 'PGT',
            'country': 'Application Self Service Issues',
            'dateofbirth': 'Taiwan',
            'service': 'Service Desk',
            'subcategory': 'Uploading documents',
            'summary': 'PGT - Application Self Service Issues - Uploading documents',
            'symptom': '''pf frog <pcleegood@gmail.com>
11月15日 週六 下午12:45 (4 天前)
寄給 pgadmissions

Dear Postgraduate Admissions Team,

I am an applicant for MSc Software Development.

Full name: Pei-Chen Lee (李沛宸)

Application Number - 01658104

Programme: MSc Software Development

I chose the option for the University to contact my referee directly.

My referee, Yihsiu,Chen (Department of Digital Content and technologies, National Cheng Chi University), received an email with the link to the online reference form and the following user ID: UOG_OAS_RECOM_143070.

However, when he follow the link and log in, the system shows a "Search Criteria – User ID" page.
Following that he clicked  "Search", the system returns the message:
"No matching values were found."'''
        }

    def setup_driver(self):
        """Setup the web driver"""
        print(f"Setting up {self.browser} driver...")

        if self.browser == 'chrome':
            options = Options()
            if self.headless:
                options.add_argument('--headless')
            options.add_argument('--no-sandbox')
            options.add_argument('--disable-dev-shm-usage')
            options.add_argument('--disable-blink-features=AutomationControlled')
            options.add_experimental_option("excludeSwitches", ["enable-automation"])
            options.add_experimental_option('useAutomationExtension', False)
            # 啟用性能日誌以捕獲網絡請求
            options.set_capability('goog:loggingPrefs', {'performance': 'ALL'})
            options.add_argument('--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36')

            try:
                if USE_WEBDRIVER_MANAGER:
                    service = ChromeService(ChromeDriverManager().install())
                    self.driver = webdriver.Chrome(service=service, options=options)
                else:
                    self.driver = webdriver.Chrome(options=options)
            except Exception as e:
                print(f"Failed to initialize Chrome driver: {e}")
                print("Please ensure Chrome and ChromeDriver are installed.")
                print("Install ChromeDriver from: https://chromedriver.chromium.org/")
                sys.exit(1)

        elif self.browser == 'edge':
            options = EdgeOptions()
            if self.headless:
                options.add_argument('--headless')
            
            # 使用臨時用戶資料目錄以避免緩存的錯誤狀態
            import tempfile
            temp_profile = tempfile.mkdtemp(prefix='edge_profile_')
            options.add_argument(f'--user-data-dir={temp_profile}')
            options.add_argument('--profile-directory=Default')
            
            # 清除緩存和 cookies
            options.add_argument('--disable-application-cache')
            options.add_argument('--disk-cache-size=0')
            options.add_argument('--media-cache-size=0')
            
            options.add_argument('--no-sandbox')
            options.add_argument('--disable-dev-shm-usage')
            options.add_argument('--disable-blink-features=AutomationControlled')
            options.add_argument('--disable-gpu')
            options.add_argument('--disable-software-rasterizer')
            options.add_argument('--disable-extensions')
            options.add_argument('--disable-infobars')
            options.add_argument('--start-maximized')
            
            # 防止 Edge 自動關閉
            options.add_experimental_option("excludeSwitches", ["enable-automation"])
            options.add_experimental_option('useAutomationExtension', False)
            options.add_experimental_option("detach", True)  # 保持瀏覽器開啟
            
            # 啟用性能日誌以捕獲網絡請求
            options.set_capability('goog:loggingPrefs', {'performance': 'ALL'})
            
            # 禁用 Edge 的自動登入和同步，清除所有緩存
            options.add_experimental_option("prefs", {
                "credentials_enable_service": False,
                "profile.password_manager_enabled": False,
                "profile.default_content_setting_values.notifications": 2,
                "profile.default_content_settings.cookies": 1,
                "profile.block_third_party_cookies": False,
                "disk-cache-size": 0,
                "browser.cache.disk.enable": False,
                "browser.cache.memory.enable": False,
                "browser.cache.offline.enable": False
            })
            
            options.add_argument('--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0')

            # 嘗試自動找到 Edge 路徑
            edge_paths = [
                r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe",
                r"C:\Program Files\Microsoft\Edge\Application\msedge.exe",
                os.path.expanduser(r"~\AppData\Local\Microsoft\Edge\Application\msedge.exe")
            ]
            edge_path = None
            for path in edge_paths:
                if os.path.exists(path):
                    edge_path = path
                    break
            
            if edge_path:
                options.binary_location = edge_path
                print(f"✓ Found Edge at: {edge_path}")
            else:
                print("⚠ Could not find Edge installation automatically")

            # 嘗試多種方式初始化 Edge 驅動程序
            driver_initialized = False
            error_messages = []

            # 方法 1: 嘗試使用 webdriver-manager（如果可用且網絡正常）
            if USE_WEBDRIVER_MANAGER and EdgeChromiumDriverManager is not None:
                try:
                    print("Attempting to download EdgeDriver using webdriver-manager...")
                    service = EdgeService(EdgeChromiumDriverManager().install())
                    self.driver = webdriver.Edge(service=service, options=options)
                    driver_initialized = True
                    print("✓ EdgeDriver downloaded and initialized successfully")
                except Exception as e:
                    error_msg = f"webdriver-manager failed: {str(e)}"
                    error_messages.append(error_msg)
                    print(f"⚠ {error_msg}")
                    if "Could not reach host" in str(e) or "offline" in str(e).lower():
                        print("  → Network issue detected. Trying alternative methods...")

            # 方法 2: 嘗試使用系統 PATH 中的 msedgedriver
            if not driver_initialized:
                try:
                    print("Attempting to use msedgedriver from system PATH...")
                    self.driver = webdriver.Edge(options=options)
                    driver_initialized = True
                    print("✓ EdgeDriver found in system PATH")
                except Exception as e:
                    error_msg = f"System PATH driver failed: {str(e)}"
                    error_messages.append(error_msg)
                    print(f"⚠ {error_msg}")

            # 方法 3: 檢查環境變數指定的路徑
            if not driver_initialized:
                env_driver_path = os.environ.get('EDGEDRIVER_PATH')
                if env_driver_path and os.path.exists(env_driver_path):
                    try:
                        print(f"Attempting to use EdgeDriver from EDGEDRIVER_PATH: {env_driver_path}")
                        service = EdgeService(env_driver_path)
                        self.driver = webdriver.Edge(service=service, options=options)
                        driver_initialized = True
                        print(f"✓ EdgeDriver initialized from EDGEDRIVER_PATH")
                    except Exception as e:
                        error_msg = f"EDGEDRIVER_PATH driver failed: {str(e)}"
                        error_messages.append(error_msg)
                        print(f"⚠ {error_msg}")

            # 方法 4: 嘗試查找常見位置的 EdgeDriver
            if not driver_initialized:
                # 獲取腳本目錄
                script_dir = os.path.dirname(os.path.abspath(__file__))
                common_driver_paths = [
                    os.path.join(os.path.expanduser("~"), ".wdm", "drivers", "edgedriver"),
                    r"C:\Program Files\Microsoft\Edge\Application\msedgedriver.exe",
                    r"C:\Program Files (x86)\Microsoft\Edge\Application\msedgedriver.exe",
                    os.path.expanduser(r"~\AppData\Local\Microsoft\Edge\Application\msedgedriver.exe"),
                    os.path.join(script_dir, "msedgedriver.exe"),
                    os.path.join(script_dir, "msedgedriver"),
                    "./msedgedriver.exe",
                    "./msedgedriver",
                ]
                
                for driver_path in common_driver_paths:
                    # 如果是目錄，查找其中的可執行文件
                    if os.path.isdir(driver_path):
                        for root, dirs, files in os.walk(driver_path):
                            for file in files:
                                if "msedgedriver" in file.lower() and (file.endswith(".exe") or not file.endswith(".zip")):
                                    full_path = os.path.join(root, file)
                                    try:
                                        print(f"Attempting to use EdgeDriver at: {full_path}")
                                        service = EdgeService(full_path)
                                        self.driver = webdriver.Edge(service=service, options=options)
                                        driver_initialized = True
                                        print(f"✓ EdgeDriver initialized from: {full_path}")
                                        break
                                    except Exception as e:
                                        continue
                        if driver_initialized:
                            break
                    elif os.path.isfile(driver_path):
                        try:
                            print(f"Attempting to use EdgeDriver at: {driver_path}")
                            service = EdgeService(driver_path)
                            self.driver = webdriver.Edge(service=service, options=options)
                            driver_initialized = True
                            print(f"✓ EdgeDriver initialized from: {driver_path}")
                            break
                        except Exception as e:
                            continue

            # 如果所有方法都失敗
            if not driver_initialized:
                print("\n" + "="*60)
                print("✗ Failed to initialize Edge driver")
                print("="*60)
                print("\nError details:")
                for i, msg in enumerate(error_messages, 1):
                    print(f"  {i}. {msg}")
                print("\nSolutions:")
                print("1. Download EdgeDriver manually:")
                print("   https://developer.microsoft.com/en-us/microsoft-edge/tools/webdriver/")
                print("   - Extract msedgedriver.exe")
                print("   - Place it in the scripts directory or add to PATH")
                print("\n2. Check your internet connection and try again")
                print("   (webdriver-manager needs internet to download drivers)")
                print("\n3. Use Chrome or Firefox instead:")
                print("   ./gla-inquiry-manager.sh auto")
                print("   ./gla-inquiry-manager.sh auto --firefox")
                print("\n4. Set EDGEDRIVER_PATH environment variable:")
                print("   export EDGEDRIVER_PATH=/path/to/msedgedriver.exe")
                sys.exit(1)

        elif self.browser == 'firefox':
            # 使用 Selenium 4.x 的方式
            options = webdriver.FirefoxOptions()
            
            if self.headless:
                options.add_argument('--headless')
            options.set_preference("general.useragent.override",
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0")

            try:
                if USE_WEBDRIVER_MANAGER and FirefoxService is not None and GeckoDriverManager is not None:
                    service = FirefoxService(GeckoDriverManager().install())
                    self.driver = webdriver.Firefox(service=service, options=options)
                else:
                    self.driver = webdriver.Firefox(options=options)
            except Exception as e:
                print(f"Failed to initialize Firefox driver: {e}")
                print("Please ensure Firefox and GeckoDriver are installed.")
                sys.exit(1)
        else:
            print(f"Unsupported browser: {self.browser}")
            print("Supported browsers: chrome, firefox, edge")
            sys.exit(1)

        self.wait = WebDriverWait(self.driver, 30)
        self.driver.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")

        print(f"✓ {self.browser.capitalize()} driver initialized successfully")

    def check_session_alive(self):
        """檢查瀏覽器會話是否仍然有效"""
        try:
            # 嘗試獲取當前 URL 來檢查會話
            self.driver.current_url
            return True
        except Exception:
            return False

    def wait_for_page_load(self, url, max_retries=3):
        """等待頁面加載並檢查是否有錯誤"""
        for attempt in range(max_retries):
            try:
                if not self.check_session_alive():
                    print(f"⚠ Browser session lost, retrying... (attempt {attempt + 1}/{max_retries})")
                    return False
                
                print(f"Navigating to: {url}")
                self.driver.get(url)
                
                # 等待頁面加載
                time.sleep(3)
                
                # 檢查頁面內容是否有 reCAPTCHA 失敗訊息
                page_source = self.driver.page_source
                page_source_lower = page_source.lower()
                
                if "recaptcha fail" in page_source_lower or '"message":"recaptcha fail"' in page_source_lower:
                    print("⚠ WARNING: Page shows 'recaptcha fail' message on load")
                    print("  This is likely from a previous submission attempt.")
                    
                    # 檢查是否有實際的表單（不只是錯誤訊息）
                    has_form = "<form" in page_source and "name=" in page_source
                    has_input_fields = page_source.count("<input") > 3  # 至少有幾個輸入欄位
                    
                    if not has_form or not has_input_fields:
                        print("  ✗ No form found - page only shows error message!")
                        print("  Attempting to clear cookies and reload...")
                        
                        try:
                            # 清除所有 cookies
                            self.driver.delete_all_cookies()
                            print("  ✓ Cookies cleared")
                            time.sleep(1)
                            
                            # 重新加載頁面
                            self.driver.get(url)
                            time.sleep(3)
                            
                            page_source = self.driver.page_source
                            page_source_lower = page_source.lower()
                            
                            if "recaptcha fail" in page_source_lower:
                                print("  ⚠ Error still present after clearing cookies")
                                print("  This might be a server-side issue. Trying one more time...")
                                
                                # 最後一次嘗試：添加隨機參數以避免緩存
                                import random
                                cache_buster = f"?_={random.randint(1000000, 9999999)}"
                                self.driver.get(url + cache_buster)
                                time.sleep(3)
                                page_source = self.driver.page_source
                                page_source_lower = page_source.lower()
                                
                                if "recaptcha fail" in page_source_lower and "<form" not in page_source:
                                    print("  ✗ Still showing error without form")
                                    print("  The server may have blocked this session.")
                                    print("  Please wait a few minutes before trying again.")
                                    return False
                                else:
                                    print("  ✓ Form now available!")
                            else:
                                print("  ✓ Error cleared after clearing cookies")
                        except Exception as clear_error:
                            print(f"  ⚠ Could not clear cookies: {clear_error}")
                            return False
                    else:
                        print("  ✓ Form is present despite error message")
                        print("  Attempting to refresh page to clear the error...")
                        
                        try:
                            self.driver.refresh()
                            time.sleep(3)
                            page_source = self.driver.page_source
                            page_source_lower = page_source.lower()
                            
                            if "recaptcha fail" in page_source_lower:
                                print("  ⚠ Error still present after refresh")
                                print("  Form should still be accessible")
                            else:
                                print("  ✓ Error cleared after refresh")
                        except Exception as refresh_error:
                            print(f"  ⚠ Could not refresh page: {refresh_error}")
                
                # 等待頁面元素
                self.wait.until(EC.presence_of_element_located((By.TAG_NAME, "body")))
                
                # 檢查頁面是否包含表單（即使有錯誤訊息）
                if "form" in page_source_lower or "<form" in page_source:
                    print("✓ Page loaded and form detected")
                    return True
                else:
                    print("⚠ Form not found on page, but continuing...")
                    return True
                
            except Exception as e:
                print(f"⚠ Page load attempt {attempt + 1} failed: {e}")
                if attempt < max_retries - 1:
                    time.sleep(2)
                    continue
                return False
        return False

    def fill_form(self):
        """Fill out the inquiry form"""
        print("Filling out the inquiry form...")

        try:
            # 檢查會話是否有效
            if not self.check_session_alive():
                print("✗ Browser session is not alive. Please restart the script.")
                return False

            # Wait for the form to load with retry
            try:
                self.wait.until(EC.presence_of_element_located((By.TAG_NAME, "form")))
            except TimeoutException:
                # 檢查頁面是否有錯誤訊息
                page_source = self.driver.page_source
                if "recaptcha fail" in page_source.lower():
                    print("⚠ Page loaded but shows reCAPTCHA error")
                    print("  This might be from a previous submission attempt.")
                    print("  Trying to find form anyway...")
                else:
                    print("✗ Form not found on page")
                    print(f"  Current URL: {self.driver.current_url}")
                    print(f"  Page title: {self.driver.title}")
                    return False

            # Fill form fields
            field_mappings = {
                'name': 'name',
                'email': 'email',
                'category': 'category',
                'country': 'country',
                'dateofbirth': 'dateofbirth',
                'service': 'service',
                'subcategory': 'subcategory',
                'summary': 'summary',
                'symptom': 'symptom'
            }

            for field_name, element_name in field_mappings.items():
                try:
                    # 在每次操作前檢查會話
                    if not self.check_session_alive():
                        print(f"✗ Browser session lost while filling {field_name}")
                        return False
                    
                    element = self.wait.until(EC.element_to_be_clickable((By.NAME, element_name)))
                    element.clear()
                    element.send_keys(self.form_data[field_name])
                    print(f"✓ Filled {field_name}")
                    time.sleep(0.5)  # Small delay to avoid detection
                except Exception as e:
                    # 檢查是否是會話錯誤
                    error_str = str(e).lower()
                    if "invalid session" in error_str or "session deleted" in error_str or "disconnected" in error_str:
                        print(f"✗ Browser session lost: {e}")
                        print("  This usually happens when Edge closes unexpectedly.")
                        print("  Try running the script again or use Chrome/Firefox.")
                        return False
                    print(f"⚠ Failed to fill {field_name}: {e}")

            print("✓ Form filling completed")

        except TimeoutException:
            print("✗ Form did not load within timeout period")
            if self.check_session_alive():
                print(f"  Current URL: {self.driver.current_url}")
                print(f"  Page title: {self.driver.title}")
            return False
        except Exception as e:
            error_str = str(e).lower()
            if "invalid session" in error_str or "session deleted" in error_str or "disconnected" in error_str:
                print(f"✗ Browser session lost: {e}")
                print("  This usually happens when Edge closes unexpectedly.")
                print("  Solutions:")
                print("  1. Try running the script again")
                print("  2. Use Chrome: ./gla-inquiry-manager.sh auto")
                print("  3. Use Firefox: ./gla-inquiry-manager.sh auto --firefox")
                print("  4. Check if Edge is up to date")
            else:
                print(f"✗ Error filling form: {e}")
            return False

        return True

    def wait_for_recaptcha_and_submit(self):
        """Wait for user to complete reCAPTCHA and submit"""
        print("\n" + "="*50)
        print("FORM IS READY!")
        print("="*50)
        print()
        print("Please complete the following steps:")
        print("1. Complete the reCAPTCHA challenge")
        print("2. Click the submit button")
        print()
        print("The script will automatically detect when you submit...")
        print()

        # Wait for either successful submission or error
        while True:
            try:
                # 檢查會話是否有效
                if not self.check_session_alive():
                    print("✗ Browser session lost while waiting for submission")
                    return None

                # Check current URL for success/error indicators
                current_url = self.driver.current_url

                # Check for success/error messages in page content
                page_source = self.driver.page_source.lower()

                # 檢查 reCAPTCHA 失敗
                if "recaptcha fail" in page_source or '"message":"recaptcha fail"' in page_source:
                    print("⚠ reCAPTCHA validation failed detected on page")
                    print("  This means the token was invalid or expired.")
                    print("  Please try again with a fresh reCAPTCHA challenge.")
                    return None

                if "thank you" in page_source or "success" in page_source:
                    print("✓ Form submission appears successful!")
                    return True
                elif "error" in page_source or "fail" in page_source:
                    # 檢查是否是 reCAPTCHA 錯誤
                    if "recaptcha" in page_source:
                        print("⚠ reCAPTCHA error detected on page")
                        return None
                    print("⚠ Possible submission error detected")
                    return False

                # Check network requests for the form submission
                try:
                    logs = self.driver.get_log('performance')
                    for log in logs:
                        log_entry = json.loads(log['message'])['message']
                        if (log_entry.get('method') == 'Network.requestWillBeSent' and
                            'send/index.html' in log_entry.get('params', {}).get('request', {}).get('url', '')):
                            print("📡 Detected form submission request!")

                            # Extract the request payload
                            request = log_entry['params']['request']
                            if 'postData' in request:
                                post_data = request['postData']
                                print("📋 Extracting reCAPTCHA token from request...")

                                # Parse form data
                                form_data = {}
                                for param in post_data.split('&'):
                                    if '=' in param:
                                        key, value = param.split('=', 1)
                                        # URL decode
                                        try:
                                            import urllib.parse
                                            key = urllib.parse.unquote(key)
                                            value = urllib.parse.unquote(value)
                                        except:
                                            pass
                                        form_data[key] = value

                                if 'g-recaptcha-response' in form_data:
                                    token = form_data['g-recaptcha-response']
                                    if token and token != '' and token != 'null':
                                        print("🎉 SUCCESS! reCAPTCHA token extracted!")
                                        print(f"Token: {token[:50]}...")
                                        return token
                                else:
                                    print("⚠ g-recaptcha-response not found in form data")
                except Exception as log_error:
                    # 忽略日誌錯誤，繼續監視
                    pass

                time.sleep(1)  # Check every second

            except Exception as e:
                error_str = str(e).lower()
                if "invalid session" in error_str or "session deleted" in error_str or "disconnected" in error_str:
                    print(f"✗ Browser session lost: {e}")
                    return None
                # 其他錯誤繼續監視
                time.sleep(1)

    def run(self):
        """Main execution method"""
        print("🚀 University of Glasgow reCAPTCHA Token Getter")
        print("="*50)

        try:
            # Setup driver
            self.setup_driver()

            # Navigate to form with improved page loading
            form_url = "https://www.gla.ac.uk/study/enquire/send/index.html"
            if not self.wait_for_page_load(form_url):
                print("✗ Failed to load form page. Exiting.")
                return None

            # 再次檢查會話
            if not self.check_session_alive():
                print("✗ Browser session lost after page load. Exiting.")
                return None

            # Fill form
            if not self.fill_form():
                print("✗ Failed to fill form. Exiting.")
                return None

            # 再次檢查會話
            if not self.check_session_alive():
                print("✗ Browser session lost after filling form. Exiting.")
                return None

            # Wait for user interaction
            token = self.wait_for_recaptcha_and_submit()

            if token:
                print("\n" + "="*50)
                print("🎊 TOKEN SUCCESSFULLY RETRIEVED!")
                print("="*50)
                print(f"Full token: {token}")
                print()
                print("💡 To use this token, run:")
                print(f"export RECAPTCHA_TOKEN='{token}'")
                print("./scripts/submit-inquiry.sh")
                print()
                print("Token saved to: recaptcha_token.txt")

                # Save token to file
                with open('recaptcha_token.txt', 'w') as f:
                    f.write(token)
                print("✓ Token also saved to recaptcha_token.txt")

                return token
            else:
                print("✗ Failed to retrieve token")
                return None

        except KeyboardInterrupt:
            print("\n⚠ Script interrupted by user")
            return None
        except Exception as e:
            print(f"✗ Unexpected error: {e}")
            return None
        finally:
            if self.driver:
                print("🧹 Cleaning up...")
                self.driver.quit()

    @staticmethod
    def show_usage():
        """Show usage instructions"""
        print("University of Glasgow reCAPTCHA Token Getter")
        print("="*45)
        print()
        print("This script automates browser interaction but requires manual reCAPTCHA solving.")
        print()
        print("Requirements:")
        print("- Python 3.x")
        print("- Chrome, Firefox, or Edge browser")
        print("- Selenium WebDriver")
        print()
        print("Installation:")
        print("pip install selenium webdriver-manager")
        print("# WebDriver will be downloaded automatically")
        print()
        print("Usage:")
        print("python get_recaptcha_automated.py              # Use Chrome")
        print("python get_recaptcha_automated.py --browser edge    # Use Edge")
        print("python get_recaptcha_automated.py --browser firefox  # Use Firefox")
        print("python get_recaptcha_automated.py --headless   # Headless mode")
        print()
        print("The script will:")
        print("1. Open browser and navigate to the form")
        print("2. Automatically fill out the form")
        print("3. Wait for you to complete reCAPTCHA and submit")
        print("4. Extract the reCAPTCHA token from the submission")
        print()


def main():
    import argparse

    parser = argparse.ArgumentParser(description='Get reCAPTCHA token for University of Glasgow inquiry form')
    parser.add_argument('--browser', choices=['chrome', 'firefox', 'edge'], default='chrome',
                       help='Browser to use (default: chrome)')
    parser.add_argument('--headless', action='store_true',
                       help='Run in headless mode (no GUI)')

    args = parser.parse_args()

    if len(sys.argv) == 1:
        RecaptchaTokenGetter.show_usage()
        return

    getter = RecaptchaTokenGetter(browser=args.browser, headless=args.headless)
    token = getter.run()

    if token:
        print(f"\n✅ Success! Token retrieved and saved.")
        sys.exit(0)
    else:
        print(f"\n❌ Failed to retrieve token.")
        sys.exit(1)


if __name__ == "__main__":
    main()
