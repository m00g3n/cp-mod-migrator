package v294

import "k8s.io/apimachinery/pkg/util/json"

type RfcAndLdap ProxyCfg

func (p *RfcAndLdap) UnmarshalJSON(text []byte) error {
	var result ProxyCfg
	if err := json.Unmarshal(text, &result); err != nil {
		return err
	}

	if result.Port == 0 {
		result.Port = DefaultPortRfcAndLdap
	}

	*p = RfcAndLdap(result)
	return nil
}

func (p *RfcAndLdap) MarshalJSON() ([]byte, error) {
	if p.Port == 0 {
		p.Port = DefaultPortHTTP
	}

	proxyCfg := ProxyCfg(*p)
	return json.Marshal(&proxyCfg)
}
