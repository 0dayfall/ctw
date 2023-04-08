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
	resp, err := httphandler.MakeRequest(req)
	defer resp.Body.Close()

	if !httphandler.IsResponseOK(resp) {
		return
	}
	json.NewDecoder(resp.Body).Decode(blockResponse)
	return
}

func BlockUserId(id string) (blockResponse BlockResponse, err error) {
	url := createBlockURL(id)
	req := httphandler.CreateGetRequest(url)
	resp, err := httphandler.MakeRequest(req)
	defer resp.Body.Close()

	if !httphandler.IsResponseOK(resp) {
		return
	}
	json.NewDecoder(resp.Body).Decode(blockResponse)
	return
}

type DeleteBlockResponse struct{}

func DeleteBlock(id string, blockedId string) (deleteBlockResponse DeleteBlockResponse, err error) {
	url := createDeleteBlockURL(id, blockedId)
	req := httphandler.CreateGetRequest(url)
	response, err := httphandler.MakeRequest(req)
	defer response.Body.Close()

	if !httphandler.IsResponseOK(response) {
		return
	}
	if err := json.NewDecoder(response.Body).Decode(&deleteBlockResponse); err != nil {
		log.Println(err)
	}
	return
}
