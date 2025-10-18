package tweet

import "time"

type DeleteIdCommand struct {
	Delete DeleteId `json:"delete"`
}

type DeleteValueCommand struct {
	Delete DeleteValue `json:"delete"`
}

type DeleteId struct {
	Ids []string `json:"ids"`
}

type DeleteValue struct {
	Values []string `json:"values"`
}

type AddCommand struct {
	Add []Add `json:"add"`
}

type Add struct {
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

type RulesResponse struct {
	Data []RuleData `json:"data"`
	Meta Meta       `json:"meta"`
}

type RuleData struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

type Summary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Valid      int `json:"valid"`
	Invalid    int `json:"invalid"`
}

type Meta struct {
	Sent    time.Time `json:"sent"`
	Summary Summary   `json:"summary,omitempty"`
}

func CreateDeleteIdCommand(ids []string) DeleteIdCommand {
	var cmd DeleteIdCommand
	cmd.Delete.Ids = ids
	return cmd
}

func CreateDeleteValueCommand(rules []string) DeleteValueCommand {
	var cmd DeleteValueCommand
	cmd.Delete.Values = rules
	return cmd
}

func CreateAddCommand(values map[string]string) AddCommand {
	var cmd AddCommand
	for value, tag := range values {
		cmd.Add = append(cmd.Add, Add{
			Value: value,
			Tag:   tag,
		})
	}
	return cmd
}

type RulesError struct {
	ClientID           string `json:"client_id"`
	RequiredEnrollment string `json:"required_enrollment"`
	RegistrationURL    string `json:"registration_url"`
	Title              string `json:"title"`
	Detail             string `json:"detail"`
	Reason             string `json:"reason"`
	Type               string `json:"type"`
}

type StreamResponse struct {
	Id   string
	Text string
}

// StreamEnvelope captures a response from the filtered stream endpoint.
type StreamEnvelope struct {
	Data     []StreamTweet  `json:"data"`
	Includes StreamIncludes `json:"includes"`
	Errors   []RulesError   `json:"errors,omitempty"`
	Meta     StreamMeta     `json:"meta,omitempty"`
}

// StreamTweet represents a Tweet entry in the filtered stream.
type StreamTweet struct {
	ID                string    `json:"id"`
	Text              string    `json:"text"`
	AuthorID          string    `json:"author_id"`
	CreatedAt         time.Time `json:"created_at"`
	PossiblySensitive bool      `json:"possibly_sensitive"`
	Source            string    `json:"source"`
	Lang              string    `json:"lang"`
}

// StreamIncludes captures expanded entities such as users.
type StreamIncludes struct {
	Users []StreamUser `json:"users,omitempty"`
}

// StreamUser represents a user record included alongside tweets.
type StreamUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// StreamMeta carries additional metadata from the endpoint.
type StreamMeta struct {
	ResultCount int `json:"result_count"`
}
