/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBigipOnboardingUnitInvalid(t *testing.T) {
	resourceName := "regkeypool_name"
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigipOnboardingInvalid(resourceName),
				ExpectError: regexp.MustCompile(" Unsupported argument: An argument named \"invalidkey\" is not expected here"),
			},
		},
	})
}

func TestAccBigipOnboardingUnitCreate(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "fc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{
    "id": "fc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "ROLLING_BACK",
        "message": "invalid config - rolling back",
        "errors": [
            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
        ]
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux = http.NewServeMux()
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "bc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{
    "id": "bc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "ROLLING_BACK",
        "message": "invalid config - rolling back",
        "errors": [
            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
        ]
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipOnboardingCreate(resourceName, server.URL),
			},
			{
				Config: testAccBigipOnboardingModify(resourceName, server.URL),
				//ExpectError: regexp.MustCompile("Error while reading the response body :map\\[class:Result code:202 errors:\\[01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\"] message:invalid config - rolling back status:ROLLING_BACK]"),
			},
		},
	})
}

func TestAccBigipOnboardingUnitUpdateError(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	count := 0
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		count++
		if count == 1 {
			_, _ = fmt.Fprintf(w, `{
    "id": "fc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
		} else if count == 2 {
			_, _ = fmt.Fprintf(w, `{
    "id": "bc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }}`)
		}
	})

	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{
    "id": "fc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "ROLLING_BACK",
        "message": "invalid config - rolling back",
        "errors": [
            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
        ]
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{
	   "id": "bc1b334d-d593-4036-b3df-37de29ecc66b",
	   "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b",
	   "result": {
	       "class": "Result",
	       "code": 202,
	       "status": "ROLLING_BACK",
	       "message": "invalid config - rolling back",
	       "errors": [
	           "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
	       ]
	   },
	   "declaration": {
	       "schemaVersion": "1.20.0",
	       "class": "Device",
	       "async": true,
	       "label": "my BIG-IP declaration for declarative onboarding",
	       "Common": {
	           "class": "Tenant",
	           "hostname": "bigip1.example.com",
	           "ravinder": {
	               "class": "User",
	               "userType": "regular",
	               "partitionAccess": {
	                   "Common": {
	                       "role": "guest"
	                   }
	               },
	               "shell": "tmsh"
	           }
	       }
	   }
	}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipOnboardingCreate(resourceName, server.URL),
				//ExpectError: regexp.MustCompile("Error while reading the response body :map\\[class:Result code:202 errors:\\[01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\"] message:invalid config - rolling back status:ROLLING_BACK]"),
			},
			{
				Config:      testAccBigipOnboardingModify(resourceName, server.URL),
				ExpectError: regexp.MustCompile("Error while reading the response body :map\\[class:Result code:202 errors:\\[01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\"] message:invalid config - rolling back status:ROLLING_BACK]"),
			},
		},
	})
}

func TestAccBigipOnboardingUnitReadError(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "ROLLING_BACK",
        "message": "invalid config - rolling back",
        "errors": [
            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
        ]
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigipOnboardingCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("Error while reading the response body :map\\[class:Result code:202 errors:\\[01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\"] message:invalid config - rolling back status:ROLLING_BACK]"),
			},
		},
	})
}

