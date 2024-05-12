package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/laugart7272/inscrips/db/sqlc"
	"github.com/laugart7272/inscrips/util"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Email string `json:"name"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshall task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueue task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal((task.Payload()), &payload); err != nil {
		return fmt.Errorf("failed unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUserEmail(ctx, payload.Email)
	if err != nil {
		return fmt.Errorf("user not exists %w", asynq.SkipRetry)
	}

	verifyemail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Name:       user.Name,
		LastName:   user.LastName,
		Email:      user.Email,
		UserID:     int32(user.ID),
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	///Send Email
	config, _ := util.LoadConfig(".")
	subject := "Bemvindo ao Inscripção do Evento - Instituto Politecnico do Saurimo"
	webServer := config.WebServer
	verifyURL := fmt.Sprintf("%s/verify_email?email_id=%d&secret_code=%s", webServer, verifyemail.ID, verifyemail.SecretCode)
	content := fmt.Sprintf(`Oi %s %s,<br/>
	Muito Obrigado por se registrar em nosso site!<br/>
	Clique <a href="%s">aqui</a> para verificar seu endereço de e-mail<br/>
	`, user.Name, user.LastName, verifyURL)
	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")

	return nil
}
