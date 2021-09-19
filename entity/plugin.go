package entity

type PluginUIComponent struct {
	Name        string   `json:"name" bson:"name"`
	Title       string   `json:"title" bson:"title"`
	Src         string   `json:"src" bson:"src"`
	Type        string   `json:"type" bson:"type"`
	Path        string   `json:"path" bson:"path"`
	ParentPaths []string `json:"parent_paths" bson:"parent_paths"`
}

type PluginUINav struct {
	Path  string   `json:"path" bson:"path"`
	Title string   `json:"title" bson:"title"`
	Icon  []string `json:"icon" bson:"icon"`
}

type PluginUIAsset struct {
	Path string `json:"path" bson:"path"`
	Type string `json:"type" bson:"type"`
}

type PluginEventKey struct {
	Include string `json:"include" bson:"include"`
	Exclude string `json:"exclude" bson:"exclude"`
}