func TestAccBigipOnboardingUnitReadErrorNew(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{
    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "ROLLING_BACK",
        "message": "invalid config - rolling back",
        "errors": [
            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
        ]
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipOnboardingCreate(resourceName, server.URL),
			},
		},
	})
}
func TestAccBigipOnboardingUnitCreateError(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		_, _ = fmt.Fprintf(w, `{
    "id": "50ce5959-a256-463d-92e9-eee11b20d229",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBigipOnboardingCreate(resourceName, server.URL),
				ExpectError: regexp.MustCompile("Timeout while polling the DO task id with result:map\\[class:Result code:202 message:processing status:RUNNING]"),
			},
		},
	})
}

func TestAccBigipOnboardingUnitLocalUpdateError(t *testing.T) {
	resourceName := "regkeypool_name"
	setup()
	log.Printf("URL:%+v", server.URL)
	count := 0
	mux.HandleFunc("mgmt/shared/authn/login", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	})
	mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		_, _ = fmt.Fprintf(w, `{}`)
	})
	//mux.HandleFunc("/mgmt/tm/net/self", func(w http.ResponseWriter, r *http.Request) {
	//	assert.Equal(t, "GET", r.Method, "Expected method 'GET', got %s", r.Method)
	//	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	//	_, _ = fmt.Fprintf(w, `{}`)
	//})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		count++
		if count == 1 {
			_, _ = fmt.Fprintf(w, `{
    "id": "fc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
		} else if count == 2 {
			_, _ = fmt.Fprintf(w, `{
    "id": "bc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "RUNNING",
        "message": "processing"
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }}`)
		}
	})

	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/fc1b334d-d593-4036-b3df-37de29ecc66b", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{
    "id": "fc1b334d-d593-4036-b3df-37de29ecc66b",
    "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/50ce5959-a256-463d-92e9-eee11b20d229",
    "result": {
        "class": "Result",
        "code": 202,
        "status": "ROLLING_BACK",
        "message": "invalid config - rolling back",
        "errors": [
            "01070734:3: Configuration error: A monitor may not default from itself \"/Common/http\""
        ]
    },
    "declaration": {
        "schemaVersion": "1.20.0",
        "class": "Device",
        "async": true,
        "label": "my BIG-IP declaration for declarative onboarding",
        "Common": {
            "class": "Tenant",
            "hostname": "bigip1.example.com",
            "ravinder": {
                "class": "User",
                "userType": "regular",
                "partitionAccess": {
                    "Common": {
                        "role": "guest"
                    }
                },
                "shell": "tmsh"
            }
        }
    }
}`)
	})
	mux.HandleFunc("/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintf(w, `{
	   "id": "bc1b334d-d593-4036-b3df-37de29ecc66b",
	   "selfLink": "https://localhost/mgmt/shared/declarative-onboarding/task/bc1b334d-d593-4036-b3df-37de29ecc66b",
       "result": {
			"class": "Result",
			"code": 202,
			"status": "RUNNING",
			"message": "processing"
    	},
	   "declaration": {
	       "schemaVersion": "1.20.0",
	       "class": "Device",
	       "async": true,
	       "label": "my BIG-IP declaration for declarative onboarding",
	       "Common": {
	           "class": "Tenant",
	           "hostname": "bigip1.example.com",
	           "ravinder": {
	               "class": "User",
	               "userType": "regular",
	               "partitionAccess": {
	                   "Common": {
	                       "role": "guest"
	                   }
	               },
	               "shell": "tmsh"
	           }
	       }
	   }
	}`)
	})
	defer teardown()
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBigipOnboardinglocalCreate(resourceName, server.URL),
			},
			{
				Config:      testAccBigipOnboardinglocalModify(resourceName, server.URL),
				ExpectError: regexp.MustCompile("Timeout while polling the DO task id with result:map\\[class:Result code:202 message:processing status:RUNNING]"),
			},
		},
	})
}

func testAccBigipOnboardingInvalid(resourceName string) string {
	return fmt.Sprintf(`
resource "bigip_do" "do-example" {
  do_json = "${file("` + folder + `/../examples/do/basic_do.json")}"
  timeout = 1
  invalidkey = "foo"
}
provider "bigip" {
  address  = "xxx.xxx.xxx.xxx"
  username = "xxx"
  password = "xxx"
}`)
}

func testAccBigipOnboardingCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_do" "do-example" {
  do_json = "${file("`+folder+`/../examples/do/basic_do.json")}"
  timeout = 1
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, url)
}

func testAccBigipOnboardingModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_do" "do-example" {
  do_json = "${file("`+folder+`/../examples/do/basic_do2.json")}"
  timeout = 1
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, url)
}

func testAccBigipOnboardinglocalCreate(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_do" "do-example" {
  do_json = "${file("`+folder+`/../examples/do/basic_do.json")}"
  bigip_address = "%s"
  bigip_user = "xxx"
  bigip_password = "xxx"
  bigip_port = "%s"
  timeout = 1
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, strings.Join(strings.Split(url, ":")[0:2], ":"), strings.Split(url, ":")[2], url)
}

func testAccBigipOnboardinglocalModify(resourceName, url string) string {
	return fmt.Sprintf(`
resource "bigip_do" "do-example" {
  do_json = "${file("`+folder+`/../examples/do/basic_do2.json")}"
  bigip_address = "%s"
  bigip_user = "xxx"
  bigip_password = "xxx"
  bigip_port = "%s"
  timeout = 1
}
provider "bigip" {
  address  = "%s"
  username = ""
  password = ""
  login_ref = ""
}`, strings.Join(strings.Split(url, ":")[0:2], ":"), strings.Split(url, ":")[2], url)
}
