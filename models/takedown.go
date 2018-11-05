package models

var (
	Content struct {
		Response map[string]int `json:"response"`
	}
	TakedownCategory = map[int]string{
		0:  "none",
		1:  "other",
		2:  "bots",
		3:  "brute_force_attacks",
		4:  "copyright_infringement_fraud",
		5:  "denial_of_service",
		6:  "disposable_email_address",
		7:  "download_leech_software",
		8:  "duplicate_accounts",
		9:  "excessive_resource_utilization",
		10: "illegal_filesharing",
		11: "keylogging",
		12: "malware_trojan_malvertising",
		13: "network_proxy",
		14: "objectionable_or_obscene_content",
		15: "phishing",
		16: "spamming",
		17: "virtual_currency_mining",
		18: "vulnerability_scanning"}
)
