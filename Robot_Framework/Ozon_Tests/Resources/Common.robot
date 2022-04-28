*** Settings ***
Library    SeleniumLibrary

*** Variables ***
${brow} =    ${BROWSER}

*** Keywords ***
Start Test
    open browser    about:blank    ${brow}

End Test
    close browser
