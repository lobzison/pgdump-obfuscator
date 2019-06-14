package main

type Configuration struct {
	Obfuscations []TargetedObfuscation
}

// TODO: read from file?
var Config *Configuration = &Configuration{
	Obfuscations: []TargetedObfuscation{
		TargetedObfuscation{
			Target{Table: "leads", Column: "tsv"},
			ScrambleBytes,
		},
		
		
		TargetedObfuscation{
			Target{Table: "customers", Column: "first_name"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "customers", Column: "last_name"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "customers", Column: "middle_name"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "customers", Column: "email"},
			ScrambleUniqueEmail,
		},
		TargetedObfuscation{
			Target{Table: "customers", Column: "hashed_email"},
			ScrambleBytes,
		},

		
		TargetedObfuscation{
			Target{Table: "customer_addresses", Column: "address_line1"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "customer_addresses", Column: "address_line2"},
			ScrambleBytes,
		},

		
		TargetedObfuscation{
			Target{Table: "customer_details", Column: "fico_score"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "customer_details", Column: "phone"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "customer_details", Column: "business_phone"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "customer_details", Column: "ssn"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "customer_details", Column: "short_ssn"},
			ScrambleDigits,
		},

		
		TargetedObfuscation{
			Target{Table: "loans", Column: "lender_name"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "loans", Column: "amount"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "loans", Column: "number"},
			ScrambleBytes,
		},
		
		
		TargetedObfuscation{
			Target{Table: "coborrowers", Column: "name"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "coborrowers", Column: "email"},
			ScrambleEmail,
		},
		TargetedObfuscation{
			Target{Table: "coborrowers", Column: "phone"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "coborrowers", Column: "business_phone"},
			ScrambleDigits,
		},
		TargetedObfuscation{
			Target{Table: "coborrowers", Column: "ssn"},
			ScrambleDigits,
		},
		
		
		TargetedObfuscation{
			Target{Table: "policies", Column: "company"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "policies", Column: "number"},
			ScrambleBytes,
		},
		TargetedObfuscation{
			Target{Table: "policies", Column: "normalized_company"},
			ScrambleBytes,
		},
		
		TargetedObfuscation{
			Target{Table: "quotes", Column: "online_bind_url"},
			ScrambleBindUrls,
		},
		
		TargetedObfuscation{
			Target{Table: "api_clients", Column: "secret"},
			ScrambleBytes,
		},
		
		TargetedObfuscation{
			Target{Table: "blacklist_contacts", Column: "email"},
			ScrambleEmail,
		},
		TargetedObfuscation{
			Target{Table: "blacklist_contacts", Column: "phone"},
			ScrambleDigits,
		},
	},
}
