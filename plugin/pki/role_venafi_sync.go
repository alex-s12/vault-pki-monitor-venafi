package pki

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) roleVenafiSync(ctx context.Context, req *logical.Request) (err error) {

	//Get role list with role sync param
	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return
	}

	if len(roles) == 0 {
		return
	}
	//name
	//sync zone name\id
	//sync endpoint


	for _, roleName := range roles {
		//	Read previous role parameters
		entry := roleEntry{
			AllowLocalhost:   true,
			AllowedDomains:   []string{"venafi.com"},
			AllowBareDomains: true,
			AllowSubdomains:  true,
			AllowGlobDomains: true,
			AllowAnyName:     true,
			EnforceHostnames: true,

			OU:            []string{"DevOps-old"},
			Organization:  []string{"Venafi-old"},
			Country:       []string{"US-old"},
			Locality:      []string{"Salt Lake-old"},
			Province:      []string{"Venafi-old"},
			StreetAddress: []string{"Venafi-old"},
			PostalCode:    []string{"122333344-old"},
		}
		//  Rewrite entry
		entry.OU = entryRewrite.OU
		entry.Organization = entryRewrite.Organization
		entry.Country = entryRewrite.Country
		entry.Locality = entryRewrite.Locality
		entry.Province = entryRewrite.Province
		entry.StreetAddress = entryRewrite.StreetAddress
		entry.PostalCode = entryRewrite.PostalCode

		// Put new entry
		// Store it
		jsonEntry, err := logical.StorageEntryJSON("role/"+roleName, entryRewrite)
		if err != nil {
			return
		}
		if err := req.Storage.Put(ctx, jsonEntry); err != nil {
			return
		}
	}

	return
}

func (b *backend) getVenafiPoliciyParams(ctx context.Context, req *logical.Request, roleName string) (entry roleEntry, err error) {
	//Get role params from TPP\Cloud
	cl, err := b.ClientVenafi(ctx, req.Storage, roleName, "role")
	if err != nil {
		return entry, fmt.Errorf("could not create venafi client: %s", err)
	}

	zone, err := cl.ReadZoneConfiguration()
	if err != nil {
		return
	}
	entry = roleEntry{
		OU:            zone.SubjectOURegexes,
		Organization:  []string{"Venafi"},
		Country:       []string{zone.Country},
		Locality:      []string{"Salt Lake"},
		Province:      []string{"Venafi"},
		StreetAddress: []string{"Venafi"},
		PostalCode:    []string{"122333344"},
	}
}

func (b *backend) getPKIRoleEntry(ctx context.Context, req *logical.Request, roleName string) (entry roleEntry, err error) {
	return entry, nil
}