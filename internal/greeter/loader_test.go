package greeter

import (
	"strings"
	"testing"
)

const sampleCSV = `"Male's First Name","Male's Last Name","Male's Nickname (If Applicable)","Male's Birthday","Male's Email Address","Male's Profile Picture","Female's First Name","Female's Last Name","Female's Nickname (If Applicable)","Female's Birthday","Female's Email Address","Female's Profile Picture","Couple's Profile Picture","Wedding Anniversary (If engaged, set the future date/ If Not Applicable Skip)"
Arnold,Ramos,,20/07/1995,arnold@example.com,http://m.example/m.jpg,Krizzia Mae,Ramos,Cha,23/12/1995,cha@example.com,http://m.example/f.jpg,http://m.example/c.jpg,27/05/2023
"Earl John",Gallarde,Ej,03/12/1986,ej@example.com,http://m.example/ejm.jpg,Sigrid,Gallarde,Sigrid,03/12/1986,sig@example.com,http://m.example/sig.jpg,http://m.example/g.jpg,28/10/2016
Carlo,Renoria,Caloy,09/09/1997,caloy@example.com,http://m.example/caloy.jpg,Irene,Renoria,Rin,10/03/1998,rin@example.com,http://m.example/rin.jpg,http://m.example/r.jpg,
`

func TestLoadOccasionsDataFromReader(t *testing.T) {
	rows, err := loadCSVFromReader(strings.NewReader(sampleCSV))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}

	r0 := rows[0]
	if r0.MaleName != "Arnold" {
		t.Errorf("row0.MaleName = %q", r0.MaleName)
	}
	if r0.MaleBirthday != "20/07/1995" {
		t.Errorf("row0.MaleBirthday = %q", r0.MaleBirthday)
	}
	if r0.FemaleNickname != "Cha" {
		t.Errorf("row0.FemaleNickname = %q", r0.FemaleNickname)
	}
	if r0.WeddingAnniversary != "27/05/2023" {
		t.Errorf("row0.WeddingAnniversary = %q", r0.WeddingAnniversary)
	}
	if r0.CoupleLastName != "Ramos" {
		t.Errorf("row0.CoupleLastName = %q", r0.CoupleLastName)
	}

	r2 := rows[2]
	if r2.WeddingAnniversary != "" {
		t.Errorf("row2.WeddingAnniversary should be empty, got %q", r2.WeddingAnniversary)
	}
}

func TestLoadOccasionsDataMissingRequired(t *testing.T) {
	csv := `Male's First Name,Female's First Name
,Joanna`
	_, err := loadCSVFromReader(strings.NewReader(csv))
	if err == nil {
		t.Fatalf("expected error for missing maleName")
	}
}
