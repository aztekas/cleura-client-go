package cleura

// Shoot cluster response data model.
type ShootClusterResponse struct {
	Metadata MetadataFields `json:"metadata"`
	Spec     SpecFields     `json:"spec"`
	Status   StatusFields   `json:"status"`
}

type ShootClusterCreateResponse struct {
	Shoot ShootClusterCreateResponseConfig `json:"shoot"`
}

type ShootClusterCreateResponseConfig struct {
	Name        string             `json:"name"`
	UID         string             `json:"uid"`
	Kubernetes  KubernetesDetails  `json:"kubernetes"`
	Provider    ProviderDetails    `json:"provider"`
	Purpose     string             `json:"purpose"`
	Region      string             `json:"region"`
	Hibernation HibernationDetails `json:"hibernation"`
}

type MetadataFields struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

type SpecFields struct {
	Purpose     string             `json:"purpose"`
	Region      string             `json:"region"`
	Provider    ProviderDetails    `json:"provider"`
	Kubernetes  KubernetesDetails  `json:"kubernetes"`
	Hibernation HibernationDetails `json:"hibernation"`
}

type HibernationDetails struct {
	Enabled                      bool                          `json:"enabled"`
	HibernationResponseSchedules []HibernationResponseSchedule `json:"schedules,omitempty"`
}

type HibernationResponseSchedule struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Location string `json:"location"`
}

type KubernetesDetails struct {
	Version string `json:"version"`
}

type StatusFields struct {
	Conditions          []Condition          `json:"conditions"`
	Hibernated          bool                 `json:"hibernated"`
	AdvertisedAddresses []AdvertisedAddress  `json:"advertisedAddresses"`
	LastOperation       LastOperationDetails `json:"lastOperation"`
}

type LastOperationDetails struct {
	Progress int16  `json:"progress"`
	State    string `json:"state"`
	Type     string `json:"type"`
}

type Condition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AdvertisedAddress struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Shoot cluster request data model.
type ShootClusterRequest struct {
	Shoot ShootClusterRequestConfig `json:"shoot"`
}

type ShootClusterRequestConfig struct {
	Name              string                `json:"name,omitempty"`
	KubernetesVersion *K8sVersion           `json:"kubernetes,omitempty"`
	Provider          *ProviderDetails      `json:"provider,omitempty"`
	Hibernation       *HibernationSchedules `json:"hibernation,omitempty"`
}
type K8sVersion struct {
	Version string `json:"version"`
}

type ProviderDetails struct {
	InfrastructureConfig InfrastructureConfigDetails `json:"infrastructureConfig"`
	Workers              []Worker                    `json:"workers"`
}

type InfrastructureConfigDetails struct {
	FloatingPoolName string `json:"floatingPoolName"`
	// Networks *WorkerNetwork `json:"networks,omitempty"`
}

/*
type WorkerNetwork struct {
	WorkersCIDR string `json:"workers,omitempty"`
}
*/

// Worker.
type Worker struct {
	Name     string         `json:"name,omitempty"`
	Minimum  int16          `json:"minimum,omitempty"`
	Maximum  int16          `json:"maximum,omitempty"`
	MaxSurge int16          `json:"maxSurge,omitempty"`
	Machine  MachineDetails `json:"machine"`
	Volume   VolumeDetails  `json:"volume"`
}

type MachineDetails struct {
	Type  string       `json:"type"`
	Image ImageDetails `json:"image"`
}
type ImageDetails struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type VolumeDetails struct {
	Size string `json:"size"`
}

type HibernationSchedules struct {
	HibernationSchedules []HibernationSchedule `json:"schedules,omitempty"`
}

type HibernationSchedule struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Worker groups.

type WorkerGroupRequest struct {
	Worker Worker `json:"worker"`
}

// Gardener Cloud Profiles

type CPMachineType struct {
	Cpu          string `json:"cpu"`
	Gpu          string `json:"gpu"`
	Memory       string `json:"memory"`
	Name         string `json:"name"`
	Usable       bool   `json:"usable"`
	Architecture string `json:"architecture"`
}
type CPMachineImage struct {
	Name     string      `json:"name"`
	Versions []CPVersion `json:"versions"`
}
type CPVersion struct {
	Version        string `json:"version"`
	ExpirationDate string `json:"expirationDate"`
	Classification string `json:"classification"`
}
type CPKubernetes struct {
	Versions []CPVersion `json:"versions"`
}
type CloudProfileSpec struct {
	Kubernetes    CPKubernetes     `json:"kubernetes"`
	MachineImages []CPMachineImage `json:"machineImages"`
	MachineTypes  []CPMachineType  `json:"machineTypes"`
}
type CloudProfile struct {
	Name string           `json:"name"`
	Spec CloudProfileSpec `json:"spec"`
}

// Openstack.

// Openstack domains.
type OpenstackDomain struct {
	Area    OpenstackArea `json:"area"`
	Id      string        `json:"id"`
	Status  string        `json:"status"`
	Name    string        `json:"name"`
	Enabled bool          `json:"enabled"`
}

type OpenstackArea struct {
	Name    string            `json:"name"`
	Tag     string            `json:"tag"`
	Regions []OpenstackRegion `json:"regions"`
}

type OpenstackRegion struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Region string `json:"region"`
}

// Openstack projects.
type OpenstackProject struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DomainId    string `json:"domain_id"`
	Enabled     bool   `json:"enabled"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}
