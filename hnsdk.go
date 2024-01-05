package hnsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Client struct {
	hostname string
}

func NewClient() *Client {
	return &Client{
		hostname: "https://hacker-news.firebaseio.com",
	}
}

type Item struct {
	ID          int    `json:"id"`
	Deleted     bool   `json:"deleted"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int    `json:"time"`
	Text        string `json:"text"`
	Dead        bool   `json:"dead"`
	Parent      int    `json:"parent"`
	Poll        int    `json:"poll"`
	Kids        []int  `json:"kids"`
	Url         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Parts       []int  `json:"parts"`
	Descendants int    `json:"descendants"`
}

type Items []Item

type Updates struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

type User struct {
	About     string `json:"about"`
	Created   int    `json:"created"`
	Delay     int    `json:"delay"`
	ID        string `json:"id"`
	Karma     int    `json:"karma"`
	Submitted []int  `json:"submitted"`
}

// Get item
func (c *Client) GetItem(ctx context.Context, id int) (Item, error) {
	return c.apiV0GetItem(ctx, id)
}

// Get user
func (c *Client) GetUser(ctx context.Context, username string) (User, error) {
	return c.apiV0GetUser(ctx, username)
}

// Get current largest item id . You can walk backward from here to discover all items.
func (c *Client) GetMaxItem(ctx context.Context) (int, error) {
	return c.apiV0GetMaxItem(ctx)
}

// Get up to 500 top stories. ID only
func (c *Client) GetTopStories(ctx context.Context, number int) ([]int, error) {
	if number < 1 || number > 500 {
		return []int{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := c.apiV0GetTopStories(ctx)
	if err != nil {
		return []int{}, err
	}

	return storyIDs[:number], nil
}

// Get up to 500 top stories. With data
func (c *Client) GetTopStoriesWithData(ctx context.Context, number int) (Items, error) {
	if number < 1 || number > 500 {
		return Items{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := c.apiV0GetTopStories(ctx)
	if err != nil {
		return Items{}, err
	}

	items, err := c.getItems(ctx, storyIDs, number)
	if err != nil {
		return Items{}, err
	}

	return items, err
}

// Get up to 500 new stories. ID only
func (c *Client) GetNewStories(ctx context.Context, number int) ([]int, error) {
	if number < 1 || number > 500 {
		return []int{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := c.apiV0GetNewStories(ctx)
	if err != nil {
		return []int{}, err
	}

	return storyIDs[:number], nil
}

// Get up to 500 new stories. With data
func (c *Client) GetNewStoriesWithData(ctx context.Context, number int) (Items, error) {
	if number < 1 || number > 500 {
		return Items{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := c.apiV0GetNewStories(ctx)
	if err != nil {
		return Items{}, err
	}

	items, err := c.getItems(ctx, storyIDs, number)
	if err != nil {
		return Items{}, err
	}

	return items, err
}

// Get up to 500 best stories. ID only
func (c *Client) GetBestStories(ctx context.Context, number int) ([]int, error) {
	if number < 1 || number > 500 {
		return []int{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := c.apiV0GetBestStories(ctx)
	if err != nil {
		return []int{}, err
	}

	return storyIDs[:number], nil
}

// Get up to 500 best stories. With data
func (c *Client) GetBestStoriesWithData(ctx context.Context, number int) (Items, error) {
	if number < 1 || number > 500 {
		return Items{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := c.apiV0GetBestStories(ctx)
	if err != nil {
		return Items{}, err
	}

	items, err := c.getItems(ctx, storyIDs, number)
	if err != nil {
		return Items{}, err
	}

	return items, err
}

func (c *Client) getItems(ctx context.Context, ids []int, number int) (items Items, err error) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error

	items = make([]Item, len(ids[:number]))
	for i, id := range ids[:number] {
		wg.Add(1)
		go func(i, id int, wg *sync.WaitGroup) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errs = append(errs, ctx.Err())
				return
			default:
				item, err := c.apiV0GetItem(ctx, id)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to fetch item %d: %v", id, err))
					return
				}

				mu.Lock()
				defer mu.Unlock()
				items[i] = item
			}
		}(i, id, &wg)
	}

	wg.Wait()

	if len(errs) > 0 {
		return items, fmt.Errorf("encountered %d errors: %v", len(errs), errs)
	}
	return items, nil
}

// Get up to 200 ask stories. ID only
func (c *Client) GetAskStories(ctx context.Context) ([]int, error) {
	return c.apiV0GetAskStories(ctx)
}

// Get up to 200 ask stories. With data
func (c *Client) GetAskStoriesWithData(ctx context.Context, number int) (Items, error) {
	if number < 1 || number > 200 {
		return Items{}, fmt.Errorf("accept number between 1 and 200 only")
	}

	ids, err := c.apiV0GetAskStories(ctx)
	if err != nil {
		return Items{}, err
	}

	items, err := c.getItems(ctx, ids, number)
	if err != nil {
		return Items{}, err
	}

	return items, err
}

// Get up to 200 show stories. ID only
func (c *Client) GetShowStories(ctx context.Context) ([]int, error) {
	return c.apiV0GetShowStories(ctx)
}

// Get up to 200 show stories. With data
func (c *Client) GetShowStoriesWithData(ctx context.Context, number int) (Items, error) {
	if number < 1 || number > 200 {
		return Items{}, fmt.Errorf("accept number between 1 and 200 only")
	}

	ids, err := c.apiV0GetShowStories(ctx)
	if err != nil {
		return Items{}, err
	}

	items, err := c.getItems(ctx, ids, number)
	if err != nil {
		return Items{}, err
	}

	return items, err
}

// Get up to 200 job stories. ID only
func (c *Client) GetJobStories(ctx context.Context) ([]int, error) {
	return c.apiV0GetJobStories(ctx)
}

// Get up to 200 job stories. With data
func (c *Client) GetJobStoriesWithData(ctx context.Context, number int) (Items, error) {
	if number < 1 || number > 200 {
		return Items{}, fmt.Errorf("accept number between 1 and 200 only")
	}

	ids, err := c.apiV0GetJobStories(ctx)
	if err != nil {
		return Items{}, err
	}

	items, err := c.getItems(ctx, ids, number)
	if err != nil {
		return Items{}, err
	}

	return items, err
}

// Get item and profiles changes
func (c *Client) GetUpdates(ctx context.Context) (Updates, error) {
	return c.apiV0GetUpdates(ctx)
}

func (c *Client) apiV0GetUser(ctx context.Context, username string) (User, error) {
	u := User{}
	bytes, err := c.apiCall(ctx, fmt.Sprintf("/v0/user/%s.json", username))
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(bytes, &u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (c *Client) apiV0GetItem(ctx context.Context, id int) (Item, error) {
	s := Item{}
	bytes, err := c.apiCall(ctx, fmt.Sprintf("/v0/item/%d.json", id))
	if err != nil {
		return Item{}, err
	}
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return Item{}, err
	}
	return s, nil
}

func (c *Client) apiV0GetMaxItem(ctx context.Context) (int, error) {
	s := 0
	bytes, err := c.apiCall(ctx, "/v0/maxitem.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (c *Client) apiV0GetTopStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := c.apiCall(ctx, "/v0/topstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) apiV0GetNewStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := c.apiCall(ctx, "/v0/newstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) apiV0GetBestStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := c.apiCall(ctx, "/v0/beststories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) apiV0GetAskStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := c.apiCall(ctx, "/v0/askstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) apiV0GetShowStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := c.apiCall(ctx, "/v0/showstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) apiV0GetJobStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := c.apiCall(ctx, "/v0/jobstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) apiV0GetUpdates(ctx context.Context) (Updates, error) {
	u := Updates{}
	bytes, err := c.apiCall(ctx, "/v0/updates.json")
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(bytes, &u)
	if err != nil {
		return Updates{}, err
	}

	return u, nil
}

func (c *Client) apiCall(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.hostname, url), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error http not 200")
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
