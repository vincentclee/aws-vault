package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/vincentclee/aws-vault/v8/prompt"
	"github.com/vincentclee/aws-vault/v8/vault"
	"github.com/vincentclee/keyring/v2"
)

type AddCommandInput struct {
	ProfileName string
	FromEnv     bool
	AddConfig   bool
}

func ConfigureAddCommand(app *kingpin.Application, a *AwsVault) {
	input := AddCommandInput{}

	cmd := app.Command("add", "Add credentials to the secure keystore.")

	cmd.Arg("profile", "Name of the profile").
		Required().
		StringVar(&input.ProfileName)

	cmd.Flag("env", "Read the credentials from the environment (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY)").
		BoolVar(&input.FromEnv)

	cmd.Flag("add-config", "Add a profile to ~/.aws/config if one doesn't exist").
		Default("true").
		BoolVar(&input.AddConfig)

	cmd.Action(func(c *kingpin.ParseContext) error {
		keyring, err := a.Keyring()
		if err != nil {
			return err
		}
		awsConfigFile, err := a.AwsConfigFile()
		if err != nil {
			return err
		}
		err = AddCommand(input, keyring, awsConfigFile)
		app.FatalIfError(err, "add")
		return nil
	})
}

func AddCommand(input AddCommandInput, keyring keyring.Keyring, awsConfigFile *vault.ConfigFile) error {
	var accessKeyID, secretKey string

	p, _ := awsConfigFile.ProfileSection(input.ProfileName)
	if p.SourceProfile != "" {
		return fmt.Errorf("Your profile has a source_profile of %s, adding credentials to %s won't have any effect",
			p.SourceProfile, input.ProfileName)
	}

	if input.FromEnv {
		if accessKeyID = os.Getenv("AWS_ACCESS_KEY_ID"); accessKeyID == "" {
			return fmt.Errorf("Missing value for AWS_ACCESS_KEY_ID")
		}
		if secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY"); secretKey == "" {
			return fmt.Errorf("Missing value for AWS_SECRET_ACCESS_KEY")
		}
	} else {
		var err error
		if accessKeyID, err = prompt.TerminalPrompt("Enter Access Key ID: "); err != nil {
			return err
		}
		if secretKey, err = prompt.TerminalSecretPrompt("Enter Secret Access Key: "); err != nil {
			return err
		}
	}

	creds := aws.Credentials{AccessKeyID: accessKeyID, SecretAccessKey: secretKey}

	ckr := &vault.CredentialKeyring{Keyring: keyring}
	if err := ckr.Set(input.ProfileName, creds); err != nil {
		return err
	}

	fmt.Printf("Added credentials to profile %q in vault\n", input.ProfileName)

	sk := &vault.SessionKeyring{Keyring: keyring}
	if n, _ := sk.RemoveForProfile(input.ProfileName); n > 0 {
		fmt.Printf("Deleted %d existing sessions.\n", n)
	}

	if _, hasProfile := awsConfigFile.ProfileSection(input.ProfileName); !hasProfile {
		if input.AddConfig {
			newProfileSection := vault.ProfileSection{
				Name: input.ProfileName,
			}
			log.Printf("Adding profile %s to config at %s", input.ProfileName, awsConfigFile.Path)
			if err := awsConfigFile.Add(newProfileSection); err != nil {
				return fmt.Errorf("Error adding profile: %w", err)
			}
		}
	}

	return nil
}
