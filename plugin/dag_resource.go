package plugin

import (
	"context"
	"fmt"

	"os"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"
	scp "github.com/zipstack/pct-provider-airflow/ssh"
)

type dagModel struct {
	Repo     string `pctsdk:"repo"`
	DagPath  string `pctsdk:"dag_path"`
	CommitId string `pctsdk:"commit_id"`
}

type dagResource struct {
	Client *scp.Client
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &dagResource{}
)

// Helper function to return a resource service instance.
func NewDagResource() schema.ResourceService {
	return &dagResource{}
}

func (r *dagResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_dag",
	}
}

func (r *dagResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
	if req.ResourceData == "" {
		return schema.ErrorResponse(fmt.Errorf("no data provided to configure resource"))
	}

	var creds map[string]any
	err := fwhelpers.Decode(req.ResourceData, &creds)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	client, err := scp.NewClient(creds["host"].(string), creds["port"].(int), creds["username"].(string), creds["password"].(string), creds["airflow_repo"].(string), creds["repo_owner"].(string), creds["target_branch"].(string), creds["authorization"].(string))
	if err != nil {
		return schema.ErrorResponse(fmt.Errorf("malformed data provided to configure resource"))
	}

	r.Client = client

	return &schema.ServiceResponse{}
}

func (r *dagResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Dag for Airflow",
		Attributes: map[string]schema.Attribute{
			"repo": &schema.StringAttribute{
				Description: "Repo for airflow",
				Required:    true,
			},
			"dag_path": &schema.StringAttribute{
				Description: "path to dags in airflow repo",
				Required:    true,
			},
			"commit_id": &schema.StringAttribute{
				Description: "Commit Id",
				Computed:    true,
			},
		},
	}

	sEnc, err := fwhelpers.Encode(s)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		SchemaContents: sEnc,
	}
}
func (r *dagResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	logger := fwhelpers.GetLogger()
	logger.Printf("coming here")
	fmt.Printf("state data in provider %#v\n", req)
	// r.Client.Create()
	ClientConfig := r.Client.SCPClient

	client := scp.NewSCPClient(r.Client.Host, &ClientConfig)

	err := client.Connect()

	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return nil
	}

	// Open a file
	f, _ := os.Open("/")

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	err = client.CopyFromFile(context.Background(), *f, "/", "0655")

	if err != nil {
		fmt.Println("Error while copying file ", err)
	}
	res := schema.ServiceResponse{}
	res.StateID = "sdajadskjsdja"
	return &res

}
func (r *dagResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()
	var state dagModel

	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		dags_details, err := r.Client.ReadDags(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.DagPath = dags_details.DagPath

		res.StateID = dags_details.CommitId
	} else {
		// No previous state exists.
		res.StateID = ""
	}

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}
	res.StateContents = stateEnc

	return &res

}

func (r *dagResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	logger := fwhelpers.GetLogger()
	logger.Print("coming here")
	fmt.Printf("Unsupported data product stage: \n")
	return nil
}
func (r *dagResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	logger := fwhelpers.GetLogger()
	logger.Print("coming here")
	fmt.Printf("Unsupported data product stage: \n")
	return nil
}
