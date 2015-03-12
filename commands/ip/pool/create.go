package pool

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/lymingtonprecision/quorra/cli"
	"github.com/lymingtonprecision/quorra/config"
	"github.com/lymingtonprecision/quorra/vsphere/ip/pool"
)

type create struct {
	FlagSet *flag.FlagSet
	Name    *string
	Addr    *string
	Size    *int
	Gateway *string
	Subnet  *string
}

func (cmd *create) CommandLine() string {
	return "[NAME] [START_ADDR/CIDR] [--size|-s SIZE] [--gateway|-gw GATEWAY] [--subnet|-sn SUBNET]"
}

func (cmd *create) Summary() string {
	return "Creates an IP Address Pool"
}

func (cmd *create) Description() string {
	return fmt.Sprintf(
		`%s

Requires a unique NAME for the pool and it's starting address and netmask
(in CIDR format):

    ip pool create MyPool 192.168.1.10/24

    Created IP Pool 'MyPool':

    Address Range: 192.168.1.1 - 192.168.1.253
    Netmask: 255.255.255.0 (/24)
    Subnet: 192.168.1.0
    Gateway: 192.168.1.254

The subnet address is assumed to be the first address in the network range
and the gateway the last address. All other addresses in the range will be
allocated from the pool.

However, this is often not desirable: you may only want to allocate from
a subset of an existing network range. In such cases the SIZE, SUBNET,
and GATEWAY can be provided to configure the pool exactly as required:

    ip pool create MyPool 192.168.1.10/24 -s 50 -gw 192.168.1.1

    Created IP Pool 'MyPool':

    Address Range: 192.168.1.10 - 192.168.1.60
    Netmask: 255.255.255.0 (/24)
    Subnet: 192.168.1.0
    Gateway: 192.168.1.1
`,
		cmd.Summary(),
	)
}

func (cmd *create) setFlags(args []string) error {
	fs := flag.NewFlagSet("ip pool create flags", flag.ContinueOnError)
	cmd.FlagSet = fs

	cmd.Name = &args[0]
	cmd.Addr = &args[1]

	cmd.Size = fs.Int("size", 0, "size of the address pool")
	fs.IntVar(cmd.Size, "s", 0, "size of the address pool")

	cmd.Gateway = fs.String("gateway", "", "network gateway address")
	fs.StringVar(cmd.Gateway, "gw", "", "network gateway address")

	cmd.Subnet = fs.String("subnet", "", "network subnet address")
	fs.StringVar(cmd.Subnet, "sn", "", "network subnet address")

	if err := fs.Parse(args[2:]); err != nil {
		return err
	}

	return nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func decIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		if ip[j] > 0 {
			ip[j]--
			break
		} else {
			ip[j] = 255
		}
	}
}

type ipRange struct {
	FirstIP net.IP
	LastIP  net.IP
	Len     int
}

func ipRangeForNet(ip net.IP, ipnet net.IPNet) ipRange {
	r := ipRange{
		FirstIP: ip.Mask(ipnet.Mask),
		LastIP:  ip.Mask(ipnet.Mask),
		Len:     0,
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		r.Len++
		r.LastIP = ip
	}

	// move back one address, to account for the broadcast address
	r.Len--
	decIP(r.LastIP)

	// trim the range of unusable addresses
	if r.FirstIP[len(r.FirstIP)-1] == 0 {
		r.Len--
		incIP(r.FirstIP)
	}

	if r.LastIP[len(r.LastIP)-1] == 0 {
		r.Len--
		decIP(r.LastIP)
	}

	if r.LastIP[len(r.LastIP)-1] == 255 {
		r.Len--
		decIP(r.LastIP)
	}

	return r
}

func (cmd *create) ipPool() (*types.IpPool, error) {
	config := types.IpPoolIpPoolConfigInfo{}

	ip, network, err := net.ParseCIDR(*cmd.Addr)
	if err != nil {
		return nil, err
	}

	iprange := ipRangeForNet(ip, *network)

	config.Netmask = net.IP(network.Mask).String()

	if len(*cmd.Subnet) > 0 {
		config.SubnetAddress = *cmd.Subnet
	} else {
		config.SubnetAddress = network.IP.String()
	}

	if len(*cmd.Gateway) > 0 {
		config.Gateway = *cmd.Gateway
	} else {
		iprange.Len--
		decIP(iprange.LastIP)
		config.Gateway = iprange.LastIP.String()
	}

	if *cmd.Size == 0 {
		cmd.Size = &iprange.Len
	}

	config.Range = iprange.FirstIP.String() + "#" + strconv.Itoa(*cmd.Size)

	config.IpPoolEnabled = true

	return &types.IpPool{Name: *cmd.Name, Ipv4Config: &config}, nil
}

func (cmd *create) Run(cl *govmomi.Client, c *config.Config, args []string) error {
	if err := cmd.setFlags(args); err != nil {
		return err
	}

	ippool, err := cmd.ipPool()
	if err != nil {
		return err
	}

	p, err := pool.Create(cl, c.Datacenter, ippool)
	if err != nil {
		return err
	}

	fmt.Printf(
		`Created IP Pool '%s':

Address Range: %s
Netmask:       %s
Subnet:        %s
Gateway:       %s
`,
		p.Object.Name,
		p.Object.Ipv4Config.Range,
		p.Object.Ipv4Config.Netmask,
		p.Object.Ipv4Config.SubnetAddress,
		p.Object.Ipv4Config.Gateway,
	)

	return nil
}

func init() {
	cli.RegisterCommand([]string{"ip", "pool", "create"}, &create{})
}
