package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

)

type ChatMessage struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionRequest struct {
	Model string `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type chatCompletionResponse struct {
	Choices []struct{
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func GenerateReply(apiKey string, model string, messages []ChatMessage) (string, error){
	reqBody := chatCompletionRequest{Model: model, Messages: messages}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil{
		return "",err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "",err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result chatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "",err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("empty response from groq")
	}
	return result.Choices[0].Message.Content, nil
}
