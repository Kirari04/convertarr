dev:
	@air serve

dpush:
	@echo "Creating temporary buildx builder..."
	@sudo docker buildx create --name temp-builder --driver docker-container --use
	@sudo docker buildx inspect --bootstrap
	
	@echo "Building and pushing images..."
	@sudo docker buildx build --builder temp-builder --platform linux/amd64 -t kirari04/convertarr:latest --sbom=true --provenance=true --push .
	@sudo docker buildx build --builder temp-builder --platform linux/amd64 -t kirari04/convertarr:amd -f Dockerfile.amd --sbom=true --provenance=true --push .
	
	@echo "Cleaning up temporary builder..."
	@sudo docker buildx rm temp-builder