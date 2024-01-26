package main

import (
	"fmt"
)

func main() {
	healthyInstances := []*NATInstance{
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
	}

	subnets := []*Subnet{
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
	}

	instancesInZone := make(map[string][]*NATInstance)
	instancesInZone = mapHealthyInstancestoZone(healthyInstances)

	allocate(instancesInZone, subnets)
	printInstances(instancesInZone)
}

type Subnet struct {
	Id   string
	Zone string
}

type NATInstance struct {
	Id      string
	Zone    string
	Subnets []*Subnet
}

func printInstances(instancesInZone map[string][]*NATInstance) {
	for _, instances := range instancesInZone {
		for _, i := range instances {
			fmt.Printf("Instance (%v-%v):\n", i.Id, i.Zone)
			for _, s := range i.Subnets {
				fmt.Printf("\tsubnet (%v-%v)\n", s.Id, s.Zone)
			}
		}
	}
}

func mapHealthyInstancestoZone(healthyInstances []*NATInstance) map[string][]*NATInstance {
	instancesInZone := make(map[string][]*NATInstance)

	for _, i := range healthyInstances {
		instancesInZone[i.Zone] = append(instancesInZone[i.Zone], i)
	}

	return instancesInZone
}

// allocate Subnets to Instances
func allocate(instancesInZone map[string][]*NATInstance, subnets []*Subnet) {
	for _, s := range subnets {
		natInstances := instancesInZone[s.Zone]
		if len(natInstances) != 0 {
			// fmt.Println(s.Zone, " ", len(instancesInZone[s.Zone]))
			// Case 1: There are some healthy NAT Instances in the same AZ

			weight := len(natInstances[0].Subnets)
			validInstance := natInstances[0]
			for _, instance := range natInstances {
				if len(instance.Subnets) < weight {
					weight = len(instance.Subnets)
					validInstance = instance
				}
			}

			validInstance.Subnets = append(validInstance.Subnets, s)
		} else {
			// fmt.Println("There is no NAT instances in ", s.Zone)
			// Case 2: there are no healthy NAT Instances in the same AZ

			weight := -1
			validInstance := &NATInstance{}

			for _, instances := range instancesInZone {
				for _, instance := range instances {
					if weight == -1 {
						weight = len(instance.Subnets)
						validInstance = instance
					}

					if len(instance.Subnets) < weight {
						weight = len(instance.Subnets)
						validInstance = instance
					}
				}
			}
			validInstance.Subnets = append(validInstance.Subnets, s)
		}
	}

}
