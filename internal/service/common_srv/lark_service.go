package common_srv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/SisyphusSQ/golib/models/dto/lark_dto"
	"github.com/SisyphusSQ/golib/utils"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	"go-starter/config"
	"go-starter/internal/lib/log"
)

type LarkService interface {
	SetBotUrl(botUrl string)
	SendBotMsg(title string, content []lark_dto.Content) error

	SendLarkMsg(ctx context.Context, req lark_dto.LarkMsgReq) (resp lark_dto.LarkMsgResp, err error)

	//loop()
	botCall(request lark_dto.BotMsg) ([]byte, error)
	sendMsg(ctx context.Context, contacts []string, msg string) ([]*larkim.CreateMessageResp, error)
}

type LarkSrvImpl struct {
	botUrl    string
	botClient *http.Client
	client    *lark.Client
}

func NewLarkService(c config.Config) LarkService {
	s := &LarkSrvImpl{
		client: lark.NewClient(c.Lark.AppID, c.Lark.AppSecret,
			lark.WithLogger(log.LarkLogger), lark.WithLogLevel(larkcore.LogLevelDebug)),

		botClient: http.DefaultClient,
	}
	return s
}

/*
func (s *LarkSrvImpl) loop() {
	fn := func() error {
		kv, err := s.configRepo.GetByKey(context.Background(), xxx.BotURL)
		if err != nil {
			return err
		}
		s.botUrl = kv.V
		return nil
	}

	err := fn()
	if err != nil {
		log.Logger.Errorf("[LarkService.loop] get bot url error: %v", err)
	}

	tk := time.NewTicker(5 * time.Minute)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			err = fn()
			if err != nil {
				log.Logger.Errorf("[LarkService.loop] get bot url error: %v", err)
			}
		}
	}
}
*/

func (s *LarkSrvImpl) SetBotUrl(botUrl string) {
	s.botUrl = botUrl
}

func (s *LarkSrvImpl) SendBotMsg(title string, content []lark_dto.Content) error {
	log.Logger.Debugf("[LarkService.SendBotMsg] title[%s] start", title)
	msg := lark_dto.NewBotMsg(title, content)
	rsp, err := s.botCall(msg)
	if err != nil {
		log.Logger.Errorf("[LarkService.SendBotMsg] botCall error: %v", err)
		return err
	}

	var resp lark_dto.BotMsgResp
	err = json.Unmarshal(rsp, &resp)
	if err != nil {
		log.Logger.Errorf("[LarkService.SendBotMsg] json.Unmarshal BotMsgResp error: %v", err)
		return err
	}

	if resp.Code != lark_dto.Success {
		log.Logger.Warnf("[LarkService.SendBotMsg] title[%s] resp code[%d] not success", title, resp.Code)
	}
	log.Logger.Infof("[LarkService.SendBotMsg] title[%s] success", title)
	return nil
}

func (s *LarkSrvImpl) botCall(request lark_dto.BotMsg) ([]byte, error) {
	if s.botUrl == "" {
		return []byte{}, errors.New("bot url not be config")
	}

	reqBody, err := json.MarshalIndent(request, "", "	   ")
	if err != nil {
		return nil, err
	}
	log.Logger.Debugf("[LarkService.botCall] request body:\n%s", string(reqBody))

	req, err := http.NewRequest("POST", s.botUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.botClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Logger.Debugf("[LarkService.botCall] raw response body:\n%s", string(respBody))
	return respBody, nil
}

func (s *LarkSrvImpl) SendLarkMsg(ctx context.Context, req lark_dto.LarkMsgReq) (resp lark_dto.LarkMsgResp, err error) {
	log.Logger.Debugf("[LarkService.SendLarkMsg] contacts[%v] start", req.Contacts)
	resp.Resp, err = s.sendMsg(ctx, req.Contacts, req.Message)
	if err != nil {
		log.Logger.Errorf("[LarkService.SendLarkMsg] sendMsg error: %v", err)
		return
	}
	log.Logger.Infof("[LarkService.SendLarkMsg] contacts[%v] success", req.Contacts)
	return
}

func (s *LarkSrvImpl) sendMsg(ctx context.Context, contacts []string, msg string) ([]*larkim.CreateMessageResp, error) {
	var (
		err   error
		resps = make([]*larkim.CreateMessageResp, 0)
	)

	for _, c := range contacts {
		uuid := utils.UUID()
		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType("email").
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(c).
				Uuid(uuid).
				MsgType("interactive").
				Content(msg).
				Build()).
			Build()

		resp, createErr := s.client.Im.Message.Create(ctx, req)
		if createErr != nil {
			log.Logger.Errorf("[LarkService.sendMsg] contact[%s] Create error: %v", c, createErr)
			err = createErr
			continue
		}

		if !resp.Success() {
			err = errors.New(resp.ErrorResp())
		}

		resps = append(resps, resp)
	}

	return resps, err
}
