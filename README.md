# Minecraft Offline

### A Minecraft launcher written in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/Minecraft-Offline/launcher)](https://goreportcard.com/report/github.com/Minecraft-Offline/launcher)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://paypal.me/JoshuaDoes)

---

## What is Minecraft Offline?

Despite the name of the launcher, which is subject to change in the future, it doesn't actually let you play a cracked version of Minecraft. Rather, Minecraft Offline is being written to handle the issues we've all faced before: mods, with multiple installed versions of the game. Perhaps you're using MultiMC to have a different dotminecraft for each profile, or maybe you're just renaming your "mods" folder manually each time you change your version profile in the official Minecraft launcher. Maybe you're constantly trying to download different versions of a mod for each Minecraft version from weird websites with sketchy ads. Minecraft Offline solves this issue by managing it all for you.

Fleshing out the idea for this launcher has taken some time, and most of anything written in this document is subject to change depending on the technical aspects of implementing these ideas. I encourage you to critique and help me implement these ideas, though.

The idea as it stands today: Repositories. They're the key to everything. It's what makes Linux package management work, so why can't we try something similar? Everything that is automatically loaded from a JSON manifest (such as versions, assets, libraries, etc) as well as my additions (such as mods, resource packs, profiles, and whatever else I come up with) can be specified inside of a repository's root manifest. The idea is to be able to add a repository, for example `https://files.minecraftforge.net/repo`, sync your cached list of packages with its specified packages in `https://files.minecraftforge.net/repo/root.json` (filename subject to change), and then go and install the latest Forge profile that's automatically configured to download their custom libraries and additional files, all straight from the repository. This would negate the requirement to go to the Forge website yourself and download the installer, which just installs the Forge profile. This can be taken a step further with the game's server direct connect arguments (for example, `--server mc.hypixel.net --port 25565`), where a server could host its own repository to pull all of the required mods and other data before launching Minecraft straight into the server. Users would also easily be able to host their own repositories, where mod authors, resource pack authors, and modded server hosts would be especially encouraged to host their own official repositories.

As of writing, the current state of the launcher is working but has an incomplete interface. These additional features on top of the vanilla launching logic flow won't have any true progress until the interface matures.

## Building from source

##### Windows: You *must* install TDM-GCC-64 or a working alternative in order to compile Go's side of webview, the driver behind the launcher interface. Further, if you would like to compile 32-bit webview DLLs or compile updated 64-bit webview DLLs, you need to install Visual Studio and run `$GOPATH\src\github.com\webview\webview\script\build.bat` to compile them. You'll find the resulting DLLs under `$GOPATH\src\github.com\webview\webview\dll\` inside of your architecture's subdirectory, which need to be placed in the same directory as Minecraft Offline.

Download and build Minecraft Offline:
```
git clone https://github.com/Minecraft-Offline/launcher.git Minecraft-Offline
cd Minecraft-Offline
go get
go build
```

## Running the build

Run Minecraft Offline:
```
# Linux and MacOS
./launcher --gameDir $PWD/.minecraft --email "Mojang account email" --password "Mojang account password" --verbosity 2

# Windows
# Keeping verbosity either 1 (debug) or 2 (trace) is important here because with a value of 0 (info, warn, error, fatal) it will automatically hide the console window (even for PowerShell and Windows Terminal, not just Command Prompt) when the webview interface opens, something we do *not* want before the actual webview interface is complete.
launcher.exe --gameDir $PWD\.minecraft --email "Mojang account email" --password "Mojang account password" --verbosity 2
```

## License
The source code for Minecraft Offline is on a to-be-determined license. Everything in this repository, from the first commit onward, is to be licensed according to the new license once it is decided upon, regardless of cloned copy or fork from this original host. You reserve the right to download, compile, and run this code at your own risk, with no warranty, neither express nor implied. You do not have the right to redistribute compiled code using the source code in this repository without first making your changes public. See LICENSE for a clear copy of this license placeholder.

## Donations
If you like what you're seeing and want me to contribute more of my time to this project, you can donate to show your support! It's okay if you don't though, Minecraft Offline is free and open-source after all. You can find the PayPal donation button at the top of this document if you're interested.