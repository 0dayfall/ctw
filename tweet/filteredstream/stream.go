package filteredstream

func CreateUrl() string {
	return "https://api.twitter.com/2/tweets/sample/stream"
}

/*func createRequest() *http.Request {
	req, err := http.NewRequest("GET", CreateUrl(), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	return req
}

func makeRequest() *http.Response {
	client := http.Client{}
	response, err := client.Do(createRequest())
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		log.Fatal(response.Status)
	}
	return response
}

func SampleStream() {
	response := makeRequest()
	var sampleStreamResponse StreamResponse
	err := json.NewDecoder(response.Body).Decode(&sampleStreamResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(sampleStreamResponse)
}
*/
