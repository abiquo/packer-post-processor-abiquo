package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/resty.v0"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	vmwcommon "github.com/hashicorp/packer/builder/vmware/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type DTO struct {
	Links []Link `json:"links,omitempty"`
}

type Link struct {
	Type  string `json:"type,omitempty"`
	Href  string `json:"href,omitempty"`
	Title string `json:"title,omitempty"`
	Rel   string `json:"rel,omitempty"`
}

type Enterprise struct {
	DTO
	Name string
}

type Repo struct {
	DTO
	Name               string
	RepositoryLocation string
}

type User struct {
	DTO
	Name  string
	Email string
}

type DiskDef struct {
	DTO
	Bootable                  bool   `json:"bootable,omitempty"`
	Sequence                  int    `json:"sequence,omitempty"`
	RequiredHDInMB            int    `json:"requiredHDInMB,omitempty"`
	DiskFileFormat            string `json:"diskFileFormat,omitempty"`
	VirtualMachineTemplateUrl string `json:"virtualMachineTemplateUrl,omitempty"`
	DiskUrl                   string `json:"diskUrl,omitempty"`
	CurrentPath               string `json:"currentPath,omitempty"`
}

type VMTDef struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	CategoryName       string `json:"categoryName,omitempty"`
	DiskFileFormat     string `json:"diskFileFormat,omitempty"`
	RequiredCpu        string `json:"requiredCpu,omitempty"`
	RequiredHDInMB     string `json:"requiredHDInMB,omitempty"`
	RequiredRamInMB    string `json:"requiredRamInMB,omitempty"`
	LoginUser          string `json:"loginUser,omitempty"`
	LoginPassword      string `json:"loginPassword,omitempty"`
	OsType             string `json:"osType,omitempty"`
	OsVersion          string `json:"osVersion,omitempty"`
	EthernetDriverType string `json:"ethernetDriverType,omitempty"`
}

type VirtualMachineTemplate struct {
	DTO
	Name                             string `json:"name,omitempty"`
	ChefEnabled                      bool   `json:"chefEnabled,omitempty"`
	CpuRequired                      int    `json:"cpuRequired,omitempty"`
	CreationDate                     string `json:"creationDate,omitempty"`
	CreationUser                     string `json:"creationUser,omitempty"`
	Description                      string `json:"description,omitempty"`
	EthernetDriverType               string `json:"ethernetDriverType,omitempty"`
	IconUrl                          string `json:"iconUrl,omitempty"`
	Id                               int    `json:"id,omitempty"`
	LoginPassword                    string `json:"loginPassword,omitempty"`
	LoginUser                        string `json:"loginUser,omitempty"`
	OsType                           string `json:"osType,omitempty"`
	OsVersion                        string `json:"osVersion,omitempty"`
	RamRequired                      int    `json:"ramRequired,omitempty"`
	State                            string `json:"state,omitempty"`
	EnableCpuHotAdd                  bool   `json:"enableCpuHotAdd,omitempty"`
	EnableRamHotAdd                  bool   `json:"enableRamHotAdd,omitempty"`
	EnableDisksHotReconfigure        bool   `json:"enableDisksHotReconfigure,omitempty"`
	EnableNicsHotReconfigure         bool   `json:"enableNicsHotReconfigure,omitempty"`
	EnableRemoteAccessHotReconfigure bool   `json:"enableRemoteAccessHotReconfigure,omitempty"`
}

type TemplateDisk struct {
	DTO
	Label              string `json:"label,omitempty"`
	Sequence           int    `json:"sequence,omitempty"`
	Path               string `json:"path,omitempty"`
	DiskFormatType     string `json:"diskFormatType,omitempty"`
	DiskFileSize       int    `json:"diskFileSize,omitempty"`
	HdRequired         int    `json:"hdRequired,omitempty"`
	State              string `json:"state,omitempty"`
	DiskControllerType string `json:"diskControllerType,omitempty"`
	DiskController     string `json:"diskController,omitempty"`
	CreationDate       string `json:"creationDate,omitempty"`
	Bootable           bool   `json:"bootable,omitempty"`
}

func (dto *DTO) GetLink(rel string) Link {
	for _, link := range dto.Links {
		if link.Rel == rel {
			return link
		}
	}
	return Link{}
}

func (d *DiskDef) ToJson() ([]byte, error) {
	str, err := json.Marshal(d)
	return str, err
}

func (template *VirtualMachineTemplate) ToJson() ([]byte, error) {
	str, err := json.Marshal(template)
	return str, err
}

func (def *VMTDef) ToJson() ([]byte, error) {
	str, err := json.Marshal(def)
	return str, err
}

