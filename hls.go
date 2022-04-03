package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type Downloader struct {
	queueCh        chan HLSInfo
	wg             sync.WaitGroup
	done           chan interface{}
	maxConcurrency int
	downloadPath   string
}

func NewDownloader() *Downloader {
	d := &Downloader{
		queueCh:        make(chan HLSInfo, 10),
		wg:             sync.WaitGroup{},
		done:           make(chan interface{}),
		maxConcurrency: 1,
		downloadPath:   filepath.Dir("/home/atori/Downloads/ffmpeg/"),
	}
	d.Start()
	return d
}

type DownloaderInfo struct {
	Waiting      int
	Concurrency  int
	DownloadPath string
}

func (d *Downloader) Start() {
	for i := 0; i < d.maxConcurrency; i++ {
		go d.download(d.queueCh)
	}
}

func (d *Downloader) Enqueue(info HLSInfo) {
	d.queueCh <- info
}

func (d *Downloader) Quit() {
	close(d.done)
	d.wg.Wait()
}

func (d *Downloader) Info() DownloaderInfo {
	return DownloaderInfo{
		Waiting:      len(d.queueCh),
		Concurrency:  d.maxConcurrency,
		DownloadPath: d.downloadPath,
	}
}

func (d *Downloader) download(ch <-chan HLSInfo) {
	d.wg.Add(1)
	defer d.wg.Done()
	for {
		select {
		case info, ok := <-ch:
			if !ok {
				return
			}

			if f, err := os.Stat(d.downloadPath); os.IsNotExist(err) || !f.IsDir() {
				os.Mkdir(d.downloadPath, 0777)
			}

			args := []string{
				"-i",
				info.URL,
				"-b:a",
				"128k",
				"-aac_coder",
				"twoloop",
				filepath.Join(d.downloadPath, info.Title),
			}

			cmd := exec.Command("ffmpeg", args...)
			cmd.Start()
			cmd.Wait()
		case <-d.done:
			return
		}
	}
}

type HLSInfo struct {
	URL   string
	Title string
}
