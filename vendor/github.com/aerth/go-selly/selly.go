// Copyright 2018 The go-selly authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Selly API wrapper for selly.gg
package selly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aerth/tgun"
)

const DefaultUserAgent = "go-selly - https://github.com/aerth/go-selly"

// Selly settings
type Selly struct {
	httpClient tgun.Client
	Email      string
	Token      string
	UserAgent  string // "Yourusername - website-using-api.com"
}

type ErrorResponse struct {
	Message string `json:"message"`
	Errors  struct {
		Title []string `json:"title"`
	} `json:"errors"`
}

func (e ErrorResponse) Error() string {
	return e.String()
}

func (e ErrorResponse) String() string {
	if len(e.Errors.Title) == 0 {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Errors.Title)
}

// Product ...
type Product struct {
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Stock           int         `json:"stock"`
	Price           string      `json:"price"`
	Currency        string      `json:"currency"`
	ProductType     int         `json:"product_type"`
	Bitcoin         bool        `json:"bitcoin"`
	Paypal          bool        `json:"paypal"`
	Stripe          bool        `json:"stripe"`
	Litecoin        bool        `json:"litecoin"`
	Dash            bool        `json:"dash"`
	Ethereum        bool        `json:"ethereum"`
	PerfectMoney    bool        `json:"perfect_money"`
	BitcoinCash     bool        `json:"bitcoin_cash"`
	Ripple          bool        `json:"ripple"`
	Private         bool        `json:"private"`
	Unlisted        bool        `json:"unlisted"`
	SellerNote      string      `json:"seller_note"`
	MaximumQuantity interface{} `json:"maximum_quantity"`
	MinimumQuantity int         `json:"minimum_quantity"`
	Custom          struct{}    `json:"custom"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// ProductGroup ...
type ProductGroup struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	ProductIds []string  `json:"product_ids"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Order ...
type Order struct {
	ID            string            `json:"id"`
	ProductID     string            `json:"product_id"`
	Email         string            `json:"email"`
	IPAddress     string            `json:"ip_address"`
	CountryCode   string            `json:"country_code"`
	UserAgent     string            `json:"user_agent"`
	Value         string            `json:"value"`
	Currency      string            `json:"currency"`
	Gateway       string            `json:"gateway"`
	RiskLevel     int               `json:"risk_level"`
	Status        int               `json:"status"`
	Delivered     string            `json:"delivered"`
	CryptoValue   interface{}       `json:"crypto_value"`
	CryptoAddress interface{}       `json:"crypto_address"`
	Referral      interface{}       `json:"referral"`
	UsdValue      string            `json:"usd_value"`
	ExchangeRate  string            `json:"exchange_rate"`
	Custom        map[string]string `json:"custom"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// Coupon ...
type Coupon struct {
	ID         int         `json:"id"`
	Code       string      `json:"code"`
	Discount   int         `json:"discount"`
	MaxUses    interface{} `json:"max_uses"`
	ProductIds []string    `json:"product_ids"`
	Uses       int         `json:"uses"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// Query ...
type Query struct {
	ID          string    `json:"id"`
	Secret      string    `json:"secret"`
	Title       string    `json:"title"`
	Email       string    `json:"email"`
	Message     string    `json:"message"`
	Status      int       `json:"status"`
	CountryCode string    `json:"country_code"`
	IPAddress   string    `json:"ip_address"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Request struct {
	Title      string `json:"title"`
	Gateway    string `json:"gateway"`
	Email      string `json:"email"`
	Value      string `json:"value"`
	Currency   string `json:"currency"`
	ReturnURL  string `json:"return_url"`
	WebhookURL string `json:"webhook_url"`
}

type Webhook struct {
	ID            string      `json:"id"`
	ProductID     string      `json:"product_id"`
	Email         string      `json:"email"`
	IPAddress     string      `json:"ip_address"`
	CountryCode   string      `json:"country_code"`
	UserAgent     string      `json:"user_agent"`
	Value         string      `json:"value"`
	Currency      string      `json:"currency"`
	Gateway       string      `json:"gateway"`
	RiskLevel     int         `json:"risk_level"`
	Status        int         `json:"status"`
	Delivered     string      `json:"delivered"`
	CryptoValue   interface{} `json:"crypto_value"`
	CryptoAddress interface{} `json:"crypto_address"`
	Referral      string      `json:"referral"`
	WebhookType   int         `json:"webhook_type"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

// New creates a new Selly instance
func New(email, token, useragent string) *Selly {
	if useragent == "" {
		useragent = DefaultUserAgent
	}
	httpclient := tgun.Client{
		UserAgent:    useragent,
		AuthUser:     email,
		AuthPassword: token,
	}
	return &Selly{
		httpClient: httpclient,
	}
}

// NewProxy creates a new Selly instance using proxy for requests
// Proxy format: socks5://127.0.0.1:1080
func NewProxy(email, token, useragent, proxy string) *Selly {
	if useragent == "" {
		useragent = DefaultUserAgent
	}
	httpclient := tgun.Client{
		UserAgent:    useragent,
		AuthUser:     email,
		AuthPassword: token,
		Proxy:        proxy,
	}
	return &Selly{
		httpClient: httpclient,
	}
}

// GetProduct returns product
func (s *Selly) GetProduct(id string) (*Product, error) {
	furl := "https://selly.gg/api/v2/products/%s"
	b, err := s.httpClient.GetBytes(fmt.Sprintf(furl, id))
	if err != nil {
		return nil, err
	}
	p := Product{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(b, &errar)
		return nil, errar
	}
	return &p, err
}

// GetProducts returns all products
func (s *Selly) GetProducts() ([]Product, error) {
	furl := "https://selly.gg/api/v2/products"
	b, err := s.httpClient.GetBytes(furl)
	if err != nil {
		return nil, err
	}
	p := []Product{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(b, &errar)
		return nil, errar
	}
	return p, err
}

// GetOrder returns specific orders
func (s *Selly) GetOrder(id string) (*Order, error) {
	furl := "https://selly.gg/api/v2/orders/%s"
	b, err := s.httpClient.GetBytes(fmt.Sprintf(furl, id))
	if err != nil {
		return nil, err
	}
	o := Order{}
	err = json.Unmarshal(b, &o)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(b, &errar)
		return nil, errar
	}
	return &o, err
}

// GetOrders returns all orders
func (s *Selly) GetOrders() ([]Order, error) {
	furl := "https://selly.gg/api/v2/orders"
	b, err := s.httpClient.GetBytes(furl)
	if err != nil {
		return nil, err
	}
	o := []Order{}
	err = json.Unmarshal(b, &o)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(b, &errar)
		return nil, errar
	}
	return o, err
}

// GetCoupon returns specific coupon
func (s *Selly) GetCoupon(id string) (*Coupon, error) {
	furl := "https://selly.gg/api/v2/coupons/%s"
	b, err := s.httpClient.GetBytes(fmt.Sprintf(furl, id))
	if err != nil {
		return nil, err
	}
	c := Coupon{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(b, &errar)
		return nil, errar
	}
	return &c, err
}

// GetCoupons returns all coupons
func (s *Selly) GetCoupons() ([]Coupon, error) {
	furl := "https://selly.gg/api/v2/coupons"
	b, err := s.httpClient.GetBytes(furl)
	if err != nil {
		return nil, err
	}
	c := []Coupon{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(b, &errar)
		return nil, errar
	}
	return c, err
}

// NewCoupon creates a new coupon and returns a coupon or error
func (s *Selly) NewCoupon(code string, discount int, productIDs []string) (*Coupon, error) {
	url := "https://selly.gg/api/v2/coupons"
	coupon := Coupon{
		Discount:   discount,
		Code:       code,
		ProductIds: productIDs,
	}
	data, err := json.Marshal(coupon)
	if err != nil {
		return nil, err
	}
	resp, err := s.postreq(url, data)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &coupon)
	if err != nil {
		errar := ErrorResponse{}
		err = json.Unmarshal(body, &errar)
		return nil, errar
	}
	return &coupon, err
}

func (s *Selly) postreq(url string, data interface{}) (*http.Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return s.httpClient.Do(req)
}
