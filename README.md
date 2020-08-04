# Survey system
[![Maintainability](https://api.codeclimate.com/v1/badges/988f3974aa0eb7372050/maintainability)](https://codeclimate.com/github/axamon/surveysystem/maintainability)

Serve a creare dei questionari da somministrare ai colleghi.

L'accesso è basato su LDAP aziendale.

# Da eseguire prima del build
    go-bindata -fs static templates

# Versione per windows che non fa parire la gui
    go build -ldflags -H=windowsgui

# setta icona per versione exe
    rsrc -ico YOUR_ICON_FILE_NAME.ico

# installare go-bindata
    go get -u github.com/go-bindata/go-bindata/...

# Build
    go generate
    go build

# Endpoints gcp
    https://europe-west6-ctio-8274d.cloudfunctions.net/SheetAppend2
# GCP endpoint

    https://europe-west6-ctio-8274d.cloudfunctions.net/SheetAppend

Update function

    gcloud functions deploy SheetAppend --region europe-west6 --trigger-http --allow-unauthenticated

# Vedere risultati
    https://docs.google.com/spreadsheets/d/1dKXJ2bm_ZYm3tlIMmFcFfM4hjtKXmqndigjekd_H_yo/edit?usp=sharing
