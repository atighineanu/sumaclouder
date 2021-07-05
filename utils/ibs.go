package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

const (
	downloadIbsLink = "http://download.suse.de/ibs/SUSE:/Maintenance:/"
)

func SyncMUChannel(mu string) error {
	mu = "SUSE:Maintenance:20223:244004"
	incidentAndReleaseAndPrefixes := strings.Split(mu, ":")
	if len(incidentAndReleaseAndPrefixes) < 4 {
		return fmt.Errorf("The MU is formatted wrong... check the MU SUSE:Maintenance:<incident_number>:<rr_number>")
	}
	compoundIbsLink := fmt.Sprintf("%s%s/", downloadIbsLink, incidentAndReleaseAndPrefixes[2])
	fmt.Println(compoundIbsLink)

	resp, err := http.Get(compoundIbsLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	regexpServer := regexp.MustCompile(`SUSE-Manager-Server_\d{1}.\d{1}_x86_64`)
	regexpProxy := regexp.MustCompile(`SUSE-Manager-Proxy_\d{1}.\d{1}_x86_64`)
	serverSuffix, _ := ProcessWebpage(*regexpServer, fmt.Sprintf("%s", string(body)))
	proxySuffix, _ := ProcessWebpage(*regexpProxy, fmt.Sprintf("%s", string(body)))

	cmdDownloadProxyStuff := []string{"wget", "--no-parent", "-r", fmt.Sprintf("%sSUSE_Updates_SLE-Module-%s/", compoundIbsLink, proxySuffix)}
	cmdDownloadServerStuff := []string{"wget", "--no-parent", "-r", fmt.Sprintf("%sSUSE_Updates_SLE-Module-%s/", compoundIbsLink, serverSuffix)}

	_, err = exec.Command(cmdDownloadProxyStuff[0], cmdDownloadProxyStuff[1:]...).CombinedOutput()
	if err != nil {
		return err
	}
	_, err = exec.Command(cmdDownloadServerStuff[0], cmdDownloadServerStuff[1:]...).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func ProcessWebpage(reg regexp.Regexp, rawOutput string) (string, error) {
	if reg.FindString(rawOutput) != "" {
		//fmt.Println(reg.FindString(rawOutput))
		return reg.FindString(rawOutput), nil
	}
	return "", nil
}
