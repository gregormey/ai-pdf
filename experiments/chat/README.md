# Grüne Wahlprogramm Chat

Eine Webanwendung zum Stellen von Fragen über das Wahlprogramm der Grünen, basierend auf Mistral und Ollama.

## Voraussetzungen

1. Go 1.21 oder höher
2. Ollama (wird automatisch installiert, falls nicht vorhanden)

### Installation

Die Anwendung handhabt die Installation automatisch:

1. Wenn Ollama nicht installiert ist, wird die Installation gestartet
2. Nach der Installation müssen Sie:
   - Warten bis die Installation abgeschlossen ist
   - Die Anwendung neu starten
3. Das LLaMA 3.3 Modell wird automatisch heruntergeladen, wenn es nicht vorhanden ist

### Manuelle Installation (optional)

Falls Sie Ollama manuell installieren möchten:

1. Installieren Sie Ollama von [ollama.ai](https://ollama.ai)
2. Laden Sie das LLaMA 3.3 Modell:
   ```bash
   ollama pull llama3.3
   ```

### Installation der Abhängigkeiten

1. Klonen Sie das Repository
2. Navigieren Sie zum Chat-Verzeichnis:
   ```bash
   cd ai-pdf/experiments/chat
   go mod tidy
   ```

## Starten der Anwendung

1. Führen Sie die Anwendung aus:
   ```bash
   go run main.go
   ```
2. Die Anwendung wird:
   - Ollama installieren (falls nötig)
   - Das Modell herunterladen (falls nötig)
   - Den Ollama-Dienst starten
   - Einen Browser mit der Anwendung öffnen

## Nutzung

1. Geben Sie Ihre Frage zum Wahlprogramm in das Textfeld ein
2. Klicken Sie auf "Frage stellen"
3. Die Antwort wird live generiert und angezeigt

## Fehlerbehebung

- Bei der Erstinstallation:
  - Folgen Sie den Anweisungen auf dem Bildschirm
  - Starten Sie die Anwendung nach der Installation neu
- Bei Verbindungsproblemen:
  - Überprüfen Sie, ob Port 8080 verfügbar ist
  - Stellen Sie sicher, dass der Ollama-Dienst läuft (`ollama serve`)
- Bei Modell-Problemen:
  - Überprüfen Sie die Installation mit `ollama list`
  - Das Modell sollte als "llama3.3" angezeigt werden

## Technischer Stack

- Go
- Ollama mit LLaMA 3.3
- WebSocket für Streaming-Antworten
- HTML/JavaScript für das Frontend
