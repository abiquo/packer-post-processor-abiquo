# Abiquo packer postprocessor

A packer postprocessor that will allow you to upload the templates to Abiquo.

## Installation

Having your Go environment ready, run:

```
go get github.com/abiquo/packer-post-processor-abiquo
```

See [packer docs](https://www.packer.io/docs/extending/plugins.html) on how to install the plugin.

## Usage

The plugin provides a new postprocessor named `abiquo` with the following available options:

| Name                          | Description                                                    |
|-------------------------------|----------------------------------------------------------------|
| `api_url`                     | The Abiquo API endpoint URL. Eg. `https://some.host/api` |
| `api_username`                | The Username for the Abiquo API. Eg. `tmpluser` |
| `api_password`                | The password for the aforementioned user. Eg. `somepass` |
| `datacenter`                  | The Abiquo datacenter name where to upload the template. Eg. `dclocation1` |
| `keep_input_artifact`         | Wheter or not to keep the input artifact. Default is `false` |
| `template_name`               | The name of the template to be created in Abiquo. Defaults to the packer VM name. |
| `description`                 | Description of the template. Defaults to the same value as `name`. |
| `category`                    | Category for the template. Defaults to `OS`. |
| `disk_format`                 | The format of the disk being uploaded. If not specified it will be guessed from the input artifact. |
| `cpu`                         | CPU cores to use by the template. Defaults to `1`. |
| `hd_mb`                       | Size of the disk being uploaded. If not specified it will be guessed from the input artifact. |
| `ram_mb`                      | The amount of RAM in MB to be used by the template. Defaults to `1024`. |
| `login_user`                  | The login username for the template. If empty, we will set the same used by packer. |
| `login_password`              | The password for the `login_user` user. If empty, we will set the same used by packer. |
| `eth_driver`                  | The NIC driver to use in the template. Defaults to `E1000`. |
| `chef_enabled`                | Wether or not the template has the Abiquo Chef agent installed. Defaults to `false`. |
| `icon_url`                    | The URL of the icon for the template. No URL will be set by default. |
| `cpu_hotadd`                  | Wether or not the guest will support CPU hotplug. Defaults to `false`. |
| `ram_hotadd`                  | Wether or not the guest will support RAM hotplug. Defaults to `false`. |
| `disk_hotadd`                 | Wether or not the guest will support disks hotplug. Defaults to `false`. |
| `nic_hotadd`                  | Wether or not the guest will support NIC hotplug. Defaults to `false`. |
| `vnc_hotadd`                  | Wether or not the guest will support hot reconfigure of VNC remote access. Defaults to `false`. |

## Example

Examples of usage:

```
"post-processors": [
    [
      {
        "type": "abiquo",
        "api_url": "https://kvmmigration.bcn.abiquo.com/api",
        "api_username": "admin",
        "api_password": "xabiquo",
        "datacenter": "Abiquo-DC",
        "template_name": "000-packertest",
        "cpu": "1",
        "description": "A packer test",
        "eth_driver": "E1000",
        "ram_mb": "1024",
        "icon_url": "{{user `icon`}}"
      }
    ]
```

# License and Authors

* Author:: Marc Cirauqui (marc.cirauqui@abiquo.com)

Copyright:: 2017, Abiquo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
