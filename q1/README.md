BONUS: What if each Subnet has a `Weight int32` attribute and we try to make total weight allocated to each NAT Instance the same no matter how subnets allocated to each NAT Instance?

ANSWER:
 - Advantages: This application functions like a Load Balancer.
 - Drawbacks: Incur more inter-AZ data transfer charges e.g Client traffic from the internet, Cross-VPC with VPC peering, across AZs. ref [https://qrlive.2am-media.tech/q1_ref]