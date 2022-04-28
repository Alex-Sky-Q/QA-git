*** Settings ***
Resource    ./PO/MainPage.robot
Resource    ./PO/StoresList.robot
Resource    ./PO/TopMenu.robot

*** Variables ***
${PAGE} =    Пункты выдачи

*** Keywords ***
Open a page
    MainPage.Open main page
    TopMenu.Click menu link

Verify stores list page is opened
    StoresList.Verify header name

Verify stores list page contains selector
    StoresList.Verify selector exists
