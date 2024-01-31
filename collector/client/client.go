package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"io/ioutil"
	"net"
	"net/http"
	"powerstore/utils"
	"time"
)

type Client struct {
	IP       string
	username string
	password string
	version  string
	baseUrl  string
	http     *http.Client
	token    string
	logger   log.Logger
}

func NewClient(config utils.Storage, logger log.Logger) (*Client, error) {
	if config.Ip == "" || config.User == "" || config.Password == "" || config.Version == "" {
		return nil, errors.New("please check config file ,Some parameters are null")
	}
	baseUrl := "https://" + config.Ip + "/api/rest/"
	var httpClient *http.Client
	httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 60 * time.Second,
	}
	client := &Client{
		IP:       config.Ip,
		username: config.User,
		password: config.Password,
		version:  config.Version,
		baseUrl:  baseUrl,
		http:     httpClient,
		logger:   logger,
	}
	return client, nil
}

func (c *Client) InitLogin() error {
	reqUrl := c.baseUrl + "cluster?select=*"
	request, err := http.NewRequest("GET", reqUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		return err
	}
	request.SetBasicAuth(c.username, c.password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	response, err := c.http.Do(request)
	if err != nil {
		level.Warn(c.logger).Log("msg", "Request URL error!")
		return err
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		c.token = response.Header.Get("Dell-Emc-Token")
		return nil
	default:
		body, err := ioutil.ReadAll(response.Body)
		level.Warn(c.logger).Log("msg", "get token error", "err", err)
		return errors.New("get token error: " + string(body))
	}
}

func (c *Client) getResource(method, uri, body string) (string, error) {
	reqUrl := c.baseUrl + uri
	request, err := http.NewRequest(method, reqUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(c.username, c.password)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("DELL-EMC-TOKEN", c.token)

	response, err := c.http.Do(request)
	if err != nil {
		level.Warn(c.logger).Log("msg", "Request URL error!")
		return "", err
	}

	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusPartialContent:
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", errors.New("get resource error: " + string(body))
		}
		return string(body), nil
	case http.StatusUnauthorized, http.StatusFound:
		level.Warn(c.logger).Log("msg", "authentication token is invalid, relogin...", "err", err)
		err = c.InitLogin()
		if err != nil {
			level.Warn(c.logger).Log("msg", "init auth error", "err", err)
			return "", err
		} else {
			return c.getResource(method, uri, body)
		}
	default:
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", errors.New("get resource error ReadAll err is not nil: " + string(body))
		}
		return "", errors.New("get resource error ReadAll err is nil: " + string(body))
	}

}
