// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package constants

import (
	"time"
)

// Messages
const (
	notused byte = iota
	EOM_MSG
	ACK_MSG
	FED_SERVER_FAULT_MSG
	AUDIT_SERVER_FAULT_MSG
	FULL_SERVER_FAULT_MSG
	COMMIT_CHAIN_MSG
	COMMIT_ENTRY_MSG
	DIRECTORY_BLOCK_SIGNATURE_MSG
	EOM_TIMEOUT_MSG
	FACTOID_TRANSACTION_MSG
	HEARTBEAT_MSG
	INVALID_ACK_MSG
	INVALID_DIRECTORY_BLOCK_MSG

	REVEAL_ENTRY_MSG
	REQUEST_BLOCK_MSG
	SIGNATURE_TIMEOUT_MSG
	MISSING_MSG
	MISSING_DATA
	DATA_RESPONSE
	MISSING_MSG_RESPONSE

	DBSTATE_MSG
	DBSTATE_MISSING_MSG
	ADDSERVER_MSG
	CHANGESERVER_KEY_MSG
	REMOVESERVER_MSG

	BOUNCE_MSG      //	test message
	BOUNCEREPLY_MSG // 	test message

	MISSING_ENTRY_BLOCKS
	ENTRY_BLOCK_RESPONSE

	INTERNALADDLEADER
	INTERNALREMOVELEADER
	INTERNALADDAUDIT
	INTERNALREMOVEAUDIT
	INTERNALTIMEOUT
	INTERNALSIG

	VOLUNTEERAUDIT
	SYNC_MSG
	LEADER_ACK_VOLUNTEER

	NUM_MESSAGES // Not used, just a counter for the number of messages.
)

const (
	// Limits for keeping inputs from flooding our execution
	INMSGQUEUE_HIGH = 1000
	INMSGQUEUE_MED  = 500
	INMSGQUEUE_LOW  = 100

	DBSTATE_REQUEST_LIM_HIGH = 200
	DBSTATE_REQUEST_LIM_MED  = 50

	// Replay -- Dynamic Replay filter based on messages as they are processed.
	INTERNAL_REPLAY = 1
	NETWORK_REPLAY  = 2
	TIME_TEST       = 4 // Checks the time_stamp;  Don't put actual hashes into the map with this.
	REVEAL_REPLAY   = 8 // Checks for Reveal Entry Replays ... No duplicate Entries within our 4 hours!

	// FReplay -- Block based Replay filter consttructed by processing the blocks, from the database
	//            then from blocks either passed to a node, or constructed by messages.
	BLOCK_REPLAY = 1 // Ensures we don't add the same transaction to multiple blocks.

	ADDRESS_LENGTH = 32 // Length of an Address or a Hash or Public Key
	// length of a Private Key
	SIGNATURE_LENGTH     = 64    // Length of a signature
	MAX_TRANSACTION_SIZE = 10240 // 10K like everything else?
	// Not sure if we need a minimum amount.  Set at 1 Factoshi

	// Database
	//==================
	// Limit on size of keys, since Maps in Go can't handle variable length keys.

	// Wallet
	//==================
	// Holds the root seeds for address generation
	// Holds the latest generated seed for each root seed.

	// Block
	//==================
	MARKER                  = 0x00                       // Byte used to mark minute boundries in Factoid blocks
	TRANSACTION_PRIOR_LIMIT = int64(12 * 60 * 60 * 1000) // Transactions prior to 12hrs before a block are invalid
	TRANSACTION_POST_LIMIT  = int64(12 * 60 * 60 * 1000) // Transactions after 12hrs following a block are invalid

	//Entry Credit Blocks (For now, everyone gets the same cap)
	EC_CAP = 5 //Number of ECBlocks we start with.
	//Administrative Block Cap for AB messages

	//Limits and Sizes
	//==================
	//Maximum size for Entry External IDs and the Data
	HASH_LENGTH = int(32) //Length of a Hash
	//Length of a signature
	//Prphan mem pool size
	//Transaction mem pool size
	//Block mem bool size
	//MY Process List size

	//Max number of entry credits per entry
	//Max number of entry credits per chain

	COMMIT_TIME_WINDOW = time.Duration(12) //Time windows for commit chain and commit entry +/- 12 hours

	//NETWORK constants
	//==================
	VERSION_0               = byte(0)
	MAIN_NETWORK_ID  uint32 = 0xFA92E5A2
	TEST_NETWORK_ID  uint32 = 0xFA92E5A3
	LOCAL_NETWORK_ID uint32 = 0xFA92E5A4
	MaxBlocksPerMsg         = 500
)

const (
	// NETWORKS:
	NETWORK_MAIN   int = iota // 0
	NETWORK_TEST              // 1
	NETWORK_LOCAL             // 2
	NETWORK_CUSTOM            // 3
)

// Slices and arrays that should not ever be modified:
//===================================================
// Used as a key in the wallet to find the current seed value.
var CURRENT_SEED = [32]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

// Entry Credit Chain
var EC_CHAINID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0c}

// Directory Chain
var D_CHAINID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0d}

// Directory Chain
var ADMIN_CHAINID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0a}

// Factoid chain
var FACTOID_CHAINID = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0f}

// Zero Hash
var ZERO_HASH = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var ZERO = []byte{0}

//---------------------------------------------------------------
// Types of entries (transactions) for Admin Block
// https://github.com/FactomProject/FactomDocs/blob/master/factomDataStructureDetails.md#adminid-bytes
//---------------------------------------------------------------
const (
	TYPE_MINUTE_NUM         uint8 = iota // 0
	TYPE_DB_SIGNATURE                    // 1
	TYPE_REVEAL_MATRYOSHKA               // 2
	TYPE_ADD_MATRYOSHKA                  // 3
	TYPE_ADD_SERVER_COUNT                // 4
	TYPE_ADD_FED_SERVER                  // 5
	TYPE_ADD_AUDIT_SERVER                // 6
	TYPE_REMOVE_FED_SERVER               // 7
	TYPE_ADD_FED_SERVER_KEY              // 8
	TYPE_ADD_BTC_ANCHOR_KEY              // 9
	TYPE_SERVER_FAULT
)

//---------------------------------------------------------------------
// Identity Status Types
//---------------------------------------------------------------------
const (
	IDENTITY_UNASSIGNED               uint8 = iota // 0
	IDENTITY_FEDERATED_SERVER                      // 1
	IDENTITY_AUDIT_SERVER                          // 2
	IDENTITY_FULL                                  // 3
	IDENTITY_PENDING_FEDERATED_SERVER              // 4
	IDENTITY_PENDING_AUDIT_SERVER                  // 5
	IDENTITY_PENDING_FULL                          // 6
	IDENTITY_SKELETON                              // 7 - Skeleton Identity
)
