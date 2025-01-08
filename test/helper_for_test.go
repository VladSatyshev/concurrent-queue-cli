package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/VladSatyshev/concurrent-queue-cli/client"
	"github.com/VladSatyshev/concurrent-queue-cli/test/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const buildPath = "./build/cli"

type cmdExecutor struct {
	binPath string
}

func newCmdExecutor(path string) *cmdExecutor {
	return &cmdExecutor{binPath: path}
}

func (c *cmdExecutor) Execute(t *testing.T, argsStr string) []byte {
	cmd := exec.Command(c.binPath, strings.Split(argsStr, " ")...)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to execute build cli cmd: %s", err.Error())
	}
	cmd.Dir = wd
	res, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to execute command: %s", err.Error())
	}
	return res
}

func prepareConfig(queuesCfg []config.QueueConfig) error {
	cfg := config.Config{
		Server: config.ServerConfig{
			Port:              ":8000",
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
	cmd := exec.Command("./build/server", "-config", config.TestConfigPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	time.Sleep(time.Millisecond * 200) // wait for server to start
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

func buildCli(t *testing.T, path string) {
	cmd := exec.Command("go", "build", "-o", path, "./../cmd/main.go")
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to build cli: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("CLI binary not found at %s after build", path)
	}
}

func configureEnvironment(t *testing.T, queuesCfg []config.QueueConfig) (*cmdExecutor, func()) {
	err := prepareConfig(queuesCfg)
	if err != nil {
		t.Fatalf("failed to configure env: %v", err)
	}

	cmd := startServer(t)

	exec := newCmdExecutor(buildPath)

	return exec, func() {
		stopServer(t, cmd)
	}
}

func assertEqualQueues(t *testing.T, expectedQueues []client.Queue, actualQueues []client.Queue) {
	assert.Equal(t, len(expectedQueues), len(actualQueues))
	for _, aq := range actualQueues {
		for _, eq := range expectedQueues {
			if aq.Name == eq.Name {
				assert.Equal(t, eq, aq)
			}
		}
	}
}

func assertEqualMessages(t *testing.T, expectedMessages []map[string]interface{}, actualMessages []map[string]interface{}) {
	assert.Equal(t, len(expectedMessages), len(actualMessages))

	for _, actualMessage := range actualMessages {
		found := false
		for _, expectedMessage := range expectedMessages {
			for _, av := range actualMessage {
				for _, ev := range expectedMessage {
					if av == ev {
						found = true
					}
				}
			}
		}
		if !found {
			t.Error("not equal messages")
		}
	}
}
