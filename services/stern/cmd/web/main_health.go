package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var healthcheckClient = http.Client{Timeout: 5 * time.Second}

func runHealthcheck(args []string) int {
	url, err := resolveHealthcheckUrl(args)
	if err != nil {
		slog.Error("Invalid healthcheck arguments", "error", err)
		return 2
	}

	if err := checkHealth(url); err != nil {
		slog.Error("Healthcheck failed", "error", err)
		return 1
	}
	return 0
}

func checkHealth(url string) error {
	response, err := healthcheckClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected status %d", response.StatusCode)
	}
	return nil
}

func resolveHealthcheckUrl(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("expected healthcheck url argument")
	}

	url := strings.TrimSpace(args[0])
	if url == "" {
		return "", fmt.Errorf("healthcheck url is required")
	}
	return url, nil
}
