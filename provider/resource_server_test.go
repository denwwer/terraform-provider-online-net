package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/src-d/terraform-provider-online-net/online"
)

func init() {
	onlineClientMock.On("SetServer", &online.Server{
		Hostname: "mvp",
		IP: []*online.Interface{
			&online.Interface{
				Address: "1.2.3.4",
				MAC:     "aa:bb:cc:dd:ee:ff",
				Reverse: "my.dns.address",
				Type:    online.Public,
			},
			&online.Interface{
				Address: "10.2.3.4",
				MAC:     "00:bb:cc:dd:ee:ff",
				Type:    online.Private,
			},
		},
	}).Return(nil)
	onlineClientMock.On("Server", 105770).Return(&online.Server{
		IP: []*online.Interface{
			&online.Interface{
				Address: "1.2.3.4",
				MAC:     "aa:bb:cc:dd:ee:ff",
				Reverse: "my.dns.address",
				Type:    online.Public,
			},
			&online.Interface{
				Address: "10.2.3.4",
				MAC:     "00:bb:cc:dd:ee:ff",
				Type:    online.Private,
			},
		},
	}, nil)
}

func TestResourceServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testAccProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{{
			ImportStateVerify: true,
			Config: `
				resource "online_server" "test" {
					server_id = 105770
					hostname = "mvp"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "hostname", "mvp"),
				resource.TestCheckResourceAttr("online_server.test", "server_id", "105770"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.0.address", "1.2.3.4"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.0.mac", "aa:bb:cc:dd:ee:ff"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.0.dns", "my.dns.address"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.address", "10.2.3.4"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.mac", "00:bb:cc:dd:ee:ff"),
			),
		}},
	})
}
