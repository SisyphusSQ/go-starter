package example_srv

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-starter/config"
	gormv2 "go-starter/internal/lib/gorm"
	"go-starter/internal/lib/log"
	redisv9 "go-starter/internal/lib/redis"
	do "go-starter/internal/models/do/mysql/example_do"
	"go-starter/internal/models/vo"
	"go-starter/internal/repository/mysql/example_repo"
	"go-starter/utils"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserService interface {
	GetByID(ctx context.Context, id int64) (do.User, error)
	List(ctx context.Context, page, pageSize int) (vo.UserListResp, error)
	Create(ctx context.Context, req vo.CreateUserReq) (do.User, error)
	Login(ctx context.Context, req vo.LoginReq) (vo.LoginResp, error)
	Logout(ctx context.Context, userID int64) (vo.UserIDResp, error)
	Update(ctx context.Context, id int64, req vo.UpdateUserReq) (vo.UserIDResp, error)
	Delete(ctx context.Context, id int64) (vo.UserIDResp, error)
}

type UserServiceImpl struct {
	repo           example_repo.UserRepository
	contextTimeout time.Duration
	config         config.Config
	cache          *redisv9.Client
}

func NewUserService(
	repo example_repo.UserRepository,
	timeout time.Duration,
	c config.Config,
	cache *redisv9.Client,
) UserService {
	if repo == nil {
		panic("UserRepository is nil")
	}
	if timeout == 0 {
		panic("Timeout is empty")
	}
	if cache == nil {
		panic("Redis cache is nil")
	}
	return &UserServiceImpl{
		repo:           repo,
		contextTimeout: timeout,
		config:         c,
		cache:          cache,
	}
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id int64) (user do.User, err error) {
	log.Logger.Debugf("[UserSrv.GetByID] id[%d] start", id)
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	user, err = s.repo.GetByID(ctx, id)
	if err != nil {
		log.Logger.Errorf("[UserSrv.GetByID] id[%d] repo.GetByID error: %v", id, err)
		return do.User{}, err
	}
	log.Logger.Infof("[UserSrv.GetByID] id[%d] success, cost: [%d] ms", id, time.Since(start).Milliseconds())
	return user, nil
}