// https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}

func (def *VMTDef) Upload(config Config, repo Repo, artifact packer.Artifact) (VirtualMachineTemplate, error) {
	var newTemplate VirtualMachineTemplate
	log.Printf("Template def is : %v", def)
	definition_json, err := def.ToJson()
	if err != nil {
		return VirtualMachineTemplate{}, err
	}
	log.Printf("Template def json is : %s", definition_json)

	file, err := getFilesFromArtifact(config, artifact, "vmdk")
	if err != nil {
		return VirtualMachineTemplate{}, err
	}

	post_url := ""
	for _, link := range repo.Links {
		if link.Rel == "applianceManagerRepositoryUri" {
			post_url = link.Href + "/templates"
		}
	}
	if post_url == "" {
		return VirtualMachineTemplate{}, errors.New("Could not find AM repo URI.")
	}

	params := map[string]string{
		"diskInfo": string(definition_json),
	}
	request, err := newfileUploadRequest(post_url, params, "diskFile", file)
	request.SetBasicAuth(config.ApiUsername, config.ApiPassword)
	if err != nil {
		log.Printf("ERROR ON UPLOAD!", err)
		return newTemplate, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(7200 * time.Second),
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("ERROR uploading file!", err)
		return newTemplate, err
	}

	location, err := resp.Location()
	if err != nil {
		log.Printf("Upload response did not had a location header!")
		log.Printf("Response code was : %d", resp.StatusCode)
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		log.Printf("Body was : %s", bodyString)
		return newTemplate, err
	}

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	rclient := resty.R().SetBasicAuth(config.ApiUsername, config.ApiPassword)
	respt, err := rclient.
		SetHeader("Accept", "application/vnd.abiquo.virtualmachinetemplate+json").
		Get(location.String())
	if err != nil {
		return newTemplate, err
	}

	json.Unmarshal(respt.Body(), &newTemplate)
	return newTemplate, err
}

func (templatedef *VMTDef) GetDiskFileInfo(config Config, artifact packer.Artifact) (string, int, error) {
	var abqFormat string

	file, err := getFilesFromArtifact(config, artifact, "vmdk")
	if err != nil {
		return abqFormat, 0, err
	}

	cmdname := "qemu-img"
	cmdargs := []string{"info", file}
	cmd := exec.Command(cmdname, cmdargs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating StdoutPipe for QEMU img command")
		return abqFormat, 0, err
	}

	if err := cmd.Start(); err != nil {
		return abqFormat, 0, err
	}

	// read command's stdout line by line
	in := bufio.NewScanner(cmdReader)

	var format string
	var sub string
	fileSize := 0

	log.Printf("Qemu-img output:")
	for in.Scan() {
		log.Printf(in.Text())
		if strings.HasPrefix(in.Text(), "file format:") {
			format = strings.TrimSpace(strings.Split(in.Text(), ":")[1])
		}
		if strings.Contains(in.Text(), "create type:") {
			sub = strings.TrimSpace(strings.Split(in.Text(), ":")[1])
		}
		if strings.HasPrefix(in.Text(), "virtual size:") {
			parsedSize := strings.TrimSpace(strings.Split(in.Text(), "(")[1])
			parsedSizeBytes := strings.TrimSpace(strings.Split(parsedSize, " ")[0])
			sizeInBytes, err := strconv.Atoi(parsedSizeBytes)
			if err != nil {
				return abqFormat, 0, err
			}
			fileSize = sizeInBytes / 1024 / 1024
		}
	}
	if err := in.Err(); err != nil {
		return abqFormat, 0, err
	}
	log.Printf("format: %s", format)
	log.Printf("type: %s", sub)
	log.Printf("size: %d", fileSize)

	switch format {
	case "vmdk":
		switch sub {
		case "monolithicSparse":
			abqFormat = "VMDK_SPARSE"
		case "streamOptimized":
			abqFormat = "VMDK_STREAM_OPTIMIZED"
		case "monolithicFlat":
			abqFormat = "VMDK_FLAT"
		}
	case "qcow2":
		abqFormat = "QCOW2_SPARSE"
	case "raw":
		abqFormat = "RAW"
	case "vdi":
		abqFormat = "VDI_SPARSE"
	}
	log.Printf("Abiquo format: %s", abqFormat)
	return abqFormat, fileSize, nil
}

