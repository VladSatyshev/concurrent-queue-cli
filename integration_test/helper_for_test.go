package test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/VladSatyshev/concurrent-queue-cli/integration_test/config"
	"gopkg.in/yaml.v2"
)

func prepareConfig(queuesCfg []config.QueueConfig) error {
	cfg := config.Config{
		Server: config.ServerConfig{
			Port:              "8000",
			Mode:              "Development",
			Timeout:           5 * time.Second,
			CtxDefaultTimeout: 10 * time.Second,
		},
		Logger: config.LoggerConfig{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
			Level:             "info",
		},
		Queues: queuesCfg,
	}

	yamlCfg, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(config.TestConfigPath, yamlCfg, 0644)
	if err != nil {
		return err
	}

	return nil
}

func startServer(t *testing.T) *exec.Cmd {
	cmd := exec.Command("go", "run", "../concurrent-queue/cmd/main.go", "-config", config.TestConfigPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	return cmd
}

func stopServer(t *testing.T, cmd *exec.Cmd) {
	err := cmd.Process.Kill()
	if err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}
}

func buildCli(t *testing.T) *exec.Cmd {
	cmd := exec.Command("go", "build", "-o", "./build/cli", "../../cmd/main.go")
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	return cmd
}

func configureEnvironment(t *testing.T, queuesCfg []config.QueueConfig) func() {
	err := prepareConfig(queuesCfg)
	if err != nil {
		t.Fatalf("failed to configure env: %v", err)
	}

	cliCmd := buildCli(t)

	cmd := startServer(t)

	return func() {
		stopServer(t, cmd)
	}
}
