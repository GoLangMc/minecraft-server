package mask

type Masking struct{}

func (m *Masking) Has(mask, field byte) bool {
	return mask&field != 0
}

func (m *Masking) Set(mask *byte, field byte, when bool) {
	if when {
		*mask |= field
	}
}
