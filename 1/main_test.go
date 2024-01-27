package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	HealthyInstances []*NATInstance
	Subnets          []*Subnet
	Expected         []Expectation
}

type Expectation struct {
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

			Expected: []Expectation{
				Expectation{
					NatInstances: []string{
						"1-us-west1-a", "2-us-west1-b", "3-us-west1-b",
					},

					AssociateSubnets: map[string][]string{
						"1-us-west1-a": []string{"1-us-west1-a", "4-us-west1-c"},
						"2-us-west1-b": []string{"2-us-west1-b"},
						"3-us-west1-b": []string{"3-us-west1-b"},
					},
				},

				Expectation{
					NatInstances: []string{
						"1-us-west1-a", "2-us-west1-b", "3-us-west1-b",
					},

					AssociateSubnets: map[string][]string{
						"1-us-west1-a": []string{"1-us-west1-a"},
						"2-us-west1-b": []string{"2-us-west1-b", "4-us-west1-c"},
						"3-us-west1-b": []string{"3-us-west1-b"},
					},
				},

				Expectation{
					NatInstances: []string{
						"1-us-west1-a", "2-us-west1-b", "3-us-west1-b",
					},

					AssociateSubnets: map[string][]string{
						"1-us-west1-a": []string{"1-us-west1-a"},
						"2-us-west1-b": []string{"2-us-west1-b"},
						"3-us-west1-b": []string{"3-us-west1-b", "4-us-west1-c"},
					},
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
					Zone: "us-west1-c",
				},
			},

			Subnets: []*Subnet{
				&Subnet{
					Id:   "1",
					Zone: "us-west1-a",
				},
				&Subnet{
					Id:   "2",
					Zone: "us-west1-a",
				},
				&Subnet{
					Id:   "3",
					Zone: "us-west1-b",
				},
			},

			Expected: []Expectation{
				Expectation{
					NatInstances: []string{
						"1-us-west1-a", "2-us-west1-b", "3-us-west1-c",
					},

					AssociateSubnets: map[string][]string{
						"1-us-west1-a": []string{"1-us-west1-a", "2-us-west1-a"},
						"2-us-west1-b": []string{"3-us-west1-b"},
						"3-us-west1-c": []string{},
					},
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
					Zone: "us-west1-c",
				},
				&Subnet{
					Id:   "4",
					Zone: "us-west1-c",
				},
			},

			Expected: []Expectation{
				Expectation{
					NatInstances: []string{
						"1-us-west1-a", "2-us-west1-b",
					},

					AssociateSubnets: map[string][]string{
						"1-us-west1-a": []string{"1-us-west1-a", "3-us-west1-c"},
						"2-us-west1-b": []string{"2-us-west1-b", "4-us-west-c"},
					},
				},
			},
		},
	}

	for index, test := range testcases {
		instancesInZone := make(map[string][]*NATInstance)
		instancesInZone = mapHealthyInstancestoZone(test.HealthyInstances)
		allocate(instancesInZone, test.Subnets)

		result := Expectation{}
		totalSubnet := 0
		longestAssociateSubnet := 0
		result.AssociateSubnets = map[string][]string{}
		for _, instances := range instancesInZone {
			for _, i := range instances {
				count := 0
				instanceName := fmt.Sprintf("%v-%v", i.Id, i.Zone)
				result.NatInstances = append(result.NatInstances, instanceName)
				for _, s := range i.Subnets {
					count += 1
					totalSubnet += 1
					subnet := fmt.Sprintf("%v-%v", s.Id, s.Zone)
					result.AssociateSubnets[instanceName] = append(result.AssociateSubnets[instanceName], subnet)
				}
				if count > longestAssociateSubnet {
					longestAssociateSubnet = count
				}
			}
		}

		fmt.Println("-> Run Test", index)
		// Check total associated subnets
		require.Equal(t, len(test.Subnets), totalSubnet)
		for index, expect := range test.Expected {
			// Check number of all NAT instances
			require.Equal(t, len(expect.NatInstances), len(result.NatInstances))

			for _, instance := range test.Expected[index].NatInstances {
				require.GreaterOrEqual(t, longestAssociateSubnet, len(expect.AssociateSubnets[instance]))
			}
		}

	}
}
