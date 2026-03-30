package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"aurora-adminui/internal/config"
)

func sendBootstrapTokenToTelegram(ctx context.Context, cfg *config.Config, message string) error {
	if cfg == nil {
		return fmt.Errorf("admin config is nil")
	}
	botToken := strings.TrimSpace(cfg.Admin.BootstrapTelegramBotToken)
	chatID := strings.TrimSpace(cfg.Admin.BootstrapTelegramChatID)
	if botToken == "" || chatID == "" {
		return fmt.Errorf("telegram bootstrap config is incomplete")
	}

	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", message)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.telegram.org/bot"+botToken+"/sendMessage",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("telegram send failed status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}
