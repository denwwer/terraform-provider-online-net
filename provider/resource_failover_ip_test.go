package provider

import (
	"testing"

	"github.com/src-d/terraform-provider-online-net/online"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	onlineClientMock.On("EditFailoverIP", "127.0.0.1", "8.8.8.8").Return(nil)
	onlineClientMock.On("EditFailoverIP", "127.0.0.1", "false").Return(nil)
	onlineClientMock.On("GenerateMACFailoverIP", "127.0.0.1", "kvm").Return("ma:ac:te:st", nil)
	onlineClientMock.On("DeleteMACFailoverIP", "127.0.0.1").Return(nil)
	onlineClientMock.On("Server", 123).Return(&online.Server{
		IP: []*online.Interface{
			&online.Interface{
				Address: "8.8.8.8",
				Type:    online.Public,
			},
		},
	}, nil)
}

func TestResourceFailoverIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				resource "online_failover_ip" "test" {
	 				"ip" = "127.0.0.1"
					"destination_server_ip" = "8.8.8.8"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", "127.0.0.1"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
				resource "online_failover_ip" "test" {
					 "ip" = "127.0.0.1"
					 "destination_server_ip" = "8.8.8.8"
					 "generate_mac" = true
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", "127.0.0.1"),
					resource.TestCheckResourceAttr("online_failover_ip.test", "mac", "ma:ac:te:st"),
				),
			},
		},
	})
}
