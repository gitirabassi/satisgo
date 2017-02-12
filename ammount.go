package satisgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
)

//Ammount is used to calculate total "sales" and refunds
type Ammount struct {
	TotalCharge int    `json:"total_charge_amount_unit"`
	TotalRefund int    `json:"total_refund_amount_unit"`
	Currency    string `json:"currency"`
}

//AmmountToday return the total ammount of charges for the past week
func (p *Satis) AmmountToday() (*Ammount, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	amm, err := p.getAmmount(today, now)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

//AmmountYesterday return the total ammount of charges for the past week
func (p *Satis) AmmountYesterday() (*Ammount, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	amm, err := p.getAmmount(today.Add(time.Duration(-24)*time.Hour), today)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

//AmmountSpecificDate return the total ammount of charges for the past week
func (p *Satis) AmmountSpecificDate(year, month, day int) (*Ammount, error) {
	now := time.Now()
	prec := time.Date(year, time.Month(month), day, 0, 0, 0, 0, now.Location())
	last := prec.Add(24 * time.Hour)
	amm, err := p.getAmmount(prec, last)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

//AmmountThisWeek return the total ammount of charges for the past week
func (p *Satis) AmmountThisWeek() (*Ammount, error) {
	now := time.Now()
	lastweek := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	lastweek = lastweek.Add(time.Duration(-167) * time.Hour)
	today := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	amm, err := p.getAmmount(lastweek, today)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

//AmmountThisMonth is cool for accountability
func (p *Satis) AmmountThisMonth() (*Ammount, error) {
	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	amm, err := p.getLongAmmount(month, today)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

//AmmountThisYear is cool for accountability
func (p *Satis) AmmountThisYear() (*Ammount, error) {
	now := time.Now()
	prec := time.Date(now.Year(), time.Month(0), 0, 0, 0, 0, 0, now.Location())
	last := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	amm, err := p.getLongAmmount(prec, last)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

//AmmountSpecificYear is cool for accountability
func (p *Satis) AmmountSpecificYear(year int) (*Ammount, error) {
	prec := time.Date(year, time.Month(0), 0, 0, 0, 0, 0, time.Now().Location())
	last := time.Date(year+1, time.Month(0), 0, 0, 0, 0, 0, time.Now().Location())
	if year == time.Now().Year() {
		now := time.Now()
		last = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	}
	amm, err := p.getLongAmmount(prec, last)
	if err != nil {
		return nil, err
	}
	return amm, nil
}

func (p *Satis) getLongAmmount(start, end time.Time) (*Ammount, error) {
	if d := end.Sub(start); d.Hours() < 168 {
		return p.getAmmount(start, end)
	}
	amm := new(Ammount)
	duration := end.Sub(start)
	hours := int(duration.Hours())
	limit := 167
	var prec, last time.Time
	prec = start
	last = prec.Add(time.Duration(limit) * time.Hour)

	if hours > limit {
		for i := 0; i < hours/limit; i++ {
			// color.Blue(fmt.Sprint(prec))
			// color.Blue(fmt.Sprint(last))
			a, err := p.getAmmount(prec, last)
			if err != nil {
				return nil, err
			}
			amm.Currency = a.Currency
			amm.TotalCharge += a.TotalCharge
			amm.TotalRefund += a.TotalRefund
			prec = last
			last = prec.Add(time.Duration(limit) * time.Hour)
		}
	}
	// color.Blue(fmt.Sprint("last -- ", prec, "-- hours: ", end.Sub(prec).Hours()))
	// color.Blue(fmt.Sprint(end))
	a, err := p.getAmmount(prec, end)
	if err != nil {
		return nil, err
	}
	amm.Currency = a.Currency
	amm.TotalCharge += a.TotalCharge
	amm.TotalRefund += a.TotalRefund
	return amm, nil
}

func (p *Satis) getAmmount(start, end time.Time) (*Ammount, error) {
	color.Red(fmt.Sprint(start))
	color.Red(fmt.Sprint(end))
	if d := end.Sub(start); d.Hours() > 168 {
		return nil, fmt.Errorf("Interval is too long")
	}
	//------------------ NEW version of the query ----------------
	q := url.Values{}
	q.Set("starting_date", putUnix(start))
	q.Set("ending_date", putUnix(end))
	query := q.Encode()
	r, err := http.NewRequest("GET", p.ammountsURL()+"?"+query, nil)
	if err != nil {
		return nil, err
	}
	//-------------------------- END -----------------------------

	//------------------ OLD version of the query ----------------
	// u, err := url.Parse(p.ammountsURL())
	// if err != nil {
	// 	return nil, err
	// }
	// q := u.Query()
	// q.Set("starting_date", putUnix(start))
	// q.Set("ending_date", putUnix(end))
	// u.RawQuery = q.Encode()
	// r := &http.Request{
	// 	Method:     "GET",
	// 	URL:        u,
	// 	Header:     make(http.Header),
	// 	Proto:      "HTTP/1.1",
	// 	ProtoMajor: 1,
	// 	ProtoMinor: 1,
	// 	Host:       u.Host,
	// }
	//-------------------------- END -----------------------------
	status, b, err := p.makeCall(r)
	if err != nil {
		return nil, fmt.Errorf("Error making the call to API: %s", err.Error())
	}
	if status != 200 {
		return nil, fmt.Errorf("Return status is %d:not compatible with the success case", status)
	}
	amm := new(Ammount)
	err = json.Unmarshal(b, amm)
	if err != nil {
		return nil, err
	}
	return amm, nil
}
