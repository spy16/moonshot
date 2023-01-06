package user

import (
	"strings"
	"time"
)

func mergeUsers(base, override User) User {
	now := time.Now()

	return User{
		ID:        base.ID,
		CreatedAt: base.CreatedAt,
		UpdatedAt: now,
		Name:      pickNonEmpty(override.Name, base.Name),
		AvatarURL: pickNonEmpty(override.AvatarURL, base.AvatarURL),
		Email:     pickNonEmpty(override.Email, base.Email),
		Gender:    pickNonEmpty(override.Gender, base.Gender),
		Locale:    pickNonEmpty(override.Locale, base.Locale),
		Location:  pickNonEmpty(override.Location, base.Location),
	}
}

func pickNonEmpty(arr ...string) string {
	for _, s := range arr {
		s = strings.TrimSpace(s)
		if s != "" {
			return s
		}
	}
	return ""
}

func mergeMaps(m1, m2 map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m1 {
		res[k] = v
	}
	for k, v := range m2 {
		res[k] = v
	}
	return res
}
