package v294

// TODO: Rename the package to the exact version (2.11)

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
	CProxyDefaultCRName           string               = "connectivity-proxy"
	CProxyDefaultCRNamespace      string               = "kyma-system"
	CProxyMigratedAnnotation      string               = "connectivityproxy.sap.com/migrated"
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
	ConnectivityService  ConnectivityService  `json:"connectivityService"`
	MultiRegionMode      MultiRegionMode      `json:"multiRegionMode"`
	Servers              Servers              `json:"servers"`
	SubaccountID         string               `json:"subaccountId"`
	SubaccountSubdomain  string               `json:"subaccountSubdomain"`
	TenantMode           TenantMode           `json:"tenantMode"`
	ServiceChannels      ServiceChannels      `json:"serviceChannels"`
}

type Spec struct {
	Config       Config       `json:"config"`
	Deployment   Deployment   `json:"deployment"`
	Ingress      Ingress      `json:"ingress"`
	SecretConfig SecretConfig `json:"secretConfig"`
}

type AuditLog struct {
	Mode                  AuditLogMode `json:"mode"`
	ServiceCredentialsKey *string      `json:"serviceCredentialsKey,omitempty"`
}

type ConnectivityService struct {
	ServiceCredentialsKey string `json:"serviceCredentialsKey"`
}

type Integration struct {
	AuditLog            AuditLog             `json:"auditlog"`
	ConnectivityService *ConnectivityService `json:"connectivityService,omitempty"`
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
	ConfigMapName *string `json:"configMapName,omitempty"`
	Enabled       bool    `json:"enabled"`
}

type Servers struct {
	BusinessDataTunnel BusinessDataTunnel `json:"businessDataTunnel"`
	Proxy              Proxy              `json:"proxy"`
}

type Proxy struct {
	Authorization *Authorization `json:"authorization,omitempty"`
	HTTP          HTTP           `json:"http"`
	RfcAndLdap    RfcAndLdap     `json:"rfcAndLdap"`
	Socks5        Socks5         `json:"socks5"`
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
	Enabled                  bool  `json:"enabled"`
	AllowRemoteConnection    *bool `json:"allowRemoteConnections,omitempty"`
	EnableProxyAuthorization bool  `json:"enableProxyAuthorization"`
	Port                     int   `json:"port"`
}

//go:generate go run ../../../../cmd/generators/proxy-conf-type-gen/main.go -type-name HTTP -port-number 2003
//go:generate go run ../../../../cmd/generators/proxy-conf-type-gen/main.go -type-name Socks5 -port-number 2004
//go:generate go run ../../../../cmd/generators/proxy-conf-type-gen/main.go -type-name RfcAndLdap -port-number 2001

type ServiceChannels struct {
	Enabled bool `json:"enabled"`
}

type Deployment struct {
	RestartWatcher RestartWatcher `json:"restartWatcher"`
}

type RestartWatcher struct {
	Enabled bool `json:"enabled"`
}

type Ingress struct {
	ClassName ClassType  `json:"className"`
	Tls       IngressTls `json:"tls"`
	Timeouts  Timeouts   `json:"timeouts"`
	Istio     *Istio     `json:"istio,omitempty"`
}

type ClassType string

const (
	ClassTypeIstio ClassType = "istio"
	ClassTypeNginx ClassType = "nginx"
)

type IngressTls struct {
	SecretName string `json:"secretName"`
}

type Timeouts struct {
	Proxy TimeoutProxy `json:"proxy"`
}

type TimeoutProxy struct {
	Connect int `json:"connect"`
	Read    int `json:"read"`
	Send    int `json:"send"`
}

type Istio struct {
	Namespace string   `json:"namespace"`
	Gateway   Gateway  `json:"gateway"`
	Tls       IstioTls `json:"tls"`
}

type IstioTls struct {
	Ciphers []string `json:"ciphers"`
}

type Gateway struct {
	Selector Selector `json:"selector"`
}

type Selector struct {
	AdditionalProperties string `json:"additionalProperties"`
}

type SecretConfig struct {
	Integration SecretConfigIntegration `json:"integration"`
}

type SecretConfigIntegration struct {
	ConnectivityService ServiceSecretConfig  `json:"connectivityService"`
	AuditLogService     *ServiceSecretConfig `json:"auditlogService,omitempty"`
}

type ServiceSecretConfig struct {
	SecretName string  `json:"secretName"`
	SecretData *string `json:"secretData,omitempty"`
}

// Migrated - checks for annotation to verify if the CR has been migrated
func (c *ConnectivityProxy) Migrated() bool {
	_, ok := c.Annotations[CProxyMigratedAnnotation]
	return ok
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
