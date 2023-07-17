package scp

import (
	"encoding/json"
)

type DagModel struct {
	Repo     string `pctsdk:"repo"`
	DagPath  string `pctsdk:"dagPath"`
	CommitId string `pctsdk:"commitId"`
}

type GitModel struct {
	Sha string `pctsdk:"sha"`
}

func (c *Client) ReadDags(commitId string) (DagModel, error) {
	dagModel := DagModel{}
	// dagModel.repo = "github.com/zipstack/pct-provider-airflow/"
	// dagModel.dagPath = "dags"
	// dagModel.commitId = "adaddadsas"

	latest_commit, _ := c.ReadLatestCommit()
	if latest_commit != commitId {
		dagModel.CommitId = latest_commit
	}

	return dagModel, nil
}

func (c *Client) ReadLatestCommit() (string, error) {

	method := "GET"
	url := "https://api.github.com/repos/" + c.RepoOwner + "/" + c.AirflowRepo + "/commits/" + c.TargetBranch

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	gitModel := GitModel{}
	if err != nil {
		return "", err
	}

	if statusCode == 200 {
		err = json.Unmarshal(b, &gitModel)
		return gitModel.Sha, err
	}

	return "", nil
}

func (c *Client) CreateDags(commitId string) (DagModel, error) {
	//c.sshUpdate()

	return DagModel{}, nil
}
