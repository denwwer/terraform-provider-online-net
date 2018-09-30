package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceRPN() *schema.Resource {
	return &schema.Resource{
		Create: resourceRPNCreate,
		Update: resourceRPNUpdate,
		Read:   resourceRPNRead,
		Delete: resourceRPNDelete,

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

func resourceRPNCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpn != nil {
		return fmt.Errorf("RPN already exists")
	}

	rpn = &online.RPNv2{
		Name: name,
		Type: online.RPNv2Type(d.Get("type").(string)),
	}

	return setRPN(c, rpn, d)
}

func resourceRPNUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpn == nil {
		return fmt.Errorf("missing RPNv2 group: %q", name)
	}

	rpn = &online.RPNv2{
		ID:   rpn.ID,
		Name: name,
		Type: online.RPNv2Type(d.Get("type").(string)),
	}

	return setRPN(c, rpn, d)
}

func setRPN(c online.Client, rpn *online.RPNv2, d *schema.ResourceData) error {
	server_ids := d.Get("server_ids").([]interface{})
	if len(server_ids) == 0 {
		return fmt.Errorf("server_ids cannot be empty")
	}

	for _, id := range server_ids {
		m := &online.Member{}
		m.Linked.ID = id.(int)
		m.VLAN = d.Get("vlan").(int)
		rpn.Members = append(rpn.Members, m)
	}

	if err := c.SetRPNv2(rpn, time.Minute); err != nil {
		return err
	}

	d.SetId(rpn.Name)

	return nil

}

func resourceRPNRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpn == nil {
		return fmt.Errorf("missing RPNv2 group: %q", name)
	}

	d.SetId(rpn.Name)

	return nil
}

func resourceRPNDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(d.Id())
	if err != nil {
		return err
	}

	if rpn == nil {
		return nil
	}

	return c.DeleteRPNv2(rpn.ID, time.Minute)
}
