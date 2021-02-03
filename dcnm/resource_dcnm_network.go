package dcnm

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ciscoecosystem/dcnm-go-client/client"
	"github.com/ciscoecosystem/dcnm-go-client/container"
	"github.com/ciscoecosystem/dcnm-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDCNMNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceDCNMNetworkCreate,
		Update: resourceDCNMNetworkUpdate,
		Read:   resourceDCNMNetworkRead,
		Delete: resourceDCNMNetworkDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDCNMNetworkImporter,
		},

		Schema: map[string]*schema.Schema{
			"fabric_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_Network_Universal",
			},

			"extension_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Default_Network_Extension_Universal",
			},

			"vrf_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NA",
			},

			"l2_only_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"vlan_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv4_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ipv6_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"secondary_gw_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"secondary_gw_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"arp_supp_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"ir_enable_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"mcast_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_1": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_2": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"dhcp_vrf": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"loopback_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"trm_enable_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"rt_both_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"l3_gateway_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"service_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"deploy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"attachments": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"serial_number": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"vlan_id": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"dot1_qvlan": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},

						"attach": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"switch_ports": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"untagged": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},

						"free_from_config": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"extension_values": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"instance_values": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getRemoteNetwork(client *client.Client, fabric, name string) (*container.Container, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s", fabric, name)

	cont, err := client.GetviaURL(durl)
	if err != nil {
		return nil, err
	}

	return cont, nil
}

