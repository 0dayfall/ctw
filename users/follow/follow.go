package user

import "strings"

const (
	following = "/2/users/:id/following"
	followers = "/2/users/:id/followers"
	unfollow  = "/2/users/:source_user_id/following/:target_user_id"
)

func createFollowingURL(id string) string {
	url := strings.Replace(following, ":id", id, 1)
	return url
}

func createFollowersURL(id string) string {
	url := strings.Replace(followers, ":id", id, 1)
	return url
}

func createUnfollowURL(sourceId, targetId string) string {
	url := strings.Replace(unfollow, ":source_user_id", sourceId, 1)
	url = strings.Replace(url, ":target_user_id", targetId, 1)
	return url
}
