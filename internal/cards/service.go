package books

import (
	"context"
	"fmt"

	"github.com/mercadolibre/fury_go-core/pkg/rusty"
	"github.com/mercadolibre/fury_go-core/pkg/transport/httpclient"
)

var (
	apiEndpoint = "https://rickandmortyapi.com/api"
)

type Service interface {
	Save()
	Find()
}

type service struct {
}

func (s *service) Save() {
	endpoint, err := rusty.NewEndpoint(httpclient.New(), apiEndpoint)
	if err != nil {
		fmt.Println(err)
	}
	endpoint.Get(context.Background())
}

func (s *service) Find() {

}
