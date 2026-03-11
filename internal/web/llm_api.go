package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// LLMService manages the vLLM server and API interactions
type LLMService struct {
	mu            sync.RWMutex
	enabled       bool
	serverRunning bool
	serverCmd     *exec.Cmd
	vllmBaseURL   string
	modelPath     string
	loraPath      string
	port          string
}

// LLMConfig holds LLM configuration
type LLMConfig struct {
	BaseModel string `json:"base_model"`
	LoRAPath  string `json:"lora_path"`
	Port      string `json:"port"`
}

// NewLLMService creates a new LLM service
func NewLLMService() *LLMService {
	return &LLMService{
		enabled:     false,
		vllmBaseURL: "http://localhost:8000",
		modelPath:   "Qwen/Qwen3-14B",
		loraPath:    "/home/exx/Desktop/fine-tune/lin_lu_train/outputs/qwen3_final_20260205_163940",
		port:        "8000",
	}
}
// NOTE: Using qwen3_final_20260205_163940 - latest model trained on final dataset
// Training format: [HH:MM] timestamp prefix, NO EOT markers, multi-turn as separate entries
// Model outputs in Chinese (中文)

// LLMStatus represents the current LLM status
type LLMStatus struct {
	Enabled       bool   `json:"enabled"`
	ServerRunning bool   `json:"server_running"`
	ModelPath     string `json:"model_path"`
	LoRAPath      string `json:"lora_path"`
	Port          string `json:"port"`
	BaseURL       string `json:"base_url"`
}

// GetStatus returns the current LLM status (checks if vLLM is actually running)
func (s *LLMService) GetStatus() LLMStatus {
	s.mu.RLock()
	baseURL := s.vllmBaseURL
	modelPath := s.modelPath
	loraPath := s.loraPath
	port := s.port
	s.mu.RUnlock()

	// Actually check if vLLM is running by pinging the models endpoint
	serverRunning := s.checkServerHealth()

	return LLMStatus{
		Enabled:       serverRunning,
		ServerRunning: serverRunning,
		ModelPath:     modelPath,
		LoRAPath:      loraPath,
		Port:          port,
		BaseURL:       baseURL,
	}
}

// checkServerHealth checks if the vLLM server is actually running
func (s *LLMService) checkServerHealth() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(s.vllmBaseURL + "/v1/models")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// StartServer starts the vLLM server
func (s *LLMService) StartServer(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.serverRunning {
		return fmt.Errorf("server already running")
	}

	// Build vLLM command with LoRA adapter
	args := []string{
		"-m", "vllm.entrypoints.openai.api_server",
		"--model", s.modelPath,
		"--port", s.port,
		"--enable-lora",
		"--lora-modules", fmt.Sprintf("linlu=%s", s.loraPath),
		"--max-lora-rank", "64",
		"--trust-remote-code",
		"--dtype", "bfloat16",
		"--gpu-memory-utilization", "0.8",
	}

	s.serverCmd = exec.CommandContext(ctx, "python3", args...)
	s.serverCmd.Stdout = os.Stdout
	s.serverCmd.Stderr = os.Stderr

	if err := s.serverCmd.Start(); err != nil {
		return fmt.Errorf("failed to start vLLM server: %v", err)
	}

	s.serverRunning = true
	s.enabled = true

	// Wait for server to be ready in background
	go s.waitForServer()

	log.Printf("vLLM server started on port %s with model %s", s.port, s.modelPath)
	return nil
}

// waitForServer polls until the server is ready
func (s *LLMService) waitForServer() {
	maxAttempts := 60 // Wait up to 5 minutes
	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(s.vllmBaseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			log.Printf("vLLM server is ready!")
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(5 * time.Second)
	}
	log.Printf("Warning: vLLM server health check timed out")
}

// StopServer stops the vLLM server
func (s *LLMService) StopServer() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.serverRunning || s.serverCmd == nil {
		s.serverRunning = false
		s.enabled = false
		return nil
	}

	if err := s.serverCmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop vLLM server: %v", err)
	}

	s.serverCmd = nil
	s.serverRunning = false
	s.enabled = false

	log.Printf("vLLM server stopped")
	return nil
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	UseLoRA     bool          `json:"use_lora,omitempty"`
}

// ChatMessage represents a single message in a chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// Chat sends a chat completion request to vLLM
func (s *LLMService) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Check if server is actually running by pinging it
	if !s.checkServerHealth() {
		return nil, fmt.Errorf("LLM server is not running")
	}
	
	s.mu.RLock()
	baseURL := s.vllmBaseURL
	s.mu.RUnlock()

	// Prepare the request body
	model := "Qwen/Qwen3-14B"
	if req.UseLoRA {
		model = "linlu" // Use the LoRA adapter name
	}

	body := map[string]interface{}{
		"model":    model,
		"messages": req.Messages,
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	} else {
		body["max_tokens"] = 512
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	} else {
		body["temperature"] = 0.7
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/v1/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vLLM returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &chatResp, nil
}

// RegisterLLMAPI registers LLM-related API endpoints
func (a *App) RegisterLLMAPI(r chi.Router) {
	// Initialize LLM service if not already done
	if a.llmService == nil {
		a.llmService = NewLLMService()
	}

	r.Get("/api/llm/status", a.handleLLMStatus())
	r.Post("/api/llm/toggle", a.handleLLMToggle())
	r.Post("/api/llm/chat", a.handleLLMChat())
	r.Post("/api/llm/config", a.handleLLMConfig())
}

func (a *App) handleLLMStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := a.llmService.GetStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func (a *App) handleLLMToggle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Enable bool `json:"enable"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		var err error
		if req.Enable {
			err = a.llmService.StartServer(r.Context())
		} else {
			err = a.llmService.StopServer()
		}

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		status := a.llmService.GetStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func (a *App) handleLLMChat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if len(req.Messages) == 0 {
			http.Error(w, "messages cannot be empty", http.StatusBadRequest)
			return
		}

		resp, err := a.llmService.Chat(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App) handleLLMConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cfg LLMConfig
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		a.llmService.mu.Lock()
		if cfg.BaseModel != "" {
			a.llmService.modelPath = cfg.BaseModel
		}
		if cfg.LoRAPath != "" {
			a.llmService.loraPath = cfg.LoRAPath
		}
		if cfg.Port != "" {
			a.llmService.port = cfg.Port
			a.llmService.vllmBaseURL = "http://localhost:" + cfg.Port
		}
		a.llmService.mu.Unlock()

		status := a.llmService.GetStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

// handleLLM renders the LLM usage page
func (a *App) handleLLM(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionData := a.getCurrentUser(r)
		loggedIn := sessionData != nil
		isDeveloper := false
		username := ""
		
		if loggedIn {
			username = sessionData.Username
			// Check if user is developer/admin/editor
			isDeveloper = sessionData.Role == "developer" || sessionData.Role == "admin" || sessionData.Role == RoleEditor
		}

		data := map[string]interface{}{
			"LoggedIn":    loggedIn,
			"CurrentUser": username,
			"IsDeveloper": isDeveloper,
			"LLMStatus":   a.llmService.GetStatus(),
		}

		if err := tpl.Execute(w, data); err != nil {
			log.Printf("Error executing llm template: %v", err)
			http.Error(w, "template error", http.StatusInternalServerError)
		}
	}
}

