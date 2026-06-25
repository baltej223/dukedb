# Duke JS Client

A lightweight JavaScript client for interacting with a Duke distributed key-value cluster.

## Installation

Copy `duke.js` into your project and import it:

```js
import Duke from "./duke.js";
```

## Usage

### Create a Client

```js
const duke = new Duke();
```

### Add Duke Nodes

```js
duke.seturl(
  "duke://localhost:8000",
  "duke://localhost:8001",
  "duke://localhost:8002",
);
```

Multiple URLs can be provided.

### Connect

Before performing any operations, call `connect()`.

```js
const connected = await duke.connect();

if (!connected) {
  console.log("No healthy Duke nodes found.");
}
```

The client checks the health of all configured nodes and stores only healthy ones for future requests.

---

## PUT

Store a key-value pair.

```js
await duke.PUT("name", "Baltej");
```

Returns:

```js
true;
```

---

## GET

Retrieve a value by key.

```js
const value = await duke.GET("name");

console.log(value);
```

Output:

```txt
Baltej
```

If the key does not exist:

```js
Error: Key not found.
```

---

## Health Checking

## The client verifies node health using db.checkHealth()

## Example

```js
import Duke from "./duke.js";

const duke = new Duke();

duke.seturl(
  "duke://localhost:8000",
  "duke://localhost:8001",
  "duke://localhost:8002",
);

const connected = await duke.connect();

if (!connected) {
  throw new Error("No Duke nodes available.");
}

await duke.PUT("username", "baltej");

const value = await duke.GET("username");

console.log(value);
```

---

## API

### `seturl(...urls)`

Registers Duke nodes.

```js
duke.seturl("duke://localhost:8000", "duke://localhost:8001");
```

---

### `connect()`

Checks node health and builds the list of working nodes.

Returns:

```js
boolean;
```

---

### `GET(key)`

Retrieves a value.

Returns:

```js
Promise<any>
```

Throws if the key is not found.

---

### `PUT(key, value)`

Stores a key-value pair.

Returns:

```js
Promise<boolean>
```

Throws if the operation fails.
