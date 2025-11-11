package hooks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type npaPrivateAppUpdateResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Status string      `json:"status,omitempty"`
}

var (
	_                                afterSuccessHook = (*npaPrivateAppUpdateResponse)(nil)
	npaPrivateAppUpdateResponseDebug bool             = true
)

func (i *npaPrivateAppUpdateResponse) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if npaPrivateAppUpdateResponseDebug {
		log.Printf("DEBUG: npaPrivateAppUpdateResponse hook called with OperationID: %s", hookCtx.OperationID)
	}

	// Only process updateNPAPrivateApp operations
	if hookCtx.OperationID != "updateNPAPrivateApp" {
		if npaPrivateAppUpdateResponseDebug {
			log.Printf("DEBUG: Skipping hook - OperationID '%s' does not match 'updateNPAPrivateApp'", hookCtx.OperationID)
		}
		return res, nil
	}

	if npaPrivateAppUpdateResponseDebug {
		log.Print("Executing AfterSuccess hook for NPA Private App update...")
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("ERROR: Unable to read response body: %v", err)
		return nil, fmt.Errorf("ERROR: Unable to read response body: %w", err)
	}

	if npaPrivateAppUpdateResponseDebug {
		log.Printf("Original response body: %s", string(body))
		log.Printf("Response status code: %d", res.StatusCode)
	}

	// First, check if the response is already an array
	var arrayTest []interface{}
	if err := json.Unmarshal(body, &arrayTest); err == nil {
		// Already an array, no transformation needed
		if npaPrivateAppUpdateResponseDebug {
			log.Printf("Response is already an array, no transformation needed")
		}
		res.Body = io.NopCloser(strings.NewReader(string(body)))
		return res, nil
	}

	// Not an array, try to unmarshal as object
	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		log.Printf("ERROR: Unable to unmarshal response: %v", err)
		// Return original response if we can't parse it
		res.Body = io.NopCloser(strings.NewReader(string(body)))
		return res, nil
	}

	if npaPrivateAppUpdateResponseDebug {
		log.Printf("Parsed response structure: %+v", responseMap)
	}

	// Response is a single object, we need to wrap it in an array
	if npaPrivateAppUpdateResponseDebug {
		log.Printf("Response is a single object, wrapping in array for compatibility")
	}

	arrayResponse := []interface{}{responseMap}
	modifiedBody, err := json.Marshal(arrayResponse)
	if err != nil {
		log.Printf("ERROR: Unable to marshal modified response: %v", err)
		res.Body = io.NopCloser(strings.NewReader(string(body)))
		return res, nil
	}

	if npaPrivateAppUpdateResponseDebug {
		log.Printf("Modified response body: %s", string(modifiedBody))
	}

	res.Body = io.NopCloser(strings.NewReader(string(modifiedBody)))
	return res, nil
}
