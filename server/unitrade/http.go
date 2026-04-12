package unitrade

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/golang-jwt/jwt/v5"
)

type apiClient struct {
	httpClient *http.Client
	baseURL    string
	token      string
	userID     uint
}

type apiEnvelope struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

func newAPIClient(host, token string, userID uint) (*apiClient, error) {
	baseURL, err := normalizeBaseURL(host)
	if err != nil {
		return nil, err
	}
	return &apiClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		token:      token,
		userID:     userID,
	}, nil
}

func normalizeBaseURL(host string) (string, error) {
	raw := strings.TrimSpace(host)
	if raw == "" {
		raw = defaultHost
	}
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("invalid host: %s", host)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")
	if parsed.Path == "" {
		parsed.Path = "/api"
	} else if !strings.HasSuffix(parsed.Path, "/api") {
		parsed.Path = strings.TrimRight(parsed.Path, "/") + "/api"
	}
	return strings.TrimRight(parsed.String(), "/"), nil
}

func (c *apiClient) doJSON(method, endpoint string, query url.Values, body any, out any) (http.Header, error) {
	var bodyBytes []byte
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBytes = data
	}

	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	reqURL.Path = path.Join(reqURL.Path, endpoint)
	if len(query) > 0 {
		reqURL.RawQuery = query.Encode()
	}

	header, retryableNotFound, err := c.doJSONOnce(method, reqURL.String(), bodyBytes, out)
	if err == nil {
		return header, nil
	}
	if !retryableNotFound {
		return nil, err
	}

	alternateURL, altOK := c.alternateURL(reqURL.String())
	if !altOK {
		return nil, err
	}
	header, _, err = c.doJSONOnce(method, alternateURL, bodyBytes, out)
	return header, err
}

func (c *apiClient) doJSONOnce(method, requestURL string, bodyBytes []byte, out any) (http.Header, bool, error) {
	var payload io.Reader
	if bodyBytes != nil {
		payload = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, requestURL, payload)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("x-token", c.token)
	}
	if c.userID != 0 {
		req.Header.Set("x-user-id", strconv.FormatUint(uint64(c.userID), 10))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, true, fmt.Errorf("request failed: %s", resp.Status)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		var envelope apiEnvelope
		if json.Unmarshal(respBody, &envelope) == nil && envelope.Msg != "" {
			return resp.Header, false, errors.New(envelope.Msg)
		}
		return resp.Header, false, fmt.Errorf("request failed: %s", resp.Status)
	}

	var envelope apiEnvelope
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return resp.Header, false, err
	}
	if envelope.Code != 0 {
		if envelope.Msg == "" {
			envelope.Msg = "request failed"
		}
		return resp.Header, false, errors.New(envelope.Msg)
	}
	if out != nil && len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, out); err != nil {
			return resp.Header, false, err
		}
	}
	return resp.Header, false, nil
}

func (c *apiClient) alternateURL(raw string) (string, bool) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	switch {
	case strings.HasPrefix(parsed.Path, "/api/"):
		parsed.Path = strings.TrimPrefix(parsed.Path, "/api")
		if parsed.Path == "" {
			parsed.Path = "/"
		}
		return parsed.String(), true
	case parsed.Path == "/api":
		parsed.Path = "/"
		return parsed.String(), true
	default:
		parsed.Path = path.Join("/api", parsed.Path)
		return parsed.String(), true
	}
}

func parseTokenClaims(token string) (*systemReq.CustomClaims, error) {
	parser := jwt.NewParser()
	claims := &systemReq.CustomClaims{}
	if _, _, err := parser.ParseUnverified(token, claims); err != nil {
		return nil, err
	}
	return claims, nil
}
