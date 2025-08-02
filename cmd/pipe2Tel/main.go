package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const maxMessageLength = 4080

func main() {
	botToken := flag.String("bot_token", "", "Telegram bot token")
	chatID := flag.String("chat_id", "", "Telegram chat ID")
	restricted := flag.Bool("rs", false, "Enable restricted mode")
	msg := flag.String("msg", "", "Message to send (provide text directly or a file path)")

	flag.Parse()

	if *botToken == "" || *chatID == "" {
		fmt.Println("Error: bot_token and chat_id are required")
		usageGuide()
		os.Exit(1)
	}

	// Determine input
	var (
		message  string
		isFile   bool
		filePath string
	)

	if *msg != "" {
		if info, err := os.Stat(*msg); err == nil && !info.IsDir() {
			filePath = *msg
			isFile = true
		} else {
			message = *msg
		}
	} else {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, os.Stdin); err != nil {
			fmt.Println("Error reading stdin:", err)
			usageGuide()
			os.Exit(1)
		}
		message = buf.String()
	}

	// Choose text or document
	if isFile || len(message) > maxMessageLength {
		if !isFile {
			tmpFile, err := os.CreateTemp("", "msg-*.txt")
			if err != nil {
				fmt.Println("Error creating temp file:", err)
				return
			}
			defer os.Remove(tmpFile.Name())
			tmpFile.WriteString(message)
			tmpFile.Close()
			filePath = tmpFile.Name()
		}
		sendDocument(*botToken, *chatID, filePath, *restricted)
	} else {
		esc := escapeMarkdownV2(message)
		sendMessage(*botToken, *chatID, esc, *restricted)
	}
}

// escapeMarkdownV2 escapes all MarkdownV2 special characters.
// See: https://core.telegram.org/bots/api#markdownv2-style
func escapeMarkdownV2(input string) string {
	// Escape backslash first
	specialChars := []string{"\\", "[", "]", "(", ")", "~", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, c := range specialChars {
		input = strings.ReplaceAll(input, c, "\\"+c)
	}
	return input
}

func sendMessage(botToken, chatID, message string, restricted bool) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", message)
	data.Set("parse_mode", "MarkdownV2")
	if restricted {
		data.Set("disable_web_page_preview", "true")
		data.Set("protect_content", "true")
	}
	resp, err := http.PostForm(apiURL, data)
	handleResponse(resp, err)
}

func sendDocument(botToken, chatID, path string, restricted bool) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", botToken)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("chat_id", chatID)
	if restricted {
		writer.WriteField("disable_web_page_preview", "true")
		writer.WriteField("protect_content", "true")
	}

	part, err := writer.CreateFormFile("document", filepath.Base(path))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	io.Copy(part, file)
	writer.Close()

	req, err := http.NewRequest("POST", apiURL, &body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	handleResponse(resp, err)
}

func handleResponse(resp *http.Response, err error) {
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: API responded with %d\n", resp.StatusCode)
		var out bytes.Buffer
		if json.Indent(&out, bodyBytes, "", "  ") == nil {
			fmt.Println("Response JSON:")
			fmt.Println(out.String())
		} else {
			fmt.Println(string(bodyBytes))
		}
	} else {
		fmt.Println("Sent successfully!")
	}
}

func usageGuide() {
	fmt.Println("Usage:")
	fmt.Println("  pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-rs] [-msg=<TEXT OR FILE_PATH>]")
	fmt.Println("  echo \"text\" | pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID>[-rs]")
}
