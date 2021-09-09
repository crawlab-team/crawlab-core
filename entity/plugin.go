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
