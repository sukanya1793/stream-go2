package stream

import (
	"fmt"
	"strings"
)

const (
	defaultActivityCopyLimit = 300
)

// RequestOption is an interface representing API request optional filters and
// parameters.
type RequestOption interface {
	String() string
}

type requestOption struct {
	key   string
	value interface{}
}

func makeRequestOption(key string, value interface{}) RequestOption {
	return requestOption{key: key, value: value}
}

func (o requestOption) String() string {
	return fmt.Sprintf("&%s=%v", o.key, o.value)
}

func withLimit(limit int) RequestOption {
	return makeRequestOption("limit", limit)
}

func withOffset(offset int) RequestOption {
	return makeRequestOption("offset", offset)
}

// GetActivitiesOption is an option usable by GetActivities methods for flat and aggregated feeds.
type GetActivitiesOption struct {
	RequestOption
}

// GetActivitiesWithLimit adds the limit parameter to API calls which support it, limiting
// the number of results in the response to the provided limit threshold.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func GetActivitiesWithLimit(limit int) GetActivitiesOption {
	return GetActivitiesOption{withLimit(limit)}
}

// GetActivitiesWithOffset adds the offset parameter to API calls which support it, getting
// results starting from the provided offset index.
// Supported operations: retrieve activities, retrieve followers, retrieve
// following.
func GetActivitiesWithOffset(offset int) GetActivitiesOption {
	return GetActivitiesOption{withOffset(offset)}
}

// GetActivitiesWithIDGTE adds the id_gte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater or
// equal than the provided id.
func GetActivitiesWithIDGTE(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_gte", id)}
}

// GetActivitiesWithIDGT adds the id_gt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID greater than
// the provided id.
func GetActivitiesWithIDGT(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_gt", id)}
}

// GetActivitiesWithIDLTE adds the id_lte parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser or equal
// than the provided id.
func GetActivitiesWithIDLTE(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_lte", id)}
}

// GetActivitiesWithIDLT adds the id_lt parameter to API calls, used when retrieving
// paginated activities from feeds, returning activities with ID lesser than the
// provided id.
func GetActivitiesWithIDLT(id string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("id_lt", id)}
}

func getActivitiesWithRanking(ranking string) GetActivitiesOption {
	return GetActivitiesOption{makeRequestOption("ranking", ranking)}
}

// GetNotificationWithMarkSeen marks as seen the given activity ids in a notification
// feed. If the all parameter is true, every activity in the feed is marked as seen.
func GetNotificationWithMarkSeen(all bool, activityIDs ...string) GetActivitiesOption {
	if all {
		return GetActivitiesOption{makeRequestOption("mark_seen", true)}
	}
	return GetActivitiesOption{makeRequestOption("mark_seen", strings.Join(activityIDs, ","))}
}

// GetNotificationWithMarkRead marks as read the given activity ids in a notification
// feed. If the all parameter is true, every activity in the feed is marked as read.
func GetNotificationWithMarkRead(all bool, activityIDs ...string) GetActivitiesOption {
	if all {
		return GetActivitiesOption{makeRequestOption("mark_read", true)}
	}
	return GetActivitiesOption{makeRequestOption("mark_read", strings.Join(activityIDs, ","))}
}

// FollowingOption is an option usable by following feed methods.
type FollowingOption struct {
	RequestOption
}

// FollowingWithFilter adds the filter parameter to API calls, used when retrieving
// following feeds, allowing the check whether certain feeds are being followed.
func FollowingWithFilter(ids ...string) FollowingOption {
	return FollowingOption{makeRequestOption("filter", strings.Join(ids, ","))}
}

// FollowingWithLimit limits the number of followings in the response to the provided limit.
func FollowingWithLimit(limit int) FollowingOption {
	return FollowingOption{withLimit(limit)}
}

// FollowingWithOffset returns followings starting from the given offset.
func FollowingWithOffset(offset int) FollowingOption {
	return FollowingOption{withOffset(offset)}
}

// FollowersOption is an option usable by followers feed methods.
type FollowersOption struct {
	RequestOption
}

// FollowersWithLimit limits the number of followers in the response to the provided limit.
func FollowersWithLimit(limit int) FollowersOption {
	return FollowersOption{withLimit(limit)}
}

// FollowersWithOffset returns followers starting from the given offset.
func FollowersWithOffset(offset int) FollowersOption {
	return FollowersOption{withOffset(offset)}
}

// UnfollowOption is an option usable with the Unfollow feed method.
type UnfollowOption struct {
	RequestOption
}

// UnfollowWithKeepHistory adds the `keep_history` parameter to API calls, used to keep
// history when unfollowing feeds, rather than purging it (default behavior).
// If the keepHistory parameter is false, nothing happens.
func UnfollowWithKeepHistory(keepHistory bool) UnfollowOption {
	if !keepHistory {
		return UnfollowOption{nop{}}
	}
	return UnfollowOption{makeRequestOption("keep_history", 1)}
}

type followFeedOptions struct {
	Target            string `json:"target,omitempty"`
	ActivityCopyLimit int    `json:"activity_copy_limit,omitempty"`
}

// FollowFeedOption is a function used to customize FollowFeed API calls.
type FollowFeedOption func(*followFeedOptions)

// FollowFeedWithActivityCopyLimit sets the activity copy threshold for Follow Feed API
// calls.
func FollowFeedWithActivityCopyLimit(activityCopyLimit int) FollowFeedOption {
	return func(o *followFeedOptions) {
		o.ActivityCopyLimit = activityCopyLimit
	}
}

type updateToTargetsOption func(*updateToTargetsRequest)

func updateToTargetsWithNew(targets ...Feed) updateToTargetsOption {
	return func(r *updateToTargetsRequest) {
		r.New = make([]string, len(targets))
		for i := range targets {
			r.New[i] = targets[i].ID()
		}
	}
}

func updateToTargetsWithAdd(targets ...Feed) updateToTargetsOption {
	return func(r *updateToTargetsRequest) {
		r.Adds = make([]string, len(targets))
		for i := range targets {
			r.Adds[i] = targets[i].ID()
		}
	}
}

func updateToTargetsWithRemove(targets ...Feed) updateToTargetsOption {
	return func(r *updateToTargetsRequest) {
		r.Removes = make([]string, len(targets))
		for i := range targets {
			r.Removes[i] = targets[i].ID()
		}
	}
}

type nop struct{}

func (nop) String() string {
	return ""
}
