package components

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/anchore/syft/syft/formats/spdx22json"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	uorspec "github.com/uor-framework/collection-spec/specs-go/v1alpha1"

	"github.com/uor-framework/uor-client-go/nodes/descriptor"
)

// InventoryToProperties updated descriptor properties with inventory information with a given path to the
// descriptor content location on disk.
func InventoryToProperties(inventory sbom.SBOM, path string, properties *descriptor.Properties) error {
	catalog := inventory.Artifacts.PackageCatalog
	pkgs := catalog.PackagesByPath(path)
	pkgLen := len(pkgs)
	if pkgLen == 0 || pkgLen > 1 {
		return fmt.Errorf("incorrect number of components found for %s, expected 1, got %d", path, pkgLen)
	}
	descriptorPkg := pkgs[0]
	var cpes = make([]string, len(descriptorPkg.CPEs))
	for i, c := range descriptorPkg.CPEs {
		cpes[i] = pkg.CPEString(c)
	}
	locations := descriptorPkg.Locations.ToSlice()
	var coordinates = make([]string, len(locations))
	for i, l := range locations {

		coordinates[i] = l.String()
	}

	var additionalMetadata json.RawMessage
	var err error
	if descriptorPkg.Metadata != nil {
		additionalMetadata, err = json.Marshal(descriptorPkg.Metadata)
		if err != nil {
			return err
		}
	}
	component := uorspec.Component{
		ID:        string(descriptorPkg.ID()),
		Name:      descriptorPkg.Name,
		Version:   descriptorPkg.Version,
		FoundBy:   descriptorPkg.FoundBy,
		Locations: coordinates,
		Licenses:  descriptorPkg.Licenses,
		Language:  string(descriptorPkg.Language),
		CPEs:      cpes,
		PURL:      descriptorPkg.PURL,
		AdditionalMetadata: map[string]json.RawMessage{
			string(descriptorPkg.MetadataType): additionalMetadata,
		},
	}

	if properties.Descriptor == nil {
		properties.Descriptor = &uorspec.DescriptorAttributes{}
	}
	properties.Descriptor.Component = component
	return nil
}

// PropertiesToInventory aggregates descriptor information to create an inventory.
func PropertiesToInventory(properties []descriptor.Properties) sbom.SBOM {
	// TODO(jpower432): Must determine how to communicate artifact relationships
	// based on aggregations
	fmt.Println("Not implemented")
	return sbom.SBOM{}
}

// InventoryToSPDXJSON writes an SPDX (22) JSON compliant document
// from an inventory.
func InventoryToSPDXJSON(out io.Writer, inventory sbom.SBOM) error {
	formatter := spdx22json.Format()
	return formatter.Encode(out, inventory)
}
