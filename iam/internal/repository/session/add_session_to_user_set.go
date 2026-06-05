package session

import "context"

func (r *repository) AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	sessionKey := makeSessionKey(sessionUUID)
	userSessionsSetKey := makeUserSessionsSetKey(userUUID)

	err := r.rdc.SAdd(ctx, userSessionsSetKey, sessionKey)
	if err != nil {
		return err
	}

	return nil
}
