FLAT
===

Dedicated to all those who believe that the world is a disc supported on the back of four large elephants that are standing on the back of a huge tortoise.

Function
--------

FLAT takes a solidity file and retrieves all its dependancies and stuffs them in a single file in the correct order for it to be verified on etherscan.

WHY?
-----
Because part of the ethos of the smart contract blockchain relies upon being able to see and verify the contracts yourself. Etherscan require single files for verification. FLAT helps you create those files.

Also - because my MacBook _hates_ **Python** and refuses to figure out how to run the version that [BLOCKCAT](https://github.com/BlockCatIO/solidity-flattener) put together.

And because it can be compiled into executables for any platform.

Installation
---

This is not intended as a library (no go getting required). 

Just clone it and build it after importing the single dependancy which is not included here : 

`> go get gopkg.in/natefinch/lumberjack.v2`

You can build it yourself or download an executable (coming soon)

built using **go v 1.10** but it would probably work with an older version if you are feeling decadent.

`> go build flat.go utils.go`

or 

`> go install  flat.go utils.go`

Usage
---

Assuming that you have the executable on your path :

Assuming your contract is in `mainfile.sol` and you want to create `consolidated.sol`

`flat -input mainfile.sol -output consolidated`

This creates a flattened version of mainfile.sol with all includes in the file `consolidated.sol` and creates a log called `consolidated.log`

_NOTE_ the output file **must not exist**. This is to prevent you from overwriting that file that you forgot to check into git after a tough night's coding.

---

Bugs / enhancements : please raise an issue or fork and issue a PR.



Dave Appleton.

* Lead Blockchain dev @ [HelloGold](https://hellogold.com)
* Senior Consultant @ [Akomba Labs](https://akombalabs.com)
* [@AppletonDave](https://twitter.com/AppletonDave) on twitter.