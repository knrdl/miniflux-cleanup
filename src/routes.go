package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	miniflux "miniflux.app/v2/client"
)

func handleRetrieveConfig(writer http.ResponseWriter, request *http.Request) {
	var payload Config
	if err := readConfig(getUsername(request), &payload); err != nil {
		apiErrorResponse(writer, err, 500)
	} else {
		apiJsonResponse(writer, &payload, 200)
	}
}

func handleSaveConfig(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var payload *Config
	if err := decoder.Decode(&payload); err != nil {
		apiErrorResponse(writer, err, 400)
	} else if err := writeConfig(getUsername(request), payload); err != nil {
		apiErrorResponse(writer, err, 500)
	} else {
		requestCronConfigUpdate()
		apiTxtResponse(writer, "", 204)
	}
}

func handlePreview(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var payload *Config
	if err := decoder.Decode(&payload); err != nil {
		apiErrorResponse(writer, err, 400)
	} else {
		if data, err := queryPreview(payload); err != nil {
			apiErrorResponse(writer, err, 500)
		} else {
			apiTxtResponse(writer, data, 200)
		}
	}
}

func queryPreview(payload *Config) (string, error) {
	output := ""
	client := miniflux.NewClient(getMinifluxUrl(), payload.ApiKey)
	var readIds []int64
	var removedIds []int64
	entryResultSet, err := client.Entries(&miniflux.Filter{Starred: miniflux.FilterNotStarred, Direction: "desc", Status: "unread", Limit: EntryProcessingLimit, Order: "id"})
	if err != nil {
		return "", err
	}
	for _, entry := range entryResultSet.Entries {
		for _, rule := range payload.Rules {
			if isEntryMatchingRule(entry, &rule) {
				if !contains(readIds, entry.ID) && !contains(removedIds, entry.ID) {
					switch rule.State {
					case "read":
						readIds = append(readIds, entry.ID)
					case "removed":
						removedIds = append(removedIds, entry.ID)
					default:
						return "", fmt.Errorf("unknown target state: %s", rule.State)
					}
					output += fmt.Sprintf("[%s] %s (%s)\n", rule.State, entry.Title, entry.Feed.Title)
				}
			}
		}
	}
	return output, nil
}
