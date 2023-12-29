package spot

type Playlist struct {
	ID   string
	Name string
}

type Track struct {
	ID      string
	Title   string
	Artists []string
	Album   string
}

type Album struct {
	ID      string
	Title   string
	Artists []string
	Tracks  []Track
}
