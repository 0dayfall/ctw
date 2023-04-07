package user

import (
	"encoding/json"
	"io"
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

func (b *BlockResponse) decode(body io.ReadCloser) *BlockResponse {
	if err := json.NewDecoder(body).Decode(b); err != nil {
		log.Println(err)
	}
	return b
}

func BlockUserIdX(id string) BlockResponse {
	url := createBlockURL(id)
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()
	if !httphandler.IsResponseOK(response) {
		return BlockResponse{}
	}
	var blockResponse BlockResponse
	if err := json.NewDecoder(response.Body).Decode(&blockResponse); err != nil {
		log.Println(err)
	}
	return blockResponse
}

func BlockUserId(id string) BlockResponse {
	url := createBlockURL(id)
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()
	if !httphandler.IsResponseOK(response) {
		return BlockResponse{}
	}
	var blockResponse BlockResponse
	return *blockResponse.decode(response.Body)
}

type DeleteBlockResponse struct{}

func DeleteBlock(id string, blockedId string) DeleteBlockResponse {
	url := createDeleteBlockURL(id, blockedId)
	req := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()
	if !httphandler.IsResponseOK(response) {
		return DeleteBlockResponse{}
	}
	var deleteBlockResponse DeleteBlockResponse
	if err := json.NewDecoder(response.Body).Decode(&deleteBlockResponse); err != nil {
		log.Println(err)
	}
	return deleteBlockResponse
}
