package main

// Play a sound while soundTimer lasts
func playSound() {
	if soundTimer > 0 {
		playFreq(440)
	} else {
		stopFreq()
	}
}

// Play a specfic frequency
func playFreq(freq int) {

}

// Stop playing a frequency
func stopFreq() {

}
