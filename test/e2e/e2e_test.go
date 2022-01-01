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

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo" //nolint
	. "github.com/onsi/gomega" //nolint

	"d7y.io/dragonfly/v2/test/e2e/e2eutil"
	_ "d7y.io/dragonfly/v2/test/e2e/manager"
)

var _ = AfterSuite(func() {
	for _, server := range servers {
		for i := 0; i < 3; i++ {
			out, err := e2eutil.KubeCtlCommand("-n", server.namespace, "get", "pod", "-l", fmt.Sprintf("component=%s", server.component),
				"-o", fmt.Sprintf("jsonpath='{.items[%d].metadata.name}'", i)).CombinedOutput()
			if err != nil {
				fmt.Printf("get pod error: %s\n", err)
				continue
			}
			podName := strings.Trim(string(out), "'")
			pod := e2eutil.NewPodExec(server.namespace, podName, server.container)

			countOut, err := e2eutil.KubeCtlCommand("-n", server.namespace, "get", "pod", "-l", fmt.Sprintf("component=%s", server.component),
				"-o", fmt.Sprintf("jsonpath='{.items[%d].status.containerStatuses[0].restartCount}'", i)).CombinedOutput()
			if err != nil {
				fmt.Printf("get pod %s restart count error: %s\n", podName, err)
				continue
			}
			rawCount := strings.Trim(string(countOut), "'")
			count, err := strconv.Atoi(rawCount)
			if err != nil {
				fmt.Printf("atoi error: %s\n", err)
				continue
			}
			fmt.Printf("pod %s restart count: %d\n", podName, count)

			if count > 0 {
				if err := e2eutil.UploadArtifactStdout(server.namespace, podName, server.logDirName, fmt.Sprintf("%s-%d-prev", server.logPrefix, i)); err != nil {
					fmt.Printf("upload pod %s artifact stdout file error: %v\n", podName, err)
				}
			}

			if err := e2eutil.UploadArtifactStdout(server.namespace, podName, server.logDirName, fmt.Sprintf("%s-%d", server.logPrefix, i)); err != nil {
				fmt.Printf("upload pod %s artifact prev stdout file error: %v\n", podName, err)
			}

			out, err = pod.Command("sh", "-c", fmt.Sprintf(`
              set -x
              cp /var/log/dragonfly/%s/core.log /tmp/artifact/%s/%s-%d-core.log
              cp /var/log/dragonfly/%s/grpc.log /tmp/artifact/%s/%s-%d-grpc.log
              cp /var/log/dragonfly/%s/gin.log /tmp/artifact/%s/%s-%d-gin.log
              `, server.logDirName, server.logDirName, server.logPrefix, i, server.logDirName, server.logDirName, server.logPrefix, i, server.logDirName, server.logDirName, server.logPrefix, i)).CombinedOutput()
			if err != nil {
				fmt.Printf("copy log output: %s, error: %s\n", string(out), err)
			}
		}
	}
})

var _ = BeforeSuite(func() {
	mode := os.Getenv("DRAGONFLY_COMPATIBILITY_E2E_TEST_MODE")
	if mode != "" {
		rawImages, err := e2eutil.KubeCtlCommand("-n", dragonflyNamespace, "get", "pod", "-l", fmt.Sprintf("component=%s", mode),
			"-o", "jsonpath='{range .items[0]}{.spec.containers[0].image}{end}'").CombinedOutput()
		image := strings.Trim(string(rawImages), "'")
		Expect(err).NotTo(HaveOccurred())
		fmt.Printf("special image name: %s\n", image)

		stableImageTag := os.Getenv("DRAGONFLY_STABLE_IMAGE_TAG")
		Expect(fmt.Sprintf("dragonflyoss/%s:%s", mode, stableImageTag)).To(Equal(image))
	}

	rawGitCommit, err := e2eutil.GitCommand("rev-parse", "--short", "HEAD").CombinedOutput()
	Expect(err).NotTo(HaveOccurred())
	gitCommit := strings.Fields(string(rawGitCommit))[0]
	fmt.Printf("git merge commit: %s\n", gitCommit)

	rawPodName, err := e2eutil.KubeCtlCommand("-n", dragonflyNamespace, "get", "pod", "-l", "component=dfdaemon",
		"-o", "jsonpath='{range .items[*]}{.metadata.name}{end}'").CombinedOutput()
	podName := strings.Trim(string(rawPodName), "'")
	Expect(err).NotTo(HaveOccurred())
	Expect(strings.HasPrefix(podName, "dragonfly-dfdaemon-")).Should(BeTrue())

	pod := e2eutil.NewPodExec(dragonflyNamespace, podName, "dfdaemon")
	rawDfgetVersion, err := pod.Command("dfget", "version").CombinedOutput()
	Expect(err).NotTo(HaveOccurred())
	dfgetGitCommit := strings.Fields(string(rawDfgetVersion))[7]
	fmt.Printf("raw dfget version: %s\n", rawDfgetVersion)
	fmt.Printf("dfget merge commit: %s\n", dfgetGitCommit)

	if mode == dfdaemonCompatibilityTestMode {
		Expect(gitCommit).NotTo(Equal(dfgetGitCommit))
		return
	}

	Expect(gitCommit).To(Equal(dfgetGitCommit))
})

// TestE2E is the root of e2e test function
func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "dragonfly e2e test suite")
}
