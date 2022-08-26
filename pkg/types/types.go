package types

type Video struct {
	Id          string `json:"_id,omitempty" bson:"_id,omitempty"`
	UniqueId    string `json:"uniqueId" bson:"uniqueId"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	PublishedAt string `json:"publishedAt" bson:"publishedAt"`
}
