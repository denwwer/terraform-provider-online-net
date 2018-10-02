package provider

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Update: resourceServerCreate,
		Read:   resourceServerRead,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"server_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"public_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     resourceInterface(),
			},
			"private_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     resourceInterface(),
			},
		},
	}
}

func resourceInterface() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hardware address of the device.",
			},
			"dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS server address.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Static IPv4 address.",
			},
		},
	}
}

func resourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	s, err := getServer(c, d)
	if err != nil {
		return err
	}

	if err := updateServerIfNeeded(c, s, d); err != nil {
		return err
	}

	return resourceServerRead(d, meta)
}

func updateServerIfNeeded(c online.Client, s *online.Server, d *schema.ResourceData) error {
	hostname := d.Get("hostname").(string)

	var changed bool
	if s.Hostname != hostname {
		changed = true
		s.Hostname = hostname
	}

	publicDNS := d.Get("public_interface.0.dns").(string)
	ip := s.InterfaceByType(online.Public)
	if ip != nil && publicDNS != "" && ip.Reverse != publicDNS {
		changed = true
		ip.Reverse = publicDNS
	}

	if !changed {
		return nil
	}

	return c.SetServer(s)
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)
	s, err := getServer(client, d)
	if err != nil {
		return err
	}

	applyServer(s, d)
	return nil
}

func getServer(c online.Client, d *schema.ResourceData) (*online.Server, error) {
	id := d.Get("server_id").(int)
	d.SetId(strconv.Itoa(id))

	return c.Server(id)
}

func applyServer(s *online.Server, d *schema.ResourceData) {
	var public, private []map[string]interface{}

	for _, iface := range s.IP {
		switch iface.Type {
		case online.Public:
			public = append(public, map[string]interface{}{
				"mac":     strings.ToLower(iface.MAC),
				"dns":     iface.Reverse,
				"address": iface.Address,
			})

		case online.Private:
			private = append(private, map[string]interface{}{
				"mac":     strings.ToLower(iface.MAC),
				"dns":     iface.Reverse,
				"address": iface.Address,
			})
		}
	}

	d.Set("public_interface", public)
	d.Set("private_interface", private)
}
