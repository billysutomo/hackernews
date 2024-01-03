package hnsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type hackernews struct {
	hostname string
}

func HN() *hackernews {
	return &hackernews{
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

func (i Item) toStory() Story {
	return Story{
		By:          i.By,
		Descendants: i.Descendants,
		ID:          i.ID,
		Kids:        i.Kids,
		Score:       i.Score,
		Time:        i.Time,
		Title:       i.Title,
		Type:        i.Type,
		URL:         i.Url,
	}
}

func (i Item) toComment() Comment {
	return Comment{
		By:     i.By,
		ID:     i.ID,
		Kids:   i.Kids,
		Parent: i.Parent,
		Text:   i.Text,
		Time:   i.Time,
		Type:   i.Type,
	}
}

func (i Item) toAsk() Ask {
	return Ask{
		By:          i.By,
		Descendants: i.Descendants,
		ID:          i.ID,
		Kids:        i.Kids,
		Score:       i.Score,
		Text:        i.Text,
		Time:        i.Time,
		Title:       i.Title,
		Type:        i.Type,
	}
}

func (i Item) toJob() Job {
	return Job{
		By:    i.By,
		ID:    i.ID,
		Score: i.Score,
		Text:  i.Text,
		Time:  i.Time,
		Title: i.Title,
		Type:  i.Type,
		URL:   i.Url,
	}
}

func (i Item) toPoll() Poll {
	return Poll{
		By:          i.By,
		Descendants: i.Descendants,
		ID:          i.ID,
		Kids:        i.Kids,
		Parts:       i.Parts,
		Score:       i.Score,
		Text:        i.Text,
		Time:        i.Time,
		Title:       i.Title,
		Type:        i.Type,
	}
}

func (i Item) toPollOpt() PollOpt {
	return PollOpt{
		By:    i.By,
		ID:    i.ID,
		Poll:  i.Poll,
		Score: i.Score,
		Text:  i.Text,
		Time:  i.Time,
		Type:  i.Type,
	}
}

type Story struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

type Stories []Story

type Comment struct {
	By     string `json:"by"`
	ID     int    `json:"id"`
	Kids   []int  `json:"kids"`
	Parent int    `json:"parent"`
	Text   string `json:"text"`
	Time   int    `json:"time"`
	Type   string `json:"type"`
}

type Ask struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Text        string `json:"text"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
}

type Job struct {
	By    string `json:"by"`
	ID    int    `json:"id"`
	Score int    `json:"score"`
	Text  string `json:"text"`
	Time  int    `json:"time"`
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type Poll struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Parts       []int  `json:"parts"`
	Score       int    `json:"score"`
	Text        string `json:"text"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
}

type PollOpt struct {
	By    string `json:"by"`
	ID    int    `json:"id"`
	Poll  int    `json:"poll"`
	Score int    `json:"score"`
	Text  string `json:"text"`
	Time  int    `json:"time"`
	Type  string `json:"type"`
}

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

func (hn hackernews) GetUser(ctx context.Context, username string) (User, error) {
	return hn.apiV0GetUser(ctx, username)
}

func (hn hackernews) GetMaxItem(ctx context.Context) (int, error) {
	return hn.apiV0GetMaxItem(ctx)
}

func (hn hackernews) GetStory(ctx context.Context, id int) (Story, error) {
	item, err := hn.apiV0GetItem(ctx, id)
	if err != nil {
		return Story{}, err
	}

	return item.toStory(), nil
}

func (hn hackernews) GetComment(ctx context.Context, id int) (Comment, error) {
	item, err := hn.apiV0GetItem(ctx, id)
	if err != nil {
		return Comment{}, err
	}

	return item.toComment(), nil
}

func (hn hackernews) GetAsk(ctx context.Context, id int) (Ask, error) {
	item, err := hn.apiV0GetItem(ctx, id)
	if err != nil {
		return Ask{}, err
	}

	return item.toAsk(), nil
}

func (hn hackernews) GetJob(ctx context.Context, id int) (Job, error) {
	item, err := hn.apiV0GetItem(ctx, id)
	if err != nil {
		return Job{}, err
	}

	return item.toJob(), nil
}

func (hn hackernews) GetPoll(ctx context.Context, id int) (Poll, error) {
	item, err := hn.apiV0GetItem(ctx, id)
	if err != nil {
		return Poll{}, err
	}

	return item.toPoll(), nil
}

func (hn hackernews) GetPollOpt(ctx context.Context, id int) (PollOpt, error) {
	item, err := hn.apiV0GetItem(ctx, id)
	if err != nil {
		return PollOpt{}, err
	}

	return item.toPollOpt(), nil
}

func (hn hackernews) GetTopStories(ctx context.Context, number int) ([]int, error) {
	if number < 1 || number > 500 {
		return []int{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := hn.apiV0GetTopStories(ctx)
	if err != nil {
		return []int{}, err
	}

	return storyIDs[:number], nil
}

func (hn hackernews) GetTopStoriesWithData(ctx context.Context, number int) (Stories, error) {
	if number < 1 || number > 500 {
		return Stories{}, fmt.Errorf("accept number between 1 and 500 only")
	}

	storyIDs, err := hn.apiV0GetTopStories(ctx)
	if err != nil {
		return Stories{}, err
	}

	stories := make(Stories, number)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var errs []error

	for i, id := range storyIDs[:number] {
		wg.Add(1)
		go func(i, id int, wg *sync.WaitGroup) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errs = append(errs, ctx.Err())
				return
			default:
				item, err := hn.apiV0GetItem(ctx, id)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to fetch item %d: %v", id, err))
					return
				}

				mu.Lock()
				defer mu.Unlock()
				stories[i] = item.toStory()
			}
		}(i, id, &wg)
	}

	wg.Wait()

	if len(errs) > 0 {
		return stories, fmt.Errorf("encountered %d errors: %v", len(errs), errs)
	}

	return stories, err
}

func (hn hackernews) GetAskStories(ctx context.Context) ([]int, error) {
	return hn.apiV0GetAskStories(ctx)
}

func (hn hackernews) GetUpdates(ctx context.Context) (Updates, error) {
	return hn.apiV0GetUpdates(ctx)
}

func (hn hackernews) apiV0GetUser(ctx context.Context, username string) (User, error) {
	u := User{}
	bytes, err := hn.apiCall(ctx, fmt.Sprintf("/v0/user/%s.json", username))
	if err != nil {
		return User{}, err
	}
	err = json.Unmarshal(bytes, &u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (hn hackernews) apiV0GetItem(ctx context.Context, id int) (Item, error) {
	s := Item{}
	bytes, err := hn.apiCall(ctx, fmt.Sprintf("/v0/item/%d.json", id))
	if err != nil {
		return Item{}, err
	}
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return Item{}, err
	}
	return s, nil
}

func (hn hackernews) apiV0GetMaxItem(ctx context.Context) (int, error) {
	s := 0
	bytes, err := hn.apiCall(ctx, "/v0/maxitem.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (hn hackernews) apiV0GetTopStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := hn.apiCall(ctx, "/v0/topstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (hn hackernews) apiV0GetAskStories(ctx context.Context) ([]int, error) {
	s := []int{}
	bytes, err := hn.apiCall(ctx, "/v0/askstories.json")
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return s, nil

}
func (hn hackernews) apiV0GetUpdates(ctx context.Context) (Updates, error) {
	u := Updates{}
	bytes, err := hn.apiCall(ctx, "/v0/updates.json")
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(bytes, &u)
	if err != nil {
		return Updates{}, err
	}

	return u, nil
}

func (hn hackernews) apiCall(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", hn.hostname, url), nil)
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
