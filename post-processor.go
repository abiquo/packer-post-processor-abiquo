package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	vmwcommon "github.com/hashicorp/packer/builder/vmware/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"

	"github.com/abiquo/api-go-client"
)

func (p *PostProcessor) GetDiskFileInfo(artifact packer.Artifact) (string, int, error) {
	var abqFormat string

	file, err := getFilesFromArtifact(p.config, artifact, "vmdk")
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

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	// Fields from config file
	ApiUrl            string `mapstructure:"api_url"`
	Insecure          bool   `mapstructure:"api_insecure"`
	ApiUsername       string `mapstructure:"api_username"`
	ApiPassword       string `mapstructure:"api_password"`
	AppKey            string `mapstructure:"app_key"`
	AppSecret         string `mapstructure:"app_secret"`
	AccessToken       string `mapstructure:"access_token"`
	AccessTokenSecret string `mapstructure:"access_token_secret"`
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

func (p *PostProcessor) getClient() *abiquo_api.AbiquoClient {
	if p.config.AppKey != "" {
		return abiquo_api.GetOAuthClient(p.config.ApiUrl, p.config.AppKey, p.config.AppSecret, p.config.AccessToken, p.config.AccessTokenSecret, p.config.Insecure)
	} else if p.config.ApiUsername != "" {
		return abiquo_api.GetClient(p.config.ApiUrl, p.config.ApiUsername, p.config.ApiPassword, p.config.Insecure)
	}
	return nil
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

	if p.config.ApiUsername == "" && p.config.AppKey == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Abiquo API username or Oauth app keys are missing!"))
	}

	if p.config.ApiUsername != "" && p.config.ApiPassword == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Abiquo API username provided but no password!"))
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

	vmxFile, err := getFilesFromArtifact(p.config, artifact, "vmx")
	if err != nil {
		ui.Message(fmt.Sprintf("err getFiles: %s", err))
	}
	vmxData, err := vmwcommon.ReadVMX(vmxFile)
	if err != nil {
		ui.Message(fmt.Sprintf("err ReadVMX: %s", err))
	}
	vmName = vmxData["displayname"]
	vmxOsType = vmxData["guestos"]
	log.Printf("DispName: '%s', Guest: '%s'", vmxData["displayname"], vmxData["guestos"])

	if p.config.Name == "" {
		p.config.Name = vmName
	}
	if p.config.Description == "" {
		p.config.Description = vmName
	}
	if p.config.GuessOsType == "" {
		p.config.GuessOsType = vmxOsType
	}

	abq := p.getClient()

	ui.Say("Looking up the repo URL for datacenter '" + p.config.Datacenter + "'")
	repo, err := p.FindRepoUrl()
	log.Printf("Repo found at: '%s'", repo.RepositoryLocation)
	if err != nil {
		return artifact, p.config.KeepInputArtifact, err
	}

	ui.Say("Checking if a template named '" + p.config.Name + "' already exists...")
	var template abiquo_api.VirtualMachineTemplate
	exists, template, err := p.CheckTemplateExists(repo)
	if err != nil {
		return artifact, p.config.KeepInputArtifact, err
	}

	if exists {
		ui.Say("Template already exists. Replacing primary disk...")
		diskdef, err := p.BuildDiskDef(template)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}

		abqdiskformat, size, err := p.GetDiskFileInfo(artifact)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}

		diskdef.RequiredHDInMB = size
		diskdef.DiskFileFormat = abqdiskformat
		file, err := getFilesFromArtifact(p.config, artifact, "vmdk")
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}
		template, err := template.ReplacePrimaryDisk(abq, diskdef, file)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}

		template_link, _ := template.GetLink("edit")
		ui.Say("Upload complete. The URL of the updated template is " + template_link.Href)
	} else {
		ui.Say("Uploading template...")
		templatedef := p.config.BuildTemplateDef()
		abqdiskformat, size, err := p.GetDiskFileInfo(artifact)
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}

		file, err := getFilesFromArtifact(p.config, artifact, "vmdk")
		if err != nil {
			return artifact, p.config.KeepInputArtifact, err
		}

		templatedef.DiskFileFormat = abqdiskformat
		templatedef.RequiredHDInMB = strconv.Itoa(size)
		templateUploaded, err := templatedef.Upload(abq, repo, file)
		log.Printf("Uploaded template : %v", templateUploaded)
		exists, template, err = p.CheckTemplateExists(repo)
		if exists {
			log.Printf("Found new template '%s'", template.Name)
		}
		template_link, _ := template.GetLink("edit")
		ui.Say("Upload complete. The URL of the new template is " + template_link.Href)
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
	template.Update(abq)

	template_link, _ := template.GetLink("edit")
	newArtifact := &Artifact{Url: template_link.Href}
	return newArtifact, false, nil
}

