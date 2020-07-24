# Survey system

Serve a creare dei questionari da somministrare ai colleghi.

L'accesso è basato su LDAP aziendale.

# Da eseguire prima del build
    go-bindata -fs static templates

# Versione per windows che non fa parire la gui
    go build -ldflags -H=windowsgui

# setta icona per versione exe
    rsrc -ico YOUR_ICON_FILE_NAME.ico