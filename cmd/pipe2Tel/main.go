package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	botToken := flag.String("bot_token", "", "Telegram bot token")
	chatID := flag.String("chat_id", "", "Telegram chat ID")
	restricted := flag.Bool("rs", false, "Enable restricted mode")
	msg := flag.String("msg", "", "Message to send (provide text directly or a file path)")

	flag.Parse()

	// Validate required flags
	if *botToken == "" || *chatID == "" {
		fmt.Println("Error: bot_token and chat_id are required")
		usageGuide()
		os.Exit(1)
	}

	// Determine the message source
	var message string
	if *msg != "" {
		// Check if msg is a file or direct text
		if fileContent, err := readFileContent(*msg); err == nil {
			message = fileContent
		} else {
			message = *msg
		}
	} else {
		// Read input from stdin
		var buffer bytes.Buffer
		_, err := io.Copy(&buffer, os.Stdin)
		if err != nil {
			fmt.Println("Error reading input message from stdin:", err)
			usageGuide()
			os.Exit(1)
		}
		message = buffer.String()
	}

	escapedMessage := escapeMarkdownV2(message)

	sendMessage(*botToken, *chatID, escapedMessage, *restricted)
}

func readFileContent(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("provided path is a directory, not a file")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// escapes special characters for MarkdownV2
func escapeMarkdownV2(input string) string {
	// Perhaps you found this useful:
	// https://github.com/telegraf/telegraf/issues/1242
	// Note that the "\\" should be the first item.
	specialChars := []string{"!", "#", "+", "-", "=", "{", "}", ".", "&"}
	for _, char := range specialChars {
		input = strings.ReplaceAll(input, char, "\\"+char)
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
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	// Check and print response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Telegram API responded with status %d\n", resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		var formattedBody bytes.Buffer
		if err := json.Indent(&formattedBody, body, "", "  "); err == nil {
			fmt.Println("Response JSON:")
			fmt.Println(formattedBody.String())
		}
	} else {
		fmt.Println("Message sent successfully!")
	}
}

// Usage guide
func usageGuide() {
	fmt.Println("Usage:")
	fmt.Println("I>   pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-restricted] [-msg=<TEXT OR FILE_PATH>]")
	fmt.Println("II>  echo \"sth\" | pipe2Tel -bot_token=<TOKEN> -chat_id=<CHAT_ID> [-restricted]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -bot_token    The Telegram bot token (required)")
	fmt.Println("  -chat_id      The Telegram chat ID (required)")
	fmt.Println("  -rs           Optional flag to enable restricted mode (no web page preview, no notification)")
	fmt.Println("  -msg          The message to send. If this is a file path, the file content is used as the message.")
	fmt.Println("                If it's not a file path, it's treated as direct text.")
	fmt.Println()
	fmt.Println("If no -msg flag is provided, the program will read the message from stdin(II).")
}
