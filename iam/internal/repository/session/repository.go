package session

import (
	"fmt"

	"github.com/voronovsg/rocket-factory/platform/pkg/cache"
)

const (
	prefixSession        = "session"
	prefixUserSessionSet = "user_session_set"
)

type repository struct {
	rdc cache.RedisClient
}

func NewRepository(rdc cache.RedisClient) *repository {
	return &repository{
		rdc: rdc,
	}
}

func makeSessionKey(sessionUUID string) string {
	return fmt.Sprintf("%s:%s", prefixSession, sessionUUID)
}

func makeUserSessionsSetKey(userUUID string) string {
	return fmt.Sprintf("%s:%s", prefixUserSessionSet, userUUID)
}
