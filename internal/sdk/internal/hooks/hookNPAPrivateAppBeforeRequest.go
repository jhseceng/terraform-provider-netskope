package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type npaPrivateAppUpdateRequest struct {
	ID                       *int                   `json:"id,omitempty"`
	AllowUnauthenticatedCors *bool                  `json:"allow_unauthenticated_cors,omitempty"`
	UribypassHeaderValue     *string                `json:"uribypass_header_value,omitempty"`
	AppOption                map[string]interface{} `json:"app_option,omitempty"`
	ClientlessAccess         *bool                  `json:"clientless_access,omitempty"`
	PrivateAppHostname       *string                `json:"host,omitempty"`
	IsUserPortalApp          *bool                  `json:"is_user_portal_app,omitempty"`
	Protocols                []interface{}          `json:"protocols,omitempty"`
	PublisherTags            []interface{}          `json:"publisher_tags,omitempty"`
	Publishers               []interface{}          `json:"publishers,omitempty"`
	RealHost                 *string                `json:"real_host,omitempty"`
	Tags                     []interface{}          `json:"tags,omitempty"`
	TrustSelfSignedCerts     *bool                  `json:"trust_self_signed_certs,omitempty"`
	UsePublisherDNS          *bool                  `json:"use_publisher_dns,omitempty"`
}

var (
	_                            beforeRequestHook = (*npaPrivateAppUpdateRequest)(nil)
	npaPrivateAppUpdateDebug     bool              = true
)

func (i *npaPrivateAppUpdateRequest) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if hookCtx.OperationID == "updateNPAPrivateApp" {
		if npaPrivateAppUpdateDebug {
			log.Print("Executing BeforeRequest hook for NPA Private App update...")
		}

		// Extract the private_app_id from the URL path
		re := regexp.MustCompile(`/steering/apps/private/(\d+)`)
		matches := re.FindStringSubmatch(req.URL.Path)
		if len(matches) < 2 {
			log.Printf("ERROR: Unable to extract private_app_id from URL path: %s", req.URL.Path)
			return req, nil // Continue without modification if ID can't be extracted
		}

		privateAppID, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Printf("ERROR: Unable to convert private_app_id to int: %v", err)
			return req, nil
		}

		if npaPrivateAppUpdateDebug {
			log.Printf("Extracted private_app_id: %d", privateAppID)
		}

		// Read the request body
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("ERROR: Unable to read request body: %v", err)
			return nil, fmt.Errorf("ERROR: Unable to read request body: %w", err)
		}

		if npaPrivateAppUpdateDebug {
			log.Printf("Original request body: %s", string(body))
		}

		// Unmarshal into our struct
		var requestMap npaPrivateAppUpdateRequest
		if err := json.Unmarshal(body, &requestMap); err != nil {
			log.Printf("ERROR: Unable to unmarshal request: %v", err)
			return nil, fmt.Errorf("ERROR: Unable to unmarshal request: %v", err)
		}

		// Add the ID field
		requestMap.ID = &privateAppID

		if npaPrivateAppUpdateDebug {
			log.Printf("Added ID field to request: %d", privateAppID)
		}

		// Marshal back to JSON
		modifiedBody, err := json.Marshal(requestMap)
		if err != nil {
			log.Printf("ERROR: Unable to marshal modified request: %v", err)
			return nil, fmt.Errorf("ERROR: Unable to marshal modified request: %v", err)
		}

		if npaPrivateAppUpdateDebug {
			log.Printf("Modified request body: %s", string(modifiedBody))
		}

		// Update the request with modified body
		req.Body = io.NopCloser(strings.NewReader(string(modifiedBody)))
		req.ContentLength = int64(len(modifiedBody))

		// Change method from PATCH to PUT
		req.Method = "PUT"

		if npaPrivateAppUpdateDebug {
			log.Printf("Changed HTTP method to: %s", req.Method)
		}

		return req, nil
	}
	return req, nil
}
