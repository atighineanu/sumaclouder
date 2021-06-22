## sumaclouder - a tool for automating the work with Google Cloud Platform

### What it does:
 - uploads the public cloud images into the bucket corresponding to your team
 - updates these images (from the bucket) to the latest version


### How to:
if not built:
  `go run main.go -h`  - for help
  
  `go run main.go imgupdate` - will update the public cloud images inside qa-css bucket (check the root.go file -> function init(), 
  the default value of --bucketname flag is "suse-manager-images"; projectID is - "suse-css-qa" in the same file)

