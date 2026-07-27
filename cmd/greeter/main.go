package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arnoldcplrms/occasion-greeter-go/internal/greeter"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file loaded (continuing): %v", err)
	}

	fmt.Println("🎉 Starting Birthday & Anniversary Greeter...")

	fmt.Println("📂 Loading occasions data...")
	occasionsData, err := greeter.LoadOccasionsData()
	if err != nil {
		log.Fatalf("✗ Error: %v", err)
	}
	fmt.Printf("✓ Loaded %d occasions\n", len(occasionsData))

	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		loc = time.FixedZone("PHT", 8*60*60)
	}
	now := time.Now().In(loc)
	if now.Day() == 1 {
		fmt.Println("📅 First day of the month — sending monthly summary...")
		occasions, month := greeter.FindOccasionsThisMonth(occasionsData)
		summary := greeter.GenerateMonthlySummaryGreeting(occasions, month)
		fmt.Printf("📝 Generated monthly summary for %s\n", month)
		result := greeter.SendMonthlySummaryEmail(summary, month)
		if !result.Success {
			log.Fatalf("✗ Failed to send monthly summary email: %s", result.Error)
		}
		fmt.Printf("✓ Monthly summary email sent successfully (ID: %s)\n", result.MessageID)
		return
	}

	fmt.Println("🔍 Searching for occasions today and tomorrow...")
	occasionsToday := greeter.FindOccasionsToday(occasionsData)
	occasionsTomorrow := greeter.FindOccasionsTomorrow(occasionsData)

	if len(occasionsToday) == 0 && len(occasionsTomorrow) == 0 {
		fmt.Println("✓ No occasions today or tomorrow")
		return
	}

	manifest, err := greeter.GetManifest()
	if err != nil {
		log.Fatalf("✗ Error: %v", err)
	}

	fmt.Printf("✓ Found %d occasion(s) today\n", len(occasionsToday))
	fmt.Printf("✓ Found %d occasion(s) tomorrow\n", len(occasionsTomorrow))

	for _, occasion := range occasionsToday {
		var logMessage string
		switch {
		case occasion.Type == greeter.OccasionTypeBirthday && occasion.Person != nil:
			logMessage = fmt.Sprintf("🎂 Processing birthday for %s", occasion.Person.Name)
		case occasion.Type == greeter.OccasionTypeAnniversary && occasion.Couple != nil:
			logMessage = fmt.Sprintf("💍 Processing anniversary for Team %s", occasion.Couple.LastName)
		}
		fmt.Printf("\n📧 %s\n", logMessage)

		greeting := greeter.GenerateGreeting(occasion, nil)
		fmt.Printf("📝 Generated greeting: %s\n", greeting)

		result := greeter.SendOccasionEmail(occasion, greeting, manifest, greeter.SendOptions{})
		if !result.Success {
			log.Fatalf("✗ Failed to send email: %s", result.Error)
		}
		fmt.Printf("✓ Email sent successfully (ID: %s)\n", result.MessageID)
	}

	if len(occasionsTomorrow) > 0 {
		fmt.Printf("\n📧 ⏰ Processing combined reminder for %d occasions tomorrow\n", len(occasionsTomorrow))
		bulkReminder := greeter.GenerateBulkReminderGreeting(occasionsTomorrow)
		fmt.Printf("📝 Generated combined reminder:\n%s\n", bulkReminder)

		result := greeter.SendBulkReminderEmail(bulkReminder)
		if !result.Success {
			log.Fatalf("✗ Failed to send combined reminder email: %s", result.Error)
		}
		fmt.Printf("✓ Combined reminder email sent successfully (ID: %s)\n", result.MessageID)
	}

	fmt.Println("\n✓ Birthday & Anniversary Greeter completed")
}
