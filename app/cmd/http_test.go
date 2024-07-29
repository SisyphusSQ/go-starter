package cmd

import (
	"fmt"
	"os"
	"testing"

	"go-starter/config"
)

func TestHttp(t *testing.T) {
	cmd := rootCmd
	cmd.SetOut(os.Stdout)
	cmd.SetArgs([]string{"http"})

	fmt.Println(os.Getwd())
	config.SetConfigFile("../../config/config.yml")

	if err := cmd.Execute(); err != nil {
		t.Fatalf("cmd.Execute() failed: %v", err)
	}
}
