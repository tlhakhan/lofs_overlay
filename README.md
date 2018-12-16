# lofs_overlay
  - Use case:  When you want to overlay custom config files (post-OS init) from a sourced directory over / filesystem, especially for Live OS where most files are on read-only RAM disk.
  - Requires: `SunOS` and was tested/working on `SunOS smos-02 5.11 joyent_20181206T012147Z i86pc i386 i86pc`
  - Dependencies: Uses filesystem `lofs` to `mount -F` on the system.

# Usage
  - Expects a string argument, either `start` or `stop`.
  - By default, assumes directory is `/usbkey/crud` as the sourced directory that will overlay ontop of / filesystem.
  - If a different directory is desired, use the `-overlayRootPath`.

```
[root@smos-02 (home) /opt/lofs_overlay]# bin/SunOS/lofs_overlay -h
Usage of bin/SunOS/lofs_overlay:
  -overlayRootPath string
        The folder path that will be considered as / when overlaying (default "/usbkey/crud")

[root@smos-02 (home) /opt/lofs_overlay]# bin/SunOS/lofs_overlay
Usage: bin/SunOS/lofs_overlay start|stop
```

# Output: start overlay
  - Start the overlay but `/usbkey/crud` doesn't exist.

```
[root@smos-02 (home) /opt/lofs_overlay]# bin/SunOS/lofs_overlay start
start the overlay process
/usbkey/crud does not exist, exiting..
```

# Output: start overlay
  - Starting the overlay.

```
[root@smos-02 (home) /opt/lofs_overlay]# bin/SunOS/lofs_overlay start
start the overlay process
walking dir /usbkey/crud
found directory /usbkey/crud/root
target path /root exists
target path /root is a directory
mkdir is not needed on /root
found file /usbkey/crud/root/.profile
target path /root/.profile exists
target path /root/.profile is a file
create file is not needed on /root/.profile
/usbkey/crud/root/.profile is now mounted to target file /root/.profile
found file /usbkey/crud/root/.vimrc
target path /root/.vimrc exists
target path /root/.vimrc is a file
create file is not needed on /root/.vimrc
/usbkey/crud/root/.vimrc is now mounted to target file /root/.vimrc
done linking the overlay
```

# Output: start overlay
  - Starting the overlay on a different overlay root path (/usbkey/crud2).

```
[root@smos-02 (home) /opt/lofs_overlay]# bin/SunOS/lofs_overlay -overlayRootPath /usbkey/crud2 start
start the overlay process
walking dir /usbkey/crud2
found directory /usbkey/crud2/root
target path /root exists
target path /root is a directory
mkdir is not needed on /root
found file /usbkey/crud2/root/.profile
target path /root/.profile exists
target path /root/.profile is a file
create file is not needed on /root/.profile
/usbkey/crud2/root/.profile is now mounted to target file /root/.profile
found file /usbkey/crud2/root/.vimrc
target path /root/.vimrc exists
target path /root/.vimrc is a file
create file is not needed on /root/.vimrc
/usbkey/crud2/root/.vimrc is now mounted to target file /root/.vimrc
done linking the overlay
```

# Output: stop overlay
  - Stopping the overlay on a running system.

```
[root@smos-02 (home) /opt/lofs_overlay]# bin/SunOS/lofs_overlay stop
stop the overlay process
unmounting the found target files: ["/root/.profile" "/root/.vimrc"]
target file /root/.profile is now unmounted
target file /root/.vimrc is now unmounted
done unlinking the overlay
```

# Appendix: Create an lofs_overlay service
  - Place the `lofs_overlay` binary with execute permissions into `/opt/custom/smf/share` directory.
  - Create the below SMF xml file in `/opt/custom/smf` directory, preferrably as `lofs_overlay.xml`.
  - The service starts immediately upon import and will start at post-OS load.
    - Use `svccfg import <service>.xml` to import the service on a running system.

```xml
<?xml version='1.0'?>
<!DOCTYPE service_bundle SYSTEM '/usr/share/lib/xml/dtd/service_bundle.dtd.1'>
<service_bundle type='manifest' name='export'>
  <service name='site/lofs_overlay' type='service' version='0'>
    <create_default_instance enabled='true'/>
    <single_instance/>
    <dependency name='filesystem' grouping='require_all' restart_on='error' type='service'>
      <service_fmri value='svc:/system/filesystem/local'/>
    </dependency>
    <method_context/>
    <exec_method name='start' type='method' exec='/opt/custom/smf/share/lofs_overlay start' timeout_seconds='60'/>
    <exec_method name='stop' type='method' exec='/opt/custom/smf/share/lofs_overlay stop' timeout_seconds='60'/>
    <property_group name='startd' type='framework'>
      <propval name='duration' type='astring' value='transient'/>
      <propval name='ignore_error' type='astring' value='core,signal'/>
    </property_group>
    <property_group name='application' type='application'/>
    <stability value='Evolving'/>
    <template>
      <common_name>
        <loctext xml:lang='C'>lofs overlay on / filesystem using /usbkey/crud as source directory</loctext>
      </common_name>
    </template>
  </service>
</service_bundle>
```
