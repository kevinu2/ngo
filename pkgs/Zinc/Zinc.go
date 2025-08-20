package Zinc

import (
	"context"

	client "github.com/zinclabs/sdk-go-zincsearch"
)

var Z *Zinc

type Zinc struct {
	Conf *client.Configuration
	Api  *client.APIClient
	CTX  *context.Context
}

func init() {
	New()
}

func New() *Zinc {
	return Z.New()
}

func (z *Zinc) New() *Zinc {
	return new(Zinc)
}

func (z *Zinc) Auth(username, password string, servers []string) {
	z.Conf = client.NewConfiguration()
	z.Conf.Servers = make(client.ServerConfigurations, 0)
	for _, v := range servers {
		z.Conf.Servers = append(z.Conf.Servers, client.ServerConfiguration{
			URL:         v,
			Description: "",
			Variables:   nil,
		})
	}

	z.Api = client.NewAPIClient(z.Conf)

	auth := context.WithValue(context.Background(), client.ContextBasicAuth, client.BasicAuth{
		UserName: username,
		Password: password,
	})
	z.CTX = &auth
}
