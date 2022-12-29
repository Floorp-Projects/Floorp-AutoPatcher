//////////////autopatcher-v1.0//////////////
/*Author: creeper-0910                    */
/*contributor: typeling1578,Comamoca      */
/*Thanks again!                           */
////////////////////////////////////////////
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/gjson"
)

var exename string

func main() {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/Floorp-Projects/Floorp/releases/latest", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyStr := string(body)
	key := gjson.Get(bodyStr, "assets.#.browser_download_url")
	key.ForEach(func(key, value gjson.Result) bool {
		if strings.Contains(value.String(), "floorp-win64.installer") {
			exename = strings.Split(value.String(), "/")[8]
			fmt.Printf("[found] %s\n", exename)
			resp, err := http.Get(value.String())
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			out, err := os.Create(exename)
			if err != nil {
				panic(err)
			}
			defer out.Close()
			bar := progressbar.NewOptions(int(resp.ContentLength),
				progressbar.OptionSetWriter(os.Stdout),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionSetDescription("[downloading] "+strings.Split(value.String(), "/")[8]),
				progressbar.OptionShowBytes(true))
			io.Copy(io.MultiWriter(out, bar), resp.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
		return true
	})
	resp, err = http.Get("https://www.7-zip.org/a/7zr.exe")
	if err != nil {
		panic(err)
	}
	out, err := os.Create("7zr.exe")
	if err != nil {
		panic(err)
	}
	bar := progressbar.NewOptions(int(resp.ContentLength),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[downloading] 7zr.exe"),
		progressbar.OptionShowBytes(true))
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	out.Close()

	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exe_dir := filepath.Dir(exe)

	if err = exec.Command(exe_dir+"/7zr.exe", "x", exename, "-x!setup.exe").Run(); err != nil {
		panic(err)
	}

	if err = exec.Command(exe_dir + "/patcher.exe").Run(); err != nil {
		panic(err)
	}
}
