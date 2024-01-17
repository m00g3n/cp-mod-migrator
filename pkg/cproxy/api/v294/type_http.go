package v294

import "k8s.io/apimachinery/pkg/util/json"

type HTTP ProxyCfg

const DefaultPortHTTP = 2003

func (p *HTTP) UnmarshalJSON(text []byte) error {
	var result ProxyCfg
	if err := json.Unmarshal(text, &result); err != nil {
		return err
	}

	if result.Port == 0 {
		result.Port = DefaultPortHTTP
	}

	*p = HTTP(result)
	return nil
}

func (p *HTTP) MarshalJSON() ([]byte, error) {
	if p.Port == 0 {
		p.Port = DefaultPortHTTP
	}

	proxyCfg := ProxyCfg(*p)
	return json.Marshal(&proxyCfg)
}
