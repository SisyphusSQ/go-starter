package common_srv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/SisyphusSQ/golib/models/do/base_do"
	"github.com/SisyphusSQ/golib/models/dto/lark_dto"
	"github.com/SisyphusSQ/golib/utils"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	"go-starter/config"
	"go-starter/internal/lib/log"
	"go-starter/internal/repository/mysql/my_common"
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

	logRepo    my_common.LarkMsgLogRepository
	configRepo my_common.ConfigKVRepository
}

func NewLarkService(c config.Config, logRepo my_common.LarkMsgLogRepository, configRepo my_common.ConfigKVRepository) LarkService {
	if logRepo == nil {
		panic("logRepo is nil")
	}

	if configRepo == nil {
		panic("configKVRepo is nil")
	}

	s := &LarkSrvImpl{
		client: lark.NewClient(c.Lark.AppID, c.Lark.AppSecret,
			lark.WithLogger(log.LarkLogger), lark.WithLogLevel(larkcore.LogLevelDebug)),

		botClient:  http.DefaultClient,
		logRepo:    logRepo,
		configRepo: configRepo,
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
		log.Logger.Error("Lark service before loop, got bot url failed, err: %v", err)
	}

	tk := time.NewTicker(5 * time.Minute)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			err = fn()
			if err != nil {
				log.Logger.Error("Lark service in loop, got bot url failed, err: %v", err)
			}
		}
	}
}
*/

func (s *LarkSrvImpl) SetBotUrl(botUrl string) {
	s.botUrl = botUrl
}

func (s *LarkSrvImpl) SendBotMsg(title string, content []lark_dto.Content) error {
	msg := lark_dto.NewBotMsg(title, content)
	rsp, err := s.botCall(msg)
	if err != nil {
		log.Logger.Errorf("SendBotMsg err:%v", err)
		return err
	}

	req, _ := json.MarshalIndent(msg, "", "    ")
	botLog := base_do.LarkMsgLog{
		SendMessage: string(req),
		Response:    string(rsp),
		Status:      base_do.Status(lark_dto.Success),
	}

	var resp lark_dto.BotMsgResp
	err = json.Unmarshal(rsp, &resp)
	if err != nil {
		log.Logger.Errorf("SendBotMsg lark.BotMsgResp json.Unmarshal err:%v", err)
		return err
	}

	if resp.Code != lark_dto.Success {
		botLog.Status = base_do.Failure
	}

	if s.logRepo != nil {
		err = s.logRepo.CreateRecord(context.Background(), botLog)
		if err != nil {
			log.Logger.Errorf("SendBotMsg logRepo.CreateRecord err:%v", err)
			return err
		}
	}
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
	log.Logger.Infof("[LarkService] reqeust body:\n%s", string(reqBody))

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
	log.Logger.Infof("[LarkService] raw response body:\n%s", string(respBody))
	return respBody, nil
}

func (s *LarkSrvImpl) SendLarkMsg(ctx context.Context, req lark_dto.LarkMsgReq) (resp lark_dto.LarkMsgResp, err error) {
	resp.Resp, err = s.sendMsg(ctx, req.Contacts, req.Message)
	if err != nil {
		log.Logger.Errorf("[LarkService] SendLarkMsg err: %v", err)
		return
	}

	return
}

func (s *LarkSrvImpl) sendMsg(ctx context.Context, contacts []string, msg string) ([]*larkim.CreateMessageResp, error) {
	var (
		err     error
		resps   = make([]*larkim.CreateMessageResp, 0)
		records = make([]base_do.LarkMsgLog, 0)
	)

	for _, c := range contacts {
		record := base_do.LarkMsgLog{
			UUID:        utils.UUID(),
			SendMessage: msg,
			Response:    "",
			Status:      base_do.Success,
			Contacts:    c,
		}

		req := larkim.NewCreateMessageReqBuilder().
			ReceiveIdType("email").
			Body(larkim.NewCreateMessageReqBodyBuilder().
				ReceiveId(c).
				Uuid(record.UUID).
				MsgType("interactive").
				Content(record.SendMessage).
				Build()).
			Build()

		resp, err := s.client.Im.Message.Create(ctx, req)
		if err != nil {
			log.Logger.Errorf("[LarkService] sendMsg err: %v", err)
			record.Status = base_do.Failure
			records = append(records, record)
			continue
		}

		if !resp.Success() {
			record.Status = base_do.Failure
			record.Response = resp.ErrorResp()
			err = errors.New(resp.ErrorResp())
		} else {
			respStr, e := json.Marshal(resp)
			if e == nil {
				record.Response = string(respStr)
			}
		}

		if resp.Data != nil {
			record.MessageID = *resp.Data.MessageId
		}

		records = append(records, record)
		resps = append(resps, resp)
	}

	go func() {
		for _, record := range records {
			errR := s.logRepo.CreateRecord(context.Background(), record)
			if errR != nil {
				log.Logger.Errorf("[LarkService] write records to db err: %v", errR)
			}
		}
	}()

	return resps, err
}
