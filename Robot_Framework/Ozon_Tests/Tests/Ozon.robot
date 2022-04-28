*** Settings ***
Documentation    First test suite
Resource    ../Resources/OzonApp.robot
Resource    ../Resources/Common.robot
Test Setup    Common.Start Test
Test Teardown    Common.End Test

*** Variables ***
${BROWSER} =    edge
${URL} =    https://www.ozon.ru

*** Test Cases ***
User should be able to open "Пункты выдачи" page
    [Tags]    Stores
    OzonApp.Open a page
    OzonApp.Verify stores list page is opened

"Пункты выдачи" page should allow to select a store
    [Tags]    Stores selector
#    Sleep    1s
    OzonApp.Open a page
    OzonApp.Verify stores list page is opened
    OzonApp.Verify stores list page contains selector
