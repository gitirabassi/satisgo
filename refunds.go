package satisgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//Refund is the type that handles charges for the user
type Refund struct {
	//ID is the unique charge_id
	ID string `json:"id,omitempty"`
	//ChargeID is the charge unique indentifier
	ChargeID string `json:"charge_id,omitempty"`
	//Description is the about the charge
	Description string `json:"description,omitempty"`
	//for now only "EUR" is supported
	Currency string `json:"currency,omitempty"`
	//Amount is expressed in EuroCents
	Amount uint64 `json:"amount,omitempty"`
	//Metadata has max 20 fields(key value storage for charges)
	Metadata map[string]string `json:"metadata,omitempty"`
	//Created is the time in UnixMilli of the creation of the refund
	Created string `json:"reason,omitempty"`
	//Reason is the reason a refund occurred
	Reason string `json:"reason,omitempty"`
}

//NewRefund generates a refund based on the given charge
func (c *Charge) NewRefund() (*Refund, error) {
	if c.ID == "" {
		return nil, fmt.Errorf("before creating a refund, you must create the charge thru the appropriate method")
	}
	r := new(Refund)
	r.ChargeID = c.ID
	r.Currency = eur
	return r, nil
}

//NewRefundWithAmount generates a refund based on the given charge and the ammount supplied
func (c *Charge) NewRefundWithAmount(a float64) (*Refund, error) {
	if c.ID == "" {
		return nil, fmt.Errorf("before creating a refund, you must create the charge thru the appropriate method")
	}
	r := new(Refund)
	r.ChargeID = c.ID
	r.Currency = eur
	err := r.SetAmmount(a)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//SetDescription helps inserting a description into the refund
//THIS IS NOT MANDATORY BUT HIGHLY SUGGESTED
func (r *Refund) SetDescription(s string) error {
	//check if string is compatible with description lenght or characters types
	r.Description = s
	return nil
}

//SetReason helps inserting a reason into the refund instance
//THIS IS NOT MANDATORY BUT HIGHLY SUGGESTED
//Only few reasons are allowed (use of package costants is suggested)
func (r *Refund) SetReason(reason string) error {
	//check if string is compatible with description lenght or characters types
	switch reason {
	case ReasonDuplicate:
		r.Reason = ReasonDuplicate
		break
	case ReasonFraud:
		r.Reason = ReasonFraud
		break
	case ReasonCustomerRequest:
		r.Reason = ReasonCustomerRequest
		break
	default:
		return fmt.Errorf("Reason provided is not supported, please read the documentation for allowes fields")
	}
	return nil
}

//SetAmmount helps inserting a description into the refund
//THIS IS MANDATORY
func (r *Refund) SetAmmount(a float64) error {
	//check if a is compatible with description lenght or characters types
	if a <= 0 {
		return fmt.Errorf("ammount to charge is negative or equal to zero")
	}
	if a >= 10000 {
		return fmt.Errorf("ammount to charge is too big: not supported by Satispay")
	}
	n := putMoney(a)
	r.Amount = n
	return nil
}

//SetMetadata is used to add metadata to a refund without making any mess around
//THere are some limits: max 20pairs, key length max 45 chars, value max 500 chars
//THIS IS NOT MANDATORY
func (r *Refund) SetMetadata(key, value string) error {
	if r.Metadata == nil {
		r.Metadata = make(map[string]string)
	}
	if len(key) > 45 || len(key) < 2 {
		return fmt.Errorf("The key has the wrong format (too long/short)")
	}
	if len(value) > 500 || len(value) < 2 {
		return fmt.Errorf("The value has the wrong format (too long/short)")
	}
	_, exist := r.Metadata[key]
	if len(r.Metadata) == 20 && !exist {
		return fmt.Errorf("Metadata is too long already")
	}
	r.Metadata[key] = value
	return nil
}

//GetRefund returns a refund provided a refund_id
func (p *Satis) GetRefund(id string) (*Refund, error) {
	r, err := http.NewRequest("GET", p.refundsURL()+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	status, b, err := p.makeCall(r)
	if err != nil {
		return nil, fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return nil, fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	c := new(Refund)
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	return c, nil
}

//UpdateRefundMetadata returns a modified Refund provided one
func (r *Refund) UpdateRefundMetadata(p *Satis) error {
	if r.Metadata == nil {
		return fmt.Errorf("metadata not initialized yet, nothing to update")
	}
	type body struct {
		Metadata map[string]string `json:"metadata"`
	}
	var bod body
	bod.Metadata = r.Metadata
	data, err := json.Marshal(&bod)
	if err != nil {
		return err
	}
	input := bytes.NewReader(data)
	req, err := http.NewRequest("PUT", p.refundsURL()+"/"+r.ID, input)
	if err != nil {
		return err
	}
	status, b, err := p.makeCall(req)
	if err != nil {
		return fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	ch := new(Refund)
	err = json.Unmarshal(b, ch)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	*r = *ch
	return nil
}

//CreateRefund is the function that makes the call to Satispay API to request a refund with the given parameters
func (r *Refund) CreateRefund(p *Satis) error {
	if r.ChargeID == "" {
		return fmt.Errorf("Charge ID cannot be empty")
	}
	if r.Amount == 0 {
		return fmt.Errorf("Amount cannot be empty")
	}
	r.Currency = eur
	if r.ID != "" {
		return fmt.Errorf("Charge ID already exist: charge already created")
	}
	if len(r.Metadata) > 20 {
		return fmt.Errorf("Metadata is too long")
	}
	//some more checking if Refund object is good
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("Error formatting Charges for creation of charge: %s", err.Error())
	}
	// fmt.Println(string(data))
	input := bytes.NewReader(data)
	req, err := http.NewRequest("POST", p.refundsURL(), input)
	if err != nil {
		return err
	}
	status, b, err := p.makeCall(req)
	if err != nil {
		return fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	ref := new(Refund)
	err = json.Unmarshal(b, ref)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	*r = *ref
	return nil
}
