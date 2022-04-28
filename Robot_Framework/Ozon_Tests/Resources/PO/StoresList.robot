*** Settings ***
Library    SeleniumLibrary

*** Variables ***
${stores_list_heading_tag} =    tag:h1
${stores_list_heading} =    пункты выдачи заказов OZON
${stores_selector_locator} =    xpath://*[@id="layoutPage"]/div[1]/div[3]/div[1]/div/div[2]/div[1]
${stores_selector_heading} =    Выбор пункта выдачи

*** Keywords ***
Verify header name
    wait until element contains    ${stores_list_heading_tag}    ${stores_list_heading}

Verify selector exists
    wait until element contains    ${stores_selector_locator}    ${stores_selector_heading}
