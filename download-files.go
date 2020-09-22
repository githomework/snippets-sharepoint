package main

import (
	"io/ioutil"
	"log"
	"os"
	"github.com/githomework/apps-app"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"
)

func getSharepointFiles() {
	emptyFolder(global.Folder + "/sharepoint")
	authCnfg := &strategy.AuthCnfg{
		SiteURL:  "https://xyz.sharepoint.com/sites/XXX/YYY",
		Username: global.options.SPUser,
		Password: global.options.SPPass,
	}
	// or using `private.json` creds source

	// authCnfg := &strategy.AuthCnfg{}
	//configPath := "./config/private.json"
	//if err := authCnfg.ReadConfig(configPath); err != nil {
	//		log.Fatalf("unable to get config: %v", err)
	//}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	// use client in raw requests or bind it with Fluent API ...

	sp := api.NewSP(client)

	
	config := &api.RequestConfig{}
	relRoot := "Shared Documents/Another Folder/Another Folder"
	web := sp.Conf(config).Web()
	folders, err := web.GetFolder(relRoot).Folders().Get()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range folders.Data() {
		d := v.Data()
		if !d.Exists || d.ItemCount < 1 || d.Name == "Ignore this folder" {
			continue
		}
		files, err := web.GetFolder(relRoot + "/" + d.Name).Files().Get()
		if err != nil {
			log.Println(err)
			return
		}

		for _, vv := range files.Data() {
			dd := vv.Data()
			if !dd.Exists {
				continue
			}
			data, err := web.GetFile(relRoot + "/" + d.Name + "/" + dd.Name).Download()
			if err != nil {
				log.Println(err)
				continue
			}
			err = ioutil.WriteFile(global.FolderAndSlash+"sharepoint/"+d.Name+"_"+dd.Name, data, 0644)
			if err != nil {
				log.Fatalf("unable to create a file: %v\n", err)
			}

		}
	}

}

func emptyFolder(folder string) {
	dirRead, _ := os.Open(folder)
	defer dirRead.Close()
	dirFiles, _ := dirRead.Readdir(0)

	// Loop over the directory's files.
	for index := range dirFiles {
		fileHere := dirFiles[index]

		// Get name of file and its full path.
		nameHere := fileHere.Name()
		fullPath := folder + "/" + nameHere

		// Remove the file.
		os.Remove(fullPath)

	}
}
