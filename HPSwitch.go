package main

import (
	"strings"
	"strconv"

)

const (
	SNMP_LOCATION = ".1.3.6.1.2.1.1.6.0"
	SNMP_NAME	  = ".1.3.6.1.2.1.1.5.0"
	SNMP_INTERFACE_NAME = ".1.3.6.1.2.1.2.2.1.2"
	SNMP_VLAN_NAME     = ".1.3.6.1.2.1.17.7.1.4.3.1.1"
	SNMP_VLAN_MEMBER   = ".1.3.6.1.2.1.17.7.1.4.3.1.2"
	SNMP_VLAN_UNTAGGED = ".1.3.6.1.2.1.17.7.1.4.3.1.4"
	SNMP_INTERFACE_ALIAS = ".1.3.6.1.2.1.31.1.1.1.18"
)
type _vlan struct {
	vlan_id	int
	vlan_name	string
}



type _HPSwitch_Port struct {
	Name 	string
	Description string
	ID	int
	TaggedVlans		[]*_vlan
	UntaggedVlan	[]*_vlan
}
type _HPSwitch struct {
	IP	string
	vlans	[]*_vlan
	ports	[]*_HPSwitch_Port
	Location	string
	Name		string
	snmp	*_SNMP_Helper
}


func HPSwitch_new(ip string, community string) (*_HPSwitch) {
	snmp := SNMP_Helper_new(ip, community)
	hpswitch := &_HPSwitch{
		IP:       ip,
		snmp:     snmp,
	}

	hpswitch.poll()

	return hpswitch
}

func (hps * _HPSwitch) get_switchport_by_id(portID int) (*_HPSwitch_Port) {
	for _, e := range hps.ports {
		if e.ID == portID {
			return e
		}
	}
	port := &_HPSwitch_Port{
		ID: portID,
	}
	hps.ports = append(hps.ports, port)
	return port
}

func (hps * _HPSwitch) get_vlan_by_id(vlan_id int) (*_vlan) {
	for _, e := range hps.vlans {
		if e.vlan_id == vlan_id {
			return e
		}
	}
	vlan := &_vlan{
		vlan_id:   vlan_id,
	}
	hps.vlans = append(hps.vlans, vlan)
	return vlan
}

func (hps * _HPSwitch) poll() {
	hps.Name = hps.snmp.get_oid_string(SNMP_NAME)
	hps.Location = hps.snmp.get_oid_string(SNMP_LOCATION)

	nextInterfaceID := 1

	name, oid := hps.snmp.get_next_oid_string(SNMP_INTERFACE_NAME)
	for {
		if !strings.HasPrefix(oid, SNMP_INTERFACE_NAME+".") {
			break
		}
		id_split := strings.Split(oid, ".")
		id := id_split[len(id_split) - 1]
		_id, _ := strconv.Atoi(id)
		if _id != nextInterfaceID {
			break
		}
		nextInterfaceID ++
		port := hps.get_switchport_by_id(_id)
		port.Name = name
		port.Description = hps.snmp.get_oid_string(SNMP_INTERFACE_ALIAS + "." + id)

		name, oid = hps.snmp.get_next_oid_string(oid)
	}

	name, oid = hps.snmp.get_next_oid_string(SNMP_VLAN_NAME)
	for {
		if !strings.HasPrefix(oid, SNMP_VLAN_NAME + ".") {
			break
		}
		id_split := strings.Split(oid, ".")
		id := id_split[len(id_split) - 1]
		_id, _ := strconv.Atoi(id)

		vlan := hps.get_vlan_by_id(_id)
		vlan.vlan_name = name
		members := hps.snmp.get_oid_hex_string(SNMP_VLAN_MEMBER + "." + id)
		untagged := hps.snmp.get_oid_hex_string(SNMP_VLAN_UNTAGGED + "." + id)
		for x := 0; x < len(members); x++ {
			for y := 0; y < 8; y++ {
				portID := 1 + x * 8 + y
				if portID > len(hps.ports) {
					break
				}
				port := hps.get_switchport_by_id(portID)
				if members[x] & (1<<uint(7 - y)) > 0 {
					if untagged[x] & (1<<uint(7 - y)) > 0 {
						port.UntaggedVlan = append(port.UntaggedVlan, vlan)
					} else {
						port.TaggedVlans = append(port.TaggedVlans, vlan)
					}
				}
			}
		}


		name, oid = hps.snmp.get_next_oid_string(oid)
	}

}