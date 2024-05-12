package mail

import (
	"testing"

	"github.com/laugart7272/inscrips/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGMail(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewEmailSender(config.EMailSenderName, config.EMailSenderAddress, config.EMailSenderPassword)

	subject := "A test EMail"
	content := `
	<h1>Hello IPS</h1>
	<p>This is a message from <a href="http://ipls.ao">Instituto Politecnico de Saurimo</a></p>
	`
	to := []string{"leolau7272@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
