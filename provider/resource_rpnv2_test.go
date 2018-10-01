package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceRPNv2(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{{
			ImportStateVerify: true,
			Config: `
				resource "online_rpnv2" "test" {
					name = "rpn-test"
					vlan = "2242"
					server_ids = ["${online_server.test.server_id}"]
				}

				resource "online_server" "test" {
					server_id = 105770
					hostname = "mvp"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2242"),
				resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "1"),
				resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", "105770"),
			),
		}},
	})
}
