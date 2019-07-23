package presets

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"

	kubermaticv1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"
)

// Presets specifies default presets for supported providers
type Presets struct {
	Digitalocean Digitalocean `json:"digitalocean,omitempty"`
	Hetzner      Hetzner      `json:"hetzner,omitempty"`
	Azure        Azure        `json:"azure,omitempty"`
	VSphere      VSphere      `json:"vsphere,omitempty"`
	AWS          AWS          `json:"aws,omitempty"`
	Openstack    Openstack    `json:"openstack,omitempty"`
	Packet       Packet       `json:"packet,omitempty"`
	GCP          GCP          `json:"gcp,omitempty"`
	Fake         Fake         `json:"fake,omitempty"`
}

type Digitalocean struct {
	Credentials []DigitaloceanCredentials `json:"credentials,omitempty"`
}

type Hetzner struct {
	Credentials []HetznerCredentials `json:"credentials,omitempty"`
}

type Azure struct {
	Credentials []AzureCredentials `json:"credentials,omitempty"`
}

type VSphere struct {
	Credentials []VSphereCredentials `json:"credentials,omitempty"`
}

type AWS struct {
	Credentials []AWSCredentials `json:"credentials,omitempty"`
}

type Openstack struct {
	Credentials []OpenstackCredentials `json:"credentials,omitempty"`
}

type Packet struct {
	Credentials []PacketCredentials `json:"credentials,omitempty"`
}

type GCP struct {
	Credentials []GCPCredentials `json:"credentials,omitempty"`
}

type Fake struct {
	Credentials []FakeCredentials `json:"credentials,omitempty"`
}

// DigitaloceanCredential defines Digitalocean credential
type DigitaloceanCredentials struct {
	Name  string `json:"name"`
	Token string `json:"token"` // Token is used to authenticate with the DigitalOcean API.
}

type HetznerCredentials struct {
	Name  string `json:"name"`
	Token string `json:"token"` // Token is used to authenticate with the Hetzner API.
}

type AzureCredentials struct {
	Name           string `json:"name"`
	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`

	ResourceGroup  string `json:"resourceGroup,omitempty"`
	VNetName       string `json:"vnet,omitempty"`
	SubnetName     string `json:"subnet,omitempty"`
	RouteTableName string `json:"routeTable,omitempty"`
	SecurityGroup  string `json:"securityGroup,omitempty"`
}

// VSphereCredentials credentials represents a credential for accessing vSphere
type VSphereCredentials struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`

	VMNetName string `json:"vmNetName,omitempty"`
}

type AWSCredentials struct {
	Name            string `json:"name"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`

	VPCID               string `json:"vpcId,omitempty"`
	RouteTableID        string `json:"routeTableId,omitempty"`
	InstanceProfileName string `json:"instanceProfileName,omitempty"`
	SecurityGroupID     string `json:"securityGroupID,omitempty"`
}

// OpenstackCredentials specifies access data to an openstack cloud.
type OpenstackCredentials struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
	TenantID string `json:"tenantID"`
	Domain   string `json:"domain"`

	Network        string `json:"network,,omitempty"`
	SecurityGroups string `json:"securityGroups,omitempty"`
	FloatingIPPool string `json:"floatingIpPool,omitempty"`
	RouterID       string `json:"routerID,omitempty"`
	SubnetID       string `json:"subnetID,omitempty"`
}

// PacketCredentials specifies access data to a Packet cloud.
type PacketCredentials struct {
	Name      string `json:"name"`
	APIKey    string `json:"apiKey"`
	ProjectID string `json:"projectId"`

	BillingCycle string `json:"billingCycle,omitempty"`
}

// GCPCredentials specifies access data to GCP.
type GCPCredentials struct {
	Name           string `json:"name"`
	ServiceAccount string `json:"serviceAccount"`

	Network    string `json:"network,omitempty"`
	Subnetwork string `json:"subnetwork,omitempty"`
}

