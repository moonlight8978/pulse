package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Mode    string // "tcp", "udp", or "http"
	Host    string
	Port    int
	Path    string
	Method  string
	Timeout time.Duration
	Silent  bool
	Verbose bool
}

func main() {
	config := parseFlags()

	if config.Verbose {
		log.Printf("Pulse healthcheck tool starting...")
		log.Printf("Mode: %s", config.Mode)
		log.Printf("Target: %s:%d", config.Host, config.Port)
		if config.Mode == "http" {
			log.Printf("Path: %s", config.Path)
			log.Printf("Method: %s", config.Method)
		}
	}

	var success bool
	var err error

	switch config.Mode {
	case "tcp":
		success, err = checkTCP(config)
	case "udp":
		success, err = checkUDP(config)
	case "http":
		success, err = checkHTTP(config)
	default:
		log.Fatalf("Invalid mode: %s. Use 'tcp', 'udp', or 'http'", config.Mode)
	}

	if err != nil {
		if config.Verbose {
			log.Printf("Error: %v", err)
		}
		os.Exit(1)
	}

	if success {
		if !config.Silent {
			fmt.Println("OK")
		}
		os.Exit(0)
	} else {
		if !config.Silent {
			fmt.Println("FAIL")
		}
		os.Exit(1)
	}
}

func parseFlags() *Config {
	config := &Config{}

	// Mode flag
	mode := flag.String("mode", "tcp", "Mode: 'tcp', 'udp', or 'http'")

	// Common flags
	host := flag.String("host", "localhost", "Host to check")
	port := flag.Int("port", 80, "Port to check")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout for the check")
	silent := flag.Bool("silent", false, "Silent mode (no output)")
	verbose := flag.Bool("verbose", false, "Verbose mode (debug output)")

	// HTTP specific flags
	path := flag.String("path", "/", "HTTP path (for http mode)")
	method := flag.String("method", "GET", "HTTP method (for http mode)")

	flag.Parse()

	config.Mode = *mode
	config.Host = *host
	config.Port = *port
	config.Timeout = *timeout
	config.Silent = *silent
	config.Verbose = *verbose
	config.Path = *path
	config.Method = strings.ToUpper(*method)

	// Validate mode
	if config.Mode != "tcp" && config.Mode != "udp" && config.Mode != "http" {
		log.Fatalf("Invalid mode: %s. Use 'tcp', 'udp', or 'http'", config.Mode)
	}

	// Validate port
	if config.Port < 1 || config.Port > 65535 {
		log.Fatalf("Invalid port: %d. Must be between 1 and 65535", config.Port)
	}

	// Validate HTTP method
	if config.Mode == "http" {
		validMethods := map[string]bool{
			"GET": true, "POST": true, "PUT": true, "DELETE": true,
			"HEAD": true, "OPTIONS": true, "PATCH": true,
		}
		if !validMethods[config.Method] {
			log.Fatalf("Invalid HTTP method: %s", config.Method)
		}
	}

	return config
}
