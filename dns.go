package google

import "fmt"

type BaseDNSZone interface {
	GetName() string

	AddRecords(func()) BaseDNSZone
	DeleteRecords(func()) BaseDNSZone
}

type GoogleDNSZone interface {
	BaseDNSZone
	ListRecords() GoogleDNSZone
}

type BaseDNSRecord interface {
	GetName() string
	GetTTL() int
	GetType() string
	GetRData() []string

	TTL(int) BaseDNSRecord
	Type(string) BaseDNSRecord
	RData(...string) BaseDNSRecord
}

func DNSZone(name string) GoogleDNSZone {
	return &googleDNSZone{
		Name: name,
	}
}

type googleDNSZone struct {
	Name string
}

func (zone *googleDNSZone) GetName() string {
	return zone.Name
}

func (zone *googleDNSZone) AddRecords(fn func()) BaseDNSZone {
	fn()
	return zone
}

func (zone *googleDNSZone) DeleteRecords(fn func()) BaseDNSZone {
	fn()
	return zone
}

func (zone *googleDNSZone) ListRecords() GoogleDNSZone {

	fmt.Printf("Zones!\n")
	return zone
}