func setNetworkAttributes(d *schema.ResourceData, cont *container.Container) *schema.ResourceData {
	d.Set("fabric_name", stripQuotes(cont.S("fabric").String()))
	d.Set("name", stripQuotes(cont.S("networkName").String()))
	d.Set("network_id", stripQuotes(cont.S("networkId").String()))
	d.Set("template", stripQuotes(cont.S("networkTemplate").String()))
	d.Set("extension_template", stripQuotes(cont.S("networkExtensionTemplate").String()))
	d.Set("vrf_name", stripQuotes(cont.S("vrf").String()))

	if cont.Exists("displayName") {
		d.Set("display_name", stripQuotes(cont.S("displayName").String()))
	}
	if cont.Exists("serviceNetworkTemplate") && stripQuotes(cont.S("serviceNetworkTemplate").String()) != "null" {
		d.Set("service_template", stripQuotes(cont.S("serviceNetworkTemplate").String()))
	}
	if cont.Exists("source") && stripQuotes(cont.S("source").String()) != "null" {
		d.Set("source", stripQuotes(cont.S("source").String()))
	}

	cont, err := cleanJsonString(stripQuotes(cont.S("networkTemplateConfig").String()))
	if err == nil {
		if cont.Exists("isLayer2Only") && stripQuotes(cont.S("isLayer2Only").String()) != "" {
			if l2, err := strconv.ParseBool(stripQuotes(cont.S("isLayer2Only").String())); err == nil {
				d.Set("l2_only_flag", l2)
			}
		} else {
			d.Set("l2_only_flag", false)
		}
		if cont.Exists("vlanId") && stripQuotes(cont.S("vlanId").String()) != "" {
			if vlan, err := strconv.Atoi(stripQuotes(cont.S("vlanId").String())); err == nil {
				d.Set("vlan_id", vlan)
			}
		}
		if cont.Exists("vlanName") {
			d.Set("vlan_name", stripQuotes(cont.S("vlanName").String()))
		}
		if cont.Exists("gatewayIpAddress") {
			d.Set("ipv4_gateway", stripQuotes(cont.S("gatewayIpAddress").String()))
		}
		if cont.Exists("gatewayIpV6Address") {
			d.Set("ipv6_gateway", stripQuotes(cont.S("gatewayIpV6Address").String()))
		}
		if cont.Exists("intfDescription") {
			d.Set("description", stripQuotes(cont.S("intfDescription").String()))
		}
		if cont.Exists("mtu") && stripQuotes(cont.S("mtu").String()) != "" {
			if mtu, err := strconv.Atoi(stripQuotes(cont.S("mtu").String())); err == nil {
				d.Set("mtu", mtu)
			}
		}
		if cont.Exists("secondaryGW1") {
			d.Set("secondary_gw_1", stripQuotes(cont.S("secondaryGW1").String()))
		}
		if cont.Exists("secondaryGW2") {
			d.Set("secondary_gw_2", stripQuotes(cont.S("secondaryGW2").String()))
		}
		if cont.Exists("suppressArp") && stripQuotes(cont.S("suppressArp").String()) != "" {
			if arp, err := strconv.ParseBool(stripQuotes(cont.S("suppressArp").String())); err == nil {
				d.Set("arp_supp_flag", arp)
			}
		} else {
			d.Set("arp_supp_flag", false)
		}
		if cont.Exists("enableIR") && stripQuotes(cont.S("enableIR").String()) != "" {
			if ir, err := strconv.ParseBool(stripQuotes(cont.S("enableIR").String())); err == nil {
				d.Set("ir_enable_flag", ir)
			}
		} else {
			d.Set("ir_enable_flag", false)
		}
		if cont.Exists("mcastGroup") {
			d.Set("mcast_group", stripQuotes(cont.S("mcastGroup").String()))
		}
		if cont.Exists("dhcpServerAddr1") {
			d.Set("dhcp_1", stripQuotes(cont.S("dhcpServerAddr1").String()))
		}
		if cont.Exists("dhcpServerAddr2") {
			d.Set("dhcp_2", stripQuotes(cont.S("dhcpServerAddr2").String()))
		}
		if cont.Exists("vrfDhcp") {
			d.Set("dhcp_vrf", stripQuotes(cont.S("vrfDhcp").String()))
		}
		if cont.Exists("loopbackId") && stripQuotes(cont.S("loopbackId").String()) != "" {
			if loopback, err := strconv.Atoi(stripQuotes(cont.S("loopbackId").String())); err == nil {
				d.Set("loopback_id", loopback)
			}
		}
		if cont.Exists("tag") {
			d.Set("tag", stripQuotes(cont.S("tag").String()))
		}
		if cont.Exists("trmEnabled") && stripQuotes(cont.S("trmEnabled").String()) != "" {
			if trm, err := strconv.ParseBool(stripQuotes(cont.S("trmEnabled").String())); err == nil {
				d.Set("trm_enable_flag", trm)
			}
		} else {
			d.Set("trm_enable_flag", false)
		}
		if cont.Exists("rtBothAuto") && stripQuotes(cont.S("rtBothAuto").String()) != "" {
			if rt, err := strconv.ParseBool(stripQuotes(cont.S("rtBothAuto").String())); err == nil {
				d.Set("rt_both_flag", rt)
			}
		} else {
			d.Set("rt_both_flag", false)
		}
		if cont.Exists("enableL3OnBorder") && stripQuotes(cont.S("enableL3OnBorder").String()) != "" {
			if l3, err := strconv.ParseBool(stripQuotes(cont.S("enableL3OnBorder").String())); err == nil {
				d.Set("l3_gateway_flag", l3)
			}
		} else {
			d.Set("l3_gateway_flag", false)
		}
	}

	d.SetId(stripQuotes(cont.S("networkName").String()))
	return d
}

func resourceDCNMNetworkImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Println("[DEBUG] Begining Importer ", d.Id())

	dcnmClient := m.(*client.Client)
	importInfo := strings.Split(d.Id(), ":")
	if len(importInfo) != 2 {
		return nil, fmt.Errorf("not getting enough arguments for the import operation")
	}
	fabricName := importInfo[0]
	network := importInfo[1]

	cont, err := getRemoteNetwork(dcnmClient, fabricName, network)
	if err != nil {
		return nil, err
	}

	stateImport := setNetworkAttributes(d, cont)

	deployed, err := checkNetworkDeploy(dcnmClient, fabricName, network)
	if err != nil {
		d.Set("deploy", false)
		return nil, err
	}
	d.Set("deploy", deployed)

	attachments, err := getNetworkAttachmentList(dcnmClient, fabricName, network)
	if err == nil {
		d.Set("attachments", attachments)
	} else {
		d.Set("attachments", make([]interface{}, 0, 1))
	}

	log.Println("[DEBUG] End of Importer ", d.Id())
	return []*schema.ResourceData{stateImport}, nil
}

func resourceDCNMNetworkCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Create method ")

	dcnmClient := m.(*client.Client)

	name := d.Get("name").(string)
	fabricName := d.Get("fabric_name").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); !ok {
			return fmt.Errorf("attachments must be configured if deploy=true")
		}
	}

	cont, err := dcnmClient.GetSegID(fmt.Sprintf("/rest/managed-pool/fabrics/%s/segments/ids", fabricName))
	if err != nil {
		return err
	}
	segID := cont.S("segmentId").String()

	network := models.Network{}
	networkProfile := models.NetworkProfileConfig{}

	network.Fabric = fabricName
	network.Name = name
	if display, ok := d.GetOk("display_name"); ok {
		network.DisplayName = display.(string)
	} else {
		network.DisplayName = name
	}
	network.NetworkId = segID
	network.Template = d.Get("template").(string)
	network.ExtensionTemplate = d.Get("extension_template").(string)
	network.VRF = d.Get("vrf_name").(string)

	if svcTemplate, ok := d.GetOk("service_template"); ok {
		network.ServiceNetworkTemplate = svcTemplate.(string)
	}
	if src, ok := d.GetOk("source"); ok {
		network.Source = src.(string)
	}

	if network.VRF == "NA" {
		networkProfile.L2OnlyFlag = true
	} else {
		networkProfile.L2OnlyFlag = false
	}
	if ipv4, ok := d.GetOk("ipv4_gateway"); ok {
		networkProfile.GatewayIpv4 = ipv4.(string)
	}
	if ipv6, ok := d.GetOk("ipv6_gateway"); ok {
		networkProfile.GatewayIPv6 = ipv6.(string)
	}
	if vlan, ok := d.GetOk("vlan_id"); ok {
		networkProfile.Vlan = vlan.(int)
	} else {
		durl := fmt.Sprintf("/rest/resource-manager/vlan/%s?vlanUsageType=TOP_DOWN_NETWORK_VLAN", fabricName)
		cont, err := dcnmClient.GetviaURL(durl)
		if err != nil {
			return err
		}
		vlan, err := strconv.Atoi(cont.String())
		if err == nil {
			networkProfile.Vlan = vlan
		}
	}
	if vlanName, ok := d.GetOk("vlan_name"); ok {
		networkProfile.VlanName = vlanName.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		networkProfile.Description = desc.(string)
	}
	if mtu, ok := d.GetOk("mtu"); ok {
		networkProfile.MTU = mtu.(int)
	}
	if secgw1, ok := d.GetOk("secondary_gw_1"); ok {
		networkProfile.SecondaryGate1 = secgw1.(string)
	}
	if secgw2, ok := d.GetOk("secondary_gw_2"); ok {
		networkProfile.SecondaryGate2 = secgw2.(string)
	}
	if arp, ok := d.GetOk("arp_supp_flag"); ok {
		networkProfile.ARPSuppFlag = arp.(bool)
	}
	if ir, ok := d.GetOk("ir_enable_flag"); ok {
		networkProfile.IRFlag = ir.(bool)
	}
	if mcast, ok := d.GetOk("mcast_group"); ok {
		networkProfile.McastGroup = mcast.(string)
	}
	if dhcp1, ok := d.GetOk("dhcp_1"); ok {
		networkProfile.DHCPServer1 = dhcp1.(string)
	}
	if dhcp2, ok := d.GetOk("dhcp_2"); ok {
		networkProfile.DHCPServer2 = dhcp2.(string)
	}
	if dhcpvrf, ok := d.GetOk("dhcp_vrf"); ok {
		networkProfile.DHCPServerVRF = dhcpvrf.(string)
	}
	if loopback, ok := d.GetOk("loopback_id"); ok {
		networkProfile.LookbackID = loopback.(int)
	}
	if tag, ok := d.GetOk("tag"); ok {
		networkProfile.Tag = tag.(string)
	}
	if trm, ok := d.GetOk("trm_enable_flag"); ok {
		networkProfile.TRMEnable = trm.(bool)
	}
	if rtBoth, ok := d.GetOk("rt_both_flag"); ok {
		networkProfile.RTBothFlag = rtBoth.(bool)
	}
	if l3enable, ok := d.GetOk("l3_gateway_flag"); ok {
		networkProfile.L3GatewayEnable = l3enable.(bool)
	}
	networkProfile.NetworkName = name
	networkProfile.SegmentID = segID

	configStr, err := json.Marshal(networkProfile)
	if err != nil {
		return err
	}
	network.Config = string(configStr)

	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks", fabricName)
	_, err = dcnmClient.Save(durl, &network)
	if err != nil {
		return err
	}
	d.SetId(name)

	//Network Deployment
	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); ok {
			attachList := make([]map[string]interface{}, 0, 1)
			for _, val := range d.Get("attachments").(*schema.Set).List() {
				attachment := val.(map[string]interface{})

				attachMap := make(map[string]interface{})

				attachMap["fabric"] = network.Fabric
				attachMap["networkName"] = network.Name
				attachMap["deployment"] = attachment["attach"].(bool)
				attachMap["serialNumber"] = attachment["serial_number"].(string)

				if attachment["vlan_id"].(int) != 0 {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				} else {
					attachMap["vlan"] = networkProfile.Vlan
				}
				if attachment["switch_ports"] != nil {
					attachMap["switchPorts"] = listToString(attachment["switch_ports"])
				}

				if attachment["dot1_qvlan"] != nil {
					attachMap["dot1QVlan"] = attachment["dot1_qvlan"].(int)
				}

				if attachment["untagged"] != nil {
					attachMap["untagged"] = attachment["untagged"].(bool)
				}

				if attachment["free_form_config"] != nil {
					attachMap["freeformConfig"] = attachment["free_form_config"].(string)
				}

				if attachment["extension_values"] != nil {
					attachMap["extensionValues"] = attachment["extension_values"].(string)
				}

				if attachment["instanceValues"] != nil {
					attachMap["instanceValues"] = attachment["instance_values"].(string)
				}

				attachList = append(attachList, attachMap)
			}

			networkAttach := models.NewNetworkAttachment(network.Name, attachList)
			durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/attachments", network.Fabric)
			cont, err := dcnmClient.SaveForAttachment(durl, networkAttach)
			if err != nil {
				d.Set("deploy", false)
				d.Set("attachments", make([]interface{}, 0, 1))
				return fmt.Errorf("Network record is created but not deployed yet. Error while attachment : %s", err)
			}

			// Network Deployment
			for _, v := range cont.Data().(map[string]interface{}) {
				if v != "SUCCESS" && v != "SUCCESS Peer attach Reponse :  SUCCESS" {
					return fmt.Errorf("Network record is created but not deployed yet. Error while attachment : %s", v)
				}
			}

			durl = fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s/deploy", network.Fabric, network.Name)
			_, err = dcnmClient.SaveAndDeploy(durl)
			if err != nil {
				d.Set("deploy", false)
			}
			time.Sleep(10 * time.Second)

		} else {
			d.Set("deploy", false)
			d.Set("attachments", make([]interface{}, 0, 1))
			return fmt.Errorf("Network record is created but not deployed yet. Either make deploy=false or provide attachments")
		}
	}

	log.Println("[DEBUG] End of Create method ", d.Id())
	return resourceDCNMNetworkRead(d, m)
}

func resourceDCNMNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Update method ", d.Id())

	dcnmClient := m.(*client.Client)

	name := d.Get("name").(string)
	fabricName := d.Get("fabric_name").(string)
	segID := d.Get("network_id").(string)

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); !ok {
			return fmt.Errorf("attachments must be configured if deploy=true")
		}
	}

	network := models.Network{}
	networkProfile := models.NetworkProfileConfig{}

	network.Fabric = fabricName
	network.Name = name
	if display, ok := d.GetOk("display_name"); ok {
		network.DisplayName = display.(string)
	} else {
		network.DisplayName = name
	}
	network.NetworkId = segID
	network.Template = d.Get("template").(string)
	network.ExtensionTemplate = d.Get("extension_template").(string)
	network.VRF = d.Get("vrf_name").(string)

	if svcTemplate, ok := d.GetOk("service_template"); ok {
		network.ServiceNetworkTemplate = svcTemplate.(string)
	}
	if src, ok := d.GetOk("source"); ok {
		network.Source = src.(string)
	}

	if network.VRF == "NA" {
		networkProfile.L2OnlyFlag = true
	} else {
		networkProfile.L2OnlyFlag = false
	}
	if ipv4, ok := d.GetOk("ipv4_gateway"); ok {
		networkProfile.GatewayIpv4 = ipv4.(string)
	}
	if ipv6, ok := d.GetOk("ipv6_gateway"); ok {
		networkProfile.GatewayIPv6 = ipv6.(string)
	}
	if vlan, ok := d.GetOk("vlan_id"); ok {
		networkProfile.Vlan = vlan.(int)
	}
	if vlanName, ok := d.GetOk("vlan_name"); ok {
		networkProfile.VlanName = vlanName.(string)
	}
	if desc, ok := d.GetOk("description"); ok {
		networkProfile.Description = desc.(string)
	}
	if mtu, ok := d.GetOk("mtu"); ok {
		networkProfile.MTU = mtu.(int)
	}
	if secgw1, ok := d.GetOk("secondary_gw_1"); ok {
		networkProfile.SecondaryGate1 = secgw1.(string)
	}
	if secgw2, ok := d.GetOk("secondary_gw_2"); ok {
		networkProfile.SecondaryGate2 = secgw2.(string)
	}
	if arp, ok := d.GetOk("arp_supp_flag"); ok {
		networkProfile.ARPSuppFlag = arp.(bool)
	}
	if ir, ok := d.GetOk("ir_enable_flag"); ok {
		networkProfile.IRFlag = ir.(bool)
	}
	if mcast, ok := d.GetOk("mcast_group"); ok {
		networkProfile.McastGroup = mcast.(string)
	}
	if dhcp1, ok := d.GetOk("dhcp_1"); ok {
		networkProfile.DHCPServer1 = dhcp1.(string)
	}
	if dhcp2, ok := d.GetOk("dhcp_2"); ok {
		networkProfile.DHCPServer2 = dhcp2.(string)
	}
	if dhcpvrf, ok := d.GetOk("dhcp_vrf"); ok {
		networkProfile.DHCPServerVRF = dhcpvrf.(string)
	}
	if loopback, ok := d.GetOk("loopback_id"); ok {
		networkProfile.LookbackID = loopback.(int)
	}
	if tag, ok := d.GetOk("tag"); ok {
		networkProfile.Tag = tag.(string)
	}
	if trm, ok := d.GetOk("trm_enable_flag"); ok {
		networkProfile.TRMEnable = trm.(bool)
	}
	if rtBoth, ok := d.GetOk("rt_both_flag"); ok {
		networkProfile.RTBothFlag = rtBoth.(bool)
	}
	if l3enable, ok := d.GetOk("l3_gateway_flag"); ok {
		networkProfile.L3GatewayEnable = l3enable.(bool)
	}
	networkProfile.NetworkName = name
	networkProfile.SegmentID = segID

	configStr, err := json.Marshal(networkProfile)
	if err != nil {
		return err
	}
	network.Config = string(configStr)

	dn := d.Id()
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s", fabricName, dn)
	_, err = dcnmClient.Update(durl, &network)
	if err != nil {
		return err
	}
	d.SetId(name)

	//Network Deployment
	if d.HasChange("deploy") && d.Get("deploy").(bool) == false {
		return fmt.Errorf("Deployed network can not be undeployed")
	}

	if deploy, ok := d.GetOk("deploy"); ok && deploy.(bool) == true {
		if _, ok := d.GetOk("attachments"); ok {
			attachList := make([]map[string]interface{}, 0, 1)
			for _, val := range d.Get("attachments").(*schema.Set).List() {
				attachment := val.(map[string]interface{})

				attachMap := make(map[string]interface{})

				attachMap["fabric"] = network.Fabric
				attachMap["networkName"] = network.Name
				attachMap["deployment"] = attachment["attach"].(bool)
				attachMap["serialNumber"] = attachment["serial_number"].(string)

				if attachment["vlan_id"].(int) != 0 {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				} else {
					attachMap["vlan"] = networkProfile.Vlan
				}

				oldAttachments, newAttachments := d.GetChange("attachments")
				sPorts, dsPorts := findDiffForPorts(oldAttachments.(*schema.Set).List(), newAttachments.(*schema.Set).List(), attachMap["serialNumber"].(string))
				if len(sPorts.([]interface{})) > 0 {
					attachMap["switchPorts"] = listToString(sPorts)
				} else {
					attachMap["switchPorts"] = ""
				}

				if len(dsPorts.([]interface{})) > 0 {
					attachMap["detachSwitchPorts"] = listToString(dsPorts)
				} else {
					attachMap["detachSwitchPorts"] = ""
				}

				if attachment["dot1_qvlan"] != nil {
					attachMap["dot1QVlan"] = attachment["dot1_qvlan"].(int)
				}

				if attachment["untagged"] != nil {
					attachMap["untagged"] = attachment["untagged"].(bool)
				}

				if attachment["free_form_config"] != nil {
					attachMap["freeformConfig"] = attachment["free_form_config"].(string)
				}

				if attachment["extension_values"] != nil {
					attachMap["extensionValues"] = attachment["extension_values"].(string)
				}

				if attachment["instanceValues"] != nil {
					attachMap["instanceValues"] = attachment["instance_values"].(string)
				}

				attachList = append(attachList, attachMap)
			}

			networkAttach := models.NewNetworkAttachment(network.Name, attachList)
			durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/attachments", network.Fabric)
			cont, err := dcnmClient.SaveForAttachment(durl, networkAttach)
			if err != nil {
				d.Set("deploy", false)
				d.Set("attachments", make([]interface{}, 0, 1))
				return fmt.Errorf("Network record is updated but not deployed yet. Error while attachment : %s", err)
			}

			// Network Deployment
			for _, v := range cont.Data().(map[string]interface{}) {
				if v != "SUCCESS" && v != "SUCCESS Peer attach Reponse :  SUCCESS" {
					return fmt.Errorf("Network record is updated but not deployed yet. Error while attachment : %s", v)
				}
			}

			durl = fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s/deploy", network.Fabric, network.Name)
			_, err = dcnmClient.SaveAndDeploy(durl)
			if err != nil {
				d.Set("deploy", false)
			}
			time.Sleep(10 * time.Second)

		} else {
			d.Set("deploy", false)
			d.Set("attachments", make([]interface{}, 0, 1))
			return fmt.Errorf("Network record is updated but not deployed yet. Either make deploy=false or provide attachments")
		}
	}

	log.Println("[DEBUG] End of Update method ", d.Id())
	return resourceDCNMNetworkRead(d, m)
}

func resourceDCNMNetworkRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Read method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	cont, err := getRemoteNetwork(dcnmClient, fabricName, dn)
	if err != nil {
		return err
	}

	setNetworkAttributes(d, cont)

	deployed, err := checkNetworkDeploy(dcnmClient, fabricName, dn)
	if err != nil {
		d.Set("deploy", false)
		return err
	}
	d.Set("deploy", deployed)

	if attaches, ok := d.GetOk("attachments"); ok {
		attachGet := make([]interface{}, 0, 1)

		durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s/attachments", fabricName, dn)
		cont, err := dcnmClient.GetviaURL(durl)
		if err != nil {
			return err
		}

		for _, val := range attaches.(*schema.Set).List() {
			attachMap := val.(map[string]interface{})
			serialNum := attachMap["serial_number"].(string)

			attachStatus, ports, vlan, err := getNetworkSwitchAttachStatus(cont, serialNum)
			if err == nil {
				attachMap["attach"] = attachStatus
				if attachMap["vlan_id"].(int) != 0 {
					attachMap["vlan_id"] = vlan
				}
				if ports != nil {
					portsAcc := interfaceToStrList(attachMap["switch_ports"])
					portsGet := ports
					if !compareStrLists(portsAcc, portsGet) {
						attachMap["switch_ports"] = ports
					}
				} else {
					attachMap["switch_ports"] = make([]string, 0, 1)
				}
			}

			attachGet = append(attachGet, attachMap)
		}

		d.Set("attachments", attachGet)
	}

	log.Println("[DEBUG] End of Read method ", d.Id())
	return nil
}

func resourceDCNMNetworkDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Begining Delete method ", d.Id())

	dcnmClient := m.(*client.Client)

	dn := d.Id()
	fabricName := d.Get("fabric_name").(string)

	if d.Get("deploy").(bool) == true {
		if attachments, ok := d.GetOk("attachments"); ok {
			attachList := make([]map[string]interface{}, 0, 1)
			for _, val := range attachments.(*schema.Set).List() {
				attachment := val.(map[string]interface{})

				attachMap := make(map[string]interface{})

				attachMap["fabric"] = fabricName
				attachMap["networkName"] = dn
				attachMap["deployment"] = false
				attachMap["serialNumber"] = attachment["serial_number"].(string)
				if attachment["vlan_id"].(int) == 0 {
					attachMap["vlan"] = d.Get("vlan_id").(int)
				} else {
					attachMap["vlan"] = attachment["vlan_id"].(int)
				}
				attachMap["detachSwitchPorts"] = ""
				attachMap["dot1QVlan"] = 0
				attachMap["extensionValues"] = ""
				attachMap["untagged"] = false
				attachMap["switchPorts"] = ""

				attachList = append(attachList, attachMap)
			}

			networkAttach := models.NewNetworkAttachment(dn, attachList)
			durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/attachments", fabricName)
			cont, err := dcnmClient.SaveForAttachment(durl, networkAttach)
			if err != nil {
				return err
			}

			// Network Deployment
			for _, v := range cont.Data().(map[string]interface{}) {
				if v != "SUCCESS" && v != "SUCCESS Peer attach Reponse :  SUCCESS" {
					return fmt.Errorf("Error while detachment : %s", v)
				}
			}
			durl = fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s/deploy", fabricName, dn)
			_, err = dcnmClient.SaveAndDeploy(durl)
			if err != nil {
				d.Set("deploy", false)
			}
			time.Sleep(10 * time.Second)
		}
	}

	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s", fabricName, dn)
	_, err := dcnmClient.Delete(durl)
	if err != nil {
		return err
	}

	d.SetId("")

	log.Println("[DEBUG] End of Delete method ", d.Id())
	return nil
}

