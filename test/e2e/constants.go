/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package e2e

const (
	proxy                 = "localhost:65001"
	hostnameFilePath      = "/etc/hostname"
	dragonflyNamespace    = "dragonfly-system"
	dragonflyE2ENamespace = "dragonfly-e2e"
)

const (
	dfdaemonCompatibilityTestMode = "dfdaemon"
)

const (
	managerServerName   = "manager"
	schedulerServerName = "scheduler"
	cdnServerName       = "cdn"
	dfdaemonServerName  = "dfdaemon"
	proxyServerName     = "proxy"
)

type server struct {
	component  string
	namespace  string
	container  string
	logDirName string
	logPrefix  string
}

var servers = map[string]server{
	managerServerName: {
		component:  managerServerName,
		namespace:  dragonflyNamespace,
		container:  managerServerName,
		logDirName: managerServerName,
		logPrefix:  managerServerName,
	},
	schedulerServerName: {
		component:  schedulerServerName,
		namespace:  dragonflyNamespace,
		container:  schedulerServerName,
		logDirName: schedulerServerName,
		logPrefix:  schedulerServerName,
	},
	cdnServerName: {
		component:  cdnServerName,
		namespace:  dragonflyNamespace,
		container:  cdnServerName,
		logDirName: cdnServerName,
		logPrefix:  cdnServerName,
	},
	dfdaemonServerName: {
		component:  dfdaemonServerName,
		namespace:  dragonflyNamespace,
		container:  dfdaemonServerName,
		logDirName: "daemon",
		logPrefix:  "daemon",
	},
	proxyServerName: {
		component:  proxyServerName,
		namespace:  dragonflyE2ENamespace,
		container:  proxyServerName,
		logDirName: "daemon",
		logPrefix:  proxyServerName,
	},
}
