package node

type TRouteTable struct {
	NodeList []TNodeInfo `json:"node_list"`
}

func (rt *TRouteTable) New() *TRouteTable {
	return nil
}
