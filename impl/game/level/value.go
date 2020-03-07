package level

func chunkIndex(x, z int) int64 {
	return (int64(z) << 0x20) | (int64(x) & 0xFFFFFFFF)
}

func sliceIndex(x, y, z int) int {
	return y<<0x08 | z<<0x04 | x
}

func blockYToSliceY(blockY int) (sliceY int) {
	sliceY = blockY >> 0x04
	return
}

func blockXZToChunkXZ(blockX, blockZ int) (chunkX, chunkZ int) {
	chunkX = blockX >> 0x04
	chunkZ = blockZ >> 0x04
	return
}

func blockLevelToSlice(levelBlockX, levelBlockY, levelBlockZ int) (sliceBlockX, sliceBlockY, sliceBlockZ int) {
	sliceBlockX = levelBlockX & 0xF
	sliceBlockY = levelBlockY & 0xF
	sliceBlockZ = levelBlockZ & 0xF

	return
}
