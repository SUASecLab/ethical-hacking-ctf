package main

type User struct {
	Uuid           string   `json:"uuid"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	VisitCardUrl   string   `json:"visitCardUrl"`
	Tags           []string `json:"tags"`
	AvailableFlags []string `json:"availableFlags"`
	CollectedFlags []string `json:"collectedFlags"`
}

type Flag struct {
	Flag        string `json:"flag"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type IndexData struct {
	Token                 string
	SuccessMessage        string
	ErrorMessage          string
	BashProgress          int
	SUASploitableProgress int
	BashFlags             []Flag
	SUASploitableFlags    []Flag
}
