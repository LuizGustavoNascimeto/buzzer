package activities

import "time"

func NewMockActivities() []*Activity {
	parentCreatedAt := time.Date(2026, 5, 27, 1, 18, 39, 679037000, time.UTC)
	parentExpiresAt := time.Date(2026, 6, 3, 1, 18, 39, 679037000, time.UTC)
	replyCreatedAt := time.Date(2026, 5, 27, 1, 18, 39, 679037000, time.UTC)
	garekCreatedAt := time.Date(2026, 5, 29, 0, 18, 39, 679037000, time.UTC)
	garekExpiresAt := time.Date(2026, 5, 29, 13, 18, 39, 679037000, time.UTC)
	worfCreatedAt := time.Date(2026, 5, 22, 1, 18, 39, 679037000, time.UTC)
	worfExpiresAt := time.Date(2026, 6, 7, 1, 18, 39, 679037000, time.UTC)

	replyToParent := "68f126b0-1ceb-4a33-88be-d90fa7109eee"

	return []*Activity{
		{
			UUID:         "68f126b0-1ceb-4a33-88be-d90fa7109eee",
			Handle:       "Andrew Brown",
			Message:      "Cloud is fun! #cloud #aws #azure #gcp",
			LikesCount:   55,
			RepliesCount: 1,
			RepostsCount: 0,
			Replies: []Activity{
				{
					UUID:                "26e12864-1c26-5c3a-9658-97a10f8fea67",
					Handle:              "Worf",
					Message:             "This post has no honor!",
					LikesCount:          0,
					RepliesCount:        0,
					RepostsCount:        0,
					ReplyToActivityUUID: &replyToParent,
					CreatedAt:           replyCreatedAt,
				},
			},
			ExpiresAt: &parentExpiresAt,
			CreatedAt: parentCreatedAt,
		},
		{
			UUID:         "66e12864-8c26-4c3a-9658-95a10f8fea67",
			Handle:       "Worf",
			Message:      "I am out of prune juice",
			LikesCount:   0,
			RepliesCount: 0,
			RepostsCount: 0,
			Replies:      []Activity{},
			ExpiresAt:    &worfExpiresAt,
			CreatedAt:    worfCreatedAt,
		},
		{
			UUID:         "248959df-3079-4947-b847-9e0892d1bab4",
			Handle:       "Garek",
			Message:      "My dear doctor, I am just simple tailor",
			LikesCount:   0,
			RepliesCount: 0,
			RepostsCount: 0,
			Replies:      []Activity{},
			ExpiresAt:    &garekExpiresAt,
			CreatedAt:    garekCreatedAt,
		},
	}
}
