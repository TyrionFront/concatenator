**Concatenator**

In main.go:

-   "PoultateStogare" creates subfolders inside of the "storage" folder and puts .txt files inside of these subfolders according to next params that are taken via CLI prompt:

    -   `dirsCount` - number of subfolders; also is used as basis for random `int` multiplier together with _minimalDelayMs_ to calculate total waiting time for files creation;
    -   `filesCount` - number of files inside of every subfolder;
    -   `minimalDelayMs` - ticker (timer) delay that defines further total duration time of the whole creation process;

-   "WaitAndProces" laucnhes ticker (timer) related processing of the whole storage in parallel to "PoultateStogare", based on checking of the `done.txt` files existence within every subfolder and inside of the main storage respectively.

As soon as the storage is populated according to the provided params and `done.txt` files are created, content from every file in storage is read and concatenated into the final resulting file (`result.txt`) that is saved inside of the main folder of the project after the `done.txt` file is created in the "storage" folder.
