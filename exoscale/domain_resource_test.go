package exoscale

import (
	"fmt"
	"testing"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDomain(t *testing.T) {
	domain := new(egoscale.DNSDomain)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDNSDomainCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSDomainExists("exoscale_domain.exo", domain),
					testAccCheckDNSDomainAttributes(domain),
					testAccCheckDNSDomainCreateAttributes("exo.exo"),
				),
			},
		},
	})
}

func testAccCheckDNSDomainExists(n string, domain *egoscale.DNSDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No domain ID is set")
		}

		client := GetDNSClient(testAccProvider.Meta())
		d, err := client.GetDomain(rs.Primary.ID)
		if err != nil {
			return err
		}

		domain.Token = d.Token
		domain.Name = d.Name
		domain.ID = d.ID

		return nil
	}
}

func testAccCheckDNSDomainAttributes(domain *egoscale.DNSDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(domain.Token) != 32 {
			return fmt.Errorf("DNS Domain: token length doesn't match")
		}

		return nil
	}
}

func testAccCheckDNSDomainCreateAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "exoscale_domain" {
				continue
			}

			if rs.Primary.ID != name {
				continue
			}

			if rs.Primary.Attributes["token"] == "" {
				return fmt.Errorf("DNS Domain: expected token to be set")
			}

			return nil
		}

		return fmt.Errorf("Could not find domain %s", name)
	}
}

func testAccCheckDNSDomainDestroy(s *terraform.State) error {
	client := GetDNSClient(testAccProvider.Meta())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "exoscale_domain" {
			continue
		}

		d, err := client.GetDomain(rs.Primary.Attributes["name"])
		if err != nil {
			return err
		}
		if d == nil {
			return nil
		}
	}
	return fmt.Errorf("DNS Domain: still exists")
}

var testAccDNSDomainCreate = `
resource "exoscale_domain" "exo" {
  name = "exo.exo"
}
`
