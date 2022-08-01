package spotify

import "fmt"

type Client struct {
	httpClient HttpClient
	id         string
	secret     string
	token      string
}

func NewClient(id, secret string) Client {
	return Client{
		id:     id,
		secret: secret,
	}
}

func (c *Client) Authorize() error {
	t, err := retrieveAuthToken(c.httpClient, c.id, c.secret)
	if err != nil {
		return fmt.Errorf("failed to authorize spotify client, %w", err)
	}

	c.token = t

	return nil
}