func checkNetworkDeploy(client *client.Client, fabricName, dn string) (bool, error) {
	durl := fmt.Sprintf("/rest/top-down/fabrics/%s/networks/%s/attachments", fabricName, dn)
	cont, err := client.GetviaURL(durl)
	if err != nil {
		return false, err
	}

	flag := false
	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		if stripQuotes(cont.Index(0).S("lanAttachState").String()) == "DEPLOYED" {
			flag = true
			break
		}
	}
	return flag, nil
}

func getNetworkSwitchAttachStatus(cont *container.Container, serial string) (bool, []string, int, error) {
	for i := 0; i < len(cont.Data().([]interface{})); i++ {
		if stripQuotes(cont.Index(i).S("switchSerialNo").String()) == serial {
			if stripQuotes(cont.Index(i).S("isLanAttached").String()) == "true" {
				var vlanAct int
				if stripQuotes(cont.Index(i).S("vlanId").String()) != "null" {
					vlanAct = int((cont.Index(i).S("vlanId").Data()).(float64))
				}
				if stripQuotes(cont.Index(i).S("portNames").String()) != "null" {
					ports := stringToList(stripQuotes(cont.Index(i).S("portNames").String()))
					return true, ports, vlanAct, nil
				}

				return true, nil, vlanAct, nil
			}

			return false, nil, 0, nil
		}
	}
	return false, nil, 0, nil
}

func findDiffForPorts(oldAttachments interface{}, newAttachments interface{}, serial string) (interface{}, interface{}) {
	oldPorts := make([]string, 0, 1)
	newPorts := make([]string, 0, 1)
	for _, val := range oldAttachments.([]interface{}) {
		attachMap := val.(map[string]interface{})

		if attachMap["serial_number"].(string) == serial {
			oldPorts = interfaceToStrList(attachMap["switch_ports"])
		}
	}

	for _, val := range newAttachments.([]interface{}) {
		attachMap := val.(map[string]interface{})

		if attachMap["serial_number"].(string) == serial {
			newPorts = interfaceToStrList(attachMap["switch_ports"])
		}
	}

	return difference(newPorts, oldPorts), difference(oldPorts, newPorts)
}

func difference(a, b []string) []interface{} {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []interface{}
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}