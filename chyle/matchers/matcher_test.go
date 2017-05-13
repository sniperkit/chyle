package matchers

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"srcd.works/go-git.v4/plumbing"
	"srcd.works/go-git.v4/plumbing/object"
)

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Close() (err error) {
	return nil
}

// encodedObject provides an encoded version of a commit
// from a commit, it's possible to fake a commit from this structure
type encodedObject struct {
	typ  plumbing.ObjectType
	size int64
	buf  *buffer
	hash plumbing.Hash
}

func (e *encodedObject) Hash() plumbing.Hash {
	return e.hash
}

func (e *encodedObject) Type() plumbing.ObjectType {
	return e.typ
}

func (e *encodedObject) SetType(typ plumbing.ObjectType) {
	e.typ = typ
}

func (e *encodedObject) Size() int64 {
	return e.size
}

func (e *encodedObject) SetSize(size int64) {
	e.size = size
}

func (e *encodedObject) Reader() (io.ReadCloser, error) {
	return e.buf, nil
}

func (e *encodedObject) Writer() (io.WriteCloser, error) {
	return e.buf, nil
}

func TestMatchersMergeCommits(t *testing.T) {
	buf1 := &buffer{}
	buf1.WriteString("parent 10a34637ad661d98ba3344717656fcc76209c2f8\n")
	buf1.WriteString("parent 3e6c06b1a28a035e21aa0a736ef80afadc43122c\n")
	encodedMergeCommit1 := &encodedObject{
		buf:  buf1,
		hash: plumbing.NewHash("3c7435cfd4e31b9be3991041c9a4f8292b752e5b"),
	}
	encodedMergeCommit1.SetType(plumbing.CommitObject)

	buf2 := &buffer{}
	buf2.WriteString("parent 2503b2b99c061ff5bac94f2c1972e4c28cf1a844\n")
	buf2.WriteString("parent ee7fbf1c52742cf4f30d00b0e9e477dde72c7e51\n")
	encodedMergeCommit2 := &encodedObject{
		buf:  buf2,
		hash: plumbing.NewHash("ecc1978dca2e31d10751ede8d8753f1cbded832e"),
	}
	encodedMergeCommit2.SetType(plumbing.CommitObject)

	commit1 := object.Commit{}
	err := commit1.Decode(encodedMergeCommit1)
	assert.NoError(t, err)

	commit2 := object.Commit{Hash: plumbing.NewHash("ea8aa7337c39b717fcdff0c858027b9778ab391a")}

	commit3 := object.Commit{}
	err = commit3.Decode(encodedMergeCommit2)
	assert.NoError(t, err)

	commit4 := object.Commit{Hash: plumbing.NewHash("f75edc98db49e7c13f818f70c418087256354303")}
	commits := []object.Commit{commit1, commit2, commit3, commit4}

	cs := Filter(&[]Matcher{mergeCommitMatcher{}}, &commits)

	assert.Len(t, *cs, 2)
	assert.Equal(t, commit1.Hash.String(), (*cs)[0]["id"])
	assert.Equal(t, commit3.Hash.String(), (*cs)[1]["id"])
}

func TestMatchersRegularCommits(t *testing.T) {
	buf1 := &buffer{}
	buf1.WriteString("parent da39a3ee5e6b4b0d3255bfef95601890afd80709\n")
	buf1.WriteString("parent 10a34637ad661d98ba3344717656fcc76209c2f8\n")
	encodedMergeCommit1 := &encodedObject{
		buf:  buf1,
		hash: plumbing.NewHash("3e6c06b1a28a035e21aa0a736ef80afadc43122c"),
	}
	encodedMergeCommit1.SetType(plumbing.CommitObject)

	buf2 := &buffer{}
	buf2.WriteString("parent ecc1978dca2e31d10751ede8d8753f1cbded832e\n")
	buf2.WriteString("parent 6110212c651287fa93aae5142e60b50edde00970\n")
	encodedMergeCommit2 := &encodedObject{
		buf:  buf2,
		hash: plumbing.NewHash("63027d7630360e4203c0e3f970ec2ffcfe5f8f1b"),
	}
	encodedMergeCommit2.SetType(plumbing.CommitObject)

	commit1 := object.Commit{}
	err := commit1.Decode(encodedMergeCommit1)
	assert.NoError(t, err)

	commit2 := object.Commit{Hash: plumbing.NewHash("f3226f91f77a87d909b8920adc91f9a301a7316b")}

	commit3 := object.Commit{}
	err = commit3.Decode(encodedMergeCommit2)
	assert.NoError(t, err)

	commit4 := object.Commit{Hash: plumbing.NewHash("ea8aa7337c39b717fcdff0c858027b9778ab391a")}
	commits := []object.Commit{commit1, commit2, commit3, commit4}

	cs := Filter(&[]Matcher{regularCommitMatcher{}}, &commits)

	assert.Len(t, *cs, 2)
	assert.Equal(t, commit2.Hash.String(), (*cs)[0]["id"])
	assert.Equal(t, commit4.Hash.String(), (*cs)[1]["id"])
}

