package main

import (
	json "github.com/JoshuaDoes/json" //JSON wrapper to handle data more conveniently

	//std necessities
	stdjson "encoding/json"
)

type AssetIndex struct {
	ID        string `json:"id"`
	SHA1      string `json:"sha1"`
	Size      int    `json:"size"`
	TotalSize int    `json:"totalSize"`
	URL       string `json:"url"`

	Objects stdjson.RawMessage `json:"objects"`
	Assets  map[string]*Asset  `json:"-"`
}

func (manifest *AssetIndex) Load() error {
	if err := cache.ObjectDownloadInstall(manifest.URL, "assets/indexes/"+manifest.ID+".json", manifest.SHA1); err != nil {
		return err
	}

	cache.Sync()

	jsonData, err := cache.ObjectRead("assets/indexes/" + manifest.ID + ".json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, manifest); err != nil {
		return err
	}

	if err := json.Unmarshal(manifest.Objects, &manifest.Assets); err != nil {
		return err
	}

	return nil
}

type Asset struct {
	Path string `json:"-"`
	Hash string `json:"hash"`
	Size int    `json:"size"`
}

func (asset *Asset) GetURL() string {
	return "http://resources.download.minecraft.net/" + asset.Hash[:2] + "/" + asset.Hash
}
