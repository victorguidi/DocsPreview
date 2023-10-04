package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type PrivateLLM struct {
	Model      string
	APIKey     string
	APIVersion string
	Resource   string
}

func NewPrivateLLM() *PrivateLLM {
	return &PrivateLLM{
		Model:      os.Getenv("DEPLOYMENT_NAME"),
		APIKey:     os.Getenv("OPENAI_API_KEY"),
		APIVersion: os.Getenv("OPENAI_API_VERSION"),
		Resource:   os.Getenv("RESOURCE"),
	}
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response represents the JSON response from the OpenAI API.
type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// PromptForBulletPoint generates a response for the given prompt and type.
func (llm *PrivateLLM) PromptForBulletPoint(prompt string, typeStr string, c *websocket.Conn) (string, error) {
	chat := []Message{
		llm.createSystemMessage(typeStr),
		llm.createUserMessage(prompt),
	}

	url := fmt.Sprintf("https://%s/openai/deployments/%s/chat/completions?api-version=%s", llm.Resource, llm.Model, llm.APIVersion)
	headers := map[string]string{
		"Content-Type": "application/json",
		"api-key":      llm.APIKey,
	}
	bodyData := map[string][]Message{
		"messages": chat,
	}

	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}

	response, err := postJSON(url, headers, jsonBody)
	if err != nil {
		return "", err
	}

	var parsedJSON Response
	if err := json.Unmarshal(response, &parsedJSON); err != nil {
		return "", err
	}

	history := parsedJSON.Choices[0].Message.Content
	chat = append(chat, Message{Role: "assistant", Content: history})

	fmt.Println(chat[len(chat)-1].Content)

	type Response struct {
		Response string `json:"response"`
	}

	resp := Response{
		Response: chat[len(chat)-1].Content,
	}

	if err = c.WriteJSON(resp); err != nil {
		log.Println(err)
	}

	return chat[len(chat)-1].Content, nil
}

func (llm *PrivateLLM) createUserMessage(message string) Message {
	return Message{Role: "user", Content: message}
}

func (llm *PrivateLLM) createSystemMessage(typeStr string) Message {
	return Message{Role: "assistant", Content: fmt.Sprintf("Você é uma IA especializada em resumir os prompts que o usuário te passar e retornar em formato markdown, mas formatando os pontos principais de acordo com o que inferir no tipo: %s", typeStr)}
}

func postJSON(url string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody.Bytes(), nil
}
