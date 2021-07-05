package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	downloadIbsLink = "http://download.suse.de/ibs/"
)

func SyncMUChannel(mu string) error {
	var maintUpdate MU
	mu = "SUSE:Maintenance:20223:244004"
	incidentAndReleaseAndPrefixes := strings.Split(mu, ":")
	if len(incidentAndReleaseAndPrefixes) < 4 {
		return fmt.Errorf("The MU is formatted wrong... check the MU SUSE:Maintenance:<incident_number>:<rr_number>")
	}
	maintUpdate.Prefix1, maintUpdate.Prefix2 = incidentAndReleaseAndPrefixes[0], incidentAndReleaseAndPrefixes[1]
	maintUpdate.Incident, maintUpdate.ReleaseRequest = incidentAndReleaseAndPrefixes[2], incidentAndReleaseAndPrefixes[3]
	compoundIbsLink := fmt.Sprintf("%s%s:/%s:/%s/", downloadIbsLink, maintUpdate.Prefix1, maintUpdate.Prefix2, maintUpdate.Incident)
	cmd := []string{"curl", compoundIbsLink}
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return err
	}
	regexpServer := regexp.MustCompile(`SUSE-Manager-Server_\d{1}.\d{1}_x86_64`)
	regexpProxy := regexp.MustCompile(`SUSE-Manager-Proxy_\d{1}.\d{1}_x86_64`)
	serverSuffix, _ := ProcessWebpage(*regexpServer, fmt.Sprintf("%s", string(out)))
	proxySuffix, _ := ProcessWebpage(*regexpProxy, fmt.Sprintf("%s", string(out)))
	//fmt.Println(servSuffix)
	//fmt.Println(fmt.Sprintf("%sSUSE_Updates_SLE-Module-%s/", compoundIbsLink, proxySuffix))
	cmdDownloadProxyStuff := []string{"wget", "--no-parent", "-r", fmt.Sprintf("%sSUSE_Updates_SLE-Module-%s/", compoundIbsLink, proxySuffix)}
	cmdDownloadServerStuff := []string{"wget", "--no-parent", "-r", fmt.Sprintf("%sSUSE_Updates_SLE-Module-%s/", compoundIbsLink, serverSuffix)}
	out, err = exec.Command(cmdDownloadProxyStuff[0], cmdDownloadProxyStuff[1:]...).CombinedOutput()
	if err != nil {
		return err
	}
	out, err = exec.Command(cmdDownloadServerStuff[0], cmdDownloadServerStuff[1:]...).CombinedOutput()
	if err != nil {
		return err
	}
	//fmt.Println(fmt.Sprintf("%s", string(out)))
	//fmt.Println(compoundIbsLink)
	//curl -LO
	//fmt.Println(maintUpdate.Prefix)
	//fmt.Println(maintUpdate.ReleaseRequest)
	//fmt.Println(maintUpdate.Incident)
	return nil
}

func ProcessWebpage(reg regexp.Regexp, rawOutput string) (string, error) {
	if reg.FindString(rawOutput) != "" {
		//fmt.Println(reg.FindString(rawOutput))
		return reg.FindString(rawOutput), nil
	}
	return "", nil
}
