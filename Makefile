dev:
	@air serve

dpush:
	@sudo docker buildx build  --platform linux/amd64 -t kirari04/convertarr:latest --sbom=true --provenance=true --push .
	@sudo docker buildx build  --platform linux/amd64 -t kirari04/convertarr:amd -f Dockerfile.amd --sbom=true --provenance=true --push .