package greeter

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Recipients   []string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
}

type SendOptions struct {
	IsReminder bool
}

type SendResult struct {
	Success   bool
	MessageID string
	Error     string
}

func LoadEmailConfig() (*EmailConfig, error) {
	recipients := []string{}
	for _, r := range strings.Split(os.Getenv("RECIPIENT_EMAIL"), ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			recipients = append(recipients, r)
		}
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil || port == 0 {
		port = 587
	}

	cfg := &EmailConfig{
		Recipients:   recipients,
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	}

	if len(cfg.Recipients) == 0 {
		return nil, fmt.Errorf("missing email configuration: recipients")
	}
	if cfg.SMTPHost == "" {
		return nil, fmt.Errorf("missing email configuration: smtpHost")
	}
	if cfg.SMTPUser == "" {
		return nil, fmt.Errorf("missing email configuration: smtpUser")
	}
	if cfg.SMTPPassword == "" {
		return nil, fmt.Errorf("missing email configuration: smtpPassword")
	}

	return cfg, nil
}

func (c *EmailConfig) newDialer() *gomail.Dialer {
	return gomail.NewDialer(c.SMTPHost, c.SMTPPort, c.SMTPUser, c.SMTPPassword)
}

func inferImageExtension(url string) string {
	cleanURL := strings.ToLower(strings.SplitN(url, "?", 2)[0])
	switch {
	case strings.HasSuffix(cleanURL, ".png"):
		return "png"
	case strings.HasSuffix(cleanURL, ".gif"):
		return "gif"
	case strings.HasSuffix(cleanURL, ".webp"):
		return "webp"
	default:
		return "jpg"
	}
}

