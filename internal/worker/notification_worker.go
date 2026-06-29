package worker

import (
	"context"
	"encoding/json"
	"log"
	"mini-issue/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

const NotificationKey = "notification:queue"

type NotificationProducer struct {
	cli *redis.Client
}

func NewNotificationProducer(cli *redis.Client) *NotificationProducer {
	return &NotificationProducer{cli}
}

func (p *NotificationProducer) PushNotificationTask(ctx context.Context, task model.NotificationTask) error {
	data, _ := json.Marshal(task)
	return p.cli.LPush(ctx, NotificationKey, data).Err()
}

type NotificationWorker struct {
	cli *redis.Client
}

func NewNotificationWorker(cli *redis.Client) *NotificationWorker {
	return &NotificationWorker{cli}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	log.Println("notification worker started")
	for {
		select {
		case <- ctx.Done():
			log.Println("notification worker stopped")
			return
		default:
		}

		result, err := w.cli.BRPop(ctx, 5 * time.Second, NotificationKey).Result()
		if err != nil {
			log.Printf("BRPOP notification queue failed: %v", err)
			time.Sleep(time.Second)
			continue
		}

		rawtask := result[1]
		var task model.NotificationTask
		json.Unmarshal([]byte(rawtask), &task)
		w.handleTask(task)
	}
}

func (w *NotificationWorker) handleTask(task model.NotificationTask) {
	switch task.Type {
	case model.Notification_IssueCreated:
		log.Printf(
			"send notification: user=%d issue=%d title=%s",
			task.UserID,
			task.IssueID,
			task.Title,
		)
	default:
		log.Printf("unknown notification type: %s", task.Type)
	}
}