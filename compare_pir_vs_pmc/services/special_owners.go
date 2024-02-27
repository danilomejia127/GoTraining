package services

import "github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/apicalls"

func UpdateSpecialOwnersKVS(keys []string) {
	for _, key := range keys {
		specialOwnersProd := apicalls.GetOriginalDataSpecialOwnersBySite(key)
		if len(specialOwnersProd) > 0 {
			specialOwnersMsg := getMsgWithSpecialOwners(key, specialOwnersProd)
			apicalls.UpdateSpecialOwnersIntoKVS(apicalls.ProductionSynchronizerStgURL, specialOwnersMsg)
		}
	}

}

func getMsgWithSpecialOwners(key string, specialOwners []string) apicalls.SpecialOwnersMsg {
	msg := apicalls.SpecialOwnersMsg{}
	msg.Msg.Key = key
	msg.Msg.Value = []struct {
		ID     int      `json:"id"`
		Values []string `json:"values"`
	}{
		{
			ID:     1,
			Values: specialOwners,
		},
	}

	return msg
}
