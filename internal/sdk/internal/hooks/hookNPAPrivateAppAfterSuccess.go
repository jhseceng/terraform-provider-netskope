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
	if hookCtx.OperationID == "updateNPAPrivateApp" {
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

		// Try to unmarshal to see the structure
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

		// Check if response has "data" field (single object response)
		if data, hasData := responseMap["data"]; hasData {
			if npaPrivateAppUpdateResponseDebug {
				log.Printf("Response has 'data' field, checking if it needs to be wrapped in array...")
			}

			// Check if data is already an array
			if _, isArray := data.([]interface{}); !isArray {
				// Wrap single object in an array
				if npaPrivateAppUpdateResponseDebug {
					log.Printf("Wrapping single data object in array for compatibility")
				}
				responseMap["data"] = []interface{}{data}

				// Marshal back to JSON
				modifiedBody, err := json.Marshal(responseMap)
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
		} else {
			// Response might already be an array at the top level
			// Check if it's an object that should be wrapped
			if _, hasStatus := responseMap["status"]; hasStatus {
				if npaPrivateAppUpdateResponseDebug {
					log.Printf("Response is a single object, wrapping in array")
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
		}

		// If no modifications needed, return original response
		if npaPrivateAppUpdateResponseDebug {
			log.Printf("No response transformation needed")
		}
		res.Body = io.NopCloser(strings.NewReader(string(body)))
		return res, nil
	}
	return res, nil
}