func generateEmailHTML(greeting, profilePictureURL string, hasImageAttachment bool, missingPhotoHint string) string {
	viewURL := html.EscapeString(profilePictureURL)
	safeGreeting := html.EscapeString(greeting)
	safeHint := ""
	if missingPhotoHint != "" {
		safeHint = html.EscapeString(missingPhotoHint)
	}

	var mediaSection string
	switch {
	case hasImageAttachment:
		mediaSection = fmt.Sprintf(`<div class="section"><img class="profile-img" src="%s" alt="Profile Picture" /><p class="hint">The profile image is attached to this email for download.</p></div>`, viewURL)
	case safeHint != "":
		mediaSection = fmt.Sprintf(`<div class="section"><p class="hint">%s</p></div>`, safeHint)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <style>
      body { font-family: Arial, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; }
      .container { background-color: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); padding: 30px; max-width: 600px; margin: 0 auto; }
      .greeting { font-size: 18px; line-height: 1.6; color: #333; white-space: pre-wrap; word-wrap: break-word; padding: 16px; background: #f9f9f9; border-radius: 6px; }
      .section { margin: 24px 0; }
      .profile-img { max-width: 100%%; border-radius: 8px; display: block; margin-bottom: 8px; }
      .footer { text-align: center; font-size: 12px; color: #999; margin-top: 20px; border-top: 1px solid #eee; padding-top: 10px; }
      .hint { font-size: 12px; color: #666; margin-top: 4px; }
    </style>
  </head>
  <body>
    <div class="container">

      <div class="section">
        <div class="greeting" id="greetingText">%s</div>
      </div>

      %s

      <div class="footer">
        <p>Sent by Birthday &amp; Anniversary Greeter 🎉</p>
      </div>
    </div>
  </body>
</html>`, safeGreeting, mediaSection)
}

type occasionEmailMeta struct {
	subject                string
	fallbackProfilePicture string
}

func getOccasionEmailMeta(o Occasion, options SendOptions) occasionEmailMeta {
	prefix := ""
	if options.IsReminder {
		prefix = "Upcoming: "
	}

	if o.Type == OccasionTypeBirthday && o.Person != nil {
		return occasionEmailMeta{
			subject:                fmt.Sprintf("%s%s's Birthday", prefix, o.Person.Name),
			fallbackProfilePicture: o.Person.ProfilePicture,
		}
	}

	if o.Type == OccasionTypeAnniversary && o.Couple != nil {
		return occasionEmailMeta{
			subject:                fmt.Sprintf("%sTeam %s Wedding Anniversary", prefix, o.Couple.LastName),
			fallbackProfilePicture: o.Couple.ProfilePicture,
		}
	}

	return occasionEmailMeta{}
}

func resolvePhotoState(o Occasion, lookup PhotoLookupResult, fallbackProfilePicture string) (string, string) {
	if lookup.Status == LookupStatusFound {
		return lookup.URL, ""
	}

	if lookup.Status == LookupStatusNotConfigured && fallbackProfilePicture != "" {
		return NormalizeManifestPhotoURL(fallbackProfilePicture), ""
	}

	if lookup.AttemptedFileName != "" {
		fmt.Printf("⚠️  Photo not found for %s: %s\n", o.Type, lookup.AttemptedFileName)
	}

	msg := lookup.Message
	if msg == "" {
		msg = "No celebrant photo available for this occasion."
	}
	return "", msg
}

func fetchURLBytes(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func SendOccasionEmail(o Occasion, greeting string, manifest *OccasionManifest, options SendOptions) SendResult {
	cfg, err := LoadEmailConfig()
	if err != nil {
		errMsg := err.Error()
		fmt.Printf("✗ Failed to send email: %s\n", errMsg)
		return SendResult{Success: false, Error: errMsg}
	}

	meta := getOccasionEmailMeta(o, options)
	shouldIncludeImage := !options.IsReminder

	var profilePictureURL, missingPhotoHint string
	if shouldIncludeImage {
		lookup := ResolveOccasionPhoto(o, manifest)
		profilePictureURL, missingPhotoHint = resolvePhotoState(o, lookup, meta.fallbackProfilePicture)
	}

	hasImageAttachment := profilePictureURL != ""
	htmlContent := generateEmailHTML(greeting, profilePictureURL, hasImageAttachment, missingPhotoHint)

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.SMTPUser)
	msg.SetHeader("To", cfg.Recipients...)
	msg.SetHeader("Subject", meta.subject)
	msg.SetBody("text/html", htmlContent)

	if shouldIncludeImage && hasImageAttachment {
		ext := inferImageExtension(profilePictureURL)
		filename := "profile-picture." + ext
		data, err := fetchURLBytes(profilePictureURL)
		if err == nil {
			msg.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(data)
				return err
			}))
		} else {
			fmt.Printf("⚠️  Failed to fetch attachment %s: %s\n", profilePictureURL, err)
		}
	}

	dialer := cfg.newDialer()
	if err := dialer.DialAndSend(msg); err != nil {
		errMsg := err.Error()
		fmt.Printf("✗ Failed to send email: %s\n", errMsg)
		return SendResult{Success: false, Error: errMsg}
	}

	return SendResult{Success: true, MessageID: generateMessageID()}
}

func SendBulkReminderEmail(greeting string) SendResult {
	cfg, err := LoadEmailConfig()
	if err != nil {
		errMsg := err.Error()
		fmt.Printf("✗ Failed to send bulk reminder email: %s\n", errMsg)
		return SendResult{Success: false, Error: errMsg}
	}

	htmlContent := generateEmailHTML(greeting, "", false, "")

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.SMTPUser)
	msg.SetHeader("To", cfg.Recipients...)
	msg.SetHeader("Subject", "Upcoming: Occasions Tomorrow")
	msg.SetBody("text/html", htmlContent)

	dialer := cfg.newDialer()
	if err := dialer.DialAndSend(msg); err != nil {
		errMsg := err.Error()
		fmt.Printf("✗ Failed to send bulk reminder email: %s\n", errMsg)
		return SendResult{Success: false, Error: errMsg}
	}

	return SendResult{Success: true, MessageID: generateMessageID()}
}

func SendMonthlySummaryEmail(greeting, _ string) SendResult {
	cfg, err := LoadEmailConfig()
	if err != nil {
		errMsg := err.Error()
		fmt.Printf("✗ Failed to send monthly summary email: %s\n", errMsg)
		return SendResult{Success: false, Error: errMsg}
	}

	htmlContent := generateEmailHTML(greeting, "", false, "")

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.SMTPUser)
	msg.SetHeader("To", cfg.Recipients...)
	msg.SetHeader("Subject", "Pauri Dgroup Occasions for the Month")
	msg.SetBody("text/html", htmlContent)

	dialer := cfg.newDialer()
	if err := dialer.DialAndSend(msg); err != nil {
		errMsg := err.Error()
		fmt.Printf("✗ Failed to send monthly summary email: %s\n", errMsg)
		return SendResult{Success: false, Error: errMsg}
	}

	return SendResult{Success: true, MessageID: generateMessageID()}
}

func generateMessageID() string {
	return fmt.Sprintf("<%d@occasion-greeter-go>", time.Now().UnixNano())
}
