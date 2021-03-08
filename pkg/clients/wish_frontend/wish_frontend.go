package wish_frontend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ContextLogic/pkg/autobots/config"
)

type WishFrontend struct {
	Client *http.Client
	Config *config.WishFrontendConfig
}

func New(config *config.WishFrontendConfig) *WishFrontend {
	return &WishFrontend{
		&http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		config,
	}
}

func (w *WishFrontend) Post(h http.Header, data []byte, uri string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/%s", w.Config.Host, uri)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header = h
	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil

}
