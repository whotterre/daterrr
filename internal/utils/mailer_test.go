package utils

import (
	"path/filepath"
	"testing"
  "fmt"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	configPath := filepath.Join("../../")
	config, err := LoadConfig(configPath)
  fmt.Print(config)
	require.NoError(t, err)

	testSender := NewGmailSender("Daterrr", config.EmailAddr, config.GmailKey)

	subject := "Test email"
	content := `
	<h1> This is a test email </h1>
	<p> Nothing to worry about here! </p>
	`

	to := []string{"darklinuxusr@gmail.com"}


  err = testSender.SendEmail(subject, content, to, nil, nil, nil)
  require.NoError(t, err)
}
