```bash

# make area for scripts/binary that the service manifest might call
mkdir -p /opt/custom/smf/share

# place the lofs_overlay.xml file here
touch /opt/custom/smf/lofs_overlay.xml

# place the lofs_overlay binary here
touch /opt/custom/smf/share/lofs_overlay

# import the smf - it will immediately start
svccfg import /opt/custom/smf/lofs_overlay.xml

```