// FakeCredentials defines fake credential for tests
type FakeCredentials struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// loadPresets loads the custom presets for supported providers
func loadPresets(path string) (*Presets, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}

	s := struct {
		Presets *Presets `json:"presets"`
	}{}

	err = yaml.UnmarshalStrict(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s.Presets, nil
}

// Manager is a object to handle presets from a predefined config
type Manager struct {
	presets *Presets
}

func New() *Manager {
	presets := &Presets{}
	return &Manager{presets: presets}
}

func NewWithPresets(presets *Presets) *Manager {
	return &Manager{presets: presets}
}

// NewFromFile returns a instance of manager with the credentials loaded from the given paths
func NewFromFile(credentialsFilename string) (*Manager, error) {
	var presets *Presets
	var err error
	if len(credentialsFilename) > 0 {
		presets, err = loadPresets(credentialsFilename)
		if err != nil {
			return nil, fmt.Errorf("failed to load presets from %s: %v", credentialsFilename, err)
		}
	}
	if presets == nil {
		presets = &Presets{
			Digitalocean: Digitalocean{},
			VSphere:      VSphere{},
			Openstack:    Openstack{},
			Hetzner:      Hetzner{},
			GCP:          GCP{},
			Azure:        Azure{},
			AWS:          AWS{},
			Packet:       Packet{},
			Fake:         Fake{},
		}
	}
	return &Manager{presets: presets}, nil
}

func (m *Manager) GetPresets() *Presets {
	return m.presets
}

func (m *Manager) SetCloudCredentials(credentialName string, cloud kubermaticv1.CloudSpec, dc *kubermaticv1.Datacenter) (*kubermaticv1.CloudSpec, error) {

	if cloud.VSphere != nil {
		return m.setVsphereCredentials(credentialName, cloud)
	}
	if cloud.Openstack != nil {
		return m.setOpenStackCredentials(credentialName, cloud, dc)
	}
	if cloud.Azure != nil {
		return m.setAzureCredentials(credentialName, cloud)
	}
	if cloud.Digitalocean != nil {
		return m.setDigitalOceanCredentials(credentialName, cloud)
	}
	if cloud.Packet != nil {
		return m.setPacketCredentials(credentialName, cloud)
	}
	if cloud.Hetzner != nil {
		return m.setHetznerCredentials(credentialName, cloud)
	}
	if cloud.AWS != nil {
		return m.setAWSCredentials(credentialName, cloud)
	}
	if cloud.GCP != nil {
		return m.setGCPCredentials(credentialName, cloud)
	}
	if cloud.Fake != nil {
		return m.setFakeCredentials(credentialName, cloud)
	}

	return nil, fmt.Errorf("can not find provider to set credentials")
}

func emptyCredentialListError(provider string) error {
	return fmt.Errorf("can not find any credential for %s provider", provider)
}

func noCredentialError(credentialName string) error {
	return fmt.Errorf("can not find %s credential", credentialName)
}

