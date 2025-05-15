package jobs

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Job struct {
	ID           string        `json:"id,omitempty"`
	RepoURL      string        `json:"repo_url"`
	Branch       string        `json:"branch,omitempty"`
	Status       string        `json:"status,omitempty"`
	CreatedAt    time.Time     `json:"created_at,omitempty"`
	FinishedAt   time.Time     `json:"finished_at,omitempty"`
	DurationMs   int64         `json:"duration_ms,omitempty"`
	LintErrors   int           `json:"lint_errors,omitempty"`
	LintWarnings int           `json:"lint_warnings,omitempty"`
	BuildLog     string        `json:"build_log,omitempty"`
	LintResults  []LintResults `json:"lint_results,omitempty"`
}

type LintResults struct {
	FilePath     string        `json:"filePath"`
	ErrorCount   int           `json:"errorCount"`
	WarningCount int           `json:"warningCount"`
	Messages     []LintMessage `json:"messages"`
}

type LintMessage struct {
	RuleID   string `json:"ruleId"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

func (j *Job) Run() {

	j.Clone()
	if j.Status == "failed" {
		return
	}

	j.Install()
	if j.Status == "failed" {
		return
	}

	j.Lint()
	if j.Status == "failed" {
		return
	}

	j.Build()
}

func (j *Job) Clone() {
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
		log.Printf("npm clone failed: %+v", err)
		log.Printf("terminal output: \n%s", string(output))
		j.Status = "failed"
		return
	}
	log.Printf("cloned! %s", output)
	j.Status = "cloned"

}

func (j *Job) Install() {
	j.Status = "installing"
	repository := "tmp/" + j.ID
	log.Println(repository)
	cmdInstall := exec.Command("npm", "install")
	cmdInstall.Dir = repository
	output, err := cmdInstall.CombinedOutput()
	if err != nil {
		log.Printf("npm install failed: %+v", err)
		log.Printf("terminal output: \n%s", string(output))
		j.Status = "failed"
		return
	}
	log.Printf("npm install success!")

	j.Status = "installing_dependencies"
	cmdInstall = exec.Command("npm", "install", "--save-dev", "eslint@8", "vite")
	output, err = cmdInstall.CombinedOutput()
	if err != nil {
		log.Printf("npm install --save-dev eslint failed: %+v", err)
		log.Printf("terminal output: \n%s", string(output))
		j.Status = "failed"
		return
	}
	log.Printf("npm install --save-dev eslint success!")

	j.Status = "install_complete"

}

func (j *Job) Lint() {
	j.Status = "linting"
	repository := "tmp/" + j.ID

	// manual rule declaration in case no eslint.rc file exists
	cmdLint := exec.Command(
		"npx", "eslint", ".",
		"--ext", ".vue",
		"--ext", ".js",
		"--ext", ".jsx",
		"--ext", ".cjs",
		"--ext", ".mjs",
		"--ext", ".ts",
		"--ext", ".tsx",
		"--ext", ".cts",
		"--ext", ".mts",
		"--fix",
		"--no-eslintrc",
		"--rule", "quotes:[\"error\",\"double\"]",
		"--format", "json",
	)
	cmdLint.Dir = repository
	output, err := cmdLint.CombinedOutput()
	if err != nil && !exitErrIsAcceptable(err) {
		log.Printf("npm run lint failed: %+v", err)
		log.Printf("terminal output: \n%s", string(output))
		j.Status = "failed"
		return
	}

	var results []LintResults
	err = json.Unmarshal(output, &results)
	if err != nil {
		log.Printf("failed to parse ESLint JSON output: %+v", err)
		j.Status = "failed"
		return
	}
	j.LintResults = results
	j.Status = "lint_complete"
}

func (j *Job) Build() {
	j.Status = "building"
	repository := "tmp/" + j.ID
	indexPath := filepath.Join(repository, "index.html")

	_, err := os.Stat(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Missing index.html at: %s", indexPath)
			addErr := addDummyIndex(j)
			if addErr != nil {
				log.Printf("Failed to add index.html: %v", addErr)
				j.Status = "failed"
				return
			}

			log.Printf("Dummy index.html injected.")
		} else {
			log.Printf("Unexpected error while checking for index.html: %v", err)
			j.Status = "failed"
			return
		}
	}

	cmdBuild := exec.Command("npm", "run", "build")
	cmdBuild.Dir = repository
	output, buildErr := cmdBuild.CombinedOutput()
	if buildErr != nil {
		log.Printf("npm run build: %+v", buildErr)
		log.Printf("terminal output: \n%s", string(output))
		j.Status = "failed"
		return
	}
	log.Printf("npm run build success!")
	j.Status = "build_complete"
}

func exitErrIsAcceptable(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 1
	}
	return false
}

func addDummyIndex(j *Job) error {
	newIndex := []byte(`<!DOCTYPE html><html><head><title>Temp</title></head><body><h1>Default Index</h1></body></html>`)
	path := filepath.Join("tmp", j.ID, "index.html")
	return os.WriteFile(path, newIndex, 0644) // 0644 is file mode
}