func (t *VirtualMachineTemplate) ReplacePrimaryDisk(config Config, diskdef DiskDef, artifact packer.Artifact) (VirtualMachineTemplate, error) {
	var newTemplate VirtualMachineTemplate

	log.Printf("Disk def is : %v", diskdef)
	diskdef_json, err := diskdef.ToJson()
	if err != nil {
		return newTemplate, err
	}
	log.Printf("Disk def json is : %s", string(diskdef_json))

	file, err := getFilesFromArtifact(config, artifact, "vmdk")
	if err != nil {
		return newTemplate, err
	}

	templateUpdateUrl := t.GetLink("templatePath").Href

	params := map[string]string{
		"diskInfo": string(diskdef_json),
	}
	request, err := newfileUploadRequest(templateUpdateUrl, params, "diskFile", file)
	request.SetBasicAuth(config.ApiUsername, config.ApiPassword)
	if err != nil {
		log.Printf("ERROR ON UPLOAD!", err)
		return newTemplate, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(7200 * time.Second),
	}
	_, err = client.Do(request)
	if err != nil {
		log.Printf("ERROR uploading file!", err)
		return newTemplate, err
	}

	rclient := resty.R().SetBasicAuth(config.ApiUsername, config.ApiPassword)
	respt, err := rclient.
		SetHeader("Accept", t.GetLink("edit").Type).
		Get(t.GetLink("edit").Href)
	if err != nil {
		return newTemplate, err
	}

	json.Unmarshal(respt.Body(), &newTemplate)
	return newTemplate, err
}

