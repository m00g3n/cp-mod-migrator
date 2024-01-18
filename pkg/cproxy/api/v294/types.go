package v294

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type HighAvailabilityMode string

const (
	HighAvailabilityModeOff       HighAvailabilityMode = "off"
	HighAvailabilityModePath      HighAvailabilityMode = "path"
	HighAvailabilityModeSubdomain HighAvailabilityMode = "subdomain"
)

type AuditLogMode string

const (
	AuditLogModeConsole AuditLogMode = "console"
	AuditLogModeService AuditLogMode = "service"
)

type TenantMode string

const (
	TenantModeDedicated = "dedicated"
	TenantModeShared    = "shared"
)

type Config struct {
	HighAvailabilityMode HighAvailabilityMode `json:"highAvailabilityMode"`
	Integration          Integration          `json:"integration"`
	ConnectivityService  ConnectivityService  `json:"connectivityService,omitempty"`
	MultiRegionMode      MultiRegionMode      `json:"multiRegionMode"`
	Servers              Servers              `json:"servers"`
	SubaccountID         string               `json:"subaccountId"`
	SubaccountSubdomain  string               `json:"subaccountSubdomain"`
	TenantMode           TenantMode           `json:"tenantMode"`
	ServiceChannels      ServiceChannels      `json:"serviceChannels"`
}

type Spec struct {
	Config Config `json:"config"`
}

type AuditLog struct {
	Mode                  AuditLogMode `json:"mode"`
	ServiceCredentialsKey string       `json:"serviceCredentialsKey,omitempty"`
}

type ConnectivityService struct {
	ServiceCredentialsKey string `json:"serviceCredentialsKey,omitempty"`
}

type Integration struct {
	AuditLog            AuditLog            `json:"auditlog"`
	ConnectivityService ConnectivityService `json:"connectivityService,omitempty"`
}

//+kubebuilder:object:root=true

type ConnectivityProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectivityProxy `json:"items"`
}

//+kubebuilder:object:root=true

type ConnectivityProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Spec `json:"spec"`
}

func (c *ConnectivityProxy) Encode() (string, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(c); err != nil {
		return "", err
	}
	data := b.Bytes()
	encode := base64.StdEncoding.EncodeToString(data)
	return encode, nil
}

type MultiRegionMode struct {
	ConfigMapName string `json:"configMapName,omitempty"`
	Enabled       bool   `json:"enabled"`
}

type Servers struct {
	BusinessDataTunnel BusinessDataTunnel `json:"businessDataTunnel"`
	Proxy              Proxy              `json:"proxy"`
}

type Proxy struct {
	Authorization Authorization `json:"authorization"`
	HTTP          HTTP          `json:"http"`
	RfcAndLdap    RfcAndLdap    `json:"rfcAndLdap"`
	Socks5        Socks5        `json:"socks5"`
}

type Authorization struct {
	OAuth OAuth `json:"oauth"`
}

type OAuth struct {
	AllowedClientId string `json:"allowedClientId"`
}

type BusinessDataTunnel struct {
	ExternalHost string `json:"externalHost"`
	ExternalPort int    `json:"externalPort"`
}

type ProxyCfg struct {
	Enabled                  bool `json:"enabled"`
	AllowRemoteConnection    bool `json:"allowRemoteConnections"`
	EnableProxyAuthorization bool `json:"enableProxyAuthorization"`
	Port                     int  `json:"port"`
}

//go:generate go run ../../../../cmd/generators/proxy-conf-type-gen/main.go -type-name HTTP -port-number 2003
//go:generate go run ../../../../cmd/generators/proxy-conf-type-gen/main.go -type-name Socks5 -port-number 2004
//go:generate go run ../../../../cmd/generators/proxy-conf-type-gen/main.go -type-name RfcAndLdap -port-number 2001

type ServiceChannels struct {
	Enabled bool `json:"enabled"`
}

// +kubebuilder:object:generate=true
// +groupName=operator.kyma-project.io

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "connectivityproxy.sap.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(&ConnectivityProxy{}, &ConnectivityProxyList{})
}
