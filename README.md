Tools
=====

## Email notify for procmail on mac ##
Use procmail+**THIS TOOL**+growlnotify to get email notified on mac.

### Install growl & growlnotify ###
Get growl & growlnotify from [http://growl.info/downloads](http://growl.info/downloads)

### Build tool ###
`go build -o /usr/local/bin/growlforprocmail growlforprocmail.go`

### Edit procmailrc ###
Add below to `~/.procmailrc`
`
:0c
|/usr/bin/formail -X From: -X Subject: |/usr/local/bin/growlforprocmail |/usr/local/bin/growlnotify -a Mail.app -t NewEmail >/dev/null 2>&1
`

### Get done ###
Now you can get notice from growl like:
* sender \<sender@sender.com\>
 email title*
