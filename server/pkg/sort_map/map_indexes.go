package sort_map

type MapIndexes interface {
	Insert(keyer Keyer)
	Remove(keyer Keyer)
	Range(func(idx int, keyer Keyer) bool)
}

type mapIndexes []Keyer

func (mi *mapIndexes) Insert(keyer Keyer) {
	idx, _ := bSearchFirstGT(*mi, keyer)
	mi.insert(idx, keyer)
}

func (mi *mapIndexes) Remove(keyer Keyer) {
	for i, k := range *mi {
		if k.Key() == keyer.Key() {
			*mi = append((*mi)[:i], (*mi)[i+1:]...)
			return
		}
	}
}

func (mi *mapIndexes) Range(f func(idx int, keyer Keyer) bool) {
	for i, keyer := range *mi {
		if !f(i, keyer) {
			return
		}
	}
}

func NewMapIndexes() MapIndexes {
	return &mapIndexes{}
}

func bSearchFirstGT(ints []Keyer, val Keyer) (int, Keyer) {
	return bSearchFirstGTInternally(ints, 0, len(ints)-1, val)
}

func bSearchFirstGTInternally(ints []Keyer, low int, high int, val Keyer) (int, Keyer) {
	if low > high {
		return -1, nil
	}

	if ints[low].GetOrder() > val.GetOrder() {
		return low, ints[low]
	}
	mid := low + ((high - low) >> 1)
	if ints[mid].GetOrder() > val.GetOrder() {
		return bSearchFirstGTInternally(ints, low, mid, val)
	} else {
		return bSearchFirstGTInternally(ints, mid+1, high, val)
	}
}

func (mi *mapIndexes) insert(index int, keyer Keyer) {
	n := len(*mi)
	*mi = append(*mi, keyer)
	if index == -1 {
		return
	}
	for i := n; i < index; i-- {
		(*mi)[i] = (*mi)[i-1]
	}
	(*mi)[index] = keyer
}
