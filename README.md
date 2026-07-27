# Birthday & Anniversary Greeter (Go)

A Go port of the [Bun/TypeScript occasion greeter](../occasion-greeter). It runs
on a cron schedule via GitHub Actions, fetches couple data from a Google Sheets
CSV, finds birthdays and anniversaries happening today, and sends personalized
greeting emails via SMTP. It also emails a reminder of tomorrow's occasions and,
on the 1st of the month, a summary of the month.

## Features

- ⏰ Runs daily at 9 AM Philippine Time via GitHub Actions cron
- 🎂 Sends personalized birthday greeting emails
- 💍 Sends wedding anniversary greeting emails
- 🖼️ Includes profile pictures in emails (downloaded from the photo manifest URL)
- 🎲 Randomly selects from pools of greeting templates
- 📦 Modular package layout with unit tests
- 🦫 Built with Go 1.23+ and the `gomail.v2` SMTP library
- 🌏 Date math is computed in `Asia/Manila` regardless of the host timezone

## Project Structure

```
occasion-greeter-go/
├── .github/workflows/
│   └── occasion-greeter.yml      # GitHub Actions workflow with cron schedule
├── cmd/
│   ├── greeter/                  # Main entrypoint
│   │   └── main.go
│   └── test-smtp/                # Verify SMTP credentials
│       └── main.go
├── internal/
│   └── greeter/                  # Library package
│       ├── constants.go          # Greeting templates
│       ├── data-loader.go        # CSV loader (Google Sheets export)
│       ├── email.go              # SMTP email service
│       ├── greeting.go           # Greeting & summary generators
│       ├── loader_test.go
│       ├── manifest.go           # Photo manifest fetcher
│       ├── matcher.go            # Today / tomorrow / month matching
│       ├── matcher_test.go
│       ├── greeting_test.go
│       ├── photoservice.go       # Google Drive photo resolver
│       ├── photoservice_test.go
│       ├── types.go              # Data structures
│       ├── utils.go              # Drive URL + name normalisation
│       └── utils_test.go
├── .env.example                  # Environment variable template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Setup

### 1. Prerequisites

- Go 1.23 or later
- A SMTP service (Gmail, SendGrid, etc.)
- A GitHub repository (for the cron workflow)

### 2. Build & Run

```Shell
# Fetch dependencies
go mod download

# Build the greeter binary
go build -o bin/greeter ./cmd/greeter

# Run the greeter (loads .env from the working directory)
./bin/greeter

