# Handin-4

## Contents

- [Contents](#Contents)
- [Setup](#Setup)

## Setup

There are 4 pre-defined ports.

All the ports are written in the file "$ports.txt$"
If you for some reason would like to alter the ports, then you can change them in the file.

## Run the program

Download the zip file.

Ensure you have go downloaded.

Navigate into the folder, and run the following: 
``` bash 
go run . $$ 
```

After this it will display that you are joining on a specific port with a certain node id and at as certain timestamp.

You will then be showed the following prompt:
``` bash
 Select action: 
  ▸ Request
    Exit
```
Here you can use your arrow keys and enter to confirm selection.

It will then show on all others if you select Request.

When you have acessed the Critacal Secetion,
You will then be showed the following prompt:
``` bash
 Do you want to quit?: 
  ▸ Yes
    No
```
Selecting "No": The program releases the Critical Section, allowing other peers to access it, and the system continues running normally.

Selecting "Yes": The connection is terminated, and the system stops.




