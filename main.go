package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"os"
	"path"
	"path/filepath"
	"regexp"
)

func main() {

	// movabletype形式のインポートファイルから画像ファイルをダウンロードする。

	importfile := "/home/yamadatt/git/make_url_list/niko-nikkori-yokunaru.txt"
	regEx := `https://ameblo.jp/(.+?).html`
	re := regexp.MustCompile(regEx)

	contents := useIoutilReadFile(importfile)

	//正規表現に合致したものを配列にわたす
	match_words := re.FindAllStringSubmatch(contents, -1)

	for _, url := range match_words {

		fmt.Println(url[0]) // デバッグ用　ダウンロードするURL

	}

}

// 読み込んだファイルの中身をstringで返す
func useIoutilReadFile(fileName string) (str string) {
	// 関数説明：読み込んだファイルの中身をstringで返す
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

//URLをダウンロード先のパスを引数にして、ファイルをダウンロードする
func DownloadImage(imgurl string, downloadpath string) string {

	response, e := http.Get(imgurl)
	if e != nil {
		log.Fatal(e)
	}

	_, filename := path.Split(imgurl)

	defer response.Body.Close()
	//open a file for writing

	file, err := os.Create(downloadpath + filename)
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Printf("Success %S!\n", filename)
	return filename
}

// ファイルを上書きする
func overwrite_file(fileName string, contents string) {
	test_file_overwrite, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer test_file_overwrite.Close()
	fmt.Fprintln(test_file_overwrite, contents)
}

//引数で与えられたパスの配下について、ファイル名を再帰的にフルパスで取得する
func Dirwalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			// Recursively calls Dirwalk in the case of a directory
			p, err := Dirwalk(filepath.Join(dir, file.Name()))
			//fmt.Println(file.Name())
			if err != nil {
				return nil, fmt.Errorf("dirwalk %s: %w", filepath.Join(dir, file.Name()), err)
			}
			// Merge into the caller's "paths" variable.
			paths = append(paths, p...)
			continue
		}
		// Now that we've reached a leaf (file) in the directory tree, we'll add it to "paths" variable.
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, nil
}
