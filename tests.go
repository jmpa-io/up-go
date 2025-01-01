package up

// This file contains common variables and functions used across tests.

import "time"

// A location used for tests.
var location = time.FixedZone("AEST", 11*60*60)
