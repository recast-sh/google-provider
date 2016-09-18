package google

// TODO move out of here!
const coreosStable = "https://www.googleapis.com/compute/v1/projects/coreos-cloud/global/images/coreos-stable-1122-2-0-v20160906"
const coreosBeta = "https://www.googleapis.com/compute/v1/projects/coreos-cloud/global/images/coreos-beta-1153-4-0-v20160910"

func LoadImage(group, kind string) string {
	switch kind {
	case "stable":
		return coreosStable
	case "beta":
		return coreosBeta
	default:
		return coreosStable
	}
}
