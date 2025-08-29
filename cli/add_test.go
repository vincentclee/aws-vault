package cli

import (
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func ExampleAddCommand() {
	f, err := os.CreateTemp("", "aws-config")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	_ = os.Setenv("AWS_CONFIG_FILE", f.Name())
	_ = os.Setenv("AWS_ACCESS_KEY_ID", "llamas")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "rock")
	_ = os.Setenv("AWS_VAULT_BACKEND", "file")
	_ = os.Setenv("AWS_VAULT_FILE_PASSPHRASE", "password")

	defer func() { _ = os.Unsetenv("AWS_ACCESS_KEY_ID") }()
	defer func() { _ = os.Unsetenv("AWS_SECRET_ACCESS_KEY") }()
	defer func() { _ = os.Unsetenv("AWS_VAULT_BACKEND") }()
	defer func() { _ = os.Unsetenv("AWS_VAULT_FILE_PASSPHRASE") }()

	app := kingpin.New(`aws-vault`, ``)
	ConfigureAddCommand(app, ConfigureGlobals(app))
	kingpin.MustParse(app.Parse([]string{"add", "--debug", "--env", "foo"}))

	// Output:
	// Added credentials to profile "foo" in vault
}
