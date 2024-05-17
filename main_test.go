package main

import (
	"doumor/metr1c/rac"
	"testing"
)

func TestCountSessionTypesSingleType(t *testing.T) {
	sessions := rac.RACQuery{}
	sessions.Output = `session                          : 8de713f9-0f21-408f-98ec-0edb40a6acbc
session-id                       : 1
infobase                         : f2bad9ba-3461-4d7a-96a2-0c05bce92369
user-name                        : Администратор
app-id                           : WebClient
last-active-at                   : 2024-05-17T11:34:53
hibernate                        : no`

	err := sessions.Parse()
	if err != nil {
		t.Error(err)
	}
	expectedActive := 1.0
	expectedHibernating := 0.0

	actualActive, actualHibernating := countSessionTypes(sessions)
	if actualActive != expectedActive || actualHibernating != expectedHibernating {
		t.Fatalf("(actual) %#v, %#v != %#v, %#v (expected)\n", actualActive, actualHibernating, expectedActive, expectedHibernating)
	}
}

func TestCountSessionTypesMultipleTypes(t *testing.T) {
	sessions := rac.RACQuery{}
	sessions.Output = `session                          : 8de713f9-0f21-408f-98ec-0edb40a6acbc
session-id                       : 1
infobase                         : f2bad9ba-3461-4d7a-96a2-0c05bce92369
user-name                        : Администратор
app-id                           : WebClient
last-active-at                   : 2024-05-17T11:34:53
hibernate                        : no

session            				 : 2ab1c7d5-97f8-4032-855a-cdac78989cc0
session-id                       : 1
infobase                         : f2bad9ba-3461-4d7a-96a2-0c05bce92369
user-name                        : Zeleboba
app-id                           : WebClient
last-active-at                   : 2024-05-17T11:34:53
hibernate                        : yes`

	err := sessions.Parse()
	if err != nil {
		t.Error(err)
	}

	expectedActive := 1.0
	expectedHibernating := 1.0

	actualActive, actualHibernating := countSessionTypes(sessions)
	if actualActive != expectedActive || actualHibernating != expectedHibernating {
		t.Fatalf("(actual) %#v, %#v != %#v, %#v (expected)\n", actualActive, actualHibernating, expectedActive, expectedHibernating)
	}
}

func TestCountLicenseTypesSingleLicense(t *testing.T) {
	licenses := rac.RACQuery{}
	licenses.Output = `session            : 8de713f9-0f21-408f-98ec-0edb40a6acbc
user-name          : Администратор
app-id             : WebClient
license-type       : soft`

	err := licenses.Parse()
	if err != nil {
		t.Error(err)
	}

	expectedSoft := 1.0
	expectedHASP := 0.0

	actualSoft, actualHASP := countLicenseTypes(licenses)
	if actualSoft != expectedSoft || actualHASP != expectedHASP {
		t.Fatalf("(actual) %#v, %#v != %#v, %#v (expected)\n", actualSoft, actualHASP, expectedSoft, expectedHASP)
	}
}

func TestCountLicenseTypesMultipleLicenseTypes(t *testing.T) {
	licenses := rac.RACQuery{}
	licenses.Output = `session            : 8de713f9-0f21-408f-98ec-0edb40a6acbc
user-name          : Администратор
app-id             : WebClient
license-type       : soft

session            : 7f21962e-3269-4c88-9c72-a6aa9d15edda
user-name          : Zeleboba
app-id             : WebClient
license-type       : HASP`

	err := licenses.Parse()
	if err != nil {
		t.Error(err)
	}

	expectedSoft := 1.0
	expectedHASP := 1.0

	actualSoft, actualHASP := countLicenseTypes(licenses)
	if actualSoft != expectedSoft || actualHASP != expectedHASP {
		t.Fatalf("(actual) %#v, %#v != %#v, %#v (expected)\n", actualSoft, actualHASP, expectedSoft, expectedHASP)
	}
}

func TestCountTotalProcMemSingleProc(t *testing.T) {
	processes := rac.RACQuery{}
	processes.Output = `process              : 7f21962e-3269-4c88-9c72-a6aa9d15edda
host                 : genlab-1c-test
port                 : 1560
pid                  : 1324
memory-size          : 1024`

	err := processes.Parse()
	if err != nil {
		t.Error(err)
	}

	expected := 1024.0
	actual, err := countTotalProcMem(processes)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Fatalf("(actual) %#v != %#v (expected)\n", actual, expected)
	}
}

func TestCountTotalProcMemMultiplePocs(t *testing.T) {
	processes := rac.RACQuery{}
	processes.Output = `process              : 7f21962e-3269-4c88-9c72-a6aa9d15edda
host                 : genlab-1c-test
port                 : 1560
pid                  : 1324
memory-size          : 1024

process              : 70f69831-810d-4d56-8e2b-de7da25ce565
host                 : genlab-1c-test
port                 : 1560
pid                  : 1324
memory-size          : 2048`

	err := processes.Parse()
	if err != nil {
		t.Error(err)
	}

	expected := 3072.0
	actual, err := countTotalProcMem(processes)
	if err != nil {
		t.Error(err)
	}

	if actual != expected {
		t.Fatalf("(actual) %#v != %#v (expected)\n", actual, expected)
	}
}
