package orders

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

type KrakenOrders struct {
	endPointPath string
	domainName   string
}

func NewDefaultKrakenOrders() *KrakenOrders {
	return &KrakenOrders{endPointPath: "/api/v3/sendorder", domainName: "https://demo-futures.kraken.com/derivatives"}
}

func generateAuthent(PostData string, endPointPath string, apiKey string) ([]byte, error) {
	sha := sha256.New()
	src := PostData + endPointPath
	sha.Write([]byte(src))

	apiDecode, err := base64.StdEncoding.DecodeString(apiKey)
	if err != nil {
		return nil, fmt.Errorf("error in decode apiKey: %w", err)
	}

	h := hmac.New(sha512.New, apiDecode)
	h.Write(sha.Sum(nil))

	result := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return []byte(result), nil
}

func (ko KrakenOrders) SendOrder(symbol string, size uint64, side string, publicApiKey string, privateApiKey string) (domain.ResponseOrder, error) {
	v := url.Values{}
	v.Add("orderType", "mkt")
	v.Add("symbol", symbol)
	v.Add("side", side)
	v.Add("size", strconv.FormatUint(size, 10))
	queryString := v.Encode()

	req, err := http.NewRequest(http.MethodPost, ko.domainName+ko.endPointPath+"?"+queryString, nil)
	if err != nil {
		return nil, fmt.Errorf("error in creating a request: %w", err)
	}

	auth, err := generateAuthent(queryString, ko.endPointPath, privateApiKey)
	if err != nil {
		return nil, fmt.Errorf("error in generate authent: %w", err)
	}

	req.Header.Add("APIKEY", publicApiKey)
	req.Header.Add("Authent", string(auth))

	//req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux i686; rv:2.0.1) Gecko/20100101 Firefox/4.0.1")

	b, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return nil, fmt.Errorf("error in dump request: %w", err)
	}
	log.Debug(string(b))

	c := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in do request: %w", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error in read response body %w", err)
	}
	log.Debug(string(data))
	jsonResp := domain.KrakenResponseOrder{}
	err = json.Unmarshal(data, &jsonResp)
	if err != nil {
		return nil, fmt.Errorf("error in dump response: %w", err)
	}
	return jsonResp, nil
}
