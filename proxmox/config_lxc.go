package proxmox

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

// LXC options for the Proxmox API
type configLxc struct {
	Ostemplate         string         `json:"ostemplate"`
        Arch               string         `json:"arch"`
        BWLimit            int            `json:"bwlimit,omitempty"`
        CMode              string         `json:"cmode"`
        Console            bool           `json:"console"`
        Cores              int            `json:"cores,omitempty"`
        CPULimit           int            `json:"cpulimit"`
        CPUUnits           int            `json:"cpuunits"`
        Description        string         `json:"description,omitempty"`
        Features           QemuDevice     `json:"features,omitempty"`
        Force              bool           `json:"force,omitempty"`
        Hookscript         string         `json:"hookscript,omitempty"`
	Hostname           string         `json:"hostname,omitempty"`
        IgnoreUnpackErrors bool           `json:"ignore-unpack-errors,omitempty"`
	Lock               string         `json:"lock,omitempty"`
        Memory             int            `json:"memory"`
        Mountpoints        QemuDevices    `json:"mountpoints,omitempty"`
        Nameserver         string         `json:"nameserver,omitempty"`
	Networks           QemuDevices    `json:"networks,omitempty"`
        OnBoot             bool           `json:"onboot"`
        OsType             string         `json:"ostype,omitempty"`
	Password           string         `json:"password,omitempty"`
	Pool               string         `json:"pool,omitempty"`
        Protection         bool           `json:"protection"`
        Restore            bool           `json:"restore,omitempty"`
	RootFs             string         `json:"rootfs,omitempty"`
	SearchDomain       string         `json:"searchdomain,omitempty"`
	SSHPublicKeys      string         `json:"ssh-public-keys,omitempty"`
        Start              bool           `json:"start"`
	Startup            string         `json:"startup,omitempty"`
	Storage            string         `json:"storage"`
        Swap               int            `json:"swap"`
	Template           bool           `json:"template,omitempty"`
        Tty                int            `json:"tty"`
	Unique             bool           `json:"unique,omitempty"`
	Unprivileged       bool           `json:"unprivileged"`
	Unused             []string       `json:"unused,omitempty"`
}

func NewConfigLxc() (configLxc) {
	return configLxc{
		Arch: "amd64",
		CMode: "tty",
		Console: true,
		CPULimit: 0,
		CPUUnits: 1024,
		Memory: 512,
		OnBoot: false,
		Protection: false,
		Start: false,
		Storage: "local",
		Swap: 512,
		Template: false,
		Tty: 2,
		Unprivileged: false,
	}
}

func NewConfigLxcFromJson(io io.Reader) (config configLxc, err error) {
	config = NewConfigLxc()
	err = json.NewDecoder(io).Decode(config)
	if err != nil {
		log.Fatal(err)
		return config, err
	}
	log.Println(config)
	return
}

// create LXC container using the Proxmox API
func (config configLxc) CreateLxc(vmr *VmRef, client *Client) (err error) {
	vmr.SetVmType("lxc")

        // convert config to map
        params, _ := json.Marshal(&config)
        var paramMap map[string]interface{}
        json.Unmarshal(params, &paramMap)

        // build list of features
	// add features as parameter list to lxc parameters
	// this overwrites the orginal formatting with a
        // comma separated list of "key=value" pairs
	featuresParam := QemuDeviceParam{}
	featuresParam = featuresParam.createDeviceParam(config.Features, nil)
	paramMap["features"] = strings.Join(featuresParam, ",")

        // build list of mountpoints
	// this does the same as for the feature list
        // except that there can be multiple of these mountpoint sets
        // and each mountpoint set comes with a new id
	for mpID, mpConfMap := range config.Mountpoints {
		mpConfParam := QemuDeviceParam{}
		mpConfParam = mpConfParam.createDeviceParam(mpConfMap, nil)

		// add mp to lxc parameters
		mpName := fmt.Sprintf("mp%v", mpID)
		paramMap[mpName] = strings.Join(mpConfParam, ",")
        }

        // build list of network parameters
	for nicID, nicConfMap := range config.Networks {
		nicConfParam := QemuDeviceParam{}
		nicConfParam = nicConfParam.createDeviceParam(nicConfMap, nil)

		// add nic to lxc parameters
		nicName := fmt.Sprintf("net%v", nicID)
		paramMap[nicName] = strings.Join(nicConfParam, ",")
        }

        // build list of unused volumes for sake of completenes,
        // even if it is not recommended to change these volumes manually
	for volID, vol := range config.Unused {
		// add volume to lxc parameters
		volName := fmt.Sprintf("unused%v", volID)
		paramMap[volName] = vol
        }

        // now that we concatenated the key value parameter
        // list for the networks, mountpoints and unused volumes,
        // remove the original keys, since the Proxmox API does
        // not know how to handle this key
        delete(paramMap, "networks")
        delete(paramMap, "mountpoints")
        delete(paramMap, "unused")

        // amend vmid
        paramMap["vmid"] = vmr.vmId

	exitStatus, err := client.CreateLxcContainer(vmr.node, paramMap)
	if err != nil {
		return fmt.Errorf("Error creating LXC container: %v, error status: %s (params: %v)", err, exitStatus, params)
	}
	return
}