func (m *Manager) setFakeCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.Fake.Credentials == nil {
		return nil, emptyCredentialListError("Fake")
	}
	for _, credential := range m.presets.Fake.Credentials {
		if credentialName == credential.Name {
			cloud.Fake.Token = credential.Token
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setGCPCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.GCP.Credentials == nil {
		return nil, emptyCredentialListError("GCP")
	}
	for _, credential := range m.presets.GCP.Credentials {
		if credentialName == credential.Name {
			cloud.GCP.ServiceAccount = credential.ServiceAccount

			cloud.GCP.Network = credential.Network
			cloud.GCP.Subnetwork = credential.Subnetwork
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setAWSCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.AWS.Credentials == nil {
		return nil, emptyCredentialListError("AWS")
	}
	for _, credential := range m.presets.AWS.Credentials {
		if credentialName == credential.Name {
			cloud.AWS.AccessKeyID = credential.AccessKeyID
			cloud.AWS.SecretAccessKey = credential.SecretAccessKey

			cloud.AWS.InstanceProfileName = credential.InstanceProfileName
			cloud.AWS.RouteTableID = credential.RouteTableID
			cloud.AWS.SecurityGroupID = credential.SecurityGroupID
			cloud.AWS.VPCID = credential.VPCID
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setHetznerCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.Hetzner.Credentials == nil {
		return nil, emptyCredentialListError("Hetzner")
	}
	for _, credential := range m.presets.Hetzner.Credentials {
		if credentialName == credential.Name {
			cloud.Hetzner.Token = credential.Token
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setPacketCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.Packet.Credentials == nil {
		return nil, emptyCredentialListError("Packet")
	}
	for _, credential := range m.presets.Packet.Credentials {
		if credentialName == credential.Name {
			cloud.Packet.ProjectID = credential.ProjectID
			cloud.Packet.APIKey = credential.APIKey

			cloud.Packet.BillingCycle = credential.BillingCycle
			if len(credential.BillingCycle) == 0 {
				cloud.Packet.BillingCycle = "hourly"
			}

			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setDigitalOceanCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.Digitalocean.Credentials == nil {
		return nil, emptyCredentialListError("Digitalocean")
	}
	for _, credential := range m.presets.Digitalocean.Credentials {
		if credentialName == credential.Name {
			cloud.Digitalocean.Token = credential.Token
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setAzureCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.Azure.Credentials == nil {
		return nil, emptyCredentialListError("Azure")
	}
	for _, credential := range m.presets.Azure.Credentials {
		if credentialName == credential.Name {
			cloud.Azure.TenantID = credential.TenantID
			cloud.Azure.ClientSecret = credential.ClientSecret
			cloud.Azure.ClientID = credential.ClientID
			cloud.Azure.SubscriptionID = credential.SubscriptionID

			cloud.Azure.ResourceGroup = credential.ResourceGroup
			cloud.Azure.RouteTableName = credential.RouteTableName
			cloud.Azure.SecurityGroup = credential.SecurityGroup
			cloud.Azure.SubnetName = credential.SubnetName
			cloud.Azure.VNetName = credential.VNetName
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setOpenStackCredentials(credentialName string, cloud kubermaticv1.CloudSpec, dc *kubermaticv1.Datacenter) (*kubermaticv1.CloudSpec, error) {
	if m.presets.Openstack.Credentials == nil {
		return nil, emptyCredentialListError("Openstack")
	}
	for _, credential := range m.presets.Openstack.Credentials {
		if credentialName == credential.Name {
			cloud.Openstack.Username = credential.Username
			cloud.Openstack.Password = credential.Password
			cloud.Openstack.Domain = credential.Domain
			cloud.Openstack.Tenant = credential.Tenant
			cloud.Openstack.TenantID = credential.TenantID

			cloud.Openstack.SubnetID = credential.SubnetID
			cloud.Openstack.Network = credential.Network
			cloud.Openstack.FloatingIPPool = credential.FloatingIPPool

			if cloud.Openstack.FloatingIPPool == "" && dc.Spec.Openstack != nil && dc.Spec.Openstack.EnforceFloatingIP {
				return nil, fmt.Errorf("preset error, no floating ip pool specified for OpenStack")
			}

			cloud.Openstack.RouterID = credential.RouterID
			cloud.Openstack.SecurityGroups = credential.SecurityGroups
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}

func (m *Manager) setVsphereCredentials(credentialName string, cloud kubermaticv1.CloudSpec) (*kubermaticv1.CloudSpec, error) {
	if m.presets.VSphere.Credentials == nil {
		return nil, emptyCredentialListError("Vsphere")
	}
	for _, credential := range m.presets.VSphere.Credentials {
		if credentialName == credential.Name {
			cloud.VSphere.Password = credential.Password
			cloud.VSphere.Username = credential.Username

			cloud.VSphere.VMNetName = credential.VMNetName
			return &cloud, nil
		}
	}
	return nil, noCredentialError(credentialName)
}
