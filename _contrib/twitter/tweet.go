package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
)

type Tweet struct {
	ID            int64  `ko:"name=id"`
	IDStr         string `ko:"name=idStr"`
	CreatedAt     string `ko:"name=createdAt"`
	FavoriteCount int64  `ko:"name=favoriteCount"`
	Lang          string `ko:"name=lang"`
	RetweetCount  int64  `ko:"name=retweetCount"`
	Source        string `ko:"name=source"`
	Text          string `ko:"name=text"`
	FullText      string `ko:"name=fullText"`
	User          *User  `ko:"name=user"`
}

func NewTweets(tweets []twitter.Tweet) []*Tweet {
	r := make([]*Tweet, len(tweets))
	for i, t := range tweets {
		r[i] = NewTweet(t)
	}
	return r
}

func NewTweet(t twitter.Tweet) *Tweet {
	return &Tweet{
		ID:            t.ID,
		IDStr:         t.IDStr,
		CreatedAt:     t.CreatedAt,
		FavoriteCount: int64(t.FavoriteCount),
		Lang:          t.Lang,
		RetweetCount:  int64(t.RetweetCount),
		Source:        t.Source,
		Text:          t.Text,
		FullText:      t.FullText,
		User:          NewUser(t.User),
	}
}

type User struct {
	CreatedAt      string `ko:"name=createdAt"`
	Description    string `ko:"name=description"`
	Email          string `ko:"name=email"`
	FavoriteCount  int64  `ko:"name=favoriteCount"`
	FollowerCount  int64  `ko:"name=followerCount"`
	FollowingCount int64  `ko:"name=followingCount"`
	ID             int64  `ko:"name=id"`
	IDStr          string `ko:"name=idStr"`
	Lang           string `ko:"name=lang"`
	Name           string `ko:"name=name"`
	Screen         string `ko:"name=screen"`
	Timezone       string `ko:"name=timezone"`
	Verified       bool   `ko:"name=verified"`
	URL            string `ko:"name=url"`
}

func NewUser(u *twitter.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		CreatedAt:      u.CreatedAt,
		Description:    u.Description,
		Email:          u.Email,
		FavoriteCount:  int64(u.FavouritesCount),
		FollowerCount:  int64(u.FollowersCount),
		FollowingCount: int64(u.FriendsCount),
		ID:             u.ID,
		IDStr:          u.IDStr,
		Lang:           u.Lang,
		Name:           u.Name,
		Screen:         u.ScreenName,
		Timezone:       u.Timezone,
		Verified:       u.Verified,
		URL:            u.URL,
	}
}
