package v294

import "k8s.io/apimachinery/pkg/util/json"

type Socks5 ProxyCfg

func (p *Socks5) UnmarshalJSON(text []byte) error {
	var result ProxyCfg
	if err := json.Unmarshal(text, &result); err != nil {
		return err
	}

	if result.Port == 0 {
		result.Port = DefaultPortSocks5
	}

	*p = Socks5(result)
	return nil
}

func (p *Socks5) MarshalJSON() ([]byte, error) {
	if p.Port == 0 {
		p.Port = DefaultPortSocks5
	}

	proxyCfg := ProxyCfg(*p)
	return json.Marshal(&proxyCfg)
}
