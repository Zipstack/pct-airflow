package plugin

import (
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"
	scp "github.com/zipstack/pct-provider-airflow/ssh"
)

// Provider implementation.
type Provider struct {
	Client           *scp.Client
	ResourceServices map[string]string
}

// Model maps the provider state as per schema.
type ProviderModel struct {
	Host     string `pctsdk:"host"`
	Port     int    `pctsdk:"port"`
	Username string `pctsdk:"username"`
	Password string `pctsdk:"password"`

	AirflowRepo   string `pctsdk:"airflow_repo"`
	RepoOwner     string `pctsdk:"repo_owner"`
	TargetBranch  string `pctsdk:"target_branch"`
	Authorization string `pctsdk:"authorization"`
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ schema.ProviderService = &Provider{}
)

// Helper function to return a provider service instance.
func NewProvider() schema.ProviderService {
	return &Provider{}
}

// Metadata returns the provider type name.
func (p *Provider) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: "airflow",
	}
}

// Schema defines the provider-level schema for configuration data.
func (p *Provider) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Airflow provider plugin",
		Attributes: map[string]schema.Attribute{
			"host": &schema.StringAttribute{
				Description: "URI for Airflow API. May also be provided via AIRFLOW_HOST environment variable.",
				Required:    true,
			},
			"port": &schema.IntAttribute{
				Description: "port number for airflow server",
				Required:    true,
			},
			"username": &schema.StringAttribute{
				Description: "ssh user name",
				Required:    true,
			},
			"password": &schema.StringAttribute{
				Description: "ssh password",
				Required:    true,
				Sensitive:   true,
			},
			"airflow_repo": &schema.StringAttribute{
				Description: "airflow repo",
				Required:    true,
			},
			"repo_owner": &schema.StringAttribute{
				Description: "repo owner",
				Required:    true,
			},
			"target_branch": &schema.StringAttribute{
				Description: "target branch",
				Required:    true,
			},
			"authorization": &schema.StringAttribute{
				Description: "authorization",
				Required:    true,
				Sensitive:   true,
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

func (p *Provider) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
	var pm ProviderModel
	err := fwhelpers.UnpackModel(req.ConfigContents, &pm)
	if err != nil {
		return schema.ErrorResponse(err)
	}
	if pm.Host == "" || pm.Port == 0 || pm.Username == "" || pm.Password == "" {
		return schema.ErrorResponse(fmt.Errorf(
			"invalid host or credentials received.\n" +
				"Provider is unable to create Airflow .",
		))
	}

	if p.Client == nil {
		client, err := scp.NewClient(
			pm.Host,
			pm.Port,
			pm.Username,
			pm.Password,
			pm.RepoOwner,
			pm.RepoOwner,
			pm.TargetBranch,
			pm.Authorization)
		if err != nil {
			return schema.ErrorResponse(err)
		}
		p.Client = client
	}

	// Make API creds available for Resource type Configure methods.
	creds := map[string]any{
		"host":          pm.Host,
		"port":          pm.Port,
		"username":      pm.Username,
		"password":      pm.Password,
		"airflow_repo":  pm.AirflowRepo,
		"repo_owner":    pm.RepoOwner,
		"target_branch": pm.TargetBranch,
		"authorization": pm.Authorization,
	}
	cEnc, err := fwhelpers.Encode(creds)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		ResourceData: cEnc,
	}
}

func (p *Provider) Resources() *schema.ServiceResponse {
	return &schema.ServiceResponse{
		ResourceServices: p.ResourceServices,
	}
}

func (p *Provider) UpdateResourceServices(resServices map[string]string) {
	if resServices != nil {
		p.ResourceServices = resServices
	}
}
