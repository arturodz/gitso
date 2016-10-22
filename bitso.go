package bitso

import (
  "fmt"
  "strconv"
  "time"
  "crypto/hmac"
  "encoding/hex"
  "crypto/sha256"
  "errors"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
)

const (
  URL = "https://api.bitso.com/v2/"
  btcmxn = "btc_mxn"
  ethmxn = "eth_mxn"
  tickerPath = "ticker"
  transactionsPath = "transactions"
  orderBookPath = "order_book"
)

type Ticker struct {
  High string
  Last string
  Timestamp string
  Volume string
  Vwap string
  Low string
  Ask string
  Bid string
}

type OrderBook struct {
  Asks [][]string
  Bids [][]string
}

type Transaction struct {
  Amount string
  Date string
  Price string
  Tid int
  Side string
}

// Client allows you to access to the Bitso API
type Client struct {
  configuration *Configuration
}

// Configuration stores the information needed to access
// to the private endpoints and the sanbox mode trigger.
type Configuration struct {
  Key       string
  Secret    string
  ClientId  string
  Sandbox   bool
}

// NewClient returns a new Bitso API client. It receives
// an optional Configuration, which is used to
// authenticate and enter in sandbox mode.
func NewClient(configuration *Configuration) *Client {
  return &Client{configuration}
}

// IsSandbox returns true if the sandbox mode is
// turned on.
func (c *Client) IsSandbox() bool {
  if c.configuration == nil {
    return false
  }
  return c.configuration.Sandbox
}

func (c *Client) Ticker(book string) (*Ticker, error) {
  if validateBook(book) == false {
    err := errors.New("Invalid book value")
    return nil, err
  }
  ticker := &Ticker{}
  v := &url.Values{}
  v.Set("book", book)
  err := c.get(tickerPath, v, ticker)
  if err != nil {
    return nil, err
  }
  return ticker, nil
}

func (c *Client) OrderBook(book string, group bool) (*OrderBook, error) {
  if validateBook(book) == false {
    err := errors.New("Invalid book value")
    return nil, err
  }
  orderBook := &OrderBook{}
  v := &url.Values{}
  v.Set("book", book)
  err := c.get(orderBookPath, v, orderBook)
  if err != nil {
    return nil, err
  }
  return orderBook, nil
}

/*
Transactions returns a list of recent trades from the specified book
and the specified time frame.

Valid time frames are hour and minute. Leaving time blank will set hour as the default frame.
*/
func (c *Client) Transactions(book string, time string) ([]*Transaction, error) {
  var transactions []*Transaction
  if validateBook(book) == false {
    err := errors.New("Invalid book value")
    return nil, err
  }
  v := &url.Values{}
  v.Set("book", book)
  v.Set("time", time)
  err := c.get(transactionsPath, v, &transactions)
  if err != nil {
    return nil, err
  }
  return transactions, nil
}

func (c *Client) getSignature() (signature, nonce string) {
  if c.validateConfiguration() == false {
    panic("can't generate a signature without configuration")
  }
  nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
  key := c.configuration.Key
  clientId := c.configuration.ClientId
  secret := c.configuration.Secret
  message := nonce + key + clientId
  fmt.Println(message)
  signature = sign(message, secret)
  return
}

func (c *Client) validateConfiguration() bool {
  if c.configuration == nil {
    return false
  }
  return true
}

func validateBook(book string) bool {
  return book == btcmxn || book == ethmxn
}

func (c *Client) get(path string, query *url.Values, schema interface{}) (error) {
  if config := c.configuration; config != nil && config.Sandbox == true {
    //Mock response
  }
  u, err := url.Parse(URL + path)
  if err != nil {
    return err
  }
  if query != nil {
    u.RawQuery = query.Encode()
  }
  resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }
  err = json.Unmarshal(body, schema)
	if err != nil {
		return err
	}
  return nil
}

func shouldBeCalled(function interface{}) bool {
  fmt.Println(function)

  fmt.Println(function.(func(string, string) string)("m", "k"))
  return false
}

func sign(message, key string) string {
  mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	bytes := mac.Sum(nil)
  s := hex.EncodeToString(bytes)
  fmt.Println(s)

  return s
}
