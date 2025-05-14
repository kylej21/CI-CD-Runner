package jobs

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Job struct {
	ID      string `json:"id,omitempty"`
	RepoURL string `json:"repo_url"`
	Branch  string `json:"branch,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (j *Job) Start() {
	j.Status = "running"
	targetPath := filepath.Join("tmp/", j.ID)
	_ = os.MkdirAll("tmp", os.ModePerm)

	folderExists := false
	_, err := os.Stat(targetPath)
	if err == nil {
		folderExists = true
	}

	if folderExists {
		log.Printf("Found existing folder. Removing: %s", targetPath)
		removeStatus := os.RemoveAll(targetPath)
		if removeStatus != nil {
			log.Printf("Error removing %s", targetPath)
			j.Status = "failed"
			return
		}
	}

	var cmd *exec.Cmd = exec.Command("git", "clone", j.RepoURL, targetPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("clone failed: %v", err)
		log.Printf("git output: %s", string(output))
		j.Status = "failed"
		return
	}
	log.Printf("cloned! %s", output)
	j.Status = "cloned"
}

func (j *Job) Run() {
	j.Status = "executing_setup"
	repository := "tmp/" + j.ID
	log.Println(repository)
	cmdInstall := exec.Command("npm", "install")
	cmdInstall.Dir = repository
	output, err := cmdInstall.CombinedOutput()
	if err != nil {
		log.Printf("npm install failed: %+v", err)
		log.Printf("terminal output: %v", output)
		j.Status = "failed"
	}
	log.Printf("npm install success!")

	cmdInstall = exec.Command("npm", "start")
	cmdInstall.Dir = repository
	output, err = cmdInstall.CombinedOutput()
	if err != nil {
		log.Printf("npm start failed: %+v", err)
		log.Printf("terminal output: %v", output)
		j.Status = "failed"
		return
	}
	log.Printf("npm start success!")
	j.Status = "built"
}
