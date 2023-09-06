// Copyright 2018 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fatedier/frp/pkg/config/types"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/consts"
)

func init() {
	RegisterCommonFlags(stcpCmd)

	stcpCmd.PersistentFlags().StringVarP(&proxyName, "proxy_name", "n", "", "proxy name")
	stcpCmd.PersistentFlags().StringVarP(&role, "role", "", "server", "role")
	stcpCmd.PersistentFlags().StringVarP(&sk, "sk", "", "", "secret key")
	stcpCmd.PersistentFlags().StringVarP(&serverName, "server_name", "", "", "server name")
	stcpCmd.PersistentFlags().StringVarP(&localIP, "local_ip", "i", "127.0.0.1", "local ip")
	stcpCmd.PersistentFlags().IntVarP(&localPort, "local_port", "l", 0, "local port")
	stcpCmd.PersistentFlags().StringVarP(&bindAddr, "bind_addr", "", "", "bind addr")
	stcpCmd.PersistentFlags().IntVarP(&bindPort, "bind_port", "", 0, "bind port")
	stcpCmd.PersistentFlags().BoolVarP(&useEncryption, "ue", "", false, "use encryption")
	stcpCmd.PersistentFlags().BoolVarP(&useCompression, "uc", "", false, "use compression")
	stcpCmd.PersistentFlags().StringVarP(&bandwidthLimit, "bandwidth_limit", "", "", "bandwidth limit")
	stcpCmd.PersistentFlags().StringVarP(&bandwidthLimitMode, "bandwidth_limit_mode", "", types.BandwidthLimitModeClient, "bandwidth limit mode")

	rootCmd.AddCommand(stcpCmd)
}

var stcpCmd = &cobra.Command{
	Use:   "stcp",
	Short: "Run frpc with a single stcp proxy",
	RunE: func(cmd *cobra.Command, args []string) error {
		clientCfg, err := parseClientCommonCfgFromCmd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pxyCfgs := make([]v1.ProxyConfigurer, 0)
		visitorCfgs := make([]v1.VisitorConfigurer, 0)

		var prefix string
		if user != "" {
			prefix = user + "."
		}

		switch role {
		case "server":
			cfg := &v1.STCPProxyConfig{}
			cfg.Name = prefix + proxyName
			cfg.Type = consts.STCPProxy
			cfg.Transport.UseEncryption = useEncryption
			cfg.Transport.UseCompression = useCompression
			cfg.Secretkey = sk
			cfg.LocalIP = localIP
			cfg.LocalPort = localPort
			cfg.Transport.BandwidthLimit, err = types.NewBandwidthQuantity(bandwidthLimit)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			cfg.Transport.BandwidthLimitMode = bandwidthLimitMode
			if err := validation.ValidateProxyConfigurerForClient(cfg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pxyCfgs = append(pxyCfgs, cfg)
		case "visitor":
			cfg := &v1.STCPVisitorConfig{}
			cfg.Name = prefix + proxyName
			cfg.Type = consts.STCPProxy
			cfg.Transport.UseEncryption = useEncryption
			cfg.Transport.UseCompression = useCompression
			cfg.SecretKey = sk
			cfg.ServerName = serverName
			cfg.BindAddr = bindAddr
			cfg.BindPort = bindPort
			if err := validation.ValidateVisitorConfigurer(cfg); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			visitorCfgs = append(visitorCfgs, cfg)
		default:
			fmt.Println("invalid role")
			os.Exit(1)
		}

		err = startService(clientCfg, pxyCfgs, visitorCfgs, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
