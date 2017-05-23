package main

type OsType struct {
	Os      string
	Version string
}

func OsTypeFromGuest(os_type string) OsType {
	var os OsType
	switch os_type {
	case "darwin14_64Guest":
		os = OsType{Os: "MACOS", Version: "14 64b"}
	case "debian8_64Guest":
		os = OsType{Os: "DEBIAN", Version: "8 64b"}
	case "debian8":
		os = OsType{Os: "DEBIAN", Version: "8"}
	case "vmkernel6Guest":
		os = OsType{Os: "ESXI", Version: "6"}
	case "coreos64Guest":
		os = OsType{Os: "LINUX", Version: "coreos"}
	case "windows9Server64Guest":
		os = OsType{Os: "WINDOWS", Version: "10 Server"}
	case "windows9_64Guest":
		os = OsType{Os: "WINDOWS", Version: "10 64b"}
	case "windows9":
		os = OsType{Os: "WINDOWS", Version: "10"}
	case "darwin12_64Guest":
		os = OsType{Os: "MACOS", Version: "12"}
	case "darwin13_64Guest":
		os = OsType{Os: "MACOS", Version: "13"}
	case "debian7_64Guest":
		os = OsType{Os: "DEBIAN", Version: "7 64b"}
	case "debian7":
		os = OsType{Os: "DEBIAN", Version: "7"}
	case "genericLinuxGuest":
		os = OsType{Os: "LINUX", Version: "generic"}
	case "other3xLinux64Guest":
		os = OsType{Os: "LINUX", Version: "3x 64b"}
	case "other3xLinux":
		os = OsType{Os: "LINUX", Version: "3x"}
	case "rhel7_64Guest":
		os = OsType{Os: "RHEL", Version: "7 64b"}
	case "rhel7":
		os = OsType{Os: "RHEL", Version: "7"}
	case "sles12_64Guest":
		os = OsType{Os: "SLES", Version: "12 64b"}
	case "sles12":
		os = OsType{Os: "SLES", Version: "12"}
	case "windowsHyperVGuest":
		os = OsType{Os: "WINDOWS", Version: "hyperv"}
	case "darwinGuest":
		os = OsType{Os: "MACOS", Version: ""}
	case "darwin64Guest":
		os = OsType{Os: "MACOS", Version: "64b"}
	case "darwin10_64Guest":
		os = OsType{Os: "MACOS", Version: "10 64b"}
	case "darwin10Guest":
		os = OsType{Os: "MACOS", Version: "10"}
	case "darwin11_64Guest":
		os = OsType{Os: "MACOS", Version: "11 64b"}
	case "darwin11Guest":
		os = OsType{Os: "MACOS", Version: "11"}
	case "solaris6Guest":
		os = OsType{Os: "SOLARIS", Version: "6"}
	case "solaris7Guest":
		os = OsType{Os: "SOLARIS", Version: "7"}
	case "solaris8Guest":
		os = OsType{Os: "SOLARIS", Version: "8"}
	case "solaris9Guest":
		os = OsType{Os: "SOLARIS", Version: "9"}
	case "solaris10Guest":
		os = OsType{Os: "SOLARIS", Version: "10"}
	case "sjdsGuest":
		os = OsType{Os: "SOLARIS", Version: "Sun Java Desktop System"}
	case "solaris11_64Guest":
		os = OsType{Os: "SOLARIS_64", Version: "11"}
	case "solaris10_64Guest":
		os = OsType{Os: "SOLARIS_64", Version: "10"}
	case "redhatGuest":
		os = OsType{Os: "RHEL", Version: ""}
	case "rhel2Guest":
		os = OsType{Os: "RHEL", Version: "2"}
	case "rhel3Guest":
		os = OsType{Os: "RHEL", Version: "3"}
	case "rhel4Guest":
		os = OsType{Os: "RHEL", Version: "4"}
	case "rhel5Guest":
		os = OsType{Os: "RHEL", Version: "5"}
	case "rhel6Guest":
		os = OsType{Os: "RHEL", Version: "6"}
	case "rhel3_64Guest":
		os = OsType{Os: "RHEL_64", Version: "3"}
	case "rhel4_64Guest":
		os = OsType{Os: "RHEL_64", Version: "4"}
	case "rhel5_64Guest":
		os = OsType{Os: "RHEL_64", Version: "5"}
	case "rhel6_64Guest":
		os = OsType{Os: "RHEL_64", Version: "6"}
	case "suseGuest":
		os = OsType{Os: "SUSE", Version: ""}
	case "opensuseGuest":
		os = OsType{Os: "SUSE", Version: "Open"}
	case "suse64Guest":
		os = OsType{Os: "SUSE_64", Version: ""}
	case "opensuse64Guest":
		os = OsType{Os: "SUSE_64", Version: "Open"}
	case "slesGuest":
		os = OsType{Os: "SLES", Version: ""}
	case "sles10Guest":
		os = OsType{Os: "SLES", Version: "10"}
	case "sles11Guest":
		os = OsType{Os: "SLES", Version: "11"}
	case "sles10_64Guest":
		os = OsType{Os: "SLES_64", Version: "10"}
	case "sles11_64Guest":
		os = OsType{Os: "SLES_64", Version: "11"}
	case "sles64Guest":
		os = OsType{Os: "SLES_64", Version: ""}
	case "oesGuest":
		os = OsType{Os: "NOVELL_OES", Version: ""}
	case "nld9Guest":
		os = OsType{Os: "NOVELL_OES", Version: "9"}
	case "netware4Guest":
		os = OsType{Os: "NOVELL_OES", Version: "NetWare 4"}
	case "netware5Guest":
		os = OsType{Os: "NOVELL_OES", Version: "NetWare 5"}
	case "netware6Guest":
		os = OsType{Os: "NOVELL_OES", Version: "NetWare 6"}
	case "mandrivaGuest":
		os = OsType{Os: "MANDRIVA", Version: ""}
	case "mandriva64Guest":
		os = OsType{Os: "MANDRIVA_64", Version: ""}
	case "turboLinuxGuest":
		os = OsType{Os: "TURBOLINUX", Version: ""}
	case "turboLinux64Guest":
		os = OsType{Os: "TURBOLINUX_64", Version: ""}
	case "ubuntuGuest":
		os = OsType{Os: "UBUNTU", Version: ""}
	case "ubuntu64Guest":
		os = OsType{Os: "UBUNTU_64", Version: ""}
	case "debian4Guest":
		os = OsType{Os: "DEBIAN", Version: "4"}
	case "debian5Guest":
		os = OsType{Os: "DEBIAN", Version: "5"}
	case "debian6Guest":
		os = OsType{Os: "DEBIAN", Version: "6"}
	case "debian4_64Guest":
		os = OsType{Os: "DEBIAN_64", Version: "4"}
	case "debian5_64Guest":
		os = OsType{Os: "DEBIAN_64", Version: "5"}
	case "debian6_64Guest":
		os = OsType{Os: "DEBIAN_64", Version: "6"}
	case "windows8Guest":
		os = OsType{Os: "WINDOWS_8", Version: ""}
	case "windows8_64Guest":
		os = OsType{Os: "WINDOWS_8_64", Version: ""}
	case "windows8Server64Guest":
		os = OsType{Os: "WINDOWS_SERVER_2012", Version: ""}
	case "other24xLinuxGuest":
		os = OsType{Os: "LINUX_2_4", Version: ""}
	case "other24xLinux64Guest":
		os = OsType{Os: "LINUX_2_4_64", Version: ""}
	case "other26xLinuxGuest":
		os = OsType{Os: "LINUX_2_6", Version: ""}
	case "other26xLinux64Guest":
		os = OsType{Os: "LINUX_2_6_64", Version: ""}
	case "centosGuest":
		os = OsType{Os: "CENTOS", Version: ""}
	case "centos64Guest":
		os = OsType{Os: "CENTOS_64", Version: ""}
	case "oracleLinuxGuest":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX", Version: ""}
	case "oracleLinux64Guest":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX_64", Version: ""}
	case "otherGuest64":
		os = OsType{Os: "OTHER_64", Version: ""}
	case "freebsdGuest":
		os = OsType{Os: "FREEBSD", Version: ""}
	case "freebsd64Guest":
		os = OsType{Os: "FREEBSD_64", Version: ""}
	case "windows7Guest":
		os = OsType{Os: "WINDOWS_7", Version: ""}
	case "windows7Server64Guest":
		os = OsType{Os: "WINDOWS_7", Version: "Server"}
	case "windows7_64Guest":
		os = OsType{Os: "WINDOWS_7", Version: "64b"}
	case "winNetStandardGuest":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: ""}
	case "winNetWebGuest":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Web"}
	case "winNetBusinessGuest":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Business"}
	case "winNetEnterpriseGuest":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Enterprise"}
	case "winNetDatacenterGuest":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Datacenter"}
	case "winNetStandard64Guest":
		os = OsType{Os: "WINDOWS_SERVER_2003_64", Version: ""}
	case "winNetEnterprise64Guest":
		os = OsType{Os: "WINDOWS_SERVER_2003_64", Version: "Enterprise"}
	case "winNetDatacenter64Guest":
		os = OsType{Os: "WINDOWS_SERVER_2003_64", Version: "Datacenter"}
	case "win31Guest":
		os = OsType{Os: "WINDOWS", Version: "3.1"}
	case "win95Guest":
		os = OsType{Os: "WINDOWS", Version: "95"}
	case "win98Guest":
		os = OsType{Os: "WINDOWS", Version: "98"}
	case "winMeGuest":
		os = OsType{Os: "WINDOWS", Version: "ME"}
	case "winNTGuest":
		os = OsType{Os: "WINDOWS", Version: "NT"}
	case "win2000ProGuest":
		os = OsType{Os: "WINDOWS", Version: "2000 Pro"}
	case "win2000AdvServGuest":
		os = OsType{Os: "WINDOWS", Version: "2000 Avd"}
	case "win2000ServGuest":
		os = OsType{Os: "WINDOWS", Version: "2000 Server"}
	case "winXPHomeGuest":
		os = OsType{Os: "WINDOWS", Version: "XP Home"}
	case "winXPPro64Guest":
		os = OsType{Os: "WINDOWS", Version: "XP Pro 64b"}
	case "winXPProGuest":
		os = OsType{Os: "WINDOWS", Version: "XP Pro"}
	case "winLonghorn64Guest":
		os = OsType{Os: "WINDOWS", Version: "Longhorn 64b"}
	case "winLonghornGuest":
		os = OsType{Os: "WINDOWS", Version: "Longhorn"}
	case "winVista64Guest":
		os = OsType{Os: "WINDOWS", Version: "Vista 64b"}
	case "winVistaGuest":
		os = OsType{Os: "WINDOWS", Version: "Vista"}
	case "otherLinuxGuest":
		os = OsType{Os: "LINUX", Version: ""}
	case "fedoraGuest":
		os = OsType{Os: "LINUX", Version: "Fedora"}
	case "otherLinux64Guest":
		os = OsType{Os: "LINUX_64", Version: ""}
	case "fedora64Guest":
		os = OsType{Os: "LINUX_64", Version: "Fedora"}
	case "vmkernel5Guest":
		os = OsType{Os: "ESXI", Version: "5"}
	case "vmkernelGuest":
		os = OsType{Os: "ESXI", Version: "4"}
	case "eComStationGuest":
		os = OsType{Os: "ECOMSTATION_32", Version: "1"}
	case "eComStation2Guest":
		os = OsType{Os: "ECOMSTATION_32", Version: "2"}
	case "otherGuest":
	default:
		os = OsType{Os: "UNRECOGNIZED", Version: ""}
	}
	return os
}
