package up

// This file contains common variables and functions used across tests.

import "time"

// A location used for tests. +11:00 is AEDT (daylight saving); AEST is +10.
var location = time.FixedZone("AEDT", 11*60*60)
