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

//func GetPageBody(url string)

func SyncMUChannel(mu string) error {
	mu = "SUSE:Maintenance:20223:244004"

	maintUpd, err := ReturnMU(mu)
	if err != nil {
		return err
	}

	compoundIbsLink := fmt.Sprintf("%s%s/", downloadIbsLink, maintUpd.Incident)

	//fmt.Println(compoundIbsLink)
	resp, err := http.Get(compoundIbsLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// finding Server and Proxy x86_64 folders in the ibs:/SUSE:/Maintenance:/Incident folder
	regexpServer := regexp.MustCompile(`"SUSE-Manager-Server_\d{1}.\d{1}_x86_64"`)
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

func ReturnMU(mu string) (MU, error) {
	var muStruct MU
	sliceMU := strings.Split(mu, ":")
	if len(sliceMU) < 4 {
		return muStruct, fmt.Errorf("The MU is formatted wrong... check the MU SUSE:Maintenance:<incident_number>:<rr_number>")
	}
	muStruct.Prefix1 = sliceMU[0]
	muStruct.Prefix2 = sliceMU[1]
	muStruct.Incident = sliceMU[2]
	muStruct.ReleaseRequest = sliceMU[3]
	return muStruct, nil
}

func ProcessWebpage(reg regexp.Regexp, rawOutput string) (string, error) {
	if reg.FindString(rawOutput) != "" {
		return reg.FindString(rawOutput), nil
	}
	return "", nil
}
