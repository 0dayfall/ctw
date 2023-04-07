package user

import (
	"strings"

	common "github.com/0dayfall/ctw/internal/data"
)

const (
	following = "/2/users/:id/following"
	followers = "/2/users/:id/followers"
	unfollow  = "/2/users/:source_user_id/following/:target_user_id"
)

func createFollowingURL(id string) string {
	return common.APIurl + strings.Replace(following, ":id", id, 1)
}

func createFollowersURL(id string) string {
	return common.APIurl + strings.Replace(followers, ":id", id, 1)
}

func createUnfollowURL(sourceId, targetId string) string {
	url := strings.Replace(unfollow, ":source_user_id", sourceId, 1)
	return common.APIurl + strings.Replace(url, ":target_user_id", targetId, 1)
}
