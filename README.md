# eamold
This server was created as a research platform for the old XML-based only (no binary XML) e-Amusement networking used on System 573.

The code is provided as-is and no support will be provided. The server is provided under a permissive license so you can freely fork and modify it for your needs, or you can reference the network requests and responses to re-implement in your own server.

Performance and security was not a focus during development so it is not recommended to use this server in a public environment.

## Running
Modify `config.yml` as required for your setup.

Run `go run main.go` from the repository folder to run from source code.

Alternatively, if you are using a pre-built executable from the Releases tab, then run the eamold-(platform)-(arch) executable directly.

## Database management

[Atlas](https://atlasgo.io/) is used for managing the database. Follow [the setup instructions](https://atlasgo.io/getting-started) to get Atlas installed. You can use the [update_db.sh](update_db.sh) script to automatically gather all schema.sql files and create and update the SQLite3 database schema as needed. If you modify any of the table schemas then run this command again to apply the changes to your SQLite3 database.

If you wish to use a database other than SQLite then you must modify update_db.sh and sqlc.yaml accordingly, as well as modify the SQL queries according to your database of choice.

How to execute update_db.sh:

`./update_db.sh server.db`

After modifying the SQL schema or queries, you must run `sqlc generate` to regenerate all of the database-related Go code.

## Support
- Drummania 7th Mix
- Drummania 7th Mix Power-Up Version
- Drummania 8th Mix
- Drummania 9th Mix
- Drummania 10th Mix
- Guitar Freaks 8th Mix
- Guitar Freaks 8th Mix Power-Up Version
- Guitar Freaks 9th Mix
- Guitar Freaks 10th Mix
- Guitar Freaks 11th Mix

Not all features of all games are implemented, but the basics for all games have been implemented. Events have the bare minimum implemented enough that all related songs should be unlocked. Not every combination of play (solo, 2 player, 1 player on left side or right side only, course modes, session modes, etc) have been tested so there's very likely to be bugs.

NOTE: Skill point-based song unlocks are saved on the player's magnetic game card and not on the server.

## Design

Main entry point of the server is `main.go`. All services must be registered using `manager.RegisterService()` here to be exposed through the HTTP server.

The `services` folder is where the main logic for the server lives. Each service has a `service.go` file which implements the `Service` interface that allows it to be registered as a service in the service manager. Within that is the list of modules that each service actually exposes through `services.get`. There's a way to resolve common shared modules from the core service for the common things like `pcbtracker`, `posevent`, etc. Each service thne implements any additional modules within its own service folder, including any database-related schemas and queries.

In the case where a module is shared between multiple games in a game series, you can split those modules out into a common package and make them accept a data provider that implements a given interface so that all individual game services can implement that interface while providing data from their own databases. An example of this can be seen in `services/gfdm_common/modules/demomusic.go` where all games that want to provide the demomusic module must also provide something that implements the DemoMusicDataProvider interface containing the GetDemoMusic function.

You can find the request and response structs in the models folder of each service. The convention used is `modulename_methodname.go`, and the struct names are `Request_ModuleName_MethodName` and `Response_ModuleName_MethodName` for the top element and an additional `_MemberName` appended to the current struct name for each member's custom struct (such as `Response_GameData_DataGet_ShopRank_Shop` which means gamedata's dataget response contains an element named "shoprank" that defines a new struct, which itself contains an element named "shop" which has its own custom struct).

The actual XML element name can be checked by looking at the XML tag. For example,
```
Prize    Response_GameData_DataGet_Prize    `xml:"gohohbi"`
```
inside the XML response the element name will be `gohohbi` but is referenced by the name `Prize` in the Go code.


## ee'mall GFDM content

The source files required to make the update data cannot be shared here due to legal reasons, but the data can be pulled from any GF11/DM10 or GF10/DM9 HDD that has ever received an ee'mall update (if there's a net_id.bin in a Dxx folder in the top directory, then it has received an ee'mall update at some point). `data_ver.bin` is also found in all of the song data folders on the HDD but the game will generate that on download so it's left out of the folders here. See [updates_work/hashes.txt](updates_work/hashes.txt) for a list of hashes for all data.

Data from the HDD images can be extracted using [this tool](https://github.com/987123879113/gobbletools/blob/master/sys573/tools/dump_pythonfs_eamuse_hdd.py).

Once the data is in the appropriate places, you can use the [tools/generate_updates.py](tools/generate_updates.py) script to generate the actual encrypted data chunks and XML files that the game will download. This tool will also spit out a static.json file containing a list of URL paths and the corresponding local data path to be exposed by the server as a static file (see config.yml's static configuration section).

An additional script [tools/generate_filelist_bin.py](tools/generate_filelist_bin.py) has also been provided which allows you to easily create a filelist.bin file for a given folder in case you modify an existing song's data or if you wish to create your own custom song data.
