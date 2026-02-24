# Hypha

Hypha (🍄) enrols new devices to an existing `cep-network` using `nebula.zip` payload quadruplets for now. 🍄 is wrapped in a basic ui running `fyne` for cross-portability. While we currently only support `Linux`, support for `Windows` and `Macos` will be coming next.

For now the user needs to supply their own `nebula.zip` archive and drop it in the below path:

`usr/home/.config/nebula-certs`

## Enrol

- To try the enrol service you should first drop your certs zip in the path below:

`usr/home/.config/nebula-certs`

- Then spin up the go app:

```
go run cmd/ui/ui.go
```

- Then click on the `Unzip Certificates` button

- If successful, you should be able to click on the `Start Nebula` and see the stdout piped from the nebula binaries
    - *NB:* Keep an eye on the terminal to be enter your pw for sudo elevation


For *Debugging* use the entry point in `app/internal directly`:

```
go run app/internal/enrolment/enrol.go
```


## DNS (WIP)

Cross platform DNS service, for now only works and tested on Linux. While some vibe-coded DNS resolvers were created already, I sure as hell don't trust them.
The UI is not integrated yet so for now you can only run in dev mode until we get the platform-specific builds done, here is how you can do it:

```
go run cmd/dns/dns-mapper.go
```
That's about it for now!

## Goal for V1.0.0:

- LAN-based retrieval of cep-standardized config (4 conf files)
- Load config to connect node into network
- All of this is wrapped by Hypha
- Set DNS service


Features:
 - Onboard device to network
 - Access app store
 - Collect Telemetry (for dist. use case)
 - Allocate mem/disk as a node