# Verify SMTP credentials
go run ./cmd/test-smtp
```

### 3. Run Tests

```Shell
go test ./...
```

## Modules

### `internal/greeter/types.go`

Defines `OccasionData`, `Occasion`, `OccasionPerson`, and `OccasionCouple` Go
structs that mirror the original TypeScript interfaces.

### `internal/greeter/utils.go`

- `ExtractDriveFileID` — Pulls a Google Drive file ID out of `?id=`,
  `/file/d/...`, or `/d/...` URLs.
- `NormalizeEmailBase` — Lowercases and trims an email address.
- `NormalizeSurnameBase` — Strips diacritics, lowercases, and converts spaces
  and hyphens to underscores.
- `NormalizeManifestPhotoURL` — Resolves any supported Google Drive URL to a
  direct viewable URL.
- `NormalizeManifestSourceURL` — Resolves any supported Google Drive URL to a
  direct download URL (used for the manifest JSON).

### `internal/greeter/data-loader.go`

- `LoadOccasionsData` — Fetches and parses the `CSV_URL` environment variable.
  The first row is treated as the header; known Google Sheets column names are
  mapped to the camelCase `OccasionData` fields. Unknown headers are passed
  through, so legacy camelCase CSVs continue to work.
- `LoadOccasionsDataFromFile` — Same parser, but reads from a local file
  (handy for tests).

### `internal/greeter/matcher.go`

- `FindOccasionsToday` — Birthdays and anniversaries happening in PH time
  today.
- `FindOccasionsTomorrow` — Same, but offset by one day.
- `FindOccasionsThisMonth` — All occasions in a given month (used for the
  monthly summary).

### `internal/greeter/greeting.go`

- `GenerateGreeting` — Picks a random birthday or anniversary template and
  replaces the `<<nickname>>` or `<<lastName>>` placeholder.
- `GenerateReminderGreeting` — Single-occasion reminder text.
- `GenerateBulkReminderGreeting` — Multi-occasion reminder email body.
- `GenerateMonthlySummaryGreeting` — Monthly summary text.

### `internal/greeter/photoservice.go` & `internal/greeter/manifest.go`

Resolve celebrant photos via the JSON manifest. The manifest shape is the
same as the original project (`birthday` and `wedding_anniversary` maps).

### `internal/greeter/email.go`

- `SendOccasionEmail` — Sends a single greeting email, including the celebrant
  profile picture as an attachment when one is available.
- `SendBulkReminderEmail` — Sends the combined tomorrow reminder.
- `SendMonthlySummaryEmail` — Sends the monthly summary.

## Data Format

The CSV must contain these header columns (Google Sheets auto-quotes the
header names):

| Column                              | Required | Notes                                   |
| ----------------------------------- | -------- | --------------------------------------- |
| `Male's First Name`                 | ✅       |                                         |
| `Male's Last Name`                  |          | Falls back to `coupleLastName` if empty |
| `Male's Nickname (If Applicable)`   |          | Falls back to first name                |
| `Male's Birthday`                   |          | `DD/MM/YYYY`                            |
| `Male's Email Address`              |          | Used to look up the celebrant photo     |
| `Male's Profile Picture`            |          | Fallback if no manifest entry           |
| `Female's First Name`               | ✅       |                                         |
| `Female's Last Name`                |          |                                         |
| `Female's Nickname (If Applicable)` |          |                                         |
| `Female's Birthday`                 |          | `DD/MM/YYYY`                            |
| `Female's Email Address`            |          |                                         |
| `Female's Profile Picture`          |          |                                         |
| `Couple's Profile Picture`          |          | Used for the anniversary photo          |
| `Wedding Anniversary (...)`         |          | `DD/MM/YYYY`                            |

CamelCase column names also work; the loader only remaps the human-readable
Google Sheets names.

### Celebrant Photo Filename Rules

- **Birthday:** manifest key is the celebrant's full email in lowercase.
- **Anniversary:** manifest key is the male's surname normalised to lowercase
  snake_case.

If a celebrant photo is not found, the email is still sent with a short
note explaining that no photo is available.

## Cron Schedule

The workflow runs at **9 AM Philippine Time** every day:

```YAML
- cron: '0 1 * * *' # 1 AM UTC = 9 AM PH (UTC+8)
```

Edit `.github/workflows/occasion-greeter.yml` to adjust.

## GitHub Secrets

Add these secrets to your repository (Settings → Secrets and variables →
Actions):

- `CSV_URL`
- `OCCASION_PHOTO_MANIFEST_URL`
- `RECIPIENT_EMAIL`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USER`
- `SMTP_PASSWORD`

## SMTP Setup

The same SMTP providers documented in the original project work here. See the
previous README for provider-specific instructions (Gmail app passwords,
SendGrid, Outlook, etc.).

## Troubleshooting

The `cmd/test-smtp` subcommand performs a quick SMTP dial test using the
current `.env` / environment variables. Use it whenever you change credentials:

```Shell
go run ./cmd/test-smtp
```

## Migration Notes

- `COUPLES_DATA` is no longer consumed directly; the CSV at `CSV_URL` is the
  single source of truth (same as the current Bun project).
- `OCCASION_PHOTO_MANIFEST_URL` is the only manifest URL — the older
  `MANIFEST_URL` aliases are not supported.
- All times are computed in `Asia/Manila`; the workflow still uses UTC cron
  for the 9 AM PH trigger.
