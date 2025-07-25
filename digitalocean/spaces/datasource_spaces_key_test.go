package spaces_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanSpacesKey_basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resources, both := testAccDataSourceDigitalOceanSpacesKeyConfig_basic(name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories:         acceptance.TestAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanSpacesKeyDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resources,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_key.key", "name", name),
					resource.TestCheckResourceAttr("digitalocean_spaces_key.key", "grant.0.bucket", "my-bucket"),
					resource.TestCheckResourceAttr("digitalocean_spaces_key.key", "grant.0.permission", "read"),
				),
			},
			{
				Config: both,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_spaces_key.key", "name", name),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_key.key", "grant.0.bucket", "my-bucket"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_key.key", "grant.0.permission", "read"),
				),
			},
		},
	})
}

func testAccDataSourceDigitalOceanSpacesKeyConfig_basic(name string) (string, string) {
	resources := fmt.Sprintf(`
resource "digitalocean_spaces_key" "key" {
  name = "%s"
  grant {
    bucket     = "my-bucket"
    permission = "read"
  }
  grant {
    bucket     = "my-bucket2"
    permission = "readwrite"
  }
}
`, name)

	both := fmt.Sprintf(`%s
data "digitalocean_spaces_key" "key" {
  name = "%s"
}
`, resources, name)

	return resources, both
}

func testAccCheckDigitalOceanSpacesKeyDestroy(s *terraform.State) error {
	return testAccCheckDigitalOceanSpacesKeyDestroyWithProvider(s, acceptance.TestAccProvider)
}

func testAccCheckDigitalOceanSpacesKeyDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_spaces_key" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is not set")
		}

		client := provider.Meta().(*config.CombinedConfig).GodoClient()

		key, resp, err := client.SpacesKeys.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("Error listing Spaces keys: %s", err)
		}
		if key.AccessKey == rs.Primary.ID {
			return fmt.Errorf("Key still exists")
		}
		return nil
	}
	return nil
}