func (t *VirtualMachineTemplate) Update(config Config) error {
	template_json, err := t.ToJson()
	if err != nil {
		return err
	}
	log.Printf("Template json is : %s", string(template_json))

	templateUrl := t.GetLink("edit").Href
	templateType := t.GetLink("edit").Type

	if os.Getenv("RESTYDEBUG") != "" {
		// Enable debug mode
		resty.SetDebug(true)

		// Using you custom log writer
		logFile, _ := os.OpenFile("/tmp/go-resty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		resty.SetLogger(logFile)
	}

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := resty.R().
		SetBasicAuth(config.ApiUsername, config.ApiPassword).
		SetHeader("Accept", templateType).
		SetHeader("Content-Type", templateType).
		SetBody(string(template_json)).
		Put(templateUrl)
	if err != nil {
		return err
	}

	json.Unmarshal(resp.Body(), &t)
	return nil
}

type AbstractCollection struct {
	Links     []Link
	TotalSize int
}

type RepoCollection struct {
	AbstractCollection
	Collection []Repo
}

type EntCollection struct {
	AbstractCollection
	Collection []Enterprise
}

type UserCollection struct {
	AbstractCollection
	Collection []User
}

type TemplateCollection struct {
	AbstractCollection
	Collection []VirtualMachineTemplate
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// Fields from config file
	ApiUrl            string `mapstructure:"api_url"`
	ApiUsername       string `mapstructure:"api_username"`
	ApiPassword       string `mapstructure:"api_password"`
	Datacenter        string `mapstructure:"datacenter"`
	KeepInputArtifact bool   `mapstructure:"keep_input_artifact"`

	Name               string `mapstructure:"template_name"`
	Description        string `mapstructure:"description"`
	CategoryName       string `mapstructure:"category"`
	DiskFileFormat     string `mapstructure:"disk_format"`
	RequiredCpu        string `mapstructure:"cpu"`
	RequiredHDInMB     string `mapstructure:"hd_mb"`
	RequiredRamInMB    string `mapstructure:"ram_mb"`
	LoginUser          string `mapstructure:"login_user"`
	LoginPassword      string `mapstructure:"login_password"`
	EthernetDriverType string `mapstructure:"eth_driver"`

	ChefEnabled                      bool   `mapstructure:"chef_enabled"`
	IconUrl                          string `mapstructure:"icon_url"`
	EnableCpuHotAdd                  bool   `mapstructure:"cpu_hotadd"`
	EnableRamHotAdd                  bool   `mapstructure:"ram_hotadd"`
	EnableDisksHotReconfigure        bool   `mapstructure:"disk_hotadd"`
	EnableNicsHotReconfigure         bool   `mapstructure:"nic_hotadd"`
	EnableRemoteAccessHotReconfigure bool   `mapstructure:"vnc_hotadd"`

	SshUser     string
	SshPass     string
	GuessOsType string

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
	}, raws...)
	if err != nil {
		return err
	}

	errs := new(packer.MultiError)

	if p.config.ApiUrl == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Abiquo API URL is missing!"))
	}

	if p.config.ApiUsername == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Abiquo API username is missing!"))
	}

	if p.config.ApiPassword == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Abiquo API password is missing!"))
	}

	if p.config.CategoryName == "" {
		p.config.CategoryName = "OS"
	}

	if p.config.RequiredCpu == "" {
		p.config.RequiredCpu = "1"
	}

	if p.config.RequiredRamInMB == "" {
		p.config.RequiredRamInMB = "1024"
	}

	if p.config.EthernetDriverType == "" {
		p.config.EthernetDriverType = "E1000"
	}

	if p.config.LoginUser == "" {
		p.config.LoginUser = p.config.SshUser
	}

	if p.config.LoginPassword == "" {
		p.config.LoginPassword = p.config.SshPass
	}

	if len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	// These are extra variables that will be made available for interpolation.
	p.config.ctx.Data = map[string]string{
		"BuildName":   p.config.PackerBuildName,
		"BuilderType": p.config.PackerBuilderType,
	}

	// If no template name, get VM name
	ui.Say("Getting config items...")
	var vmName string
	var vmxOsType string
	if p.config.Name == "" {
		vmxFile, err := getFilesFromArtifact(p.config, artifact, "vmx")
		if err != nil {
			ui.Message(fmt.Sprintf("err: %s", err))
		}
		vmxData, err := vmwcommon.ReadVMX(vmxFile)
		if err != nil {
			ui.Message(fmt.Sprintf("err: %s", err))
		}
		vmName = vmxData["displayname"]
		vmxOsType = vmxData["guestos"]
	}
	if vmName != "" {
		p.config.Name = vmName
		if p.config.Description == "" {
			p.config.Description = vmName
		}
	}
	if p.config.GuessOsType == "" {
		p.config.GuessOsType = vmxOsType
	}
	log.Printf("Config is : %v", p.config)

	ui.Say("Looking up the repo URL for datacenter '" + p.config.Datacenter + "'")
	repo, err := p.config.FindRepoUrl()
	if err != nil {
		return artifact, p.config.KeepInputArtifact, err
	}

	ui.Say("Checking if a template named '" + p.config.Name + "' already exists...")
	var template VirtualMachineTemplate
	exists, template, err := p.config.CheckTemplateExists(repo)
	if exists {
		ui.Say("Template already exists. Replacing primary disk...")
		diskdef := p.config.BuildDiskDef(template)
		templatedef := p.config.BuildTemplateDef()
		abqdiskformat, size, err := templatedef.GetDiskFileInfo(p.config, artifact)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}
		diskdef.RequiredHDInMB = size
		diskdef.DiskFileFormat = abqdiskformat
		template, err := template.ReplacePrimaryDisk(p.config, diskdef, artifact)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}
		ui.Say("Upload complete. The URL of the updated template is " + template.GetLink("edit").Href)
	} else {
		ui.Say("Uploading template...")
		templatedef := p.config.BuildTemplateDef()
		abqdiskformat, size, err := templatedef.GetDiskFileInfo(p.config, artifact)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}
		templatedef.DiskFileFormat = abqdiskformat
		templatedef.RequiredHDInMB = strconv.Itoa(size)
		templateUploaded, err := templatedef.Upload(p.config, repo, artifact)
		log.Printf("Uploaded template : %v", templateUploaded)
		exists, template, err = p.config.CheckTemplateExists(repo)
		if exists {
			log.Printf("Found new template '%s'", template.Name)
		}
		ui.Say("Upload complete. The URL of the new template is " + template.GetLink("edit").Href)
	}

	ui.Say("Updating template with extra attributes...")

	template.ChefEnabled = p.config.ChefEnabled
	template.Description = p.config.Description
	template.IconUrl = p.config.IconUrl
	template.EnableCpuHotAdd = p.config.EnableCpuHotAdd
	template.EnableRamHotAdd = p.config.EnableRamHotAdd
	template.EnableDisksHotReconfigure = p.config.EnableDisksHotReconfigure
	template.EnableNicsHotReconfigure = p.config.EnableNicsHotReconfigure
	template.EnableRemoteAccessHotReconfigure = p.config.EnableRemoteAccessHotReconfigure
	template.Update(p.config)

	newArtifact := &Artifact{Url: template.GetLink("edit").Href}
	return newArtifact, false, nil
}

