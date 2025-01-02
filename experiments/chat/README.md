# Grüne Wahlprogramm Chat

Eine Webanwendung zum Stellen von Fragen über das Wahlprogramm der Grünen, basierend auf LLaMA 2 und Ollama.

## Voraussetzungen

1. Go 1.21 oder höher
2. Ollama installiert und konfiguriert

### Installation von Ollama

1. Installieren Sie Ollama von [ollama.ai](https://ollama.ai) (optional, die Anwendung kann es auch automatisch installieren).
2. Öffnen Sie ein Terminal und laden Sie das LLaMA 2 Modell:
   ```bash
   ollama pull llama2
   ```

### Automatische Installation

Die Anwendung überprüft, ob `ollama` installiert ist. Wenn nicht, wird versucht, `ollama` automatisch zu installieren. Stellen Sie sicher, dass Sie über die erforderlichen Berechtigungen verfügen, um Skripte auszuführen.

- **Für Windows**: Die Anwendung verwendet PowerShell, um das Installationsskript herunterzuladen und auszuführen.
- **Für Unix-ähnliche Systeme**: Die Anwendung verwendet `curl`, um das Installationsskript herunterzuladen und auszuführen.

### Installation der Abhängigkeiten

1. Klonen Sie das Repository
2. Navigieren Sie zum Chat-Verzeichnis:

```bash
go mod tidy
```

## Starten der Anwendung

1. Stellen Sie sicher, dass der Ollama-Dienst läuft (oder lassen Sie die Anwendung es automatisch installieren).
2. Führen Sie die Anwendung aus:
```bash
go run main.go
```
3. Ein Browser-Fenster öffnet sich automatisch mit der Anwendung unter `http://localhost:8080`

## Nutzung

1. Geben Sie Ihre Frage zum Wahlprogramm in das Textfeld ein
2. Klicken Sie auf "Frage stellen"
3. Die Antwort wird live generiert und angezeigt

## Fehlerbehebung

- Wenn Sie einen Fehler mit Ollama erhalten, stellen Sie sicher, dass der Dienst läuft:
```bash
ollama serve
```
- Bei Verbindungsproblemen überprüfen Sie, ob Port 8080 verfügbar ist
- Bei Modell-Fehlern stellen Sie sicher, dass LLaMA 2 korrekt installiert ist:
```bash
ollama list
```

## Technischer Stack

- Go
- Ollama (LLaMA 2)
- WebSocket für Streaming-Antworten
- HTML/JavaScript für das Frontend
