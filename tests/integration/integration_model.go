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

type transactions struct {
	trxs []model.Transaction
	mx   sync.Mutex
}

func (t *transactions) add(trxs ...model.Transaction) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.trxs = append(t.trxs, trxs...)
}

func (t *transactions) get() []model.Transaction {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.trxs
}

type transfer struct {
	trfs []model.TransactionTransfer
	mx   sync.Mutex
}

func (t *transfer) add(trfs ...model.TransactionTransfer) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.trfs = append(t.trfs, trfs...)
}
func (t *transfer) get() []model.TransactionTransfer {
	t.mx.Lock()
	defer t.mx.Unlock()
	return t.trfs
}
