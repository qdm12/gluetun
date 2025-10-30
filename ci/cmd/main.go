package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " <command>")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	switch os.Args[1] {
	case "run-container":
		err := runContainerTest(ctx)
		stop()
		if err != nil {
			fmt.Println("❌ Test failed:", err)
			os.Exit(1)
		}
		fmt.Println("✅ Test completed successfully.")
	default:
		fmt.Println("Unknown command:", os.Args[1])
		stop()
		os.Exit(1)
	}
}

func runContainerTest(ctx context.Context) error {
	secrets, err := readSecrets(ctx)
	if err != nil {
		return fmt.Errorf("reading secrets: %w", err)
	}

	const timeout = 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("creating Docker client: %w", err)
	}
	defer client.Close()

	config := &container.Config{
		Image:       "qmcgaw/gluetun",
		StopTimeout: ptrTo(3),
		Env: []string{
			"VPN_SERVICE_PROVIDER=mullvad",
			"VPN_TYPE=wireguard",
			"LOG_LEVEL=debug",
			"SERVER_COUNTRIES=USA",
			"WIREGUARD_PRIVATE_KEY=" + secrets.mullvadWireguardPrivateKey,
			"WIREGUARD_ADDRESSES=" + secrets.mullvadWireguardAddress,
		},
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
		CapAdd:     []string{"NET_ADMIN", "NET_RAW"},
	}
	networkConfig := (*network.NetworkingConfig)(nil)
	platform := (*v1.Platform)(nil)
	const containerName = "" // auto-generated name

	response, err := client.ContainerCreate(ctx, config, hostConfig, networkConfig, platform, containerName)
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}
	for _, warning := range response.Warnings {
		fmt.Println("Warning during container creation:", warning)
	}
	containerID := response.ID
	defer stopContainer(client, containerID)

	beforeStartTime := time.Now()

	err = client.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("starting container: %w", err)
	}

	return waitForLogLine(ctx, client, containerID, beforeStartTime)
}

func ptrTo[T any](v T) *T { return &v }

type secrets struct {
	mullvadWireguardPrivateKey string
	mullvadWireguardAddress    string
}

func readSecrets(ctx context.Context) (secrets, error) {
	expectedSecrets := [...]string{
		"Mullvad Wireguard private key",
		"Mullvad Wireguard address",
	}

	scanner := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0, len(expectedSecrets))

	for i := range expectedSecrets {
		fmt.Println("🤫 reading", expectedSecrets[i], "from Stdin...")
		if !scanner.Scan() {
			break
		}
		lines = append(lines, strings.TrimSpace(scanner.Text()))
		if ctx.Err() != nil {
			return secrets{}, ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return secrets{}, fmt.Errorf("reading secrets from stdin: %w", err)
	}

	if len(lines) < len(expectedSecrets) {
		return secrets{}, fmt.Errorf("expected %d secrets via Stdin, but only received %d",
			len(expectedSecrets), len(lines))
	}
	for i, line := range lines {
		if line == "" {
			return secrets{}, fmt.Errorf("secret on line %d/%d was empty", i+1, len(lines))
		}
	}

	return secrets{
		mullvadWireguardPrivateKey: lines[0],
		mullvadWireguardAddress:    lines[1],
	}, nil
}

func stopContainer(client *client.Client, containerID string) {
	const stopTimeout = 5 * time.Second // must be higher than 3s, see above [container.Config]'s StopTimeout field
	stopCtx, stopCancel := context.WithTimeout(context.Background(), stopTimeout)
	defer stopCancel()

	err := client.ContainerStop(stopCtx, containerID, container.StopOptions{})
	if err != nil {
		fmt.Println("failed to stop container:", err)
	}
}

var successRegexp = regexp.MustCompile(`^.+Public IP address is .+$`)

func waitForLogLine(ctx context.Context, client *client.Client, containerID string,
	beforeStartTime time.Time,
) error {
	logOptions := container.LogsOptions{
		ShowStdout: true,
		Follow:     true,
		Since:      beforeStartTime.Format(time.RFC3339Nano),
	}

	reader, err := client.ContainerLogs(ctx, containerID, logOptions)
	if err != nil {
		return fmt.Errorf("error getting container logs: %w", err)
	}
	defer reader.Close()

	var linesSeen []string
	scanner := bufio.NewScanner(reader)
	for ctx.Err() == nil {
		if scanner.Scan() {
			line := scanner.Text()
			if len(line) > 8 { // remove Docker log prefix
				line = line[8:]
			}
			linesSeen = append(linesSeen, line)
			if successRegexp.MatchString(line) {
				fmt.Println("✅ Success line logged")
				return nil
			}
			continue
		}
		err := scanner.Err()
		if err != nil && err != io.EOF {
			logSeenLines(linesSeen)
			return fmt.Errorf("reading log stream: %w", err)
		}

		// The scanner is either done or cannot read because of EOF
		fmt.Println("The log scanner stopped")
		logSeenLines(linesSeen)

		// Check if the container is still running
		inspect, err := client.ContainerInspect(ctx, containerID)
		if err != nil {
			return fmt.Errorf("inspecting container: %w", err)
		}
		if !inspect.State.Running {
			return fmt.Errorf("container stopped unexpectedly while waiting for log line. Exit code: %d", inspect.State.ExitCode)
		}
	}

	return ctx.Err()
}

func logSeenLines(lines []string) {
	fmt.Println("Logs seen so far:")
	for _, line := range lines {
		fmt.Println("  " + line)
	}
}
