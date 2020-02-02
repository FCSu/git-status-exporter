package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Open an existing repository in a specific folder.
func main() {
	CheckArgs("<repo_path>", "<textfile_colletor_path>")
	repoPath := os.Args[1]
	outputPath := os.Args[2]

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(repoPath)
	CheckIfError(err)

	// We'll get the last commit just like execute `git log -1`
	Info("git log -1")

	// ... retrieving the HEAD reference
	ref, err := r.Head()
	CheckIfError(err)

	// We can verify the current status of the worktree using the method Status.
	w, err := r.Worktree()
	CheckIfError(err)
	status, err := w.Status()
	CheckIfError(err)

	tagrefs, err := r.Tags()
	CheckIfError(err)

	// Find the current tag if there is one
	tagName := ""
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if t.Hash() == ref.Hash() {
			fmt.Println("current tag:", t.Name().Short())
			tagName = t.Name().Short()
			tagrefs.Close()
		}
		return nil
	})
	CheckIfError(err)

	// ... retrieving the commit object and gather info
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)
	containsLocalChange := strings.Contains(status.String(), " M ")
	branch := ref.Name().Short()
	commitHash := commit.Hash.String()[:7]
	message := strings.TrimSuffix(commit.Message, "\n")

	if containsLocalChange {
		fmt.Println("Local Changes:")
		fmt.Println(strings.TrimSuffix(status.String(), "\n"))
	}
	fmt.Println(branch, commitHash, message)

	// ... writing to status.prom
	versionStatus := fmt.Sprintf("version_status{branch=\"%s\",revision=\"%s\",message=\"%s\",version=\"%s\",local_changes=\"%t\"} 1\n", branch, commitHash, message, tagName, containsLocalChange)
	_ = os.Mkdir(outputPath, 0700)
	err = ioutil.WriteFile(outputPath+"/status.prom", []byte(versionStatus), 0644)
	CheckIfError(err)
}
