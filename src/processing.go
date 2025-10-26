package main

import (
	"log"
	"regexp"
	"strings"

	miniflux "miniflux.app/v2/client"
)

func updateEntries(client *miniflux.Client, readIds *[]int64, removedIds *[]int64) error {
	if len(*readIds) > 0 {
		if err := client.UpdateEntries(*readIds, "read"); err != nil {
			return err
		}
	}
	if len(*removedIds) > 0 {
		if err := client.UpdateEntries(*removedIds, "removed"); err != nil {
			return err
		}
	}
	return nil
}

func fmtRegex(regex string) string {
	return "(?i)" + strings.ReplaceAll(regex, " ", `\s+`)
}

func isEntryMatchingRule(entry *miniflux.Entry, rule *FilterRule) bool {
	if rule.UrlValue != "" {
		var url string
		switch rule.UrlType {
		case "site":
			url = entry.Feed.SiteURL
		case "entry":
			url = entry.URL
		case "feed":
			url = entry.Feed.FeedURL
		default:
			log.Println("unknown url value: ", rule.UrlValue)
			return false
		}
		switch rule.UrlMode {
		case "full":
			if strings.ToLower(rule.UrlValue) != strings.ToLower(url) {
				return false
			}
		case "start":
			if !strings.HasPrefix(strings.ToLower(url), strings.ToLower(rule.UrlValue)) {
				return false
			}
		case "regex":
			if match, _ := regexp.MatchString(fmtRegex(rule.UrlValue), url); !match {
				return false
			}
		default:
			log.Println("unknown url mode: ", rule.UrlMode)
			return false
		}
	}
	if rule.TitleRegex != "" {
		switch rule.FilterMode {
		case "clean":
			if match, _ := regexp.MatchString(fmtRegex(rule.TitleRegex), entry.Title); !match {
				return false
			}
		case "keep":
			if match, _ := regexp.MatchString(fmtRegex(rule.TitleRegex), entry.Title); match {
				return false
			}
		default:
			log.Println("unknown filter mode: ", rule.FilterMode)
			return false
		}
	}
	if rule.ContentRegex != "" {
		switch rule.FilterMode {
		case "clean":
			if match, _ := regexp.MatchString(fmtRegex(rule.ContentRegex), entry.Content); !match {
				return false
			}
		case "keep":
			if match, _ := regexp.MatchString(fmtRegex(rule.ContentRegex), entry.Content); match {
				return false
			}
		default:
			log.Println("unknown filter mode: ", rule.FilterMode)
			return false
		}
	}
	if rule.CategoryRegex != "" {
		if match, _ := regexp.MatchString(fmtRegex(rule.CategoryRegex), entry.Feed.Category.Title); !match {
			return false
		}
	}
	return true
}

func getHighestId(client *miniflux.Client) int64 {
	filter := &miniflux.Filter{Order: "id", Limit: 1, Direction: "desc"}
	newestEntry, err := client.Entries(filter)
	if err != nil {
		log.Panicf("failed to fetch entries: %s", err)
	}
	if newestEntry.Total == 0 {
		return -1
	} else {
		return newestEntry.Entries[0].ID
	}
}
