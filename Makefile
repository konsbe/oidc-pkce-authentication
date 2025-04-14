build-img:
	docker run -it -dp 8080:8080 -e KC_BOOTSTRAP_ADMIN_USERNAME=admin -e KC_BOOTSTRAP_ADMIN_PASSWORD=admin --network host quay.io/keycloak/keycloak:26.1.4 start-dev

tag-img: 
	sudo docker tag c8bfd0bc98f0 backend-oidc-pkce-keycloak-image:keycloak-server
