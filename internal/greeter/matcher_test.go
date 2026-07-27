package greeter

import "testing"

func TestFindOccasionsForDate_Match(t *testing.T) {
	data := []OccasionData{
		{
			MaleName:           "Arnold",
			MaleBirthday:       "20/07/1995",
			MaleEmail:          "a@x.com",
			MaleProfilePicture: "http://m/a.jpg",
			FemaleName:         "Krizzia",
			FemaleBirthday:     "23/12/1995",
			CoupleLastName:     "Ramos",
			MaleLastName:       "Ramos",
			WeddingAnniversary: "27/05/2023",
		},
	}
	target := dateParts{Day: 20, Month: 7}
	got := findOccasionsForDate(data, target)
	if len(got) != 1 {
		t.Fatalf("expected 1 occasion, got %d", len(got))
	}
	if got[0].Type != OccasionTypeBirthday || got[0].Person.Name != "Arnold" {
		t.Errorf("unexpected: %+v", got[0])
	}
}

func TestFindOccasionsForDate_Anniversary(t *testing.T) {
	data := []OccasionData{
		{
			MaleName:             "Arnold",
			MaleLastName:         "Ramos",
			CoupleLastName:       "Ramos",
			WeddingAnniversary:   "27/05/2023",
			WeddingProfilePicture: "http://m/c.jpg",
		},
	}
	target := dateParts{Day: 27, Month: 5}
	got := findOccasionsForDate(data, target)
	if len(got) != 1 {
		t.Fatalf("expected 1 occasion, got %d", len(got))
	}
	if got[0].Type != OccasionTypeAnniversary || got[0].Couple.LastName != "Ramos" {
		t.Errorf("unexpected: %+v", got[0])
	}
}

func TestFindOccasionsForDate_SkipsEmptyDates(t *testing.T) {
	data := []OccasionData{
		{MaleName: "Arnold", MaleBirthday: "", CoupleLastName: "Ramos"},
	}
	got := findOccasionsForDate(data, dateParts{Day: 1, Month: 1})
	if len(got) != 0 {
		t.Errorf("expected no occasions, got %d", len(got))
	}
}

func TestFindOccasionsForDate_InvalidDate(t *testing.T) {
	data := []OccasionData{
		{MaleName: "Arnold", MaleBirthday: "1995-07-20", CoupleLastName: "Ramos"},
	}
	got := findOccasionsForDate(data, dateParts{Day: 20, Month: 7})
	if len(got) != 0 {
		t.Errorf("invalid date should not match, got %d", len(got))
	}
}

func TestFindOccasionsForMonth_Deduplicates(t *testing.T) {
	data := []OccasionData{
		{MaleName: "Arnold", MaleBirthday: "20/07/1995", CoupleLastName: "Ramos"},
		{MaleName: "Arnold", MaleBirthday: "20/07/1990", CoupleLastName: "Ramos"},
	}
	got, name := findOccasionsForMonth(data, 7)
	if name != "July" {
		t.Errorf("expected month 'July', got %q", name)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 unique occasion, got %d", len(got))
	}
}

func TestAddDaysToDateParts(t *testing.T) {
	got := addDaysToDateParts(dateParts{Day: 31, Month: 12}, 1)
	if got.Day != 1 || got.Month != 1 {
		t.Errorf("expected 1/1, got %d/%d", got.Day, got.Month)
	}
}
