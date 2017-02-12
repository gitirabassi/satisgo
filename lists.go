package satisgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/buger/jsonparser"
)

//GetRefundFromChargeID returns all charges from the beginning
func (p *Satis) GetRefundFromChargeID(chargeID string) (*[]Refund, error) {
	total := make([]Refund, 0, 100)
	temp := make([]Refund, 0, 100)
	last := ""
	stopper := true
	for stopper {
		q := url.Values{}
		q.Set("limit", "100")
		q.Set("charge_id", chargeID)
		if last != "" {
			q.Set("starting_after", last)
		}
		query := q.Encode()
		more, err := p.getList(&temp, p.refundsURL(), query)
		if err != nil {
			return nil, err
		}
		total = append(total, temp...)
		stopper = more
		if len(temp) > 0 {
			last = temp[len(temp)-1].ID
		}
	}
	return &total, nil
}

//GetRefundSinceChargeID returns all charges from the beginning
func (p *Satis) GetRefundSinceChargeID(chargeID string) (*[]Refund, error) {
	total := make([]Refund, 0, 100)
	temp := make([]Refund, 0, 100)
	last := "chargeID"
	stopper := true
	for stopper {
		q := url.Values{}
		q.Set("limit", "100")
		if last != "" {
			q.Set("starting_after", last)
		}
		query := q.Encode()
		more, err := p.getList(&temp, p.refundsURL(), query)
		if err != nil {
			return nil, err
		}
		total = append(total, temp...)
		stopper = more
		if len(temp) > 0 {
			last = temp[len(temp)-1].ID
		}
	}
	return &total, nil
}

//GetAllRefunds returns all charges from the beginning
func (p *Satis) GetAllRefunds() (*[]Refund, error) {
	total := make([]Refund, 0, 100)
	temp := make([]Refund, 0, 100)
	last := ""
	stopper := true
	for stopper {
		q := url.Values{}
		q.Set("limit", "100")
		if last != "" {
			q.Set("starting_after", last)
		}
		query := q.Encode()
		more, err := p.getList(&temp, p.refundsURL(), query)
		if err != nil {
			return nil, err
		}
		total = append(total, temp...)
		stopper = more
		if len(temp) > 0 {
			last = temp[len(temp)-1].ID
		}
	}
	return &total, nil
}

//GetAllUsers returns all charges from the beginning
func (p *Satis) GetAllUsers() (*[]User, error) {
	total := make([]User, 0, 100)
	temp := make([]User, 0, 100)
	last := ""
	stopper := true
	for stopper {
		q := url.Values{}
		q.Set("limit", "100")
		if last != "" {
			q.Set("starting_after", last)
		}
		query := q.Encode()
		more, err := p.getList(&temp, p.usersURL(), query)
		if err != nil {
			return nil, err
		}
		total = append(total, temp...)
		stopper = more
		last = temp[len(temp)-1].ID
	}
	return &total, nil
}

//GetAllCharges returns all charges from the beginning
func (p *Satis) GetAllCharges() (*[]Charge, error) {
	total := make([]Charge, 0, 100)
	temp := make([]Charge, 0, 100)
	last := new(Charge)
	stopper := true
	for stopper {
		q := url.Values{}
		q.Set("limit", "100")
		if last.ID != "" {
			q.Set("starting_after", last.ID)
		}
		query := q.Encode()
		more, err := p.getList(&temp, p.chargesURL(), query)
		if err != nil {
			return nil, err
		}
		total = append(total, temp...)
		stopper = more
		last.ID = temp[len(temp)-1].ID
	}
	return &total, nil
}

//getList is used to manage general lists in the satispay API. the bool in the return indicates if there are more where this came from
func (p *Satis) getList(list interface{}, baseURL, query string) (bool, error) {
	//maybe some checking into the baseURL and query string can be done but since this is an internal function will leave it be wild nad young
	var uri string
	if baseURL == "" {
		return false, fmt.Errorf("no baseURL provided")
	}
	if query == "" {
		uri = baseURL
	} else {
		uri = baseURL + "?" + query
	}
	r, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return false, err
	}
	status, b, err := p.makeCall(r)
	if err != nil {
		return false, fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return false, fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	data, _, _, err := jsonparser.Get(b, "list")
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(data, list)
	if err != nil {
		return false, err
	}
	hasMore, err := jsonparser.GetBoolean(b, "has_more")
	if err != nil {
		return false, err
	}
	return hasMore, nil
}
