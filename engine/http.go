package engine

import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return string(body), nil
}
