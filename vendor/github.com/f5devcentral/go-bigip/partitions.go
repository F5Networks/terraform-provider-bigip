package bigip

// TMPartitions contains a list of all partitions on the BIG-IP system.
type TMPartitions struct {
        TMPartitions []*TMPartition `json:"items"`
}

type TMPartition struct {
        Name               string `json:"name,omitempty"`
        Kind               string `json:"kind,omitempty"`
        DefaultRouteDomain int    `json:"defaultRouteDomain,omitempty"`
        FullPath           string `json:"fullPath,omitempty"`
        SelfLink           string `json:"selfLink,omitempty"`
}

// TMPartitions returns a list of partitions.
func (b *BigIP) TMPartitions() (*TMPartitions, error) {
        var pList TMPartitions
        if err, _ := b.getForEntity(&pList, "auth", "tmPartition"); err != nil {
                return nil, err
        }
        return &pList, nil
}
