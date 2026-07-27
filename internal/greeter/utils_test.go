package greeter

import "testing"

func TestExtractDriveFileID(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want string
	}{
		{"query id", "https://drive.google.com/uc?export=view&id=ABC123", "ABC123"},
		{"file d path", "https://drive.google.com/file/d/XYZ_9/view?usp=sharing", "XYZ_9"},
		{"open id", "https://drive.google.com/open?id=ABCD-EFG", "ABCD-EFG"},
		{"uc id", "https://drive.google.com/uc?export=download&id=ZZZ", "ZZZ"},
		{"non-drive", "https://example.com/foo.png", ""},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ExtractDriveFileID(c.url)
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestNormalizeSurnameBase(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "Ramos", "ramos"},
		{"with diacritic", "Ramos", "ramos"},
		{"spaces to underscore", "San Pedro", "san_pedro"},
		{"hyphens to underscore", "San-Pedro", "san_pedro"},
		{"punctuation removed", "O'Brien", "obrien"},
		{"trim", "  Smith  ", "smith"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NormalizeSurnameBase(c.in)
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestNormalizeEmailBase(t *testing.T) {
	got := NormalizeEmailBase("  Foo@Example.COM  ")
	want := "foo@example.com"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNormalizeManifestPhotoURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"open id", "https://drive.google.com/open?id=ABC", "https://drive.google.com/uc?export=view&id=ABC"},
		{"file d", "https://drive.google.com/file/d/DEF/view?usp=sharing", "https://drive.google.com/uc?export=view&id=DEF"},
		{"passthrough", "https://example.com/pic.jpg", "https://example.com/pic.jpg"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NormalizeManifestPhotoURL(c.in)
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestNormalizeManifestSourceURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"file d", "https://drive.google.com/file/d/DEF/view?usp=sharing", "https://drive.google.com/uc?export=download&id=DEF"},
		{"open id", "https://drive.google.com/open?id=ABC", "https://drive.google.com/uc?export=download&id=ABC"},
		{"passthrough", "https://example.com/manifest.json", "https://example.com/manifest.json"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NormalizeManifestSourceURL(c.in)
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}
