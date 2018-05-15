package exoscale

import (
	"context"
	"fmt"
	"net"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/terraform/helper/schema"
)

func secondaryIPResource() *schema.Resource {
	return &schema.Resource{
		Create: createSecondaryIP,
		Exists: existsSecondaryIP,
		Read:   readSecondaryIP,
		Delete: deleteSecondaryIP,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"compute_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Elastic IP address",
				ValidateFunc: ValidateIPv4String,
			},
			"nic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createSecondaryIP(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetComputeClient(meta)

	virtualMachineID := d.Get("compute_id").(string)

	resp, err := client.RequestWithContext(ctx, &egoscale.ListNics{
		VirtualMachineID: virtualMachineID,
	})
	if err != nil {
		return err
	}

	nics := resp.(*egoscale.ListNicsResponse)
	if nics.Count == 0 {
		return fmt.Errorf("The VM has no NIC %v", virtualMachineID)
	}

	// XXX Fragile
	nic := nics.Nic[0]
	resp, err = client.RequestWithContext(ctx, &egoscale.AddIPToNic{
		NicID:     nic.ID,
		IPAddress: net.ParseIP(d.Get("ip_address").(string)),
	})
	if err != nil {
		return err
	}

	secondaryIP := resp.(*egoscale.NicSecondaryIP)

	d.SetId(secondaryIP.ID)
	// XXX this is fragile
	d.Set("nic_id", nic.ID)
	return nil
}

func existsSecondaryIP(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	client := GetComputeClient(meta)

	nicID := d.Get("nic_id").(string)
	virtualMachineID := d.Get("compute_id").(string)
	resp, err := client.RequestWithContext(ctx, &egoscale.ListNics{
		NicID:            nicID,
		VirtualMachineID: virtualMachineID,
	})

	if err != nil {
		// XXX Check the root cause of that error to tell
		//     using pkg/errors.
		return err != nil, err
	}

	nics := resp.(*egoscale.ListNicsResponse)
	if nics.Count == 0 {
		return false, nil
	}

	return true, nil
}

func readSecondaryIP(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	client := GetComputeClient(meta)

	nicID := d.Get("nic_id").(string)
	virtualMachineID := d.Get("compute_id").(string)
	resp, err := client.RequestWithContext(ctx, &egoscale.ListNics{
		NicID:            nicID,
		VirtualMachineID: virtualMachineID,
	})

	if err != nil {
		return err
	}

	nics := resp.(*egoscale.ListNicsResponse)
	if len(nics.Nic) == 0 {
		// No nics, means the VM is gone.
		d.SetId("")
		return nil
	}

	nic := nics.Nic[0]
	for _, ip := range nic.SecondaryIP {
		if d.Id() == "" || ip.ID == d.Id() {
			err := applySecondaryIP(d, ip)
			if err != nil {
				return err
			}
			// fix fix
			d.Set("nic_id", nic.ID)
			d.Set("network_id", nic.NetworkID)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func deleteSecondaryIP(d *schema.ResourceData, meta interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()

	client := GetComputeClient(meta)

	return client.BooleanRequestWithContext(ctx, &egoscale.RemoveIPFromNic{
		ID: d.Id(),
	})
}

func applySecondaryIP(d *schema.ResourceData, secondaryIP egoscale.NicSecondaryIP) error {
	d.SetId(secondaryIP.ID)
	if secondaryIP.IPAddress != nil {
		d.Set("ip_address", secondaryIP.IPAddress.String())
	} else {
		d.Set("ip_address", "")
	}
	d.Set("network_id", secondaryIP.NetworkID)
	d.Set("nic_id", secondaryIP.NicID)

	return nil
}
