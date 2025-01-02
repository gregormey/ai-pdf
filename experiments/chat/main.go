package main

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

//go:embed content/Wahlprogramm_2024_Zukunft_Erleben_29.01.2024.txt
var fileContent string

// Add this struct for JSON response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//go:embed template/main.html
var htmlTemplate string

//go:embed prompt/prompt.txt
var promtTemplate string

type PageData struct {
	Question string
	Answer   string
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func checkOllamaAndModel() error {
	// Check if ollama is installed
	_, err := exec.LookPath("ollama")
	if err != nil {
		// Attempt to install ollama
		if installErr := installOllama(); installErr != nil {
			return fmt.Errorf("ollama is not installed and could not be installed: %v", installErr)
		}
		return nil // Retry after installation
	}

	// Check if ollama service is running
	cmd := exec.Command("ollama", "list")
	if err := cmd.Run(); err != nil {
		// Try to start ollama service
		fmt.Println("Starting ollama service...")
		cmd = exec.Command("ollama", "serve")
		cmd.Start()
		// Wait for service to start
		time.Sleep(2 * time.Second)
	}

	// Check if llama2 model is pulled
	cmd = exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check models: %v", err)
	}

	if !strings.Contains(string(output), "llama2") {
		fmt.Println("Pulling llama2 model... This may take a while...")
		cmd = exec.Command("ollama", "pull", "llama2")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to pull llama2 model: %v", err)
		}
	}

	return nil
}

// Function to install ollama
func installOllama() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows command to download and install ollama
		// Download the OllamaSetup.exe file
		cmd = exec.Command("powershell", "-Command", "Invoke-WebRequest -Uri https://ollama.com/download/OllamaSetup.exe -OutFile OllamaSetup.exe")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to download OllamaSetup.exe: %v, output: %s", err, output)
		}

		// Execute the downloaded installer
		cmd = exec.Command("OllamaSetup.exe")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to execute OllamaSetup.exe: %v", err)
		}
		return nil
	} else {
		// Unix-like command to download and install ollama
		cmd = exec.Command("bash", "-c", "curl -sSL https://ollama.ai/install.sh | sh")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install ollama: %v, output: %s", err, output)
	}
	return nil
}

func main() {
	// Check ollama and model before starting
	if err := checkOllamaAndModel(); err != nil {
		log.Fatal(err)
	}

	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Fehler beim Rendern der Seite", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		for {
			_, question, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				break
			}

			log.Printf("Neue Frage: %s", question)

			// Send processing status
			conn.WriteJSON(Response{
				Status:  "processing",
				Message: "Verarbeite Anfrage...",
			})

			prompt := prompts.NewPromptTemplate(promtTemplate, []string{"content", "question"})

			promptText, err := prompt.Format(map[string]any{
				"content":  fileContent,
				"question": string(question),
			})
			if err != nil {
				conn.WriteJSON(Response{
					Status:  "error",
					Message: "Fehler bei der Verarbeitung der Frage",
				})
				continue
			}

			ctx := context.Background()
			answer, err := llm.Call(ctx, promptText,
				llms.WithTemperature(0.8),
				llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
					return conn.WriteJSON(Response{
						Status:  "streaming",
						Message: string(chunk),
					})
				}),
			)
			if err != nil {
				conn.WriteJSON(Response{
					Status:  "error",
					Message: "Fehler: " + err.Error(),
				})
				continue
			}

			log.Printf("Neue Antwort: %s", answer)
		}
	})

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 2 * time.Minute, // Increased timeout for LLM responses
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		fmt.Println("Server startet auf http://localhost:8080")
		log.Fatal(server.ListenAndServe())
	}()

	// Warte kurz, bis der Server gestartet ist
	time.Sleep(100 * time.Millisecond)

	// Öffne den Browser
	if err := openBrowser("http://localhost:8080"); err != nil {
		log.Printf("Fehler beim Öffnen des Browsers: %v", err)
	}

	// Halte das Programm am Laufen
	select {}
}
