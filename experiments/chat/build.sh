#!/bin/bash

# Ordner für die Binaries erstellen, falls er nicht existiert
mkdir -p bin

# Lokales Binary erstellen
echo "Erstelle lokales Binary..."
go build -o ./bin/chat
echo "Lokales Binary erstellt: ./bin/chat"

# Windows-Binary erstellen
echo "Erstelle Windows-Binary..."
GOOS=windows GOARCH=amd64 go build -o ./bin/chat.exe
echo "Windows-Binary erstellt: ./bin/chat.exe"


# Abschlussmeldung
echo "Alle Binaries wurden im Ordner ./bin erstellt."
