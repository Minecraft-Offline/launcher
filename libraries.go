package main

type Library struct {
	Downloads struct {
		Artifact    Artifact    `json:"artifact"`
		Classifiers Classifiers `json:"classifiers"`
	} `json:"downloads"`
	Extract struct {
		Exclude []string `json:"exclude"`
	} `json:"extract"`
	Name    string  `json:"name"`
	Natives Natives `json:"natives"`
	Rules   []Rule  `json:"rules"`

	Path string `json:"-"`
}

type Artifact struct {
	Path string `json:"path"`
	SHA1 string `json:"sha1"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

type Classifiers struct {
	JavaDoc        *Artifact `json:"javadoc,omitempty"`
	NativesLinux   *Artifact `json:"natives-linux,omitempty"`
	NativesMacOS   *Artifact `json:"natives-macos,omitempty"`
	NativesWindows *Artifact `json:"natives-windows,omitempty"`
	Sources        *Artifact `json:"sources,omitempty"`
}

type Natives struct {
	Linux   string `json:"linux,omitempty"`
	OSX     string `json:"osx,omitempty"`
	Windows string `json:"windows,omitempty"`
}
