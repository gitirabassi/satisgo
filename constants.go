package satisgo

const (
	prod     = "https://authservices.satispay.com"
	sand     = "https://staging.authservices.satispay.com"
	auth     = "/wally-services/protocol/authenticated"
	users    = "/online/v1/users"
	charges  = "/online/v1/charges"
	refunds  = "/online/v1/refunds"
	ammounts = "/online/v1/amounts"
	dev      = "staging"
	eur      = "EUR"
)

const (
	//Required represent a status waiting for a payment
	Required = "REQUIRED"
	//Success represent "payment successfull"
	Success = "SUCCESS"
	//Failure represent payment Failure
	Failure = "FAILURE"
	//ErrDeclined represent a cancellation from the user
	ErrDeclined = "DECLINED_BY_PAYER"
	//ErrFalseRequest represent a cancellation by the wrong user
	ErrFalseRequest = "DECLINED_BY_PAYER_NOT_REQUIRED"
	//ErrNewer is the cancellation of the charge because the user received a new one
	ErrNewer = "CANCEL_BY_NEW_CHARGE"
	//ErrInternal represent a failure from the server
	ErrInternal = "INTERNAL_FAILURE"
	//ErrExpired the user took to long to respond to the charge request
	ErrExpired = "EXPIRED"
)

const (
	//ReasonDuplicate happens when two of the same charge appears
	ReasonDuplicate = "DUPLICATE"
	//ReasonFraud happens when someone else creates and approves a transaction for a user.
	ReasonFraud = "FRAUDULENT"
	//ReasonCustomerRequest if the customer requested a refund (if allowed by company policy)
	ReasonCustomerRequest = "REQUESTED_BY_CUSTOMER"
)

var (
	debug = true
)

func (p *Satis) verificationURL() string {
	if p.env == dev {
		return sand + auth
	}
	return prod + auth
}

func (p *Satis) usersURL() string {
	if p.env == dev {
		return sand + users
	}
	return prod + users
}

func (p *Satis) chargesURL() string {
	if p.env == dev {
		return sand + charges
	}
	return prod + charges
}

func (p *Satis) refundsURL() string {
	if p.env == dev {
		return sand + refunds
	}
	return prod + refunds
}

func (p *Satis) ammountsURL() string {
	if p.env == dev {
		return sand + ammounts
	}
	return prod + ammounts
}
