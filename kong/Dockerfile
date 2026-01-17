FROM kong:latest

# Ensure any patching steps are executed as root user
USER root

# Add custom plugin to the image
#COPY ./kong/plugins/my-plugin /usr/local/share/lua/5.1/kong/plugins/my-plugin
#ENV KONG_PLUGINS=bundled,my-plugin

# Ensure kong user is selected for image execution
USER kong

# Expose Kong ports
EXPOSE 8000 8443 8001 8444
