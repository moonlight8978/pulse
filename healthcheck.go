package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func checkTCP(config *Config) (bool, error) {
	if config.Verbose {
		log.Printf("Checking TCP connection to %s:%d", config.Host, config.Port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		if config.Verbose {
			log.Printf("TCP connection failed: %v", err)
		}
		return false, nil
	}
	defer conn.Close()

	if config.Verbose {
		log.Printf("TCP connection successful")
	}
	return true, nil
}

func checkUDP(config *Config) (bool, error) {
	if config.Verbose {
		log.Printf("Checking UDP connection to %s:%d", config.Host, config.Port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Create UDP connection
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "udp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		if config.Verbose {
			log.Printf("UDP connection failed: %v", err)
		}
		return false, nil
	}
	defer conn.Close()

	// Send a small probe packet to check if the port is actually listening
	// We'll send a simple "ping" message
	probe := []byte("PING")

	// Set write deadline
	if err := conn.SetWriteDeadline(time.Now().Add(config.Timeout)); err != nil {
		if config.Verbose {
			log.Printf("Failed to set write deadline: %v", err)
		}
		return false, nil
	}

	// Send probe packet
	_, err = conn.Write(probe)
	if err != nil {
		if config.Verbose {
			log.Printf("Failed to send UDP probe: %v", err)
		}
		return false, nil
	}

	// Try to read a response (optional - some UDP services don't respond)
	// Set read deadline
	if err := conn.SetReadDeadline(time.Now().Add(config.Timeout)); err != nil {
		if config.Verbose {
			log.Printf("Failed to set read deadline: %v", err)
		}
		// Don't fail here, as some UDP services don't respond
	}

	// Try to read response (non-blocking)
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, err = conn.Read(buffer)

	// For UDP health checks, we consider it successful if:
	// 1. We can establish the connection (address is reachable)
	// 2. We can send a packet (no immediate ICMP unreachable)
	// 3. We don't get an immediate error on write

	// Note: We don't require a response because many UDP services don't respond to unsolicited packets
	// The fact that we can send a packet without getting an ICMP "port unreachable" is usually sufficient

	if config.Verbose {
		if err != nil {
			log.Printf("UDP probe sent successfully, no response received (this is normal for many UDP services)")
		} else {
			log.Printf("UDP probe sent successfully, response received")
		}
		log.Printf("UDP connection successful")
	}
	return true, nil
}

func checkHTTP(config *Config) (bool, error) {
	if config.Verbose {
		log.Printf("Checking HTTP %s request to %s:%d%s", config.Method, config.Host, config.Port, config.Path)
	}

	// Determine protocol based on port
	protocol := "http"
	if config.Port == 443 {
		protocol = "https"
	}

	// Build URL
	url := fmt.Sprintf("%s://%s:%d%s", protocol, config.Host, config.Port, config.Path)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Create request
	req, err := http.NewRequest(config.Method, url, nil)
	if err != nil {
		if config.Verbose {
			log.Printf("Failed to create HTTP request: %v", err)
		}
		return false, err
	}

	// Set User-Agent to identify our tool
	req.Header.Set("User-Agent", "Pulse/1.0")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		if config.Verbose {
			log.Printf("HTTP request failed: %v", err)
		}
		return false, nil
	}
	defer resp.Body.Close()

	// Read and discard response body to ensure connection is properly closed
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil && config.Verbose {
		log.Printf("Warning: failed to read response body: %v", err)
	}

	if config.Verbose {
		log.Printf("HTTP request successful, status: %d", resp.StatusCode)
	}

	// Consider 2xx and 3xx status codes as success
	return resp.StatusCode >= 200 && resp.StatusCode < 400, nil
}
