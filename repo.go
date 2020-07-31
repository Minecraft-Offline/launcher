package main

type Repository struct {
	Name        string `json:"name"`
	LastUpdated string `json:"lastUpdated"`

	Assets    []*Asset   `json:"assets"`
	Libraries []*Library `json:"libraries"`
	//Mods          []*Mod          `json:"mods"`
	//ResourcePacks []*ResourcePack `json:"resourcePacks"`
	Versions []*Version `json:"versions"`
}

//https://files.minecraftforge.net/repo/
// - manifest.json
// - assets/
// - libraries/
// - mods/
// - resourcepacks/
// - versions/
