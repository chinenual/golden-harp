rsrc -manifest app.manifest -ico=icon_app.ico,icon_download.ico,icon_upload.ico -o rsrc.syso
go build -ldflags="-H windowsgui"
