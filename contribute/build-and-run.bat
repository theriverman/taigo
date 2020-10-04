@ECHO off

cls
go build -ldflags "-X main.taigaUsername=my-username -X main.taigaPassword=my-password -X main.sandboxProjectSlug=my-sandbox-slug -X main.sandboxEpicID=1234567890 -X main.sandboxFileUploadPath=C:\tmp\bad-puns-make-me-sic.png" -o taigo-dev.exe || ECHO Failed to build the binary && exit 1
.\taigo-dev.exe
