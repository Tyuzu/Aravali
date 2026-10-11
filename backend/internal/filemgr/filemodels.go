// File: internal/filemgr/filemodels.go

package filemgr

type FileMetadata struct {
	ID        string              `db:"fileid,omitempty"`
	Hash      string              `db:"hash"`
	UserPosts map[string][]string `db:"userPosts"` // Maps userID to an array of postIDs
	PostURLs  map[string]string   `db:"postUrls"`  // Maps postID to its corresponding URL
}
