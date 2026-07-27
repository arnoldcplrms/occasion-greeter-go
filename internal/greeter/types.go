package greeter

type OccasionData struct {
	MaleName              string
	MaleLastName          string
	FemaleName            string
	FemaleLastName        string
	MaleNickname          string
	FemaleNickname        string
	MaleBirthday          string
	FemaleBirthday        string
	MaleEmail             string
	FemaleEmail           string
	MaleProfilePicture    string
	FemaleProfilePicture  string
	CoupleLastName        string
	WeddingAnniversary    string
	WeddingProfilePicture string
}

type OccasionPerson struct {
	Name           string
	Nickname       string
	Email          string
	ProfilePicture string
}

type OccasionCouple struct {
	LastName       string
	MaleLastName   string
	ProfilePicture string
}

type Occasion struct {
	Type   string
	Day    int
	Person *OccasionPerson
	Couple *OccasionCouple
}

const (
	OccasionTypeBirthday     = "birthday"
	OccasionTypeAnniversary  = "anniversary"
)
