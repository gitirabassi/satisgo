package satisgo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
)

//User is the parsed object from a call to user list
type User struct {
	ID    string `json:"id"`
	Phone string `json:"phone_number"`
}

//UserFromPhone is the way to get an identifier with a phone number
func (p *Satis) UserFromPhone(phone string) (*User, error) {
	reader := strings.NewReader(fmt.Sprintf(`{"phone_number":"%s"}`, phone))
	r, err := http.NewRequest("POST", p.usersURL(), reader)
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
	id, err := jsonparser.GetString(b, "id")
	if err != nil {
		return nil, fmt.Errorf("Error Parsing ID from body: %s", err.Error())
	}
	ud, err := jsonparser.GetString(b, "uuid")
	if err != nil {
		return nil, fmt.Errorf("Error Parsing UUID from body: %s", err.Error())
	}
	ph, err := jsonparser.GetString(b, "phone_number")
	if err != nil {
		return nil, fmt.Errorf("Error Parsing Phone_number from body: %s", err.Error())
	}
	if phone != ph {
		return nil, fmt.Errorf("the sent number phone is not the same when it came back")
	}
	if ud != id {
		return nil, fmt.Errorf("uuid and id in body of response are not equal")
	}
	u := new(User)
	u.ID = id
	u.Phone = phone
	return u, nil
}

//UserFromID is the way to get a phone number with an id
func (p *Satis) UserFromID(id string) (*User, error) {
	r, err := http.NewRequest("GET", p.usersURL()+"/"+id, nil)
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
	ud, err := jsonparser.GetString(b, "id")
	if err != nil {
		return nil, fmt.Errorf("Error Parsing ID from body: %s", err.Error())
	}
	phone, err := jsonparser.GetString(b, "phone_number")
	if err != nil {
		return nil, fmt.Errorf("Error Parsing Phone_number from body: %s", err.Error())
	}
	if ud != id {
		return nil, fmt.Errorf("The provided id and the one that came back are not equal")
	}
	u := new(User)
	u.ID = id
	u.Phone = phone
	return u, nil
}
