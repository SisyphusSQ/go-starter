package common_srv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/SisyphusSQ/go-starter/v2/internal/lib/log"
	"github.com/SisyphusSQ/go-starter/v2/internal/models/dto/lark_dto"
	"github.com/SisyphusSQ/go-starter/v2/utils"
)

type LarkService interface {
	SetBotURL(botURL string)
	SendBotMessage(ctx context.Context, title string, content []lark_dto.LarkContent) error
	SendMessage(ctx context.Context, req lark_dto.LarkMessageRequest) (lark_dto.LarkMessageResponse, error)
}

type LarkServiceImpl struct {
	botURL    string
	botClient *http.Client
	client    *lark.Client
}

func NewLarkService(c config.Config) LarkService {
	return &LarkServiceImpl{
		client: lark.NewClient(
			c.Lark.AppID,
			c.Lark.AppSecret,
			lark.WithLogger(log.LarkLogger),
			lark.WithLogLevel(larkcore.LogLevelWarn),
			lark.WithHttpClient(&http.Client{Timeout: c.Lark.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}),
		),
		botClient: &http.Client{Timeout: c.Lark.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }},
	}
}

func (s *LarkServiceImpl) SetBotURL(botURL string) {
	s.botURL = strings.TrimSpace(botURL)
}

func (s *LarkServiceImpl) SendBotMessage(ctx context.Context, title string, content []lark_dto.LarkContent) error {
	if s.botURL == "" {
		return errors.New("lark bot URL is not configured")
	}

	reqBody, err := json.Marshal(lark_dto.NewBotMessage(title, content))
	if err != nil {
		return fmt.Errorf("marshal lark bot message: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.botURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create lark bot request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.botClient.Do(req)
	if err != nil {
		return errors.New("send lark bot request failed")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read lark bot response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("lark bot returned HTTP %d", resp.StatusCode)
	}

	var result lark_dto.BotMessageResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("decode lark bot response: %w", err)
	}
	if result.Code != 0 {
		return fmt.Errorf("lark bot returned code %d: %s", result.Code, result.Msg)
	}
	return nil
}

func (s *LarkServiceImpl) SendMessage(ctx context.Context, req lark_dto.LarkMessageRequest) (lark_dto.LarkMessageResponse, error) {
	responses := make([]*larkim.CreateMessageResp, 0, len(req.Contacts))
	errs := make([]error, 0)
	for _, contact := range req.Contacts {
		request := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType("email").
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(contact).
				Uuid(utils.UUID()).
				MsgType("interactive").
				Content(req.Message).
				Build()).
			Build()

		resp, err := s.client.Im.Message.Create(ctx, request)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if !resp.Success() {
			errs = append(errs, errors.New(resp.ErrorResp()))
			continue
		}
		responses = append(responses, resp)
	}
	return lark_dto.LarkMessageResponse{Responses: responses}, errors.Join(errs...)
}
