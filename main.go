package ruomu_uc

import (
	"encoding/json"

	"github.com/hashicorp/go-plugin"
	"github.com/yockii/ruomu-core/shared"

	"github.com/yockii/ruomu-uc/controller"
)

type UC struct{}

func (UC) Initial(params map[string]string) error {
	return nil
}

func (UC) InjectCall(code string, value []byte) ([]byte, error) {
	switch code {
	case CodeUserAdd:
		user, err := controller.UserController.Add(value)
		if err != nil {
			return nil, err
		} else {
			bs, err := json.Marshal(user)
			return bs, err
		}
	}

	return nil, nil
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"uc": &shared.CommunicatePlugin{Impl: &UC{}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
