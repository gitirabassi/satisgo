package satisgo

import (
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/fatih/color"
)

func (p *Satis) makeCall(req *http.Request) (int, []byte, error) {
	var insecure *tls.Config
	if p.env == dev {
		insecure = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
	} else {
		insecure = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	insecure.BuildNameToCertificate()
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       insecure,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.bearer))
	if req.Method == http.MethodPost {
		req.Header.Set("Idempotency-Key", generateUUID())
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := checkIntegrity(resp)
	if err != nil {
		return -1, nil, err
	}

	// this block needs to be taken out ------------------------
	err = handleHeader(resp.StatusCode)
	if err != nil {
		//Parse body to find error code and message
		code, e := jsonparser.GetInt(body, "code")
		if e != nil {
			return -1, nil, err
		}
		msg, e := jsonparser.GetString(body, "message")
		if e != nil {
			return -1, nil, err
		}
		return -1, nil, fmt.Errorf("%s --> CODE %d: %s", err.Error(), code, msg)
	}
	// till here -------------------------------------------------
	return resp.StatusCode, body, nil
}

func checkIntegrity(r *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	color.Yellow(string(body))
	color.Blue(fmt.Sprint(r.Header))
	//checkin content lenght
	lenght, ok := r.Header["Content-Length"]
	if ok != true {
		return nil, fmt.Errorf("Content-Lenght value in header does not exist")
	}
	if len(lenght) != 1 {
		return nil, fmt.Errorf("Multiple Content-Lenght values in header")
	}
	l, err := strconv.Atoi(lenght[0])
	if err != nil {
		return nil, err
	}
	if len(body) != l {
		return nil, fmt.Errorf("Content-Lenght value in header is not true")
	}
	if r.StatusCode == 204 {
		return []byte(""), nil
	}
	//checking content type
	t, ok := r.Header["Content-Type"]
	if ok != true {
		return nil, fmt.Errorf("Content-Type value in header does not exist")
	}
	if len(t) != 1 {
		return nil, fmt.Errorf("Multiple Content-Type values in header")
	}
	if t[0] != "application/json" && t[0] != "application/json;charset=utf-8" {
		return nil, fmt.Errorf("Content-Type value in header is not correct")
	}
	//check digest
	digest, ok := r.Header["Digest"]
	if ok != true {
		return nil, fmt.Errorf("Content-Type value in header does not exist")
	}
	if len(digest) != 1 {
		return nil, fmt.Errorf("Multiple Content-Type values in header")
	}
	hash := sha512.New()
	_, err = hash.Write(body)
	if err != nil {
		return nil, err
	}
	dig := "SHA-512=" + base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if digest[0] != dig {
		return nil, fmt.Errorf("Digest is incorrect")
	}
	//check wlt
	wlt, ok := r.Header["X-Satispay-Cid"]
	if ok != true {
		return nil, fmt.Errorf("X-Satispay-Cid value in header does not exist")
	}
	if len(wlt) != 1 {
		return nil, fmt.Errorf("Multiple X-Satispay-Cid values in header")
	}
	w, err := jsonparser.GetString(body, "wlt")
	if err != nil {
		if err.Error() != "Key path not found" {
			return nil, fmt.Errorf("Error Parsing WLT from body: %s", err.Error())
		}
	}
	if err == nil {
		if w != wlt[0] {
			return nil, fmt.Errorf("WLT checks has gone wrong")
		}
	}
	//Still NOT checking
	//   - SERVER header
	//   - DATE HEADER
	return body, nil
}
