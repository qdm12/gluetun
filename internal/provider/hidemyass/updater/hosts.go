package updater

func getUniqueHosts(tcpHostToURL, udpHostToURL map[string]string) (
	hosts []string) {
	uniqueHosts := make(map[string]struct{}, len(tcpHostToURL))
	for host := range tcpHostToURL {
		uniqueHosts[host] = struct{}{}
	}
	for host := range udpHostToURL {
		uniqueHosts[host] = struct{}{}
	}

	hosts = make([]string, 0, len(uniqueHosts))
	for host := range uniqueHosts {
		hosts = append(hosts, host)
	}

	return hosts
}
