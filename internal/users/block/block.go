package user

import (
	"encoding/json"
	"log"
	"strings"

	common "github.com/0dayfall/ctw/internal/data"
	"github.com/0dayfall/ctw/internal/httphandler"
)

const (
	blockURL       = "/2/users/:id/blocking"
	deleteBlockURL = "/2/users/:source_user_id/blocking/:target_user_id"
)

func createBlockURL(id string) string {
	return common.APIurl + strings.Replace(blockURL, ":id", id, 1)
}

func createDeleteBlockURL(sourceId, targetId string) string {
	url := strings.Replace(deleteBlockURL, ":source_user_id", sourceId, 1)
	url = strings.Replace(url, ":target_user_id", targetId, 1)
	url = common.APIurl + url
	return url
}

type BlockResponse struct {
}

func BlockUserIdX(id string) (blockResponse BlockResponse, err error) {
	url := createBlockURL(id)
	req := httphandler.CreateGetRequest(url)
	httpResponse, err := httphandler.MakeRequest(req)
	defer httphandler.CloseBody(httpResponse.Body)

	if !httphandler.IsResponseOK(httpResponse) {
		return
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&blockResponse); err != nil {
		log.Println(err)
	}
	return
}

func BlockUserId(id string) (blockResponse BlockResponse, err error) {
	url := createBlockURL(id)
	req := httphandler.CreateGetRequest(url)
	httpResponse, err := httphandler.MakeRequest(req)
	defer httphandler.CloseBody(httpResponse.Body)

	if !httphandler.IsResponseOK(httpResponse) {
		return
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&blockResponse); err != nil {
		log.Println(err)
	}
	return
}

type DeleteBlockResponse struct{}

func DeleteBlock(id string, blockedId string) (deleteBlockResponse DeleteBlockResponse, err error) {
	url := createDeleteBlockURL(id, blockedId)
	req := httphandler.CreateGetRequest(url)
	httpResponse, err := httphandler.MakeRequest(req)
	defer httphandler.CloseBody(httpResponse.Body)

	if !httphandler.IsResponseOK(httpResponse) {
		return
	}
	if err := json.NewDecoder(httpResponse.Body).Decode(&deleteBlockResponse); err != nil {
		log.Println(err)
	}
	return
}
