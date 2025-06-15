Output ./assets/mash-demo.gif

Set WindowBar Colorful
Set WindowBarSize 70
Set Theme "nord"
Set Shell "bash"
Set FontSize 24
Set Width 1400
Set Height 1000


Type "mash"
Sleep 1s
Enter
Sleep 1s
Enter
Wait+Screen /Playwright/
Sleep 1s
Type "/"
Sleep 1s
Type "Playwright"
Sleep 400ms
Enter
Sleep 400ms
Enter
Hide
Wait+Screen /Playwright-cli-select/
Show
Sleep 2s

Hide
Ctrl+c
Sleep 1s
Type "clear"
Enter

Show

Type 'mash --global --tag'
Sleep 1s
Enter
Sleep 1s
Type "/"
Sleep 1s
Type "Goodbye"
Sleep 1s
Enter
Sleep 1s
Enter
Sleep 2s

Hide
Ctrl+c
Sleep 1s
Type "clear"
Enter

Show

Type 'mash --global --tag="clock;cypress"'
Sleep 1s
Enter
Sleep 1s
Down
Sleep 1s
Enter
Sleep 3s

Hide
Ctrl+c
Sleep 1s
Type "clear"
Enter

Show

Type 'mash --global --tree'
Sleep 1s
Enter
Sleep 3s

Hide
Ctrl+c
Sleep 1s
Type "clear"
Enter

Show

Type "mash --help"
Sleep 1s
Enter
Sleep 3s

Hide
Ctrl+c
Sleep 1s
Type "clear"
Enter

Show 

