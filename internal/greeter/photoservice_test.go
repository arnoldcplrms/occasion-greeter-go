package greeter

import "testing"

func TestResolveOccasionPhoto_BirthdayFound(t *testing.T) {
	manifest := &OccasionManifest{
		Birthday: map[string]string{
			"foo@example.com": "https://drive.google.com/open?id=ABC",
		},
	}
	o := Occasion{
		Type:   OccasionTypeBirthday,
		Person: &OccasionPerson{Email: "Foo@Example.com"},
	}
	got := ResolveOccasionPhoto(o, manifest)
	if got.Status != LookupStatusFound {
		t.Fatalf("expected found, got %s: %s", got.Status, got.Message)
	}
	if got.URL != "https://drive.google.com/uc?export=view&id=ABC" {
		t.Errorf("unexpected url: %s", got.URL)
	}
}

func TestResolveOccasionPhoto_AnniversaryFallback(t *testing.T) {
	manifest := &OccasionManifest{
		WeddingAnnivesary: map[string]string{
			"romero_sa": "https://drive.google.com/open?id=ZZZ",
		},
	}
	o := Occasion{
		Type:   OccasionTypeAnniversary,
		Couple: &OccasionCouple{MaleLastName: "Romero-Sa"},
	}
	got := ResolveOccasionPhoto(o, manifest)
	if got.Status != LookupStatusFound {
		t.Fatalf("expected found via fallback key, got %s", got.Status)
	}
}

func TestResolveOccasionPhoto_NotConfiguredUsesFallback(t *testing.T) {
	o := Occasion{
		Type:   OccasionTypeBirthday,
		Person: &OccasionPerson{Email: "a@b.com", ProfilePicture: "https://drive.google.com/open?id=DEF"},
	}
	got := ResolveOccasionPhoto(o, nil)
	if got.Status != LookupStatusNotConfigured {
		t.Fatalf("expected not-configured, got %s", got.Status)
	}
	url, _ := resolvePhotoState(o, got, o.Person.ProfilePicture)
	if url != "https://drive.google.com/uc?export=view&id=DEF" {
		t.Errorf("unexpected fallback url: %s", url)
	}
}

func TestResolveOccasionPhoto_Missing(t *testing.T) {
	manifest := &OccasionManifest{
		Birthday: map[string]string{},
	}
	o := Occasion{
		Type:   OccasionTypeBirthday,
		Person: &OccasionPerson{Email: "missing@b.com"},
	}
	got := ResolveOccasionPhoto(o, manifest)
	if got.Status != LookupStatusMissing {
		t.Fatalf("expected missing, got %s", got.Status)
	}
}
