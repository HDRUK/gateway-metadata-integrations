package pull

import (
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

type Pull struct {
	Uri         string
	Username    string
	Password    string
	AccessToken string
	Verbose     bool
}

// NewPull Creates a new instance of Pull
func NewPull(uri string, username string, password string, accessToken string, verbose bool) *Pull {
	pull := Pull{
		Uri:         uri,
		Username:    username,
		Password:    password,
		AccessToken: accessToken,
		Verbose:     verbose,
	}

	return &pull
}

func (p *Pull) Auth() {

}

func Run() {
	fmt.Printf("Pull Running\r\n")
}