func TestMatchersWithCommitMessage(t *testing.T) {
	re, err := regexp.Compile(`whatever.*`)

	assert.NoError(t, err)

	commit1 := object.Commit{Message: "test1"}
	commit2 := object.Commit{Message: "whatever1"}
	commit3 := object.Commit{Message: "test2"}
	commit4 := object.Commit{Message: "whatever2"}
	commits := []object.Commit{commit1, commit2, commit3, commit4}

	cs := Filter(&[]Matcher{messageMatcher{re}}, &commits)

	assert.Len(t, *cs, 2)
	assert.Equal(t, commit2.Message, (*cs)[0]["message"])
	assert.Equal(t, commit4.Message, (*cs)[1]["message"])
}

func TestMatchersWithAuthor(t *testing.T) {
	re, err := regexp.Compile(".*whatever.*")

	assert.NoError(t, err)

	commit1 := object.Commit{Author: object.Signature{Email: "test@test.com"}}
	commit2 := object.Commit{Author: object.Signature{Email: "whatever@test.com"}}
	commit3 := object.Commit{Author: object.Signature{Email: "test2@test.com"}}
	commit4 := object.Commit{Author: object.Signature{Email: "whatever2@test.com"}}
	commits := []object.Commit{commit1, commit2, commit3, commit4}

	cs := Filter(&[]Matcher{authorMatcher{re}}, &commits)
	assert.Len(t, *cs, 2)
	assert.Equal(t, commit2.Author.Email, (*cs)[0]["authorEmail"])
	assert.Equal(t, commit4.Author.Email, (*cs)[1]["authorEmail"])
}

func TestMatchersWithCommitter(t *testing.T) {
	re, err := regexp.Compile(".*whatever.*")

	assert.NoError(t, err)

	commit1 := object.Commit{Committer: object.Signature{Email: "test@test.com"}}
	commit2 := object.Commit{Committer: object.Signature{Email: "whatever@test.com"}}
	commit3 := object.Commit{Committer: object.Signature{Email: "test2@test.com"}}
	commit4 := object.Commit{Committer: object.Signature{Email: "whatever2@test.com"}}
	commits := []object.Commit{commit1, commit2, commit3, commit4}

	cs := Filter(&[]Matcher{committerMatcher{re}}, &commits)
	assert.Len(t, *cs, 2)
	assert.Equal(t, commit2.Committer.Email, (*cs)[0]["committerEmail"])
	assert.Equal(t, commit4.Committer.Email, (*cs)[1]["committerEmail"])
}

func TestTransformCommitsToMap(t *testing.T) {
	commit1 := object.Commit{}
	commits := []object.Commit{commit1}
	commitMaps := transformCommitsToMap(&commits)

	expected := map[string]interface{}{
		"id":             commit1.ID().String(),
		"authorName":     commit1.Author.Name,
		"authorEmail":    commit1.Author.Email,
		"authorDate":     commit1.Author.When.String(),
		"committerName":  commit1.Committer.Name,
		"committerEmail": commit1.Committer.Email,
		"committerDate":  commit1.Committer.When.String(),
		"message":        commit1.Message,
		"type":           "regular",
	}

	assert.Len(t, *commitMaps, 1)
	assert.Equal(t, expected, (*commitMaps)[0])
}

func TestCreateMatchers(t *testing.T) {
	matchers := map[string]string{}
	matchers["TYPE"] = "regular"
	matchers["MESSAGE"] = ".*"
	matchers["AUTHOR"] = ".*"
	matchers["COMMITTER"] = ".*"

	m := CreateMatchers(matchers)

	assert.Len(t, *m, 4, "Must contain 4 matchers")
}