func (p *PostProcessor) FindRepoUrl() (abiquo_api.Repo, error) {
	var repo abiquo_api.Repo
	var ent abiquo_api.Enterprise

	abq := p.getClient()
	log.Printf("Trying to find the repo URL for Datacenter '%s'\n", p.config.Datacenter)

	user, err := abq.Login()
	if err != nil {
		return repo, err
	}

	ent, err = user.GetEnterprise(abq)
	if err != nil {
		return repo, err
	}

	repos, err := ent.GetRepos(abq)
	for _, repo := range repos {
		dc_lnk, _ := repo.GetLink("datacenter")
		if dc_lnk.Title == p.config.Datacenter {
			log.Printf("Found URL '%s' (%s) for DC '%s'", repo.Name, repo.RepositoryLocation, p.config.Datacenter)
			return repo, err
		}
	}

	errorMsg := fmt.Sprintf("Could not find the URL for DC '%s'!", p.config.Datacenter)
	return repo, errors.New(errorMsg)
}

func (p *PostProcessor) CheckTemplateExists(repo abiquo_api.Repo) (bool, abiquo_api.VirtualMachineTemplate, error) {
	var template abiquo_api.VirtualMachineTemplate
	log.Printf("Checking if a template with name '%s' already exists.", p.config.Name)
	abq := p.getClient()

	templates, err := repo.GetTemplates(abq)
	if err != nil {
		return false, template, err
	}

	for _, t := range templates {
		if t.Name == p.config.Name {
			template_lnk, _ := t.GetLink("edit")
			log.Printf("Found template in URL: '%s'", template_lnk.Href)
			return true, t, nil
		}
	}
	log.Printf("Not found.")
	return false, template, nil
}

func (p *PostProcessor) BuildDiskDef(template abiquo_api.VirtualMachineTemplate) (abiquo_api.DiskDef, error) {
	var disk abiquo_api.DiskDef
	var primaryDisk abiquo_api.Disk
	abq := p.getClient()

	template_link, _ := template.GetLink("edit")
	template_url := template_link.Href
	disks, err := template.GetDisks(abq)
	if err != nil {
		return disk, err
	}

	for _, d := range disks {
		if d.Sequence == 0 {
			primaryDisk = d
		}
	}
	disk_link, _ := primaryDisk.GetLink("edit")
	diskUrl := disk_link.Href

	disk = abiquo_api.DiskDef{
		Bootable:                  true,
		Sequence:                  0,
		RequiredHDInMB:            0,
		DiskFileFormat:            "",
		VirtualMachineTemplateUrl: template_url,
		DiskUrl:                   diskUrl,
		CurrentPath:               primaryDisk.Path,
	}

	return disk, nil
}

func (config *Config) BuildTemplateDef() abiquo_api.VMTDef {
	log.Printf("Trying to parse guest os type '%s'", config.GuessOsType)
	ostype := OsTypeFromGuest(config.GuessOsType)
	log.Printf("My best guess... '%s', version '%s'", ostype.Os, ostype.Version)

	definition := abiquo_api.VMTDef{
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
