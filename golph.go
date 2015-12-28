package golph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.1.0"
	defaultBaseURL = "https://secure.phabricator.url/"
	userAgent      = "golph/" + libraryVersion
	mediaType      = "application/json"
	postMediaType  = "application/x-www-form-urlencoded; param=value"
)

// Client manages communication to Phabricator
// Since Phabricator is not versioned, this is a bit wonky in the case of breaking changes.
type Client struct {
	// HTTP client used to communicate with the Phabricator API.
	client *http.Client

	// API Token for authentication
	apiToken string

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Conduit connections, see https://secure.phabricator.com/conduit/
	Projects ProjectsService
	//Tasks    TasksService
	//Users    UsersService

	// Optional function called after every successful request made to the Phabricator APIs
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// Response is a DigitalOcean response. This wraps the standard http.Response returned from DigitalOcean.
type Response struct {
	*http.Response

	// Monitoring URI
	Monitor string
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error message
	Message string
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// NewClient returns a new Phabricator API client.
func NewClient(apiToken string, phabricatorUrl string, httpClient *http.Client) *Client {
	if phabricatorUrl == "" {
		phabricatorUrl = defaultBaseURL
	}

	baseURL, _ := url.Parse(phabricatorUrl)

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, apiToken: apiToken, BaseURL: baseURL, UserAgent: userAgent}
	c.Projects = &ProjectsServiceOp{client: c}
	//c.Tasks = &TasksServiceOp{client: c}
	//c.Users = &UsersServiceOp{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := strings.NewReader("")
	if body != nil {
		buf = strings.NewReader(structToValues(body).Encode())
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", postMediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

/***** \/\/\/\/\/

Most of what follows is cargoculted from digitalocean/godo

^^^^^^^^^^ *****/

// OnRequestCompleted sets the Phabricator API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}

	return &response
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	//req.ParseForm()
	//req.Form.Add("api.token", c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err := io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err := json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

func structToValues(i interface{}) (values url.Values) {
	values = url.Values{}
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		// You ca use tags here...
		// tag := typ.Field(i).Tag.Get("tagname")
		// Convert each type into a string for the url.Values string map
		var v string
		switch f.Interface().(type) {
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}
		fieldName := typ.Field(i).Name
		if typ.Field(i).Tag.Get("form") != "" {
			fieldName = typ.Field(i).Tag.Get("form")
		}
		values.Set(fieldName, v)
	}
	return
}