func (config *Config) FindRepoUrl() (Repo, error) {
	log.Printf("Trying to find the repo URL for Datacenter '%s'\n", config.Datacenter)
	login_url := config.ApiUrl + "/login"
	log.Printf("Login URL: %s\n", login_url)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	rclient := resty.R().SetBasicAuth(config.ApiUsername, config.ApiPassword)
	resp, err := rclient.
		SetHeader("Accept", "application/vnd.abiquo.user+json").
		Get(login_url)
	log.Printf("Request done.\n")
	if err != nil {
		log.Printf("Error on request: %s\n", err)
		return Repo{}, err
	}

	var user User
	json.Unmarshal(resp.Body(), &user)
	for _, link := range user.Links {
		if link.Rel == "enterprise" {
			resp, err = rclient.SetHeader("Accept", link.Type).Get(link.Href)
			if err != nil {
				return Repo{}, err
			}

			var ent Enterprise
			json.Unmarshal(resp.Body(), &ent)

			for _, link = range ent.Links {
				link_rel := link.Rel
				if link_rel == "datacenterrepositories" {
					resp, err = rclient.SetHeader("Accept", link.Type).Get(link.Href)
					if err != nil {
						return Repo{}, err
					}

					var repos RepoCollection
					json.Unmarshal([]byte(resp.Body()), &repos)

					for _, repo := range repos.Collection {
						for _, repolink := range repo.Links {
							if repolink.Rel == "datacenter" && repolink.Title == config.Datacenter {
								return repo, err
								log.Printf("Found URL '%s' (%s) for DC '%s'", repo.Name, repo.RepositoryLocation, config.Datacenter)
							}
						}
					}
				}
			}
		}
	}
	log.Printf("Could not find the URL for DC '%s'!", config.Datacenter)
	return Repo{}, err
}

func (config *Config) CheckTemplateExists(repo Repo) (bool, VirtualMachineTemplate, error) {
	var template VirtualMachineTemplate
	log.Printf("Checking if a template with name '%s' already exists.", config.Name)
	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	rclient := resty.R().SetBasicAuth(config.ApiUsername, config.ApiPassword)

	for _, link := range repo.Links {
		if link.Rel == "virtualmachinetemplates" {
			resp, err := rclient.
				SetHeader("Accept", link.Type).
				SetQueryParam("limit", "0").
				Get(link.Href)
			if err != nil {
				return false, template, err
			}

			var templatescol TemplateCollection
			json.Unmarshal([]byte(resp.Body()), &templatescol)
			for _, tmpl := range templatescol.Collection {
				if tmpl.Name == config.Name {
					log.Printf("Found.")
					tmpljson, err := tmpl.ToJson()
					log.Printf("Template '%s' with DTO : %s", tmpl.Name, tmpljson)
					return true, tmpl, err
				}
			}
		}
	}
	log.Printf("Not found.")
	return false, template, nil
}

func (config *Config) BuildDiskDef(template VirtualMachineTemplate) DiskDef {
	var disk DiskDef
	template_url := template.GetLink("edit").Href
	diskUrl := template.GetLink("disk0").Href

	resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	rclient := resty.R().SetBasicAuth(config.ApiUsername, config.ApiPassword)
	resp, err := rclient.
		SetHeader("Accept", template.GetLink("disk0").Type).
		Get(template.GetLink("disk0").Href)
	if err != nil {
		return disk
	}

	var primaryDisk TemplateDisk
	json.Unmarshal([]byte(resp.Body()), &primaryDisk)

	disk = DiskDef{
		Bootable:                  true,
		Sequence:                  0,
		RequiredHDInMB:            0,
		DiskFileFormat:            "",
		VirtualMachineTemplateUrl: template_url,
		DiskUrl:                   diskUrl,
		CurrentPath:               primaryDisk.Path,
	}

	return disk
}

func (config *Config) BuildTemplateDef() VMTDef {
	log.Printf("Trying to parse guest os type %s", config.GuessOsType)
	ostype := OsTypeFromGuest(config.GuessOsType)
	log.Printf("My best guess... %s, version %s", ostype.Os, ostype.Version)

	definition := VMTDef{
		Name:               config.Name,
		Description:        config.Description,
		CategoryName:       config.CategoryName,
		DiskFileFormat:     config.DiskFileFormat,
		RequiredCpu:        config.RequiredCpu,
		RequiredHDInMB:     config.RequiredHDInMB,
		RequiredRamInMB:    config.RequiredRamInMB,
		LoginUser:          config.LoginUser,
		LoginPassword:      config.LoginPassword,
		OsType:             ostype.Os,
		OsVersion:          ostype.Version,
		EthernetDriverType: config.EthernetDriverType,
	}

	return definition
}

func getFilesFromArtifact(config Config, artifact packer.Artifact, suffix string) (string, error) {
	log.Printf("Got the following files: %v", artifact.Files())
	log.Printf("Builder is: %s", config.PackerBuilderType)

	for _, file := range artifact.Files() {
		if strings.HasSuffix(file, suffix) {
			log.Printf("Selecting file %s", file)
			return file, nil
		}
	}
	return "", errors.New("Could not get the disks file from packer artifact.")
}
