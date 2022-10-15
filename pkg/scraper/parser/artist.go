package parser

// ArtistAlbums represents a mapping artist -> albums
type ArtistAlbums map[string][]string

func (m ArtistAlbums) ResolveArtist(album string) (string, bool) {
	for artist, albums := range m {
		for _, a := range albums {
			if album == a {
				return artist, true
			}
		}
	}

	return "", false
}
