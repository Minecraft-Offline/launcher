package main

import (
	"github.com/JoshuaDoes/json"

	//std necessities
	"fmt"
)

type Cache struct {
	CacheDir string         `json:"cacheDir"` //The root directory of the cache
	RootDir  string         `json:"rootDir"`  //The directory to install objects to
	Objects  []*CacheObject `json:"objects"`  //Every known object in the cache
}

//LoadCache loads a cache manifest into memory, or creates a new one if it doesn't exist
func LoadCache(cacheDir, rootDir string) (*Cache, error) {
	cache := &Cache{
		CacheDir: cacheDir,
		RootDir:  rootDir,
		Objects:  make([]*CacheObject, 0),
	}

	if fileExists(cacheDir + "/root.json") {
		cacheJSON, err := fileRead(cacheDir + "/root.json")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(cacheJSON, cache)
		if err != nil {
			return nil, err
		}

		validObjects := make([]*CacheObject, 0)
		for i := 0; i < len(cache.Objects); i++ {
			object := cache.Objects[i]

			sha1, err := hash_file_sha1(cacheDir + "/objects/" + object.SHA1)
			if err != nil {
				return nil, err
			}
			if sha1 != object.SHA1 {
				rm(cacheDir + "/objects/" + object.SHA1) //Delete the old object
				continue                                 //Don't keep deleted objects
			}

			if object.URL != "" && object.Path != "" && object.SHA1 != "" {
				validObjects = append(validObjects, object) //Only add objects if they're valid
			}
		}
		cache.Objects = validObjects

		//Update changed paths
		cache.CacheDir = cacheDir
		cache.RootDir = rootDir

		return cache, nil
	}

	return cache, nil
}

func (cache *Cache) Sync() error {
	defer log.Trace("Synced cache to " + cache.CacheDir + "/root.json")
	cacheJSON, err := json.Marshal(cache, false)
	if err != nil {
		return err
	}

	err = fileWrite(cacheJSON, cache.CacheDir+"/root.json")
	return err
}

func (cache *Cache) ObjectDownload(url, objectPath, sha1 string) (*CacheObject, error) {
	if sha1 == "" {
		err = download(url, cache.CacheDir+"/objects/tmp", "")

		sha1, err = hash_file_sha1(cache.CacheDir + "/objects/tmp")
		if err != nil {
			return nil, err
		}

		err = fileMove(cache.CacheDir+"/objects/tmp", cache.CacheDir+"/objects/"+sha1)
		if err != nil {
			return nil, err
		}
	} else {
		err = download(url, cache.CacheDir+"/objects/"+sha1, sha1)
		if err != nil {
			return nil, err
		}
	}

	object := &CacheObject{
		URL:  url,
		Path: objectPath,
		SHA1: sha1,
	}

	for i, cachedObject := range cache.Objects {
		if cachedObject.Path == object.Path {
			cache.Objects[i] = object
			return object, nil
		}
	}

	cache.Objects = append(cache.Objects, object)

	return object, nil
}

func (cache *Cache) ObjectGet(objectPath string) (*CacheObject, error) {
	var object *CacheObject
	for _, cachedObject := range cache.Objects {
		if cachedObject.Path == objectPath {
			object = cachedObject
			break
		}
	}
	if object == nil {
		return nil, fmt.Errorf("cache: object " + objectPath + " does not exist in cache manifest")
	}

	if !fileExists(cache.CacheDir + "/objects/" + object.SHA1) {
		return nil, fmt.Errorf("cache: object "+objectPath+" does not exist on filesystem (%v)", object)
		err = download(object.URL, cache.CacheDir+"/objects/"+object.SHA1, object.SHA1)
		if err != nil {
			return nil, err
		}
		return object, nil
	}

	return object, nil
}

func (cache *Cache) ObjectInstall(objectPath string) error {
	object, err := cache.ObjectGet(objectPath)
	if err != nil {
		return err
	}

	if fileExists(cache.RootDir + "/" + objectPath) {
		sha1, err := hash_file_sha1(cache.RootDir + "/" + objectPath)
		if err != nil {
			return err
		}

		if sha1 == object.SHA1 {
			return nil
		}
	}

	return fileCopy(cache.CacheDir+"/objects/"+object.SHA1, cache.RootDir+"/"+objectPath)
}

func (cache *Cache) ObjectRead(objectPath string) ([]byte, error) {
	object, err := cache.ObjectGet(objectPath)
	if err != nil {
		return nil, err
	}

	return fileRead(cache.CacheDir + "/objects/" + object.SHA1)
}

func (cache *Cache) ObjectCopy(src, dst string) error {
	object, err := cache.ObjectGet(src)
	if err != nil {
		return err
	}

	//Check if the object has previously been copied before
	_, err = cache.ObjectGet(dst)
	if err == nil {
		return nil
	}

	newObject := &CacheObject{
		URL:  object.URL,
		Path: dst,
		SHA1: object.SHA1,
	}
	cache.Objects = append(cache.Objects, newObject)

	return nil
}

//ObjectDownloadInstall conveniently downloads the object to the cache and installs it to the root directory
func (cache *Cache) ObjectDownloadInstall(url, objectPath, sha1 string) error {
	_, err = cache.ObjectDownload(url, objectPath, sha1)
	if err != nil {
		return err
	}

	err = cache.ObjectInstall(objectPath)
	return err
}

func (cache *Cache) Reset() {
	cache.Objects = make([]*CacheObject, 0)
}

type CacheObject struct {
	URL  string `json:"url"`  //The URL from which the cached object was downloaded
	Path string `json:"path"` //The relative path for the cached object to be installed
	SHA1 string `json:"sha1"` //The SHA1 hash of the cached object
}
