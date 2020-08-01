package main

import (
	"github.com/JoshuaDoes/logger" //Advanced logging
	flag "github.com/spf13/pflag"  //Unix-like command-line flags

	//std necessities
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

//Various command-line flags
var (
	gameDir         string //The folder to use instead of .minecraft
	cacheDir        string //The root cache folder, useful to change if you have faster I/O on another drive
	customArgs      string //The additional things to add to gameArgs
	email           string //Mojang account username
	password        string //Mojang account password
	server          string //The server to connect to on launch, empty to ignore
	targetVersion   string //The version to run if not the latest
	versionManifest string //A URL to the JSON manifest to use for fetching game versions
	verbosity       int    //0 = default (info, warning, error), 1 = 0 + debug, 2 = 1 + trace

	ihavebadinternet bool //Enable SHA1 hash checking of all files, in case of frequent internet issues preventing completed downloads
)

//Logging stuff
var (
	log       *logger.Logger //Holds the advanced logger's logging methods
	logPrefix string         = "LAUNCHER"
)

//Global things for anyone to use
var (
	err      error
	clientID string = "gomc"

	auth            *Auth
	cache           *Cache
	selectedVersion *Version
	versions        *VersionManifest
	webView         *Webview
)

func init() {
	switch runtime.GOOS {
	case "windows":
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		gameDir = strings.Replace(home, "\\", "/", -1) + "/AppData/Roaming/.minecraft"
	case "linux":
		home := os.Getenv("HOME")
		gameDir = home + "/.minecraft"
	case "darwin":
		home := os.Getenv("HOME")
		gameDir = home + "/Library/Application Support/minecraft"
	default:
		log.Fatal("Unsupported platform: ", runtime.GOOS)
	}

	oldGameDir := gameDir

	//Apply the command-line flags
	flag.StringVar(&gameDir, "gameDir", gameDir, "the directory to hold launcher and game data")
	flag.StringVar(&cacheDir, "cacheDir", gameDir+"/cache", "the directory to hold cache objects")
	flag.StringVar(&customArgs, "customArgs", "", "the additional arguments to pass to the game")
	flag.StringVar(&email, "email", "", "your mojang account email")
	flag.StringVar(&password, "password", "", "your mojang account password")
	flag.StringVar(&server, "server", "", "ip:port")
	flag.StringVar(&targetVersion, "version", "", "the target game version")
	flag.StringVar(&versionManifest, "versionManifest", "https://launchermeta.mojang.com/mc/game/version_manifest.json", "the version manifest to fetch game versions")
	flag.IntVar(&verbosity, "verbosity", 0, "sets the verbosity level; 0 = default, 1 = debug, 2 = trace")

	flag.BoolVar(&ihavebadinternet, "ihavebadinternet", false, "enables SHA1 hashing of all files")

	flag.Parse()

	//Update changed paths when gameDir changes
	if oldGameDir != gameDir {
		cacheDir = gameDir + "/cache"
	}

	//Create game directory
	mkdir(gameDir)

	//Initialize auth
	auth = &Auth{
		Email:    email,
		Password: password,
	}

	//Create the logger
	log = logger.NewLogger(logPrefix, verbosity)
}

func main() {
	log.Trace("--- main() ---")

	log.Info("Minecraft Offline © JoshuaDoes: 2020.")
	log.Debug("Platform: ", runtime.GOOS)

	go doWebsrv()

	doCleanup()
	doInitCache()
	//doLogin()
	//doFetchVersions()
	doWebview()
	/*
		doDownloadVersion()
		doDownloadAssets()
		doDownloadLibraries()
		doGameStart()
	*/

	log.Info("Good-bye!")
}

func doCleanup() {
	log.Debug("Cleaning out game directory...")
	rm(gameDir + "/assets")
	rm(gameDir + "/libraries")
	rm(gameDir + "/natives")
	rm(gameDir + "/versions")
}

func doInitCache() {
	log.Debug("Initializing cache...")
	cache, err = LoadCache(cacheDir, gameDir)
	if err != nil {
		log.Fatal(err)
	}
}

func doWebsrv() {
	log.Debug("Initializing webserver...")
	err = StartWebsrv()
	if err != nil {
		log.Fatal(err)
	}
}

func doWebview() {
	log.Debug("Initializing webview...")
	webView, err = NewWebview("http://localhost:25580/")
	if err != nil {
		log.Fatal(err)
	}

	if verbosity == 0 {
		hideConsole()
	}

	log.Debug("Running webview...")
	webView.Run()

	log.Debug("Destroying webview...")
	webView.Destroy()
}

func doLoadToken() {
	if fileExists(gameDir + "/token") {
		log.Debug("Loading token...")
		auth.LoadToken(gameDir + "/token")
	}
}

func doLogin() error {
	log.Info("Logging in...")
	err = auth.Login()
	if err != nil {
		return err
	}

	log.Debug("Saving authentication token...")
	err = auth.SaveToken(gameDir + "/token")
	if err != nil {
		return err
	}

	return nil
}

func doFetchVersions() {
	log.Debug("Fetching game versions...")
	versions, err = NewVersionManifest(versionManifest)
	if err != nil {
		log.Fatal(err)
	}
}

func doDownloadVersion() {
	if targetVersion != "" {
		for _, gameVersion := range versions.Versions {
			if gameVersion.ID == targetVersion {
				selectedVersion = gameVersion
				log.Info("Manually selected version ", selectedVersion.ID)
				break
			}
		}
		if selectedVersion == nil {
			log.Fatal("Invalid version: ", targetVersion)
		}
	} else {
		selectedVersion = versions.Versions[0]
		log.Info("Automatically selected version ", selectedVersion.ID)
	}

	log.Debug("Downloading and installing version manifest...")
	versionPath := fmt.Sprintf("versions/%s/%s.json", selectedVersion.ID, selectedVersion.ID)
	err = cache.ObjectDownloadInstall(selectedVersion.URL, versionPath, "")
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Loading version manifest ", versionPath)
	err = selectedVersion.Load(gameDir + "/" + versionPath)
	if err != nil {
		log.Fatal(err)
	}

	clientPath := fmt.Sprintf("versions/%s/%s.jar", selectedVersion.ID, selectedVersion.ID)

	//Download the game client
	log.Info("Downloading and installing client...")
	err = cache.ObjectDownloadInstall(selectedVersion.Downloads.Client.URL, clientPath, selectedVersion.Downloads.Client.SHA1)
	if err != nil {
		log.Fatal(err)
	}
	selectedVersion.Path = gameDir + "/" + clientPath

	cache.Sync()
}

func doDownloadAssets() {
	counter := 0
	length := strconv.Itoa(len(selectedVersion.AssetIndex.Assets))
	log.Info("Downloading " + length + " assets...")
	for path, asset := range selectedVersion.AssetIndex.Assets {
		objectPath := fmt.Sprintf("assets/objects/%s/%s", asset.Hash[:2], asset.Hash)
		legacyPath := fmt.Sprintf("assets/virtual/legacy/%s", path)

		counter++
		log.Debug("Downloading and installing asset " + path + " (" + strconv.Itoa(counter) + "/" + length + ")")
		err = cache.ObjectDownloadInstall(asset.GetURL(), objectPath, asset.Hash)
		if err != nil {
			log.Fatal(err)
		}

		//log.Debug("Duplicating cached asset ", path, " to ", legacyPath)
		err = cache.ObjectCopy(objectPath, legacyPath)
		if err != nil {
			log.Fatal(err)
		}

		//log.Debug("Installing legacy asset ", legacyPath)
		err = cache.ObjectInstall(legacyPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	cache.Sync()
}

func doDownloadLibraries() {
	log.Info("Downloading libraries...")
	for i := 0; i < len(selectedVersion.Libraries); i++ {
		if len(selectedVersion.Libraries[i].Rules) > 0 {
			valid := false
			for _, rule := range selectedVersion.Libraries[i].Rules {
				valid = rule.Valid()

				if !valid {
					break
				}
			}
			if !valid {
				log.Debug("Skipping library ", selectedVersion.Libraries[i].Name)
				continue
			}
		}

		if selectedVersion.Libraries[i].Downloads.Artifact.URL != "" {
			libraryPath := fmt.Sprintf("libraries/%s", selectedVersion.Libraries[i].Downloads.Artifact.Path)

			log.Debug("Downloading and installing library ", selectedVersion.Libraries[i].Name)
			err = cache.ObjectDownloadInstall(selectedVersion.Libraries[i].Downloads.Artifact.URL, libraryPath, selectedVersion.Libraries[i].Downloads.Artifact.SHA1)
			if err != nil {
				log.Fatal(err)
			}
			selectedVersion.Libraries[i].Path = gameDir + "/" + libraryPath
		}

		switch runtime.GOOS {
		case "windows":
			if native := selectedVersion.Libraries[i].Downloads.Classifiers.NativesWindows; native != nil {
				nativePath := fmt.Sprintf("libraries/%s", selectedVersion.Libraries[i].Downloads.Classifiers.NativesWindows.Path)

				log.Debug("Downloading and installing Windows native ", selectedVersion.Libraries[i].Name)
				err = cache.ObjectDownloadInstall(selectedVersion.Libraries[i].Downloads.Classifiers.NativesWindows.URL, nativePath, selectedVersion.Libraries[i].Downloads.Classifiers.NativesWindows.SHA1)
				if err != nil {
					log.Fatal(err)
				}

				log.Debug("Extracting Windows native ", selectedVersion.Libraries[i].Name)
				err = extract(gameDir+"/"+nativePath, gameDir+"/natives/windows", selectedVersion.Libraries[i].Extract.Exclude...)
				if err != nil {
					log.Fatal(err)
				}
			}
		case "linux":
			if native := selectedVersion.Libraries[i].Downloads.Classifiers.NativesLinux; native != nil {
				nativePath := fmt.Sprintf("libraries/%s", selectedVersion.Libraries[i].Downloads.Classifiers.NativesLinux.Path)

				log.Debug("Downloading and installing Linux native ", selectedVersion.Libraries[i].Name)
				err = cache.ObjectDownloadInstall(selectedVersion.Libraries[i].Downloads.Classifiers.NativesLinux.URL, nativePath, selectedVersion.Libraries[i].Downloads.Classifiers.NativesLinux.SHA1)
				if err != nil {
					log.Fatal(err)
				}

				log.Debug("Extracting Linux native ", selectedVersion.Libraries[i].Name)
				err = extract(gameDir+"/"+nativePath, gameDir+"/natives/linux", selectedVersion.Libraries[i].Extract.Exclude...)
				if err != nil {
					log.Fatal(err)
				}
			}
		case "darwin":
			if native := selectedVersion.Libraries[i].Downloads.Classifiers.NativesMacOS; native != nil {
				nativePath := fmt.Sprintf("libraries/%s", selectedVersion.Libraries[i].Downloads.Classifiers.NativesMacOS.Path)

				log.Debug("Downloading and installing MacOS native ", selectedVersion.Libraries[i].Name)
				err = cache.ObjectDownloadInstall(selectedVersion.Libraries[i].Downloads.Classifiers.NativesMacOS.URL, nativePath, selectedVersion.Libraries[i].Downloads.Classifiers.NativesMacOS.SHA1)
				if err != nil {
					log.Fatal(err)
				}

				log.Debug("Extracting MacOS native ", selectedVersion.Libraries[i].Name)
				err = extract(gameDir+"/"+nativePath, gameDir+"/natives/osx", selectedVersion.Libraries[i].Extract.Exclude...)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	cache.Sync()
}

func doGameStart() {
	separator := ":"
	if runtime.GOOS == "windows" {
		separator = ";"
	}
	libraries := ""
	for _, library := range selectedVersion.Libraries {
		if library.Path == "" {
			continue
		}
		libraries += library.Path + separator
	}
	libraries += selectedVersion.Path

	nativesDir := ""
	librariesDir := ""
	switch runtime.GOOS {
	case "windows":
		nativesDir = gameDir + "\\natives\\windows"
		librariesDir = gameDir + "\\libraries"
	case "linux":
		nativesDir = gameDir + "/natives/linux"
		librariesDir = gameDir + "/libraries"
	case "darwin":
		nativesDir = gameDir + "/natives/osx"
		librariesDir = gameDir + "/libraries"
	}

	gameArgs := ""
	jvmArgs := fmt.Sprintf("-Djava.library.path=%s\x00-Dminecraft.launcher.brand=gomc\x00-Dminecraft.launcher.version=0.0.0", nativesDir)

	if selectedVersion.LegacyArguments != "" {
		gameArgs = strings.Replace(selectedVersion.LegacyArguments, " ", "\x00", -1)
	} else {
		jvmArgs = selectedVersion.Arguments.JVMArgs()
		jvmArgs = strings.Replace(jvmArgs, "${natives_directory}", nativesDir, -1)
		jvmArgs = strings.Replace(jvmArgs, "${launcher_name}", "gomc", -1)
		jvmArgs = strings.Replace(jvmArgs, "${launcher_version}", "0.0.0", -1)
		jvmArgs = strings.Replace(jvmArgs, "${classpath}", gameDir+"/libraries", -1)

		gameArgs = selectedVersion.Arguments.GameArgs()
	}

	gameArgs = strings.Replace(gameArgs, "${assets_root}", gameDir+"/assets", -1)
	gameArgs = strings.Replace(gameArgs, "${assets_index_name}", selectedVersion.Assets, -1)
	gameArgs = strings.Replace(gameArgs, "${auth_access_token}", auth.AccessToken, -1)
	gameArgs = strings.Replace(gameArgs, "${auth_session}", auth.DecodeToken.YGGT, -1)
	gameArgs = strings.Replace(gameArgs, "${auth_player_name}", auth.Username, -1)
	gameArgs = strings.Replace(gameArgs, "${auth_uuid}", auth.ID, -1)
	gameArgs = strings.Replace(gameArgs, "${game_assets}", gameDir+"/assets/virtual/legacy", -1)
	gameArgs = strings.Replace(gameArgs, "${game_directory}", gameDir, -1)
	gameArgs = strings.Replace(gameArgs, "${user_properties}", "{}", -1)
	gameArgs = strings.Replace(gameArgs, "${user_type}", "mojang", -1)
	gameArgs = strings.Replace(gameArgs, "${version_name}", selectedVersion.ID, -1)
	gameArgs = strings.Replace(gameArgs, "${version_type}", selectedVersion.Type, -1)

	launchArgs := jvmArgs + "\x00-cp\x00" + libraries + "\x00" + selectedVersion.MainClass + "\x00" + gameArgs
	if customArgs != "" {
		launchArgs += "\x00" + strings.Replace(customArgs, " ", "\x00", -1)
	}

	if server != "" {
		host := strings.Split(server, ":")
		switch len(host) {
		case 0:
			break
		case 1:
			launchArgs += "\x00--server\x00" + host[0] + "\x00--port\x0025565"
		case 2:
			if _, err := strconv.Atoi(host[1]); err != nil {
				log.Fatal(err)
			}

			launchArgs += "\x00--server\x00" + host[0] + "\x00--port\x00" + host[1]
		default:
			log.Fatal("Server string invalid")
		}
	}

	execArgs := strings.Split(launchArgs, "\x00")
	if verbosity > 0 {
		for _, execArg := range execArgs {
			log.Trace("Argument: ", execArg)
		}
	}

	gameProcess := exec.Command("java", strings.Split(launchArgs, "\x00")...)
	gameProcess.Dir = gameDir
	gameProcess.Stdout = os.Stdout
	gameProcess.Stderr = os.Stderr

	log.Info("Starting Minecraft ", selectedVersion.ID, "...")
	err = gameProcess.Run()
	if err != nil {
		log.Error(err)
	}
}
