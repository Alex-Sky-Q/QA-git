*** Settings ***
Library    SeleniumLibrary

*** Variables ***
${menu_link} =    ${PAGE}
${topmenu_stores_link_css_locator} =    css=#layoutPage > div.gt > div.c9q > div > ul > li:nth-child(7) > div > a

*** Keywords ***
Click menu link
#    click link    ${menu_link}
    click element    ${topmenu_stores_link_css_locator}
