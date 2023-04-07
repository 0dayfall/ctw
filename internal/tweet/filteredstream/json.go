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
	Data []Data `json:"data"`
	Meta Meta   `json:"meta"`
}

type Data struct {
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
