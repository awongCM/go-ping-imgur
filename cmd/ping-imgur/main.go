package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/awongCM/go-ping-imgur/internal/imgur"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "ping-imgur",
		Short: "Ping Imgur image URLs with simple HTTP requests",
		Long: `Ping Imgur images by sending HTTP HEAD requests.

No Imgur API key is required. This tool is meant for learning and for
checking that your hosted images still respond.`,
	}

	var file string
	var timeout time.Duration

	pingCmd := &cobra.Command{
		Use:   "ping [target...]",
		Short: "Ping one or more Imgur URLs or image IDs",
		Example: strings.TrimSpace(`
ping-imgur ping cvWgXFc
ping-imgur ping https://i.imgur.com/cvWgXFc.jpg
ping-imgur ping -f images.txt`),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			targets, err := collectTargets(args, file)
			if err != nil {
				return err
			}
			if len(targets) == 0 {
				return fmt.Errorf("provide targets as arguments or with --file")
			}

			client := &http.Client{Timeout: timeout}
			okCount := 0

			for _, target := range targets {
				result := imgur.Ping(client, target)
				printResult(result)
				if result.OK() {
					okCount++
				}
			}

			fmt.Printf("\n%d/%d targets responded OK\n", okCount, len(targets))
			if okCount < len(targets) {
				return fmt.Errorf("one or more targets failed")
			}
			return nil
		},
	}

	pingCmd.Flags().StringVarP(&file, "file", "f", "", "file with one URL or image ID per line")
	pingCmd.Flags().DurationVarP(&timeout, "timeout", "t", 15*time.Second, "HTTP request timeout")

	root.AddCommand(pingCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func collectTargets(args []string, file string) ([]string, error) {
	var targets []string
	targets = append(targets, args...)

	if file == "" {
		return targets, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		targets = append(targets, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}

func printResult(result imgur.Result) {
	status := "OK"
	if !result.OK() {
		status = "FAIL"
	}

	line := fmt.Sprintf("[%s] %s", status, result.Target)
	if result.URL != "" && result.URL != result.Target {
		line += fmt.Sprintf(" -> %s", result.URL)
	}
	if result.StatusCode > 0 {
		line += fmt.Sprintf(" (%d)", result.StatusCode)
	}
	line += fmt.Sprintf(" in %s", result.Duration.Round(time.Millisecond))
	if result.Err != nil {
		line += fmt.Sprintf(" — %v", result.Err)
	}

	fmt.Println(line)
}
