





## *RunConfig*
Root of the config file


### Fields

<dl>
<dt>
	<code>name</code>  <strong>string</strong>  - required
</dt>

<dd>

Name of the upgrade.


<blockquote>It is recommended to use an upgrade version as a name.</blockquote>
</dd>

<dt>
	<github.com/zeta</code>  <strong><a href="#zetaconfig">ZetaConfig</a></strong>  - required
</dt>

<dd>

Configuration of a Zeta node.

</dd>

<dt>
	<code>data_node</code>  <strong><a href="#datanodeconfig">DataNodeConfig</a></strong>  - optional
</dt>

<dd>

Configuration of a Data node.

</dd>



### Complete example


```hcl
name = "v1.65.0"

[zeta]
 [zeta.binary]
  path = "/path/zeta-binary"
  args = ["--arg1", "val1", "--arg2"]
 [zeta.rpc]
  socketPath = "/path/socket.sock"
  httpPath = "/rpc"

```


</dl>

---


## *ZetaConfig*
Allows to configure Zeta binary and it's arguments.


### Fields

<dl>
<dt>
	<code>binary</code>  <strong><a href="#binaryconfig">BinaryConfig</a></strong>  - required
</dt>

<dd>

Configuration of Zeta binary to be run.

</dd>

<dt>
	<code>rpc</code>  <strong><a href="#rpcconfig">RPCConfig</a></strong>  - required
</dt>

<dd>

Visor communicates with the core node via RPC API that runs over UNIX socket.
This parameter allows to configure the UNIX socket to match the core node configuration.


</dd>



### Complete example


```hcl
[zeta]
 [zeta.binary]
  path = "/path/zeta-binary"
  args = ["--arg1", "val1", "--arg2"]
 [zeta.rpc]
  socketPath = "/path/socket.sock"
  httpPath = "/rpc"

```


</dl>

---


## *DataNodeConfig*
Allows to configure Data node binary and it's arguments.


### Fields

<dl>
<dt>
	<code>binary</code>  <strong><a href="#binaryconfig">BinaryConfig</a></strong>  - required
</dt>

<dd>



</dd>



### Complete example


```hcl
[data_node]
 [data_node.binary]
  path = "/path/data-node-binary"
  args = ["--arg1", "val1", "--arg2"]

```


</dl>

---


## *BinaryConfig*
Allows to configure binary and it's arguments.


### Fields

<dl>
<dt>
	<code>path</code>  <strong>string</strong>  - required
</dt>

<dd>

Path to the binary.


<blockquote>Both absolute or relative path can be used.
Relative path is relative to a parent folder of this config file.
</blockquote>
</dd>

<dt>
	<code>args</code>  <strong>[]string</strong>  - required
</dt>

<dd>

Arguments that will be applied to the binary.


<blockquote>Each element the list represents one space seperated argument.
</blockquote>
</dd>



### Complete example


```hcl
path = "/path/binary"
args = ["--arg1", "val1", "--arg2"]

```


</dl>

---


## *RPCConfig*
Allows to configure connection to core node exposed UNIX socket RPC API.


### Fields

<dl>
<dt>
	<code>socketPath</code>  <strong>string</strong>  - required
</dt>

<dd>

Path of the mounted socket.


<blockquote>This path can be configured in Zeta core node configuration.</blockquote>
</dd>

<dt>
	<code>httpPath</code>  <strong>string</strong>  - required
</dt>

<dd>

HTTP path of the socket path.


<blockquote>This path can be configured in Zeta core node configuration.</blockquote>
</dd>



### Complete example


```hcl
[zeta.rpc]
 socketPath = "/path/socket.sock"
 httpPath = "/rpc"

```


</dl>

---


