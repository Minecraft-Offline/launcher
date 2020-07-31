package main

import (
	json "github.com/JoshuaDoes/json" //JSON wrapper to handle data more conveniently
)

type VersionManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []*Version `json:"versions"`
}

func NewVersionManifest(url string) (manifest *VersionManifest, err error) {
	err = cache.ObjectDownloadInstall(url, "versions/versions.json", "")
	if err != nil {
		return nil, err
	}

	jsonData, err := cache.ObjectRead("versions/versions.json")
	if err != nil {
		return nil, err
	}

	manifest = &VersionManifest{}
	err = json.Unmarshal(jsonData, manifest)
	return manifest, err
}

func (manifest *VersionManifest) GetVersionsByType(versionType string) []*Version {
	versions := make([]*Version, 0)
	for _, version := range manifest.Versions {
		if version.Type == versionType {
			versions = append(versions, version)
		}
	}
	return versions
}
func (manifest *VersionManifest) GetReleases() []*Version {
	return manifest.GetVersionsByType("release")
}
func (manifest *VersionManifest) GetSnapshots() []*Version {
	return manifest.GetVersionsByType("snapshot")
}

type Version struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`

	ReleaseTime string `json:"releaseTime"`
	Time        string `json:"time"`

	Arguments       Arguments `json:"arguments"`
	LegacyArguments string    `json:"minecraftArguments"`
	MainClass       string    `json:"mainClass"`

	Path string `json:"-"`

	Assets     string     `json:"assets"`
	AssetIndex AssetIndex `json:"assetIndex"`

	Downloads Downloads `json:"downloads"`
	Libraries []Library `json:"libraries"`
}

func (version *Version) Load(path string) error {
	jsonData, err := fileRead(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, version)
	if err != nil {
		return err
	}

	//Parse argument values into readable values
	version.Arguments.ParseValues()

	//Load asset manifest
	if err := version.AssetIndex.Load(); err != nil {
		return err
	}

	return nil
}

type Downloads struct {
	Client         Download `json:"client"`
	ClientMappings Download `json:"client_mappings"`
	Server         Download `json:"server"`
	ServerMappings Download `json:"server_mappings"`
}
type Download struct {
	SHA1 string `json:"sha1"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}
