package domain

type Profile struct {
	Created             int64    `bson:"created" json:"created"`
	Modified            int64    `bson:"modified" json:"modified"`
	ProfileDescription  string   `json:"description"`
	Id                  string   `json:"id"`
	ProfileName         string   `json:"name"`
	Profilemanufacturer string   `json:"manufacturer"`
	Profilemodel        string   `json:"model"`
	ProfileLabels       []string `json:"labels"`
}
