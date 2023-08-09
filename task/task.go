package task

import (
	"context"
	"echoapptpl/dbo"
	"github.com/axengine/utils/log"
)

type Task struct {
	db   *dbo.DBO
	exit chan struct{}
}

func New(db *dbo.DBO) *Task {
	return &Task{db: db,
		exit: make(chan struct{}, 1)}
}

func (t *Task) Start(ctx context.Context) {
	go t.start(ctx)
}

func (t *Task) Wait() {
	<-t.exit
}

func (t *Task) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Logger.Info("Train Task exit...")
			t.exit <- struct{}{}
			return
		default:
			// dosomething
		}
	}
}
