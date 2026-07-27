package greeter

import (
	"strconv"
	"strings"
	"time"
)

type dateParts struct {
	Day   int
	Month int
}

func parseDate(dateStr string) *dateParts {
	if strings.TrimSpace(dateStr) == "" {
		return nil
	}

	parts := strings.Split(dateStr, "/")
	if len(parts) != 3 {
		return nil
	}

	day, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil
	}
	month, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return nil
	}

	if day < 1 || day > 31 || month < 1 || month > 12 {
		return nil
	}

	return &dateParts{Day: day, Month: month}
}

func getDatePH(dayOffset int) dateParts {
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		loc = time.FixedZone("PHT", 8*60*60)
	}
	now := time.Now().In(loc)
	if dayOffset != 0 {
		now = now.AddDate(0, 0, dayOffset)
	}
	return dateParts{
		Day:   now.Day(),
		Month: int(now.Month()),
	}
}

func addDaysToDateParts(d dateParts, daysToAdd int) dateParts {
	t := time.Date(2000, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	t = t.AddDate(0, 0, daysToAdd)
	return dateParts{
		Day:   t.Day(),
		Month: int(t.Month()),
	}
}

func findOccasionsForDate(data []OccasionData, target dateParts) []Occasion {
	var occasions []Occasion

	for _, d := range data {
		if mb := parseDate(d.MaleBirthday); mb != nil &&
			mb.Day == target.Day && mb.Month == target.Month {
			nick := d.MaleNickname
			if nick == "" {
				nick = d.MaleName
			}
			occasions = append(occasions, Occasion{
				Type: OccasionTypeBirthday,
				Person: &OccasionPerson{
					Name:           d.MaleName,
					Nickname:       nick,
					Email:          d.MaleEmail,
					ProfilePicture: d.MaleProfilePicture,
				},
			})
		}

		if fb := parseDate(d.FemaleBirthday); fb != nil &&
			fb.Day == target.Day && fb.Month == target.Month {
			nick := d.FemaleNickname
			if nick == "" {
				nick = d.FemaleName
			}
			occasions = append(occasions, Occasion{
				Type: OccasionTypeBirthday,
				Person: &OccasionPerson{
					Name:           d.FemaleName,
					Nickname:       nick,
					Email:          d.FemaleEmail,
					ProfilePicture: d.FemaleProfilePicture,
				},
			})
		}

		if ann := parseDate(d.WeddingAnniversary); ann != nil &&
			ann.Day == target.Day && ann.Month == target.Month {
			occasions = append(occasions, Occasion{
				Type: OccasionTypeAnniversary,
				Couple: &OccasionCouple{
					LastName:       d.CoupleLastName,
					MaleLastName:   d.MaleLastName,
					ProfilePicture: d.WeddingProfilePicture,
				},
			})
		}
	}

	return occasions
}

func findOccasionsForMonth(data []OccasionData, month int) ([]Occasion, string) {
	seen := make(map[string]bool)
	var occasions []Occasion

	for _, d := range data {
		if mb := parseDate(d.MaleBirthday); mb != nil && mb.Month == month {
			key := "birthday:" + strings.ToLower(d.MaleName)
			if !seen[key] {
				seen[key] = true
				nick := d.MaleNickname
				if nick == "" {
					nick = d.MaleName
				}
				occasions = append(occasions, Occasion{
					Type: OccasionTypeBirthday,
					Day:  mb.Day,
					Person: &OccasionPerson{
						Name:           d.MaleName,
						Nickname:       nick,
						Email:          d.MaleEmail,
						ProfilePicture: d.MaleProfilePicture,
					},
				})
			}
		}

		if fb := parseDate(d.FemaleBirthday); fb != nil && fb.Month == month {
			key := "birthday:" + strings.ToLower(d.FemaleName)
			if !seen[key] {
				seen[key] = true
				nick := d.FemaleNickname
				if nick == "" {
					nick = d.FemaleName
				}
				occasions = append(occasions, Occasion{
					Type: OccasionTypeBirthday,
					Day:  fb.Day,
					Person: &OccasionPerson{
						Name:           d.FemaleName,
						Nickname:       nick,
						Email:          d.FemaleEmail,
						ProfilePicture: d.FemaleProfilePicture,
					},
				})
			}
		}

		if ann := parseDate(d.WeddingAnniversary); ann != nil && ann.Month == month {
			key := "anniversary:" + strings.ToLower(d.CoupleLastName)
			if !seen[key] {
				seen[key] = true
				occasions = append(occasions, Occasion{
					Type: OccasionTypeAnniversary,
					Day:  ann.Day,
					Couple: &OccasionCouple{
						LastName:       d.CoupleLastName,
						MaleLastName:   d.MaleLastName,
						ProfilePicture: d.WeddingProfilePicture,
					},
				})
			}
		}
	}

	monthName := time.Date(2000, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("January")
	return occasions, monthName
}

func FindOccasionsToday(data []OccasionData) []Occasion {
	return findOccasionsForDate(data, getDatePH(0))
}

func FindOccasionsTomorrow(data []OccasionData) []Occasion {
	today := getDatePH(0)
	tomorrow := addDaysToDateParts(today, 1)
	return findOccasionsForDate(data, tomorrow)
}

func FindOccasionsThisMonth(data []OccasionData) ([]Occasion, string) {
	now := getDatePH(0)
	return findOccasionsForMonth(data, now.Month)
}
