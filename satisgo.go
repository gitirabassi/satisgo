package satisgo

import (
	"fmt"
	"net/http"
)

//Satis is the base unit for a payment/action with the satispay API
type Satis struct {
	bearer   string
	env      string
	verified bool
}

//New is the generator for a basic interaction with the API
func New(bearer, env string) (*Satis, error) {
	p := new(Satis)
	switch env {
	case "staging":
		p.env = env
		break
	case "production":
		p.env = env
		break
	default:
		return nil, fmt.Errorf("Wrong status string (only 'production' and 'staging' allowed)")
	}
	//find some parameters to check the string-validity of bearer
	//mybe only allow a subset of characters
	p.bearer = bearer
	return p, nil
}

//Verify is used to make sure the token is correct
func (p *Satis) Verify() error {
	r, err := http.NewRequest("GET", p.verificationURL(), nil)
	if err != nil {
		return err
	}
	status, _, err := p.makeCall(r)
	if err != nil {
		return err
	}
	if status != 204 {
		return fmt.Errorf("Return status is %d:not compatible with the verification", status)
	}
	return nil
}
