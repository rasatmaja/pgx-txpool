package integration

import (
	"sync"

	"github.com/rasatmaja/pgx-txpool/tests/integration/model"
)

type users struct {
	usrs []model.User
	mx   sync.Mutex
}

func (u *users) add(users ...model.User) {
	u.mx.Lock()
	defer u.mx.Unlock()
	u.usrs = append(u.usrs, users...)
}

func (u *users) get() []model.User {
	u.mx.Lock()
	defer u.mx.Unlock()
	return u.usrs
}
