package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/virzz/splunk-go"
	"github.com/virzz/vlog"
)

var rootCmd = &cobra.Command{
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	SilenceErrors:     true,
	SilenceUsage:      true,
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(userHome, ".config", "splunk.auth")
	buf, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	auth := splunk.Config{}
	err = json.Unmarshal(buf, &auth)
	if err != nil {
		return err
	}
	return splunk.Init(cmd.Context(), &auth)
}

var authCmd = &cobra.Command{
	Use: "auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		endpoint, _ := f.GetString("endpoint")
		username, _ := f.GetString("username")
		password, _ := f.GetString("password")

		auth := splunk.Config{
			Host:     endpoint,
			Username: username,
			Password: password,
		}
		err := splunk.Init(cmd.Context(), &auth)
		if err != nil {
			return err
		}
		if !splunk.AuthCheck() {
			return errors.New("Host/Username/Password Invalid")
		}
		userHome, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath := filepath.Join(userHome, ".config")
		os.MkdirAll(configPath, 0755)
		buf, err := json.Marshal(&auth)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(configPath, "splunk.auth"), buf, 0755)
	},
}

func init() {
	authCmd.Flags().StringP("endpoint", "e", "", "Splunk Endpoint")
	authCmd.Flags().StringP("username", "u", "", "Splunk Username")
	authCmd.Flags().StringP("password", "p", "", "Splunk Password")
}

func main() {
	for _, cmd := range rootCmd.Commands() {
		cmd.PersistentPreRunE = persistentPreRunE
	}
	rootCmd.AddCommand(authCmd)
	err := rootCmd.Execute()
	if err != nil {
		vlog.Error(err.Error())
	}
}
