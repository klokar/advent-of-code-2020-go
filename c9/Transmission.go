package c9

type Transmission struct {
	Data    []int
	pointer int
}

type PreambleTransmission struct {
	Transmission
	Length       int
	preamble     []int
	calculations []int
}

type ContiguousTransmission struct {
	Transmission
	SearchedValue int
	offset        int
}

func (st *Transmission) ReadNext() bool {
	if st.pointer+1 < len(st.Data) {
		st.pointer++
		return true
	}

	return false
}

func (st Transmission) Current() int {
	return st.Data[st.pointer]
}

func (pt *PreambleTransmission) ReadNext() bool {
	if !pt.Transmission.ReadNext() {
		return false
	}

	if len(pt.preamble) == 0 {
		pt.pointer = pt.Length - 1
		pt.updatePreamble(pt.Length)
	} else {
		pt.updatePreamble(1)
	}

	return true
}

func (pt *PreambleTransmission) updatePreamble(length int) {
	// Set new preamble
	pt.preamble = pt.Data[pt.pointer-pt.Length+1 : pt.pointer+1]

	// Remove old calculations from the front that represents deprecated number
	// or specify empty calculations if starting
	newCalculationStartIndex := length*pt.Length - 1
	if len(pt.calculations) > newCalculationStartIndex {
		pt.calculations = pt.calculations[newCalculationStartIndex:]
	} else {
		pt.calculations = make([]int, 0, pt.Length*pt.Length)
	}

	// Add calculations for current line at the end
	for i := pt.Length - length; i < pt.Length; i++ {
		for _, compared := range pt.preamble {
			pt.calculations = append(pt.calculations, compared+pt.preamble[i])
		}
	}
}

func (pt PreambleTransmission) IsValid() bool {
	matches := false
	for _, calculation := range pt.calculations {
		if calculation == pt.preamble[pt.Length-1] {
			matches = true
		}
	}

	return matches
}

func (ct *ContiguousTransmission) Find() (bool, []int, int, int) {
	ct.offset = 1
	sum := ct.Current()
	min := ct.Current()
	max := ct.Current()
	for ct.offset = 1; sum < ct.SearchedValue; ct.offset++ {
		v := ct.Data[ct.pointer+ct.offset]

		if v < min {
			min = v
		}

		if v > max {
			max = v
		}

		sum += v
	}

	if sum == ct.SearchedValue {
		return true, ct.Data[ct.pointer : ct.pointer+ct.offset], min, max
	}

	return false, []int{}, 0, 0
}
