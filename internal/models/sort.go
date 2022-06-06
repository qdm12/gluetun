package models

import "sort"

var _ sort.Interface = (*SortableServers)(nil)

type SortableServers []Server

func (s SortableServers) Len() int {
	return len(s)
}

func (s SortableServers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableServers) Less(i, j int) bool {
	a, b := s[i], s[j]

	if a.Country == b.Country { //nolint:nestif
		if a.Region == b.Region {
			if a.City == b.City {
				if a.ServerName == b.ServerName {
					if a.Number == b.Number {
						if a.Hostname == b.Hostname {
							if a.ISP == b.ISP {
								return a.VPN < b.VPN
							}
							return a.ISP < b.ISP
						}
						return a.Hostname < b.Hostname
					}
					return a.Number < b.Number
				}
				return a.ServerName < b.ServerName
			}
			return a.City < b.City
		}
		return a.Region < b.Region
	}
	return a.Country < b.Country
}
