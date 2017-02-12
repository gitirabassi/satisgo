package satisgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//Charge is the type that handles charges for the user
type Charge struct {
	//ID is the unique charge_id
	ID string `json:"id,omitempty"`
	//Description is the about the charge
	Description string `json:"description,omitempty"`
	//for now only "EUR" is supported
	Currency string `json:"currency,omitempty"`
	//Amount is expressed in EuroCents
	Amount uint64 `json:"amount,omitempty"`
	//Status can have one of 3 states: REQUIRED,SUCCESS,FAILURE
	Status string `json:"status,omitempty"`
	//StatusDetails is helpfull identifing the problem when FAILURE is display as status
	StatusDetails string `json:"status_detail,omitempty"`
	//UserID is the user unique indentifier
	UserID string `json:"user_id,omitempty"`
	//UserShortName is given by webButton
	UserShortName string `json:"user_short_name,omitempty"`
	//Metadata has max 20 fields(key value storage for charges)
	Metadata map[string]string `json:"metadata,omitempty,omitempty"`
	//Return a string of "true" or "false"
	Paid bool `json:"paid,omitempty"`
	//ChargeDate is the date in which the payment has been made
	ChargeDate string `json:"charge_date,omitempty"`
	//Refund is the ammount of the charge that has gone in a Refund
	Refund uint64 `json:"refund_amount,omitempty"`
	//EmailOnSuccess takes "true" or "false" and send an email after a payment has occurred (true by default with this library)
	EmailOnSuccess bool `json:"required_success_email,omitempty"`
	//ExpireIn represent the number of seconds the user has to approve the payment before it expires
	ExpireIn int `json:"expire_in,omitempty"`
	//Expire date is the rielaboration of the ExpireIN from the server
	ExpireDate string `json:"expire_date,omitempty"`
	//Is given to get notified when a "status" change in a charge
	CallbackURL string `json:"callback_url,omitempty"`
}

//NewCharge generates a charge based on the obtained user
func (u *User) NewCharge() (*Charge, error) {
	if u.ID == "" {
		return nil, fmt.Errorf("not possible to create a charge if user_id is empty")
	}
	c := new(Charge)
	c.Currency = eur
	c.UserID = u.ID
	c.EmailOnSuccess = true
	return c, nil
}

//GetCharge returns a charge provided a charge_id
func (p *Satis) GetCharge(id string) (*Charge, error) {
	r, err := http.NewRequest("GET", p.chargesURL()+"/"+id, nil)
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
	c := new(Charge)
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	return c, nil
}