func (s *UserServiceImpl) List(ctx context.Context, page, pageSize int) (resp vo.UserListResp, err error) {
	log.Logger.Debugf("[UserSrv.List] page[%d] pageSize[%d] start", page, pageSize)
	start := time.Now()
	if err = vo.ValidateBaseList(page, pageSize); err != nil {
		log.Logger.Warnf("[UserSrv.List] page[%d] pageSize[%d] validate error: %v", page, pageSize, err)
		return resp, err
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	conditions := &gormv2.DBConditions{
		NeedCount: true,
		Order:     "id DESC",
		Limit:     pageSize,
		Offset:    (page - 1) * pageSize,
	}
	users, err := s.repo.ListByConditions(ctx, conditions)
	if err != nil {
		log.Logger.Errorf("[UserSrv.List] page[%d] pageSize[%d] repo.ListByConditions error: %v", page, pageSize, err)
		return resp, err
	}
	resp = vo.UserListResp{
		Total: conditions.Count,
		List:  users,
	}
	log.Logger.Infof("[UserSrv.List] page[%d] pageSize[%d] total[%d] listLen[%d] success, cost: [%d] ms",
		page, pageSize, conditions.Count, len(users), time.Since(start).Milliseconds())
	return resp, nil
}

func (s *UserServiceImpl) Create(ctx context.Context, req vo.CreateUserReq) (user do.User, err error) {
	log.Logger.Debugf("[UserSrv.Create] email[%s] start", req.Email)
	if req.Name == "" || req.Email == "" || req.Password == "" {
		log.Logger.Warnf("[UserSrv.Create] email[%s] invalid param, name/email/password empty", req.Email)
		return do.User{}, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Logger.Errorf("[UserSrv.Create] email[%s] bcrypt.GenerateFromPassword error: %v", req.Email, err)
		return do.User{}, err
	}
	user = do.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(password),
	}
	err = s.repo.Create(ctx, &user)
	if err != nil {
		log.Logger.Errorf("[UserSrv.Create] email[%s] repo.Create error: %v", req.Email, err)
		return do.User{}, err
	}
	log.Logger.Infof("[UserSrv.Create] email[%s] id[%d] success", req.Email, user.ID)
	return user, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req vo.LoginReq) (resp vo.LoginResp, err error) {
	log.Logger.Debugf("[UserSrv.Login] email[%s] start", req.Email)
	if req.Email == "" || req.Password == "" {
		log.Logger.Warnf("[UserSrv.Login] invalid param")
		return resp, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Logger.Warnf("[UserSrv.Login] email[%s] user not found", req.Email)
			return resp, ErrInvalidCredentials
		}
		log.Logger.Errorf("[UserSrv.Login] email[%s] repo.GetByEmail error: %v", req.Email, err)
		return resp, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Logger.Warnf("[UserSrv.Login] email[%s] invalid credentials", req.Email)
		return resp, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Email, s.config.Key.JWT)
	if err != nil {
		log.Logger.Errorf("[UserSrv.Login] email[%s] GenerateToken error: %v", req.Email, err)
		return resp, err
	}

	cacheKey := fmt.Sprintf("jwt:user:%d", user.ID)
	expire := time.Duration(s.config.Key.JWT.Expire) * time.Second
	if err = s.cache.Set(ctx, cacheKey, token, expire).Err(); err != nil {
		log.Logger.Errorf("[UserSrv.Login] userID[%d] cache.Set error: %v", user.ID, err)
		return resp, err
	}

	resp = vo.LoginResp{
		Token:  token,
		Expire: s.config.Key.JWT.Expire,
	}
	log.Logger.Infof("[UserSrv.Login] email[%s] userID[%d] success", req.Email, user.ID)
	return resp, nil
}

func (s *UserServiceImpl) Logout(ctx context.Context, userID int64) (resp vo.UserIDResp, err error) {
	log.Logger.Debugf("[UserSrv.Logout] userID[%d] start", userID)
	if userID <= 0 {
		log.Logger.Warnf("[UserSrv.Logout] userID[%d] invalid param", userID)
		return resp, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	cacheKey := fmt.Sprintf("jwt:user:%d", userID)
	if err = s.cache.Del(ctx, cacheKey).Err(); err != nil {
		log.Logger.Errorf("[UserSrv.Logout] userID[%d] cache.Del error: %v", userID, err)
		return resp, err
	}

	resp = vo.UserIDResp{
		ID: userID,
	}
	log.Logger.Infof("[UserSrv.Logout] userID[%d] success", userID)
	return resp, nil
}

func (s *UserServiceImpl) Update(ctx context.Context, id int64, req vo.UpdateUserReq) (resp vo.UserIDResp, err error) {
	log.Logger.Debugf("[UserSrv.Update] id[%d] start", id)
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if len(updates) == 0 {
		log.Logger.Warnf("[UserSrv.Update] id[%d] no valid updates", id)
		return resp, utils.ErrBadParamInput
	}

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	if err = s.repo.UpdateByID(ctx, id, updates); err != nil {
		log.Logger.Errorf("[UserSrv.Update] id[%d] repo.UpdateByID error: %v", id, err)
		return resp, err
	}
	resp = vo.UserIDResp{
		ID: id,
	}
	log.Logger.Infof("[UserSrv.Update] id[%d] success", id)
	return resp, nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int64) (resp vo.UserIDResp, err error) {
	log.Logger.Debugf("[UserSrv.Delete] id[%d] start", id)
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	if err = s.repo.DeleteByID(ctx, id); err != nil {
		log.Logger.Errorf("[UserSrv.Delete] id[%d] repo.DeleteByID error: %v", id, err)
		return resp, err
	}
	resp = vo.UserIDResp{
		ID: id,
	}
	log.Logger.Infof("[UserSrv.Delete] id[%d] success", id)
	return resp, nil
}
