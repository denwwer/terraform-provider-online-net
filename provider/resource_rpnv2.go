package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceRPNv2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRPNv2Create,
		Update: resourceRPNv2Update,
		Read:   resourceRPNv2Read,
		Delete: resourceRPNv2Delete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  online.Standard,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"server_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceRPNv2Create(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	c := meta.(online.Client)
	rpnv2, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpnv2 != nil {
		return fmt.Errorf("RPN already exists")
	}

	rpnv2 = &online.RPNv2{
		Name: name,
		Type: online.RPNv2Type(d.Get("type").(string)),
	}

	return setRPNv2(c, rpnv2, d)
}

func resourceRPNv2Update(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	c := meta.(online.Client)
	rpnv2, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpnv2 == nil {
		return fmt.Errorf("missing RPNv2 group: %q", name)
	}

	rpnv2 = &online.RPNv2{
		ID:   rpnv2.ID,
		Name: name,
		Type: online.RPNv2Type(d.Get("type").(string)),
	}

	return setRPNv2(c, rpnv2, d)
}

func setRPNv2(c online.Client, rpnv2 *online.RPNv2, d *schema.ResourceData) error {
	server_ids := d.Get("server_ids").([]interface{})
	if len(server_ids) == 0 {
		return fmt.Errorf("server_ids cannot be empty")
	}

	for _, id := range server_ids {
		m := &online.Member{}
		m.Linked.ID = id.(int)
		m.VLAN = d.Get("vlan").(int)
		rpnv2.Members = append(rpnv2.Members, m)
	}

	if err := c.SetRPNv2(rpnv2, time.Minute); err != nil {
		return err
	}

	d.SetId(rpnv2.Name)

	return nil

}

func resourceRPNv2Read(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	c := meta.(online.Client)
	rpnv2, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpnv2 == nil {
		return fmt.Errorf("missing RPNv2 group: %q", name)
	}

	d.SetId(rpnv2.Name)

	return nil
}

func resourceRPNv2Delete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	c := meta.(online.Client)
	rpnv2, err := c.RPNv2ByName(d.Id())
	if err != nil {
		return err
	}

	if rpnv2 == nil {
		return nil
	}

	return c.DeleteRPNv2(rpnv2.ID, time.Minute)
}