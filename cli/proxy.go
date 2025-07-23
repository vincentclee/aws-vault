package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/vincentclee/aws-vault/v7/server"
)

func ConfigureProxyCommand(app *kingpin.Application) {
	stop := false

	cmd := app.Command("proxy", "Start a proxy for the ec2 instance role server locally.").
		Alias("server").
		Hidden()

	cmd.Flag("stop", "Stop the proxy").
		BoolVar(&stop)

	cmd.Action(func(*kingpin.ParseContext) error {
		if stop {
			server.StopProxy()
			return nil
		}
		handleSigTerm()
		return server.StartProxy()
	})
}

func handleSigTerm() {
	// shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Shutdown()
		os.Exit(1)
	}()
}
