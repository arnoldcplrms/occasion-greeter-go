package greeter

import (
	"strings"
	"testing"
)

func TestGenerateGreeting_Birthday(t *testing.T) {
	o := Occasion{
		Type: OccasionTypeBirthday,
		Person: &OccasionPerson{
			Name:     "Arnold",
			Nickname: "Arnold",
		},
	}
	got := GenerateGreeting(o, nil)
	if got == "" {
		t.Fatal("expected non-empty greeting")
	}
	if !strings.Contains(got, "Arnold") {
		t.Errorf("expected greeting to contain nickname, got %q", got)
	}
	if strings.Contains(got, "<<nickname>>") {
		t.Errorf("placeholder not replaced: %q", got)
	}
}

func TestGenerateGreeting_Anniversary(t *testing.T) {
	o := Occasion{
		Type: OccasionTypeAnniversary,
		Couple: &OccasionCouple{LastName: "Ramos"},
	}
	got := GenerateGreeting(o, nil)
	if got == "" {
		t.Fatal("expected non-empty greeting")
	}
	if !strings.Contains(got, "Ramos") {
		t.Errorf("expected greeting to contain last name, got %q", got)
	}
	if strings.Contains(got, "<<lastName>>") {
		t.Errorf("placeholder not replaced: %q", got)
	}
}

func TestGenerateBulkReminderGreeting(t *testing.T) {
	occasions := []Occasion{
		{Type: OccasionTypeBirthday, Person: &OccasionPerson{Name: "Arnold"}},
		{Type: OccasionTypeAnniversary, Couple: &OccasionCouple{LastName: "Ramos"}},
	}
	got := GenerateBulkReminderGreeting(occasions)
	if !strings.Contains(got, "Arnold's birthday") {
		t.Errorf("missing birthday line: %q", got)
	}
	if !strings.Contains(got, "Team Ramos's wedding anniversary") {
		t.Errorf("missing anniversary line: %q", got)
	}
}

func TestGenerateMonthlySummaryGreeting_Empty(t *testing.T) {
	got := GenerateMonthlySummaryGreeting(nil, "January")
	if !strings.Contains(got, "no occasions") {
		t.Errorf("expected 'no occasions' message, got %q", got)
	}
}

func TestGenerateMonthlySummaryGreeting_WithOccasions(t *testing.T) {
	occasions := []Occasion{
		{Type: OccasionTypeBirthday, Day: 20, Person: &OccasionPerson{Name: "Arnold"}},
		{Type: OccasionTypeAnniversary, Day: 27, Couple: &OccasionCouple{LastName: "Ramos"}},
	}
	got := GenerateMonthlySummaryGreeting(occasions, "July")
	if !strings.Contains(got, "Arnold's birthday") {
		t.Errorf("missing birthday entry: %q", got)
	}
	if !strings.Contains(got, "Ramos") {
		t.Errorf("missing anniversary entry: %q", got)
	}
	if !strings.Contains(got, "Total: 2") {
		t.Errorf("missing total: %q", got)
	}
}
