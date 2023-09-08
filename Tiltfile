## Metadata Federation Tiltfile
##
## Loki Sinclair <loki.sinclair@hdruk.ac.uk>
##

cfg = read_json('tiltconf.json')

docker_build(
    ref='hdruk/' + cfg.get('name'),
    context='.'
)

k8s_yaml('chart/' + cfg.get('name') + '/deployment.yaml')
k8s_yaml('chart/' + cfg.get('name') + '/service.yaml')
k8s_resource(
    cfg.get('name'),
    port_forwards=9889
)