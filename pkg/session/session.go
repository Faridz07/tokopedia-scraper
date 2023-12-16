package session

import (
	"context"
	"errors"
	"tokopedia-scraper/infrastructure/log"

	"github.com/google/uuid"
)

const (
	SessionKey = "sess_key"
)

type Session struct {
	Ctx        context.Context
	Error      error
	Logger     log.Log
	Message    string
	ErrMessage string
	Xid        uuid.UUID
}

func NewSession(logger log.Log,
) *Session {
	return &Session{
		Ctx:    context.Background(),
		Xid:    uuid.New(),
		Logger: logger,
	}
}

func (sess *Session) SetContext(ctx context.Context) {
	sess.Ctx = ctx
}

func (sess *Session) SetError(message string, err error) {
	if message != "" && err == nil {
		err = errors.New(message)
	}

	sess.ErrMessage = message
	sess.Error = err

	sys := &log.Model{
		XID:   sess.GetXID(),
		Error: err.Error(),
	}

	sess.Logger.Error(sys, message)
}

func (sess *Session) SetErrorWithData(message string, err error) {

	if message != "" && err == nil {
		err = errors.New(message)
	}

	sess.ErrMessage = message
	sess.Error = err

	sys := &log.Model{
		XID:   sess.GetXID(),
		Error: err.Error(),
	}

	sess.Logger.Error(sys, message)
}

func (sess *Session) SetInfo(message string) {
	sess.Message = message

	sys := &log.Model{
		XID:   sess.GetXID(),
		Error: "",
	}

	sess.Logger.Info(sys, message)
}

func GetSession(ctx context.Context) *Session {
	sess, _ := ctx.Value(SessionKey).(*Session)
	return sess
}

func (sess *Session) GetXID() uuid.UUID {
	if sess.Xid == uuid.Nil {
		sess.Xid = uuid.New()
	}
	return sess.Xid
}

func (sess *Session) GetContext() context.Context {
	return sess.Ctx
}

func (sess *Session) GetError() error {
	return sess.Error
}
