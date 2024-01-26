package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	HealthyInstances []*NATInstance
	Subnets          []*Subnet
}

type Result struct {
	NatInstances     []string
	AssociateSubnets map[string][]string
}

// Function to check if an element exists in a list
func elementExists(element string, list []string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func Test_allocate(t *testing.T) {
	testcases := []TestCase{
		TestCase{
			HealthyInstances: []*NATInstance{
				&NATInstance{
					Id:   "1",
					Zone: "us-west1-a",
				},
				&NATInstance{
					Id:   "2",
					Zone: "us-west1-b",
				},
				&NATInstance{
					Id:   "3",
					Zone: "us-west1-b",
				},
			},

			Subnets: []*Subnet{
				&Subnet{
					Id:   "1",
					Zone: "us-west1-a",
				},
				&Subnet{
					Id:   "2",
					Zone: "us-west1-b",
				},
				&Subnet{
					Id:   "3",
					Zone: "us-west1-b",
				},
				&Subnet{
					Id:   "4",
					Zone: "us-west1-c",
				},
			},
		},

		TestCase{
			HealthyInstances: []*NATInstance{
				&NATInstance{
					Id:   "1",
					Zone: "us-west1-a",
				},
				&NATInstance{
					Id:   "2",
					Zone: "us-west1-b",
				},
				&NATInstance{
					Id:   "3",
					Zone: "us-west1-b",
				},
			},

			Subnets: []*Subnet{
				&Subnet{
					Id:   "1",
					Zone: "us-west1-a",
				},
				&Subnet{
					Id:   "2",
					Zone: "us-west1-b",
				},
				&Subnet{
					Id:   "3",
					Zone: "us-west1-b",
				},
			},
		},
	}

	expectations := []Result{
		Result{
			NatInstances: []string{
				"1-us-west1-a", "2-us-west1-b", "3-us-west1-b",
			},

			AssociateSubnets: map[string][]string{
				"1-us-west1-a": []string{"1-us-west1-a", "4-us-west1-c"},
				"2-us-west1-b": []string{"2-us-west1-b"},
				"3-us-west1-b": []string{"3-us-west1-b"},
			},
		},

		Result{
			NatInstances: []string{
				"1-us-west1-a", "2-us-west1-b", "3-us-west1-b",
			},

			AssociateSubnets: map[string][]string{
				"1-us-west1-a": []string{"1-us-west1-a", "4-us-west1-c"},
				"2-us-west1-b": []string{"2-us-west1-b"},
				"3-us-west1-b": []string{"3-us-west1-b"},
			},
		},
	}

	for index, test := range testcases {
		instancesInZone := make(map[string][]*NATInstance)
		instancesInZone = mapHealthyInstancestoZone(test.HealthyInstances)
		allocate(instancesInZone, test.Subnets)

		result := Result{}
		result.AssociateSubnets = map[string][]string{}
		for _, instances := range instancesInZone {
			for _, i := range instances {
				instanceName := fmt.Sprintf("%v-%v", i.Id, i.Zone)
				result.NatInstances = append(result.NatInstances, instanceName)
				for _, s := range i.Subnets {
					subnet := fmt.Sprintf("%v-%v", s.Id, s.Zone)
					result.AssociateSubnets[instanceName] = append(result.AssociateSubnets[instanceName], subnet)
				}
			}
		}

		// Check results

		require.Equal(t, len(expectations[index].NatInstances), len(result.NatInstances))

		// Check associated subnets to a specific NAT instance
		// for _, instance := range expectations[index].NatInstances {
		// 	expectList := expectations[index].AssociateSubnets[instance]
		// 	require.Equal(t, elementExists(instance, expectList), true)
		// }

		for _, instance := range expectations[index].NatInstances {
			require.Equal(t, len(expectations[index].AssociateSubnets[instance]), len(result.AssociateSubnets[instance]))
		}
	}
}
