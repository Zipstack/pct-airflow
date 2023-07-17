package scp

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	HTTPClient    *http.Client
	SCPClient     ssh.ClientConfig
	Host          string
	Port          int
	Username      string
	Password      string
	AirflowRepo   string
	RepoOwner     string
	TargetBranch  string
	Authorization string
}

func NewClient(host string, port int, username string, password string, airflow_repo string, repo_owner string, target_branch string, auth string) (*Client, error) {
	//var hostKey ssh.PublicKey

	c := Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
		SCPClient: ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,

		AirflowRepo:   airflow_repo,
		RepoOwner:     repo_owner,
		TargetBranch:  target_branch,
		Authorization: auth,
	}
	return &c, nil
}

func (c *Client) sshUpdate(method string, url string, body []byte, headers map[string]string) error {

	//sshClientConfig := c.SSHClinet

	//client = NewConfigurer(c.Host, sshClientConfig).Create()
	return nil
}

func (c *Client) doRequest(method string, url string, body []byte, headers map[string]string) ([]byte, int, string, map[string][]string, error) {
	payload := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, 500, "500 Internal Server Error", nil, err
	}

	req.Header.Add("Authorization", c.Authorization)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "PCT")
	req.Header.Add("Content-Type", "application/json")

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 500, "500 Internal Server Error", nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 500, "500 Internal Server Error", nil, err
	}

	defer res.Body.Close()
	return b, res.StatusCode, res.Status, res.Header, nil
}
