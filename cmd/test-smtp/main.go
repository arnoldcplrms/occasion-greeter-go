package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file loaded (continuing)")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	if portStr == "" {
		portStr = "587"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587
	}
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")

	fmt.Println("Testing SMTP connection...")
	fmt.Printf("SMTP Host: %s\n", host)
	fmt.Printf("SMTP Port: %d\n", port)
	fmt.Printf("SMTP User: %s\n", user)
	fmt.Printf("Password length: %d characters\n", len(pass))
	fmt.Println()

	dialer := gomail.NewDialer(host, port, user, pass)
	closer, err := dialer.Dial()
	if err != nil {
		fmt.Println("❌ SMTP connection error:")
		fmt.Println(err)
		return
	}
	defer closer.Close()
	fmt.Println("✅ SMTP connection verified successfully!")
}
