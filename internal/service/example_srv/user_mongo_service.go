package example_srv

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go-starter/internal/lib/log"
	do "go-starter/internal/models/do/mongo/example_do"
	"go-starter/internal/models/vo"
	"go-starter/internal/repository/mongo/example_repo"
	"go-starter/utils"
)

type UserMongoService interface {
	IsAvailable() bool
	GetByID(ctx context.Context, id string) (*do.User, error)
	List(ctx context.Context, page, pageSize int) (vo.UserMongoListResp, error)
	Create(ctx context.Context, req vo.CreateUserReq) (*do.User, error)
	Update(ctx context.Context, id string, req vo.UpdateUserReq) (vo.UserMongoIDResp, error)
	Delete(ctx context.Context, id string) (vo.UserMongoIDResp, error)
}

type UserMongoServiceImpl struct {
	repo           example_repo.UserRepository
	contextTimeout time.Duration
}

func NewUserMongoService(repo example_repo.UserRepository, timeout time.Duration) UserMongoService {
	if repo == nil {
		panic("UserRepository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	return &UserMongoServiceImpl{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (s *UserMongoServiceImpl) IsAvailable() bool {
	return s.repo.IsAvailable()
}

func (s *UserMongoServiceImpl) ensureAvailable() error {
	if s.repo.IsAvailable() {
		return nil
	}
	return example_repo.ErrMongoUnavailable
}

func (s *UserMongoServiceImpl) GetByID(ctx context.Context, id string) (user *do.User, err error) {
	log.Logger.Debugf("[UserMongoSrv.GetByID] id[%s] start", id)
	if err = s.ensureAvailable(); err != nil {
		log.Logger.Warnf("[UserMongoSrv.GetByID] mongo unavailable")
		return nil, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Logger.Warnf("[UserMongoSrv.GetByID] id[%s] invalid objectID: %v", id, err)
		return nil, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	user, err = s.repo.GetByID(ctx, objectID)
	if err != nil {
		log.Logger.Errorf("[UserMongoSrv.GetByID] id[%s] repo.GetByID error: %v", id, err)
		return nil, err
	}
	log.Logger.Infof("[UserMongoSrv.GetByID] id[%s] success", id)
	return user, nil
}

func (s *UserMongoServiceImpl) List(ctx context.Context, page, pageSize int) (resp vo.UserMongoListResp, err error) {
	log.Logger.Debugf("[UserMongoSrv.List] page[%d] pageSize[%d] start", page, pageSize)
	if err = s.ensureAvailable(); err != nil {
		log.Logger.Warnf("[UserMongoSrv.List] mongo unavailable")
		return resp, err
	}

	if err = vo.ValidateBaseList(page, pageSize); err != nil {
		log.Logger.Warnf("[UserMongoSrv.List] page[%d] pageSize[%d] validate error: %v", page, pageSize, err)
		return resp, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	offset := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	cond := bson.M{"isDelete": bson.M{"$ne": true}}
	users, total, err := s.repo.GetByCondAndPage(ctx, cond, offset, limit)
	if err != nil {
		log.Logger.Errorf("[UserMongoSrv.List] page[%d] pageSize[%d] repo.GetByCondAndPage error: %v", page, pageSize, err)
		return resp, err
	}
	resp = vo.UserMongoListResp{
		Total: total,
		List:  users,
	}
	log.Logger.Infof("[UserMongoSrv.List] page[%d] pageSize[%d] total[%d] success", page, pageSize, total)
	return resp, nil
}

func (s *UserMongoServiceImpl) Create(ctx context.Context, req vo.CreateUserReq) (user *do.User, err error) {
	log.Logger.Debugf("[UserMongoSrv.Create] email[%s] start", req.Email)
	if err = s.ensureAvailable(); err != nil {
		log.Logger.Warnf("[UserMongoSrv.Create] mongo unavailable")
		return nil, err
	}

	if req.Name == "" || req.Email == "" {
		log.Logger.Warnf("[UserMongoSrv.Create] invalid param")
		return nil, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	user = &do.User{
		Name:     req.Name,
		Email:    req.Email,
		IsDelete: false,
	}
	err = s.repo.Create(ctx, user)
	if err != nil {
		log.Logger.Errorf("[UserMongoSrv.Create] email[%s] repo.Create error: %v", req.Email, err)
		return nil, err
	}
	log.Logger.Infof("[UserMongoSrv.Create] email[%s] id[%s] success", req.Email, user.Id.Hex())
	return user, nil
}

func (s *UserMongoServiceImpl) Update(ctx context.Context, id string, req vo.UpdateUserReq) (resp vo.UserMongoIDResp, err error) {
	log.Logger.Debugf("[UserMongoSrv.Update] id[%s] start", id)
	if err = s.ensureAvailable(); err != nil {
		log.Logger.Warnf("[UserMongoSrv.Update] mongo unavailable")
		return resp, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Logger.Warnf("[UserMongoSrv.Update] id[%s] invalid objectID: %v", id, err)
		return resp, utils.ErrBadParamInput
	}

	updates := bson.M{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if len(updates) == 0 {
		log.Logger.Warnf("[UserMongoSrv.Update] id[%s] no valid updates", id)
		return resp, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	if err = s.repo.UpdateByID(ctx, objectID, updates); err != nil {
		log.Logger.Errorf("[UserMongoSrv.Update] id[%s] repo.UpdateByID error: %v", id, err)
		return resp, err
	}
	resp = vo.UserMongoIDResp{
		ID: id,
	}
	log.Logger.Infof("[UserMongoSrv.Update] id[%s] success", id)
	return resp, nil
}

func (s *UserMongoServiceImpl) Delete(ctx context.Context, id string) (resp vo.UserMongoIDResp, err error) {
	log.Logger.Debugf("[UserMongoSrv.Delete] id[%s] start", id)
	if err = s.ensureAvailable(); err != nil {
		log.Logger.Warnf("[UserMongoSrv.Delete] mongo unavailable")
		return resp, err
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Logger.Warnf("[UserMongoSrv.Delete] id[%s] invalid objectID: %v", id, err)
		return resp, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	if err = s.repo.DeleteByID(ctx, objectID); err != nil {
		log.Logger.Errorf("[UserMongoSrv.Delete] id[%s] repo.DeleteByID error: %v", id, err)
		return resp, err
	}
	resp = vo.UserMongoIDResp{
		ID: id,
	}
	log.Logger.Infof("[UserMongoSrv.Delete] id[%s] success", id)
	return resp, nil
}
