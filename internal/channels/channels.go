package channels

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ipfred/cc-hooks-notify/internal/config"
)

type Sender interface {
	Send(title, body, event string) error
}

type baseSender struct {
	client  *http.Client
	config  config.ChannelConfig
	webhook string
}

func New(name string, cfg config.ChannelConfig) (Sender, error) {
	base := baseSender{
		client:  &http.Client{Timeout: 10 * time.Second},
		config:  cfg,
		webhook: strings.TrimSpace(cfg.Webhook),
	}
	switch strings.ToLower(name) {
	case "dingtalk":
		return dingtalkSender{baseSender: base}, nil
	case "feishu":
		return feishuSender{baseSender: base}, nil
	default:
		return nil, fmt.Errorf("unsupported channel %q", name)
	}
}

func (s baseSender) enabled() bool {
	return s.webhook != "" && (s.config.Enabled == nil || *s.config.Enabled || s.config.Webhook != "")
}

func (s baseSender) post(rawURL string, payload any) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	textBytes, _ := io.ReadAll(resp.Body)
	text := string(textBytes)
	if resp.StatusCode != http.StatusOK {
		return text, fmt.Errorf("http %d: %s", resp.StatusCode, text)
	}
	return text, nil
}

func emoji(event string) string {
	switch event {
	case "permission":
		return "🔐"
	case "idle":
		return "💤"
	case "stop":
		return "✅"
	default:
		return "🔔"
	}
}

type dingtalkSender struct{ baseSender }

func (s dingtalkSender) Send(title, body, event string) error {
	if !s.enabled() {
		return errors.New("dingtalk webhook is empty or disabled")
	}
	useMarkdown := true
	if s.config.Markdown != nil {
		useMarkdown = *s.config.Markdown
	}
	var payload any
	if useMarkdown {
		payload = map[string]any{
			"msgtype": "markdown",
			"markdown": map[string]string{
				"title": title,
				"text":  fmt.Sprintf("#### %s %s\n\n%s", emoji(event), title, body),
			},
		}
	} else {
		payload = map[string]any{
			"msgtype": "text",
			"text": map[string]string{
				"content": fmt.Sprintf("%s %s\n\n%s", emoji(event), title, body),
			},
		}
	}
	text, err := s.post(s.signedURL(), payload)
	if err != nil {
		return err
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil
	}
	if code, ok := result["errcode"].(float64); ok && code == 0 {
		return nil
	}
	return fmt.Errorf("dingtalk returned %s", text)
}

func (s dingtalkSender) signedURL() string {
	if strings.TrimSpace(s.config.Secret) == "" {
		return s.webhook
	}
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	mac := hmac.New(sha256.New, []byte(s.config.Secret))
	_, _ = mac.Write([]byte(timestamp + "\n" + s.config.Secret))
	sign := url.QueryEscape(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	separator := "&"
	if !strings.Contains(s.webhook, "?") {
		separator = "?"
	}
	return s.webhook + separator + "timestamp=" + timestamp + "&sign=" + sign
}

type feishuSender struct{ baseSender }

func (s feishuSender) Send(title, body, event string) error {
	if !s.enabled() {
		return errors.New("feishu webhook is empty or disabled")
	}
	useCard := true
	if s.config.Card != nil {
		useCard = *s.config.Card
	}
	payload := map[string]any{}
	if useCard {
		payload = map[string]any{
			"msg_type": "interactive",
			"card": map[string]any{
				"config": map[string]any{"wide_screen_mode": true},
				"header": map[string]any{
					"template": feishuColor(event),
					"title": map[string]string{
						"content": fmt.Sprintf("%s %s", emoji(event), title),
						"tag":     "plain_text",
					},
				},
				"elements": []map[string]any{
					{
						"tag": "div",
						"text": map[string]string{
							"content": body,
							"tag":     "lark_md",
						},
					},
				},
			},
		}
	} else {
		payload = map[string]any{
			"msg_type": "text",
			"content": map[string]string{
				"text": fmt.Sprintf("%s %s\n\n%s", emoji(event), title, body),
			},
		}
	}
	if strings.TrimSpace(s.config.Secret) != "" {
		timestamp := fmt.Sprintf("%d", time.Now().Unix())
		payload["timestamp"] = timestamp
		payload["sign"] = s.sign(timestamp)
	}
	text, err := s.post(s.webhook, payload)
	if err != nil {
		return err
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil
	}
	if code, ok := result["code"].(float64); ok && code == 0 {
		return nil
	}
	return fmt.Errorf("feishu returned %s", text)
}

func (s feishuSender) sign(timestamp string) string {
	mac := hmac.New(sha256.New, []byte(s.config.Secret))
	_, _ = mac.Write([]byte(timestamp + "\n" + s.config.Secret))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func feishuColor(event string) string {
	switch event {
	case "permission":
		return "orange"
	case "stop":
		return "green"
	default:
		return "blue"
	}
}
