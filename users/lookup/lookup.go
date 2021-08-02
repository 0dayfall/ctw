package lookup

import (
	"encoding/json"
	"log"
	"strings"

	common "github.com/0dayfall/carboncopy/data"
	httphandler "github.com/0dayfall/carboncopy/httphandler"
)

func createIDsLookupURL() string {
	return "https://api.twitter.com/2/users"
}

func createIDLookupURL() string {
	return "https://api.twitter.com/2/users/"
}

func createUsernamesLookupURL() string {
	return "https://api.twitter.com/2/users/by"
}

func createUsernameLookupURL() string {
	return "https://api.twitter.com/2/users/by/username/"
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

func LookupID(ids string) {
	url := createIDLookupURL() + ids
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var userResponse User
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		log.Println(err)
	}
	httphandler.PrettyPrint(userResponse)
	log.Println(userResponse)
}

func LookupIDs(users []string) {
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
	httphandler.PrettyPrint(userResponse)
	log.Println(userResponse)
}

func LookupUsername(user string) {
	url := createUsernameLookupURL() + "user"
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var userResponse User
	if err := json.NewDecoder(response.Body).Decode(&userResponse); err != nil {
		log.Println(err)
	}
	httphandler.PrettyPrint(userResponse)
	log.Println(response)
}

func LookupUsernames(users []string) {
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
	httphandler.PrettyPrint(userResponse)
}
