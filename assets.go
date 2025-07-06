package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	slic := []byte{}
	buf := bytes.NewBuffer(slic)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running cmd: %w", err)
	}
	probeDeets := FfprobeOut{}
	err = json.Unmarshal(buf.Bytes(), &probeDeets)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling ffprobe: %w", err)
	}
	if probeDeets.Streams[0].Width/probeDeets.Streams[0].Height == 16/9 {
		return "landscape", nil
	} else if probeDeets.Streams[0].Height/probeDeets.Streams[0].Width == 16/9 {
		return "portrait", nil
	}
	return "other", nil

}
