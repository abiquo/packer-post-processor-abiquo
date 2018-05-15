package main

type OsType struct {
	Os      string
	Version string
}

func OsTypeFromGuest(osType string) OsType {
	var os OsType
	switch osType {
	// version 6.5
	case "centos6-64":
		os = OsType{Os: "CENTOS_64", Version: "6"}
	case "centos7-64":
		os = OsType{Os: "CENTOS_64", Version: "7"}
	case "centos6":
		os = OsType{Os: "CENTOS", Version: "6"}
	case "centos7":
		os = OsType{Os: "CENTOS", Version: "7"}
	case "darwin15-64":
		os = OsType{Os: "MACOS", Version: "15"}
	case "darwin16-64":
		os = OsType{Os: "MACOS", Version: "16"}
	case "debian10-64":
		os = OsType{Os: "DEBIAN_64", Version: "10"}
	case "debian10":
		os = OsType{Os: "DEBIAN", Version: "10"}
	case "debian9-64":
		os = OsType{Os: "DEBIAN_64", Version: "9"}
	case "debian9":
		os = OsType{Os: "DEBIAN", Version: "9"}
	case "oraclelinux6-64":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX_64", Version: "6"}
	case "oraclelinux6":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX", Version: "6"}
	case "oraclelinux7-64":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX_64", Version: "7"}
	case "oraclelinux7":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX", Version: "7"}
	case "vmkernel65":
		os = OsType{Os: "ESXI", Version: "6.5"}
	case "vmware-photon-64":
		os = OsType{Os: "ESXI", Version: "photon"}
	// version 6
	case "darwin14-64":
		os = OsType{Os: "MACOS", Version: "14 64b"}
	case "debian8-64":
		os = OsType{Os: "DEBIAN", Version: "8 64b"}
	case "debian8":
		os = OsType{Os: "DEBIAN", Version: "8"}
	case "vmkernel6":
		os = OsType{Os: "ESXI", Version: "6"}
	case "coreos-64":
		os = OsType{Os: "LINUX", Version: "coreos"}
	case "windows9srv-64":
		os = OsType{Os: "WINDOWS", Version: "Server 2016"}
	case "windows9-64":
		os = OsType{Os: "WINDOWS", Version: "10 64b"}
	case "windows9":
		os = OsType{Os: "WINDOWS", Version: "10"}
	// version 55
	case "darwin12-64":
		os = OsType{Os: "MACOS", Version: "12"}
	case "darwin13-64":
		os = OsType{Os: "MACOS", Version: "13"}
	case "debian7-64":
		os = OsType{Os: "DEBIAN", Version: "7 64b"}
	case "debian7":
		os = OsType{Os: "DEBIAN", Version: "7"}
	case "genericLinuxGuest":
		os = OsType{Os: "LINUX", Version: "generic"}
	case "other3xlinux-64":
		os = OsType{Os: "LINUX", Version: "3x 64b"}
	case "other3xlinux":
		os = OsType{Os: "LINUX", Version: "3x"}
	case "rhel7-64":
		os = OsType{Os: "RHEL", Version: "7 64b"}
	case "rhel7":
		os = OsType{Os: "RHEL", Version: "7"}
	case "sles12-64":
		os = OsType{Os: "SLES", Version: "12 64b"}
	case "sles12":
		os = OsType{Os: "SLES", Version: "12"}
	case "winhyperv":
		os = OsType{Os: "WINDOWS", Version: "hyperv"}
	// 51
	case "darwin":
		os = OsType{Os: "MACOS", Version: ""}
	case "darwin-64":
		os = OsType{Os: "MACOS", Version: "64b"}
	case "darwin10-64":
		os = OsType{Os: "MACOS", Version: "10 64b"}
	case "darwin10":
		os = OsType{Os: "MACOS", Version: "10"}
	case "darwin11-64":
		os = OsType{Os: "MACOS", Version: "11 64b"}
	case "darwin11":
		os = OsType{Os: "MACOS", Version: "11"}
	case "solaris6":
		os = OsType{Os: "SOLARIS", Version: "6"}
	case "solaris7":
		os = OsType{Os: "SOLARIS", Version: "7"}
	case "solaris8":
		os = OsType{Os: "SOLARIS", Version: "8"}
	case "solaris9":
		os = OsType{Os: "SOLARIS", Version: "9"}
	case "solaris10":
		os = OsType{Os: "SOLARIS", Version: "10"}
	case "sjds":
		os = OsType{Os: "SOLARIS", Version: "Sun Java Desktop System"}
	case "solaris11-64":
		os = OsType{Os: "SOLARIS_64", Version: "11"}
	case "solaris10-64":
		os = OsType{Os: "SOLARIS_64", Version: "10"}
	case "redhat":
		os = OsType{Os: "RHEL", Version: ""}
	case "rhel2":
		os = OsType{Os: "RHEL", Version: "2"}
	case "rhel3":
		os = OsType{Os: "RHEL", Version: "3"}
	case "rhel4":
		os = OsType{Os: "RHEL", Version: "4"}
	case "rhel5":
		os = OsType{Os: "RHEL", Version: "5"}
	case "rhel6":
		os = OsType{Os: "RHEL", Version: "6"}
	case "rhel3-64":
		os = OsType{Os: "RHEL_64", Version: "3"}
	case "rhel4-64":
		os = OsType{Os: "RHEL_64", Version: "4"}
	case "rhel5-64":
		os = OsType{Os: "RHEL_64", Version: "5"}
	case "rhel6-64":
		os = OsType{Os: "RHEL_64", Version: "6"}
	case "suse":
		os = OsType{Os: "SUSE", Version: ""}
	case "opensuse":
		os = OsType{Os: "SUSE", Version: "Open"}
	case "suse64":
		os = OsType{Os: "SUSE_64", Version: ""}
	case "opensuse64":
		os = OsType{Os: "SUSE_64", Version: "Open"}
	case "sles":
		os = OsType{Os: "SLES", Version: ""}
	case "sles10":
		os = OsType{Os: "SLES", Version: "10"}
	case "sles11":
		os = OsType{Os: "SLES", Version: "11"}
	case "sles10-64":
		os = OsType{Os: "SLES_64", Version: "10"}
	case "sles11-64":
		os = OsType{Os: "SLES_64", Version: "11"}
	case "sles-64":
		os = OsType{Os: "SLES_64", Version: ""}
	case "oes":
		os = OsType{Os: "NOVELL_OES", Version: ""}
	case "nld9":
		os = OsType{Os: "NOVELL_OES", Version: "9"}
	case "netware4":
		os = OsType{Os: "NOVELL_OES", Version: "NetWare 4"}
	case "netware5":
		os = OsType{Os: "NOVELL_OES", Version: "NetWare 5"}
	case "netware6":
		os = OsType{Os: "NOVELL_OES", Version: "NetWare 6"}
	case "mandriva":
		os = OsType{Os: "MANDRIVA", Version: ""}
	case "mandriva-64":
		os = OsType{Os: "MANDRIVA_64", Version: ""}
	case "turbolinux":
		os = OsType{Os: "TURBOLINUX", Version: ""}
	case "turbolinux-64":
		os = OsType{Os: "TURBOLINUX_64", Version: ""}
	case "ubuntu":
		os = OsType{Os: "UBUNTU", Version: ""}
	case "ubuntu-64":
		os = OsType{Os: "UBUNTU_64", Version: ""}
	case "debian4":
		os = OsType{Os: "DEBIAN", Version: "4"}
	case "debian5":
		os = OsType{Os: "DEBIAN", Version: "5"}
	case "debian6":
		os = OsType{Os: "DEBIAN", Version: "6"}
	case "debian4-64":
		os = OsType{Os: "DEBIAN_64", Version: "4"}
	case "debian5-64":
		os = OsType{Os: "DEBIAN_64", Version: "5"}
	case "debian6-64":
		os = OsType{Os: "DEBIAN_64", Version: "6"}
	case "windows8":
		os = OsType{Os: "WINDOWS_8", Version: ""}
	case "windows8-64":
		os = OsType{Os: "WINDOWS_8_64", Version: ""}
	case "windows8srv-64":
		os = OsType{Os: "WINDOWS_SERVER_2012", Version: ""}
	case "other24xlinux":
		os = OsType{Os: "LINUX_2_4", Version: ""}
	case "other24xlinux-64":
		os = OsType{Os: "LINUX_2_4_64", Version: ""}
	case "other26xlinux":
		os = OsType{Os: "LINUX_2_6", Version: ""}
	case "other26xlinux-64":
		os = OsType{Os: "LINUX_2_6_64", Version: ""}
	case "centos":
		os = OsType{Os: "CENTOS", Version: ""}
	case "centos-64":
		os = OsType{Os: "CENTOS_64", Version: ""}
	case "oraclelinux":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX", Version: ""}
	case "oraclelinux-64":
		os = OsType{Os: "ORACLE_ENTERPRISE_LINUX_64", Version: ""}
	case "other-64":
		os = OsType{Os: "OTHER_64", Version: ""}
	case "freebsd":
		os = OsType{Os: "FREEBSD", Version: ""}
	case "freebsd-64":
		os = OsType{Os: "FREEBSD_64", Version: ""}
	case "windows7":
		os = OsType{Os: "WINDOWS_7", Version: ""}
	case "windows7srv-64":
		os = OsType{Os: "WINDOWS_7", Version: "Server"}
	case "windows7-64":
		os = OsType{Os: "WINDOWS_7", Version: "64b"}
	case "winnetstandard":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: ""}
	case "winnetweb":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Web"}
	case "winnetbussiness":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Business"}
	case "winnetenterprise":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Enterprise"}
	case "winnetdatacenter":
		os = OsType{Os: "WINDOWS_SERVER_2003", Version: "Datacenter"}
	case "winnetstandard-64":
		os = OsType{Os: "WINDOWS_SERVER_2003_64", Version: ""}
	case "winnetenterprise-64":
		os = OsType{Os: "WINDOWS_SERVER_2003_64", Version: "Enterprise"}
	case "winnetdatacenter-64":
		os = OsType{Os: "WINDOWS_SERVER_2003_64", Version: "Datacenter"}
	case "win31":
		os = OsType{Os: "WINDOWS", Version: "3.1"}
	case "win95":
		os = OsType{Os: "WINDOWS", Version: "95"}
	case "win98":
		os = OsType{Os: "WINDOWS", Version: "98"}
	case "winme":
		os = OsType{Os: "WINDOWS", Version: "ME"}
	case "winnt":
		os = OsType{Os: "WINDOWS", Version: "NT"}
	case "win2000pro":
		os = OsType{Os: "WINDOWS", Version: "2000 Pro"}
	case "win2000advserv":
		os = OsType{Os: "WINDOWS", Version: "2000 Avd"}
	case "win2000serv":
		os = OsType{Os: "WINDOWS", Version: "2000 Server"}
	case "winxphome":
		os = OsType{Os: "WINDOWS", Version: "XP Home"}
	case "winxppro-64":
		os = OsType{Os: "WINDOWS", Version: "XP Pro 64b"}
	case "winxppro":
		os = OsType{Os: "WINDOWS", Version: "XP Pro"}
	case "longhorn-64":
		os = OsType{Os: "WINDOWS", Version: "Longhorn 64b"}
	case "longhorn":
		os = OsType{Os: "WINDOWS", Version: "Longhorn"}
	case "winvista-64":
		os = OsType{Os: "WINDOWS", Version: "Vista 64b"}
	case "winvista":
		os = OsType{Os: "WINDOWS", Version: "Vista"}
	case "otherlinux":
		os = OsType{Os: "LINUX", Version: ""}
	case "fedora":
		os = OsType{Os: "LINUX", Version: "Fedora"}
	case "otherlinux-64":
		os = OsType{Os: "LINUX_64", Version: ""}
	case "fedora-64":
		os = OsType{Os: "LINUX_64", Version: "Fedora"}
	case "vmkernel5":
		os = OsType{Os: "ESXI", Version: "5"}
	case "vmkernel":
		os = OsType{Os: "ESXI", Version: "4"}
	case "ecomstation":
		os = OsType{Os: "ECOMSTATION_32", Version: "1"}
	case "ecomstation2":
		os = OsType{Os: "ECOMSTATION_32", Version: "2"}
	// UNRECOGNIZED
	case "other":
	default:
		os = OsType{Os: "UNRECOGNIZED", Version: ""}
	}

	return os
}
