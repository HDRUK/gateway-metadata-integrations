## Metadata Federation Tiltfile
##
## Loki Sinclair <loki.sinclair@hdruk.ac.uk>
##

cfg = read_json('tiltconf.json')

docker_build(
    ref='hdruk/' + cfg.get('name'),
    context='.',
    live_update=[
        sync('.', '/app'),
        run('go mod download', trigger='./go.mod'),
        run('go build --ldflags="-s -w" -o metadata_federation_service'),
    ]
)

k8s_yaml('chart/' + cfg.get('name') + '/deployment.yaml')
k8s_yaml('chart/' + cfg.get('name') + '/service.yaml')
k8s_resource(
    cfg.get('name'),
    port_forwards=9889,
    labels=["API"]
)