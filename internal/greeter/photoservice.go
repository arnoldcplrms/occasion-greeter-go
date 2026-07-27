package greeter

type LookupStatus string

const (
	LookupStatusFound           LookupStatus = "found"
	LookupStatusMissing         LookupStatus = "missing"
	LookupStatusNotConfigured   LookupStatus = "not-configured"
	LookupStatusError           LookupStatus = "error"
)

type PhotoLookupResult struct {
	Status            LookupStatus
	URL               string
	Message           string
	AttemptedFileName string
}

type OccasionManifest struct {
	Birthday          map[string]string `json:"birthday"`
	WeddingAnniversary map[string]string `json:"wedding_anniversary"`
	WeddingAnnivesary map[string]string `json:"wedding_annivesary"`
}

func lookupFromManifest(group string, lookupKey string, manifest *OccasionManifest) PhotoLookupResult {
	if manifest == nil {
		return PhotoLookupResult{
			Status: LookupStatusNotConfigured,
			Message: "Occasion manifest is not configured. Falling back to existing profile picture when available.",
		}
	}

	var bucket map[string]string
	switch group {
	case "birthday":
		bucket = manifest.Birthday
	case "wedding_anniversary":
		if len(manifest.WeddingAnniversary) > 0 {
			bucket = manifest.WeddingAnniversary
		} else {
			bucket = manifest.WeddingAnnivesary
		}
	default:
		bucket = nil
	}

	if bucket == nil {
		bucket = map[string]string{}
	}

	rawURL := bucket[lookupKey]
	url := NormalizeManifestPhotoURL(rawURL)
	if url != "" {
		return PhotoLookupResult{
			Status:            LookupStatusFound,
			URL:               url,
			AttemptedFileName: lookupKey,
		}
	}

	return PhotoLookupResult{
		Status:            LookupStatusMissing,
		URL:               "",
		AttemptedFileName: lookupKey,
		Message:           "No celebrant photo available for this occasion.",
	}
}

func ResolveOccasionPhoto(occasion Occasion, manifest *OccasionManifest) PhotoLookupResult {
	if occasion.Type == OccasionTypeBirthday && occasion.Person != nil {
		lookupKey := NormalizeEmailBase(occasion.Person.Email)
		if lookupKey == "" {
			return PhotoLookupResult{
				Status:  LookupStatusMissing,
				URL:     "",
				Message: "No celebrant photo available for this occasion.",
			}
		}
		return lookupFromManifest("birthday", lookupKey, manifest)
	}

	if occasion.Type == OccasionTypeAnniversary && occasion.Couple != nil {
		lookupKey := NormalizeSurnameBase(occasion.Couple.MaleLastName)
		if lookupKey == "" {
			return PhotoLookupResult{
				Status:  LookupStatusMissing,
				URL:     "",
				Message: "No celebrant photo available for this occasion.",
			}
		}
		return lookupFromManifest("wedding_anniversary", lookupKey, manifest)
	}

	return PhotoLookupResult{
		Status:  LookupStatusError,
		URL:     "",
		Message: "Invalid occasion type for photo lookup.",
	}
}
