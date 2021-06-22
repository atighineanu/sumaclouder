## sumaclouder - a tool for automating the work with Google Cloud Platform

### What it does:
 - uploads the public cloud images into the bucket corresponding to your team
 - updates these images (from the bucket) to the latest version


### How to:
if not built:
  `go run main.go -h`  - for help

  `go run main.go imgupdate` - will update the public cloud images inside qa-css bucket (check the root.go file -> function init(), 
  the default value of --bucketname flag is "suse-manager-images"; projectID is - "suse-css-qa" in the same file)

  `go run main.go listimages --bucketname <some-bucket>` - will list all the files in a given bucket (default bucket is the same
  like the one specified above)

if you built the binary:
`./sumaclouder <feature>`


#### To make it work:
- you have to have go installed, GOPATH and GOROOT properly set
- you have to have VPN to our intranet running
- you have to have the intranet DNS/resolv.conf properly set
