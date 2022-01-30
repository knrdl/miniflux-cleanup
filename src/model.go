package main

type FilterRule struct {
	UrlType       string `json:"url_type"`
	UrlMode       string `json:"url_mode"`
	UrlValue      string `json:"url_value"`
	FilterMode    string `json:"filter_mode"`
	TitleRegex    string `json:"title_regex"`
	ContentRegex  string `json:"content_regex"`
	CategoryRegex string `json:"category_regex"`
	State         string `json:"state"`
}

type Config struct {
	Rules  []FilterRule `json:"rules"`
	ApiKey string       `json:"api_key"`
}
