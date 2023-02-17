package utils

import (
	"os"
	"time"
)

func GetEnvVars(key, defultValue string) string {
	envVal := os.Getenv(key)
	if envVal == "" {
		return defultValue
	}
	return envVal
}

func Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return time.Now().In(loc)
}

func GetMapValidations() map[string]interface{} {
	return map[string]interface{}{
		"id":                  "",
		"created_at":          "",
		"updated_at":          "",
		"name":                "stringlength(2|15)",
		"description":         "maxstringlength(3000)",
		"amount_of_employees": "numeric",
		"registered":          "type(bool)",
		"type":                "in(Corporations|NonProfit|Cooperative|Sole Proprietorship)",
	}
}