//CancelCharge cancel a charge not yet approved by client
func (c *Charge) CancelCharge(p *Satis) error {
	input := strings.NewReader(`{"charge_state":"CANCELED"}`)
	r, err := http.NewRequest("PUT", p.chargesURL()+"/"+c.ID, input)
	if err != nil {
		return err
	}
	status, b, err := p.makeCall(r)
	if err != nil {
		return fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	ch := new(Charge)
	err = json.Unmarshal(b, ch)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	*c = *ch
	return nil
}

//UpdateChargeDescription returns a modified Charge provided one
func (c *Charge) UpdateChargeDescription(p *Satis) error {
	type body struct {
		Description string `json:"description"`
	}
	var bod body
	bod.Description = c.Description
	data, err := json.Marshal(&bod)
	if err != nil {
		return err
	}
	input := bytes.NewReader(data)
	r, err := http.NewRequest("PUT", p.chargesURL()+"/"+c.ID, input)
	if err != nil {
		return err
	}
	status, b, err := p.makeCall(r)
	if err != nil {
		return fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	ch := new(Charge)
	err = json.Unmarshal(b, ch)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	*c = *ch
	return nil
}

//UpdateChargeMetadata returns a modified Charge provided one
func (c *Charge) UpdateChargeMetadata(p *Satis) error {
	if c.Metadata == nil {
		return fmt.Errorf("metadata not initialized yet, nothing to update")
	}
	type body struct {
		Metadata map[string]string `json:"metadata"`
	}
	var bod body
	bod.Metadata = c.Metadata
	data, err := json.Marshal(&bod)
	if err != nil {
		return err
	}
	input := bytes.NewReader(data)
	req, err := http.NewRequest("PUT", p.chargesURL()+"/"+c.ID, input)
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
	ch := new(Charge)
	err = json.Unmarshal(b, ch)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	*c = *ch
	return nil
}

//SetDescription helps inserting a description into the charge
//THIS IS NOT MANDATORY BUT HIGHLY SUGGESTED
func (c *Charge) SetDescription(s string) error {
	//check if string is compatible with description lenght or characters types
	c.Description = s
	return nil
}

//SetAmmount helps inserting a description into the charge
//THIS IS MANDATORY
func (c *Charge) SetAmmount(a float64) error {
	//check if a is compatible with description lenght or characters types
	if a <= 0 {
		return fmt.Errorf("ammount to charge is negative or equal to zero")
	}
	if a >= 10000 {
		return fmt.Errorf("ammount to charge is too big: not supported by Satispay")
	}
	n := putMoney(a)
	c.Amount = n
	return nil
}

//SetExpiration is used to change the default 15 minutes expiration time for approving a charge
//THIS IS NOT MANDATORY
func (c *Charge) SetExpiration(d time.Duration) error {
	if d < time.Minute {
		return fmt.Errorf("ExpireIn Time too short")
	}
	if d > time.Hour {
		return fmt.Errorf("ExpireIn Time too long")
	}
	c.ExpireIn = int(d.Seconds())
	return nil
}

//SetCallbackURL is used to change the default 15 minutes expiration time for approving a charge
//THIS IS NOT MANDATORY
func (c *Charge) SetCallbackURL(s string) error {
	c.CallbackURL = s
	return nil
}

//SetMetadata is used to add metadata to a charge without making any mess around
//THere are some limits: max 20pairs, key length max 45 chars, value max 500 chars
//THIS IS NOT MANDATORY
func (c *Charge) SetMetadata(key, value string) error {
	if len(key) > 45 || len(key) < 2 {
		return fmt.Errorf("The key has the wrong format (too long/short)")
	}
	if len(value) > 500 {
		return fmt.Errorf("The value has the wrong format (too long/short)")
	}
	if c.Metadata == nil {
		c.Metadata = make(map[string]string)
	}
	_, exist := c.Metadata[key]
	if exist && value == "" {
		delete(c.Metadata, key)
		return nil
	}
	if len(c.Metadata) == 20 && !exist {
		return fmt.Errorf("Metadata is too long already")
	}
	c.Metadata[key] = value
	return nil
}

//CreateCharge is the function that makes the call to Satispay API
func (c *Charge) CreateCharge(p *Satis) error {
	if c.UserID == "" {
		return fmt.Errorf("User_ID cannot be empty")
	}
	if c.Amount == 0 {
		return fmt.Errorf("Amount cannot be empty")
	}
	c.Currency = eur
	c.EmailOnSuccess = true
	if c.ID != "" {
		return fmt.Errorf("Charge ID already exist: charge already created")
	}
	if len(c.Metadata) > 20 {
		return fmt.Errorf("Metadata is too long")
	}
	if c.ExpireDate != "" {
		return fmt.Errorf("Charge expire_date already exist: charge already created")
	}
	if c.CallbackURL == "" {
		return fmt.Errorf("CallbackURL cannot be empty")
	}
	//some more checking if Charge object is good
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Error formatting Charges for creation of charge: %s", err.Error())
	}
	// fmt.Println(string(data))
	input := bytes.NewReader(data)
	r, err := http.NewRequest("POST", p.chargesURL(), input)
	if err != nil {
		return err
	}
	status, b, err := p.makeCall(r)
	if err != nil {
		return fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	charg := new(Charge)
	err = json.Unmarshal(b, charg)
	if err != nil {
		return fmt.Errorf("Error unmarshaling response to Charge: %s", err.Error())
	}
	*c = *charg
	return nil
}
