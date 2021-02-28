go build -o taigo-dev.exe -ldflags "-X main.taigaUsername=$env:taigaUsername -X main.taigaPassword=$env:taigaPassword -X main.sandboxProjectSlug=$env:sandboxProjectSlug -X main.sandboxEpicID=$env:sandboxEpicID -X main.sandboxFileUploadPath=$env:sandboxFileUploadPath"
.\taigo-dev.exe
