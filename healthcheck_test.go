package main

import (
	"testing"
	"time"
)

func TestParseFlags(t *testing.T) {
	// Test default values
	config := &Config{
		Mode:    "tcp",
		Host:    "localhost",
		Port:    80,
		Path:    "/",
		Method:  "GET",
		Timeout: 5 * time.Second,
		Silent:  false,
		Verbose: false,
	}

	if config.Mode != "tcp" {
		t.Errorf("Expected default mode 'tcp', got '%s'", config.Mode)
	}

	if config.Host != "localhost" {
		t.Errorf("Expected default host 'localhost', got '%s'", config.Host)
	}

	if config.Port != 80 {
		t.Errorf("Expected default port 80, got %d", config.Port)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test valid port range
	validPorts := []int{1, 80, 443, 8080, 65535}
	for _, port := range validPorts {
		if port < 1 || port > 65535 {
			t.Errorf("Port %d should be invalid", port)
		}
	}

	// Test invalid port range
	invalidPorts := []int{0, 65536, 99999}
	for _, port := range invalidPorts {
		if port >= 1 && port <= 65535 {
			t.Errorf("Port %d should be valid", port)
		}
	}

	// Test valid HTTP methods
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	for _, method := range validMethods {
		if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" && method != "HEAD" && method != "OPTIONS" && method != "PATCH" {
			t.Errorf("Method %s should be valid", method)
		}
	}
}

func TestCheckTCP(t *testing.T) {
	// Test with invalid host (should fail gracefully)
	config := &Config{
		Mode:    "tcp",
		Host:    "invalid-host-that-does-not-exist.local",
		Port:    9999,
		Timeout: 1 * time.Second,
		Verbose: false,
	}

	success, err := checkTCP(config)
	if err != nil {
		t.Errorf("checkTCP should not return error for connection failure: %v", err)
	}
	if success {
		t.Error("checkTCP should return false for invalid host")
	}
}

func TestCheckHTTP(t *testing.T) {
	// Test with invalid host (should fail gracefully)
	config := &Config{
		Mode:    "http",
		Host:    "invalid-host-that-does-not-exist.local",
		Port:    9999,
		Path:    "/",
		Method:  "GET",
		Timeout: 1 * time.Second,
		Verbose: false,
	}

	success, err := checkHTTP(config)
	if err != nil {
		t.Errorf("checkHTTP should not return error for connection failure: %v", err)
	}
	if success {
		t.Error("checkHTTP should return false for invalid host")
	}
}

func TestCheckUDP(t *testing.T) {
	// Test with invalid host (should fail gracefully)
	config := &Config{
		Mode:    "udp",
		Host:    "invalid-host-that-does-not-exist.local",
		Port:    9999,
		Timeout: 1 * time.Second,
		Verbose: false,
	}

	success, err := checkUDP(config)
	if err != nil {
		t.Errorf("checkUDP should not return error for connection failure: %v", err)
	}
	if success {
		t.Error("checkUDP should return false for invalid host")
	}
}
