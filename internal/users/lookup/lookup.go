package user

import (
	"encoding/json"
	"log"
	"strings"

	common "github.com/0dayfall/ctw/internal/data"
	"github.com/0dayfall/ctw/internal/httphandler"
)

const (
	users           = "/2/users"
	userById        = "/2/users/:id"
	usersByUsername = "/2/users/by"
	userByUsername  = "/2/users/by/username/:username"
)

func createIDsLookupURL() string {
	return common.APIurl + users
}

func createLookupUserByIdURL(id string) string {
	return common.APIurl + strings.Replace(userById, ":id", id, 1)
}

func createUsernamesLookupURL() string {
	return common.APIurl + users
}

func createUsernameLookupURL(username string) string {
	return common.APIurl + strings.Replace(userByUsername, ":username", username, 1)
}

type User struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	UserName        string          `json:"username"`
	CreatedAt       string          `json:"created_at"`
	Description     string          `json:"description"`
	Entities        common.Entities `json:"entities"`
	Location        string          `json:"location"`
	PinnedTweetID   string          `json:"pinned_tweet_id"`
	ProfileImageURL string          `json:"profile_image_url"`
	Protected       bool            `json:"protected"`
	PublicMetrics   UserMetrics     `json:"public_metrics"`
	URL             string          `json:"url"`
	Verified        bool            `json:"verified"`
	WithHeld        WithHeld        `json:"withheld"`
}

type WithHeld struct {
	Copyright    bool     `json:"copyright"`
	CountryCodes []string `json:"country_codes"`
}

// UserMetricsObj contains details about activity for this user
type UserMetrics struct {
	Followers int `json:"followers_count"`
	Following int `json:"following_count"`
	Tweets    int `json:"tweet_count"`
	Listed    int `json:"listed_count"`
}

func LookupID(id string) User {
	url := createLookupUserByIdURL(id)
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var userResponse User
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		log.Println(err)
	}
	return userResponse
}

func LookupIDs(users []string) User {
	url := createUsernamesLookupURL()
	req := httphandler.CreateGetRequest(url)
	q := req.URL.Query()
	var userNames string
	for _, user := range users {
		userNames += user + ", "
	}
	q.Add("usernames", userNames)
	req.URL.RawQuery = q.Encode()
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var userResponse User
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		log.Println(err)
	}
	return userResponse
}

func LookupUsername(user string) User {
	url := createUsernameLookupURL(user)
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var userResponse User
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		log.Println(err)
	}
	return userResponse
}

func LookupUsernames(users []string) User {
	url := createUsernamesLookupURL()
	req := httphandler.CreateGetRequest(url)
	q := req.URL.Query()
	userNames := strings.Join(users[:], ",")
	q.Add("usernames", userNames)
	req.URL.RawQuery = q.Encode()
	log.Println("GET " + req.URL.String())
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()
	httphandler.IsResponseOK(response)
	var userResponse User
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		log.Println(err)
	}
	return userResponse
}
