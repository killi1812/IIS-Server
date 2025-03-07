package apiq

type InstagramUsername struct {
	Status   bool   `json:"status" xml:"status"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	Attempts string `json:"attempts"`
}

//Example
/*
{
    "status": true,
    "username": "javan",
    "user_id": "18527",
    "attempts": "1"
}
*/
type HdProfilePicURLInfo struct {
	URL string `json:"url"`
}

type BioLinks struct {
	LinkType string `json:"link_type"`
	LynxURL  string `json:"lynx_url"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

type BiographyWithEntities struct {
	Entities []Entities `json:"entities"`
}

type Entities struct {
	Hashtag interface{} `json:"hashtag"`
	User    struct {
		Username string `json:"username"`
		ID       string `json:"id"`
	} `json:"user"`
}

type UserInfo struct {
	Status                         bool                   `json:"status"`
	FullName                       string                 `json:"full_name"`
	IsMemorialized                 bool                   `json:"is_memorialized"`
	IsPrivate                      bool                   `json:"is_private"`
	HasStoryArchive                interface{}            `json:"has_story_archive"`
	Username                       string                 `json:"username"`
	IsRegulatedC18                 bool                   `json:"is_regulated_c18"`
	RegulatedNewsInLocations       []interface{}          `json:"regulated_news_in_locations"`
	TextPostAppBadgeLabel          string                 `json:"text_post_app_badge_label"`
	ShowTextPostAppBadge           bool                   `json:"show_text_post_app_badge"`
	Pk                             string                 `json:"pk"`
	LiveBroadcastVisibility        interface{}            `json:"live_broadcast_visibility"`
	LiveBroadcastID                interface{}            `json:"live_broadcast_id"`
	ProfilePicURL                  string                 `json:"profile_pic_url"`
	HdProfilePicURLInfo            *HdProfilePicURLInfo   `json:"hd_profile_pic_url_info"`
	IsUnpublished                  bool                   `json:"is_unpublished"`
	MutualFollowersCount           interface{}            `json:"mutual_followers_count"`
	ProfileContextLinksWithUserIds interface{}            `json:"profile_context_links_with_user_ids"`
	BiographyWithEntities          *BiographyWithEntities `json:"biography_with_entities"`
	AccountBadges                  []interface{}          `json:"account_badges"`
	BioLinks                       []BioLinks             `json:"bio_links"`
	ExternalLynxURL                string                 `json:"external_lynx_url"`
	ExternalURL                    string                 `json:"external_url"`
	HasChaining                    interface{}            `json:"has_chaining"`
	FbidV2                         string                 `json:"fbid_v2"`
	SupervisionInfo                interface{}            `json:"supervision_info"`
	InteropMessagingUserFbid       string                 `json:"interop_messaging_user_fbid"`
	AccountType                    int                    `json:"account_type"`
	Biography                      string                 `json:"biography"`
	IsEmbedsDisabled               bool                   `json:"is_embeds_disabled"`
	ShowAccountTransparencyDetails bool                   `json:"show_account_transparency_details"`
	IsVerified                     bool                   `json:"is_verified"`
	IsProfessionalAccount          interface{}            `json:"is_professional_account"`
	FollowerCount                  int                    `json:"follower_count"`
	AddressStreet                  interface{}            `json:"address_street"`
	CityName                       interface{}            `json:"city_name"`
	IsBusiness                     bool                   `json:"is_business"`
	Zip                            interface{}            `json:"zip"`
	Category                       string                 `json:"category"`
	ShouldShowCategory             bool                   `json:"should_show_category"`
	TransparencyLabel              interface{}            `json:"transparency_label"`
	TransparencyProduct            interface{}            `json:"transparency_product"`
	FollowingCount                 int                    `json:"following_count"`
	MediaCount                     int                    `json:"media_count"`
	LatestReelMedia                interface{}            `json:"latest_reel_media"`
	TotalClipsCount                int                    `json:"total_clips_count"`
	LatestBestiesReelMedia         interface{}            `json:"latest_besties_reel_media"`
	ReelMediaSeenTimestamp         interface{}            `json:"reel_media_seen_timestamp"`
	ID                             string                 `json:"id"`
	Attempts                       string                 `json:"attempts"`
}
