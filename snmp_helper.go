package main
import (
        "github.com/alouca/gosnmp"
	"log"
)

type _SNMP_Helper struct {
	IP string
	Community string
	snmp	*gosnmp.GoSNMP

}

func SNMP_Helper_new(IP string, Community string) (*_SNMP_Helper) {
	snmp, err := gosnmp.NewGoSNMP(IP, Community, gosnmp.Version1, 1)
	if (err != nil) {
		log.Fatalln(err)
	}

	helper := &_SNMP_Helper{
		IP:        IP,
		Community: Community,
		snmp:      snmp,
	}

	return helper
}


func (snmp * _SNMP_Helper) get_oid_int(oid string) (int) {
	reply, err := snmp.snmp.Get(oid)
	if (err != nil) {
		log.Fatalln(reply)
	}
	value := reply.Variables[0]

	var int_val int
	switch (value.Type) {
	case gosnmp.Integer:
		int_val = int((value.Value).(int))
	case gosnmp.Counter32:
		int_val = int((value.Value).(uint64))
	case gosnmp.Gauge32:
		int_val = int((value.Value).(uint64))
	case gosnmp.TimeTicks:
		int_val = int((value.Value).(int))
	default:
		log.Fatalln("Value of oid", oid, "Is no number!")
	}
	return int_val

}
func (snmp * _SNMP_Helper) get_oid_string(oid string) (string) {
	reply, err := snmp.snmp.Get(oid)
	if (err != nil) {
		log.Fatalln(reply)
	}
	value := reply.Variables[0]
	if (value.Type != gosnmp.OctetString) {
		return ""
	}


	return (value.Value).(string)

}

func (snmp * _SNMP_Helper) get_oid_hex_string(oid string) ([]uint8) {
	reply, err := snmp.snmp.Get(oid)
	if (err != nil) {
		log.Fatalln(reply)
	}
	value := reply.Variables[0]

	if (value.Type != gosnmp.OctetString) {
		log.Println("Value is no hex-String")
	}


	return []byte((value.Value).(string))

}

func (snmp * _SNMP_Helper) get_next_oid_int(oid string) (int, string) {
	reply, err := snmp.snmp.GetNext(oid)
	if (err != nil) {
		log.Fatalln(reply)
	}
	value := reply.Variables[0]

	var int_val int
	switch (value.Type) {
	case gosnmp.Integer:
		int_val = int((value.Value).(int))
	case gosnmp.Counter32:
		int_val = int((value.Value).(uint64))
	case gosnmp.Gauge32:
		int_val = int((value.Value).(uint64))
	case gosnmp.TimeTicks:
		int_val = int((value.Value).(int))
	default:
		log.Println("Value of oid", oid, "Is no number!")
		return 0, reply.Variables[0].Name
	}
	return int_val, reply.Variables[0].Name

}
func (snmp * _SNMP_Helper) get_next_oid_string(oid string) (string, string) {
	reply, err := snmp.snmp.GetNext(oid)
	if (err != nil) {
		log.Fatalln(reply)
	}
	value := reply.Variables[0]

	if (value.Type != gosnmp.OctetString) {
		return "", value.Name
	}


	return (value.Value).(string), value.Name

}