# simple-_blockhain_implementation
In this application we have developed a blockchain application with basic blockchain functions
there is a blockchain that is centralized in this scenario, and when the center node starts to work, it creates the genesis block. 
After that, it starts to listen in 3000 ports to respond to incoming requests.
When the miner starts to work, the miner node sends a request to join the network.
If the central node responds positively, the miner joins the network.
Then the miner node gets all the blocks from the Central node and waits to receive transaction list.
When we run the Send_TransactionList the transaction list goes to the central node.
The central node sends the incoming transactions to the miner node.
Miner node finds the new block with proof of work and sends the center node.
The central node checks the correctness of the block, if it is true it appends it to its database and sends it to the miner nodes.




