from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as ec
from selenium.webdriver.support.wait import WebDriverWait

options = webdriver.FirefoxOptions()
options.browser_version = '91.6.0'

# driver = webdriver.Remote(command_executor='http://127.0.0.1:4444', options=options)
driver = webdriver.Firefox(options=options)

try:
    driver.get('https://the-internet.herokuapp.com')

    el = driver.find_element(By.LINK_TEXT, 'Notification Messages')

    wait = WebDriverWait(driver, 5)
    wait.until(ec.presence_of_element_located((By.CSS_SELECTOR, '#page-footer')))

    els = driver.find_elements(By.TAG_NAME, 'a')
    print(f'There are {len(els)} "a" tags')
    els = driver.find_elements(By.TAG_NAME, 'myTag')
    print(f'There are {len(els)} "myTag" tags')

    el.click()
    flash = wait.until(ec.presence_of_element_located((By.CSS_SELECTOR, '#flash')))

    assert 'Action' in flash.text

finally:
    driver.quit()
