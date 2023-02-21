package cpanel

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"golang.org/x/net/publicsuffix"
)

type Args map[string]interface{}

type API struct {
	URL      *url.URL
	Username string
	Token    string
	cl       *http.Client
	ctx      *context.Context
}

type DNSRecordData struct {
	Line   int
	Domain string
	Type   string
	Record string
	TTL    int
}

type DNSRecordResponse struct {
	Newserial string
	StatusMsg string
	Status    bool
}

func NewClient(cpURL, username, token string) (*API, error) {
	PTransport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if cpURL == "" || username == "" || token == "" {
		return nil, errors.New("not all required details (URL, username, token) were provided")
	}
	u, err := url.Parse(cpURL)
	if err != nil || u.Hostname() == "" {
		return nil, errors.Errorf("'%s' is not a URL", cpURL)

	}

	return &API{
		URL:      u,
		Username: username,
		Token:    token,
		cl:       &http.Client{Transport: PTransport},
	}, nil
}

func (cpa *API) Request(module, function string, arguments Args) ([]gjson.Result, error) {
	var reqURL = cpa.URL
	var req *http.Request
	var err error

	ctx := context.Background()

	reqArgs := url.Values{}
	for k, v := range arguments {
		reqArgs.Add(k, fmt.Sprintf("%v", v))
	}

	reqArgs.Add("cpanel_jsonapi_user", cpa.Username)
	reqArgs.Add("cpanel_jsonapi_apiversion", "2")
	reqArgs.Add("cpanel_jsonapi_module", module)
	reqArgs.Add("cpanel_jsonapi_func", function)

	reqURL.Path = "/json-api/cpanel"

	reqURL.RawQuery = reqArgs.Encode()
	req, err = http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	auth := fmt.Sprintf("cpanel %s:%s", cpa.Username, cpa.Token)

	req.Header.Add("Authorization", auth)
	req.Header.Set("User-Agent", "DNS Updater")

	req = req.WithContext(ctx)

	if cpa.cl == nil {
		cpa.cl = http.DefaultClient
	}

	resp, err := cpa.cl.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "API request %s:%s", module, function)
	}
	defer resp.Body.Close()

	// Buffer the full response. This costs more memory but we want to to report the contents if
	// it's not valid JSON.
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read full API response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("API request %s/%s failed: HTTP %s",
			module, function, resp.Status)
	}

	/* on error
	   "data": {
	     "reason":"Could not find function 'fetchzone_records1' in module 'ZoneEdit'",
	     "result": 0
	     }
	*/

	// fetch için
	reason := gjson.GetBytes(buf, "cpanelresult.data.reason").Value()
	if reason != nil {
		return nil, errors.New(reason.(string))

	}

	value := gjson.GetBytes(buf, "cpanelresult.data")
	return value.Array(), nil

}

func (cpa *API) FindDNSRecord(domain, recType string) (*DNSRecordData, error) {
	var resp []gjson.Result
	var rec DNSRecordData

	zone, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimRight(domain, "."))

	arguments := Args{}
	arguments["domain"] = zone
	if !(recType == "A" || recType == "TXT") {
		return &rec, errors.New(recType + " type is not implemented in cPanel FindDNSRecord function")
	}
	arguments["type"] = recType

	arguments["name"] = domain

	if resp, err = cpa.Request("ZoneEdit", "fetchzone_records", arguments); err != nil {

		return nil, err
	}

	if len(resp) == 0 {

		return &rec, nil
	}

	rec.Line = int(resp[0].Get("line").Int())
	rec.Domain = resp[0].Get("name").String()
	rec.Type = resp[0].Get("type").String()
	rec.Record = resp[0].Get("record").String()
	rec.TTL = int(resp[0].Get("ttl").Int())

	return &rec, err
}

func (cpa *API) GetBaseDomains() (*[]string, error) {
	var resp []gjson.Result
	var err error
	var domainList []string

	if resp, err = cpa.Request("DomainLookup", "getbasedomains", nil); err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, errors.New("Empty response")
	}

	for _, b := range resp {
		domainList = append(domainList, b.Get("domain").String())
	}
	return &domainList, err
}

func (cpa *API) UpdateTXTRecord(rec DNSRecordData, value string) (z *DNSRecordResponse, err error) {
	var resp []gjson.Result

	zone, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimRight(rec.Domain, "."))

	arguments := Args{}
	arguments["domain"] = zone
	arguments["type"] = rec.Type
	arguments["name"] = rec.Domain
	arguments["txtdata"] = value
	arguments["line"] = rec.Line
	if rec.TTL > 0 {
		arguments["ttl"] = rec.TTL
	}

	if resp, err = cpa.Request("ZoneEdit", "edit_zone_record", arguments); err != nil {
		return z, err
	}
	if len(resp) == 0 {
		return nil, errors.New("Empty response")
	}

	var response DNSRecordResponse

	response.StatusMsg = resp[0].Get("result.statusmsg").String() // : [{"result":{"statusmsg":"","newserial":"2021050400","status":1}}]
	response.Newserial = resp[0].Get("result.newserial").String()
	response.Status = resp[0].Get("result.status").Bool()

	return &response, err
}

func (cpa *API) AddTXTRecords(domain, value string, ttl int) (z *DNSRecordResponse, err error) {
	var resp []gjson.Result

	zone, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimRight(domain, "."))

	arguments := Args{}
	arguments["domain"] = zone
	arguments["type"] = "TXT"
	arguments["name"] = domain
	arguments["txtdata"] = value
	if ttl > 0 {
		arguments["ttl"] = ttl
	}

	if resp, err = cpa.Request("ZoneEdit", "add_zone_record", arguments); err != nil {
		return z, err
	}
	if len(resp) == 0 {
		return nil, errors.New("Empty response")
	}

	var response DNSRecordResponse

	response.StatusMsg = resp[0].Get("result.statusmsg").String() // : [{"result":{"statusmsg":"","newserial":"2021050400","status":1}}]
	response.Newserial = resp[0].Get("result.newserial").String()
	response.Status = resp[0].Get("result.status").Bool()

	return &response, err
}

func (cpa *API) RemoveTXTRecord(rec DNSRecordData) (z *DNSRecordResponse, err error) {
	var resp []gjson.Result

	zone, err := publicsuffix.EffectiveTLDPlusOne(strings.TrimRight(rec.Domain, "."))

	arguments := Args{}
	arguments["domain"] = zone
	arguments["line"] = rec.Line

	if resp, err = cpa.Request("ZoneEdit", "remove_zone_record", arguments); err != nil {
		return z, err
	}
	if len(resp) == 0 {
		return nil, errors.New("Empty response")
	}

	var response DNSRecordResponse

	response.StatusMsg = resp[0].Get("result.statusmsg").String()
	response.Newserial = resp[0].Get("result.newserial").String()
	response.Status = resp[0].Get("result.status").Bool()

	return &response, err
}
