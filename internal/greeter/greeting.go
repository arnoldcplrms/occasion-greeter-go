package greeter

import (
	"fmt"
	"math/rand"
	"strings"
)

func pickRandom(templates []string, rng *rand.Rand) string {
	if len(templates) == 0 {
		return ""
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(int64(len(templates))))
	}
	return templates[rng.Intn(len(templates))]
}

func GenerateGreeting(o Occasion, rng *rand.Rand) string {
	if o.Type == OccasionTypeBirthday && o.Person != nil {
		return strings.Replace(pickRandom(BirthdayGreetings, rng), "<<nickname>>", o.Person.Nickname, 1)
	}
	if o.Type == OccasionTypeAnniversary && o.Couple != nil {
		return strings.Replace(pickRandom(WeddingAnniversaryGreetings, rng), "<<lastName>>", o.Couple.LastName, 1)
	}
	return ""
}

func GenerateReminderGreeting(o Occasion) string {
	if o.Type == OccasionTypeBirthday && o.Person != nil {
		return fmt.Sprintf("Reminder: Tomorrow is %s's birthday. Please prepare to greet and pray for %s.", o.Person.Name, o.Person.Nickname)
	}
	if o.Type == OccasionTypeAnniversary && o.Couple != nil {
		return fmt.Sprintf("Reminder: Tomorrow is Team %s's wedding anniversary. Please prepare to celebrate and pray for their marriage.", o.Couple.LastName)
	}
	return ""
}

func formatReminderListItem(o Occasion) string {
	if o.Type == OccasionTypeBirthday && o.Person != nil {
		return fmt.Sprintf("%s's birthday", o.Person.Name)
	}
	if o.Type == OccasionTypeAnniversary && o.Couple != nil {
		return fmt.Sprintf("Team %s's wedding anniversary", o.Couple.LastName)
	}
	return "An upcoming occasion"
}

func GenerateBulkReminderGreeting(occasions []Occasion) string {
	lines := []string{"Reminder: These occasions are happening tomorrow:", ""}
	for _, o := range occasions {
		lines = append(lines, "- "+formatReminderListItem(o))
	}
	return strings.Join(lines, "\n")
}

func GenerateMonthlySummaryGreeting(occasions []Occasion, monthName string) string {
	header := fmt.Sprintf("📋 Monthly Summary for %s", monthName)

	if len(occasions) == 0 {
		return strings.Join([]string{
			header,
			"",
			"There are no occasions scheduled for this month.",
		}, "\n")
	}

	shortMonth := monthName
	if len(shortMonth) > 3 {
		shortMonth = shortMonth[:3]
	}

	var lines []string
	for _, o := range occasions {
		if o.Type == OccasionTypeBirthday && o.Person != nil {
			lines = append(lines, fmt.Sprintf("🎂 %s's birthday (%s %d)", o.Person.Name, shortMonth, o.Day))
		}
		if o.Type == OccasionTypeAnniversary && o.Couple != nil {
			lines = append(lines, fmt.Sprintf("💍 Team %s's wedding anniversary (%s %d)", o.Couple.LastName, shortMonth, o.Day))
		}
	}

	return strings.Join([]string{
		header,
		"",
		strings.Join(lines, "\n"),
		"",
		fmt.Sprintf("Total: %d occasion(s) this month.", len(occasions)),
	}, "\n")
}
